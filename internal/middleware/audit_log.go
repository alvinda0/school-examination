package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/alvindashahrul/my-app/internal/services"
	"github.com/google/uuid"
)

type auditContextKey string

const auditDataKey auditContextKey = "audit_data"

// AuditData adalah data tambahan yang bisa diset oleh handler
type AuditData struct {
	Action      string
	EntityID    *uuid.UUID
	EntityType  string
	NewData     model.JSONB
	OldData     model.JSONB
	Changes     model.JSONB
	DeletedData model.JSONB
}

// SetAuditData menyimpan audit data ke context via pointer yang sudah ada,
// sehingga middleware bisa membacanya setelah handler selesai
func SetAuditData(r *http.Request, data *AuditData) *http.Request {
	// Jika sudah ada pointer di context, update isinya (mutate)
	if existing, ok := r.Context().Value(auditDataKey).(*AuditData); ok && existing != nil {
		*existing = *data
		return r
	}
	// Fallback: inject pointer baru (tidak akan terbaca middleware, tapi aman)
	ctx := context.WithValue(r.Context(), auditDataKey, data)
	return r.WithContext(ctx)
}

// responseWriter wraps http.ResponseWriter untuk menangkap status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// AuditLogMiddleware mencatat setiap request yang masuk ke audit_logs
func AuditLogMiddleware(authService services.AuthService, auditLogService services.AuditLogService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := newResponseWriter(w)

			// Inject pointer kosong ke context SEBELUM handler jalan
			// Handler akan mutate pointer ini via SetAuditData
			auditData := &AuditData{}
			ctx := context.WithValue(r.Context(), auditDataKey, auditData)
			r = r.WithContext(ctx)

			// Jalankan handler
			next.ServeHTTP(wrapped, r)

			// Ambil token dari header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return
			}

			// Ambil claims dari token (fallback ke DB jika token lama)
			claims, err := authService.GetFullClaimsFromToken(parts[1])
			if err != nil {
				log.Printf("[AuditLog] gagal parse token: %v", err)
				return
			}

			userID, err := uuid.Parse(claims["user_id"])
			if err != nil {
				log.Printf("[AuditLog] gagal parse user_id '%s': %v", claims["user_id"], err)
				return
			}

			roleID, err := uuid.Parse(claims["role_id"])
			if err != nil {
				log.Printf("[AuditLog] gagal parse role_id '%s': %v", claims["role_id"], err)
				return
			}

			// Hitung durasi request
			durationMs := int(time.Since(start).Milliseconds())

			// Ambil IP address
			ipAddress := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ipAddress = strings.Split(forwarded, ",")[0]
			}

			// Ambil user agent
			userAgent := r.UserAgent()

			// Buat audit log entry
			logEntry := &model.AuditLog{
				ID:         uuid.New(),
				UserID:     userID,
				FullName:   claims["full_name"],
				RoleID:     roleID,
				RoleName:   claims["role"],
				Method:     r.Method,
				Endpoint:   r.URL.Path,
				StatusCode: wrapped.statusCode,
				IPAddress:  &ipAddress,
				UserAgent:  &userAgent,
				DurationMs: &durationMs,
				CreatedAt:  time.Now(),
			}

			// Baca data yang diset handler via pointer
			if auditData.Action != "" {
				logEntry.Action = &auditData.Action
			}
			logEntry.EntityID = auditData.EntityID
			if auditData.EntityType != "" {
				logEntry.EntityType = &auditData.EntityType
			}
			logEntry.NewData = auditData.NewData
			logEntry.OldData = auditData.OldData
			logEntry.Changes = auditData.Changes
			logEntry.DeletedData = auditData.DeletedData

			// Gunakan context baru yang terpisah dari request context
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if saveErr := auditLogService.CreateAuditLog(ctx, logEntry); saveErr != nil {
					log.Printf("[AuditLog] gagal simpan: %v", saveErr)
				}
			}()
		}
	}
}

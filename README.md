# School Examination API

Sistem ujian online berbasis Go (Gin + GORM + PostgreSQL).

## Struktur Project

```
school-examination/
├── config/             # Konfigurasi aplikasi (.env loader)
├── database/           # Koneksi DB & auto-migrate
├── internal/
│   ├── handlers/       # HTTP handlers (controller)
│   ├── middleware/     # JWT auth & role middleware
│   ├── models/         # GORM models & request structs
│   ├── repository/     # Database queries
│   ├── services/       # Business logic
│   └── utils/          # JWT, password, response helpers
├── .env
├── go.mod
└── main.go
```

## Fitur

- **Role-based access**: super_admin, admin, teacher, student
- **Bank soal**: PG (pilihan ganda), esai, benar/salah
- **Jadwal ujian** per kelas & mata pelajaran
- **Timer otomatis & auto-submit** saat waktu habis
- **Acak soal & pilihan jawaban**
- **Koreksi otomatis** untuk PG dan benar/salah
- **Penilaian manual** untuk esai
- **Rekap nilai** per ujian
- **Anti-cheat flag** (full screen mode)
- JWT Authentication
- CORS support

## Setup

1. Pastikan PostgreSQL berjalan di port 5433
2. Buat database `school`
3. Copy `.env` dan sesuaikan konfigurasi
4. Jalankan:

```bash
go run main.go
```

Server akan otomatis:
- Migrasi semua tabel
- Seed super admin: `superadmin@school.com` / `password`

## API Endpoints

### Auth
| Method | Endpoint | Akses |
|--------|----------|-------|
| POST | `/api/v1/auth/register` | Public |
| POST | `/api/v1/auth/login` | Public |

### Users
| Method | Endpoint | Akses |
|--------|----------|-------|
| GET | `/api/v1/me` | Semua |
| GET | `/api/v1/users` | super_admin, admin |
| PUT | `/api/v1/users/:id` | super_admin, admin |
| DELETE | `/api/v1/users/:id` | super_admin |

### Mata Pelajaran
| Method | Endpoint | Akses |
|--------|----------|-------|
| GET | `/api/v1/subjects` | Semua |
| POST | `/api/v1/subjects` | super_admin, admin |

### Bank Soal
| Method | Endpoint | Akses |
|--------|----------|-------|
| GET | `/api/v1/questions` | super_admin, admin, teacher |
| POST | `/api/v1/questions` | super_admin, admin, teacher |
| PUT | `/api/v1/questions/:id` | super_admin, admin, teacher (mapel sendiri) |
| DELETE | `/api/v1/questions/:id` | super_admin, admin, teacher (mapel sendiri) |

### Kelas
| Method | Endpoint | Akses |
|--------|----------|-------|
| GET | `/api/v1/classes` | Semua |
| POST | `/api/v1/classes` | super_admin, admin |
| GET | `/api/v1/classes/:id/students` | super_admin, admin, teacher |
| POST | `/api/v1/classes/assign` | super_admin, admin |

### Ujian
| Method | Endpoint | Akses |
|--------|----------|-------|
| POST | `/api/v1/exams` | super_admin, admin, teacher |
| GET | `/api/v1/exams` | super_admin, admin, teacher |
| GET | `/api/v1/exams/:id/results` | super_admin, admin, teacher |
| POST | `/api/v1/exams/grade-essay` | super_admin, admin, teacher |

### Student
| Method | Endpoint | Akses |
|--------|----------|-------|
| GET | `/api/v1/student/exams` | student |
| POST | `/api/v1/student/exams/:id/start` | student |
| POST | `/api/v1/student/submissions/:id/answer` | student |
| POST | `/api/v1/student/submissions/:id/submit` | student |
| GET | `/api/v1/student/results` | student |

## Contoh Request

### Login
```json
POST /api/v1/auth/login
{
  "email": "superadmin@school.com",
  "password": "password"
}
```

### Buat Soal PG
```json
POST /api/v1/questions
Authorization: Bearer <token>
{
  "subject_id": 1,
  "type": "multiple_choice",
  "content": "Ibu kota Indonesia adalah?",
  "points": 2,
  "options": [
    { "content": "Jakarta", "is_correct": true },
    { "content": "Surabaya", "is_correct": false },
    { "content": "Bandung", "is_correct": false },
    { "content": "Medan", "is_correct": false }
  ]
}
```

### Buat Ujian
```json
POST /api/v1/exams
Authorization: Bearer <token>
{
  "title": "UTS Matematika Kelas X",
  "subject_id": 1,
  "class_id": 1,
  "start_time": "2026-05-20T08:00:00Z",
  "end_time": "2026-05-20T10:00:00Z",
  "duration_minutes": 90,
  "shuffle_questions": true,
  "shuffle_options": true,
  "anti_cheat": true,
  "passing_score": 70,
  "question_ids": [1, 2, 3, 4, 5]
}
```

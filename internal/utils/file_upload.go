package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UploadFile handles file upload and returns the file path
func UploadFile(file multipart.File, header *multipart.FileHeader, uploadDir string) (string, error) {
	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return "", fmt.Errorf("format file tidak didukung. Gunakan: jpg, jpeg, png, gif, webp")
	}

	// Validate file size (max 5MB)
	if header.Size > 5*1024*1024 {
		return "", fmt.Errorf("ukuran file terlalu besar. Maksimal 5MB")
	}

	// Create upload directory if not exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("gagal membuat direktori upload: %v", err)
	}

	// Generate unique filename
	timestamp := time.Now().Format("20060102150405")
	uniqueID := uuid.New().String()[:8]
	filename := fmt.Sprintf("%s_%s%s", timestamp, uniqueID, ext)
	filepath := filepath.Join(uploadDir, filename)

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("gagal membuat file: %v", err)
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("gagal menyimpan file: %v", err)
	}

	// Return relative path for URL
	return "/uploads/teachers/" + filename, nil
}

// DeleteFile deletes a file from the filesystem
func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	// Convert URL path to filesystem path
	// Remove leading slash and convert to OS-specific path
	fsPath := strings.TrimPrefix(filePath, "/")
	fsPath = filepath.FromSlash(fsPath)

	// Check if file exists
	if _, err := os.Stat(fsPath); os.IsNotExist(err) {
		return nil // File doesn't exist, no error
	}

	// Delete the file
	return os.Remove(fsPath)
}

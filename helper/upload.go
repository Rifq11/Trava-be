package helper

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	uploadsDir    = "public/uploads"
	maxFileSize   = 10 * 1024 * 1024 // 10MB
)

var allowedMimes = []string{
	"image/jpeg",
	"image/jpg",
	"image/png",
	"image/webp",
}

func init() {
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		panic(fmt.Sprintf("Failed to create uploads directory: %v", err))
	}
}

func GetFileUrl(filename string) string {
	if filename == "" {
		return ""
	}
	if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") {
		return filename
	}
	return fmt.Sprintf("/uploads/%s", filename)
}

// delete
func DeleteFile(filename string) error {
	if filename == "" {
		return nil
	}

	var filePath string
	if strings.Contains(filename, "/uploads/") {
		filePath = filepath.Join("public", strings.TrimPrefix(filename, "/uploads/"))
	} else {
		filePath = filepath.Join(uploadsDir, filename)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}

	return os.Remove(filePath)
}

// genereate file name, example: name-timestamp-random.ext
func generateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)
	uniqueSuffix := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int63n(1e9))
	return fmt.Sprintf("%s-%s%s", name, uniqueSuffix, ext)
}

func validateFileType(mimeType string) bool {
	for _, allowed := range allowedMimes {
		if mimeType == allowed {
			return true
		}
	}
	return false
}

// single
func UploadSingle(fieldName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile(fieldName)
		if err != nil {
			c.Next()
			return
		}

		if file.Size > maxFileSize {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "File size exceeds maximum limit of 10MB",
			})
			c.Abort()
			return
		}

		src, err := file.Open()
		if err != nil {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "Failed to read file",
			})
			c.Abort()
			return
		}
		defer src.Close()

		buffer := make([]byte, 512)
		_, err = src.Read(buffer)
		if err != nil && err != io.EOF {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "Failed to read file",
			})
			c.Abort()
			return
		}

		src.Seek(0, 0)

		mimeType := http.DetectContentType(buffer)
		if !validateFileType(mimeType) {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "Invalid file type. Only JPEG, PNG, and WebP images are allowed.",
			})
			c.Abort()
			return
		}

		filename := generateUniqueFilename(file.Filename)

		dst := filepath.Join(uploadsDir, filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Failed to save file",
			})
			c.Abort()
			return
		}

		c.Set("uploaded_file", filename)
		c.Next()
	}
}

// multiple
func UploadMultiple(fieldName string, maxCount int) gin.HandlerFunc {
	return func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.Next()
			return
		}

		files := form.File[fieldName]
		if len(files) == 0 {
			c.Next()
			return
		}

		if len(files) > maxCount {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Maximum %d files allowed", maxCount),
			})
			c.Abort()
			return
		}

		var filenames []string

		for _, file := range files {
			if file.Size > maxFileSize {
				c.JSON(400, gin.H{
					"status":  "error",
					"message": "File size exceeds maximum limit of 10MB",
				})
				c.Abort()
				return
			}

			src, err := file.Open()
			if err != nil {
				c.JSON(400, gin.H{
					"status":  "error",
					"message": "Failed to read file",
				})
				c.Abort()
				return
			}

			buffer := make([]byte, 512)
			_, err = src.Read(buffer)
			src.Close()

			if err != nil && err != io.EOF {
				c.JSON(400, gin.H{
					"status":  "error",
					"message": "Failed to read file",
				})
				c.Abort()
				return
			}

			mimeType := http.DetectContentType(buffer)
			if !validateFileType(mimeType) {
				c.JSON(400, gin.H{
					"status":  "error",
					"message": "Invalid file type. Only JPEG, PNG, and WebP images are allowed.",
				})
				c.Abort()
				return
			}

			filename := generateUniqueFilename(file.Filename)

			dst := filepath.Join(uploadsDir, filename)
			if err := c.SaveUploadedFile(file, dst); err != nil {
				c.JSON(500, gin.H{
					"status":  "error",
					"message": "Failed to save file",
				})
				c.Abort()
				return
			}

			filenames = append(filenames, filename)
		}

		c.Set("uploaded_files", filenames)
		c.Next()
	}
}


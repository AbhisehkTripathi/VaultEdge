package handlers

import (
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Document represents a document structure
type Document struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Type        string    `json:"type"`
	FolderID    int       `json:"folder_id"`
	UploadedAt  time.Time `json:"uploaded_at"`
	URL         string    `json:"url"`
}

// Folder represents a folder structure
type Folder struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	ParentID    *int      `json:"parent_id"`
	CreatedAt   time.Time `json:"created_at"`
	DocumentCount int     `json:"document_count"`
}

// Master represents master data structure
type Master struct {
	ID          int       `json:"id"`
	Type        string    `json:"type"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
}

// HealthCheck returns the health status of the application
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now(),
		"service":   "upload-document-saas",
		"version":   "1.0.0",
	})
}

// UploadDocument handles document upload
func UploadDocument(c *fiber.Ctx) error {
	// Parse multipart form
	file, err := c.FormFile("document")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Get folder ID from form (optional)
	folderIDStr := c.FormValue("folder_id")
	folderID := 1 // default folder
	if folderIDStr != "" {
		if id, err := strconv.Atoi(folderIDStr); err == nil {
			folderID = id
		}
	}

	// Validate file
	if err := validateFile(file); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Save file (in production, you'd save to cloud storage)
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	uploadPath := fmt.Sprintf("./uploads/%s", filename)
	
	if err := c.SaveFile(file, uploadPath); err != nil {
		log.Printf("Error saving file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Create document record (in production, save to database)
	document := Document{
		ID:         int(time.Now().Unix()), // Simple ID generation
		Name:       file.Filename,
		Size:       file.Size,
		Type:       filepath.Ext(file.Filename),
		FolderID:   folderID,
		UploadedAt: time.Now(),
		URL:        fmt.Sprintf("/uploads/%s", filename),
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Document uploaded successfully",
		"document": document,
	})
}

// GetDocumentByID retrieves a document by ID
func GetDocumentByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid document ID",
		})
	}

	// In production, fetch from database
	document := Document{
		ID:         id,
		Name:       "sample-document.pdf",
		Size:       1024000,
		Type:       ".pdf",
		FolderID:   1,
		UploadedAt: time.Now().Add(-24 * time.Hour),
		URL:        "/uploads/sample-document.pdf",
	}

	return c.JSON(fiber.Map{
		"document": document,
	})
}

// ListDocuments retrieves all documents with pagination
func ListDocuments(c *fiber.Ctx) error {
	// Get query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	folderID := c.Query("folder_id")

	// Sample documents (in production, fetch from database)
	documents := []Document{
		{
			ID:         1,
			Name:       "document1.pdf",
			Size:       1024000,
			Type:       ".pdf",
			FolderID:   1,
			UploadedAt: time.Now().Add(-24 * time.Hour),
			URL:        "/uploads/document1.pdf",
		},
		{
			ID:         2,
			Name:       "document2.docx",
			Size:       512000,
			Type:       ".docx",
			FolderID:   1,
			UploadedAt: time.Now().Add(-12 * time.Hour),
			URL:        "/uploads/document2.docx",
		},
	}

	// Filter by folder if specified
	if folderID != "" {
		// In production, filter by folder_id in database query
		log.Printf("Filtering by folder ID: %s", folderID)
	}

	return c.JSON(fiber.Map{
		"documents": documents,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": len(documents),
		},
	})
}

// ListFolders retrieves all folders
func ListFolders(c *fiber.Ctx) error {
	// Sample folders (in production, fetch from database)
	folders := []Folder{
		{
			ID:            1,
			Name:          "Documents",
			ParentID:      nil,
			CreatedAt:     time.Now().Add(-7 * 24 * time.Hour),
			DocumentCount: 5,
		},
		{
			ID:            2,
			Name:          "Images",
			ParentID:      nil,
			CreatedAt:     time.Now().Add(-5 * 24 * time.Hour),
			DocumentCount: 3,
		},
		{
			ID:            3,
			Name:          "Contracts",
			ParentID:      &[]int{1}[0],
			CreatedAt:     time.Now().Add(-3 * 24 * time.Hour),
			DocumentCount: 2,
		},
	}

	return c.JSON(fiber.Map{
		"folders": folders,
	})
}

// ListMasters retrieves master data
func ListMasters(c *fiber.Ctx) error {
	masterType := c.Query("type")

	// Sample master data (in production, fetch from database)
	masters := []Master{
		{
			ID:          1,
			Type:        "document_type",
			Value:       "pdf",
			Description: "PDF Document",
			IsActive:    true,
		},
		{
			ID:          2,
			Type:        "document_type",
			Value:       "docx",
			Description: "Word Document",
			IsActive:    true,
		},
		{
			ID:          3,
			Type:        "folder_type",
			Value:       "public",
			Description: "Public Folder",
			IsActive:    true,
		},
	}

	// Filter by type if specified
	if masterType != "" {
		var filteredMasters []Master
		for _, master := range masters {
			if master.Type == masterType {
				filteredMasters = append(filteredMasters, master)
			}
		}
		masters = filteredMasters
	}

	return c.JSON(fiber.Map{
		"masters": masters,
	})
}

// SendKafkaTestMessage sends a test message to Kafka
func SendKafkaTestMessage(c *fiber.Ctx) error {
	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON payload",
		})
	}

	// In production, send to Kafka
	log.Printf("Kafka test message: %+v", payload)

	return c.JSON(fiber.Map{
		"message": "Test message sent to Kafka",
		"payload": payload,
		"timestamp": time.Now(),
	})
}

// validateFile validates uploaded file
func validateFile(file *multipart.FileHeader) error {
	// Check file size (10MB limit)
	maxSize := int64(10 * 1024 * 1024)
	if file.Size > maxSize {
		return fmt.Errorf("file size exceeds 10MB limit")
	}

	// Check file extension
	allowedExts := map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".txt":  true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}

	ext := filepath.Ext(file.Filename)
	if !allowedExts[ext] {
		return fmt.Errorf("file type %s not allowed", ext)
	}

	return nil
}

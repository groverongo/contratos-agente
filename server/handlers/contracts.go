package handlers

import (
	"context"
	"net/http"
	"net/url" // Added this
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

// CreateContract handles creating a contract metadata and returning a Presigned PUT URL
func (h *Handler) CreateContract(c echo.Context) error {
	// 1. Get user ID from Auth
	userID := c.Request().Header.Get("X-Stack-Auth-User-Id")
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User ID missing"})
	}

	type Req struct {
		Title       string `json:"title"`
		Filename    string `json:"filename"`
		ContentType string `json:"content_type"`
		Size        int64  `json:"size"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.Title == "" || req.Filename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title and Filename required"})
	}

	// Default parameters
	if req.ContentType == "" {
		req.ContentType = "application/pdf"
	}

	bucketName := "contracts"
	// Ensure bucket exists
	exists, err := h.MinIO.BucketExists(context.Background(), bucketName)
	if err == nil && !exists {
		h.MinIO.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	}

	objectName := uuid.New().String() + "-" + req.Filename

	// 2. Generate Presigned PUT URL
	// We can set expiry (e.g., 15 minutes)
	expiry := time.Minute * 15
	presignedURL, err := h.MinIO.PresignedPutObject(context.Background(), bucketName, objectName, expiry)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate upload URL"})
	}

	// 3. Create Contract in DB
	contract := Contract{
		ID:       uuid.New(),
		Title:    req.Title,
		AuthorID: userID,
		Status:   "Pending Upload", // Status indicating file is not yet confirmed uploaded
	}

	if err := h.DB.Create(&contract).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create contract"})
	}

	// 4. Create Contract Version
	version := ContractVersion{
		ID:            uuid.New(),
		ContractID:    contract.ID,
		VersionNumber: 1,
		FilePath:      objectName,
	}
	if err := h.DB.Create(&version).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create version"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"contract":   contract,
		"upload_url": presignedURL.String(),
	})
}

// ListContracts lists contracts for the user (as author or recipient)
func (h *Handler) ListContracts(c echo.Context) error {
	userID := c.Request().Header.Get("X-Stack-Auth-User-Id")

	var contracts []Contract
	// Simple query: where author_id = userID OR exists in recipients
	// For simplicity, just author for now or join
	err := h.DB.Preload("Recipients").Preload("Versions").Where("author_id = ?", userID).Find(&contracts).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Db error"})
	}

	return c.JSON(http.StatusOK, contracts)
}

// GetContractFileUrl generates a presigned URL for the latest version
func (h *Handler) GetContractFileUrl(c echo.Context) error {
	id := c.Param("id")
	var contract Contract
	if err := h.DB.Preload("Versions").First(&contract, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Contract not found"})
	}

	// Get latest version
	if len(contract.Versions) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "No file versions"})
	}
	latest := contract.Versions[len(contract.Versions)-1] // Assuming order, otherwise sort

	// Presign URL
	reqParams := make(url.Values)
	presignedURL, err := h.MinIO.PresignedGetObject(context.Background(), "contracts", latest.FilePath, time.Hour*24, reqParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to sign URL"})
	}

	return c.JSON(http.StatusOK, map[string]string{"url": presignedURL.String()})
}

// UpdateRecipients sets the recipients
func (h *Handler) UpdateRecipients(c echo.Context) error {
	id := c.Param("id")
	type Req struct {
		Emails []string `json:"emails"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Transaction to update recipients
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// Verify ownership?
		// Delete existing? Or just add new? "Set" usually means replace.
		// For simplicity, we just add new ones or ignore if exist.
		// Detailed logic: Remove those not in list, add new ones.

		// 1. Delete all current recipients
		if err := tx.Delete(&ContractRecipient{}, "contract_id = ?", id).Error; err != nil {
			return err
		}

		// 2. Add new
		for _, email := range req.Emails {
			r := ContractRecipient{
				ContractID:     uuid.MustParse(id),
				RecipientEmail: email,
				Status:         "Pending",
			}
			if err := tx.Create(&r).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update recipients"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

// SignContract handles a user signing the contract
func (h *Handler) SignContract(c echo.Context) error {
	id := c.Param("id")

	// Identify who is signing
	// userID := c.Request().Header.Get("X-Stack-Auth-User-Id")
	// email := ... (fetching user email from StackAuth or DB)
	// For MVP, passing email in body
	type Req struct {
		Email string `json:"email"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email required"})
	}

	var recipient ContractRecipient
	if err := h.DB.First(&recipient, "contract_id = ? AND recipient_email = ?", id, req.Email).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Recipient not found"})
	}

	recipient.Status = "Signed"
	now := time.Now()
	recipient.SignedAt = &now

	h.DB.Save(&recipient)

	return c.JSON(http.StatusOK, recipient)
}

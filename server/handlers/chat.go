package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
)

// AskAI proxies the question to the AI Agent service
func (h *Handler) AskAI(c echo.Context) error {
	type Req struct {
		ContractID string `json:"contract_id"`
		Question   string `json:"question"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Retrieve file path from Contract to pass to AI?
	// Or AI agent just needs the file URL?
	// The implementation plan says: "AI assistant will have context about the contract by automatically uploading the pdf file as soon as the firsts message is sent".
	// We can pass the Presigned URL to the AI agent.

	// 1. Get Contract
	var contract Contract
	if err := h.DB.Preload("Versions").First(&contract, "id = ?", req.ContractID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Contract not found"})
	}
	if len(contract.Versions) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Contract has no file"})
	}
	latest := contract.Versions[len(contract.Versions)-1]

	// 2. Generate Presigned URL
	// Presign URL
	reqParams := make(url.Values)
	presignedURL, err := h.MinIO.PresignedGetObject(context.Background(), "contracts", latest.FilePath, time.Hour*1, reqParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to sign URL"})
	}

	aiReq := map[string]interface{}{
		"contract_id": req.ContractID,
		"question":    req.Question,
		"file_url":    presignedURL.String(),
	}

	jsonData, _ := json.Marshal(aiReq)
	resp, err := http.Post("http://ai-agent:3000/ask", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "AI Agent unavailable"})
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	return c.JSON(http.StatusOK, map[string]string{"response": string(body)})
}

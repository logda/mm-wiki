package controllers

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/astaxie/beego"
	"github.com/phachon/mm-wiki/app/utils"
)

type AIChatController struct {
	beego.Controller
}

type AIChatRequest struct {
	Document string `json:"document"`
	Message  string `json:"message"`
}

type AIChatResponse struct {
	Reply string `json:"reply"`
}

func (c *AIChatController) Post() {
	body, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.handleError(400, "Invalid request: "+err.Error(), "Failed to read request body")
		return
	}

	c.Ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var req AIChatRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.handleError(400, "Invalid request: "+err.Error(), "Failed to parse JSON body")
		return
	}

	if req.Document == "" || req.Message == "" {
		c.handleError(400, "Invalid request", "Both 'document' and 'message' fields are required")
		return
	}

	reply, err := callLargeLanguageModel(req.Document, req.Message)
	if err != nil {
		c.handleError(500, "Internal server error", err.Error())
		return
	}

	c.Data["json"] = AIChatResponse{Reply: reply}
	c.ServeJSON()
}

func (c *AIChatController) handleError(status int, errorMsg, details string) {
	beego.Error(errorMsg + ": " + details)
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]interface{}{
		"error":   errorMsg,
		"details": details,
	}
	c.ServeJSON()
}

func callLargeLanguageModel(document, message string) (string, error) {
	// Read from the configuration

	apiKey := beego.AppConfig.String("llm_key")
	model := beego.AppConfig.String("llm_name")
	chatURL := beego.AppConfig.String("llm_url")

	client := utils.NewDashScopeClient(apiKey, model, chatURL)

	response, err := client.ChatCompletion(document, message)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}

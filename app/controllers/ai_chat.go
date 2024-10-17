package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

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

func (this *AIChatController) Post() {
	// Log the raw request body
	body, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		beego.Error("Failed to read request body:", err)
		this.Ctx.Output.SetStatus(400)
		this.Data["json"] = map[string]interface{}{
			"error":   "Invalid request: " + err.Error(),
			"details": "Failed to read request body",
		}
		this.ServeJSON()
		return
	}
	// beego.Info("Request Body:", string(body))

	// Restore the body to the request so it can be read again
	this.Ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var req AIChatRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		beego.Error("Failed to unmarshal request body:", err)
		this.Ctx.Output.SetStatus(400)
		this.Data["json"] = map[string]interface{}{
			"error":   "Invalid request: " + err.Error(),
			"details": "Failed to parse JSON body",
		}
		this.ServeJSON()
		return
	}

	// Log the parsed request
	// beego.Info("Parsed request:", req)

	// Validate required fields
	if req.Document == "" || req.Message == "" {
		this.Ctx.Output.SetStatus(400)
		this.Data["json"] = map[string]interface{}{
			"error":   "Invalid request",
			"details": "Both 'document' and 'message' fields are required",
		}
		this.ServeJSON()
		return
	}

	// Call the large language model
	reply, err := callLargeLanguageModel(req.Document, req.Message)
	if err != nil {
		beego.Error("Error calling large language model:", err)
		this.Ctx.Output.SetStatus(500)
		this.Data["json"] = map[string]interface{}{
			"error":   "Internal server error",
			"details": err.Error(),
		}
		this.ServeJSON()
		return
	}

	this.Data["json"] = AIChatResponse{Reply: reply}
	this.ServeJSON()
}

// func callLargeLanguageModel(document, message string) (string, error) {

// 	client := utils.NewDashScopeClient("sk-db3c48935a704e1994258b928aa3ef24")

// 	response, err := client.ChatCompletion(document, message)
// 	if err != nil {
// 		return "", err
// 	}

// 	return response.Choices[0].Message.Content, nil
// }

func callLargeLanguageModel(document, message string) (string, error) {
	// Read the API key from the configuration
	apiKey := beego.AppConfig.String("dashscope_api_key")

	client := utils.NewDashScopeClient(apiKey)

	response, err := client.ChatCompletion(document, message)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}

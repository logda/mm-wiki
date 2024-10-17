package utils

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

const DashScopeAPIURL = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"

type DashScopeClient struct {
    APIKey string
}

type DashScopeRequest struct {
    Model    string    `json:"model"`
    Messages []Message `json:"messages"`
}

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type DashScopeResponse struct {
    Choices []struct {
        Message struct {
            Role    string `json:"role"`
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
}

func NewDashScopeClient(apiKey string) *DashScopeClient {
    return &DashScopeClient{APIKey: apiKey}
}

func (c *DashScopeClient) ChatCompletion(document, message string) (*DashScopeResponse, error) {
    req := DashScopeRequest{
        Model: "qwen-plus",
        Messages: []Message{
            {Role: "system", Content: "You are a helpful assistant. Use the following document as context for answering the user's question: " + document},
            {Role: "user", Content: message},
        },
    }

    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    httpReq, err := http.NewRequest("POST", DashScopeAPIURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }

    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

    client := &http.Client{}
    resp, err := client.Do(httpReq)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
    }

    var dashScopeResp DashScopeResponse
    err = json.Unmarshal(body, &dashScopeResp)
    if err != nil {
        return nil, err
    }

    return &dashScopeResp, nil
}
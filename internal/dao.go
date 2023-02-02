package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/calebtracey/ai-interaction-api/external"
	"io"
	"net/http"
	"os"
)

const (
	imageURL    = "https://api.openai.com/v1/images/generations"
	contentType = "application/json"
)

var aiError = errors.New("internal server error")

type DAOI interface {
	GenerateImage(ctx context.Context, apiRequest *http.Request) (apiResp external.APIResponse)
}
type DAO struct {
	Client *http.Client
}

func (s DAO) GenerateImage(ctx context.Context, apiRequest *http.Request) (apiResp external.APIResponse) {
	req, httpErr := http.NewRequest(http.MethodPost, imageURL, io.NopCloser(apiRequest.Body))

	if httpErr != nil {
		apiResp.Message.ErrorLog = external.ErrorLogs{{
			ExceptionType: "status bad request",
			StatusCode:    "400",
			Trace:         fmt.Sprintf("GenerateImage: error: %v", httpErr),
			RootCause:     httpErr.Error(),
		}}
		return apiResp
	}

	req.Header.Add("Authorization", os.Getenv("API_KEY"))
	req.Header.Add("Content-Type", "application/json")

	clientResp, clientErr := s.Client.Do(req.WithContext(ctx))

	defer clientResp.Body.Close()

	if clientErr != nil {
		apiResp.Message.ErrorLog = external.ErrorLogs{{
			ExceptionType: "internal server error",
			StatusCode:    "500",
			Trace:         fmt.Sprintf("GenerateImage: error: %v", clientErr),
			RootCause:     clientErr.Error(),
		}}
		return apiResp
	}

	var aiResp external.AIResponse

	if jsonErr := json.NewDecoder(clientResp.Body).Decode(&aiResp); jsonErr != nil {
		apiResp.Message.ErrorLog = external.ErrorLogs{{
			ExceptionType: "internal server error",
			StatusCode:    "500",
			Trace:         fmt.Sprintf("GenerateImage: error: %v", jsonErr),
			RootCause:     jsonErr.Error(),
		}}
		return apiResp
	}

	apiResp.Result = aiResp

	return apiResp
}

package facade

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/calebtracey/ai-interaction-api/external"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strconv"
)

const (
	imageURL    = "https://api.openai.com/v1/images/generations"
	apiKey      = "API_KEY"
	contentType = "application/json"
)

type ServiceI interface {
	GenerateImage(ctx context.Context, apiRequest external.APIRequest) (apiResp external.APIResponse)
}

type Service struct {
	Client *http.Client
}

func addHeaders(req *http.Request) *http.Request {
	req.Header.Add("Authorization", os.Getenv(apiKey))
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Access-Control-Allow-Origin", "*")
	return req
}

func (s Service) GenerateImage(ctx context.Context, apiRequest external.APIRequest) (apiResp external.APIResponse) {
	var aiResp external.AIResponse

	if errs := apiRequest.Validate(); len(errs) > 0 {
		for _, err := range errs {
			apiResp.Message.ErrorLog = append(apiResp.Message.ErrorLog, errorLog(err, 400, "Validate"))
		}
		return apiResp
	}

	// make dao call here
	requestBuf := new(bytes.Buffer)
	if jsonErr := json.NewDecoder(requestBuf).Decode(&apiRequest); jsonErr != nil {
		return responseWithError(apiResp, jsonErr, http.StatusBadRequest, "Decode")

	}

	request, httpErr := http.NewRequest(http.MethodPost, imageURL, io.NopCloser(requestBuf))
	if httpErr != nil {
		return responseWithError(apiResp, httpErr, http.StatusBadRequest, "GenerateImage")
	}

	clientResp, clientErr := s.Client.Do(addHeaders(request.WithContext(ctx)))

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Error(err)
		}
	}(clientResp.Body)

	if clientErr != nil {
		return responseWithError(apiResp, clientErr, http.StatusInternalServerError, "GenerateImage")
	}

	if jsonErr := json.NewDecoder(clientResp.Body).Decode(&aiResp); jsonErr != nil {
		return responseWithError(apiResp, jsonErr, http.StatusInternalServerError, "GenerateImage")
	}

	apiResp.Result = aiResp

	return apiResp
}

func errorLog(err error, code int, trace string) external.ErrorLog {
	return external.ErrorLog{
		ExceptionType: http.StatusText(code),
		StatusCode:    strconv.Itoa(code),
		Trace:         fmt.Sprintf("%s: error: %v", trace, err),
		RootCause:     err.Error(),
	}
}

// responseWithError adds an error log and returns the response
func responseWithError(resp external.APIResponse, err error, code int, trace string) external.APIResponse {
	resp.Message.ErrorLog = external.ErrorLogs{{
		ExceptionType: http.StatusText(code),
		StatusCode:    strconv.Itoa(code),
		Trace:         fmt.Sprintf("%s: error: %v", trace, err),
		RootCause:     err.Error(),
	}}
	return resp
}

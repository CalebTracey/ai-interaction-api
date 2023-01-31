package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/calebtracey/ai-interaction-api/external"
	"net/http"
	"os"
)

const (
	imageURL    = "https://api.openai.com/v1/images/generations"
	contentType = "application/json"
)

type DAOI interface {
	GenerateImage(ctx context.Context, apiRequest external.APIRequest) (apiResp external.AIResponse, err external.Errors)
}
type DAO struct {
	Client *http.Client
}

func (s DAO) GenerateImage(ctx context.Context, apiRequest external.APIRequest) (apiResp external.AIResponse, err external.Errors) {
	buf := new(bytes.Buffer)
	if jsonErr := json.NewEncoder(buf).Encode(apiRequest); err != nil {
		return apiResp, external.Errors{
			{
				Error:      jsonErr.Error(),
				StatusCode: http.StatusBadRequest,
				Trace:      "GenerateImage",
			},
		}
	}

	req, httpErr := http.NewRequestWithContext(ctx, http.MethodPost, imageURL, buf)
	if httpErr != nil {
		return apiResp, external.Errors{
			{
				Error:      httpErr.Error(),
				StatusCode: http.StatusBadRequest,
				Trace:      "GenerateImage",
			},
		}
	}

	req.Header.Add("Authorization", os.Getenv("API_KEY"))
	req.Header.Add("Content-Type", "application/json")

	aiResp, clientErr := s.Client.Do(req)
	if clientErr != nil {
		return apiResp, external.Errors{
			{
				Error:      clientErr.Error(),
				StatusCode: http.StatusBadRequest,
				Trace:      "GenerateImage",
			},
		}
	}

	if aiResp.StatusCode != http.StatusOK {
		return apiResp, external.Errors{
			{
				Error:      clientErr.Error(),
				StatusCode: http.StatusBadRequest,
				Trace:      "GenerateImage",
			},
		}
	}

	if jsonErr := json.NewDecoder(aiResp.Body).Decode(&apiResp); jsonErr != nil {
		return apiResp, external.Errors{
			{
				Error:      jsonErr.Error(),
				StatusCode: http.StatusBadRequest,
				Trace:      "GenerateImage",
			},
		}
	}

	return apiResp, nil
}

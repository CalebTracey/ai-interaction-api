package openai

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
)

const (
	imageURL    = "https://api.openai.com/v1/images/generations"
	apiKey      = "API_KEY"
	contentType = "application/json"
)

type DAOI interface {
	ImageRequest(ctx context.Context, apiRequest *external.APIRequest) (resp external.AIResponse, err error)
}

type DAO struct {
	Client *http.Client
}

func (s DAO) ImageRequest(ctx context.Context, apiRequest *external.APIRequest) (resp external.AIResponse, err error) {
	reqBytes, jsonErr := json.Marshal(apiRequest)
	if jsonErr != nil {
		return resp, fmt.Errorf("ImageRequest: %w", jsonErr)
	}

	request, reqErr := http.NewRequest(http.MethodPost, imageURL, io.NopCloser(bytes.NewReader(reqBytes)))
	if reqErr != nil {
		return resp, fmt.Errorf("ImageRequest: %w", reqErr)
	}

	clientResp, clientErr := s.Client.Do(addHeaders(request.WithContext(ctx)))

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Error(err)
		}
	}(clientResp.Body)

	if clientErr != nil {
		return resp, fmt.Errorf("ImageRequest: %w", clientErr)
	}

	if jsonErr = json.NewDecoder(clientResp.Body).Decode(&resp); jsonErr != nil {
		return resp, fmt.Errorf("ImageRequest: %w", jsonErr)
	}

	return resp, nil
}

func addHeaders(req *http.Request) *http.Request {
	req.Header.Add("Authorization", os.Getenv(apiKey))
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Access-Control-Allow-Origin", "*")
	return req
}

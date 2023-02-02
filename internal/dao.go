package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/calebtracey/ai-interaction-api/external"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"os"
	"strconv"
)

const (
	imageURL    = "https://api.openai.com/v1/images/generations"
	contentType = "application/json"
)

var aiError = errors.New("internal server error")

type DAOI interface {
	GenerateImage(ctx context.Context, apiRequest external.APIRequest) (apiResp external.APIResponse)
}
type DAO struct {
	Client *http.Client
}

func (s DAO) GenerateImage(ctx context.Context, apiRequest external.APIRequest) (apiResp external.APIResponse) {
	buf := new(bytes.Buffer)
	g, ctx := errgroup.WithContext(ctx)

	if jsonErr := json.NewEncoder(buf).Encode(&apiRequest); jsonErr != nil {
		apiResp.Message.ErrorLog = external.ErrorLogs{{
			ExceptionType: "status bad request",
			StatusCode:    "400",
			Trace:         fmt.Sprintf("GenerateImage: error: %v", jsonErr),
			RootCause:     jsonErr.Error(),
		}}
		return apiResp
	}

	req, httpErr := http.NewRequestWithContext(ctx, http.MethodPost, imageURL, buf)
	if httpErr != nil {
		apiResp.Message.ErrorLog = external.ErrorLogs{{
			ExceptionType: "status bad request",
			StatusCode:    "400",
			Trace:         fmt.Sprintf("GenerateImage: error: %v", httpErr),
			RootCause:     httpErr.Error(),
		}}
		return apiResp
	}

	key := os.Getenv("API_KEY")
	log.Infof("=== Auth Header: %v", key)
	req.Header.Add("Authorization", key)
	req.Header.Add("Content-Type", "application/json")

	responseChan := make(chan external.APIResponse, 1)

	g.Go(func() error {
		defer close(responseChan)
		var tempResp external.APIResponse
		aiResp, clientErr := s.Client.Do(req)

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Error(err)
			}
		}(aiResp.Body)

		if clientErr != nil {
			return clientErr
		}

		if aiResp.StatusCode != http.StatusOK {
			tempResp.Message.ErrorLog = external.ErrorLogs{{
				StatusCode: strconv.Itoa(aiResp.StatusCode),
				Trace:      "GenerateImage: error",
			}}
			return aiError
		}

		if jsonErr := json.NewDecoder(aiResp.Body).Decode(&tempResp); jsonErr != nil {
			tempResp.Message.ErrorLog = external.ErrorLogs{{
				ExceptionType: "internal server error",
				StatusCode:    "500",
				Trace:         fmt.Sprintf("GenerateImage: error: %v", jsonErr),
				RootCause:     jsonErr.Error(),
			}}
			return aiError
		}

		responseChan <- tempResp
		return nil
	})

	g.Go(func() error {
		for res := range responseChan {
			apiResp = res
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		apiResp.Message.ErrorLog = external.ErrorLogs{{
			ExceptionType: "status internal server error",
			StatusCode:    "500",
			Trace:         fmt.Sprintf("GenerateImage: error: %v", err),
			RootCause:     err.Error(),
		}}
		return apiResp
	}

	//if clientErr != nil {
	//	apiResp.Message.ErrorLog = external.ErrorLogs{{
	//		ExceptionType: "status internal server error",
	//		StatusCode:    "500",
	//		Trace:         fmt.Sprintf("GenerateImage: error: %v", clientErr),
	//		RootCause:     clientErr.Error(),
	//	}}
	//	return apiResp
	//}

	return apiResp
}

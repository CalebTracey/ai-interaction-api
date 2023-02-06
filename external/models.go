package external

import (
	"errors"
	"fmt"
	"net/http"
)

const REQUEST_LIMIT = 10

var (
	MissingSizeError   = errors.New("request parameter 'size' is required")
	MissingPromptError = errors.New("request parameter 'prompt' is required")
	InvalidSizeError   = errors.New("request parameter 'n' cannot equal 0")
	AboveLimitError    = errors.New(fmt.Sprintf("request parameter 'n' cannot be higher than request limit: %d", REQUEST_LIMIT))
)

type APIRequest struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

func (req *APIRequest) Validate() (errs []error) {
	if req.Size == "" {
		errs = append(errs, MissingSizeError)
	}
	if req.Prompt == "" {
		errs = append(errs, MissingPromptError)
	}
	if req.N == 0 {
		errs = append(errs, InvalidSizeError)
	}
	if req.N >= REQUEST_LIMIT {
		errs = append(errs, AboveLimitError)
	}
	return errs
}

type APIResponse struct {
	Result  AIResponse `json:"result"`
	Message Message    `json:"message"`
}

func (res *APIResponse) StatusCode() (code int) {
	if len(res.Message.ErrorLog) > 0 {
		return res.Message.ErrorLog.GetHTTPStatus(len(res.Result.Data))
	}
	return http.StatusOK
}

type AIResponse struct {
	Created int `json:"created"`
	Data    []struct {
		Url string `json:"url"`
	} `json:"data"`
}

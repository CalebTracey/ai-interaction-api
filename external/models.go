package external

import (
	"errors"
	"fmt"
	"net/http"
)

const AmountLimit = 20

const (
	LargeImage  = "1024x1024"
	MediumImage = "512x512"
	SmallImage  = "256x256"
)

var (
	MissingSizeError   = errors.New("request parameter 'size' is required")
	MissingPromptError = errors.New("request parameter 'prompt' is required")
	InvalidSizeError   = errors.New("request parameter 'n' is invalid")
	AboveLimitError    = errors.New(fmt.Sprintf("request parameter 'n' cannot be higher than request limit: %d", AmountLimit))
)

type APIRequest struct {
	Prompt string `json:"prompt,omitempty"`
	N      int    `json:"n,omitempty"`
	Size   string `json:"size,omitempty"`
}

type Errors []error

type BadRequest struct {
	Errors
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
	Created int       `json:"created"`
	Data    GenImages `json:"data"`
}

type GenImages []GenImage

type GenImage struct {
	Url string `json:"url"`
}

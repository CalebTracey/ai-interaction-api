package external

type APIRequest struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type APIResponse struct {
	Result  AIResponse `json:"result"`
	Message Message    `json:"message"`
}

func (res *APIResponse) StatusCode() (code int) {
	if len(res.Message.ErrorLog) > 0 {
		code = res.Message.ErrorLog.GetHTTPStatus(len(res.Result.Data))
	}
	return code
}

type AIResponse struct {
	Created int `json:"created"`
	Data    []struct {
		Url string `json:"url"`
	} `json:"data"`
}

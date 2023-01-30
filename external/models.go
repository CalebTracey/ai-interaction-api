package external

type Errors []Error
type Error struct {
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
	Trace      string `json:"trace,omitempty"`
}
type APIRequest struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type AIResponse struct {
	Created int `json:"created"`
	Data    []struct {
		Url string `json:"url"`
	} `json:"data"`
}

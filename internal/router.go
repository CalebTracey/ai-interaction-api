package internal

import (
	"encoding/json"
	"github.com/calebtracey/ai-interaction-api/external"
	"github.com/calebtracey/ai-interaction-api/internal/facade"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Handler struct {
	Service facade.ServiceI
}

func (h Handler) InitializeRoutes() *gin.Engine {
	r := gin.Default()

	r.Group("/v1").POST("/image", h.imageHandler())
	r.Group("/v1").POST("/test", h.testHandler())

	return r
}
func (h Handler) imageHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sw := time.Now()

		var apiResponse external.APIResponse
		var apiRequest external.APIRequest

		if err := json.NewDecoder(ctx.Request.Body).Decode(&apiRequest); err != nil {
			apiResponse.Message.AddMessageDetails(sw)
			ctx.JSON(http.StatusBadRequest, apiResponse)
			return
		}

		if apiResponse = h.Service.GenerateImage(ctx, apiRequest); apiResponse.StatusCode() == http.StatusOK {
			apiResponse.Message.AddMessageDetails(sw)
			ctx.JSON(http.StatusOK, apiResponse)
			return

		} else {

			log.Errorf("imageHandler: error: %v", apiResponse.Message.ErrorLog)
			apiResponse.Message.AddMessageDetails(sw)
			ctx.JSON(apiResponse.StatusCode(), apiResponse)
			return
		}
	}
}

func (h Handler) testHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp := external.APIResponse{
			Result: external.AIResponse{
				Created: 2,
				Data: []struct {
					Url string `json:"url"`
				}{{Url: "https://oaidalleapiprodscus.blob.core.windows.net/private/org-BOMcU14BHoW1rBIBVWcFPDPn/user-XDlrwYWHipMjIdRv4dmtwNOV/img-Swt7XBlTf6srv63bankSiew6.png?st=2023-01-29T20%3A18%3A52Z&se=2023-01-29T22%3A18%3A52Z&sp=r&sv=2021-08-06&sr=b&rscd=inline&rsct=image/png&skoid=6aaadede-4fb3-4698-a8f6-684d7786b067&sktid=a48cca56-e6da-484e-a814-9c849652bcb3&skt=2023-01-29T18%3A29%3A13Z&ske=2023-01-30T18%3A29%3A13Z&sks=b&skv=2021-08-06&sig=FF3xaUHl5CQlrHRWF55F%2BOzATeGDN8qsNmRRDJcppbQ%3D"}, {Url: "https://oaidalleapiprodscus.blob.core.windows.net/private/org-BOMcU14BHoW1rBIBVWcFPDPn/user-XDlrwYWHipMjIdRv4dmtwNOV/img-Nvym99je5OLe5ikGEgzcMJ3L.png?st=2023-01-29T20%3A18%3A52Z&se=2023-01-29T22%3A18%3A52Z&sp=r&sv=2021-08-06&sr=b&rscd=inline&rsct=image/png&skoid=6aaadede-4fb3-4698-a8f6-684d7786b067&sktid=a48cca56-e6da-484e-a814-9c849652bcb3&skt=2023-01-29T18%3A29%3A13Z&ske=2023-01-30T18%3A29%3A13Z&sks=b&skv=2021-08-06&sig=wjlq5B2ZjZMM%2BS%2B16J8rOHQxq/HmnD7zjZ8FnUFaCA4%3D"}},
			},
			Message: external.Message{
				ErrorLog:  nil,
				HostName:  "test",
				Status:    "200",
				TimeTaken: "not enabled for '/test",
				Count:     2,
			},
		}
		time.Sleep(5 * time.Second)
		ctx.JSON(http.StatusOK, resp)
	}
}

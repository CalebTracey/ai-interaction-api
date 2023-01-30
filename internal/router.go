package internal

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"images-ai/external"
	"net/http"
)

type Handler struct {
	DAO DAOI
}

func (h Handler) InitializeRoutes() *gin.Engine {

	r := gin.Default()
	r.POST("/image", h.imageHandler())

	return r
}

type AIResponse struct {
	Created string `json:"created"`
	Data    images `json:"data"`
}

type images []image

type image struct {
	Url string `json:"url"`
}

func (h Handler) imageHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var apiRequest external.APIRequest

		if err := json.NewDecoder(ctx.Request.Body).Decode(&apiRequest); err != nil {
			panic(err)
		}

		apiResp, daoErrs := h.DAO.GenerateImage(ctx, apiRequest)
		if daoErrs != nil {
			log.Error(daoErrs)
			ctx.JSON(daoErrs[0].StatusCode, daoErrs)
		}

		ctx.JSON(http.StatusOK, apiResp.Data)
	}
}

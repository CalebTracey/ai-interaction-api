package main

import (
	"github.com/NYTimes/gziphandler"
	"github.com/calebtracey/ai-interaction-api/internal"
	config "github.com/calebtracey/config-yaml"
	"log"
)

const configPath = "config.yaml"

type builder struct{}

func (b *builder) initializeDAO() (*internal.DAO, error) {
	//if err := godotenv.Load(); err != nil {
	//	return nil, err
	//}

	cfg := config.New(configPath)

	openAiSvc, err := cfg.Service("openAi")
	if err != nil {
		return nil, err
	}

	return &internal.DAO{
		Client: openAiSvc.Client,
	}, nil
}

func main() {
	if svc, err := new(builder).initializeDAO(); err != nil {
		log.Fatal(err)
	} else {
		run(svc)
	}
}

func run(svc *internal.DAO) {
	handler := internal.Handler{DAO: svc}

	log.Fatal(listenAndServe("8080", gziphandler.GzipHandler(corsHandler().Handler(handler.InitializeRoutes()))))
}

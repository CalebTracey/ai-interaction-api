package main

import (
	"github.com/NYTimes/gziphandler"
	"github.com/calebtracey/ai-interaction-api/internal"
	config "github.com/calebtracey/config-yaml"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const configPath = "config.yaml"

type builder struct{}

func (b *builder) initializeDAO() (*internal.DAO, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

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
	defer panicQuit()

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

func panicQuit() {
	if r := recover(); r != nil {
		log.Errorf("I panicked and am quitting: %v", r)
		log.Error("I should be alerting someone...")
	}
}

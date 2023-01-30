package main

import (
	config "github.com/calebtracey/config-yaml"
	"github.com/joho/godotenv"
	"images-ai/internal"
	"log"
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
	if svc, err := new(builder).initializeDAO(); err != nil {
		log.Fatal(err)
	} else {
		run(svc)
	}
}

func run(svc *internal.DAO) {
	log.Fatal(internal.Handler{DAO: svc}.InitializeRoutes().Run())
}

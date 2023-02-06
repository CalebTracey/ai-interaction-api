package main

import (
	"github.com/calebtracey/ai-interaction-api/internal/facade"
	config "github.com/calebtracey/config-yaml"
	"github.com/joho/godotenv"
)

const OpenaiApi = "openAi"

func initializeDAO(appConfig *config.Config) (facade.Service, []error) {
	var errs []error
	if err := godotenv.Load(); err != nil {
		errs = append(errs, err)
	}

	openAiSvc, err := appConfig.Service(OpenaiApi)
	if err != nil {
		errs = append(errs, err)
	}

	return facade.Service{
		Client: openAiSvc.Client,
	}, errs
}

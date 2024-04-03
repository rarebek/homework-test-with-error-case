package test

import (
	"EXAM3/api-gateway/config"
	"EXAM3/api-gateway/services"
	"testing"
)

func RunApiTests(t *testing.T) {
	cfg := config.Load()
	service, err := services.NewServiceManager(&cfg)
}

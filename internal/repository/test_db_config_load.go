package repository

import (
	"chopper/internal/domain"
	"fmt"
	"net/url"
	"os"
)

func TestDbConfigLoad() (domain.DataBaseConfig, error) {
	var testDataBaseConfig domain.DataBaseConfig
	user := os.Getenv("TEST_DB_USER")
	password := os.Getenv("TEST_DB_PASSWORD")
	host := os.Getenv("TEST_DB_HOST")
	port := os.Getenv("TEST_DB_PORT")
	name := os.Getenv("TEST_DB_NAME")
	if user == "" || password == "" || host == "" || port == "" || name == "" {
		return domain.DataBaseConfig{}, fmt.Errorf("lack of environment variables")
	}
	testDataBaseConfig.User = user
	testDataBaseConfig.Password = url.QueryEscape(password)
	testDataBaseConfig.Host = host
	testDataBaseConfig.Port = port
	testDataBaseConfig.Name = name
	return testDataBaseConfig, nil
}

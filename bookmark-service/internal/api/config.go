package api

import (
	"log"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)


type Config struct {
	AppPort     	string 		`default:"8080" envconfig:"APP_PORT"`
	ServiceName 	string 		`default:"bookmark_service" envconfig:"SERVICE_NAME"`
	InstanceId  	string 		`default:"" envconfig:"INSTANCE_ID"`
	LogLevel 		string 		`default:"info" envconfig:"LOG_LEVEL"`
	BasePath 		string 		`default:"/" envconfig:"BASE_PATH"`
}

// NewConfig creates a new Config instance and loads the configuration from environment variables or a .env file.
func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{}
	
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	//thread-safety
	if cfg.InstanceId == "" {
		cfg.InstanceId = uuid.New().String()
	}else{
		//verify if the provided instance ID is a valid UUID
		if _, err := uuid.Parse(cfg.InstanceId); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
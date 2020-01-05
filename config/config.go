package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBEngine string
}

func CollectConfig() (config Config) {
	var missingEnv []string

	// DB_ENGINE
	config.DBEngine = os.Getenv("DB_ENGINE")
	if config.DBEngine == "" {
		missingEnv = append(missingEnv, "DB_ENGINE")
	}

	// Validation
	if len(missingEnv) > 0 {
		msg := fmt.Sprintf("Environment variables missing: %v", missingEnv)
		panic(fmt.Sprint(msg))
	}

	return
}

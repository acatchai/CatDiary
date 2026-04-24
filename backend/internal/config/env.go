package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if p := strings.TrimSpace(os.Getenv("CATDIARY_ENV_FILE")); p != "" {
		_ = godotenv.Load(p)
		return
	}
	_ = godotenv.Load(".env")
	_ = godotenv.Load(filepath.Join("..", ".env"))
}

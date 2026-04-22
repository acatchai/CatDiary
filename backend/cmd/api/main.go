package main

import (
	"log"
	"os"
	"strings"

	"github.com/acatchai/catdiary/backend/internal/repository"
	"github.com/acatchai/catdiary/backend/internal/router"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	dsn := strings.TrimSpace(os.Getenv("CATDIARY_MYSQL_DSN"))
	if dsn == "" {
		log.Fatal("CATDIARY_MYSQL_DSN is required")
	}

	repository.InitDB(dsn)

	h := server.New(server.WithHostPorts("0.0.0.0:8080"))
	router.Register(h)
	h.Spin()
}

package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/acatchai/catdiary/backend/internal/config"
	"github.com/acatchai/catdiary/backend/internal/repository"
	"github.com/acatchai/catdiary/backend/internal/router"
	"github.com/acatchai/catdiary/backend/internal/worker"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/go-playground/validator/v10"
)

func main() {

	config.LoadEnv()

	dsn := strings.TrimSpace(os.Getenv("CATDIARY_MYSQL_DSN"))
	if dsn == "" {
		log.Fatal("CATDIARY_MYSQL_DSN is required")
	}

	repository.InitDB(dsn)

	redisAddr := strings.TrimSpace(os.Getenv("CATDIARY_REDIS_ADDR"))
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}
	redisPassword := os.Getenv("CATDIARY_REDIS_PASSWORD")
	redisDB := 0
	if v := strings.TrimSpace(os.Getenv("CATDIARY_REDIS_DB")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			redisDB = n
		}
	}
	repository.InitRedis(redisAddr, redisPassword, redisDB)
	worker.StartDraftFlusher()

	vd := validator.New(validator.WithRequiredStructEnabled())

	h := server.New(server.WithHostPorts("0.0.0.0:8080"),
		server.WithCustomValidatorFunc(func(_ *protocol.Request, req any) error {
			return vd.Struct(req)
		}))

	router.Register(h)
	h.Spin()
}

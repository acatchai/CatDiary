package main

import (
	"log"
	"os"
	"strings"

	"github.com/acatchai/catdiary/backend/internal/config"
	"github.com/acatchai/catdiary/backend/internal/repository"
	"github.com/acatchai/catdiary/backend/internal/router"
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

	vd := validator.New(validator.WithRequiredStructEnabled())

	h := server.New(server.WithHostPorts("0.0.0.0:8080"),
		server.WithCustomValidatorFunc(func(_ *protocol.Request, req any) error {
			return vd.Struct(req)
		}))

	router.Register(h)
	h.Spin()
}

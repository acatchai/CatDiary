package main

import "github.com/acatchai/catdiary/backend/internal/repository"

func main() {
	dsn := "catdiary:123456@tcp(127.0.0.1:3306)/catdiary?charset=utf8mb4&parseTime=True&loc=Local"
	repository.InitDB(dsn)
}

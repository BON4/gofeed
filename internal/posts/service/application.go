package service

import (
	"github.com/BON4/gofeed/internal/posts/app"
	"github.com/BON4/gofeed/internal/posts/config"
)

func NewApplication(cfg config.ServerConfig) *app.Application {
	return &app.Application{}
}

package app

import (
	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/saintbyte/vk_group_getter/internal/config"
	"gorm.io/gorm"
)

type App struct {
	Config *config.Config
	VK     *api.VK
	DB     *gorm.DB
}

func New(cfg *config.Config, vk *api.VK, db *gorm.DB) *App {
	return &App{
		Config: cfg,
		VK:     vk,
		DB:     db,
	}
}

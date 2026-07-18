package main

import (
	"log"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/saintbyte/vk_group_getter/docs"
	"github.com/saintbyte/vk_group_getter/internal/app"
	"github.com/saintbyte/vk_group_getter/internal/config"
	"github.com/saintbyte/vk_group_getter/internal/model"
	"github.com/saintbyte/vk_group_getter/internal/processor"
	"github.com/saintbyte/vk_group_getter/internal/queue"
	"github.com/saintbyte/vk_group_getter/internal/server"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title           VK Group Getter API
// @version         1.0
// @description     API for processing VK group addresses.
// @BasePath        /
// @schemes         http
func main() {
	cfg := config.Load()

	host := cfg.Host
	if host == "" {
		host = "localhost"
	}
	docs.SwaggerInfo.Host = host + ":" + cfg.Port

	vk := api.NewVK(cfg.VKToken)

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	err = db.AutoMigrate(&model.Group{}, &model.VkPost{}, &model.VkPeople{}, &model.GroupPeople{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}
	application := app.New(cfg, vk, db)

	q := queue.New(100, processor.NewGroupHandler(application))
	q.Start()

	srv := server.New(application, q)
	log.Fatal(srv.Listen(cfg.Host + ":" + cfg.Port))
}

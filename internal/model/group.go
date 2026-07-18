package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Address           string `json:"address" gorm:"uniqueIndex"`
	ExternalID        string `json:"external_id" gorm:"column:external_id;type:text"`
	Type              string `json:"type" gorm:"column:type;type:text;not null;default:'vk_group'"`
	PostsCount        int    `json:"posts_count" gorm:"default:0"`
	FetchedPostsCount int    `json:"fetched_posts_count" gorm:"default:0"`
}

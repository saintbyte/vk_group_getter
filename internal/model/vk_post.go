package model

import (
	"time"

	"gorm.io/gorm"
)

type VkPost struct {
	gorm.Model
	GroupID        uint       `json:"group_id" gorm:"index"`
	VkID           int        `json:"vk_id" gorm:"index"`
	OwnerID        int        `json:"owner_id"`
	FromID         int        `json:"from_id"`
	Date           time.Time  `json:"date"`
	Text           string     `json:"text" gorm:"type:text"`
	PostType       string     `json:"post_type"`
	AccessKey      string     `json:"access_key"`
	ReplyOwnerID   int        `json:"reply_owner_id"`
	ReplyPostID    int        `json:"reply_post_id"`
	FriendsOnly    int        `json:"friends_only"`
	CommentsCount  int        `json:"comments_count"`
	LikesCount     int        `json:"likes_count"`
	RepostsCount   int        `json:"reposts_count"`
	ViewsCount     int        `json:"views_count"`
	SignerID       int        `json:"signer_id"`
	CanPin         int        `json:"can_pin"`
	CanDelete      int        `json:"can_delete"`
	CanEdit        int        `json:"can_edit"`
	IsPinned       int        `json:"is_pinned"`
	IsFavorite     int        `json:"is_favorite"`
	IsArchived     int        `json:"is_archived"`
	IsDeleted      int        `json:"is_deleted"`
	MarkedAsAds    int        `json:"marked_as_ads"`
	Edited         *time.Time `json:"edited,omitempty"`
	PostID         int        `json:"post_id"`
	PostponedID    int        `json:"postponed_id"`
	ShortTextRate  float64    `json:"short_text_rate"`
	CarouselOffset int        `json:"carousel_offset"`
	Hash           string     `json:"hash"`
	PostSource     string     `json:"post_source" gorm:"type:text"`
	Attachments    string     `json:"attachments" gorm:"type:text"`
	Geo            string     `json:"geo" gorm:"type:text"`
	CopyHistory    string     `json:"copy_history" gorm:"type:text"`
	Copyright      string     `json:"copyright" gorm:"type:text"`
	Donut          string     `json:"donut" gorm:"type:text"`
	Header         string     `json:"header" gorm:"type:text"`
	ParentsStack   string     `json:"parents_stack" gorm:"type:text"`
}

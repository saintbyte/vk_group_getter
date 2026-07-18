package model

import (
	"time"

	"gorm.io/gorm"
)

type GroupPeople struct {
	gorm.Model
	GroupID  uint       `json:"group_id" gorm:"index;uniqueIndex:idx_group_people"`
	PeopleID uint       `json:"people_id" gorm:"index;uniqueIndex:idx_group_people"`
	Status   string     `json:"status" gorm:"size:20;default:'active'"` // active, left, kicked
	BrokeAt  *time.Time `json:"broke_at"`
}

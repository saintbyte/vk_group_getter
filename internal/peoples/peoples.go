package peoples

import (
	"time"

	"github.com/saintbyte/vk_group_getter/internal/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) UpsertPeople(people *model.VkPeople) error {
	var existing model.VkPeople
	err := r.DB.Where("vk_id = ?", people.VkID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.DB.Create(people).Error
	}
	if err != nil {
		return err
	}
	people.ID = existing.ID
	return r.DB.Save(people).Error
}

func (r *Repository) UpsertGroupPeople(gp *model.GroupPeople) error {
	var existing model.GroupPeople
	err := r.DB.Where("group_id = ? AND people_id = ?", gp.GroupID, gp.PeopleID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.DB.Create(gp).Error
	}
	if err != nil {
		return err
	}
	gp.ID = existing.ID
	gp.CreatedAt = existing.CreatedAt
	return r.DB.Save(gp).Error
}

func (r *Repository) MarkInactiveMembers(groupID uint, activePeopleIDs []uint) error {
	if len(activePeopleIDs) == 0 {
		return nil
	}
	now := time.Now()
	return r.DB.Model(&model.GroupPeople{}).
		Where("group_id = ? AND status = ? AND people_id NOT IN ?", groupID, "active", activePeopleIDs).
		Updates(map[string]interface{}{
			"status":     "left",
			"broke_at":   now,
			"updated_at": now,
		}).Error
}

func (r *Repository) GetActivePeopleIDs(groupID uint) ([]uint, error) {
	var ids []uint
	err := r.DB.Model(&model.GroupPeople{}).
		Where("group_id = ? AND status = ?", groupID, "active").
		Pluck("people_id", &ids).Error
	return ids, err
}

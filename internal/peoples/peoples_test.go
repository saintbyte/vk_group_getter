package peoples

import (
	"testing"

	"github.com/saintbyte/vk_group_getter/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	err = db.AutoMigrate(&model.VkPeople{}, &model.GroupPeople{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestUpsertPeople_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)

	people := &model.VkPeople{VkID: 100, FirstName: "Test"}
	if err := repo.UpsertPeople(people); err != nil {
		t.Fatal(err)
	}

	if people.ID == 0 {
		t.Fatal("expected ID to be set after create")
	}

	var count int64
	db.Model(&model.VkPeople{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 record, got %d", count)
	}
}

func TestUpsertPeople_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)

	people := &model.VkPeople{VkID: 100, FirstName: "Old"}
	err := repo.UpsertPeople(people)
	if err != nil {
		t.Fatal(err)
	}
	people.FirstName = "New"
	if err := repo.UpsertPeople(people); err != nil {
		t.Fatal(err)
	}

	var saved model.VkPeople
	db.Where("vk_id = ?", 100).First(&saved)
	if saved.FirstName != "New" {
		t.Fatalf("expected 'New', got %q", saved.FirstName)
	}

	var count int64
	db.Model(&model.VkPeople{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 record after upsert, got %d", count)
	}
}

func TestUpsertGroupPeople_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)

	gp := &model.GroupPeople{GroupID: 1, PeopleID: 10, Status: "active"}
	if err := repo.UpsertGroupPeople(gp); err != nil {
		t.Fatal(err)
	}

	if gp.ID == 0 {
		t.Fatal("expected ID to be set after create")
	}
}

func TestUpsertGroupPeople_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)

	gp := &model.GroupPeople{GroupID: 1, PeopleID: 10, Status: "active"}
	err := repo.UpsertGroupPeople(gp)
	if err != nil {
		t.Fatal(err)
	}
	gp.Status = "left"
	if err := repo.UpsertGroupPeople(gp); err != nil {
		t.Fatal(err)
	}

	var saved model.GroupPeople
	db.Where("group_id = ? AND people_id = ?", 1, 10).First(&saved)
	if saved.Status != "left" {
		t.Fatalf("expected 'left', got %q", saved.Status)
	}

	var count int64
	db.Model(&model.GroupPeople{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 record after upsert, got %d", count)
	}
}

func TestMarkInactiveMembers(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)

	db.Create(&model.GroupPeople{GroupID: 1, PeopleID: 10, Status: "active"})
	db.Create(&model.GroupPeople{GroupID: 1, PeopleID: 20, Status: "active"})
	db.Create(&model.GroupPeople{GroupID: 1, PeopleID: 30, Status: "active"})

	if err := repo.MarkInactiveMembers(1, []uint{10, 20}); err != nil {
		t.Fatal(err)
	}

	var active []model.GroupPeople
	db.Where("group_id = ? AND status = ?", 1, "active").Find(&active)
	if len(active) != 2 {
		t.Fatalf("expected 2 active, got %d", len(active))
	}

	var left []model.GroupPeople
	db.Where("group_id = ? AND status = ?", 1, "left").Find(&left)
	if len(left) != 1 {
		t.Fatalf("expected 1 left, got %d", len(left))
	}
	if left[0].PeopleID != 30 {
		t.Fatalf("expected people_id 30 to be left, got %d", left[0].PeopleID)
	}
	if left[0].BrokeAt == nil {
		t.Fatal("expected broke_at to be set")
	}
}

func TestMarkInactiveMembers_EmptyList(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)

	db.Create(&model.GroupPeople{GroupID: 1, PeopleID: 10, Status: "active"})

	if err := repo.MarkInactiveMembers(1, []uint{}); err != nil {
		t.Fatal(err)
	}

	var count int64
	db.Model(&model.GroupPeople{}).Where("status = ?", "active").Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 active, got %d", count)
	}
}

func TestGetActivePeopleIDs(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)

	db.Create(&model.GroupPeople{GroupID: 1, PeopleID: 10, Status: "active"})
	db.Create(&model.GroupPeople{GroupID: 1, PeopleID: 20, Status: "active"})
	db.Create(&model.GroupPeople{GroupID: 1, PeopleID: 30, Status: "left"})
	db.Create(&model.GroupPeople{GroupID: 2, PeopleID: 10, Status: "active"})

	ids, err := repo.GetActivePeopleIDs(1)
	if err != nil {
		t.Fatal(err)
	}

	if len(ids) != 2 {
		t.Fatalf("expected 2 active IDs, got %d", len(ids))
	}

	idMap := map[uint]bool{}
	for _, id := range ids {
		idMap[id] = true
	}
	if !idMap[10] || !idMap[20] {
		t.Errorf("expected IDs 10 and 20, got %v", ids)
	}
}

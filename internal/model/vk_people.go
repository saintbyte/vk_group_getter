package model

import (
	"time"

	"gorm.io/gorm"
)

type VkPeople struct {
	gorm.Model
	VkID                   int       `json:"vk_id" gorm:"uniqueIndex"`
	FirstName              string    `json:"first_name"`
	LastName               string    `json:"last_name"`
	Deactivated            string    `json:"deactivated"`
	IsClosed               int       `json:"is_closed"`
	CanAccessClosed        int       `json:"can_access_closed"`
	Sex                    int       `json:"sex"`
	Bdate                  string    `json:"bdate"`
	City                   string    `json:"city" gorm:"type:text"`
	Country                string    `json:"country" gorm:"type:text"`
	Photo50                string    `json:"photo_50"`
	Photo100               string    `json:"photo_100"`
	Photo200               string    `json:"photo_200"`
	PhotoMax               string    `json:"photo_max"`
	Photo200Orig           string    `json:"photo_200_orig"`
	Photo400Orig           string    `json:"photo_400_orig"`
	PhotoMaxOrig           string    `json:"photo_max_orig"`
	PhotoID                string    `json:"photo_id"`
	Online                 int       `json:"online"`
	OnlineMobile           int       `json:"online_mobile"`
	Nickname               string    `json:"nickname"`
	Domain                 string    `json:"domain"`
	ScreenName             string    `json:"screen_name"`
	MaidenName             string    `json:"maiden_name"`
	HasMobile              int       `json:"has_mobile"`
	MobilePhone            string    `json:"mobile_phone"`
	HomePhone              string    `json:"home_phone"`
	Skype                  string    `json:"skype"`
	Site                   string    `json:"site"`
	Status                 string    `json:"status" gorm:"type:text"`
	StatusAudio            string    `json:"status_audio" gorm:"type:text"`
	LastSeen               string    `json:"last_seen" gorm:"type:text"`
	University             int       `json:"university"`
	UniversityName         string    `json:"university_name"`
	Faculty                int       `json:"faculty"`
	FacultyName            string    `json:"faculty_name"`
	Graduation             int       `json:"graduation"`
	EducationForm          string    `json:"education_form"`
	EducationStatus        string    `json:"education_status"`
	Universities           string    `json:"universities" gorm:"type:text"`
	Schools                string    `json:"schools" gorm:"type:text"`
	FriendStatus           int       `json:"friend_status"`
	CanPost                int       `json:"can_post"`
	CanSeeAllPosts         int       `json:"can_see_all_posts"`
	CanSeeAudio            int       `json:"can_see_audio"`
	CanWritePrivateMessage int       `json:"can_write_private_message"`
	CanSendFriendRequest   int       `json:"can_send_friend_request"`
	CanCallFromGroup       int       `json:"can_call_from_group"`
	CanBeInvitedGroup      int       `json:"can_be_invited_group"`
	CommonCount            int       `json:"common_count"`
	HasPhoto               int       `json:"has_photo"`
	IsFriend               int       `json:"is_friend"`
	IsFavorite             int       `json:"is_favorite"`
	IsHiddenFromFeed       int       `json:"is_hidden_from_feed"`
	Verified               int       `json:"verified"`
	Trending               int       `json:"trending"`
	Blacklisted            int       `json:"blacklisted"`
	BlacklistedByMe        int       `json:"blacklisted_by_me"`
	Relation               int       `json:"relation"`
	RelationPartner        string    `json:"relation_partner" gorm:"type:text"`
	Relatives              string    `json:"relatives" gorm:"type:text"`
	Counters               string    `json:"counters" gorm:"type:text"`
	HomeTown               string    `json:"home_town"`
	Personal               string    `json:"personal" gorm:"type:text"`
	Interests              string    `json:"interests" gorm:"type:text"`
	Music                  string    `json:"music" gorm:"type:text"`
	Activities             string    `json:"activities" gorm:"type:text"`
	Movies                 string    `json:"movies" gorm:"type:text"`
	Tv                     string    `json:"tv" gorm:"type:text"`
	Books                  string    `json:"books" gorm:"type:text"`
	Games                  string    `json:"games" gorm:"type:text"`
	About                  string    `json:"about" gorm:"type:text"`
	Quotes                 string    `json:"quotes" gorm:"type:text"`
	WallDefault            string    `json:"wall_default"`
	Timezone               int       `json:"timezone"`
	Occupation             string    `json:"occupation" gorm:"type:text"`
	Career                 string    `json:"career" gorm:"type:text"`
	Military               string    `json:"military" gorm:"type:text"`
	Exports                string    `json:"exports" gorm:"type:text"`
	CropPhoto              string    `json:"crop_photo" gorm:"type:text"`
	FollowersCount         int       `json:"followers_count"`
	OnlineInfo             string    `json:"online_info" gorm:"type:text"`
	TrackCode              string    `json:"track_code"`
	FoundWith              int       `json:"found_with"`
	Type                   string    `json:"type"`
	Lists                  string    `json:"lists" gorm:"type:text"`
	Connections            string    `json:"connections" gorm:"type:text"`
	ImageStatus            string    `json:"image_status" gorm:"type:text"`
	LastModified           time.Time `json:"last_modified"`
}

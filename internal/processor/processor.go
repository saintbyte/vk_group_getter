package processor

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/object"
	"github.com/saintbyte/vk_group_getter/internal/app"
	"github.com/saintbyte/vk_group_getter/internal/model"
	"github.com/saintbyte/vk_group_getter/internal/peoples"
)

func boolToInt(v object.BaseBoolInt) int {
	if v {
		return 1
	}
	return 0
}

func marshal(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func editedTime(ts int) *time.Time {
	if ts == 0 {
		return nil
	}
	t := time.Unix(int64(ts), 0)
	return &t
}

func sleepBetweenRequests() {
	time.Sleep(time.Duration(1120) * time.Millisecond)
}

func NewGroupHandler(a *app.App) func(id uint) {
	return func(id uint) {
		var group model.Group
		if err := a.DB.First(&group, id).Error; err != nil {
			log.Printf("queue: failed to get group %d: %v", id, err)
			return
		}
		log.Printf("queue: processing group %d (%s)", group.ID, group.Address)

		response, err := a.VK.UtilsResolveScreenName(api.Params{
			"screen_name": group.Address,
			"lang":        "ru",
		})
		if err != nil {
			log.Printf("queue: failed to resolve group %s: %v", group.Address, err)
			return
		}

		if response.ObjectID == 0 {
			log.Printf("queue: group %s not found", group.Address)
			return
		}

		if response.Type != "group" {
			log.Printf("queue: %s is not a group (type: %s)", group.Address, response.Type)
			return
		}

		group.ExternalID = strconv.Itoa(response.ObjectID)
		if err := a.DB.Save(&group).Error; err != nil {
			log.Printf("queue: failed to save group %d: %v", id, err)
			return
		}

		log.Printf("queue: group %s resolved to ID %s", group.Address, group.ExternalID)

		ownerID := -response.ObjectID
		offset := 0
		const count = 100
		const maxPosts = 1200
		totalFetched := 0
		totalPostsOnWall := 0

		for {
			resp, err := a.VK.WallGet(api.Params{
				"owner_id": ownerID,
				"offset":   offset,
				"count":    count,
				"filter":   "all",
				"lang":     "ru",
			})
			if err != nil {
				log.Printf("queue: failed to get wall posts for %s (offset %d): %v", group.Address, offset, err)
				return
			}

			for _, post := range resp.Items {
				if totalFetched >= maxPosts {
					break
				}
				vkPost := model.VkPost{
					GroupID:        group.ID,
					VkID:           post.ID,
					OwnerID:        post.OwnerID,
					FromID:         post.FromID,
					Date:           time.Unix(int64(post.Date), 0),
					Text:           post.Text,
					PostType:       post.PostType,
					AccessKey:      post.AccessKey,
					ReplyOwnerID:   post.ReplyOwnerID,
					ReplyPostID:    post.ReplyPostID,
					FriendsOnly:    post.FriendsOnly,
					CommentsCount:  post.Comments.Count,
					LikesCount:     post.Likes.Count,
					RepostsCount:   post.Reposts.Count,
					ViewsCount:     post.Views.Count,
					SignerID:       post.SignerID,
					CanPin:         boolToInt(post.CanPin),
					CanDelete:      boolToInt(post.CanDelete),
					CanEdit:        boolToInt(post.CanEdit),
					IsPinned:       boolToInt(post.IsPinned),
					IsFavorite:     boolToInt(post.IsFavorite),
					IsArchived:     boolToInt(post.IsArchived),
					IsDeleted:      boolToInt(post.IsDeleted),
					MarkedAsAds:    boolToInt(post.MarkedAsAds),
					Edited:         editedTime(post.Edited),
					PostID:         post.PostID,
					PostponedID:    post.PostponedID,
					ShortTextRate:  post.ShortTextRate,
					CarouselOffset: post.CarouselOffset,
					Hash:           post.Hash,
					PostSource:     marshal(post.PostSource),
					Attachments:    marshal(post.Attachments),
					Geo:            marshal(post.Geo),
					CopyHistory:    marshal(post.CopyHistory),
					Copyright:      marshal(post.Copyright),
					Donut:          marshal(post.Donut),
					Header:         marshal(post.Header),
					ParentsStack:   marshal(post.ParentsStack),
				}

				if err := a.DB.Create(&vkPost).Error; err != nil {
					log.Printf("queue: failed to save post %d for group %s: %v", post.ID, group.Address, err)
				}
				totalFetched++
			}

			log.Printf("queue: fetched %d posts (offset %d/%d) for %s", len(resp.Items), offset, resp.Count, group.Address)

			totalPostsOnWall = resp.Count

			offset += count
			if offset >= resp.Count || totalFetched >= maxPosts {
				break
			}

			sleepBetweenRequests()
		}

		log.Printf("queue: finished downloading posts for %s", group.Address)

		if group.PostsCount == 0 {
			group.PostsCount = totalPostsOnWall
		}
		group.FetchedPostsCount = totalFetched
		if err := a.DB.Save(&group).Error; err != nil {
			log.Printf("queue: failed to save group stats for %s: %v", group.Address, err)
		}

		log.Printf("queue: downloading members for %s", group.Address)
		fetchMembers(a, group)
	}
}

func fetchMembers(a *app.App, group model.Group) {
	repo := peoples.New(a.DB)
	groupID := group.ExternalID
	offset := 0
	const memberCount = 1000
	totalMembers := 0
	var activePeopleIDs []uint

	for {
		resp, err := a.VK.GroupsGetMembersFields(api.Params{
			"group_id": groupID,
			"fields":   "sex,bdate,city,country,photo_50,photo_100,photo_200_orig,photo_200,photo_400_orig,photo_max,photo_max_orig,online,online_mobile,lists,domain,has_mobile,contacts,connections,site,education,universities,schools,can_post,can_see_all_posts,can_see_audio,can_write_private_message,status,last_seen,common_count,relation,relatives,counters",
			"count":    memberCount,
			"offset":   offset,
			"sort":     "id_asc",
			"lang":     "ru",
		})
		if err != nil {
			log.Printf("queue: failed to get members for %s (offset %d): %v", group.Address, offset, err)
			return
		}

		for _, user := range resp.Items {
			people := vkUserToPeople(user)
			if err := repo.UpsertPeople(&people); err != nil {
				log.Printf("queue: failed to save people %d for group %s: %v", user.ID, group.Address, err)
				continue
			}
			activePeopleIDs = append(activePeopleIDs, people.ID)

			gp := model.GroupPeople{
				GroupID:  group.ID,
				PeopleID: people.ID,
				Status:   "active",
			}
			if err := repo.UpsertGroupPeople(&gp); err != nil {
				log.Printf("queue: failed to save group_people link for %d/%d: %v", group.ID, people.ID, err)
			}
		}

		totalMembers += len(resp.Items)
		log.Printf("queue: fetched %d members (offset %d/%d) for %s", len(resp.Items), offset, resp.Count, group.Address)

		offset += memberCount
		if offset >= resp.Count {
			break
		}
		sleepBetweenRequests()
	}

	if err := repo.MarkInactiveMembers(group.ID, activePeopleIDs); err != nil {
		log.Printf("queue: failed to mark inactive members for %s: %v", group.Address, err)
	}

	log.Printf("queue: finished downloading %d members for %s", totalMembers, group.Address)
}

func vkUserToPeople(user object.UsersUser) model.VkPeople {
	return model.VkPeople{
		VkID:                   user.ID,
		FirstName:              user.FirstName,
		LastName:               user.LastName,
		Deactivated:            user.Deactivated,
		IsClosed:               boolToInt(user.IsClosed),
		CanAccessClosed:        boolToInt(user.CanAccessClosed),
		Sex:                    user.Sex,
		Bdate:                  user.Bdate,
		City:                   marshal(user.City),
		Country:                marshal(user.Country),
		Photo50:                user.Photo50,
		Photo100:               user.Photo100,
		Photo200:               user.Photo200,
		PhotoMax:               user.PhotoMax,
		Photo200Orig:           user.Photo200Orig,
		Photo400Orig:           user.Photo400Orig,
		PhotoMaxOrig:           user.PhotoMaxOrig,
		PhotoID:                user.PhotoID,
		Online:                 boolToInt(user.Online),
		OnlineMobile:           boolToInt(user.OnlineMobile),
		Nickname:               user.Nickname,
		Domain:                 user.Domain,
		ScreenName:             user.ScreenName,
		MaidenName:             user.MaidenName,
		HasMobile:              boolToInt(user.HasMobile),
		MobilePhone:            user.MobilePhone,
		HomePhone:              user.HomePhone,
		Skype:                  user.Skype,
		Site:                   user.Site,
		Status:                 user.Status,
		StatusAudio:            marshal(user.StatusAudio),
		LastSeen:               marshal(user.LastSeen),
		University:             user.University,
		UniversityName:         user.UniversityName,
		Faculty:                user.Faculty,
		FacultyName:            user.FacultyName,
		Graduation:             user.Graduation,
		EducationForm:          user.EducationForm,
		EducationStatus:        user.EducationStatus,
		Universities:           marshal(user.Universities),
		Schools:                marshal(user.Schools),
		FriendStatus:           user.FriendStatus,
		CanPost:                boolToInt(user.CanPost),
		CanSeeAllPosts:         boolToInt(user.CanSeeAllPosts),
		CanSeeAudio:            boolToInt(user.CanSeeAudio),
		CanWritePrivateMessage: boolToInt(user.CanWritePrivateMessage),
		CanSendFriendRequest:   boolToInt(user.CanSendFriendRequest),
		CanCallFromGroup:       boolToInt(user.CanCallFromGroup),
		CanBeInvitedGroup:      boolToInt(user.CanBeInvitedGroup),
		CommonCount:            user.CommonCount,
		HasPhoto:               boolToInt(user.HasPhoto),
		IsFriend:               boolToInt(user.IsFriend),
		IsFavorite:             boolToInt(user.IsFavorite),
		IsHiddenFromFeed:       boolToInt(user.IsHiddenFromFeed),
		Verified:               boolToInt(user.Verified),
		Trending:               boolToInt(user.Trending),
		Blacklisted:            boolToInt(user.Blacklisted),
		BlacklistedByMe:        boolToInt(user.BlacklistedByMe),
		Relation:               user.Relation,
		RelationPartner:        marshal(user.RelationPartner),
		Relatives:              marshal(user.Relatives),
		Counters:               marshal(user.Counters),
		HomeTown:               user.HomeTown,
		Personal:               marshal(user.Personal),
		Interests:              user.Interests,
		Music:                  user.Music,
		Activities:             user.Activities,
		Movies:                 user.Movies,
		Tv:                     user.Tv,
		Books:                  user.Books,
		Games:                  user.Games,
		About:                  user.About,
		Quotes:                 user.Quotes,
		WallDefault:            user.WallDefault,
		Timezone:               user.Timezone,
		Occupation:             marshal(user.Occupation),
		Career:                 marshal(user.Career),
		Military:               marshal(user.Military),
		Exports:                marshal(user.Exports),
		CropPhoto:              marshal(user.CropPhoto),
		FollowersCount:         user.FollowersCount,
		OnlineInfo:             marshal(user.OnlineInfo),
		TrackCode:              user.TrackCode,
		FoundWith:              user.FoundWith,
		Type:                   user.Type,
		Lists:                  marshal(user.Lists),
		ImageStatus:            marshal(user.ImageStatus),
		LastModified:           time.Now(),
	}
}

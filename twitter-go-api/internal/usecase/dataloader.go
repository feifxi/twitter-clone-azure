package usecase

import (
	"context"
	"database/sql"

	"github.com/chanombude/twitter-go-api/internal/db"
)

// populateTweetItems acts as a simple DataLoader to batch fetch authors and parent/original tweets,
// resolving the N+1 query problem.
func (u *Usecase) populateTweetItems(ctx context.Context, tweets []db.Tweet, viewerID *int64) ([]TweetItem, error) {
	if len(tweets) == 0 {
		return []TweetItem{}, nil
	}

	var vID sql.NullInt64
	if viewerID != nil {
		vID = sql.NullInt64{Int64: *viewerID, Valid: true}
	}

	// 1. Collect unique IDs
	userIDsMap := make(map[int64]bool)
	tweetIDsMap := make(map[int64]bool)

	for _, t := range tweets {
		userIDsMap[t.UserID] = true
		if t.ParentID.Valid {
			tweetIDsMap[t.ParentID.Int64] = true
		}
		if t.RetweetID.Valid {
			tweetIDsMap[t.RetweetID.Int64] = true
		}
	}

	// Convert to slices
	userIDs := make([]int64, 0, len(userIDsMap))
	for id := range userIDsMap {
		userIDs = append(userIDs, id)
	}

	tweetIDs := make([]int64, 0, len(tweetIDsMap))
	for id := range tweetIDsMap {
		tweetIDs = append(tweetIDs, id)
	}

	// 2. Fetch Authors
	users, err := u.store.GetUsersByIDs(ctx, db.GetUsersByIDsParams{
		UserIds:  userIDs,
		ViewerID: vID,
	})
	if err != nil {
		return nil, err
	}

	usersMap := make(map[int64]UserItem)
	for _, rawUser := range users {
		usersMap[rawUser.ID] = UserItem{
			User: db.User{
				ID:             rawUser.ID,
				Username:       rawUser.Username,
				Email:          rawUser.Email,
				DisplayName:    rawUser.DisplayName,
				Bio:            rawUser.Bio,
				AvatarUrl:      rawUser.AvatarUrl,
				Role:           rawUser.Role,
				Provider:       rawUser.Provider,
				FollowersCount: rawUser.FollowersCount,
				FollowingCount: rawUser.FollowingCount,
				CreatedAt:      rawUser.CreatedAt,
				UpdatedAt:      rawUser.UpdatedAt,
			},
			IsFollowing: rawUser.IsFollowing,
		}
	}

	// 3. Fetch Referenced Tweets (Parents / Originals)
	var refTweetsMap map[int64]TweetItem
	if len(tweetIDs) > 0 {
		rawRefTweets, err := u.store.GetTweetsByIDs(ctx, db.GetTweetsByIDsParams{
			TweetIds: tweetIDs,
			ViewerID: vID,
		})
		if err != nil {
			return nil, err
		}

		// Because a referenced tweet also needs an author, we recursively call ourselves!
		// But to prevent infinite recursion, we assume ref tweets don't need *their* parents fully populated
		// down an infinite tree, just the immediate parent/original.
		// For a simple DataLoader, one level of depth is usually enough.
		refTweetsSlice := make([]db.Tweet, 0, len(rawRefTweets))
		for _, rt := range rawRefTweets {
			refTweetsSlice = append(refTweetsSlice, db.Tweet{
				ID:           rt.ID,
				UserID:       rt.UserID,
				Content:      rt.Content,
				MediaType:    rt.MediaType,
				MediaUrl:     rt.MediaUrl,
				ParentID:     rt.ParentID,
				RetweetID:    rt.RetweetID,
				ReplyCount:   rt.ReplyCount,
				RetweetCount: rt.RetweetCount,
				LikeCount:    rt.LikeCount,
				CreatedAt:    rt.CreatedAt,
				UpdatedAt:    rt.UpdatedAt,
			})
		}

		// To avoid deep recursion, we'll manually attach authors to the ref tweets here.
		// We need to fetch any missing authors first.
		missingAuthorsMap := make(map[int64]bool)
		for _, rt := range refTweetsSlice {
			if _, ok := usersMap[rt.UserID]; !ok {
				missingAuthorsMap[rt.UserID] = true
			}
		}

		if len(missingAuthorsMap) > 0 {
			missingAuthorIDs := make([]int64, 0, len(missingAuthorsMap))
			for id := range missingAuthorsMap {
				missingAuthorIDs = append(missingAuthorIDs, id)
			}
			moreUsers, err := u.store.GetUsersByIDs(ctx, db.GetUsersByIDsParams{
				UserIds:  missingAuthorIDs,
				ViewerID: vID,
			})
			if err == nil {
				for _, rawUser := range moreUsers {
					usersMap[rawUser.ID] = UserItem{
						User: db.User{
							ID:             rawUser.ID,
							Username:       rawUser.Username,
							Email:          rawUser.Email,
							DisplayName:    rawUser.DisplayName,
							Bio:            rawUser.Bio,
							AvatarUrl:      rawUser.AvatarUrl,
							Role:           rawUser.Role,
							Provider:       rawUser.Provider,
							FollowersCount: rawUser.FollowersCount,
							FollowingCount: rawUser.FollowingCount,
							CreatedAt:      rawUser.CreatedAt,
							UpdatedAt:      rawUser.UpdatedAt,
						},
						IsFollowing: rawUser.IsFollowing,
					}
				}
			}
		}

		refTweetsMap = make(map[int64]TweetItem)
		for i, rt := range refTweetsSlice {
			raw := rawRefTweets[i]
			item := TweetItem{
				Tweet:       rt,
				Author:      usersMap[rt.UserID],
				IsLiked:     raw.IsLiked,
				IsRetweeted: raw.IsRetweeted,
				IsFollowing: raw.IsFollowing,
			}
			refTweetsMap[rt.ID] = item
		}
	}

	// 4. Assemble final result
	result := make([]TweetItem, 0, len(tweets))

	// Fast lookup helper for the caller side mapper (since the caller passes db.Tweet but they are actually generic structs with is_liked etc embedded)
	// We will just map the base db.Tweet logic correctly. The caller needs to map IsLiked, IsRetweeted, IsFollowing manually if they didn't pass it into this function.
	// We'll update this function signature to take an interface or intermediate strict, but for now we'll assemble the base.

	for _, t := range tweets {
		item := TweetItem{
			Tweet:  t,
			Author: usersMap[t.UserID],
		}

		if t.ParentID.Valid {
			if parentTweet, ok := refTweetsMap[t.ParentID.Int64]; ok {
				item.ParentUsername = &parentTweet.Author.Username
			}
		}

		if t.RetweetID.Valid {
			if originalTweet, ok := refTweetsMap[t.RetweetID.Int64]; ok {
				item.OriginalTweet = &originalTweet
			}
		}

		result = append(result, item)
	}

	return result, nil
}

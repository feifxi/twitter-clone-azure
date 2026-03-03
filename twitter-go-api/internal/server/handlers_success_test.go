package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/chanombude/twitter-go-api/internal/config"
	"github.com/chanombude/twitter-go-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type mockTweetUC struct {
	createTweetFn func(ctx context.Context, input usecase.CreateTweetInput) (usecase.TweetItem, error)
}

func (m *mockTweetUC) CreateTweet(ctx context.Context, input usecase.CreateTweetInput) (usecase.TweetItem, error) {
	return m.createTweetFn(ctx, input)
}
func (m *mockTweetUC) DeleteTweet(context.Context, int64, int64) error { return nil }
func (m *mockTweetUC) GetTweet(context.Context, int64, *int64) (usecase.TweetItem, error) {
	return usecase.TweetItem{}, nil
}
func (m *mockTweetUC) ListReplies(context.Context, int64, int32, int32, *int64) ([]usecase.TweetItem, error) {
	return nil, nil
}
func (m *mockTweetUC) LikeTweet(context.Context, int64, int64) error   { return nil }
func (m *mockTweetUC) UnlikeTweet(context.Context, int64, int64) error { return nil }
func (m *mockTweetUC) Retweet(context.Context, int64, int64) (usecase.TweetItem, error) {
	return usecase.TweetItem{}, nil
}
func (m *mockTweetUC) UndoRetweet(context.Context, int64, int64) error { return nil }

type mockUserUC struct {
	updateProfileFn func(ctx context.Context, userID int64, input usecase.UpdateProfileInput) (usecase.UserItem, error)
	listFollowersFn func(ctx context.Context, targetUserID int64, page, size int32, viewerID *int64) ([]usecase.UserItem, error)
}

func (m *mockUserUC) GetUser(context.Context, int64, *int64) (usecase.UserItem, error) {
	return usecase.UserItem{}, nil
}
func (m *mockUserUC) UpdateProfile(ctx context.Context, userID int64, input usecase.UpdateProfileInput) (usecase.UserItem, error) {
	return m.updateProfileFn(ctx, userID, input)
}
func (m *mockUserUC) FollowUser(context.Context, int64, int64) (bool, error) { return true, nil }
func (m *mockUserUC) UnfollowUser(context.Context, int64, int64) error       { return nil }
func (m *mockUserUC) ListFollowers(ctx context.Context, targetUserID int64, page, size int32, viewerID *int64) ([]usecase.UserItem, error) {
	if m.listFollowersFn != nil {
		return m.listFollowersFn(ctx, targetUserID, page, size, viewerID)
	}
	return nil, nil
}
func (m *mockUserUC) ListFollowing(context.Context, int64, int32, int32, *int64) ([]usecase.UserItem, error) {
	return nil, nil
}

type mockFeedUC struct {
	getGlobalFeedFn func(ctx context.Context, page, size int32, viewerID *int64) ([]usecase.TweetItem, error)
}

func (m *mockFeedUC) GetGlobalFeed(ctx context.Context, page, size int32, viewerID *int64) ([]usecase.TweetItem, error) {
	return m.getGlobalFeedFn(ctx, page, size, viewerID)
}
func (m *mockFeedUC) GetFollowingFeed(context.Context, int64, int32, int32) ([]usecase.TweetItem, error) {
	return nil, nil
}
func (m *mockFeedUC) GetUserFeed(context.Context, int64, int32, int32, *int64) ([]usecase.TweetItem, error) {
	return nil, nil
}

type mockAuthUC struct {
	refreshSessionFn func(ctx context.Context, refreshToken string) (usecase.AuthResult, error)
	logoutFn         func(ctx context.Context, userID *int64, refreshToken *string)
}

func (m *mockAuthUC) LoginWithGoogle(context.Context, string) (usecase.AuthResult, error) {
	return usecase.AuthResult{}, nil
}
func (m *mockAuthUC) RefreshSession(ctx context.Context, refreshToken string) (usecase.AuthResult, error) {
	return m.refreshSessionFn(ctx, refreshToken)
}
func (m *mockAuthUC) Logout(ctx context.Context, userID *int64, refreshToken *string) {
	if m.logoutFn != nil {
		m.logoutFn(ctx, userID, refreshToken)
	}
}
func (m *mockAuthUC) GetMe(context.Context, int64) (usecase.UserItem, error) {
	return usecase.UserItem{}, nil
}

func TestCreateTweetSuccessUsesDetectedMediaType(t *testing.T) {
	t.Parallel()

	var gotType string
	var gotBytes []byte
	called := false
	mock := &mockTweetUC{
		createTweetFn: func(_ context.Context, input usecase.CreateTweetInput) (usecase.TweetItem, error) {
			called = true
			if input.Media == nil {
				t.Fatal("expected media input")
			}
			gotType = input.Media.ContentType
			b, err := io.ReadAll(input.Media.Reader)
			if err != nil {
				t.Fatalf("failed to read media reader: %v", err)
			}
			gotBytes = b
			return usecase.TweetItem{
				ID:        123,
				Content:   input.Content,
				CreatedAt: time.Now(),
				Author: usecase.UserItem{
					ID:       input.UserID,
					Username: "tester",
					Email:    "tester@example.com",
				},
			}, nil
		},
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("content", "hello")
	png := append([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, bytes.Repeat([]byte{0x00}, 16)...)
	part, _ := writer.CreateFormFile("media", "photo.png")
	_, _ = part.Write(png)
	_ = writer.Close()

	ctx, rec := newHandlerTestContext(http.MethodPost, "/api/v1/tweets", &body, writer.FormDataContentType())
	setAuthorizedUser(ctx, 7)

	s := &Server{
		config:  config.Config{MaxMediaBytes: 10 << 20},
		tweetUC: mock,
	}
	s.createTweet(ctx)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", rec.Code, rec.Body.String())
	}
	if !called {
		t.Fatal("expected createTweet usecase to be called")
	}
	if gotType != "image/png" {
		t.Fatalf("expected image/png, got %q", gotType)
	}
	if !bytes.Equal(gotBytes, png) {
		t.Fatal("media bytes changed before reaching usecase")
	}
	if !strings.Contains(rec.Body.String(), `"id":123`) {
		t.Fatalf("unexpected response body: %s", rec.Body.String())
	}
}

func TestUpdateProfileSuccessUsesDetectedAvatarType(t *testing.T) {
	t.Parallel()

	var gotType string
	called := false
	mock := &mockUserUC{
		updateProfileFn: func(_ context.Context, userID int64, input usecase.UpdateProfileInput) (usecase.UserItem, error) {
			called = true
			if userID != 5 {
				t.Fatalf("expected userID=5, got %d", userID)
			}
			if input.Avatar == nil {
				t.Fatal("expected avatar input")
			}
			gotType = input.Avatar.ContentType
			return usecase.UserItem{
				ID:       userID,
				Username: "tester",
				Email:    "tester@example.com",
			}, nil
		},
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("displayName", "new name")
	png := append([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, bytes.Repeat([]byte{0x00}, 8)...)
	part, _ := writer.CreateFormFile("avatar", "avatar.png")
	_, _ = part.Write(png)
	_ = writer.Close()

	ctx, rec := newHandlerTestContext(http.MethodPut, "/api/v1/users/profile", &body, writer.FormDataContentType())
	setAuthorizedUser(ctx, 5)

	s := &Server{
		config: config.Config{MaxAvatarBytes: 5 << 20},
		userUC: mock,
	}
	s.updateProfile(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rec.Code, rec.Body.String())
	}
	if !called {
		t.Fatal("expected updateProfile usecase to be called")
	}
	if gotType != "image/png" {
		t.Fatalf("expected image/png, got %q", gotType)
	}
}

func TestRefreshTokenSuccessSetsSessionCookies(t *testing.T) {
	t.Parallel()

	mock := &mockAuthUC{
		refreshSessionFn: func(_ context.Context, refreshToken string) (usecase.AuthResult, error) {
			if refreshToken != "r-old" {
				t.Fatalf("unexpected refresh token: %s", refreshToken)
			}
			return usecase.AuthResult{
				AccessToken:  "a-new",
				RefreshToken: "r-new",
				User: usecase.UserItem{
					ID:       9,
					Username: "tester",
					Email:    "tester@example.com",
				},
			}, nil
		},
	}

	ctx, rec := newHandlerTestContext(http.MethodPost, "/api/v1/auth/refresh", nil, "")
	ctx.Request.AddCookie(&http.Cookie{Name: "refresh_token", Value: "r-old"})

	s := &Server{
		config: config.Config{
			TokenDurationMinutes:     15,
			RefreshTokenDurationDays: 30,
			CookieSameSite:           "Lax",
		},
		authUC: mock,
	}
	s.refreshToken(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rec.Code, rec.Body.String())
	}

	cookies := rec.Result().Cookies()
	var accessSet, refreshSet bool
	for _, c := range cookies {
		if c.Name == "access_token" && c.Value == "a-new" {
			accessSet = true
		}
		if c.Name == "refresh_token" && c.Value == "r-new" {
			refreshSet = true
		}
	}
	if !accessSet || !refreshSet {
		t.Fatalf("expected both session cookies to be set, got: %+v", cookies)
	}
}

func TestLogoutWithAuthPayloadCallsUsecaseWithUserID(t *testing.T) {
	t.Parallel()

	called := false
	mock := &mockAuthUC{
		logoutFn: func(_ context.Context, userID *int64, refreshToken *string) {
			called = true
			if userID == nil || *userID != 77 {
				t.Fatalf("expected userID=77, got %+v", userID)
			}
			if refreshToken != nil {
				t.Fatalf("expected nil refresh token, got %v", *refreshToken)
			}
		},
	}

	ctx, rec := newHandlerTestContext(http.MethodPost, "/api/v1/auth/logout", nil, "")
	setAuthorizedUser(ctx, 77)

	s := &Server{
		config: config.Config{CookieSameSite: "Lax"},
		authUC: mock,
	}
	s.logout(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !called {
		t.Fatal("expected logout usecase to be called")
	}
}

func TestGetGlobalFeedReturnsItemsAndHasNext(t *testing.T) {
	t.Parallel()

	var gotSize int32
	mock := &mockFeedUC{
		getGlobalFeedFn: func(_ context.Context, _ int32, size int32, _ *int64) ([]usecase.TweetItem, error) {
			gotSize = size
			now := time.Now()
			return []usecase.TweetItem{
				{ID: 1, CreatedAt: now, Author: usecase.UserItem{ID: 11, Username: "u1", Email: "u1@example.com"}},
				{ID: 2, CreatedAt: now, Author: usecase.UserItem{ID: 12, Username: "u2", Email: "u2@example.com"}},
				{ID: 3, CreatedAt: now, Author: usecase.UserItem{ID: 13, Username: "u3", Email: "u3@example.com"}},
			}, nil
		},
	}

	ctx, rec := newHandlerTestContext(http.MethodGet, "/api/v1/feeds/global?size=2", nil, "")
	s := &Server{feedUC: mock}
	s.getGlobalFeed(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rec.Code, rec.Body.String())
	}
	if gotSize != 3 {
		t.Fatalf("expected usecase size=3 (size+1), got %d", gotSize)
	}

	var got struct {
		Items []struct {
			ID int64 `json:"id"`
		} `json:"items"`
		HasNext    bool    `json:"hasNext"`
		NextCursor *string `json:"nextCursor"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if len(got.Items) != 2 {
		t.Fatalf("expected 2 items after trim, got %d", len(got.Items))
	}
	if !got.HasNext {
		t.Fatal("expected hasNext=true")
	}
	if got.NextCursor == nil || *got.NextCursor == "" {
		t.Fatal("expected nextCursor to be set")
	}
	if got.Items[0].ID != 1 || got.Items[1].ID != 2 {
		t.Fatalf("unexpected trimmed items: %+v", got.Items)
	}
}

func TestListFollowersReturnsItemsAndHasNext(t *testing.T) {
	t.Parallel()

	var gotSize int32
	mock := &mockUserUC{
		listFollowersFn: func(_ context.Context, targetUserID int64, _ int32, size int32, _ *int64) ([]usecase.UserItem, error) {
			gotSize = size
			if targetUserID != 99 {
				t.Fatalf("expected targetUserID=99, got %d", targetUserID)
			}
			return []usecase.UserItem{
				{ID: 1, Username: "a", Email: "a@example.com"},
				{ID: 2, Username: "b", Email: "b@example.com"},
				{ID: 3, Username: "c", Email: "c@example.com"},
			}, nil
		},
	}

	ctx, rec := newHandlerTestContext(http.MethodGet, "/api/v1/users/99/followers?size=2", nil, "")
	ctx.Params = gin.Params{{Key: "id", Value: "99"}}
	s := &Server{userUC: mock}
	s.listFollowers(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rec.Code, rec.Body.String())
	}
	if gotSize != 3 {
		t.Fatalf("expected usecase size=3 (size+1), got %d", gotSize)
	}

	var got struct {
		Items []struct {
			ID int64 `json:"id"`
		} `json:"items"`
		HasNext    bool    `json:"hasNext"`
		NextCursor *string `json:"nextCursor"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if len(got.Items) != 2 {
		t.Fatalf("expected 2 items after trim, got %d", len(got.Items))
	}
	if !got.HasNext {
		t.Fatal("expected hasNext=true")
	}
	if got.NextCursor == nil || *got.NextCursor == "" {
		t.Fatal("expected nextCursor to be set")
	}
}

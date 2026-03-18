package engagement

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Mock implementations (shared across all _test.go files in this package)
// ---------------------------------------------------------------------------

type mockFavoriteRepo struct {
	addFn        func(ctx context.Context, userID, recipeID uuid.UUID) error
	removeFn     func(ctx context.Context, userID, recipeID uuid.UUID) error
	listByUserFn func(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]Favorite, int64, error)
	checkFn      func(ctx context.Context, userID uuid.UUID, recipeIDs []uuid.UUID) (map[uuid.UUID]bool, error)
}

func (m *mockFavoriteRepo) Add(ctx context.Context, userID, recipeID uuid.UUID) error {
	if m.addFn != nil {
		return m.addFn(ctx, userID, recipeID)
	}
	return nil
}

func (m *mockFavoriteRepo) Remove(ctx context.Context, userID, recipeID uuid.UUID) error {
	if m.removeFn != nil {
		return m.removeFn(ctx, userID, recipeID)
	}
	return nil
}

func (m *mockFavoriteRepo) ListByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]Favorite, int64, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID, page, pageSize)
	}
	return nil, 0, nil
}

func (m *mockFavoriteRepo) Check(ctx context.Context, userID uuid.UUID, recipeIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	if m.checkFn != nil {
		return m.checkFn(ctx, userID, recipeIDs)
	}
	return nil, nil
}

type mockViewRepo struct {
	recordFn func(ctx context.Context, view *ViewHistory) error
}

func (m *mockViewRepo) Record(ctx context.Context, view *ViewHistory) error {
	if m.recordFn != nil {
		return m.recordFn(ctx, view)
	}
	return nil
}

type mockRatingRepo struct {
	upsertFn      func(ctx context.Context, rating *Rating) error
	getByRecipeFn func(ctx context.Context, recipeID uuid.UUID) ([]Rating, error)
}

func (m *mockRatingRepo) Upsert(ctx context.Context, rating *Rating) error {
	if m.upsertFn != nil {
		return m.upsertFn(ctx, rating)
	}
	return nil
}

func (m *mockRatingRepo) GetByRecipe(ctx context.Context, recipeID uuid.UUID) ([]Rating, error) {
	if m.getByRecipeFn != nil {
		return m.getByRecipeFn(ctx, recipeID)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// Tests: EngagementService.AddFavorite
// ---------------------------------------------------------------------------

func TestEngagementService_AddFavorite_Success(t *testing.T) {
	userID := uuid.New()
	recipeID := uuid.New()

	svc := NewEngagementService(
		&mockFavoriteRepo{
			addFn: func(_ context.Context, u, r uuid.UUID) error {
				assert.Equal(t, userID, u)
				assert.Equal(t, recipeID, r)
				return nil
			},
		},
		&mockViewRepo{},
		&mockRatingRepo{},
	)

	resp, err := svc.AddFavorite(context.Background(), userID, recipeID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, recipeID.String(), resp.RecipeID)
	assert.False(t, resp.CreatedAt.IsZero())
}

func TestEngagementService_AddFavorite_RepoError(t *testing.T) {
	repoErr := errors.New("duplicate key")

	svc := NewEngagementService(
		&mockFavoriteRepo{
			addFn: func(_ context.Context, _, _ uuid.UUID) error {
				return repoErr
			},
		},
		&mockViewRepo{},
		&mockRatingRepo{},
	)

	resp, err := svc.AddFavorite(context.Background(), uuid.New(), uuid.New())

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Contains(t, err.Error(), "EngagementService.AddFavorite")
	assert.Nil(t, resp)
}

// ---------------------------------------------------------------------------
// Tests: EngagementService.RemoveFavorite
// ---------------------------------------------------------------------------

func TestEngagementService_RemoveFavorite_Success(t *testing.T) {
	userID := uuid.New()
	recipeID := uuid.New()

	svc := NewEngagementService(
		&mockFavoriteRepo{
			removeFn: func(_ context.Context, u, r uuid.UUID) error {
				assert.Equal(t, userID, u)
				assert.Equal(t, recipeID, r)
				return nil
			},
		},
		&mockViewRepo{},
		&mockRatingRepo{},
	)

	err := svc.RemoveFavorite(context.Background(), userID, recipeID)

	assert.NoError(t, err)
}

func TestEngagementService_RemoveFavorite_RepoError(t *testing.T) {
	repoErr := errors.New("not found")

	svc := NewEngagementService(
		&mockFavoriteRepo{
			removeFn: func(_ context.Context, _, _ uuid.UUID) error {
				return repoErr
			},
		},
		&mockViewRepo{},
		&mockRatingRepo{},
	)

	err := svc.RemoveFavorite(context.Background(), uuid.New(), uuid.New())

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Contains(t, err.Error(), "EngagementService.RemoveFavorite")
}

// ---------------------------------------------------------------------------
// Tests: EngagementService.ListFavorites
// ---------------------------------------------------------------------------

func TestEngagementService_ListFavorites_Success(t *testing.T) {
	userID := uuid.New()
	now := time.Now().UTC()
	favorites := []Favorite{
		{ID: uuid.New(), UserID: userID, RecipeID: uuid.New(), CreatedAt: now},
		{ID: uuid.New(), UserID: userID, RecipeID: uuid.New(), CreatedAt: now.Add(-time.Hour)},
	}

	svc := NewEngagementService(
		&mockFavoriteRepo{
			listByUserFn: func(_ context.Context, u uuid.UUID, page, pageSize int) ([]Favorite, int64, error) {
				assert.Equal(t, userID, u)
				assert.Equal(t, 1, page)
				assert.Equal(t, 20, pageSize)
				return favorites, 2, nil
			},
		},
		&mockViewRepo{},
		&mockRatingRepo{},
	)

	got, total, err := svc.ListFavorites(context.Background(), userID, 1, 20)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, got, 2)
	assert.Equal(t, favorites[0].RecipeID.String(), got[0].RecipeID)
	assert.Equal(t, favorites[1].RecipeID.String(), got[1].RecipeID)
}

func TestEngagementService_ListFavorites_Empty(t *testing.T) {
	svc := NewEngagementService(
		&mockFavoriteRepo{
			listByUserFn: func(_ context.Context, _ uuid.UUID, _, _ int) ([]Favorite, int64, error) {
				return []Favorite{}, 0, nil
			},
		},
		&mockViewRepo{},
		&mockRatingRepo{},
	)

	got, total, err := svc.ListFavorites(context.Background(), uuid.New(), 1, 20)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, got)
}

func TestEngagementService_ListFavorites_RepoError(t *testing.T) {
	repoErr := errors.New("db error")

	svc := NewEngagementService(
		&mockFavoriteRepo{
			listByUserFn: func(_ context.Context, _ uuid.UUID, _, _ int) ([]Favorite, int64, error) {
				return nil, 0, repoErr
			},
		},
		&mockViewRepo{},
		&mockRatingRepo{},
	)

	got, total, err := svc.ListFavorites(context.Background(), uuid.New(), 1, 20)

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Contains(t, err.Error(), "EngagementService.ListFavorites")
	assert.Nil(t, got)
	assert.Equal(t, int64(0), total)
}

// ---------------------------------------------------------------------------
// Tests: EngagementService.CheckFavorites
// ---------------------------------------------------------------------------

func TestEngagementService_CheckFavorites_Success(t *testing.T) {
	userID := uuid.New()
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	favoriteMap := map[uuid.UUID]bool{
		id1: true,
		id2: false,
		id3: true,
	}

	svc := NewEngagementService(
		&mockFavoriteRepo{
			checkFn: func(_ context.Context, u uuid.UUID, ids []uuid.UUID) (map[uuid.UUID]bool, error) {
				assert.Equal(t, userID, u)
				assert.Len(t, ids, 3)
				return favoriteMap, nil
			},
		},
		&mockViewRepo{},
		&mockRatingRepo{},
	)

	got, err := svc.CheckFavorites(context.Background(), userID, []uuid.UUID{id1, id2, id3})

	assert.NoError(t, err)
	assert.True(t, got[id1])
	assert.False(t, got[id2])
	assert.True(t, got[id3])
}

func TestEngagementService_CheckFavorites_RepoError(t *testing.T) {
	repoErr := errors.New("db error")

	svc := NewEngagementService(
		&mockFavoriteRepo{
			checkFn: func(_ context.Context, _ uuid.UUID, _ []uuid.UUID) (map[uuid.UUID]bool, error) {
				return nil, repoErr
			},
		},
		&mockViewRepo{},
		&mockRatingRepo{},
	)

	got, err := svc.CheckFavorites(context.Background(), uuid.New(), []uuid.UUID{uuid.New()})

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, got)
}

// ---------------------------------------------------------------------------
// Tests: EngagementService.RecordView
// ---------------------------------------------------------------------------

func TestEngagementService_RecordView_Success(t *testing.T) {
	userID := uuid.New()
	recipeID := uuid.New()
	sessionID := "session-abc-123"
	source := "search"

	var capturedView *ViewHistory
	svc := NewEngagementService(
		&mockFavoriteRepo{},
		&mockViewRepo{
			recordFn: func(_ context.Context, v *ViewHistory) error {
				capturedView = v
				return nil
			},
		},
		&mockRatingRepo{},
	)

	err := svc.RecordView(context.Background(), &userID, sessionID, recipeID, source)

	assert.NoError(t, err)
	assert.NotNil(t, capturedView)
	assert.Equal(t, &userID, capturedView.UserID)
	assert.Equal(t, sessionID, capturedView.SessionID)
	assert.Equal(t, recipeID, capturedView.RecipeID)
	assert.Equal(t, source, capturedView.Source)
	assert.False(t, capturedView.ViewedAt.IsZero())
	assert.NotEqual(t, uuid.Nil, capturedView.ID)
}

func TestEngagementService_RecordView_AnonymousUser(t *testing.T) {
	recipeID := uuid.New()

	var capturedView *ViewHistory
	svc := NewEngagementService(
		&mockFavoriteRepo{},
		&mockViewRepo{
			recordFn: func(_ context.Context, v *ViewHistory) error {
				capturedView = v
				return nil
			},
		},
		&mockRatingRepo{},
	)

	err := svc.RecordView(context.Background(), nil, "anon-session", recipeID, "direct")

	assert.NoError(t, err)
	assert.NotNil(t, capturedView)
	assert.Nil(t, capturedView.UserID)
	assert.Equal(t, recipeID, capturedView.RecipeID)
}

func TestEngagementService_RecordView_RepoError(t *testing.T) {
	repoErr := errors.New("write failed")

	svc := NewEngagementService(
		&mockFavoriteRepo{},
		&mockViewRepo{
			recordFn: func(_ context.Context, _ *ViewHistory) error {
				return repoErr
			},
		},
		&mockRatingRepo{},
	)

	err := svc.RecordView(context.Background(), nil, "session", uuid.New(), "search")

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Contains(t, err.Error(), "EngagementService.RecordView")
}

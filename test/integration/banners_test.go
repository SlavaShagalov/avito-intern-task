package integration

import (
	"context"
	pBanner "github.com/SlavaShagalov/avito-intern-task/internal/banner"
	pBannerRepo "github.com/SlavaShagalov/avito-intern-task/internal/banner/repository"
	bannerRepository "github.com/SlavaShagalov/avito-intern-task/internal/banner/repository/pgx"
	bannerUsecase "github.com/SlavaShagalov/avito-intern-task/internal/banner/usecase"
	"github.com/SlavaShagalov/avito-intern-task/internal/models"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/config"
	pErrors "github.com/SlavaShagalov/avito-intern-task/internal/pkg/errors"
	pLog "github.com/SlavaShagalov/avito-intern-task/internal/pkg/log/zap"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/storage/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"log"
	"testing"
)

var dbBanners = []models.Banner{
	{
		ID:        1,
		TagIDs:    []int64{1, 2, 3},
		FeatureID: 1,
		Content: map[string]any{
			"info":  "banner_1 info",
			"title": "banner_1",
		},
		IsActive: true,
	},
	{
		ID:        2,
		TagIDs:    []int64{4, 5},
		FeatureID: 2,
		Content: map[string]any{
			"info":  "banner_2 info",
			"title": "banner_2",
		},
		IsActive: true,
	},
	{
		ID:        3,
		TagIDs:    []int64{4},
		FeatureID: 1,
		Content: map[string]any{
			"info":  "banner_3 info",
			"title": "banner_3",
		},
		IsActive: false,
	},
}

type BannerSuite struct {
	suite.Suite
	pgxPool *pgxpool.Pool
	log     *zap.Logger
	uc      pBanner.Usecase
	ctx     context.Context
}

func (s *BannerSuite) SetupSuite() {
	s.ctx = context.Background()

	s.log = pLog.NewDev()

	config.SetTestPostgresConfig()
	config.SetTestRedisConfig()
	var err error
	s.pgxPool, err = postgres.NewPgx(s.log)
	s.Require().NoError(err)

	bannerRepo := bannerRepository.New(s.pgxPool, s.log)
	s.uc = bannerUsecase.New(bannerRepo, s.log)
}

func (s *BannerSuite) TearDownSuite() {
	s.pgxPool.Close()
	s.log.Info("Postgres connection closed")

	err := s.log.Sync()
	if err != nil {
		log.Println(err)
	}
}

func (s *BannerSuite) TestCreate() {
	type testCase struct {
		params *pBannerRepo.CreateParams
		err    error
	}

	tests := map[string]testCase{
		"normal": {
			params: &pBannerRepo.CreateParams{
				TagIDs:    []int64{1, 2, 3},
				FeatureID: 3,
				Content:   map[string]any{"title": "banner 5"},
				IsActive:  true,
			},
			err: nil,
		},
		"banner already exists": {
			params: &pBannerRepo.CreateParams{
				TagIDs:    []int64{2},
				FeatureID: 1,
				Content:   map[string]any{"title": "banner 5"},
				IsActive:  true,
			},
			err: pErrors.ErrBannerAlreadyExists,
		},
		"tag ids is nil": {
			params: &pBannerRepo.CreateParams{
				TagIDs:    nil,
				FeatureID: 1,
				Content:   map[string]any{},
				IsActive:  true,
			},
			err: pErrors.ErrBadTagIDsField,
		},
		"tag ids have zero id": {
			params: &pBannerRepo.CreateParams{
				TagIDs:    []int64{0, 1},
				FeatureID: 1,
				Content:   map[string]any{},
				IsActive:  true,
			},
			err: pErrors.ErrBadTagIDsField,
		},
		"tag ids have negative id": {
			params: &pBannerRepo.CreateParams{
				TagIDs:    []int64{1, -1},
				FeatureID: 1,
				Content:   map[string]any{},
				IsActive:  true,
			},
			err: pErrors.ErrBadTagIDsField,
		},
		"feature id is zero": {
			params: &pBannerRepo.CreateParams{
				TagIDs:    []int64{1},
				FeatureID: 0,
				Content:   map[string]any{},
				IsActive:  true,
			},
			err: pErrors.ErrBadFeatureIDField,
		},
		"feature id is negative": {
			params: &pBannerRepo.CreateParams{
				TagIDs:    []int64{1},
				FeatureID: 0,
				Content:   map[string]any{},
				IsActive:  true,
			},
			err: pErrors.ErrBadFeatureIDField,
		},
		"content is nil": {
			params: &pBannerRepo.CreateParams{
				TagIDs:    []int64{1},
				FeatureID: 1,
				Content:   nil,
				IsActive:  true,
			},
			err: pErrors.ErrBadContentField,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			bannerID, err := s.uc.Create(context.Background(), test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				// check banner in db
				banners, err := s.uc.List(context.Background(), &pBannerRepo.FilterParams{
					FeatureID: test.params.FeatureID,
					TagID:     test.params.TagIDs[0],
				})
				assert.NoError(s.T(), err, "failed to fetch banners from db")
				assert.Equal(s.T(), bannerID, banners[0].ID, "incorrect ID")
				assert.Equal(s.T(), test.params.FeatureID, banners[0].FeatureID, "incorrect FeatureID")
				assert.Equal(s.T(), test.params.TagIDs, banners[0].TagIDs, "incorrect TagIDs")
				assert.Equal(s.T(), test.params.Content, banners[0].Content, "incorrect Content")
				assert.Equal(s.T(), test.params.IsActive, banners[0].IsActive, "incorrect IsActive")

				// reset changes in db
				err = s.uc.Delete(context.Background(), bannerID)
				assert.NoError(s.T(), err, "failed to delete created banner")
			}
		})
	}
}

func (s *BannerSuite) TestList() {
	type testCase struct {
		params  *pBannerRepo.FilterParams
		banners []models.Banner
		err     error
	}

	tests := map[string]testCase{
		"no params": {
			params:  &pBannerRepo.FilterParams{},
			banners: dbBanners,
			err:     nil,
		},
		"filter by feature id": {
			params: &pBannerRepo.FilterParams{FeatureID: 1},
			banners: []models.Banner{
				dbBanners[0],
				dbBanners[2],
			},
			err: nil,
		},
		"filter by tag id": {
			params: &pBannerRepo.FilterParams{TagID: 4},
			banners: []models.Banner{
				dbBanners[1],
				dbBanners[2],
			},
			err: nil,
		},
		"limit": {
			params: &pBannerRepo.FilterParams{Limit: 2},
			banners: []models.Banner{
				dbBanners[0],
				dbBanners[1],
			},
			err: nil,
		},
		"offset": {
			params:  &pBannerRepo.FilterParams{Offset: 2},
			banners: []models.Banner{dbBanners[2]},
			err:     nil,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {

			banners, err := s.uc.List(context.Background(), test.params)

			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				assert.Equal(s.T(), len(test.banners), len(banners), "incorrect banners length")
				for i := 0; i < len(test.banners); i++ {
					assert.Equal(s.T(), test.banners[i].ID, banners[i].ID, "incorrect ID")
					assert.Equal(s.T(), test.banners[i].FeatureID, banners[i].FeatureID, "incorrect FeatureID")
					assert.Equal(s.T(), test.banners[i].TagIDs, banners[i].TagIDs, "incorrect TagIDs")
					assert.Equal(s.T(), test.banners[i].Content, banners[i].Content, "incorrect Content")
					assert.Equal(s.T(), test.banners[i].IsActive, banners[i].IsActive, "incorrect IsActive")
				}
			}
		})
	}
}

func (s *BannerSuite) TestGet() {
	type testCase struct {
		params  *pBannerRepo.GetParams
		content map[string]any
		err     error
	}

	tests := map[string]testCase{
		"normal": {
			params: &pBannerRepo.GetParams{
				FeatureID: 1,
				TagID:     1,
				IsAdmin:   true,
			},
			content: dbBanners[0].Content,
			err:     nil,
		},
		"banner not found": {
			params: &pBannerRepo.GetParams{
				FeatureID: 5,
				TagID:     1,
				IsAdmin:   true,
			},
			content: nil,
			err:     pErrors.ErrBannerNotFound,
		},
		"inactive banner for user": {
			params: &pBannerRepo.GetParams{
				FeatureID: 1,
				TagID:     4,
				IsAdmin:   false,
			},
			content: nil,
			err:     pErrors.ErrBannerDisabled,
		},
		"inactive banner for admin": {
			params: &pBannerRepo.GetParams{
				FeatureID: 1,
				TagID:     4,
				IsAdmin:   true,
			},
			content: dbBanners[2].Content,
			err:     nil,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			content, err := s.uc.Get(context.Background(), test.params)

			assert.ErrorIs(s.T(), err, test.err, "unexpected error")
			if err == nil {
				assert.Equal(s.T(), test.content, content, "incorrect content")
			}
		})
	}
}

func (s *BannerSuite) TestPartialUpdate() {
	type testCase struct {
		params *pBannerRepo.PartialUpdateParams
		banner *models.Banner
		err    error
	}

	fullUpdated := models.Banner{
		ID:        3,
		TagIDs:    []int64{5},
		FeatureID: 3,
		Content:   map[string]any{"info": "updated"},
		IsActive:  true,
	}

	resetParams := pBannerRepo.PartialUpdateParams{
		ID:        dbBanners[2].ID,
		TagIDs:    dbBanners[2].TagIDs,
		FeatureID: &dbBanners[2].FeatureID,
		Content:   dbBanners[2].Content,
		IsActive:  &dbBanners[2].IsActive,
	}

	tests := map[string]testCase{
		"full update": {
			params: &pBannerRepo.PartialUpdateParams{
				ID:        dbBanners[2].ID,
				TagIDs:    fullUpdated.TagIDs,
				FeatureID: &fullUpdated.FeatureID,
				Content:   fullUpdated.Content,
				IsActive:  &fullUpdated.IsActive,
			},
			banner: &fullUpdated,
			err:    nil,
		},
		"only tags update": {
			params: &pBannerRepo.PartialUpdateParams{
				ID:     dbBanners[2].ID,
				TagIDs: fullUpdated.TagIDs,
			},
			banner: &models.Banner{
				ID:        dbBanners[2].ID,
				TagIDs:    fullUpdated.TagIDs,
				FeatureID: dbBanners[2].FeatureID,
				Content:   dbBanners[2].Content,
				IsActive:  dbBanners[2].IsActive,
			},
			err: nil,
		},
		"only feature update": {
			params: &pBannerRepo.PartialUpdateParams{
				ID:        dbBanners[2].ID,
				FeatureID: &fullUpdated.FeatureID,
			},
			banner: &models.Banner{
				ID:        dbBanners[2].ID,
				TagIDs:    dbBanners[2].TagIDs,
				FeatureID: fullUpdated.FeatureID,
				Content:   dbBanners[2].Content,
				IsActive:  dbBanners[2].IsActive,
			},
			err: nil,
		},
		"only content update": {
			params: &pBannerRepo.PartialUpdateParams{
				ID:      dbBanners[2].ID,
				Content: fullUpdated.Content,
			},
			banner: &models.Banner{
				ID:        dbBanners[2].ID,
				TagIDs:    dbBanners[2].TagIDs,
				FeatureID: dbBanners[2].FeatureID,
				Content:   fullUpdated.Content,
				IsActive:  dbBanners[2].IsActive,
			},
			err: nil,
		},
		"only is_active update": {
			params: &pBannerRepo.PartialUpdateParams{
				ID:       dbBanners[2].ID,
				IsActive: &fullUpdated.IsActive,
			},
			banner: &models.Banner{
				ID:        dbBanners[2].ID,
				TagIDs:    dbBanners[2].TagIDs,
				FeatureID: dbBanners[2].FeatureID,
				Content:   dbBanners[2].Content,
				IsActive:  fullUpdated.IsActive,
			},
			err: nil,
		},
		"update feature and tags": {
			params: &pBannerRepo.PartialUpdateParams{
				ID:        dbBanners[2].ID,
				TagIDs:    fullUpdated.TagIDs,
				FeatureID: &fullUpdated.FeatureID,
			},
			banner: &models.Banner{
				ID:        dbBanners[2].ID,
				TagIDs:    fullUpdated.TagIDs,
				FeatureID: fullUpdated.FeatureID,
				Content:   dbBanners[2].Content,
				IsActive:  dbBanners[2].IsActive,
			},
			err: nil,
		},
		"banner with such feature and tag already exists": {
			params: &pBannerRepo.PartialUpdateParams{
				ID:     dbBanners[2].ID,
				TagIDs: []int64{5},
				FeatureID: func() *int64 {
					var id int64 = 2
					return &id
				}(),
			},
			banner: nil,
			err:    pErrors.ErrBannerAlreadyExists,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			err := s.uc.PartialUpdate(context.Background(), test.params)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if err == nil {
				// check banner in db
				banners, err := s.uc.List(context.Background(), &pBannerRepo.FilterParams{
					FeatureID: test.banner.FeatureID,
					TagID:     test.banner.TagIDs[0],
				})
				assert.NoError(s.T(), err, "failed to fetch banners from db")
				assert.Equal(s.T(), test.banner.ID, banners[0].ID, "incorrect ID")
				assert.Equal(s.T(), test.banner.FeatureID, banners[0].FeatureID, "incorrect FeatureID")
				assert.Equal(s.T(), test.banner.TagIDs, banners[0].TagIDs, "incorrect TagIDs")
				assert.Equal(s.T(), test.banner.Content, banners[0].Content, "incorrect Content")
				assert.Equal(s.T(), test.banner.IsActive, banners[0].IsActive, "incorrect IsActive")

				// reset banner
				err = s.uc.PartialUpdate(context.Background(), &resetParams)
				assert.NoError(s.T(), err, "failed to reset banner in db")
			}
		})
	}
}

func (s *BannerSuite) TestDelete() {
	type testCase struct {
		featureID   int64
		tagID       int64
		setupBanner func() (int64, error)
		err         error
	}

	tests := map[string]testCase{
		"normal": {
			featureID: 5,
			tagID:     1,
			setupBanner: func() (int64, error) {
				return s.uc.Create(context.Background(), &pBannerRepo.CreateParams{
					TagIDs:    []int64{1, 2, 3},
					FeatureID: 5,
					Content:   map[string]any{"title": "tmp banner"},
					IsActive:  true,
				})
			},
			err: nil,
		},
		"banner not found": {
			setupBanner: func() (int64, error) {
				return 999, nil
			},
			err: pErrors.ErrBannerNotFound,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			id, err := test.setupBanner()
			s.Require().NoError(err)

			err = s.uc.Delete(context.Background(), id)
			assert.ErrorIs(s.T(), err, test.err, "unexpected error")

			if test.err == nil {
				_, err = s.uc.Get(context.Background(), &pBannerRepo.GetParams{
					FeatureID: test.featureID,
					TagID:     test.tagID,
					IsAdmin:   true,
				})
				assert.ErrorIs(s.T(), err, pErrors.ErrBannerNotFound, "banner should be deleted")
			}
		})
	}
}

func TestBannerSuite(t *testing.T) {
	suite.Run(t, new(BannerSuite))
}

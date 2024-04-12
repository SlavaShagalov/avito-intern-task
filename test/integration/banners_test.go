package integration

import (
	"context"
	pBanner "github.com/SlavaShagalov/avito-intern-task/internal/banner"
	pBannerRepo "github.com/SlavaShagalov/avito-intern-task/internal/banner/repository"
	bannerRepository "github.com/SlavaShagalov/avito-intern-task/internal/banner/repository/v3/pgx"
	bannerUsecase "github.com/SlavaShagalov/avito-intern-task/internal/banner/usecase"
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
				TagIDs:    []int{1, 2, 3},
				FeatureID: 3,
				Content:   map[string]any{"title": "banner 5"},
				IsActive:  true,
			},
			err: nil,
		},
		"banner already exists": {
			params: &pBannerRepo.CreateParams{
				TagIDs:    []int{2},
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
				TagIDs:    []int{0, 1},
				FeatureID: 1,
				Content:   map[string]any{},
				IsActive:  true,
			},
			err: pErrors.ErrBadTagIDsField,
		},
		"tag ids have negative id": {
			params: &pBannerRepo.CreateParams{
				TagIDs:    []int{1, -1},
				FeatureID: 1,
				Content:   map[string]any{},
				IsActive:  true,
			},
			err: pErrors.ErrBadTagIDsField,
		},
		"feature id is zero": {
			params: &pBannerRepo.CreateParams{
				TagIDs:    []int{},
				FeatureID: 0,
				Content:   map[string]any{},
				IsActive:  true,
			},
			err: pErrors.ErrBadFeatureIDField,
		},
		"feature id is negative": {
			params: &pBannerRepo.CreateParams{
				TagIDs:    []int{},
				FeatureID: 0,
				Content:   map[string]any{},
				IsActive:  true,
			},
			err: pErrors.ErrBadFeatureIDField,
		},
		"content is nil": {
			params: &pBannerRepo.CreateParams{
				TagIDs:    []int{},
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
				_ = bannerID
				//assert.Equal(s.T(), test.params., board.WorkspaceID, "incorrect WorkspaceID")
				//assert.Equal(s.T(), test.params.Title, board.Title, "incorrect Title")
				//assert.Equal(s.T(), test.params.Description, board.Description, "incorrect Description")
				//
				//getBoard, err := s.uc.Get(ctx, board.ID)
				//assert.NoError(s.T(), err, "failed to fetch board from the database")
				//assert.Equal(s.T(), board.ID, getBoard.ID, "incorrect boardID")
				//assert.Equal(s.T(), test.params.WorkspaceID, getBoard.WorkspaceID, "incorrect WorkspaceID")
				//assert.Equal(s.T(), test.params.Title, getBoard.Title, "incorrect Title")
				//assert.Equal(s.T(), test.params.Description, getBoard.Description, "incorrect Description")
				//
				//err = s.uc.Delete(ctx, board.ID)
				//assert.NoError(s.T(), err, "failed to delete created board")
			}
		})
	}
}

func TestBoardSuite(t *testing.T) {
	suite.Run(t, new(BannerSuite))
}

package std

import (
	"context"
	"fmt"
	"github.com/SlavaShagalov/avito-intern-task/internal/models"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/avito-intern-task/internal/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"

	pBannerRepo "github.com/SlavaShagalov/avito-intern-task/internal/banner/repository"
)

type repository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func New(pool *pgxpool.Pool, log *zap.Logger) pBannerRepo.Repository {
	return &repository{
		pool: pool,
		log:  log,
	}
}

const createBannerCmd = `
	INSERT INTO banners (feature_id, content, is_active)
	VALUES ($1, $2, $3)
	RETURNING id;`

func (r *repository) Create(ctx context.Context, params *pBannerRepo.CreateParams) (int64, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.log.Error(constants.DBError, zap.Error(err))
		return 0, pErrors.ErrDb
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, createBannerCmd,
		params.FeatureID,
		params.Content,
		params.IsActive,
	)
	var bannerID int64
	err = row.Scan(&bannerID)
	if err != nil {
		r.log.Error(constants.DBError, zap.Error(err))
		return 0, pErrors.ErrDb
	}

	var valueStrings []string
	var args []any
	for _, tagID := range params.TagIDs {
		valueString := fmt.Sprintf("($%d, $%d)", len(args)+1, len(args)+2)
		args = append(args, bannerID, tagID)
		valueStrings = append(valueStrings, valueString)
	}

	createBannerTagsCmd := fmt.Sprintf(`
		INSERT INTO banner_tags (banner_id, tag_id)
		VALUES %s;`,
		strings.Join(valueStrings, ", "))

	r.log.Debug("", zap.String("", createBannerTagsCmd))

	_, err = tx.Exec(ctx, createBannerTagsCmd, args...)
	if err != nil {
		r.log.Error(constants.DBError, zap.Error(err))
		return 0, pErrors.ErrDb
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.log.Error(constants.DBError, zap.Error(err))
		return 0, pErrors.ErrDb
	}

	r.log.Debug("Banner created", zap.Int64("banner_id", bannerID))
	return bannerID, nil
}

func (r *repository) List(ctx context.Context, params *pBannerRepo.FilterParams) ([]models.Banner, error) {
	listCmd := `
	SELECT b.id,
		   ARRAY_AGG(bt.tag_id) AS tag_ids,
		   b.feature_id,
		   b.content,
		   b.is_active,
		   b.created_at,
		   b.updated_at
	FROM banners b
			 JOIN
		 banner_tags bt ON b.id = bt.banner_id`

	conditions := []string{}
	args := []any{}
	if params.FeatureID > 0 {
		conditions = append(conditions, fmt.Sprintf("b.feature_id = $%d", len(args)+1))
		args = append(args, params.FeatureID)
	}
	if params.TagID > 0 {
		conditions = append(conditions, fmt.Sprintf("bt.tag_id = $%d", len(args)+1))
		args = append(args, params.TagID)
	}

	if len(conditions) > 0 {
		listCmd += " WHERE " + fmt.Sprintf("%s", conditions[0])
		for _, condition := range conditions[1:] {
			listCmd += " AND " + fmt.Sprintf("%s", condition)
		}
	}

	listCmd += " GROUP BY b.id"
	if params.Limit > 0 {
		listCmd += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, params.Limit)
	}
	if params.Offset > 0 {
		listCmd += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, params.Offset)
	}

	r.log.Debug("Cmd", zap.String("", listCmd))
	r.log.Debug("Args", zap.Any("", args))

	rows, err := r.pool.Query(ctx, listCmd, args...)
	if err != nil {
		r.log.Error(constants.DBError, zap.Error(err))
		return nil, pErrors.ErrDb
	}
	defer rows.Close()

	banners := []models.Banner{}
	var banner models.Banner
	for rows.Next() {
		err = rows.Scan(
			&banner.ID,
			&banner.TagIDs,
			&banner.FeatureID,
			&banner.Content,
			&banner.IsActive,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)
		if err != nil {
			r.log.Error(constants.DBError, zap.Error(err))
			return nil, pErrors.ErrDb
		}

		banners = append(banners, banner)
	}

	return banners, nil
}

const getCmd = `
	SELECT b.content
	FROM banners b
			 JOIN banner_tags bt ON b.id = bt.banner_id
	WHERE bt.tag_id = $1 AND b.feature_id = $2
	GROUP BY b.id;`

func (r *repository) Get(ctx context.Context, params *pBannerRepo.GetParams) (map[string]any, error) {
	row := r.pool.QueryRow(ctx, getCmd, params.TagID, params.FeatureID)

	var content map[string]any
	err := row.Scan(&content)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pErrors.ErrBannerNotFound
		}

		r.log.Error(constants.DBError, zap.Error(err))
		return nil, pErrors.ErrDb
	}

	return content, nil
}

const delTagsCmd = `
	DELETE FROM banner_tags
	WHERE banner_id = $1;`

func (r *repository) PartialUpdate(ctx context.Context, params *pBannerRepo.PartialUpdateParams) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.log.Error(constants.DBError, zap.Error(err))
		return pErrors.ErrDb
	}
	defer tx.Rollback(ctx)

	var setValues []string
	var args []any
	if params.FeatureID != nil {
		setValue := fmt.Sprintf("feature_id = $%d", len(args)+1)
		args = append(args, *params.FeatureID)
		setValues = append(setValues, setValue)
	}
	if params.Content != nil {
		setValue := fmt.Sprintf("content = $%d", len(args)+1)
		args = append(args, params.Content)
		setValues = append(setValues, setValue)
	}
	if params.IsActive != nil {
		setValue := fmt.Sprintf("is_active = $%d", len(args)+1)
		args = append(args, *params.IsActive)
		setValues = append(setValues, setValue)
	}

	if len(setValues) != 0 {
		partialUpdateCmd := "UPDATE banners SET " + setValues[0]
		for _, setValue := range setValues[1:] {
			partialUpdateCmd += ", " + setValue
		}
		partialUpdateCmd += fmt.Sprintf(", updated_at = now() WHERE id = $%d;", len(args)+1)
		args = append(args, params.ID)

		r.log.Debug("partialUpdateCmd", zap.String("", partialUpdateCmd))
		r.log.Debug("args", zap.Any("", args))

		result, err := r.pool.Exec(ctx, partialUpdateCmd, args...)
		if err != nil {
			r.log.Error(constants.DBError, zap.Error(err))
			return pErrors.ErrDb
		}

		rowsAffected := result.RowsAffected()
		if rowsAffected == 0 {
			return pErrors.ErrBannerNotFound
		}
	}

	if params.TagIDs != nil {
		_, err = r.pool.Exec(ctx, delTagsCmd, params.ID)
		if err != nil {
			r.log.Error(constants.DBError, zap.Error(err))
			return pErrors.ErrDb
		}

		var valueStrings []string
		var insArgs []any
		for _, tagID := range params.TagIDs {
			valueString := fmt.Sprintf("($%d, $%d)", len(insArgs)+1, len(insArgs)+2)
			insArgs = append(insArgs, params.ID, tagID)
			valueStrings = append(valueStrings, valueString)
		}

		createBannerTagsCmd := fmt.Sprintf(`
		INSERT INTO banner_tags (banner_id, tag_id)
		VALUES %s;`,
			strings.Join(valueStrings, ", "))

		r.log.Debug("createBannerTagsCmd", zap.String("", createBannerTagsCmd))
		r.log.Debug("insArgs", zap.Any("", insArgs))

		_, err = tx.Exec(ctx, createBannerTagsCmd, insArgs...)
		if err != nil {
			r.log.Error(constants.DBError, zap.Error(err))
			return pErrors.ErrDb
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.log.Error(constants.DBError, zap.Error(err))
		return pErrors.ErrDb
	}

	r.log.Debug("Banner updated", zap.Int64("banner_id", params.ID))
	return nil
}

const deleteCmd = `
	DELETE FROM banners
	WHERE id = $1;`

func (r *repository) Delete(ctx context.Context, id int64) error {
	result, err := r.pool.Exec(ctx, deleteCmd, id)
	if err != nil {
		r.log.Error(constants.DBError, zap.Error(err))
		return errors.Wrap(pErrors.ErrDb, err.Error())
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return pErrors.ErrBannerNotFound
	}

	r.log.Debug("Banner deleted", zap.Int64("banner_id", id))
	return nil
}

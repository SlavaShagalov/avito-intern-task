package std

import (
	"context"
	"fmt"
	"github.com/SlavaShagalov/avito-intern-task/internal/models"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/avito-intern-task/internal/pkg/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
INSERT INTO banners (content, is_active)
VALUES ($1, $2)
RETURNING id;`

const createBannerReferencesCmd = `
INSERT INTO banner_references(banner_id, feature_id, tag_id)
VALUES %s;`

func (r *repository) Create(ctx context.Context, params *pBannerRepo.CreateParams) (int64, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.log.Error(constants.DBError, zap.Error(err))
		return 0, pErrors.ErrDb
	}
	defer tx.Rollback(ctx) // nolint

	row := tx.QueryRow(ctx, createBannerCmd,
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
		valueString := fmt.Sprintf("($%d, $%d, $%d)", len(args)+1, len(args)+2, len(args)+3)
		args = append(args, bannerID, params.FeatureID, tagID)
		valueStrings = append(valueStrings, valueString)
	}

	cmd := fmt.Sprintf(createBannerReferencesCmd, strings.Join(valueStrings, ", "))
	_, err = tx.Exec(ctx, cmd, args...)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		if pgErr.Code == pgerrcode.UniqueViolation {
			return 0, pErrors.ErrBannerAlreadyExists
		}
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

const listCmd = `
SELECT b.id,
       ARRAY_AGG(br.tag_id) AS tag_ids,
       br.feature_id,
       b.content,
       b.is_active,
       b.created_at,
       b.updated_at
FROM banners b
         JOIN
     banner_references br ON b.id = br.banner_id
%s
GROUP BY b.id, br.feature_id
ORDER BY b.id
%s;`

func (r *repository) List(ctx context.Context, params *pBannerRepo.FilterParams) ([]models.Banner, error) {
	conditions := make([]string, 0, 2)
	args := make([]any, 0, 4)
	if params.FeatureID > 0 {
		conditions = append(conditions, fmt.Sprintf("feature_id = $%d", len(args)+1))
		args = append(args, params.FeatureID)
	}
	if params.TagID > 0 {
		conditions = append(conditions, fmt.Sprintf("tag_id = $%d", len(args)+1))
		args = append(args, params.TagID)
	}

	var conditionPart string
	if len(conditions) > 0 {
		condition := strings.Join(conditions, " AND ")
		conditionPart = fmt.Sprintf(`
WHERE b.id IN (SELECT banner_id
               FROM banner_references
               WHERE %s)`, condition)
	}

	var limitPart string
	if params.Limit > 0 {
		limitPart = fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, params.Limit)
	}
	if params.Offset > 0 {
		limitPart += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, params.Offset)
	}

	cmd := fmt.Sprintf(listCmd, conditionPart, limitPart)

	r.log.Debug("Cmd", zap.String("", cmd))
	r.log.Debug("Args", zap.Any("", args))

	rows, err := r.pool.Query(ctx, cmd, args...)
	if err != nil {
		r.log.Error(constants.DBError, zap.Error(err))
		return nil, pErrors.ErrDb
	}
	defer rows.Close()

	banners := make([]models.Banner, 0, 4)
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
SELECT b.id,
       ARRAY_AGG(br.tag_id) AS tag_ids,
       br.feature_id,
       b.content,
       b.is_active,
       b.created_at,
       b.updated_at
FROM banners b
         JOIN banner_references br ON b.id = br.banner_id
WHERE b.id = (SELECT banner_id
              FROM banner_references
              WHERE tag_id = $1
                AND feature_id = $2)
GROUP BY b.id, br.feature_id;`

func (r *repository) Get(ctx context.Context, params *pBannerRepo.GetParams) (*models.Banner, error) {
	row := r.pool.QueryRow(ctx, getCmd, params.TagID, params.FeatureID)

	banner := new(models.Banner)
	err := row.Scan(
		&banner.ID,
		&banner.TagIDs,
		&banner.FeatureID,
		&banner.Content,
		&banner.IsActive,
		&banner.CreatedAt,
		&banner.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pErrors.ErrBannerNotFound
		}
		r.log.Error(constants.DBError, zap.Error(err))
		return nil, pErrors.ErrDb
	}

	return banner, nil
}

const updateBannerCmd = `
UPDATE banners
SET %s,
    updated_at = now()
WHERE id = $%d;`

const delBannerReferencesCmd = `
DELETE
FROM banner_references
WHERE banner_id = $1
RETURNING feature_id;`

const updateBannerFeatureCmd = `
UPDATE banner_references
SET feature_id = $2
WHERE banner_id = $1;`

func (r *repository) PartialUpdate(ctx context.Context, params *pBannerRepo.PartialUpdateParams) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.log.Error(constants.DBError, zap.Error(err))
		return pErrors.ErrDb
	}
	defer tx.Rollback(ctx) // nolint

	setValues := make([]string, 0, 2)
	args := make([]any, 0, 3)
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
	if len(setValues) > 0 {
		setValuesPart := strings.Join(setValues, ", ")
		cmd := fmt.Sprintf(updateBannerCmd, setValuesPart, len(args)+1)
		args = append(args, params.ID)

		r.log.Debug("partialUpdateCmd", zap.String("", cmd))
		r.log.Debug("args", zap.Any("", args))

		res, err := tx.Exec(ctx, cmd, args...)
		if err != nil {
			r.log.Error(constants.DBError, zap.Error(err))
			return pErrors.ErrDb
		}
		if res.RowsAffected() == 0 {
			return pErrors.ErrBannerNotFound
		}
	}

	if params.TagIDs != nil {
		row := tx.QueryRow(ctx, delBannerReferencesCmd, params.ID)
		var featureID int64
		if err = row.Scan(&featureID); err != nil {
			r.log.Error(constants.DBError, zap.Error(err))
			return pErrors.ErrDb
		}
		if params.FeatureID != nil {
			featureID = *params.FeatureID
		}

		valueStrings := make([]string, 0, 3)
		insArgs := make([]any, 0, 3)
		for _, tagID := range params.TagIDs {
			valueString := fmt.Sprintf("($%d, $%d, $%d)", len(insArgs)+1, len(insArgs)+2, len(insArgs)+3)
			valueStrings = append(valueStrings, valueString)
			insArgs = append(insArgs, params.ID, featureID, tagID)
		}

		cmd := fmt.Sprintf(createBannerReferencesCmd, strings.Join(valueStrings, ", "))

		r.log.Debug("createBannerTagsCmd", zap.String("", cmd))
		r.log.Debug("insArgs", zap.Any("", insArgs))

		res, err := tx.Exec(ctx, cmd, insArgs...)
		if err != nil {
			pgErr := err.(*pgconn.PgError)
			if pgErr.Code == pgerrcode.UniqueViolation {
				return pErrors.ErrBannerAlreadyExists
			}
			r.log.Error(constants.DBError, zap.Error(err))
			return pErrors.ErrDb
		}
		if res.RowsAffected() == 0 {
			return pErrors.ErrBannerNotFound
		}
	} else if params.FeatureID != nil {
		res, err := tx.Exec(ctx, updateBannerFeatureCmd, params.ID, *params.FeatureID)
		if err != nil {
			pgErr := err.(*pgconn.PgError)
			if pgErr.Code == pgerrcode.UniqueViolation {
				return pErrors.ErrBannerAlreadyExists
			}
			r.log.Error(constants.DBError, zap.Error(err))
			return pErrors.ErrDb
		}
		if res.RowsAffected() == 0 {
			return pErrors.ErrBannerNotFound
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
	res, err := r.pool.Exec(ctx, deleteCmd, id)
	if err != nil {
		r.log.Error(constants.DBError, zap.Error(err))
		return pErrors.ErrDb
	}
	if res.RowsAffected() == 0 {
		return pErrors.ErrBannerNotFound
	}
	r.log.Debug("Banner deleted", zap.Int64("banner_id", id))
	return nil
}

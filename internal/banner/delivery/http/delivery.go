package http

import (
	"encoding/json"
	pBannerRepo "github.com/SlavaShagalov/avito-intern-task/internal/banner/repository"
	mw "github.com/SlavaShagalov/avito-intern-task/internal/middleware"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/avito-intern-task/internal/pkg/errors"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	pBanner "github.com/SlavaShagalov/avito-intern-task/internal/banner"

	pHTTP "github.com/SlavaShagalov/avito-intern-task/internal/pkg/http"
)

const (
	FeatureIDKey = "feature_id"
	TagIDKey     = "tag_id"
	LimitKey     = "limit"
	OffsetKey    = "offset"
)

type delivery struct {
	uc  pBanner.Usecase
	log *zap.Logger
}

func RegisterHandlers(mux *mux.Router, uc pBanner.Usecase, log *zap.Logger, checkAuth mw.Middleware, adminAccess mw.Middleware) {
	dlv := delivery{
		uc:  uc,
		log: log,
	}

	const (
		bannersPath    = constants.ApiPrefix + "/banner"
		bannerPath     = bannersPath + "/{id}"
		userBannerPath = constants.ApiPrefix + "/user_banner"
	)

	mux.HandleFunc(bannersPath, checkAuth(adminAccess(dlv.create))).Methods(http.MethodPost)
	mux.HandleFunc(bannersPath, checkAuth(adminAccess(dlv.list))).Methods(http.MethodGet)
	mux.HandleFunc(userBannerPath, checkAuth(dlv.get)).Methods(http.MethodGet)
	mux.HandleFunc(bannerPath, checkAuth(adminAccess(dlv.partialUpdate))).Methods(http.MethodPatch)
	mux.HandleFunc(bannerPath, checkAuth(adminAccess(dlv.delete))).Methods(http.MethodDelete)
}

func (d *delivery) create(w http.ResponseWriter, r *http.Request) {
	body, err := pHTTP.ReadBody(r, d.log)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	var request createRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	params := pBannerRepo.CreateParams{
		TagIDs:    request.TagIDs,
		FeatureID: request.FeatureID,
		Content:   request.Content,
		IsActive:  request.IsActive,
	}

	bannerID, err := d.uc.Create(r.Context(), &params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newCreateResponse(bannerID)
	pHTTP.SendJSON(w, r, http.StatusCreated, response)
}

func (d *delivery) list(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	featureID, err := strconv.Atoi(queryParams.Get(FeatureIDKey))
	if queryParams.Get(FeatureIDKey) != "" && err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrBadFeatureIDParam)
		return
	}
	tagID, err := strconv.Atoi(queryParams.Get(TagIDKey))
	if queryParams.Get(TagIDKey) != "" && err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrBadTagIDParam)
		return
	}
	limit, err := strconv.Atoi(queryParams.Get(LimitKey))
	if queryParams.Get(LimitKey) != "" && err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrBadLimitParam)
		return
	}
	offset, err := strconv.Atoi(queryParams.Get(OffsetKey))
	if queryParams.Get(OffsetKey) != "" && err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrBadOffsetParam)
		return
	}

	params := pBannerRepo.FilterParams{
		FeatureID: featureID,
		TagID:     tagID,
		Limit:     limit,
		Offset:    offset,
	}

	banners, err := d.uc.List(r.Context(), &params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	response := newListResponse(banners)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

func (d *delivery) get(w http.ResponseWriter, r *http.Request) {
	isAdmin, ok := r.Context().Value(mw.ContextIsAdmin).(bool)
	if !ok {
		d.log.Error("is_admin field not found")
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	queryParams := r.URL.Query()
	tagID, err := strconv.Atoi(queryParams.Get(TagIDKey))
	if err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrBadTagIDParam)
		return
	}
	featureID, err := strconv.Atoi(queryParams.Get(FeatureIDKey))
	if err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrBadFeatureIDParam)
		return
	}

	params := pBannerRepo.GetParams{
		FeatureID: featureID,
		TagID:     tagID,
		IsAdmin:   isAdmin,
	}

	content, err := d.uc.Get(r.Context(), &params)
	if err != nil {
		if errors.Is(err, pErrors.ErrBannerDisabled) {
			w.WriteHeader(http.StatusForbidden)
		} else {
			pHTTP.HandleError(w, r, err)
		}
		return
	}

	pHTTP.SendJSON(w, r, http.StatusOK, content)
}

func (d *delivery) partialUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bannerID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	body, err := pHTTP.ReadBody(r, d.log)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	var request partialUpdateRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		d.log.Error(constants.FailedReadRequestBody, zap.Error(err))
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	params := pBannerRepo.PartialUpdateParams{
		ID:        bannerID,
		TagIDs:    request.TagIDs,
		FeatureID: request.FeatureID,
		Content:   request.Content,
		IsActive:  request.IsActive,
	}

	err = d.uc.PartialUpdate(r.Context(), &params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (d *delivery) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bannerID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	err = d.uc.Delete(r.Context(), bannerID)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

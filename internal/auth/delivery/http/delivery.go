package http

import (
	"encoding/json"
	"github.com/SlavaShagalov/avito-intern-task/internal/auth"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/avito-intern-task/internal/pkg/errors"
	pHTTP "github.com/SlavaShagalov/avito-intern-task/internal/pkg/http"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

const (
	authPrefix = "/auth"
	signInPath = constants.ApiPrefix + authPrefix + "/signin"
	signUpPath = constants.ApiPrefix + authPrefix + "/signup"
)

type delivery struct {
	uc  auth.Usecase
	log *zap.Logger
}

func RegisterHandlers(mux *mux.Router, uc auth.Usecase, log *zap.Logger) {
	del := delivery{
		uc:  uc,
		log: log,
	}

	mux.HandleFunc(signUpPath, del.signup).Methods(http.MethodPost)
	mux.HandleFunc(signInPath, del.signin).Methods(http.MethodPost)
}

// signup godoc
//
//	@Summary		Creates new user and returns authentication cookie.
//	@Description	Creates new user and returns authentication cookie.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			signUpParams	body		SignUpRequest	true	"Sign up params."
//	@Success		200				{object}	SignUpResponse	"Successfully created user."
//	@Failure		400				{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/auth/signup [post]
func (d *delivery) signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := pHTTP.ReadBody(r, d.log)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	var request SignUpRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		pHTTP.HandleError(w, r, pErrors.ErrReadBody)
		return
	}

	params := auth.SignUpParams{
		Username: request.Username,
		Password: request.Password,
	}

	user, authToken, err := d.uc.SignUp(ctx, &params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}
	w.Header().Add("token", authToken)

	response := newSignUpResponse(user)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

// signin godoc
//
//	@Summary		Logs in and returns the authentication cookie
//	@Description	Logs in and returns the authentication cookie
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			signInParams	body		SignInRequest	true	"Successfully authenticated."
//	@Success		200				{object}	SignInResponse	"successfully auth"
//	@Failure		400				{object}	http.JSONError
//	@Failure		404				{object}	http.JSONError
//	@Failure		405
//	@Failure		500
//	@Router			/auth/signin [post]
func (d *delivery) signin(w http.ResponseWriter, r *http.Request) {
	body, err := pHTTP.ReadBody(r, d.log)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	var request SignInRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}

	params := auth.SignInParams{
		Username: request.Username,
		Password: request.Password,
	}

	user, authToken, err := d.uc.SignIn(r.Context(), &params)
	if err != nil {
		pHTTP.HandleError(w, r, err)
		return
	}
	w.Header().Add("token", authToken)

	response := newSignInResponse(user)
	pHTTP.SendJSON(w, r, http.StatusOK, response)
}

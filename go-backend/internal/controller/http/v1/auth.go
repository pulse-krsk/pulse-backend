package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	cuserr "github.com/kurochkinivan/pulskrsk/internal/customErrors"
	"github.com/kurochkinivan/pulskrsk/internal/usecase"
)

type authHandler struct {
	auth       usecase.Auth
	bytesLimit int64
	signingkey string
}

func NewAuthHandler(auth usecase.Auth, bytesLimit int64, signingkey string) *authHandler {
	return &authHandler{
		auth:       auth,
		bytesLimit: bytesLimit,
		signingkey: signingkey,
	}
}

const basePath string = "/api/v1"
const baseAuthPath string = basePath + "/auth"

func (h *authHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("%s %s/login", http.MethodPost, baseAuthPath), errMdw(logMdw(h.loginUser)))
	mux.HandleFunc(fmt.Sprintf("%s %s/refresh", http.MethodGet, baseAuthPath), errMdw(logMdw(h.refreshTokens)))
	mux.HandleFunc(fmt.Sprintf("%s %s/logout", http.MethodPost, baseAuthPath), errMdw(logMdw(h.logoutUser)))
}

type (
	loginRequest struct {
		OauthToken string `json:"token"`
	}

	loginResponse struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Avatar    string `json:"avatar"`
	}
)

func (h *authHandler) loginUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	const op string = "authHandler.loginUser"

	reqData, err := io.ReadAll(io.LimitReader(r.Body, h.bytesLimit))
	if err != nil {
		return cuserr.ErrReadRequestBody.WithErr(fmt.Errorf("%s: %w", op, err))
	}
	defer r.Body.Close()

	var req loginRequest
	err = json.Unmarshal(reqData, &req)
	if err != nil {
		return cuserr.ErrSerializeData.WithErr(fmt.Errorf("%s: %w", op, err))
	}

	if req.OauthToken == "" {
		return cuserr.ErrNotAllFieldsProvided
	}

	accessToken, refreshToken, user, err := h.auth.LoginUser(r.Context(), req.OauthToken)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	loginResp := loginResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
	}
	err = json.NewEncoder(w).Encode(loginResp)
	if err != nil {
		return cuserr.ErrSerializeData.WithErr(fmt.Errorf("%s: %w", op, err))
	}

	return nil
}

func (h *authHandler) refreshTokens(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	const op string = "authHandler.refreshTokens"

	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		return cuserr.ErrGetCookie.WithErr(fmt.Errorf("%s: %w", op, err))
	}

	accessToken, newRefreshToken, err := h.auth.RefreshTokens(r.Context(), refreshToken.Value)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: newRefreshToken, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode})
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: accessToken, Path: "/", SameSite: http.SameSiteLaxMode})

	return nil
}

func (h *authHandler) logoutUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	const op string = "authHandler.logoutUser"

	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		return cuserr.ErrGetCookie.WithErr(fmt.Errorf("%s: %w", op, err))
	}

	err = h.auth.LogoutUser(r.Context(), refreshToken.Value)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

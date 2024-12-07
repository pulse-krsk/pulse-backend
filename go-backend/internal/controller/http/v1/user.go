package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	cuserr "github.com/kurochkinivan/pulskrsk/internal/customErrors"
	"github.com/kurochkinivan/pulskrsk/internal/usecase"
)

type userHandler struct {
	bytesLimit  int64
	signingKey  string
	userUseCase usecase.User
}

func NewUserHandler(userUseCase usecase.User, bytesLimit int64, signingKey string) *userHandler {
	return &userHandler{
		userUseCase: userUseCase,
		bytesLimit:  bytesLimit,
		signingKey:  signingKey,
	}
}

func (h *userHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("%s %s/users/event-types", http.MethodPost, basePath), errMdw(authMdw(logMdw(h.addFavoriteEventTypes), h.signingKey)))
	mux.HandleFunc(fmt.Sprintf("%s %s/users/profile", http.MethodGet, basePath), errMdw(authMdw(logMdw(h.getUserWithTypes), h.signingKey)))
}

type (
	addFavoriteEventTypesRequest struct {
		EventTypes []string `json:"types"`
	}
)

func (h *userHandler) addFavoriteEventTypes(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	const op string = "userHandler.addFavoriteEventTypes"

	reqData, err := io.ReadAll(io.LimitReader(r.Body, h.bytesLimit))
	if err != nil {
		return cuserr.ErrReadRequestBody.WithErr(fmt.Errorf("%s: %w", op, err))
	}
	defer r.Body.Close()

	var req addFavoriteEventTypesRequest
	err = json.Unmarshal(reqData, &req)
	if err != nil {
		return cuserr.ErrSerializeData.WithErr(fmt.Errorf("%s: %w", op, err))
	}

	err = h.userUseCase.AddFavoriteEventTypes(r.Context(), r.Header.Get("user_id"), req.EventTypes)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (h *userHandler) getUserWithTypes(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	const op string = "userHandler.getUserWithTypes"

	user, err := h.userUseCase.GetUserWithTypes(r.Context(), r.Header.Get("user_id"))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return cuserr.ErrSerializeData.WithErr(fmt.Errorf("%s: %w", op, err))
	}

	return nil
}

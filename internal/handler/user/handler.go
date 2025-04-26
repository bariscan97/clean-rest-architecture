package user

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"github.com/bariscan97/clean-rest-architecture/pkg/token"
	repo "github.com/bariscan97/clean-rest-architecture/internal/repository/user"
	"github.com/bariscan97/clean-rest-architecture/internal/utils"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type authKey token.AuthKey 

type Handler struct {
	repository repo.IUserRepository
	TokenMaker *token.JWTMaker
}

func NewUserHandler(repository repo.IUserRepository, secretKey string) *Handler {
	return &Handler{
		repository: repository,
		TokenMaker: token.NewJWTMaker(secretKey),
	}
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}
	user, err := h.repository.GetUserByIdentifier(r.Context(), id.String())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toUserRes(user))

}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u RegisterUserReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	hashed, err := utils.HashPassword(u.Password)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}
	u.Password = hashed

	created, err := h.repository.CreateUser(r.Context(), CreateReqToDomain(u))
	if err != nil {
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}

	res := toUserRes(created)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	pageStr, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limitStr, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	users, err := h.repository.ListUsers(r.Context(), pageStr, limitStr)

	if err != nil {
		http.Error(w, "error listing users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ListUserRes(users))
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var u UpdateUserReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(authKey{}).(*token.UserClaims)

	if u.Password != "" {
		hashed, err := utils.HashPassword(u.Password)
		if err != nil {
			http.Error(w, "error hashing password", http.StatusInternalServerError)
			return
		}
		u.Password = hashed
	}
	
	currentUserID := claims.ID
	

	if err := h.repository.UpdateUserByID(r.Context(), currentUserID, utils.StructToMap(u)); err != nil {
		http.Error(w, "error updating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	currentUserID := r.Context().Value(authKey{}).(*token.UserClaims).ID

	if err := h.repository.DeleteUserByID(r.Context(), currentUserID); err != nil {
		http.Error(w, "error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var u LoginUserReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}
	user, err := h.repository.GetUserByIdentifier(r.Context(), u.Identifier)

	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}

	if err = utils.CheckPassword(u.Password, user.Password); err != nil {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}

	accessToken, accessClaims, err := h.TokenMaker.CreateToken(user.ID, user.UserName, user.Email, 15*time.Minute)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	res := LoginUserRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
		User:                 toUserRes(user),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

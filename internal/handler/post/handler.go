package post

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/bariscan97/clean-rest-architecture/pkg/token"
	repo "github.com/bariscan97/clean-rest-architecture/internal/repository/post"
	"github.com/bariscan97/clean-rest-architecture/internal/utils"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type authKey token.AuthKey 

type Handler struct {
	repository repo.IPostRepository
}

func NewPostHandler(repository repo.IPostRepository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (h *Handler) GetCommentByPostID(w http.ResponseWriter, r *http.Request) {
	parentID := chi.URLParam(r, "id")
	var parsedParentID *uuid.UUID
	if parentID != "" {
		uid, err := uuid.Parse(parentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}
		parsedParentID = &uid
	}

	pageStr, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limitStr, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	posts, err := h.repository.ListPosts(r.Context(), nil ,parsedParentID, pageStr, limitStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ListPostRes(posts))
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	var parseduserID *uuid.UUID
	if userID != "" {
		uid, err := uuid.Parse(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}
		parseduserID = &uid
	}

	pageStr, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limitStr, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	posts, err := h.repository.ListPosts(r.Context(), parseduserID, nil, pageStr, limitStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ListPostRes(posts))
}

func (h *Handler) DeletePostByID(w http.ResponseWriter, r *http.Request) {
	postID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}
	currentUserID := r.Context().Value(authKey{}).(*token.UserClaims).ID

	if err := h.repository.DeletePostByID(r.Context(), currentUserID, postID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
	
func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	postID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	var p UpdatePostReq
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}
	
	currentUserID := r.Context().Value(authKey{}).(*token.UserClaims).ID

	if err := h.repository.UpdatePost(r.Context(), postID, currentUserID, utils.StructToMap(p)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	parentID := r.URL.Query().Get("parent_id")
	var parsedParentID *uuid.UUID
	if parentID != "" {
		uid, err := uuid.Parse(parentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}
		parsedParentID = &uid
	}
	var p CreatePostReq
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	
	currendUserID := r.Context().Value(authKey{}).(*token.UserClaims).ID

	created, err := h.repository.CreatePost(r.Context(), parsedParentID, currendUserID, CreateReqToDomain(p))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	res := toCreatePostRes(created)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

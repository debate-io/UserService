package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/debate-io/service-auth/internal/interface/server/middleware"
	"github.com/debate-io/service-auth/internal/registry"
	"github.com/go-chi/chi"
	"github.com/ztrue/tracerr"
	"go.uber.org/zap"
)

type RestHandler struct {
	logger   *zap.Logger
	usecases *registry.UseCases
}

type Url string

const (
	ImageUrl Url = "/user/{id}/image"
	PingUrl  Url = "/ping"
)

func NewRestHandler(
	logger *zap.Logger,
	usecases *registry.UseCases,
	isDebug bool,
) *RestHandler {
	return &RestHandler{
		logger:   logger,
		usecases: usecases,
	}
}

func (h *RestHandler) GetImageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "invalid url path", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(middleware.JwtClaimsKey).(*model.Claims)
	if claims == nil {
		http.Error(w, repo.ErrUnauthorized.Unwrap().Error(), http.StatusUnauthorized)
		return
	}

	role := claims.Role
	if role != model.RoleAdmin && role != model.RoleContentManager && role != model.RoleDefaultUser {
		http.Error(w, repo.ErrUnauthorized.Unwrap().Error(), http.StatusUnauthorized)
		return
	}

	image, contentType, err := h.usecases.Users.DownloadImage(r.Context(), int(id))
	if err != nil {
		fmt.Println(err)
		http.Error(w, tracerr.Unwrap(err).Error(), http.StatusNoContent)
		return
	}

	w.Header().Add("Content-Type", fmt.Sprintf("image/%s", contentType))
	w.WriteHeader(http.StatusOK)
	w.Write(image)
}

func (h *RestHandler) PutImageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "invalid url path", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(middleware.JwtClaimsKey).(*model.Claims)
	if claims == nil {
		http.Error(w, repo.ErrUnauthorized.Unwrap().Error(), http.StatusUnauthorized)
		return
	}

	role, claimsId := claims.Role, claims.UserID
	if role != model.RoleAdmin && int(id) != claimsId {
		http.Error(w, repo.ErrUnauthorized.Unwrap().Error(), http.StatusUnauthorized)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Не удалось получить файл: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Не удалось прочитать файл: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.usecases.Users.UploadImage(r.Context(), int(id), data)
	if err != nil {
		http.Error(w, tracerr.Unwrap(err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Изображение успешно загружено"))
}

func (rh *RestHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

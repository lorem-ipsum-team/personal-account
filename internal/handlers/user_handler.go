package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kerilOvs/profile_sevice/internal/models"
	"github.com/kerilOvs/profile_sevice/internal/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req struct {
		Name        string             `json:"name"`
		Surname     string             `json:"surname"`
		AboutMyself *string            `json:"about_myself,omitempty"`
		Gender      *models.UserGender `json:"gender,omitempty"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body"))
	}

	user, err := h.service.CreateUser(req.Name, req.Surname, req.AboutMyself, req.Gender)
	if err != nil {
		return errorResponseWithCode(c, err, http.StatusBadRequest)
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid user ID"))
	}

	if err := h.service.DeleteUser(id); err != nil {
		return errorResponseWithCode(c, err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandler) GetUserById(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid user ID"))
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		return errorResponseWithCode(c, err, http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUserProfile(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid user ID"))
	}

	var req models.UserProfileUpdate
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body"))
	}

	if err := h.service.UpdateUserProfile(id, req); err != nil {
		return errorResponseWithCode(c, err, http.StatusBadRequest)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandler) UpdateUserAbout(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid user ID"))
	}

	var req struct {
		AboutMyself string `json:"about_myself"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body"))
	}

	if err := h.service.UpdateUserAbout(id, req.AboutMyself); err != nil {
		return errorResponseWithCode(c, err, http.StatusBadRequest)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandler) UpdateUserName(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid user ID"))
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body"))
	}

	if err := h.service.UpdateUserName(id, req.Name); err != nil {
		return errorResponseWithCode(c, err, http.StatusBadRequest)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandler) GetUserPhotos(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid user ID"))
	}

	photos, err := h.service.GetUserPhotos(id)
	if err != nil {
		return errorResponseWithCode(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, photos)
}

func (h *UserHandler) AddUserPhoto(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid user ID"))
	}

	var req struct {
		URL string `json:"url"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body"))
	}

	photo, err := h.service.AddUserPhoto(id, req.URL)
	if err != nil {
		return errorResponseWithCode(c, err, http.StatusBadRequest)
	}

	return c.JSON(http.StatusCreated, photo)
}

func (h *UserHandler) RemoveUserPhoto(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid user ID"))
	}

	photoID, err := uuid.Parse(c.Param("photoId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid photo ID"))
	}

	if err := h.service.RemoveUserPhoto(userID, photoID); err != nil {
		return errorResponseWithCode(c, err, http.StatusNotFound)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandler) UpdatePrimaryPhoto(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid user ID"))
	}

	var req struct {
		ID *uuid.UUID `json:"id,omitempty"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body"))
	}

	if req.ID == nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Photo ID is required"))
	}

	if err := h.service.SetPrimaryPhoto(userID, *req.ID); err != nil {
		return errorResponseWithCode(c, err, http.StatusNotFound)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandler) UpdateUserSurname(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid user ID"))
	}

	var req struct {
		Surname string `json:"surname"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body"))
	}

	if err := h.service.UpdateUserSurname(id, req.Surname); err != nil {
		return errorResponseWithCode(c, err, http.StatusBadRequest)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandler) AddUserTag(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid user ID"))
	}

	var req struct {
		Tag string `json:"tag"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body"))
	}

	tag, err := h.service.AddUserTag(id, req.Tag)
	if err != nil {
		return errorResponseWithCode(c, err, http.StatusBadRequest)
	}

	return c.JSON(http.StatusCreated, tag)
}

func (h *UserHandler) GetUserTags(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid user ID"))
	}

	tags, err := h.service.GetUserTags(id)
	if err != nil {
		return errorResponseWithCode(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, tags)
}

func (h *UserHandler) RemoveUserTag(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid user ID"))
	}

	tagID, err := uuid.Parse(c.Param("tagId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid tag ID"))
	}

	if err := h.service.RemoveUserTag(userID, tagID); err != nil {
		return errorResponseWithCode(c, err, http.StatusNotFound)
	}

	return c.NoContent(http.StatusNoContent)
}

func errorResponse(msg string) map[string]interface{} {
	return map[string]interface{}{"error": msg}
}

func errorResponseWithCode(c echo.Context, err error, code int) error {
	return c.JSON(code, errorResponse(err.Error()))
}

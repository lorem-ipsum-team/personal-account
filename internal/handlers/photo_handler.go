package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kerilOvs/profile_sevice/internal/service"
	"github.com/labstack/echo/v4"
)

type PhotoHandler struct {
	userService  *service.UserService
	photoService *service.PhotoService
}

func NewPhotoHandler(userService *service.UserService, photoService *service.PhotoService) *PhotoHandler {
	return &PhotoHandler{
		userService:  userService,
		photoService: photoService,
	}
}

// @Summary Загрузить фото
// @Accept  multipart/form-data
// @Param   photo formData file true "Фото пользователя"
// @Success 201 {object} models.UserPhoto
func (h *PhotoHandler) UploadPhoto(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("invalid user id"))
	}

	// Получаем файл из формы
	file, err := c.FormFile("photo")
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("photo is required"))
	}

	// Проверяем тип файла
	if !isValidImageType(file.Header.Get("Content-Type")) {
		return c.JSON(http.StatusBadRequest, errorResponse("only jpeg/png images are allowed"))
	}

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("failed to read file"))
	}
	defer src.Close()

	// Загружаем фото в MinIO
	objectName, err := h.photoService.UploadPhoto(c.Request().Context(), src, file.Size)
	if err != nil {
		return errorResponseWithCode(c, err, http.StatusInternalServerError)
	}

	// Получаем URL для доступа к фото
	photoURL, err := h.photoService.GetPhotoURL(objectName, 24*time.Hour)
	if err != nil {
		return errorResponseWithCode(c, err, http.StatusInternalServerError)
	}

	// Сохраняем информацию о фото в БД
	photo, err := h.userService.AddUserPhoto(userID, photoURL)
	if err != nil {
		return errorResponseWithCode(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, photo)
}

// @Summary Получить фото
// @Produce json
// @Param   id path string true "ID фото"
// @Success 200
// @Header  200 {string} Content-Type "image/jpeg"
func (h *PhotoHandler) GetPhoto(c echo.Context) error {
	objectName := c.Param("id")
	url, err := h.photoService.GetPhotoURL(objectName, 24*time.Hour)
	if err != nil {
		return errorResponseWithCode(c, err, http.StatusNotFound)
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// Вспомогательная функция для проверки типа изображения
func isValidImageType(contentType string) bool {
	return contentType == "image/jpeg" ||
		contentType == "image/png" ||
		contentType == "image/jpg"
}

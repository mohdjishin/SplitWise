package handlers

import (
	"encoding/json"
	e "errors"
	"net/http"

	jUtil "github.com/mohdjishin/SplitWise/helper/jwt"
	"github.com/mohdjishin/SplitWise/helper/validate"
	"github.com/mohdjishin/SplitWise/internal/db"
	"github.com/mohdjishin/SplitWise/internal/errors"
	"github.com/mohdjishin/SplitWise/internal/models"
	"github.com/mohdjishin/SplitWise/internal/models/dto"
	"github.com/mohdjishin/SplitWise/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register handles user registration
// @Summary Register a new user
// @Description Registers a new user with email, password, and name. Returns conflict error if email already exists.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body dto.RegisterRequest true "User details"
// @Success 200 {object} map[string]string "User registered"
// @Failure 400 {object} errors.Error" Bad Request"
// @Failure 409 {object} errors.Error "Conflict"
// @Router /auth/register [post]
func Register(w http.ResponseWriter, r *http.Request) {

	var input dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.LoggerInstance.Error("Error decoding request body", zap.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validate.ValidateStruct(input); err != nil {
		logger.LoggerInstance.Error("Error validating request body", zap.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.ErrValidationFailed(err.Error()))
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := models.User{Email: input.Email, Password: string(hashedPassword), Name: input.Name}

	tx := db.GetDb().Create(&user)
	if tx.Error != nil {
		logger.LoggerInstance.Error("Error creating user", zap.Any("error", tx.Error))
		if e.Is(tx.Error, gorm.ErrDuplicatedKey) {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(errors.ErrUserAlreadyExists)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errors.ErrInternalServerError)
		}
		return
	}
	logger.LoggerInstance.Info("User registered", zap.String("email", user.Email))
	// TODO: Creating a Response Model.
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "User registered"})
}

// Login handles user login
// @Summary Login a user
// @Description Logs in a user with email and password.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.LoginRequest true "User credentials"
// @Success 200 {object} map[string]string "User logged in successfully, returns token"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized - Invalid credentials"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.LoggerInstance.Error("Error decoding request body", zap.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode((errors.ErrBadRequest))
		return
	}
	if err := validate.ValidateStruct(input); err != nil {
		logger.LoggerInstance.Error("Error validating request body", zap.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.ErrValidationFailed(err.Error()))
		return
	}
	var user models.User
	if err := db.GetDb().Where("email = ?", input.Email).First(&user).Error; err != nil {
		logger.LoggerInstance.Error("Error fetching user", zap.Any("error", err))
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errors.ErrInvalidCredential)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		logger.LoggerInstance.Error("Error comparing password", zap.Any("error", err))
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errors.ErrInvalidCredential)
		return
	}

	token, err := jUtil.GenerateToken(user.ID)

	if err != nil {
		logger.LoggerInstance.Error("Error generating token", zap.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}

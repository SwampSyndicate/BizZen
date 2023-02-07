package handlers

import (
	"encoding/json"
	"net/http"
	"server/config"
	"server/models"
	"server/utils"

	"github.com/gorilla/mux"
)

// TODO: Add comment documentation (type Credentials)
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TODO: Add comment documentation (func CreateUser)
func (db Handler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	user := models.User{}

	// ? Duplicative code block for decoding request body and error checking/response.
	// TODO:  Create new function to consolidate duplicative code (decoding request body / error handling).
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&user); err != nil {
		utils.RespondWithError(writer, http.StatusBadRequest, err.Error())
		return
	}
	defer request.Body.Close()

	// ? Should error handling and response be handled by the called function instead?
	if err := user.HashPassword(user.Password); err != nil {
		utils.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	// ? Should error handling and response be handled by the called function instead?
	// TODO: Create wrapper function for Handler type/struct to encapsulate "gorm.DB" logic
	if err := db.DB.Create(&user).Error; err != nil {
		utils.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(
		writer,
		http.StatusCreated,
		user)
}

// TODO: Add comment documentation (func Authenticate)
func (db Handler) Authenticate(writer http.ResponseWriter, request *http.Request) {
	var credentials Credentials

	// ? Duplicative code block for decoding request body and error checking/response.
	// TODO:  Create new function to consolidate duplicative code (decoding request body / error handling).
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&credentials); err != nil {
		utils.RespondWithError(writer, http.StatusBadRequest, err.Error())
		return
	}

	defer request.Body.Close()

	// ? Duplicative code block for checking if user exists. Should create a new type/struct to store user state and reference that instead.
	// TODO: Consolidate block into separate function and/or store user exists check in a new type/struct.
	user, err := db.checkIfUserExists(credentials.Email, writer, request)
	if err != nil {
		return
	}

	// ? Should error handling and response be handled by the called function instead?
	if err := user.CheckPassword(credentials.Password); err != nil {
		utils.RespondWithError(
			writer,
			http.StatusBadRequest,
			"Incorrect password.")

		return
	}

	// ? Should error handling and response be handled by the called function instead?
	validToken, err := GenerateToken(user.Email, user.AccountType, config.AppConfig.GetSigningKey())
	if err != nil {
		utils.RespondWithError(
			writer,
			http.StatusInternalServerError,
			err.Error())
	}

	utils.RespondWithJSON(
		writer,
		http.StatusOK,
		validToken)
}

// TODO: Add comment documentation (func checkIfUserExists)
func (db Handler) checkIfUserExists(userEmail string, writer http.ResponseWriter, request *http.Request) (*models.User, error) {
	var user models.User

	// ? Should error handling and response be handled by the called function instead?
	// TODO: Create wrapper function for Handler type/struct to encapsulate "gorm.DB" logic (First / checkIfUserExists)
	if err := db.DB.First(&user, models.User{Email: userEmail}).Error; err != nil {
		utils.RespondWithError(writer, http.StatusNotFound, "User does not exist.")
		return nil, err
	}

	return &user, nil
}

// TODO: Add comment documentation (func GetUser)
func (db Handler) GetUser(writer http.ResponseWriter, request *http.Request) {
	userEmail := mux.Vars(request)["email"]

	// ? Duplicative code block for checking if user exists. Should create a new type/struct to store user state and reference that instead.
	// TODO: Consolidate block into separate function and/or store user exists check in a new type/struct.
	user, err := db.checkIfUserExists(userEmail, writer, request)
	if err != nil {
		return
	}

	utils.RespondWithJSON(
		writer,
		http.StatusOK,
		user)
}

// TODO: Add comment documentation (func UpdateUser)
func (db Handler) UpdateUser(writer http.ResponseWriter, request *http.Request) {
	userEmail := mux.Vars(request)["email"]

	// ? Duplicative code block for checking if user exists. Should create a new type/struct to store user state and reference that instead.
	// TODO: Consolidate block into separate function and/or store user exists check in a new type/struct.
	user, err := db.checkIfUserExists(userEmail, writer, request)
	if err != nil {
		return
	}

	// ? Duplicative code block for decoding request body and error checking/response.
	// TODO:  Create new function to consolidate duplicative code (decoding request body / error handling).
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		utils.RespondWithError(writer, http.StatusBadRequest, err.Error())
		return
	}

	defer request.Body.Close()

	// ? Should error handling and response be handled by the called function instead?
	// TODO: Create wrapper function for Handler type/struct to encapsulate "gorm.DB" logic (Save / UpdateUser)
	if err := db.DB.Save(&user).Error; err != nil {
		utils.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(
		writer,
		http.StatusOK,
		user)

}

// TODO: Add comment documentation (func DeleteUser)
func (db Handler) DeleteUser(writer http.ResponseWriter, request *http.Request) {
	userEmail := mux.Vars(request)["email"]

	// ? Duplicative code block for checking if user exists. Should create a new type/struct to store user state and reference that instead.
	// TODO: Consolidate block into separate function and/or store user exists check in a new type/struct.
	user, err := db.checkIfUserExists(userEmail, writer, request)
	if err != nil {
		return
	}

	// ? Should error handling and response be handled by the called function instead?
	// TODO: Create wrapper function for Handler type/struct to encapsulate "gorm.DB" logic (DeleteUser)
	if err := db.DB.Delete(&user).Error; err != nil {
		utils.RespondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(
		writer,
		http.StatusOK,
		user)
}

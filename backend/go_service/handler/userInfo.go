package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"

	//"log"
	customErrors "appBE/errors"
	"appBE/model"
	"appBE/service"
	"net/http"

	"github.com/gorilla/mux"
)

func saveUserInfoHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one save user info request")

	// Check data type
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}

	// Decode the body into a User struct
	var user model.UserInfo
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
		return
	}

	// Call service to process and save the post
	if err := service.SaveUserInfo(&user); err != nil {
		if errors.Is(err, customErrors.ErrUsernameAlreadyExists) {
			http.Error(w, "user already exists", http.StatusNotFound)
		} else {
			log.Printf("Error saving user info: %v", err)
			// For all other errors, return internal server error
			http.Error(w, "Failed to save user from backend",
				http.StatusInternalServerError)
		}
		return
	}

	// Response
	fmt.Fprintf(w, "User info saved successfully\n")
	fmt.Fprintf(w, "The saved user %s is \n", user.Username)
}

func getUserInfoByIDHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get user info request")

	w.Header().Set("Content-Type", "application/json")
	userIDStr := mux.Vars(r)["userID"]

	// 1. process data
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid userID provided", http.StatusBadRequest)
		return
	}

	// 2. call service level to get user info
	user, err := service.GetUserInfoByID(userID)
	if err != nil {
		// Check if the error is due to the user not being found
		if errors.Is(err, customErrors.ErrUserNotFound){
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			// For all other errors, return internal server error
			http.Error(w, "Failed to search user by ID from backend",
				http.StatusInternalServerError)
		}
		return
	}

	// 3. format json response
	js, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to parse user into JSON format",
			http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func getUserIDByNameHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get user id request")

	// 1. process data
	
	w.Header().Set("Content-Type", "application/json")
	username := mux.Vars(r)["username"]

	// 2. call service level to get user id
	userID, err := service.GetUserIDByName(username)
	if err != nil {
		// Check if the error is due to the user not being found
		if errors.Is(err, customErrors.ErrUserNotFound){
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			// For all other errors, return internal server error
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
		}
		return
	}

	// 3. format json response
	js, err := json.Marshal(userID)
	if err != nil {
		http.Error(w, "Failed to parse user into JSON format",
			http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func getUsernameByIDHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get username request")

	// 1. process data
	
	w.Header().Set("Content-Type", "application/json")
	userIDStr := mux.Vars(r)["userID"]

	// 1. process data
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid userID provided", http.StatusBadRequest)
		return
	}
	// 2. call service level to get user id
	username, err := service.GetUsernameByID(userID)
	if err != nil {
		// Check if the error is due to the user not being found
		if errors.Is(err, customErrors.ErrUserNotFound){
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			// For all other errors, return internal server error
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
		}
		return
	}

	// 3. format json response
	js, err := json.Marshal(username)
	if err != nil {
		http.Error(w, "Failed to parse user into JSON format",
			http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-content-service/config"
	"github.com/temuka-content-service/models"
	"gorm.io/gorm"
)

var db *gorm.DB = config.GetDBInstance()

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Title       string `json:"title"`
		Description string `json:"desc"`
		UserID      int    `json:"user_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newPost := models.Post{
		Title:       requestBody.Title,
		Description: requestBody.Description,
		UserID:      requestBody.UserID,
	}
	db.Create(&newPost)

	response := struct {
		Message string      `json:"message"`
		Data    models.Post `json:"data"`
	}{
		Message: "Post has been created",
		Data:    newPost,
	}
	json.NewEncoder(w).Encode(response)
}

func GetTimelinePosts(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// userIDstr := vars["user_id"]

	// userID, err := strconv.Atoi(userIDstr)
	// if err != nil {
	// 	http.Error(w, "Invalid user id", http.StatusBadRequest)
	// 	return
	// }

	// var currentUser models.User
	// if err := db.First(&currentUser, "id = ?", userID).Error; err != nil {
	// 	http.Error(w, "Cannot retrieve the data because the user was not found", http.StatusNotFound)
	// 	return
	// }

	// var userPosts []models.Post
	// if err := db.Where("user_id = ?", currentUser.ID).Find(&userPosts).Error; err != nil {
	// 	http.Error(w, "Error retrieving user posts", http.StatusInternalServerError)
	// 	return
	// }

	// var friendPosts []models.Post
	// for _, friendID := range currentUser.Followings {
	// 	var posts []models.Post
	// 	if err := db.Where("user_id = ?", friendID).Find(&posts).Error; err != nil {
	// 		http.Error(w, "Error retrieving friend posts", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	friendPosts = append(friendPosts, posts...)
	// }

	// timelinePosts := append(userPosts, friendPosts...)

	// response := struct {
	// 	Message string        `json:"message"`
	// 	Data    []models.Post `json:"data"`
	// }{
	// 	Message: "Timeline posts has been retrieved",
	// 	Data:    timelinePosts,
	// }

	// json.NewEncoder(w).Encode(response)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDstr := vars["id"]

	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	if err := db.Delete(&models.Post{}, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error deleting post", http.StatusInternalServerError)
		}
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Post has been deleted",
	}

	json.NewEncoder(w).Encode(response)
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDstr := vars["id"]

	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		UserID int `json:"user_id"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving post", http.StatusBadRequest)
		}
		return
	}

	alreadyLiked := false
	for _, userID := range post.Likes {
		if userID == requestBody.UserID {
			alreadyLiked = true
			break
		}
	}

	if !alreadyLiked {
		post.Likes = append(post.Likes, requestBody.UserID)
		if err := db.Save(&post).Error; err != nil {
			http.Error(w, "Error liking post", http.StatusInternalServerError)
			return
		}

		response := struct {
			Message string `json:"message"`
		}{
			Message: "You have liked this post",
		}

		json.NewEncoder(w).Encode(response)
	}
}
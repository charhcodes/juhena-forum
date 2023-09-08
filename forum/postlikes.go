package forum

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

var postID int
var userID string

func getPostID(w http.ResponseWriter, r *http.Request) int {
	// Extract post ID from URL
	postIDStr := strings.TrimPrefix(r.URL.Path, "/post-like/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}
	return postID
}

func getUserID(w http.ResponseWriter, r *http.Request) string {
	// Get user ID from session cookie
	_, userID, err := GetCookieValue(r)
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusBadRequest)
	}
	return userID
}

// Handler for handling like and dislike actions
func HandleLikesDislikes(w http.ResponseWriter, r *http.Request) {
	// Check session cookie
	checkCookies(w, r)

	// Get post id from the URL path
	postIDStr := strings.TrimPrefix(r.URL.Path, "/post/")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	var comments []Comment

	// Assuming you have a function to retrieve the logged-in user's ID, if available
	// If not, you can set it to a default value or handle it as you see fit.
	// userID := "user1"

	// Get the post data by calling the getPostByID function or fetching it from the database
	post, err := getPostByID(postIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Get comments by postID
	comments, err = getCommentsByPostID(postIDStr)
	if err != nil {
		http.Error(w, "Could not fetch comments", http.StatusInternalServerError)
		return
	}

	// Get like and dislike counts for the post
	likeCount := GetLikeCount(postID, 0, "posts")
	dislikeCount := GetDislikeCount(postID, 0, "posts")

	// Assuming your Post struct has a field named PostID
	var data struct {
		PostID     int
		Post       Post
		Comments   []Comment
		Success    bool // Add the Success field to indicate if the comment was successfully posted
		Like       int  // Like count
		Dislike    int  // Dislike count
		IsLiked    bool // Indicate if the user has liked this post
		IsDisliked bool // Indicate if the user has disliked this post
	}

	data.PostID = postID
	data.Post = *post // Use the dereferenced post pointer
	data.Comments = comments
	data.Success = r.URL.Query().Get("success") == "1"
	data.Like = likeCount
	data.Dislike = dislikeCount

	// Get the user's ID and check if they have liked/disliked this post
	userID := getUserID(w, r)
	isLiked, isDisliked, _ := checkUserLikeDislike(userID, postID)
	data.IsLiked = isLiked
	data.IsDisliked = isDisliked

	tmpl, err := template.ParseFiles("postPage.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	action := r.FormValue("action")

	// Handle the like/dislike action based on user's input
	if action == "like" {
		err = addLike(userID, postID)

		if err != nil {
			http.Error(w, "Error adding like", http.StatusInternalServerError)
			return
		}
	} else if action == "dislike" {
		err = addDislike(userID, postID)
		if err != nil {
			http.Error(w, "Error adding dislike", http.StatusInternalServerError)
			return
		}
	}

	fmt.Println("Like/Dislike successful!")

	err = tmpl.ExecuteTemplate(w, "postPage.html", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect back to the post page
	postIDStr = strings.TrimPrefix(r.URL.Path, "/post-like/")
	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}

// Other database functions (e.g., addLike, removeLike, checkUserLikeDislike) should be implemented accordingly.

// checkUserLikeDislike checks if a user has already liked or disliked a post.
// func checkUserLikeDislike(userID string, postID int) (liked bool, disliked bool, err error) {
// 	// Execute a query to check if the user has liked the post
// 	row := DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = 1", userID, postID)
// 	var likeCount int
// 	if err := row.Scan(&likeCount); err != nil {
// 		return false, false, err
// 	}
// 	liked = likeCount > 0

// 	// Execute a query to check if the user has disliked the post
// 	row = DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = -1", userID, postID)
// 	var dislikeCount int
// 	if err := row.Scan(&dislikeCount); err != nil {
// 		return false, false, err
// 	}
// 	disliked = dislikeCount > 0

// 	return liked, disliked, nil
// }

func checkUserLikeDislike(userID string, postID int) (liked bool, disliked bool, err error) {
	// Execute a query to check if the user has liked or disliked the post
	row := DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = 1", userID, postID)
	var likeCount int
	if err := row.Scan(&likeCount); err != nil {
		return false, false, err
	}
	liked = likeCount > 0

	row = DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = -1", userID, postID)
	var dislikeCount int
	if err := row.Scan(&dislikeCount); err != nil {
		return false, false, err
	}
	disliked = dislikeCount > 0

	// check for data entries that have not been liked or disliked
	row = DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = 0", userID, postID)
	var neitherCount int
	if err := row.Scan(&neitherCount); err != nil {
		return false, false, err
	}

	// If there is a "neither" entry, treat it as a like and a dislike
	if neitherCount > 0 {
		liked = true
		disliked = true
	}

	return liked, disliked, nil
}

func addLike(userID string, postID int) error {
	existingLike, existingDislike, err := checkUserLikeDislike(userID, postID)
	if err != nil {
		return err
	}

	if existingDislike {
		// Replace value '-1' with '1'
		_, err = DB.Exec("UPDATE postlikes SET type = 1 WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		fmt.Println("Replaced Dislike with Like")
	} else if existingLike {
		// Toggle value '1' to '0'
		_, err = DB.Exec("UPDATE postlikes SET type = 0 WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		fmt.Println("Toggled Like to Neither")
	} else {
		// Insert new entry with value '1'
		_, err = DB.Exec("INSERT INTO postlikes (user_id, post_id, type) VALUES (?, ?, 1)", userID, postID)
		if err != nil {
			return err
		}
		fmt.Println("Added Like")
	}
	return nil
}

func addDislike(userID string, postID int) error {
	existingLike, existingDislike, err := checkUserLikeDislike(userID, postID)
	if err != nil {
		return err
	}

	if existingLike {
		// Replace value '1' with '-1'
		_, err = DB.Exec("UPDATE postlikes SET type = -1 WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		fmt.Println("Replaced Like with Dislike")
	} else if existingDislike {
		// Toggle value '-1' to '0'
		_, err = DB.Exec("UPDATE postlikes SET type = 0 WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		fmt.Println("Toggled Dislike to Neither")
	} else {
		// Insert new entry with value '-1'
		_, err = DB.Exec("INSERT INTO postlikes (user_id, post_id, type) VALUES (?, ?, -1)", userID, postID)
		if err != nil {
			return err
		}
		fmt.Println("Added Dislike")
	}
	return nil
}

func GetLikeCount(postID int, commentID int, contentType string) int {
	var count int
	if contentType == "post" {
		// Query the database to get the like count for the post with the given postID
		err := DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = 1", userID, postID).Scan(&count)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error fetching the like count")
			// No likes found for the post
		}
	} else if contentType == "comments" {
		err := DB.QueryRow("SELECT COUNT(*) FROM reactions WHERE user_id = ? AND post_id = ? AND comment_id = ? AND type = 1", userID, postID, commentID).Scan(&count)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error fetching the like count")
		}
		// An error occurred while querying the database

	}
	return count
}

// GetDislikeCount retrieves the dislike count for a specific post from the database.
func GetDislikeCount(postID int, commentID int, contentType string) int {
	var count int
	if contentType == "post" {
		// Query the database to get the like count for the post with the given postID
		err := DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = -1", userID, postID).Scan(&count)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error fetching the dislike count")
			// No likes found for the post
		}
	} else if contentType == "comments" {
		err := DB.QueryRow("SELECT COUNT(*) FROM reactions WHERE user_id = ? AND post_id = ? AND comment_id = ? AND type = -1", userID, postID, commentID).Scan(&count)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error fetching the dislike count")
		}
		// An error occurred while querying the database

	}
	return count
}

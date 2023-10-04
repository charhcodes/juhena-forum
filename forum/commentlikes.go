package forum

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// getCommentPostID get the post id from the URL
func getCommentPostID(w http.ResponseWriter, r *http.Request) int {
	// Extract post ID from URL
	postIDStr := strings.TrimPrefix(r.URL.Path, "/comment-like/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}
	return postID
}

// get comment ID
func getCommentID(postID int) (int, error) {
	query := "SELECT id FROM comments WHERE post_id = ?"
	var commentID int
	err := DB.QueryRow(query, postID).Scan(&commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Return 0 for no comment found, without an error
		}
		return 0, err
	}
	return commentID, nil
}

func CommentLikesHandler(w http.ResponseWriter, r *http.Request) {
	postID := getCommentPostID(w, r)
	sessionID := GetSessionIDFromRequest(r)
	if sessionID == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	postID = getCommentPostID(w, r)
	userID, err := getUserIDFromSessionID(sessionID)
	if err != nil {
		http.Error(w, "Error getting session ID", http.StatusInternalServerError)
		return
	}

	commentID, err := getCommentID(postID)
	if err != nil {
		log.Fatal(err)
	}

	r.ParseForm()
	action := r.FormValue("comment-action")

	// Check if the user has already liked/disliked the post
	_, _, _ = checkCommentLikeDislike(userID, postID, commentID)
	if err != nil {
		http.Error(w, "Error checking user like/dislike", http.StatusInternalServerError)
		return
	}

	// Handle the like/dislike action based on user's previous interaction
	if action == "like" {
		err = addCommentLike(userID, postID, commentID)
		if err != nil {
			http.Error(w, "Error adding like", http.StatusInternalServerError)
			return
		}
	} else if action == "dislike" {
		err = addCommentDislike(userID, postID, commentID)
		if err != nil {
			http.Error(w, "Error adding dislike", http.StatusInternalServerError)
			return
		}
	}

	fmt.Println("Comment Like/Dislike successful!")

	// Redirect back to the post page
	postIDStr := strings.TrimPrefix(r.URL.Path, "/comment-like/")
	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}

// check if user has previously liked or disliked a comment
func checkCommentLikeDislike(userID int, postID int, commentID int) (liked bool, disliked bool, err error) {
	// check user id, post id and comment id for whether they've previously liked or disliked the post
	row := DB.QueryRow("SELECT COUNT(*) FROM reactions WHERE user_id = ? AND post_id = ? AND comment_id = ? AND type = 1", userID, postID, commentID)
	var likeCount int
	if err := row.Scan(&likeCount); err != nil {
		return false, false, err
	}
	liked = likeCount > 0

	row = DB.QueryRow("SELECT COUNT(*) FROM reactions WHERE user_id = ? AND post_id = ? AND comment_id = ? AND type = -1", userID, postID, commentID)
	var dislikeCount int
	if err := row.Scan(&dislikeCount); err != nil {
		return false, false, err
	}
	disliked = dislikeCount > 0

	// check for data entries that have not been liked or disliked
	row = DB.QueryRow("SELECT COUNT(*) FROM reactions WHERE user_id = ? AND post_id = ? AND comment_id = ? AND type = 0", userID, postID, commentID)
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

func addCommentLike(userID int, postID int, commentID int) error {
	fmt.Println("addCommentLike function triggered")
	existingLike, existingDislike, err := checkCommentLikeDislike(userID, postID, commentID)
	if err != nil {
		return err
	}

	if existingDislike {
		// Replace value '-1' with '1'
		_, err = DB.Exec("UPDATE reactions SET type = 1 WHERE user_id = ? AND post_id = ? AND comment_id = ?", userID, postID, commentID)
		if err != nil {
			return err
		}
		fmt.Println("Replaced Dislike with Like")
	} else if existingLike {
		// Toggle value '1' to '0'
		_, err = DB.Exec("UPDATE reactions SET type = 0 WHERE user_id = ? AND post_id = ? AND comment_id = ?", userID, postID, commentID)
		if err != nil {
			return err
		}
		fmt.Println("Toggled Like to Neither")
	} else {
		// Insert new entry with value '1'
		_, err = DB.Exec("INSERT INTO reactions (user_id, post_id, comment_id, type) VALUES (?, ?, ?, 1)", userID, postID, commentID)
		if err != nil {
			return err
		}
		fmt.Println("Added Like")
	}
	// Update total likes and dislikes counts in the comments table
	err = addTotalLikesCommentDislikes(userID, postID, commentID)
	if err != nil {
		return err
	}
	return nil
}

func addCommentDislike(userID int, postID int, commentID int) error {
	existingLike, existingDislike, err := checkCommentLikeDislike(userID, postID, commentID)
	if err != nil {
		return err
	}

	if existingLike {
		// Replace value '1' with '-1'
		_, err = DB.Exec("UPDATE reactions SET type = -1 WHERE user_id = ? AND post_id = ? AND comment_id = ?", userID, postID, commentID)
		if err != nil {
			return err
		}
		fmt.Println("Replaced Like with Dislike")
	} else if existingDislike {
		// Toggle value '-1' to '0'
		_, err = DB.Exec("UPDATE reactions SET type = 0 WHERE user_id = ? AND post_id = ? AND comment_id = ?", userID, postID, commentID)
		if err != nil {
			return err
		}
		fmt.Println("Toggled Dislike to Neither")
	} else {
		// Insert new entry with value '-1'
		_, err = DB.Exec("INSERT INTO reactions (user_id, post_id, comment_id, type) VALUES (?, ?, ?, -1)", userID, postID, commentID)
		if err != nil {
			return err
		}
		fmt.Println("Added Dislike")
	}
	// Update total likes and dislikes counts in the comments table
	err = addTotalLikesCommentDislikes(userID, postID, commentID)
	if err != nil {
		return err
	}
	return nil
}

// create a function to add total likes from the database postlikes table
func getTotalCommentLikes(userID int, postID int, commentID int) (int, error) {
	//from the postlikes table, get the total number of likes for a post
	row := DB.QueryRow("SELECT COUNT(*) FROM reactions WHERE comment_id = ? AND type = 1", postID, commentID, userID)
	var likeCount int
	if err := row.Scan(&likeCount); err != nil {
		return 0, err
	}
	return likeCount, nil
}

// create a function to add total dislikes from the database postlikes table
func getTotalCommentDislikes(userID int, postID int, commentID int) (int, error) {
	//from the postlikes table, get the total number of dislikes for a post
	row := DB.QueryRow("SELECT COUNT(*) FROM reactions WHERE comment_id = ? AND type = -1", postID, commentID, userID)
	var dislikeCount int
	if err := row.Scan(&dislikeCount); err != nil {
		return 0, err
	}
	return dislikeCount, nil
}

// Add total likes and dislikes to the comments table
func addTotalLikesCommentDislikes(userID int, postID int, commentID int) error {
	// Get the total number of likes for a comment
	likeCount, err := getTotalCommentLikes(postID, commentID, userID)
	if err != nil {
		return err
	}

	// Get the total number of dislikes for a comment
	dislikeCount, err := getTotalCommentDislikes(postID, commentID, userID)
	if err != nil {
		return err
	}

	// Update the likes_count and dislikes_count columns in the comments table
	_, err = DB.Exec("UPDATE reactions SET likes_count = ?, dislikes_count = ? WHERE id = ?", likeCount, dislikeCount, commentID)
	if err != nil {
		return err
	}

	return nil
}

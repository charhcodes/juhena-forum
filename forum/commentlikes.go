package forum

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func getCommentPostID(w http.ResponseWriter, r *http.Request) int {
	// Extract post ID from URL
	postIDStr := strings.TrimPrefix(r.URL.Path, "/comment-like/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}
	return postID
}

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
	fmt.Println(postID)
	userID := getUserID(w, r)
	fmt.Println(userID)
	// commentID, err := getCommentID(postID)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(commentID)

	r.ParseForm()
	action := r.FormValue("comment-action")
	commentIDstr := r.FormValue("reaction-id")
	fmt.Println("hellooo", commentIDstr)
	commentID, err := strconv.Atoi(commentIDstr)
	if err != nil {
		fmt.Println("Error converting commentID to int")
	}

	// Check if the user has already liked/disliked the post
	_, _, err = checkCommentLikeDislike(userID, postID, commentID)
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

		// Update the likes count after liking
		err = updateCommentLikesCount(postID, commentID)
		if err != nil {
			http.Error(w, "Error updating likes count", http.StatusInternalServerError)
			return
		}

	} else if action == "dislike" {
		err = addCommentDislike(userID, postID, commentID)
		if err != nil {
			http.Error(w, "Error adding dislike", http.StatusInternalServerError)
			return
		}

		// Update the dislikes count after disliking
		err = updateCommentDislikesCount(postID, commentID)
		if err != nil {
			http.Error(w, "Error updating dislikes count", http.StatusInternalServerError)
			return
		}
	}

	fmt.Println("Comment Like/Dislike successful!")

	// Redirect back to the post page
	postIDStr := strconv.Itoa(postID)
	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}

// check if user has previously liked or disliked a comment
func checkCommentLikeDislike(userID string, postID int, commentID int) (liked bool, disliked bool, err error) {
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
func addCommentLike(userID string, postID int, commentID int) error {
	existingLike, existingDislike, err := checkCommentLikeDislike(userID, postID, commentID)
	if err != nil {
		return err
	}

	if existingLike {
		// If user has already liked, remove the like.
		_, err = DB.Exec("DELETE FROM reactions WHERE user_id = ? AND post_id = ? AND comment_id = ? AND type = 1", userID, postID, commentID)
	} else if existingDislike {
		// If user has disliked, remove the dislike and add a like.
		_, err = DB.Exec("UPDATE reactions SET type = 1 WHERE user_id = ? AND post_id = ? AND comment_id = ?", userID, postID, commentID)
	} else {
		// If user has neither liked nor disliked, insert a new like.
		_, err = DB.Exec("INSERT INTO reactions (user_id, post_id, comment_id, type) VALUES (?, ?, ?, 1)", userID, postID, commentID)
	}
	return err
}

func addCommentDislike(userID string, postID int, commentID int) error {
	existingLike, existingDislike, err := checkCommentLikeDislike(userID, postID, commentID)
	if err != nil {
		return err
	}

	if existingDislike {
		// If user has already disliked, remove the dislike.
		_, err = DB.Exec("DELETE FROM reactions WHERE user_id = ? AND post_id = ? AND comment_id = ? AND type = -1", userID, postID, commentID)
	} else if existingLike {
		// If user has liked, remove the like and add a dislike.
		_, err = DB.Exec("UPDATE reactions SET type = -1 WHERE user_id = ? AND post_id = ? AND comment_id = ?", userID, postID, commentID)
	} else {
		// If user has neither liked nor disliked, insert a new dislike.
		_, err = DB.Exec("INSERT INTO reactions (user_id, post_id, comment_id, type) VALUES (?, ?, ?, -1)", userID, postID, commentID)
	}
	return err
}

// func addCommentLike(userID string, postID int, commentID int) error {
// 	fmt.Println("addCommentLike function triggered")
// 	existingLike, existingDislike, err := checkCommentLikeDislike(userID, postID, commentID)
// 	if err != nil {
// 		return err
// 	}

// 	if existingDislike {
// 		// Replace value '-1' with '1'
// 		_, err = DB.Exec("UPDATE reactions SET type = 1 WHERE user_id = ? AND post_id = ? AND comment_id = ?", userID, postID, commentID)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Println("Replaced Dislike with Like")
// 	} else if existingLike {
// 		// Toggle value '1' to '0'
// 		_, err = DB.Exec("UPDATE reactions SET type = 0 WHERE user_id = ? AND post_id = ? AND comment_id = ?", userID, postID, commentID)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Println("Toggled Like to Neither")
// 	} else {
// 		// Insert new entry with value '1'
// 		_, err = DB.Exec("INSERT INTO reactions (user_id, post_id, comment_id, type) VALUES (?, ?, ?, 1)", userID, postID, commentID)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Println("Added Like")
// 	}
// 	return nil
// }

// func addCommentDislike(userID string, postID int, commentID int) error {
// 	existingLike, existingDislike, err := checkCommentLikeDislike(userID, postID, commentID)
// 	if err != nil {
// 		return err
// 	}

// 	if existingLike {
// 		// Replace value '1' with '-1'
// 		_, err = DB.Exec("UPDATE reactions SET type = -1 WHERE user_id = ? AND post_id = ? AND comment_id = ?", userID, postID, commentID)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Println("Replaced Like with Dislike")
// 	} else if existingDislike {
// 		// Toggle value '-1' to '0'
// 		_, err = DB.Exec("UPDATE reactions SET type = 0 WHERE user_id = ? AND post_id = ? AND comment_id = ?", userID, postID, commentID)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Println("Toggled Dislike to Neither")
// 	} else {
// 		// Insert new entry with value '-1'
// 		_, err = DB.Exec("INSERT INTO reactions (user_id, post_id, comment_id, type) VALUES (?, ?, ?, -1)", userID, postID, commentID)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Println("Added Dislike")
// 	}
// 	return nil
// }

// // create a function to add total likes from the database reaction table
func updateCommentLikesCount(postID int, commentID int) error {
	// Query to calculate total likes from the reactions table
	row := DB.QueryRow("SELECT COUNT(*) FROM reactions WHERE post_id = ? AND comment_id = ? AND type = 1", postID, commentID)
	var likeCount int
	if err := row.Scan(&likeCount); err != nil {
		return err
	}

	// Update the comments table
	_, err := DB.Exec("UPDATE comments SET commentlikes_count = ? WHERE post_id = ? AND id = ?", likeCount, postID, commentID)
	if err != nil {
		return err
	}
	fmt.Println("Updated commentlikes_count for the comment")
	return nil
}

// Similarly for dislikes
func updateCommentDislikesCount(postID int, commentID int) error {
	// Query to calculate total dislikes from the reactions table
	row := DB.QueryRow("SELECT COUNT(*) FROM reactions WHERE post_id = ? AND comment_id = ? AND type = -1", postID, commentID)
	var dislikeCount int
	if err := row.Scan(&dislikeCount); err != nil {
		return err
	}

	// Update the comments table
	_, err := DB.Exec("UPDATE comments SET commentdislikes_count = ? WHERE post_id = ? AND id = ?", dislikeCount, postID, commentID)
	if err != nil {
		return err
	}
	fmt.Println("Updated commentdislikes_count for the comment")
	return nil
}


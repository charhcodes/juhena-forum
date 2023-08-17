package forum

import (
	"net/http"
	"strconv"
	"strings"
)

// var counter int

// // insert like into table when like button is clicked on post
// func likePostButton(w http.ResponseWriter, r *http.Request) {
// 	// if post has already been liked
// 	if counter == 1 {
// 		// remove previous like
// 		_, err := DB.Exec("INSERT INTO postlikes (type) VALUES (0)")
// 		if err != nil {
// 			counter = 1
// 			log.Println(err)
// 			http.Error(w, "Could not remove like from post", http.StatusInternalServerError)
// 			return
// 		} else {
// 			counter = 0
// 			return
// 		}
// 	}
// 	// like post
// 	_, err := DB.Exec("INSERT INTO postlikes (type) VALUES (1)")
// 	if err != nil {
// 		counter = 0
// 		log.Println(err)
// 		http.Error(w, "Could not like post", http.StatusInternalServerError)
// 		return
// 	} else {
// 		counter = 1
// 		fmt.Println("Post like successful")
// 		return
// 	}
// }

// // insert dislike into table when dislike button is clicked on post
// func dislikePostButton(w http.ResponseWriter, r *http.Request) {
// 	// if post has already been disliked
// 	if counter == -1 {
// 		// remove previous dislike
// 		_, err := DB.Exec("INSERT INTO postlikes (type) VALUES (0)")
// 		if err != nil {
// 			counter = -1
// 			log.Println(err)
// 			http.Error(w, "Could not remove dislike from post", http.StatusInternalServerError)
// 			return
// 		} else {
// 			counter = 0
// 			return
// 		}
// 	}
// 	// dislike post
// 	_, err := DB.Exec("INSERT INTO postlikes (type) VALUES (-1)")
// 	if err != nil {
// 		counter = 0
// 		log.Println(err)
// 		http.Error(w, "Could not dislike post", http.StatusInternalServerError)
// 		return
// 	} else {
// 		counter = -1
// 		fmt.Println("Post dislike successful")
// 		return
// 	}
// }

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
func handleLikeDislike(w http.ResponseWriter, r *http.Request) {
	postID = getPostID(w, r)
	userID = getUserID(w, r)

	// Check if the user has already liked/disliked the post
	liked, disliked, err := checkUserLikeDislike(userID, postID)
	if err != nil {
		http.Error(w, "Error checking user like/dislike", http.StatusInternalServerError)
		return
	}

	// Handle the like/dislike action based on user's previous interaction
	if r.Method == http.MethodPost {
		if liked {
			// Handle removing like
			err := removeLike(userID, postID)
			if err != nil {
				http.Error(w, "Error removing like", http.StatusInternalServerError)
				return
			}
		} else if disliked {
			// Handle removing dislike and adding like
			err := removeDislike(userID, postID)
			if err != nil {
				http.Error(w, "Error removing dislike", http.StatusInternalServerError)
				return
			}
			err = addLike(userID, postID)
			if err != nil {
				http.Error(w, "Error adding like", http.StatusInternalServerError)
				return
			}
		} else {
			// Handle adding like
			err := addLike(userID, postID)
			if err != nil {
				http.Error(w, "Error adding like", http.StatusInternalServerError)
				return
			}
		}
	}

	// Redirect back to the post page
	postIDStr := strings.TrimPrefix(r.URL.Path, "/post-like/")
	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}

// Other database functions (e.g., addLike, removeLike, checkUserLikeDislike) should be implemented accordingly.

// checkUserLikeDislike checks if a user has already liked or disliked a post.
func checkUserLikeDislike(userID string, postID int) (liked bool, disliked bool, err error) {
	// Execute a query to check if the user has liked the post
	row := DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = 1", userID, postID)
	var likeCount int
	if err := row.Scan(&likeCount); err != nil {
		return false, false, err
	}
	liked = likeCount > 0

	// Execute a query to check if the user has disliked the post
	row = DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = -1", userID, postID)
	var dislikeCount int
	if err := row.Scan(&dislikeCount); err != nil {
		return false, false, err
	}
	disliked = dislikeCount > 0

	return liked, disliked, nil
}

func addLike(userID string, postID int) error {
	_, err := DB.Exec("INSERT INTO postlikes (user_id, post_id, type) VALUES (?, ?, 1)", userID, postID)
	if err != nil {
		return err
	}
	return nil
}

func addDislike(userID string, postID int) error {
	_, err := DB.Exec("INSERT INTO postlikes (user_id, post_id, type) VALUES (?, ?, -1)", userID, postID)
	if err != nil {
		return err
	}
	return nil
}

// func PostLikesHandler(w http.ResponseWriter, r *http.Request) {
// 	// Check session cookie
// 	checkCookies(w, r)

// 	// Get postID from URL path
// 	postIDStr := strings.TrimPrefix(r.URL.Path, "/post-like/")
// 	postID, err := strconv.Atoi(postIDStr)
// 	if err != nil {
// 		http.Error(w, "Invalid post ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Retrieve the sessionID and userID from the cookie
// 	_, userId, err := CommentgetCookieValue(r)
// 	if err != nil {
// 		http.Error(w, "cookie not found", http.StatusBadRequest)
// 		return
// 	}

// 	// insert data into table
// 	_, err = DB.Exec("INSERT INTO postlikes (user_id, post_id) VALUES (?, ?)",
// 		userId, postID)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, "Could not insert data to table", http.StatusInternalServerError)
// 		return
// 	}

// }

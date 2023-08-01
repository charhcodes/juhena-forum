package forum

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CREATE COMMENTS FUNCTION
func PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Check session cookie
	checkCookies(w, r)

	// Get postID from URL path
	postIDStr := strings.TrimPrefix(r.URL.Path, "/post-comment/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Assuming you have a function to retrieve the logged-in user's ID, if available
	// If not, you can set it to a default value or handle it as you see fit.
	userID := "user1"
	fmt.Println(userID)

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	postComment := r.Form.Get("commentContent")
	fmt.Println(postComment)

	if postComment == "" {
		fmt.Fprintln(w, "Error - please ensure comment box is not empty!")
		return
	}

	dateCreated := time.Now()
	updatedAt := dateCreated
	fmt.Println(dateCreated)

	_, err = DB.Exec("INSERT INTO comments (post_id, user_id, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		postID, userID, postComment, dateCreated, updatedAt)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not post comment", http.StatusInternalServerError)
		return
	}

	fmt.Println("Comment successfully posted!")

	http.Redirect(w, r, "/post-comment/"+postIDStr+"?success=1", http.StatusFound)

	// // Check session cookie
	// checkCookies(w, r)

	// // Get postID from URL path
	// postIDStr := strings.TrimPrefix(r.URL.Path, "/post/")
	// postID, err := strconv.Atoi(postIDStr)
	// if err != nil {
	// 	http.Error(w, "Invalid post ID", http.StatusBadRequest)
	// 	return
	// }

	// err = r.ParseForm()
	// if err != nil {
	// 	http.Error(w, "Could not parse form", http.StatusBadRequest)
	// 	return
	// }

	// postComment := r.Form.Get("commentContent") // Change to "commentContent"

	// if postComment == "" {
	// 	fmt.Fprintln(w, "Error - please ensure fields aren't empty!")
	// 	return
	// }

	// dateCreated := time.Now()

	// // Use the postID in the SQL query to associate the comment with the correct post
	// _, err = DB.Exec("INSERT INTO comments (post_id, content, created_at) VALUES (?, ?, ?)", postID, postComment, dateCreated)
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, "Could not post comment", http.StatusInternalServerError)
	// 	return
	// }

	// fmt.Println("Comment successfully posted!")
	// http.Redirect(w, r, "/post/"+postIDStr, http.StatusFound) // Redirect back to the post page

	// // Check session cookie
	// checkCookies(w, r)

	// err := r.ParseForm()
	// if err != nil {
	// 	http.Error(w, "Could not parse form", http.StatusBadRequest)
	// 	return
	// }

	// postID := r.Form.Get("postID")
	// commentContent := r.Form.Get("commentContent")

	// if postID == "" || commentContent == "" {
	// 	fmt.Fprintln(w, "Error - please ensure fields aren't empty!")
	// 	return
	// }

	// userID := 0 // Get the userID of the currently logged-in user, if available

	// dateCreated := time.Now()

	// _, err = DB.Exec("INSERT INTO comments (user_id, post_id, content, created_at) VALUES (?, ?, ?, ?)", userID, postID, commentContent, dateCreated)
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, "Could not post comment", http.StatusInternalServerError)
	// 	return
	// }

	// fmt.Println("Comment successfully posted!")
	// http.Redirect(w, r, "/post/"+postID, http.StatusFound)
}

// view and execute comments
// func executeComments(postID int) ([]Comment, error) {
// 	rows, err := DB.Query("SELECT id, user_id, content, created_at FROM comments WHERE post_id = ?", postID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var comment Comment
// 		// add user_id
// 		err := rows.Scan(&comment.Content, &comment.Time)
// 		if err != nil {
// 			return nil, err
// 		}
// 		// Format the datetime string
// 		t, err := time.Parse("2006-01-02T15:04:05.999999999-07:00", comment.Time)
// 		if err != nil {
// 			return nil, err
// 		}
// 		comment.Time = t.Format("January 2, 2006, 15:04:05")
// 		comments = append(comments, comment)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	// reverse posts
// 	comments = Creverse(comments)
// 	return comments, nil
// }

// // reverse comments section
// func Creverse(s []Comment) []Comment {
// 	//runes := []rune(s)
// 	length := len(s)
// 	for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
// 		s[i], s[j] = s[j], s[i]
// 	}
// 	return s
// }

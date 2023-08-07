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

	// Retrieve the sessionID and userID from the cookie
	_, userId, err := CommentgetCookieValue(r)
	if err != nil {
		http.Error(w, "cookie not found", http.StatusBadRequest)
		return
	}

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
	fmt.Println(userId)
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Could not convert", http.StatusInternalServerError)
		return
	}

	// Use userID and postID to create a new comment
	//user_ID gets excecuted to the database
	_, err = DB.Exec("INSERT INTO comments (post_id, user_id, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		postID, userIdInt, postComment, dateCreated, dateCreated)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not post comment", http.StatusInternalServerError)
		return
	}

	fmt.Println("Comment successfully posted!")

	http.Redirect(w, r, "/post/"+postIDStr+"?success=1", http.StatusFound)
}

func CommentgetCookieValue(r *http.Request) (string, string, error) {
	//split to cookie and value
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", "", err
	}
	value := strings.Split(cookie.Value, "&")

	return value[0], value[1], nil
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
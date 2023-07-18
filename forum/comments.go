package forum

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// CREATE COMMENTS FUNCTION
func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Check session cookie
	checkCookies(w, r)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	postComment := r.Form.Get("postComment")

	if postComment == "" {
		fmt.Fprintln(w, "Error - please ensure fields aren't empty!")
		return
	}

	dateCreated := time.Now()

	_, err = DB.Exec("INSERT INTO comments (content, created_at) VALUES (?, ?)", postComment, dateCreated)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not post comment", http.StatusInternalServerError)
		return
	}

	fmt.Println("Comment successfully posted!")
	http.Redirect(w, r, "/submit-comment", http.StatusFound)
}

// comment handler that uses JSON. don't use, i'm just keeping this for posterity and just in case

// type CommentResponse struct {
// 	Success bool   `json:"success"`
// 	Message string `json:"message"`
// }

// func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
// 	// Check session cookie
// 	checkCookies(w, r)

// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	err := r.ParseForm()
// 	if err != nil {
// 		http.Error(w, "Could not parse form", http.StatusBadRequest)
// 		return
// 	}

// 	postComment := r.PostForm.Get("postComment")

// 	if postComment == "" {
// 		response := CommentResponse{
// 			Success: false,
// 			Message: "Error - please ensure fields aren't empty!",
// 		}
// 		jsonResponse, err := json.Marshal(response)
// 		if err != nil {
// 			http.Error(w, "Could not marshal JSON response", http.StatusInternalServerError)
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write(jsonResponse)
// 		return
// 	}

// 	dateCreated := time.Now()

// 	_, err = DB.Exec("INSERT INTO comments (content, created_at) VALUES (?, ?)", postComment, dateCreated)
// 	if err != nil {
// 		response := CommentResponse{
// 			Success: false,
// 			Message: "Could not post comment",
// 		}
// 		jsonResponse, err := json.Marshal(response)
// 		if err != nil {
// 			http.Error(w, "Could not marshal JSON response", http.StatusInternalServerError)
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write(jsonResponse)
// 		return
// 	}

// 	response := CommentResponse{
// 		Success: true,
// 		Message: "Comment successfully posted!",
// 	}
// 	jsonResponse, err := json.Marshal(response)
// 	if err != nil {
// 		http.Error(w, "Could not marshal JSON response", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(jsonResponse)

// }

// view and execute comments
func executeComments() ([]Comment, error) {
	rows, err := DB.Query("SELECT id, user_id, content, created_at FROM comments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		// add user_id
		err := rows.Scan(&comment.Content, &comment.Time)
		if err != nil {
			return nil, err
		}
		// Format the datetime string
		t, err := time.Parse("2006-01-02T15:04:05.999999999-07:00", comment.Time)
		if err != nil {
			return nil, err
		}
		comment.Time = t.Format("January 2, 2006, 15:04:05")
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// reverse posts
	comments = Creverse(comments)
	return comments, nil
}

// reverse comments section
func Creverse(s []Comment) []Comment {
	//runes := []rune(s)
	length := len(s)
	for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

package forum

//ADDED
import (
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

type LikeDislikeCounts struct {
	LikeCount    int
	DislikeCount int
}

func PostLikeCounterHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := strings.TrimPrefix(r.URL.Path, "/post-like-counter/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	likeCount, dislikeCount, err := fetchLikeDislikeCounts(postID)
	if err != nil {
		http.Error(w, "Error fetching like/dislike counts", http.StatusInternalServerError)
		return
	}

	counts := LikeDislikeCounts{
		LikeCount:    likeCount,
		DislikeCount: dislikeCount,
	}

	tmpl, err := template.ParseFiles("postPage.html") // Replace with your template file
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, counts)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func fetchLikeDislikeCounts(postID int) (int, int, error) {
	// Execute queries to count likes, dislikes, and "neithers"
	likeRow := DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE post_id = ? AND type = 1", postID)
	var likeCount int
	if err := likeRow.Scan(&likeCount); err != nil {
		return 0, 0, err
	}

	dislikeRow := DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE post_id = ? AND type = -1", postID)
	var dislikeCount int
	if err := dislikeRow.Scan(&dislikeCount); err != nil {
		return 0, 0, err
	}

	neitherRow := DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE post_id = ? AND type = 0", postID)
	var neitherCount int
	if err := neitherRow.Scan(&neitherCount); err != nil {
		return 0, 0, err
	}

	// If there is a "neither" entry, treat it as a like and a dislike
	if neitherCount > 0 {
		likeCount += neitherCount
		dislikeCount += neitherCount
	}

	return likeCount, dislikeCount, nil
}

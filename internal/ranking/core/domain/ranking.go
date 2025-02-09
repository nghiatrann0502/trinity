package domain

const (
	ViewWeight    = 1
	LikeWeight    = 2
	CommentWeight = 2
	ShareWeight   = 3
	WatchWeight   = 0.01667 // about 1 point per minute
)

func ValidateAction(action string) bool {
	switch action {
	case "view", "like", "comment", "share", "watch":
		return true
	default:
		return false
	}
}

type VideoDetail struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	Duration    int    `json:"duration"`
	URL         string `json:"url"`
}

type VideoRanking struct {
	ID      int `json:"id"`
	VideoID int `json:"video_id"`
	Score   int `json:"score"`
}

type VideoRankingUpdate struct {
	VideoID int    `json:"video_id" example:"1"`
	Action  string `json:"action" example:"view"` // view, like, comment, share, watch
	Value   int    `json:"value" example:"1"`
}

package domain

type Video struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	Duration    int    `json:"duration"`
	URL         string `json:"url"`
}

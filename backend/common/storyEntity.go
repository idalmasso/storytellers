package common

//Story is the actual struct containing the data of a story
type Story struct {
	User       string `json:"user"`
	Title      string `json:"title"`
	Text       string `json:"text"`
	ImagesPath string `json:"images-path"`
}

package common

//Story is the actual struct containing the data of a story
type Story struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	User       string `json:"user"`
	Title      string `json:"title"`
	Text       string `json:"text"`
	ImagesPath string `json:"images-path"`
}

type StoryIdentifier interface {
	HasValue() bool
	SetValue(...interface{}) error
	SetFunctionRetrieve(func(db interface{}) (Story, error))
	Retrieve() func(db interface{}) (Story, error)
}

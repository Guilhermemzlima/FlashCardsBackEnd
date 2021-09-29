package deck

type DeckPreview struct {
	Id       string `json:"_id" bson:"_id" validate:"required"`
	ImageURL string `json:"imageURL" bson:"imageURL" validate:"required,max=100"`
	Name     string `json:"name" bson:"name" validate:"required,max=100"`
	UserId   string `json:"userId" bson:"userId" validate:"required"`
}

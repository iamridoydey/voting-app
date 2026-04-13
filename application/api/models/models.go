package models

type Language struct {
    Name  string `json:"name"  bson:"name"`
    Votes int    `json:"votes" bson:"votes"`
    Image string `json:"image" bson:"image"`
}
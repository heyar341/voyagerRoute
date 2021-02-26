package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserData struct {
	ID               primitive.ObjectID   `json:"id" bson:"_id"`
	UserName         string               `json:"username" bson:"username"`
	Email            string               `json:"email" bson:"email"`
	Password         []byte               `json:"password" bson:"password"`
	MultiRouteTitles map[string]time.Time `json:"multi_route_titles" bson:"multi_route_titles"`
}

type SessionData struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	SessionId string             `json:"session_id" bson:"session_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
}
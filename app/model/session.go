package model

import (
	"app/dbhandler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID        primitive.ObjectID `bson:"_id"`
	SessionID string             `bson:"session_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
}

func FindSession(sessionID string) (bson.M, error) {
	s := bson.M{"session_id": sessionID}
	d, err := dbhandler.Find("googroutes", "sessions", s, nil)
	return d, err
}

func DeleteSession(sessionID string) error {
	d := bson.M{"session_id": sessionID}
	err := dbhandler.Delete("googroutes", "sessions", d)
	return err
}

func CreateNewSession(sessionID string, userID primitive.ObjectID) error {
	d := bson.D{
		{"session_id", sessionID},
		{"user_id", userID},
	}
	_, err := dbhandler.Insert("googroutes", "sessions", d)
	return err
}

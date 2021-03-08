package model

import (
	"app/dbhandler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type route interface {
}

type MultiRoute struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id"`
	Title  string           `json:"title" bson:"title"`
	Routes map[string]route `json:"routes" bson:"routes"`
}

func SaveRoute(userID primitive.ObjectID, multiRoute *MultiRoute) error {
	routeDocument := bson.D{
		{"user_id", userID},
		{"title", multiRoute.Title},
		{"routes", multiRoute.Routes},
	}
	_, err := dbhandler.Insert("googroutes", "routes", routeDocument)
	return err
}

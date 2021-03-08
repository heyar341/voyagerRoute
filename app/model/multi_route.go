package model

import (
	"app/dbhandler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Route interface {
}

type MultiRoute struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id"`
	Title  string             `json:"title" bson:"title"`
	Routes map[string]Route   `json:"routes" bson:"routes"`
}

func (m *MultiRoute) SaveRoute(userID primitive.ObjectID) error {
	routeDocument := bson.D{
		{"user_id", userID},
		{"title", m.Title},
		{"routes", m.Routes},
	}
	_, err := dbhandler.Insert("googroutes", "routes", routeDocument)
	return err
}

func (m *MultiRoute) UpdateRoute() error{
	//routes collectionに保存
	routeDoc := bson.M{"_id": m.ID}
	updateDoc := bson.D{
		{"title", m.Title},
		{"routes", m.Routes},
	}
	err := dbhandler.UpdateOne("googroutes", "routes", "$set", routeDoc, updateDoc)
	return err
}

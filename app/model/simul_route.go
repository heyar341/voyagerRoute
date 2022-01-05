package model

import (
	"app/dbhandler"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SimulRouteData struct {
	Destination string `json:"destination" bson:"destination"`
	//Distance and Duration are stored as string here. They are processed as a number
	//at DoSimulSearch() in app/controllers/simulsearch/simul_search.go
	Distance string `json:"distance" bson:"distance"`
	Duration string `json:"duration" bson:"duration"`
}

type SimulRoute struct {
	ID          primitive.ObjectID        `json:"_id" bson:"_id"`
	Title       string                    `json:"title" bson:"title"`
	SimulRoutes map[string]SimulRouteData `json:"simul_routes" bson:"simul_routes"`
}

func (s *SimulRoute) SaveRoute(userID primitive.ObjectID) error {
	routeDocument := bson.D{
		{"user_id", userID},
		{"title", s.Title},
		{"simul_routes", s.SimulRoutes},
	}
	_, err := dbhandler.Insert("googroutes", "simulroutes", routeDocument)
	return err
}

func (s *SimulRoute) UpdateSimulRoute() error {
	//routes collectionに保存
	routeDoc := bson.M{"_id": s.ID}
	updateDoc := bson.D{
		{"title", s.Title},
		{"simul_routes", s.SimulRoutes},
	}
	err := dbhandler.UpdateOne("googroutes", "simulroutes", "$set", routeDoc, updateDoc)
	return err
}

func FindSimulRoute(userID primitive.ObjectID, title string) (bson.M, error) {
	//routes collectionから取得
	simulRouteDoc := bson.D{{"user_id", userID}, {"title", title}}
	r, err := dbhandler.Find("googroutes", "simulroutes", simulRouteDoc, nil)
	return r, err
}

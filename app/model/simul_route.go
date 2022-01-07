package model

import (
	"app/dbhandler"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DestinationData struct {
	PlaceId string `json:"place_id" bson:"place_id"`
	Address string `json:"address" bson:"address"`
	//Distance and Duration are stored as string here. They are processed as a number
	//at DoSimulSearch() in app/controllers/simulsearch/simul_search.go
	Distance string `json:"distance" bson:"distance"`
	Duration string `json:"duration" bson:"duration"`
}

//SimulParamの緯度と経度
type LatLng struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type SimulRoute struct {
	ID            primitive.ObjectID         `json:"_id" bson:"_id"`
	Title         string                     `json:"title" bson:"title"`
	Origin        string                     `json:"origin" bson:"origin"`
	OriginAddress string                     `json:"origin_address" bson:"origin_address"`
	Mode          string                     `json:"mode" bson:"mode"`
	DepartureTime string                     `json:"departure_time" bson:"departure_time"`
	LatLng        LatLng                     `json:"latlng" bson:"latLng"`
	Avoid         []string                   `json:"avoid" bson:"avoid"`
	Destinations  map[string]DestinationData `json:"destinations" bson:"destinations"`
}

type RouteUpdateRequest struct {
	SimulRoute
	PreviousTitle string `json:"previous_title" bson:"previous_title"`
}

func (s *SimulRoute) SaveRoute(userID primitive.ObjectID) error {
	routeDocument := bson.D{
		{"user_id", userID},
		{"title", s.Title},
		{"origin", s.Origin},
		{"origin_address", s.OriginAddress},
		{"mode", s.Mode},
		{"departure_time", s.DepartureTime},
		{"latlng", s.LatLng},
		{"avoid", s.Avoid},
		{"destinations", s.Destinations},
	}
	_, err := dbhandler.Insert("googroutes", "simulroutes", routeDocument)
	return err
}

func (s *SimulRoute) UpdateSimulRoute() error {
	//routes collectionに保存
	routeDoc := bson.M{"_id": s.ID}
	updateDoc := bson.D{
		{"title", s.Title},
		{"origin", s.Origin},
		{"origin_address", s.OriginAddress},
		{"mode", s.Mode},
		{"departure_time", s.DepartureTime},
		{"latlng", s.LatLng},
		{"avoid", s.Avoid},
		{"destinations", s.Destinations},
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

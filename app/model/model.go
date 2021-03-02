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

//ルート編集保存requestのフィールドを保存するstruct
type RouteUpdateRequest struct {
	ID            primitive.ObjectID     `json:"id" bson:"_id"`
	Title         string                 `json:"title" bson:"title"`
	PreviousTitle string                 `json:"previous_title" bson:"previous_title"`
	Routes        map[string]interface{} `json:"routes" bson:"routes"`
}

//同時検索のリクエストパラメータ
type SimulParams struct {
	Origin        string            `json:"origin"`
	Destinations  map[string]string `json:"destinations"`
	Mode          string            `json:"mode"`
	DepartureTime string            `json:"departure_time"`
	LatLng        LatLng            `json:"latlng"`
	Avoid         []string          `json:"avoid"`
}

//SimulParamの緯度と経度
type LatLng struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

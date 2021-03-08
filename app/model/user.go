package model

import (
	"app/dbhandler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID   `json:"id" bson:"_id"`
	UserName         string               `json:"username" bson:"username"`
	Email            string               `json:"email" bson:"email"`
	Password         []byte               `json:"password" bson:"password"`
	//MultiRouteTitles map[string]time.Time `json:"multi_route_titles" bson:"multi_route_titles"`
}

func UpdateMultiRouteTitles(userID primitive.ObjectID, routeTitle, operator string, updateVal interface{}) error {
		userDoc := bson.M{"_id": userID}
		updateField := bson.M{"multi_route_titles." + routeTitle: updateVal} //nested fieldsは.(ドット表記)で繋いで書く
		err := dbhandler.UpdateOne("googroutes", "users", operator, userDoc, updateField)
		return err
}
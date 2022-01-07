package model

import (
	"app/dbhandler"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	UserName string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Password []byte             `json:"password" bson:"password"`
	//MultiRouteTitles map[string]time.Time `json:"multi_route_titles" bson:"multi_route_titles"`
}

func UpdateMultiRouteTitles(userID primitive.ObjectID, routeTitle, operator string, updateVal interface{}) error {
	userDoc := bson.M{"_id": userID}
	updateField := bson.M{"multi_route_titles." + routeTitle: updateVal} //nested fieldsは.(ドット表記)で繋いで書く
	err := dbhandler.UpdateOne("googroutes", "users", operator, userDoc, updateField)
	return err
}

func UpdateSimulRouteTitles(userID primitive.ObjectID, routeTitle, operator string, updateVal interface{}) error {
	userDoc := bson.M{"_id": userID}
	updateField := bson.M{"simul_route_titles." + routeTitle: updateVal} //nested fieldsは.(ドット表記)で繋いで書く
	err := dbhandler.UpdateOne("googroutes", "users", operator, userDoc, updateField)
	return err
}

func FindUser(fieldName string, fieldValue interface{}) (bson.M, error) {
	//取得するドキュメントの条件
	userDoc := bson.M{fieldName: fieldValue}
	//DBから取得
	d, err := dbhandler.Find("googroutes", "users", userDoc, nil)
	return d, err
}

func UpdateUser(userID primitive.ObjectID, fieldName string, fieldValue interface{}) error {
	//user documentを更新
	userDoc := bson.M{"_id": userID}
	updateDoc := bson.D{{fieldName, fieldValue}}
	err := dbhandler.UpdateOne("googroutes", "users", "$set", userDoc, updateDoc)
	return err
}

func SaveNewUser(userName, email string, password []byte) (primitive.ObjectID, error) {
	userDoc := bson.D{{"username", userName},
		{"email", email},
		{"password", password},
		{"multi_route_titles", map[string]time.Time{}},
	}

	userID, err := dbhandler.Insert("googroutes", "users", userDoc)
	return userID, err
}

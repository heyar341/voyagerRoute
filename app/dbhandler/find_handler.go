package dbhandler

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

///optionDoc(フィールド指定)は順番関係ないからtypeはDでなくM
func Find(dbName, collectionName string, document interface{}, optionDoc bson.M) (bson.M, error) {
	client, ctx, cancel, err := connectDB()
	defer cancel()
	if err != nil {
		return nil, err
	}
	//処理終了後に切断
	defer client.Disconnect(ctx)
	database := client.Database(dbName)
	collection := database.Collection(collectionName)
	opts := options.FindOne().SetProjection(optionDoc)
	//DBからのレスポンスを挿入する変数
	var response bson.M
	err = collection.FindOne(ctx, document, opts).Decode(&response)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Printf("Error while finding a document: %v", err)
		}
		return nil, err
	}
	return response, nil
}

func FindAll(dbName, collectionName string, document interface{}, optionDoc bson.M) (interface{}, error) {
	client, ctx, cancel, err := connectDB()
	defer cancel()
	if err != nil {
		return nil, err
	}
	//処理終了後に切断
	defer client.Disconnect(ctx)
	database := client.Database(dbName)
	collection := database.Collection(collectionName)
	opts := options.Find().SetProjection(optionDoc)
	cursor, err := collection.Find(ctx, document, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Error while Finding documents: %v", err)
		}
		return nil, err
	}

	//DBからのレスポンスを挿入する変数
	var response []bson.D

	err = cursor.All(context.TODO(), &response)
	if err != nil {
		log.Printf("Error while decoding documents: %v", err)
	}

	return response, nil
}

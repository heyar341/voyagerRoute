package dbhandler

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Insert(dbName, collectionName string, document interface{}) (primitive.ObjectID, error) {
	//objectIDを取得するには、１番目の帰り値のInsertedIDフィールドを取得する
	client, ctx, cancel, err := connectDB()
	defer cancel()
	if err != nil {
		return primitive.NilObjectID, err
	}
	//処理終了後に切断
	defer client.Disconnect(ctx)
	database := client.Database(dbName)
	collection := database.Collection(collectionName)
	insertRes, err := collection.InsertOne(ctx, document)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return insertRes.InsertedID.(primitive.ObjectID), nil
}

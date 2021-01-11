package dbhandler

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func Find(dbName, collectionName string, document interface{}) (interface{}, error) {
	client, ctx, err := connectDB()
	if err != nil {
		return nil, err
	}
	//処理終了後に切断
	defer client.Disconnect(ctx)
	database := client.Database(dbName)
	collection := database.Collection(collectionName)
	//DBからのレスポンスを挿入する変数
	var response bson.D
	err = collection.FindOne(ctx, document).Decode(&response)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Fatal("ドキュメントが見つかりません")
		}
		return nil, err
	}
	return response, nil
}

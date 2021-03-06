package dbhandler

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Delete(dbName, collectionName string, document interface{}) error {
	//objectIDを取得するには、１番目の帰り値のInsertedIDフィールドを取得する
	client, ctx, cancel, err := connectDB()
	defer cancel()
	if err != nil {
		return err
	}
	//処理終了後に切断
	defer client.Disconnect(ctx)
	database := client.Database(dbName)
	collection := database.Collection(collectionName)
	//DBからのレスポンスを挿入する変数
	var deletedDocument bson.M
	err = collection.FindOneAndDelete(ctx, document).Decode(&deletedDocument)
	if err != nil && err == mongo.ErrNoDocuments {
		log.Println("During deleting a document: ", err)
		return err
	}
	return nil
}

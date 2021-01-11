package dbhandler

import "go.mongodb.org/mongo-driver/mongo"

func Insert(dbName,collectionName string, document interface{}) (*mongo.InsertOneResult, error){
	//objectIDを取得するには、１番目の帰り値のInsertedIDフィールドを取得する
	client, ctx, err := connectDB()
	if err != nil {
		return nil, err
	}
	//処理終了後に切断
	defer client.Disconnect(ctx)
	database := client.Database(dbName)
	collection := database.Collection(collectionName)
	insertRes, err := collection.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return insertRes, nil
}

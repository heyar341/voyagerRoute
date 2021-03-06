package dbhandler

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateOne(dbName, collectionName, updateMode string, document, updateField interface{}) error {

	client, ctx, cancel, err := connectDB()
	defer cancel()
	if err != nil {
		return err
	}
	//処理終了後に切断
	defer client.Disconnect(ctx)

	database := client.Database(dbName)
	collection := database.Collection(collectionName)
	update := bson.D{{updateMode, updateField}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	var updatedDocument bson.D
	err = collection.FindOneAndUpdate(ctx, document, update, opts).Decode(&updatedDocument)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Error while updating the document: %v", err)
			return err
		}
		log.Println(err)
		return err
	}

	return nil
}

package bsonconv

import (
	"app/customerr"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

//DocToStruct converts a bson document to a struct
func DocToStruct(d bson.M, s interface{}, e *error, modelName string) {
	if *e != nil {
		return
	}
	b, err := bson.Marshal(d)
	if err != nil {
		*e = customerr.BaseErr{
			Op:  "convert BSON document to struct",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while bson marshaling %s: %w", modelName, err),
		}
	}
	err = bson.Unmarshal(b, s)
	if err != nil {
		*e = customerr.BaseErr{
			Op:  "convert BSON document to struct",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while bson unmarshaling %s: %w", modelName, err),
		}
	}
}

package dbhandler

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func Connect()(*mongo.Client, context.Context, error){
	env_err := godotenv.Load("env/dev.env")
	if env_err != nil{
		panic("Can't load env file")
	}
	//envファイルからDB情報取得
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://"+DB_USER+":"+DB_PASSWORD+"@"+DB_HOST+":"+DB_PORT))
	if err != nil {
		return nil, nil, fmt.Errorf("データベースの情報が間違っています。")
	}
	ctx, _ := context.WithTimeout(context.Background(), 20 * time.Second)

	err = client.Connect(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("データベースに接続できません。")
	}
	//clientを返してDB操作をできるようにする
	return client,ctx,nil
}

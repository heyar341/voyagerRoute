package mypage

import (
	"app/dbhandler"
	"app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"sort"
	"time"
)

//user documentのmulti_route_titlesフィールドの値を入れるstruct
type TitleMap struct {
	TitleName string
	TimeStamp time.Time
}

type TitileSlice []TitleMap

func (t TitileSlice) Len() int      { return len(t) }
func (t TitileSlice) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

//TimestampのAfterメソッドで、ソート時に最新のタイトルが先頭に来るようにする
func (t TitileSlice) Less(i, j int) bool { return t[i].TimeStamp.After(t[j].TimeStamp) }

func RouteTitles(userID primitive.ObjectID) []string {
	result, err := dbhandler.Find("googroutes", "users", bson.D{{"_id", userID}}, bson.M{"username": 1, "multi_route_titles": 1})
	if err != nil {
		return []string{}
	}
	bsonByte, err := bson.Marshal(result)
	if err != nil {
		log.Println("Error while json marshaling: %v", err)
	}

	var user model.UserData
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &user)
	titles := user.MultiRouteTitles

	tSlice := make(TitileSlice, len(titles))
	i := 0
	for k, v := range titles {
		tSlice[i] = TitleMap{k, v}
		i++
	}
	//保存日時順にソート
	sort.Sort(tSlice)
	//タイトル名を入れるsliceを作成
	titleNames := make([]string, 0, len(titles))
	for _, tMap := range tSlice {
		titleNames = append(titleNames, tMap.TitleName)
	}

	return titleNames
}

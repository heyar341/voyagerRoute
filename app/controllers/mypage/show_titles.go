package mypage

import (
	"app/dbhandler"
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
	userDoc := bson.D{{"_id", userID}}
	bsonDoc, err := dbhandler.Find("googroutes", "users", userDoc, nil)
	if err != nil {
		return []string{}
	}

	titlesM := bsonDoc["multi_route_titles"].(primitive.M) //bson M型 (map[string]interface{})

	var titles = make(map[string]time.Time)
	for title, tStamp := range titlesM {
		mongoTS, ok := tStamp.(primitive.DateTime)
		if !ok {
			log.Println("Assertion error at checking timestamp type")
			return []string{}
		}
		timeStamp := mongoTS.Time() //time.Time型に変換
		titles[title] = timeStamp
	}

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

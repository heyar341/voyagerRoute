package mypage

import (
	"app/model"
	"log"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func routeTitles(userID primitive.ObjectID) []string {
	b, err := model.FindUser("_id", userID)
	if err != nil {
		return []string{}
	}
	titlesM := b["multi_route_titles"].(primitive.M) //bson M型 (map[string]interface{})

	var titles = make(map[string]time.Time)
	for title, tStamp := range titlesM {
		t, ok := tStamp.(primitive.DateTime)
		if !ok {
			log.Println("Assertion error at checking timestamp type")
			return []string{}
		}
		timeStamp := t.Time() //time.Time型に変換
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

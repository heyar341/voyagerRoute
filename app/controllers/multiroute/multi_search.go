package multiroute

import (
	"app/model"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strings"
	"time"
)

type routesData struct {
	routesInfo model.MultiRoute
	userID     primitive.ObjectID
	err        error
}

type multiSearchErr struct {
	Op  string
	Msg string
	Err error
}

func (s multiSearchErr) Error() string {
	return s.Err.Error()
}

func getRoutesInfo(req *http.Request) *routesData {
	reqFields, ok := req.Context().Value("reqFields").(model.MultiRoute)
	if !ok {
		return &routesData{
			err: multiSearchErr{
				Op:  "Get req routes info from context",
				Msg: "エラーが発生しました。",
				Err: fmt.Errorf("error while getting request fields from reuest's context"),
			},
		}
	}
	return &routesData{
		routesInfo: reqFields,
		err:        nil,
	}
}

func (r *routesData) getUserID(req *http.Request) {
	if r.err != nil {
		return
	}
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		r.err = multiSearchErr{
			Op:  "Get user from context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting user from reuest's context"),
		}
		return
	}
	r.userID = user.ID
}

func (r *routesData) saveRoute() {
	err := model.SaveRoute(r.userID, &r.routesInfo)
	if err != nil {
		e := multiSearchErr{
			Op:  "Save new multi route",
			Err: fmt.Errorf("error while saving multi route: %w", err),
		}

		if strings.Contains(err.Error(), "(BSONObjectTooLarge)") {
			e.Msg = "ルートのデータサイズが大きすぎるため、保存できません。\nルートの分割をお願いします。"
			r.err = e
		} else {
			e.Msg = "エラーが発生しました。"
			r.err = e
		}
	}
}

//user documentのmult_route_titlesに新しいルート名とタイムスタンプを追加する関数
func (r *routesData) updateRouteTitles() {
	if r.err != nil {
		return
	}
	now := time.Now().UTC() //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
	err := model.UpdateMultiRouteTitles(r.userID, r.routesInfo.Title, "$set", now)
	if err != nil {
		r.err = multiSearchErr{
			Op:  "update user document's multi_route_titles field",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while updating user's multi_route_titles %w", err),
		}
		return
	}
}

func SaveNewRoute(w http.ResponseWriter, req *http.Request) {
	//バリデーション完了後のrequest Fieldsを取得
	r := getRoutesInfo(req)
	//Auth middlewareからuserIDを取得
	r.getUserID(req)
	//routes collectionに保存
	r.saveRoute()
	/*users collectionのmulti_route_titlesフィールドにルート名と作成時刻を追加($set)する。
	  作成時刻はルート名取得時に作成時刻でソートするため*/
	r.updateRouteTitles()

	if r.err != nil {
		e := r.err.(multiSearchErr)
		http.Error(w, e.Msg, http.StatusInternalServerError)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	msgJSON := map[string]string{"msg": "OK"}
	json.NewEncoder(w).Encode(&msgJSON)
}

package multiroute

import (
	"app/controllers"
	"app/customerr"
	"app/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type routesData struct {
	routesInfo model.MultiRoute
	user       model.User
	err        error
}

//getRoutesInfo gets multiRoute from request's context
func getRoutesInfo(req *http.Request) *routesData {
	reqFields, ok := req.Context().Value("reqFields").(model.MultiRoute)
	if !ok {
		return &routesData{
			err: customerr.BaseErr{
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

//saveRoute saves route document to routes collection
func (r *routesData) saveRoute() {
	err := r.routesInfo.SaveRoute(r.user.ID)
	if err != nil {
		e := customerr.BaseErr{
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

//addRouteTitle adds mult_route_titles field to user document
func (r *routesData) addRouteTitle() {
	if r.err != nil {
		return
	}
	now := time.Now().UTC() //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
	err := model.UpdateMultiRouteTitles(r.user.ID, r.routesInfo.Title, "$set", now)
	if err != nil {
		r.err = customerr.BaseErr{
			Op:  "update user document's multi_route_titles field",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while updating user's multi_route_titles %w", err),
		}
		return
	}
}

func SaveNewRoute(w http.ResponseWriter, req *http.Request) {
	r := getRoutesInfo(req)
	r.user, r.err = controllers.GetUserFromCtx(req)
	r.saveRoute()
	/*users collectionのmulti_route_titlesフィールドにルート名と作成時刻を追加($set)する。
	  作成時刻はルート名取得時に作成時刻でソートするため*/
	r.addRouteTitle()

	if r.err != nil {
		e := r.err.(customerr.BaseErr)
		http.Error(w, e.Msg, http.StatusInternalServerError)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	msgJSON := map[string]string{"msg": "OK"}
	json.NewEncoder(w).Encode(&msgJSON)
}

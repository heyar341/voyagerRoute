package multiroute

import (
	"app/contexthandler"
	"app/customerr"
	"app/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type routeData struct {
	multiRoute model.MultiRoute
	user       model.User
	err        error
}

//getMultiRouteFromCtx gets multiRoute from request's context
func (r *routeData) getMultiRouteFromCtx(req *http.Request) {
	reqFields, ok := req.Context().Value("reqFields").(model.MultiRoute)
	if !ok {
		r.err = customerr.BaseErr{
			Op:  "Get req routes info from context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting request fields from reuest's context"),
		}
	}
	r.multiRoute = reqFields
}

//saveRoute saves route document to routes collection
func (r *routeData) saveRoute() {
	err := r.multiRoute.SaveRoute(r.user.ID)
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
func (r *routeData) addRouteTitle() {
	if r.err != nil {
		return
	}
	now := time.Now().UTC() //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
	err := model.UpdateMultiRouteTitles(r.user.ID, r.multiRoute.Title, "$set", now)
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
	var r = &routeData{}
	r.getMultiRouteFromCtx(req)
	contexthandler.GetUserFromCtx(req, &r.user, &r.err)
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

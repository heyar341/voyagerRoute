package simulsearch

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

type updateRouteData struct {
	simulRoute    model.SimulRoute
	previousTitle string
	user          model.User
	err           error
}

//getUpdateRouteFromCtx gets validated RouteUpdateRequest from context
func (u *updateRouteData) getUpdateRouteFromCtx(req *http.Request) {
	reqFields, ok := req.Context().Value("reqFields").(model.RouteUpdateRequest)
	if !ok {
		u.err = customerr.BaseErr{
			Op:  "Get req routes info from context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting request fields from reuest's context"),
		}
	}
	u.previousTitle = reqFields.PreviousTitle
	u.simulRoute = model.SimulRoute{
		ID:            reqFields.ID,
		Title:         reqFields.Title,
		Origin:        reqFields.Origin,
		OriginAddress: reqFields.OriginAddress,
		Mode:          reqFields.Mode,
		DepartureTime: reqFields.DepartureTime,
		LatLng:        reqFields.LatLng,
		Avoid:         reqFields.Avoid,
		Destinations:  reqFields.Destinations,
	}
}

//updateRoute updates route document in routes collection
func (u *updateRouteData) updateRoute() {
	if u.err != nil {
		return
	}
	err := u.simulRoute.UpdateSimulRoute()
	if err != nil {
		e := customerr.BaseErr{
			Op:  "Save new simul route",
			Err: fmt.Errorf("error while updating simul route: %w", err),
		}

		if strings.Contains(err.Error(), "(BSONObjectTooLarge)") {
			e.Msg = "ルートのデータサイズが大きすぎるため、保存できません。\nルートの分割をお願いします。"
			u.err = e
		} else {
			e.Msg = "エラーが発生しました。"
			u.err = e
		}
	}
}

//updateRouteTitles updates timestamp of simul_route_titles field in users collection, and delete previous title if title was changed.
func (u *updateRouteData) updateRouteTitles() {
	if u.err != nil {
		return
	}
	now := time.Now().UTC() //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
	if u.simulRoute.Title != u.previousTitle {
		//「元のルート名をuser documentから削除」
		//documentではなく、document内のフィールドを削除する場合、Deleteではなく、Update operatorの$unsetを使って削除する
		//公式ドキュメントURL: https://docs.mongodb.com/manual/reference/operator/update/unset/
		err := model.UpdateSimulRouteTitles(u.user.ID, u.previousTitle, "$unset", "")
		if err != nil {
			u.err = customerr.BaseErr{
				Op:  "Remove previous simul route title",
				Msg: "エラーが発生しました。",
				Err: fmt.Errorf("error while removing previous simul route title: %w", err),
			}
			return
		}
	}
	//「タイムスタンプを更新または追加」
	err := model.UpdateSimulRouteTitles(u.user.ID, u.simulRoute.Title, "$set", now)
	if err != nil {
		u.err = customerr.BaseErr{
			Op:  "Set new simul route title and timestamp",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while setting new simul route title and timestamp: %w", err),
		}
		return
	}
}

//ルートを更新保存するための関数
func UpdateRoute(w http.ResponseWriter, req *http.Request) {
	var u = &updateRouteData{}
	u.getUpdateRouteFromCtx(req)
	contexthandler.GetUserFromCtx(req, &u.user, &u.err)
	u.updateRoute()
	u.updateRouteTitles()

	if u.err != nil {
		e := u.err.(customerr.BaseErr)
		http.Error(w, e.Msg, http.StatusInternalServerError)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	respMsg := map[string]string{"msg": "ok"}
	json.NewEncoder(w).Encode(&respMsg)
}

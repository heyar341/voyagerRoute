package simulsearch

import (
	"app/contexthandler"
	"app/customerr"
	"app/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type simulRouteData struct {
	simulRoute model.SimulRoute
	user       model.User
	err        error
}

//getSimulRouteFromCtx gets simulRoute from request's context
func (s *simulRouteData) getSimulRouteFromCtx(req *http.Request) {
	simulRouteFields, ok := req.Context().Value("simulRouteFields").(model.SimulRoute)
	if !ok {
		s.err = customerr.BaseErr{
			Op:  "Get req routes info from context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting request fields from reuest's context"),
		}
	}
	s.simulRoute = simulRouteFields
}

//saveRoute saves route document to routes collection
func (s *simulRouteData) saveRoute() {
	err := s.simulRoute.SaveRoute(s.user.ID)
	if err != nil {
		e := customerr.BaseErr{
			Op:  "Save new simul route",
			Err: fmt.Errorf("error while saving simul route: %w", err),
		}

		e.Msg = "エラーが発生しました。"
		s.err = e
	}
}

//addRouteTitle adds simulroute_titles field to user document
func (s *simulRouteData) addRouteTitle() {
	if s.err != nil {
		return
	}
	now := time.Now().UTC() //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
	err := model.UpdateSimulRouteTitles(s.user.ID, s.simulRoute.Title, "$set", now)
	if err != nil {
		s.err = customerr.BaseErr{
			Op:  "update user document's simul_route_titles field",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while updating user's simul_route_titles %w", err),
		}
		return
	}
}

func SaveNewRoute(w http.ResponseWriter, req *http.Request) {
	var s = &simulRouteData{}
	s.getSimulRouteFromCtx(req)
	contexthandler.GetUserFromCtx(req, &s.user, &s.err)
	s.saveRoute()
	/*users collectionのsimul_route_titlesフィールドにルート名と作成時刻を追加($set)する。
	  作成時刻はルート名取得時に作成時刻でソートするため*/
	s.addRouteTitle()

	if s.err != nil {
		e := s.err.(customerr.BaseErr)
		http.Error(w, e.Msg, http.StatusInternalServerError)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	msgJSON := map[string]string{"msg": "OK"}
	json.NewEncoder(w).Encode(&msgJSON)
}

package simulsearch

import (
	"app/bsonconv"
	"app/contexthandler"
	"app/cookiehandler"
	"app/customerr"
	"app/model"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

var showRouteTpl *template.Template

func init() {
	showRouteTpl = template.Must(template.Must(template.ParseGlob("templates/simul_search/show_and_edit/simul_search_show.html")).ParseGlob("templates/includes/*.html"))
}

type editRoute struct {
	user       model.User
	routeModel model.SimulRoute
	routeJSON  string
	err        error
}

//getRouteFromDB gets route document from DB
func (eR *editRoute) getRouteFromDB(title string) bson.M {
	if eR.err != nil {
		return nil
	}
	d, err := model.FindSimulRoute(eR.user.ID, title)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			eR.err = customerr.BaseErr{
				Op:  "Finding simul route document",
				Msg: "ご指定いただいたルートがありません。",
				Err: fmt.Errorf("error while finding route document from routes collection: %w", err),
			}
			return nil
		} else {
			eR.err = customerr.BaseErr{
				Op:  "Finding route document",
				Msg: "エラーが発生しました。",
				Err: fmt.Errorf("error while finding route document from routes collection: %w", err),
			}
			return nil
		}
	}
	return d
}

//convertStructToJSON makes JSON object from simulRoute struct
func (eR *editRoute) convertStructToJSON() {
	if eR.err != nil {
		return
	}
	//レスポンス作成
	jsonEnc, err := json.Marshal(eR.routeModel)
	if err != nil {
		eR.err = customerr.BaseErr{
			Op:  "json marshaling simulRoute struct",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while json marshaling: %w", err),
		}
		return
	}
	//JSONのバイナリ形式のままだとtemplateで読み込めないので、stringに変換
	eR.routeJSON = string(jsonEnc)
}

func ShowAndEditSimulRoutesTpl(w http.ResponseWriter, req *http.Request) {
	var eR editRoute
	contexthandler.GetUserFromCtx(req, &eR.user, &eR.err)
	routeTitle := req.URL.Query().Get("route_title")
	d := eR.getRouteFromDB(routeTitle)
	bsonconv.DocToStruct(d, &eR.routeModel, &eR.err, "simul route")
	eR.convertStructToJSON()
	if eR.err != nil {
		e := eR.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage/simul_search/show_routes")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	data := contexthandler.GetLoginStateFromCtx(req)
	data["routeInfo"] = eR.routeJSON
	nineIterator := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	data["nineIterator"] = nineIterator
	showRouteTpl.ExecuteTemplate(w, "simul_search_show.html", data)
}

func init() {
	simulSearchTpl = template.Must(template.Must(template.ParseGlob("templates/simul_search/simul_search.html")).ParseGlob("templates/includes/*.html"))
}

func SimulSearchTpl(w http.ResponseWriter, req *http.Request) {
	msg := url.QueryEscape("エラーが発生しました。しばらく経ってからもう一度ご利用ください。")
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		log.Printf("Error whle gettibg data from context")
		http.Redirect(w, req, "/?msg="+msg, http.StatusInternalServerError)
		return
	}
	nineIterator := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	data["nineIterator"] = nineIterator
	simulSearchTpl.ExecuteTemplate(w, "simul_search.html", data)
}

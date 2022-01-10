package simulsearch

import (
	"app/bsonconv"
	"app/controllers"
	"app/cookiehandler"
	"app/customerr"
	"app/errormsg"
	"app/model"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type simulSearchController struct {
	controllers.Controller
	simulRoute model.SimulRoute
	user       model.User
}

func init() {
	showRouteTpl = template.Must(template.Must(template.ParseGlob("templates/simul_search/show_and_edit/simul_search_show.html")).ParseGlob("templates/includes/*.html"))
	simulSearchTpl = template.Must(template.Must(template.ParseGlob("templates/simul_search/simul_search.html")).ParseGlob("templates/includes/*.html"))
}

var simulSearchTpl *template.Template

//Index executes template for simul search page.
func Index(w http.ResponseWriter, req *http.Request) {
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		log.Printf("Error whle gettibg data from context")
		http.Redirect(w, req, "/?msg="+errormsg.TriAgain, http.StatusInternalServerError)
		return
	}
	nineIterator := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	data["nineIterator"] = nineIterator
	simulSearchTpl.ExecuteTemplate(w, "simul_search.html", data)
}

//getSimulRouteFromCtx gets simulRoute struct from request's context.
func (s *simulSearchController) getSimulRouteFromCtx(req *http.Request) {
	simulRouteFields, ok := req.Context().Value("simulRouteFields").(model.SimulRoute)
	if !ok {
		s.Err = customerr.BaseErr{
			Op:  "Get req routes info from context",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while getting request fields from reuest's context"),
		}
	}
	s.simulRoute = simulRouteFields
}

//saveRoute saves simulroute document to simulroutes collection.
func (s *simulSearchController) saveRoute() {
	err := s.simulRoute.SaveRoute(s.user.ID)
	if err != nil {
		s.Err = customerr.BaseErr{
			Op:  "Save new simul route",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while saving simul route: %w", err),
		}
	}
	return
}

//addRouteTitle adds title field and timestamp value to simul_route_titles field in users collection.
//The reason timestamp added to db is for sorting route titles in the descending order
//when showing them at mypage.
func (s *simulSearchController) addRouteTitle() {
	if s.Err != nil {
		return
	}
	now := time.Now().UTC() //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
	err := model.UpdateSimulRouteTitles(s.user.ID, s.simulRoute.Title, "$set", now)
	if err != nil {
		s.Err = customerr.BaseErr{
			Op:  "update user document's simul_route_titles field",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while updating user's simul_route_titles %w", err),
		}
	}
	return
}

func Save(w http.ResponseWriter, req *http.Request) {
	var s simulSearchController
	s.getSimulRouteFromCtx(req)
	s.GetUserFromCtx(req, &s.user)
	s.saveRoute()
	s.addRouteTitle()

	if s.Err != nil {
		e := s.Err.(customerr.BaseErr)
		http.Error(w, e.Msg, http.StatusInternalServerError)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	msgJSON := map[string]string{"msg": "OK"}
	json.NewEncoder(w).Encode(&msgJSON)
}

var showRouteTpl *template.Template

//getRouteFromDB gets route document from DB
func (s *simulSearchController) getRouteFromDB(title string) bson.M {
	if s.Err != nil {
		return nil
	}
	d, err := model.FindSimulRoute(s.user.ID, title)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.Err = customerr.BaseErr{
				Op:  "finding simul route document",
				Msg: "ご指定いただいたルートがありません。",
				Err: fmt.Errorf("error while finding route document from routes collection: %w", err),
			}
			return nil
		} else {
			s.Err = customerr.BaseErr{
				Op:  "Finding simul route document",
				Msg: errormsg.SomethingBad,
				Err: fmt.Errorf("error while finding route document from routes collection: %w", err),
			}
			return nil
		}
	}
	return d
}

//convertStructToJSON makes JSON object from simulRoute struct
func (s *simulSearchController) convertStructToJSON() string {
	if s.Err != nil {
		return ""
	}
	//レスポンス作成
	jsonEnc, err := json.Marshal(s.simulRoute)
	if err != nil {
		s.Err = customerr.BaseErr{
			Op:  "json marshaling simulRoute struct",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while json marshaling: %w", err),
		}
		return ""
	}
	//JSONのバイナリ形式のままだとtemplateで読み込めないので、stringに変換
	return string(jsonEnc)
}

func Show(w http.ResponseWriter, req *http.Request) {
	var s simulSearchController
	s.GetUserFromCtx(req, &s.user)
	routeTitle := req.URL.Query().Get("route_title")
	d := s.getRouteFromDB(routeTitle)
	bsonconv.DocToStruct(d, &s.simulRoute, &s.Err, "simul route")
	var routeJSON string
	routeJSON = s.convertStructToJSON()

	if s.Err != nil {
		e := s.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage/simul_search/show_routes")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	data := s.GetLoginStateFromCtx(req)
	data["routeInfo"] = routeJSON
	nineIterator := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	data["nineIterator"] = nineIterator
	showRouteTpl.ExecuteTemplate(w, "simul_search_show.html", data)
}

var previousTitle string

//getUpdateRouteFromCtx gets validated RouteUpdateRequest from context
func (s *simulSearchController) getUpdateRouteFromCtx(req *http.Request) {
	reqFields, ok := req.Context().Value("reqFields").(model.RouteUpdateRequest)
	if !ok {
		s.Err = customerr.BaseErr{
			Op:  "Get req routes info from context",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while getting request fields from reuest's context"),
		}
	}
	previousTitle = reqFields.PreviousTitle
	s.simulRoute = model.SimulRoute{
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
func (s *simulSearchController) updateRouteDoc() {
	if s.Err != nil {
		return
	}
	err := s.simulRoute.UpdateSimulRoute()
	if err != nil {
		var m string
		if strings.Contains(err.Error(), "(BSONObjectTooLarge)") {
			m = errormsg.TooLargeData
		} else {
			m = errormsg.SomethingBad
		}
		s.Err = customerr.BaseErr{
			Op:  "save new simul route",
			Msg: m,
			Err: fmt.Errorf("error while updating simul route: %w", err),
		}
	}
}

//updateRouteTitles updates timestamp value of simul_route_titles field in users collection,
//and delete previous title if title was changed.
//When you delete field in field (not field in document),
//you should use update operator $unset instead of delete.
//Official document URL: https://docs.mongodb.com/manual/reference/operator/update/unset/
func (s *simulSearchController) updateRouteTitles() {
	if s.Err != nil {
		return
	}
	now := time.Now().UTC() //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
	if s.simulRoute.Title != previousTitle {
		err := model.UpdateSimulRouteTitles(s.user.ID, previousTitle, "$unset", "")
		if err != nil {
			s.Err = customerr.BaseErr{
				Op:  "Remove previous simul route title",
				Msg: errormsg.SomethingBad,
				Err: fmt.Errorf("error while removing previous simul route title: %w", err),
			}
			return
		}
	}
	//「タイムスタンプを更新または追加」
	err := model.UpdateSimulRouteTitles(s.user.ID, s.simulRoute.Title, "$set", now)
	if err != nil {
		s.Err = customerr.BaseErr{
			Op:  "Set new simul route title and timestamp",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while setting new simul route title and timestamp: %w", err),
		}
		return
	}
}

//ルートを更新保存するための関数
func Update(w http.ResponseWriter, req *http.Request) {
	var s simulSearchController
	s.getUpdateRouteFromCtx(req)
	s.GetUserFromCtx(req, &s.user)
	s.updateRouteDoc()
	s.updateRouteTitles()

	if s.Err != nil {
		e := s.Err.(customerr.BaseErr)
		http.Error(w, e.Msg, http.StatusInternalServerError)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	respMsg := map[string]string{"msg": "ok"}
	json.NewEncoder(w).Encode(&respMsg)
}

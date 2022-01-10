package multiroute

import (
	"app/controllers"
	"app/internal/bsonconv"
	"app/internal/cookiehandler"
	"app/internal/customerr"
	"app/internal/errormsg"
	"app/internal/jsonconv"
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

type multiSearchController struct {
	controllers.Controller
	multiRoute model.MultiRoute
	user       model.User
}

var previousTitleAtUpdate = ""
var multiSearchTpl, showRouteTpl *template.Template

func init() {
	multiSearchTpl = template.Must(template.Must(template.ParseGlob("templates/multi_search/search/multi_search.html")).ParseGlob("templates/includes/*.html"))
	showRouteTpl = template.Must(template.Must(template.ParseGlob("templates/multi_search/show_and_edit/multi_route_show.html")).ParseGlob("templates/includes/*.html"))
}

/***
Methods related to Save
***/

//getMultiRouteFromCtx gets multiRoute from request's context
func (m *multiSearchController) getMultiRouteFromCtx(req *http.Request) {
	reqFields, ok := req.Context().Value("reqFields").(model.MultiRoute)
	if !ok {
		m.Err = customerr.BaseErr{
			Op:  "Get req routes info from context",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while getting request fields from reuest's context"),
		}
	}
	m.multiRoute = reqFields
}

//saveRoute saves route document to routes collection
func (m *multiSearchController) saveRoute() {
	err := m.multiRoute.SaveRoute(m.user.ID)
	if err != nil {
		e := customerr.BaseErr{
			Op:  "Save new multi route",
			Err: fmt.Errorf("error while saving multi route: %w", err),
		}

		if strings.Contains(err.Error(), "(BSONObjectTooLarge)") {
			e.Msg = errormsg.RouteDataTooLarge
			m.Err = e
		} else {
			e.Msg = errormsg.SomethingBad
			m.Err = e
		}
	}
}

//addRouteTitleToUserDoc adds multi_route_titles field to user document
func (m *multiSearchController) addRouteTitleToUserDoc() {
	if m.Err != nil {
		return
	}
	now := time.Now().UTC() //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
	err := model.UpdateMultiRouteTitles(m.user.ID, m.multiRoute.Title, "$set", now)
	if err != nil {
		m.Err = customerr.BaseErr{
			Op:  "update user document's multi_route_titles field",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while updating user's multi_route_titles %w", err),
		}
		return
	}
}

/***
Methods related to Update
***/

//getUpdateRouteFromCtx gets validated RouteUpdateRequest from context
func (m *multiSearchController) getUpdateRouteFromCtx(req *http.Request) {
	reqFields, ok := req.Context().Value("reqFields").(model.MultiRouteUpdateRequest)
	if !ok {
		m.Err = customerr.BaseErr{
			Op:  "Get req routes info from context",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while getting request fields from reuest's context"),
		}
	}
	m.multiRoute.ID = reqFields.ID
	m.multiRoute.Title = reqFields.Title
	m.multiRoute.Routes = reqFields.Routes
	previousTitleAtUpdate = reqFields.PreviousTitle
}

//updateRoute updates route document in routes collection
func (m *multiSearchController) updateRoute() {
	if m.Err != nil {
		return
	}
	err := m.multiRoute.UpdateRoute()
	if err != nil {
		e := customerr.BaseErr{
			Op:  "Save new multi route",
			Err: fmt.Errorf("error while updating multi route: %w", err),
		}

		if strings.Contains(err.Error(), "(BSONObjectTooLarge)") {
			e.Msg = errormsg.RouteDataTooLarge
			m.Err = e
		} else {
			e.Msg = errormsg.SomethingBad
			m.Err = e
		}
	}
}

//updateRouteTitles updates timestamp of multi_route_titles field in users collection, and delete previous title if title was changed.
func (m *multiSearchController) updateRouteTitles() {
	if m.Err != nil {
		return
	}
	now := time.Now().UTC() //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
	if m.multiRoute.Title != previousTitleAtUpdate {
		//「元のルート名をuser documentから削除」
		//documentではなく、document内のフィールドを削除する場合、Deleteではなく、Update operatorの$unsetを使って削除する
		//公式ドキュメントURL: https://docs.mongodb.com/manual/reference/operator/update/unset/
		err := model.UpdateMultiRouteTitles(m.user.ID, previousTitleAtUpdate, "$unset", "")
		if err != nil {
			m.Err = customerr.BaseErr{
				Op:  "Remove previous multi route title",
				Msg: errormsg.SomethingBad,
				Err: fmt.Errorf("error while removing previous multi route title: %w", err),
			}
			return
		}
	}
	//「タイムスタンプを更新または追加」
	err := model.UpdateMultiRouteTitles(m.user.ID, m.multiRoute.Title, "$set", now)
	if err != nil {
		m.Err = customerr.BaseErr{
			Op:  "Set new multi route title and timestamp",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while setting new multi route title and timestamp: %w", err),
		}
		return
	}
}

//getRouteFromDB gets route document from DB
func (m *multiSearchController) getRouteFromDB(title string) bson.M {
	if m.Err != nil {
		return nil
	}
	d, err := model.FindRoute(m.user.ID, title)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.Err = customerr.BaseErr{
				Op:  "Finding route document",
				Msg: "ご指定いただいたルートがありません。",
				Err: fmt.Errorf("error while finding route document from routes collection: %w", err),
			}
			return nil
		} else {
			m.Err = customerr.BaseErr{
				Op:  "Finding route document",
				Msg: errormsg.SomethingBad,
				Err: fmt.Errorf("error while finding route document from routes collection: %w", err),
			}
			return nil
		}
	}
	return d
}

/***
Method related to Delete
***/

//deleteRoute delete multi_route_title field from user document
func (m *multiSearchController) deleteRoute(routeTitle string) {
	if m.Err != nil {
		return
	}
	err := model.UpdateMultiRouteTitles(m.user.ID, routeTitle, "$unset", "")
	if err != nil {
		m.Err = customerr.BaseErr{
			Op:  "Deleting route title",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while deleting %v from multi_route_titles: %w", routeTitle, err),
		}
		return
	}
}

func Index(w http.ResponseWriter, req *http.Request) {
	msg := "エラーが発生しました。もう一度操作を行ってください。"
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		http.Redirect(w, req, "/mypage/show_routes/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting data from context: %v", ok)
		return
	}
	multiSearchTpl.ExecuteTemplate(w, "multi_search.html", data)
}

func Save(w http.ResponseWriter, req *http.Request) {
	var r multiSearchController
	r.getMultiRouteFromCtx(req)
	r.GetUserFromCtx(req, &r.user)
	r.saveRoute()
	/*users collectionのmulti_route_titlesフィールドにルート名と作成時刻を追加($set)する。
	  作成時刻はルート名取得時に作成時刻でソートするため*/
	r.addRouteTitleToUserDoc()

	if r.Err != nil {
		e := r.Err.(customerr.BaseErr)
		http.Error(w, e.Msg, http.StatusInternalServerError)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	msgJSON := map[string]string{"msg": "OK"}
	json.NewEncoder(w).Encode(&msgJSON)
}

//ルートを更新保存するための関数
func Update(w http.ResponseWriter, req *http.Request) {
	var m multiSearchController
	m.getUpdateRouteFromCtx(req)
	m.GetUserFromCtx(req, &m.user)
	m.updateRoute()
	m.updateRouteTitles()

	if m.Err != nil {
		e := m.Err.(customerr.BaseErr)
		http.Error(w, e.Msg, http.StatusInternalServerError)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	respMsg := map[string]string{"msg": "ok"}
	json.NewEncoder(w).Encode(&respMsg)
}

func Show(w http.ResponseWriter, req *http.Request) {
	var m multiSearchController
	m.GetUserFromCtx(req, &m.user)
	routeTitle := req.URL.Query().Get("route_title")
	d := m.getRouteFromDB(routeTitle)
	bsonconv.DocToStruct(d, &m.multiRoute, &m.Err, "multi route")
	routeJSON := jsonconv.StructToJSON(m.multiRoute, &m.Err)
	if m.Err != nil {
		e := m.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage/show_routes")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	data := m.GetLoginStateFromCtx(req)
	data["routeInfo"] = routeJSON
	showRouteTpl.ExecuteTemplate(w, "multi_route_show.html", data)
}

func Delete(w http.ResponseWriter, req *http.Request) {
	var m multiSearchController
	controllers.CheckHTTPMethod(req, &m.Err)
	m.GetUserFromCtx(req, &m.user)
	routeTitle := req.FormValue("title")
	m.deleteRoute(routeTitle)

	//レスポンス作成
	if m.Err != nil {
		e := m.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage/show_routes/?search_type=multi_search")
		log.Printf("Error while deleting route title: %v", e.Err)
		return
	}

	successMsg := "ルート「" + routeTitle + "」を削除しました。"
	cookiehandler.MakeCookieAndRedirect(w, req, "success", successMsg, "/mypage/show_routes/?search_type=multi_search")
	log.Printf("User [%v] deleted route [%v]", m.user.UserName, routeTitle)
}

////user collectionから削除 エラー解析に使うかもしれないので、route自体は削除せずに残しておく

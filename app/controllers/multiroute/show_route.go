package multiroute

import (
	"app/controllers"
	"app/cookiehandler"
	"app/customerr"
	"app/model"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var multiSearchTpl, showRouteTpl *template.Template

func init() {
	multiSearchTpl = template.Must(template.Must(template.ParseGlob("templates/multi_search/search/multi_search.html")).ParseGlob("templates/includes/*.html"))
	showRouteTpl = template.Must(template.Must(template.ParseGlob("templates/multi_search/show_and_edit/multi_route_show.html")).ParseGlob("templates/includes/*.html"))
}

func MultiSearchTpl(w http.ResponseWriter, req *http.Request) {
	msg := "エラーが発生しました。もう一度操作を行ってください。"
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		http.Redirect(w, req, "/mypage/show_routes/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting data from context: %v", ok)
		return
	}
	multiSearchTpl.ExecuteTemplate(w, "multi_search.html", data)
}

type editRoute struct {
	user       model.User
	routeModel model.MultiRoute
	routeJSON  string
	err        error
}

//getRouteFromDB gets route document from DB
func (eR *editRoute) getRouteFromDB(title string) bson.M {
	if eR.err != nil {
		return nil
	}
	d, err := model.FindRoute(eR.user.ID, title)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			eR.err = customerr.BaseErr{
				Op:  "Finding route document",
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

//convertStructToJSON makes JSON object from multiRoute struct
func (eR *editRoute) convertStructToJSON() {
	if eR.err != nil {
		return
	}
	//レスポンス作成
	jsonEnc, err := json.Marshal(eR.routeModel)
	if err != nil {
		eR.err = customerr.BaseErr{
			Op:  "json marshaling multiRoute struct",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while json marshaling: %w", err),
		}
		return
	}
	//JSONのバイナリ形式のままだとtemplateで読み込めないので、stringに変換
	eR.routeJSON = string(jsonEnc)
}

//getDataFromCtx gets data for executing template from request's context
func (eR *editRoute) getDataFromCtx(req *http.Request) map[string]interface{} {
	if eR.err != nil {
		return nil
	}
	d, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		eR.err = customerr.BaseErr{
			Op:  "Get data map from context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting data from context"),
		}
		return nil
	}
	return d
}

func ShowAndEditRoutesTpl(w http.ResponseWriter, req *http.Request) {
	var eR editRoute
	eR.user, eR.err = controllers.GetUserFromCtx(req)
	routeTitle := req.URL.Query().Get("route_title")
	d := eR.getRouteFromDB(routeTitle)
	eR.err = controllers.ConvertDucToStruct(d, &eR.routeModel, "multi route")
	eR.convertStructToJSON()
	//contextからデータ取得
	data := eR.getDataFromCtx(req)

	if eR.err != nil {
		e := eR.err.(customerr.BaseErr)

		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage/show_routes")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	data["routeInfo"] = eR.routeJSON
	showRouteTpl.ExecuteTemplate(w, "multi_route_show.html", data)
}

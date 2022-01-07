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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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

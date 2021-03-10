package reqvalidator

import (
	"app/controllers"
	"app/controllers/multiroute"
	"app/customerr"
	"app/model"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type multiRouteValidator struct {
	err error
}

//checkHTTPMethod checks request's Content-Type and HTTP methods
func (m *multiRouteValidator) checkHTTPMethod(req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" || req.Method != "POST" {
		m.err = customerr.BaseErr{
			Op:  "check HTTP method",
			Msg: "HTTPメソッドが不正です。",
			Err: fmt.Errorf("invalid HTTP method access"),
		}
		return
	}
}

//checkContainedChar checks if routeTitle contains . or $
func (m *multiRouteValidator) checkContainedChar(title string) {
	if m.err != nil {
		return
	}
	if strings.ContainsAny(title, ".$") {
		m.err = customerr.BaseErr{
			Op:  "check contained character in route's title",
			Msg: "ルート名にご使用いただけない文字が含まれています。",
			Err: fmt.Errorf(". or $ was contained in route title"),
		}
		return
	}
}

//convertJSONToStruct converts request's JSON data to multiRoute struct
func (m *multiRouteValidator) convertJSONToStruct(req *http.Request, reqFields interface{}) {
	if m.err != nil {
		return
	}
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, reqFields)
	if err != nil {
		m.err = customerr.BaseErr{
			Op:  "json unmarshal multi route request",
			Msg: "入力に不正があります。",
			Err: fmt.Errorf("error while json unmarshaling multiroute request: %w", err),
		}
		return
	}
}

func SaveRoutesValidator(SaveRoutes http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var m multiRouteValidator
		m.checkHTTPMethod(req)
		//convertJSONToStructの第２引数はinterfaceなので、変数を宣言してポインタを渡す必要がある
		var reqFields model.MultiRoute
		m.checkContainedChar(reqFields.Title)
		m.convertJSONToStruct(req, &reqFields)

		if m.err != nil {
			e := m.err.(customerr.BaseErr)
			http.Error(w, e.Msg, http.StatusBadRequest)
			log.Printf("operation: %s, error: %v", e.Op, e.Err)
			return
		}

		//contextに各フィールドの値を追加
		ctx := req.Context()
		ctx = context.WithValue(ctx, "reqFields", reqFields)
		SaveRoutes.ServeHTTP(w, req.WithContext(ctx))
	}
}

func UpdateRouteValidator(UpdateRoute http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var m multiRouteValidator
		controllers.CheckHTTPMethod(req, &m.err)
		//convertJSONToStructの第２引数はinterfaceなので、変数を宣言してポインタを渡す必要がある
		var reqFields multiroute.RouteUpdateRequest
		m.convertJSONToStruct(req, &reqFields)
		m.checkContainedChar(reqFields.Title)

		if m.err != nil {
			e := m.err.(customerr.BaseErr)
			http.Error(w, e.Msg, http.StatusBadRequest)
			log.Printf("operation: %s, error: %v", e.Op, e.Err)
			return
		}

		//contextに各フィールドの値を追加
		ctx := req.Context()
		ctx = context.WithValue(ctx, "reqFields", reqFields)
		UpdateRoute.ServeHTTP(w, req.WithContext(ctx))
	}
}

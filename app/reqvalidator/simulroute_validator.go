package reqvalidator

import (
	"app/internal/customerr"
	"app/model"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type simulRouteValidator struct {
	err error
}

//checkHTTPMethod checks request's Content-Type and HTTP methods
func (s *simulRouteValidator) checkHTTPMethod(req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" || req.Method != "POST" {
		s.err = customerr.BaseErr{
			Op:  "check HTTP method",
			Msg: "HTTPメソッドが不正です。",
			Err: fmt.Errorf("invalid HTTP method access"),
		}
		return
	}
}

//checkContainedChar checks if routeTitle contains . or $
func (s *simulRouteValidator) checkContainedChar(title string) {
	if s.err != nil {
		return
	}
	if strings.ContainsAny(title, ".$") {
		s.err = customerr.BaseErr{
			Op:  "check contained character in route's title",
			Msg: "ルート名にご使用いただけない文字が含まれています。",
			Err: fmt.Errorf(". or $ was contained in route title"),
		}
		return
	}
}

//convertJSONToStruct converts request's JSON data to multiRoute struct
func (s *simulRouteValidator) convertJSONToStruct(req *http.Request, reqFields interface{}) {
	if s.err != nil {
		return
	}
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, reqFields)
	if err != nil {
		s.err = customerr.BaseErr{
			Op:  "json unmarshal multi route request",
			Msg: "入力に不正があります。",
			Err: fmt.Errorf("error while json unmarshaling multiroute request: %w", err),
		}
		return
	}
}

func SaveSimulRouteValidator(SaveRoutes http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var s simulRouteValidator
		s.checkHTTPMethod(req)
		//convertJSONToStructの第２引数はinterfaceなので、変数を宣言してポインタを渡す必要がある
		var reqFields model.SimulRoute
		s.convertJSONToStruct(req, &reqFields)
		s.checkContainedChar(reqFields.Title)

		if s.err != nil {
			e := s.err.(customerr.BaseErr)
			http.Error(w, e.Msg, http.StatusBadRequest)
			log.Printf("operation: %s, error: %v", e.Op, e.Err)
			return
		}

		//contextに各フィールドの値を追加
		ctx := req.Context()
		ctx = context.WithValue(ctx, "simulRouteFields", reqFields)
		SaveRoutes.ServeHTTP(w, req.WithContext(ctx))
	}
}

func UpdateSimulRouteValidator(Update http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var s simulRouteValidator
		//convertJSONToStructの第２引数はinterfaceなので、変数を宣言してポインタを渡す必要がある
		var reqFields model.RouteUpdateRequest
		s.convertJSONToStruct(req, &reqFields)
		s.checkContainedChar(reqFields.Title)

		if s.err != nil {
			e := s.err.(customerr.BaseErr)
			http.Error(w, e.Msg, http.StatusBadRequest)
			log.Printf("operation: %s, error: %v", e.Op, e.Err)
			return
		}

		//contextに各フィールドの値を追加
		ctx := req.Context()
		ctx = context.WithValue(ctx, "reqFields", reqFields)
		Update.ServeHTTP(w, req.WithContext(ctx))
	}
}

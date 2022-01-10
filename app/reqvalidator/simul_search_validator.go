package reqvalidator

import (
	"app/controllers/simulsearch"
	"app/internal/customerr"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type simulSearchValidator struct {
	reqParams simulsearch.SimulParams
	err       error
}

//checkHTTPMethod checks request's Content-Type and HTTP methods
func (s *simulSearchValidator) checkHTTPMethod(req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" || req.Method != "POST" {
		s.err = customerr.BaseErr{
			Op:  "check HTTP method",
			Msg: "HTTPメソッドが不正です。",
			Err: fmt.Errorf("invalid HTTP method access"),
		}
		return
	}
}

//convertJSONToStruct converts request's JSON data to simulSearch struct
func (s *simulSearchValidator) convertJSONToStruct(req *http.Request) {
	if s.err != nil {
		return
	}
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, &s.reqParams)
	if err != nil {
		s.err = customerr.BaseErr{
			Op:  "json unmarshal multi route request",
			Msg: "入力に不正があります。",
			Err: fmt.Errorf("error while json unmarshaling multiroute request: %w", err),
		}
		return
	}
}

//checkAndModifyOrigin checks origin length and add prefix(place_id:)
func (s *simulSearchValidator) checkAndModifyOrigin() {
	if s.err != nil {
		return
	}
	if s.reqParams.Origin == "" {
		s.err = customerr.BaseErr{
			Op:  "checking origin input of simul search",
			Msg: "出発地を入力してください。",
			Err: fmt.Errorf("origin was empty at simul search"),
		}
		return
	}
	//place_id:を追加
	s.reqParams.Origin = "place_id:" + s.reqParams.Origin
}

//addPrefixToDestinations adds prefix(place_id:) to destinations
func (s *simulSearchValidator) addPrefixToDestinations() {
	if s.err != nil {
		return
	}
	for i := 1; i < 10; i++ {
		if s.reqParams.Destinations[strconv.Itoa(i)] == "" {
			continue
		}
		s.reqParams.Destinations[strconv.Itoa(i)] = "place_id:" + s.reqParams.Destinations[strconv.Itoa(i)]
	}

}

func SimulSearchValidator(Search http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var s simulSearchValidator
		s.checkHTTPMethod(req)
		s.convertJSONToStruct(req)
		s.checkAndModifyOrigin()
		s.addPrefixToDestinations()

		if s.err != nil {
			e := s.err.(customerr.BaseErr)
			http.Error(w, e.Msg, http.StatusBadRequest)
			log.Printf("operation: %s, error: %v", e.Op, e.Err)
			return
		}

		//contextに各フィールドの値を追加
		ctx := req.Context()
		ctx = context.WithValue(ctx, "reqParams", s.reqParams)
		Search.ServeHTTP(w, req.WithContext(ctx))
	}
}

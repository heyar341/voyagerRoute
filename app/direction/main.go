package direction

import (
	"googlemaps.github.io/maps"
	"net/http"
	"strings"
)

var Req = &maps.DirectionsRequest{
	Language:      "ja",
	Region:        "JP",
}
func SearchRoutes(res http.ResponseWriter, req *http.Request){
	//必須パラメータ
	Req.Origin = req.FormValue("origin")
	Req.Destination = req.FormValue("destination")
	Req.DepartureTime = req.FormValue("departureTime")
	Req.ArrivalTime = req.FormValue("arrivalTime")
	//オプションパラメータ
	if req.FormValue("mode") != "" { lookupMode(req.FormValue("mode"), Req) }
	if req.FormValue("waypoints") != "" { Req.Waypoints = strings.Split(req.FormValue("waypoints"), "|") }
	if req.FormValue("alternatives") != "" { lookupAlternatives(req.FormValue("alternatives"), Req) }
	if req.FormValue("transitRoutingPreference") != "" {lookupTransitRoutingPreference(req.FormValue("transitRoutingPreference"), Req)}
	if req.FormValue("trafficModel") != "" { lookupTrafficModel(req.FormValue("trafficModel"), Req) }
	if req.FormValue("avoid") != "" { lookupAvoid(req.FormValue("avoid"), Req) }
	if req.FormValue("transitMode") != "" { lookupTransitMode(req.FormValue("transitMode"), Req) }
}
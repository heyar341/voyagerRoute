package simuloptions

import (
	"googlemaps.github.io/maps"
	"log"
	"strings"
)

//徒歩、運転、乗り換えのモード選択
func LookupMode(mode string, r *maps.DirectionsRequest) {
	switch mode {
	case "driving":
		r.Mode = maps.TravelModeDriving
	case "walking":
		r.Mode = maps.TravelModeWalking
	case "bicycling":
		r.Mode = maps.TravelModeBicycling
	case "transit":
		r.Mode = maps.TravelModeTransit
	case "":
		// ignore
	default:
		log.Printf("Unknown mode '%s'", mode)
	}
}

//trueの場合複数ルートを探し、falseの場合１つのルートのみ返す
func LookupAlternatives(alternatives string, r *maps.DirectionsRequest) {
	if alternatives == "true" {
		r.Alternatives = true
	} else {
		r.Alternatives = false
	}
}

//乗り換え数の少なさを優先するか、歩行距離の短さを優先するか選択
func LookupTransitRoutingPreference(transitRoutingPreference string, r *maps.DirectionsRequest) {
	switch transitRoutingPreference {
	case "fewer_transfers":
		r.TransitRoutingPreference = maps.TransitRoutingPreferenceFewerTransfers
	case "less_walking":
		r.TransitRoutingPreference = maps.TransitRoutingPreferenceLessWalking
	case "":
		// ignore
	default:
		log.Printf("Unknown transit routing preference %s", transitRoutingPreference)
	}
}

//最速時間、過去のデータからの最適予測時間、最も遅い場合の予測のどれか選択
func LookupTrafficModel(trafficModel string, r *maps.DirectionsRequest) {
	switch trafficModel {
	case "optimistic":
		r.TrafficModel = maps.TrafficModelOptimistic
	case "best_guess":
		r.TrafficModel = maps.TrafficModelBestGuess
	case "pessimistic":
		r.TrafficModel = maps.TrafficModelPessimistic
	case "":
		// ignore
	default:
		log.Printf("Unknown traffic mode %s", trafficModel)
	}
}

//有料道路、高速道路、フェリーを除外する場合選択
func LookupAvoid(avoid string, r *maps.DirectionsRequest) {
	for _, a := range strings.Split(avoid, "|") {
		switch a {
		case "tolls":
			r.Avoid = append(r.Avoid, maps.AvoidTolls)
		case "highways":
			r.Avoid = append(r.Avoid, maps.AvoidHighways)
		case "ferries":
			r.Avoid = append(r.Avoid, maps.AvoidFerries)
		default:
			log.Printf("Unknown avoid restriction %s", a)
		}
	}
}

//交通手段を選択
func LookupTransitMode(transitMode string, r *maps.DirectionsRequest) {
	for _, t := range strings.Split(transitMode, "|") {
		switch t {
		case "bus":
			r.TransitMode = append(r.TransitMode, maps.TransitModeBus)
		case "subway":
			r.TransitMode = append(r.TransitMode, maps.TransitModeSubway)
		case "train":
			r.TransitMode = append(r.TransitMode, maps.TransitModeTrain)
		case "tram":
			r.TransitMode = append(r.TransitMode, maps.TransitModeTram)
		case "rail":
			r.TransitMode = append(r.TransitMode, maps.TransitModeRail)
		}
	}
}

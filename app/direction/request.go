package direction

import (
	"googlemaps.github.io/maps"
	"log"
	"strings"
)

func lookupMode(mode string, r *maps.DirectionsRequest) {
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
		log.Fatalf("Unknown mode '%s'", mode)
	}
}

func lookupAlternatives(alternatives string, r *maps.DirectionsRequest)  {
	if alternatives == "true" {
		r.Alternatives = true
	} else {
		r.Alternatives = false
	}
}
func lookupTransitRoutingPreference(transitRoutingPreference string, r *maps.DirectionsRequest) {
	switch transitRoutingPreference {
	case "fewer_transfers":
		r.TransitRoutingPreference = maps.TransitRoutingPreferenceFewerTransfers
	case "less_walking":
		r.TransitRoutingPreference = maps.TransitRoutingPreferenceLessWalking
	case "":
		// ignore
	default:
		log.Fatalf("Unknown transit routing preference %s", transitRoutingPreference)
	}
}

func lookupTrafficModel(trafficModel string, r *maps.DirectionsRequest) {
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
		log.Fatalf("Unknown traffic mode %s", trafficModel)
	}
}

func lookupAvoid(avoid string, r *maps.DirectionsRequest)  {
	for _, a := range strings.Split(avoid, "|") {
		switch a {
		case "tolls":
			r.Avoid = append(r.Avoid, maps.AvoidTolls)
		case "highways":
			r.Avoid = append(r.Avoid, maps.AvoidHighways)
		case "ferries":
			r.Avoid = append(r.Avoid, maps.AvoidFerries)
		default:
			log.Fatalf("Unknown avoid restriction %s", a)
		}
	}
}
func lookupTransitMode(transitMode string, r *maps.DirectionsRequest)  {
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
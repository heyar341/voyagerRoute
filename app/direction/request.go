package direction

import (
	"flag"
	"googlemaps.github.io/maps"
	"log"
	"strings"
)

var (
	origin                   = flag.String("origin", "", "The address or textual latitude/longitude value from which you wish to calculate directions.")
	destination              = flag.String("destination", "", "The address or textual latitude/longitude value from which you wish to calculate directions.")
	mode                     = flag.String("mode", "", "The travel mode for this directions request.")
	departureTime            = flag.String("departure_time", "", "The depature time for transit mode directions request.")
	arrivalTime              = flag.String("arrival_time", "", "The arrival time for transit mode directions request.")
	waypoints                = flag.String("waypoints", "", "The waypoints for driving directions request, | separated.")
	alternatives             = flag.Bool("alternatives", false, "Whether the Directions service may provide more than one route alternative in the response.")
	avoid                    = flag.String("avoid", "", "Indicates that the calculated route(s) should avoid the indicated features, | separated.")
	language                 = flag.String("language", "", "Specifies the language in which to return results.")
	region                   = flag.String("region", "", "Specifies the region code, specified as a ccTLD (\"top-level domain\") two-character value.")
	transitMode              = flag.String("transit_mode", "", "Specifies one or more preferred modes of transit, | separated. This parameter may only be specified for transit directions.")
	transitRoutingPreference = flag.String("transit_routing_preference", "", "Specifies preferences for transit routes.")
	trafficModel             = flag.String("traffic_model", "", "Specifies traffic prediction model when request future directions. Valid values are optimistic, best_guess, and pessimistic. Optional.")
)
var Req = &maps.DirectionsRequest{
	Origin:        *origin,
	Destination:   *destination,
	DepartureTime: *departureTime,
	ArrivalTime:   *arrivalTime,
	Alternatives:  *alternatives,
	Language:      *language,
	Region:        *region,
}
func MakeRequest() *maps.DirectionsRequest {

	
	lookupMode(*mode, Req)
	lookupTransitRoutingPreference(*transitRoutingPreference, Req)
	lookupTrafficModel(*trafficModel, Req)

	if *waypoints != "" {
		Req.Waypoints = strings.Split(*waypoints, "|")
	}

	if *avoid != "" {
		for _, a := range strings.Split(*avoid, "|") {
			switch a {
			case "tolls":
				Req.Avoid = append(Req.Avoid, maps.AvoidTolls)
			case "highways":
				Req.Avoid = append(Req.Avoid, maps.AvoidHighways)
			case "ferries":
				Req.Avoid = append(Req.Avoid, maps.AvoidFerries)
			default:
				log.Fatalf("Unknown avoid restriction %s", a)
			}
		}
	}
	if *transitMode != "" {
		for _, t := range strings.Split(*transitMode, "|") {
			switch t {
			case "bus":
				Req.TransitMode = append(Req.TransitMode, maps.TransitModeBus)
			case "subway":
				Req.TransitMode = append(Req.TransitMode, maps.TransitModeSubway)
			case "train":
				Req.TransitMode = append(Req.TransitMode, maps.TransitModeTrain)
			case "tram":
				Req.TransitMode = append(Req.TransitMode, maps.TransitModeTram)
			case "rail":
				Req.TransitMode = append(Req.TransitMode, maps.TransitModeRail)
			}
		}
	}
	return Req
}

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
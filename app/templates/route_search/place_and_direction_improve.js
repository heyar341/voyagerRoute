// This example requires the Places library. Include the libraries=places
// parameter when you first load the API. For example:
// <script
// src="https://maps.googleapis.com/maps/api/js?key=YOUR_API_KEY&libraries=places">
function initMap() {
    const map = new google.maps.Map(document.getElementById("map"), {
        mapTypeControl: false,
        center: { lat: 35.6816228, lng: 139.765199 },
        zoom: 10,
        //streetViewを無効化
        streetViewControl: false,
    });
    new AutocompleteDirectionsHandler(map);
}

class AutocompleteDirectionsHandler {
    constructor(map) {
        this.resp = {};
        this.map = map;
        this.originPlaceId = "";
        this.destinationPlaceId = "";
        this.travelMode = google.maps.TravelMode.WALKING;
        this.directionsService = new google.maps.DirectionsService();
        this.directionsRenderer = new google.maps.DirectionsRenderer();
        this.directionsRenderer.setMap(map);
        this.directionsRenderer.setPanel(document.getElementById("right-panel"));
        const originInput = document.getElementById("origin-input");
        const destinationInput = document.getElementById("destination-input");
        const modeSelector = document.getElementById("mode-selector");
        const originAutocomplete = new google.maps.places.Autocomplete(originInput);
        // Auto Completeを画面のmap上でない場所で使うときのドキュメントURL:
        //https://developers.google.com/maps/documentation/javascript/places-autocomplete?hl=en#places-searchbox
        //SearchBoxを使う場合、Detailまで返ってくるから、料金高くなる可能性あり
        // var input = document.getElementById('searchTextField');
        // var searchBox = new google.maps.places.SearchBox(input);
        // Specify just the place data fields that you need.
        originAutocomplete.setFields(["place_id"]);
        const destinationAutocomplete = new google.maps.places.Autocomplete(
            destinationInput
        );
        // Specify just the place data fields that you need.
        destinationAutocomplete.setFields(["place_id"]);
        this.setupClickListener(
            "changemode-walking",
            google.maps.TravelMode.WALKING
        );
        this.setupClickListener(
            "changemode-transit",
            google.maps.TravelMode.TRANSIT
        );
        this.setupClickListener(
            "changemode-driving",
            google.maps.TravelMode.DRIVING
        );
        this.setupPlaceChangedListener(originAutocomplete, "ORIG");
        this.setupPlaceChangedListener(destinationAutocomplete, "DEST");
    }
    // Sets a listener on a radio button to change the filter type on Places
    // Autocomplete.
    setupClickListener(id, mode) {
        const radioButton = document.getElementById(id);
        radioButton.addEventListener("click", () => {
            this.travelMode = mode;
            this.route();
        });
    }
    setupPlaceChangedListener(autocomplete, mode) {
        autocomplete.bindTo("bounds", this.map);
        autocomplete.addListener("place_changed", () => {
            const place = autocomplete.getPlace();
            if (!place.place_id) {
                window.alert("Please select an option from the dropdown list.");
                return;
            }

            if (mode === "ORIG") {
                this.originPlaceId = place.place_id;
            } else {
                this.destinationPlaceId = place.place_id;
            }
            this.route();
        });
    }
    route() {
        if (!this.originPlaceId || !this.destinationPlaceId) {
            return;
        }
        const me = this;
        this.directionsService.route(
            {
                origin: { placeId: this.originPlaceId },
                destination: { placeId: this.destinationPlaceId },
                travelMode: this.travelMode,
                //↓複数ルートを返す場合、指定
                // provideRouteAlternatives: true,
            },
            (response, status) => {
                if (status === "OK") {
                    //複数ルートがある場合、subRouteRendererで各ルートを薄く表示
                    if(response.routes.length > 1) {
                        for (var i = 0; i < response.routes.length; i++) {
                            //jsではObjectは参照渡しなので、Object.assignを使って、値渡しにする
                            var sub_res = Object.assign({}, response);
                            sub_res.routes = [response.routes[i]];
                            // 各ルートごとに、Renderオブジェクトを作成する必要がある
                            var subRouteRenderer = new google.maps.DirectionsRenderer();
                            //ドキュメントURL: https://developers.google.com/maps/documentation/javascript/reference/directions#DirectionsRendererOptions
                            subRouteRenderer.setOptions({
                                suppressMarkers: false,
                                suppressPolylines: false,
                                suppressInfoWindows: false,
                                //Colorとopacity(不透明度)と太さを設定
                                polylineOptions: {
                                    strokeColor: '#808080',
                                    strokeOpacity: 0.5,
                                    strokeWeight: 7
                                }
                            });
                            subRouteRenderer.setDirections(sub_res)
                            subRouteRenderer.setMap(this.map)
                        }
                    }
                    //responseをRendererに渡して、ルートを描画
                    //ドキュメントURL: https://developers.google.com/maps/documentation/javascript/reference/marker#MarkerOptions
                    me.directionsRenderer.setOptions({
                            suppressMarkers: false,
                            suppressPolylines: false,
                            suppressInfoWindows: false,
                            polylineOptions: {
                                strokeColor: '#00bfff',
                                strokeOpacity: 1.0,
                                strokeWeight: 7
                            }
                        });
                    this.resp = response;
                    console.log(response)
                    console.log(typeof response)
                    me.directionsRenderer.setDirections(response);
                } else {
                    window.alert("Directions request failed due to " + status);
                }
            }
        );
    }
}
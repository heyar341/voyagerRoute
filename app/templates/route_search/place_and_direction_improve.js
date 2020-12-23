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
    var today = new Date();
    var yyyy = today.getFullYear();
    var mm = ("0"+(today.getMonth()+1)).slice(-2);
    var dd = ("0"+today.getDate()).slice(-2);
    document.getElementById("date").value=yyyy+'-'+mm+'-'+dd;
    new AutocompleteDirectionsHandler(map);
}

class AutocompleteDirectionsHandler {
    constructor(map) {
        this.resp = {};
        this.map = map;
        this.directionsRequest = {};
        this.originPlaceId = "";
        this.destinationPlaceId = "";
        this.transitTime = new Date();
        this.travelMode = google.maps.TravelMode.WALKING;
        this.directionsService = new google.maps.DirectionsService();
        this.directionsRenderer = new google.maps.DirectionsRenderer();
        this.directionsRenderer.setMap(map);
        this.directionsRenderer.setPanel(document.getElementById("right-panel"));
        const originInput = document.getElementById("origin-input");
        const destinationInput = document.getElementById("destination-input");
        const originAutocomplete = new google.maps.places.Autocomplete(originInput);
        // Specify just the place data fields that you need.
        originAutocomplete.setFields(["place_id"]);
        const destinationAutocomplete = new google.maps.places.Autocomplete(
            destinationInput
        );
        // Specify just the place data fields that you need.
        destinationAutocomplete.setFields(["place_id"]);
        //EventListenerの設定
        this.setupClickListener("changemode-walking",google.maps.TravelMode.WALKING);
        this.setupClickListener("changemode-transit", google.maps.TravelMode.TRANSIT);
        this.setupClickListener("changemode-driving", google.maps.TravelMode.DRIVING);
        this.setupPlaceChangedListener(originAutocomplete, "ORIG");
        this.setupPlaceChangedListener(destinationAutocomplete, "DEST");
        this.setupOptionListener("date");
        this.setupOptionListener("time");
        this.setupOptionListener("avoid-toll");
        this.setupOptionListener("avoid-highway");
        

    }

    //経路オプションのラジオボタンが押されたら発火
    setupClickListener(id, mode) {
        const radioButton = document.getElementById(id);
        radioButton.addEventListener("click", () => {
            if(id == "changemode-transit"){
                document.getElementById("transit-time").style.display = "block"
            }
            else if(id != "changemode-transit"){
                document.getElementById("transit-time").style.display = "none"
            }
            if(id == "changemode-driving"){
                document.getElementById("driving-option").style.display = "block"
            }
            else if(id != "changemode-driving"){
                document.getElementById("driving-option").style.display = "none"
            }
            this.travelMode = mode;
            this.route();
        });
    }
    //出発地と目的地の入力があった場合、発火
    setupPlaceChangedListener(autocomplete, mode) {
        autocomplete.bindTo("bounds", this.map);
        autocomplete.addListener("place_changed", () => {
            const place = autocomplete.getPlace();
            if (!place.place_id) {
                window.alert("表示された選択肢の中から選んでください。");
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
    setupOptionListener(id) {
        const optionChange = document.getElementById(id);
        optionChange.addEventListener("change", ()=>{
           this.route();
        });
    }
    
    route() {
        if (!this.originPlaceId || !this.destinationPlaceId) {
            return;
        }
        const me = this;
        this.directionsRequest =
            {
            origin: { placeId: this.originPlaceId },
            destination: { placeId: this.destinationPlaceId },
            travelMode: this.travelMode,
            //↓複数ルートを返す場合、指定
            // provideRouteAlternatives: true,
            }
        if(document.getElementById("specify-route-options").checked){
            if(document.getElementById("changemode-transit").checked){
                this.directionsRequest.transitOptions = {}
                if(document.getElementById("depart-time").checked){
                    console.log(new Date(document.getElementById("date").value +
                        "T" +    document.getElementById("time").value));
                    console.log(document.getElementById("time").value);
                    this.directionsRequest.transitOptions.departureTime =
                        new Date(document.getElementById("date").value +
                        "T" +    document.getElementById("time").value);
                }
                else if(document.getElementById("arrival-time").checked){
                    this.directionsRequest.transitOptions.arrivalTime =
                        new Date(document.getElementById("date").value +
                        "T" +    document.getElementById("time").value);
                }
                this.directionsRequest.transitOptions.routingPreference = 'FEWER_TRANSFERS';
            }
            else if(document.getElementById("changemode-driving").checked){
                if(document.getElementById("avoid-toll").checked){
                    this.directionsRequest.avoidTolls = true;
                }
                if(document.getElementById("avoid-highway").checked){
                    this.directionsRequest.avoidHighways = true;
                }
            }
        }
        this.directionsService.route(
            this.directionsRequest,
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
                    me.directionsRenderer.setDirections(response);

                } else {
                    window.alert("検索結果はありません");
                }
            }
        );
    }
}
//ルート番号
var routeID = 0;
//ルート描画時の線の色を指定
const colorMap = {0:"#00bfff",1:"#c8e300",2:"#9543de",3:"#00db30",4:"#b09856",5:"#00deda",6:"#eb86d5",7:"#83b300",8:"#ffb300",9:"#de0000"}

function initMap() {
    const map = new google.maps.Map(document.getElementById("map"), {
        mapTypeControl: false,
        center: { lat: 35.6816228, lng: 139.765199 },
        zoom: 10,
        //streetViewを無効化
        streetViewControl: false,
    });
    //現在時刻を取得し、時間指定要素に挿入
    var today = new Date();
    var yyyy = today.getFullYear();
    var mm = ("0"+(today.getMonth()+1)).slice(-2);
    var dd = ("0"+today.getDate()).slice(-2);

    //１番目のルート要素をHTMLに追加
    $('#search-box').append(genSearchBox(routeID,colorMap[routeID]));
    document.getElementById("date" + String(routeID)).value=yyyy+'-'+mm+'-'+dd;
    document.getElementById("date" + String(routeID)).min=yyyy+'-'+mm+'-'+dd;
    //AutocompleteとDiretionsServiceのインスタンス化
    new AutocompleteDirectionsHandler(map,String(routeID));
    $('.toggle-title').on('click', function () {
        $(this).toggleClass('active');
        $(this).next().slideToggle();
    });

    //ボタンが押されたら２番目以降のルート要素をHTMLに追加
    $("#add-route").on("click", function () {
        routeID++;
        $('#search-box').append(genSearchBox(routeID,colorMap[routeID]));
        document.getElementById("date" + String(routeID)).value=yyyy+'-'+mm+'-'+dd;
        document.getElementById("date" + String(routeID)).min=yyyy+'-'+mm+'-'+dd;
        new AutocompleteDirectionsHandler(map,String(routeID));
        $('.toggle-title').on('click', function () {
            $('.toggle-title').off("click");
            $('.toggle-title').on('click', function () {
                //activeの場合、CSSで下矢印を表示する
                $(this).toggleClass('active');
                $(this).next().slideToggle();
            });
        });
    });
}

class AutocompleteDirectionsHandler {
    constructor(map,routeNum) {
        this.routeNum = routeNum;
        this.colorCode = colorMap[parseInt(routeNum)];
        this.map = map;
        this.directionsRequest = {};
        this.originPlaceId = "";
        this.destinationPlaceId = "";
        this.poly = [];
        this.travelMode = google.maps.TravelMode.TRANSIT;
        this.directionsService = new google.maps.DirectionsService();
        this.directionsRenderer = new google.maps.DirectionsRenderer();
        this.directionsRenderer.setMap(map);
        this.directionsRenderer.setPanel(document.getElementById("route-detail-panel" + this.routeNum));
        const originInput = document.getElementById("origin-input" + this.routeNum);
        const destinationInput = document.getElementById("destination-input" + this.routeNum);
        const originAutocomplete = new google.maps.places.Autocomplete(originInput);
        // Specify just the place data fields that you need.
        originAutocomplete.setFields(["place_id"]);
        const destinationAutocomplete = new google.maps.places.Autocomplete(
            destinationInput
        );
        // Specify just the place data fields that you need.
        destinationAutocomplete.setFields(["place_id"]);
        //EventListenerの設定
        this.setupClickListener("changemode-walking" + this.routeNum,google.maps.TravelMode.WALKING);
        this.setupClickListener("changemode-transit" + this.routeNum, google.maps.TravelMode.TRANSIT);
        this.setupClickListener("changemode-driving" + this.routeNum, google.maps.TravelMode.DRIVING);
        this.setupPlaceChangedListener(originAutocomplete, "ORIG");
        this.setupPlaceChangedListener(destinationAutocomplete, "DEST");
        this.setupOptionListener("date" + this.routeNum);
        this.setupOptionListener("time" + this.routeNum);
        this.setupOptionListener("avoid-toll" + this.routeNum);
        this.setupOptionListener("avoid-highway" + this.routeNum);
        this.setUpRouteSelectedListener(this,this.directionsRenderer);

    }

    //経路オプションのラジオボタンが押されたら発火
    setupClickListener(id, mode) {
        const radioButton = document.getElementById(id);
        radioButton.addEventListener("click", () => {
            if(id == "changemode-transit" + this.routeNum){
                document.getElementById("transit-time" + this.routeNum).style.display = "block"
            }
            else if(id != "changemode-transit" + this.routeNum){
                document.getElementById("transit-time" + this.routeNum).style.display = "none"
            }
            if(id == "changemode-driving" + this.routeNum){
                document.getElementById("driving-option" + this.routeNum).style.display = "block"
            }
            else if(id != "changemode-driving" + this.routeNum){
                document.getElementById("driving-option" + this.routeNum).style.display = "none"
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
        optionChange.addEventListener("change", ()=> {
                this.route();
            });
    }

    //複数ルートがある場合、パネルのルートを押したら発火
    setUpRouteSelectedListener(obj,directionsRenderer) {
        google.maps.event.addListener(directionsRenderer, 'routeindex_changed', function () {
            // console.log(directionsRenderer.directions);
            console.log(directionsRenderer);

            var target = directionsRenderer.getRouteIndex();
            for(var i = 0; i < obj.poly.length; i++){
                if(i == target){
                    obj.poly[i].setOptions({
                        //選択したルートの場合、色をcolorCodeに従って変更
                        polylineOptions: {
                                strokeColor: obj.colorCode,
                                strokeOpacity: 1.0,
                                strokeWeight: 7,
                            //青ラインを一番上に表示するため、zIndexを他のルートより大きくする。
                            zIndex: 1
                        }
                    });
                }
                else{
                    obj.poly[i].setOptions({
                        //選択したルート以外の場合、色を#808080に設定(選択されている場合青だから、
                        //元に戻すためには、全てのルートについて#808080に設定する必要あり。)
                        polylineOptions: {
                            strokeColor: '#808080',
                            strokeOpacity: 0.7,
                            strokeWeight: 7,
                            //青ラインを一番上に表示するため、zIndexを最小にする
                            zIndex: 0
                        }
                    });
                }
                obj.poly[i].setMap(obj.map);
            }
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
            provideRouteAlternatives: true,
            }
            if(document.getElementById("changemode-transit" + this.routeNum).checked){
                this.directionsRequest.transitOptions = {}
                if(document.getElementById("depart-time" + this.routeNum).checked){
                    console.log(new Date(document.getElementById("date" + this.routeNum).value +
                        "T" +    document.getElementById("time" + this.routeNum).value));
                    console.log(document.getElementById("time" + this.routeNum).value);
                    this.directionsRequest.transitOptions.departureTime =
                        new Date(document.getElementById("date" + this.routeNum).value +
                        "T" +    document.getElementById("time" + this.routeNum).value);
                }
                else if(document.getElementById("arrival-time" + this.routeNum).checked){
                    this.directionsRequest.transitOptions.arrivalTime =
                        new Date(document.getElementById("date" + this.routeNum).value +
                        "T" +    document.getElementById("time" + this.routeNum).value);
                }
                this.directionsRequest.transitOptions.routingPreference = 'FEWER_TRANSFERS';
            }
            else if(document.getElementById("changemode-driving" + this.routeNum).checked){
                if(document.getElementById("avoid-toll" + this.routeNum).checked){
                    this.directionsRequest.avoidTolls = true;
                }
                if(document.getElementById("avoid-highway" + this.routeNum).checked){
                    this.directionsRequest.avoidHighways = true;
                }
            }

        //Directions Serviceを使ったルート検索メソッド
        this.directionsService.route(
            this.directionsRequest,
            (response, status) => {
                if (status === "OK") {
                    //複数ルートがある場合、subRouteRendererで各ルートを薄く表示
                    if(me.poly.length > 0){
                        for(var i = 0;i<me.poly.length;i++){
                            me.poly[i].setMap(null);
                        }
                        me.poly = [];
                    }
                    for (var i = 0; i < response.routes.length; i++) {
                        //jsではObjectは参照渡しなので、Object.assignを使って、値渡しにする
                        var sub_res = Object.assign({}, response);
                        sub_res.routes = [response.routes[i]];
                        // Rendererは１つのラインしか描画できないので、各ルートごとにRenderオブジェクトを作成する必要がある
                        var subRouteRenderer = new google.maps.DirectionsRenderer();
                        //ドキュメントURL: https://developers.google.com/maps/documentation/javascript/reference/directions#DirectionsRendererOptions
                        subRouteRenderer.setOptions({
                            //Colorとopacity(不透明度)と太さを設定
                            polylineOptions: {
                                strokeColor: '#808080',
                                strokeOpacity: 0.5,
                                strokeWeight: 7
                            }
                        });
                        subRouteRenderer.setDirections(sub_res);
                        subRouteRenderer.setMap(this.map);
                        me.poly.push(subRouteRenderer);
                    }
                    //responseをRendererに渡して、パネルにルートを表示
                    me.directionsRenderer.setOptions({
                            suppressPolylines: true,
                        });
                    me.directionsRenderer.setDirections(response);
                    console.log(me.directionsRenderer);

                } else {
                    window.alert("検索結果はありません");
                }
            }
        );
    }
}
//ルート番号
var routeID = 0;
//ルート描画時の線の色を指定
const colorMap = {0:"#00bfff",1:"#c8e300",2:"#9543de",3:"#00db30",4:"#4586b5",5:"#00deda",6:"#eb86d5",7:"#83b300",8:"#ffb300",9:"#de0000"}
const routeMap = {
                    "title":"",
                    "routes": {}
                }

//Ajax通信
$(function(){
    $("#save-route").click(function() {
        routeMap["title"] = document.getElementById("route-name").value;
        // 多重送信を防ぐため通信完了までボタンをdisableにする
        var button = $(this);
        button.attr("disabled", true);

        $.ajax({
            url: "/routes_save", // 通信先のURL
            type: "POST",		// 使用するHTTPメソッド
            data: JSON.stringify(routeMap),
            contentType: 'application/json',
            dataType: "json", // responseのデータの種類
            timespan: 1000,	// 通信のタイムアウトの設定(ミリ秒)
            //通信成功
        }).done(function(data,textStatus,jqXHR){
            alert("ルートの保存に成功しました。");
            //通信失敗
        }).fail(function(xhr, status, error){// HTTPエラー時
            alert("Server Error. Pleasy try again later.");
            //通信終了後
        }).always(function(arg1, status, arg2){
            //status が "success" の場合は always(data, status, xhr) となるが
            //、"success" 以外の場合は always(xhr, status, error)となる。
            button.attr("disabled", false);  // ボタンを再び enableにする
        });
    });
});


function initMap() {
    const map = new google.maps.Map(document.getElementById("map"), {
        mapTypeControl: false,
        center: { lat: 35.6816228, lng: 139.765199 },
        zoom: 10,
        //streetViewを無効化
        streetViewControl: false,
        options: {
            gestureHandling: 'greedy' //地図埋め込み時Ctrボタン要求の無効化
        },
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
        if(routeID == 9){
            document.getElementById("add-route").style.display = "none";
        }
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
        /**
         * Assign the project to an employee.
         * @param {String} routeNum - ルートのIndex番号
         * @param {string} colorCode - ルートごとのカラーコード
         * @param {Object} map - google mapオブジェクト
         * @param {Object} directionRequest - directionServiceの引数に指定するオブジェクト
         * @param {String} originPlaceId - Autocomplete Serviceで地名から変換された地点ID
         * @param {String} destinationPlaceId - Autocomplete Serviceで地名から変換された地点ID
         * @param {Array} poly - ルートごとのdirectionRendererオブジェクトの配列
         * @param {Object} travelMode - directionsRequestのオプションフィールド
         * @param {Object} directionsService - google maps API Javascriptのオブジェクト
         * @param {Object} directionsRenderer - レンダリング機能を提供するオブジェクト
         */
        this.routeNum = routeNum;
        this.colorCode = colorMap[parseInt(routeNum)];
        this.map = map;
        this.directionsRequest = {};
        this.originPlaceId = "";
        this.destinationPlaceId = "";
        this.poly = [];
        this.travelMode = google.maps.TravelMode.WALKING;
        this.directionsService = new google.maps.DirectionsService();
        this.directionsRenderer = new google.maps.DirectionsRenderer();
        //初期設定
        this.directionsRenderer.setMap(map);
        this.directionsRenderer.setPanel(document.getElementById("route-detail-panel" + this.routeNum));
        const originInput = document.getElementById("origin-input" + this.routeNum);
        const destinationInput = document.getElementById("destination-input" + this.routeNum);
        const originAutocomplete = new google.maps.places.Autocomplete(originInput);
        //Places detailは高額料金がかかるので、必要なフィールドを指定して、料金を下げる
        originAutocomplete.setFields(["place_id"]);
        const destinationAutocomplete = new google.maps.places.Autocomplete(
            destinationInput
        );
        //Places detailは高額料金がかかるので、必要なフィールドを指定して、料金を下げる
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
        this.setUpDecideRouteListener(this,this.directionsRenderer);

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
        //documentに明記されていない
        google.maps.event.addListener(directionsRenderer, 'routeindex_changed', function () {
            document.getElementById("route-decide" + obj.routeNum).style.display = "block";
            var target = directionsRenderer.getRouteIndex();
            for(var i = 0; i < obj.poly.length; i++){
                if(i == target){
                    obj.poly[i].setOptions({
                        //選択したルートの場合、色をcolorCodeに従って変更
                        polylineOptions: {
                                strokeColor: obj.colorCode,
                                strokeOpacity: 1.0,
                                strokeWeight: 7,
                            //色付きラインを一番上に表示するため、zIndexを他のルートより大きくする。
                            zIndex: parseInt(obj.routeNum) + 1
                        }
                    });
                }
                else{
                    obj.poly[i].setOptions({
                        //選択したルート以外の場合、色を#808080に設定(選択されている場合、色付きだから、
                        //元に戻すためには、全てのルートについて#808080に設定する必要あり。)
                        polylineOptions: {
                            strokeColor: '#808080',
                            strokeOpacity: 0.7,
                            strokeWeight: 7,
                            //色付きラインを一番上に表示するため、zIndexを小さくする
                            zIndex: parseInt(obj.routeNum)
                        }
                    });
                }
                obj.poly[i].setMap(obj.map);
            }
        });
        }
    setUpDecideRouteListener(obj,directionsRenderer) {
            document.getElementById("route-decide" + obj.routeNum).addEventListener("click",function () {
                var target = directionsRenderer.getRouteIndex();
                //ルートを決定したら、toggleを閉じる
                $('#toggle-'+obj.routeNum).next().slideToggle();
                //directionsRendererから目的のルート情報を取得してrouteObjインスタンスを作成
                var ruoteOjb = {
                    geocoded_waypoints: directionsRenderer.directions.geocoded_waypoints,
                    request: directionsRenderer.directions.request,
                    routes: [directionsRenderer.directions.routes[target]],
                    status: directionsRenderer.directions.status,
                    __proto__: directionsRenderer.directions.__proto__
                }
                //選択したルートオブジェクトをrouteMapに追加
                routeMap["routes"][obj.routeNum] = ruoteOjb;
                for(var i = 0; i < obj.poly.length; i++){
                    if(i != target){
                        obj.poly[i].setMap(null);
                    }
                }
        });
    }
    //directions Serviceを使用し、ルート検索
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
            //公共交通機関を選択した場合
            if(document.getElementById("changemode-transit" + this.routeNum).checked){
                this.directionsRequest.transitOptions = {}
                //出発時間を指定した場合
                if(document.getElementById("depart-time" + this.routeNum).checked){
                    console.log(new Date(document.getElementById("date" + this.routeNum).value +
                        "T" +    document.getElementById("time" + this.routeNum).value));
                    console.log(document.getElementById("time" + this.routeNum).value);
                    this.directionsRequest.transitOptions.departureTime =
                        new Date(document.getElementById("date" + this.routeNum).value +
                        "T" +    document.getElementById("time" + this.routeNum).value);
                }
                //到着時間を指定した場合
                else if(document.getElementById("arrival-time" + this.routeNum).checked){
                    this.directionsRequest.transitOptions.arrivalTime =
                        new Date(document.getElementById("date" + this.routeNum).value +
                        "T" +    document.getElementById("time" + this.routeNum).value);
                }
                this.directionsRequest.transitOptions.routingPreference = 'FEWER_TRANSFERS';
            }
            //自動車ルートを指定した場合
            else if(document.getElementById("changemode-driving" + this.routeNum).checked){
                //有料道路不使用の場合
                if(document.getElementById("avoid-toll" + this.routeNum).checked){
                    this.directionsRequest.avoidTolls = true;
                }
                //高速道路不使用の場合
                if(document.getElementById("avoid-highway" + this.routeNum).checked){
                    this.directionsRequest.avoidHighways = true;
                }
            }

        //Directions Serviceを使ったルート検索メソッド
        this.directionsService.route(
            this.directionsRequest,
            (response, status) => {
                if (status === "OK") {
                    //検索結果表示前に、現在の表示を全て削除
                    if(me.poly.length > 0){
                        for(var i = 0;i<me.poly.length;i++){
                            me.poly[i].setMap(null);
                        }
                        me.poly = [];
                    }

                    if(response.request.travelMode == "TRANSIT" && response.routes[0].legs[0].start_address.match(/日本/)){
                        document.getElementById("route-decide"+me.routeNum).style.display = "none";
                        alert("日本国内の公共交通機関情報はご利用いただけません。")
                        return
                    }
                    //複数ルートが帰ってきた場合、それぞれについて、ラインを描画する
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
                    console.log(response.routes[0].summary)
                    console.log(response.routes[0].legs[0].distance.text)
                    console.log(response.routes[0].legs[0].duration.text)

                    //ルートが１つのみの場合、detail-panelが表示されないので、span要素で距離、所要時間を表示する
                    if(response.routes.length == 1){
                        document.getElementById("one-result-panel").style.display = "block"
                        document.getElementById("one-result-text").innerText = "ルート: " +
                                                    response.routes[0].summary +" ," +
                                                    response.routes[0].legs[0].distance.text + " ," +
                                                    response.routes[0].legs[0].duration.text
                    }
                    else {
                        //ルートが２つ以上の場合、必要ないので、表示しない
                        document.getElementById("one-result-panel").style.display = "none"
                    }
                } else {
                    document.getElementById("route-decide"+me.routeNum).style.display = "none";
                    if(this.directionsRequest.travelMode === google.maps.TravelMode.TRANSIT){
                        window.alert("出発地と目的地の距離が遠すぎる場合、結果が表示されない場合があります。\n" +
                            "また、日本国内の公共交通機関情報はご利用いただけません。");
                    }
                    else {
                        window.alert("出発地と目的地の距離が遠すぎる場合、結果が表示されない場合があります。");
                    }
                }
            }
        );
    }
}


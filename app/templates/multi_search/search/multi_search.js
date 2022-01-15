//ルート番号
let routeID = 0;
//ルート描画時の線の色を指定
const colorMap = {
  0: "#00bfff",
  1: "#c8e300",
  2: "#9543de",
  3: "#00db30",
  4: "#4586b5",
  5: "#00deda",
  6: "#eb86d5",
  7: "#83b300",
  8: "#ffb300",
  9: "#de0000",
};
const multiSearchReq = {
  title: "",
  routes: {},
};
// 入力中のルート番号を入れる
let currRouteNum = "0";

//Ajax通信
$(function () {
  $("#save-route").click(function () {
    const keys = Object.keys(multiSearchReq.routes);
    if (keys.length === 0) {
      window.alert("ルートを１つ以上設定してください。");
      return;
    }
    if ($("#route-name").val() === "") {
      window.alert("ルート名は１文字以上入力してください。");
      return;
    }
    if (/[\.\$]/.test(document.getElementById("route-name").value)) {
      window.alert(".または$はルート名に使用できません。");
      return;
    }
    multiSearchReq["title"] = document.getElementById("route-name").value;
    // 多重送信を防ぐため通信完了までボタンをdisableにする
    const button = $(this);
    button.attr("disabled", true);

    $.ajax({
      url: "/routes_save", // 通信先のURL
      type: "POST", // 使用するHTTPメソッド
      data: JSON.stringify(multiSearchReq),
      contentType: "application/json",
      dataType: "json", // responseのデータの種類
      timespan: 1000, // 通信のタイムアウトの設定(ミリ秒)
    })
      //通信成功
      .done(function (data, textStatus, jqXHR) {
        window.location.href = "/multi_search";
        //通信失敗
      })
      .fail(function (xhr, status, error) {
        // HTTPエラー時
        switch (xhr.status) {
          case 401:
            alert(xhr.responseText);
            break;
          case 500:
            alert(xhr.responseText);
        }

        //通信終了後
      })
      .always(function (arg1, status, arg2) {
        //status が "success" の場合は always(data, status, xhr) となるが
        //、"success" 以外の場合は always(xhr, status, error)となる。
        button.attr("disabled", false); // ボタンを再び enableにする
      });
  });
});

//現在時刻を取得し、時間指定要素に挿入
let today = new Date();
let yyyy = today.getFullYear();
let mm = ("0" + (today.getMonth() + 1)).slice(-2); //getMonthは0 ~ 11
let dd = ("0" + today.getDate()).slice(-2);
let ymd = yyyy + "-" + mm + "-" + dd;
let hr = ("0" + today.getHours()).slice(-2);
let minu = ("0" + today.getMinutes()).slice(-2);
let clock = hr + ":" + minu + ":00";

//ブラウザのタイムゾーンのUTCからの時差をminutes単位で取得
let tzoneOffsetminu = today.getTimezoneOffset();

//Google Maps API実行ファイル読み込み
window.onload = function () {
  fetch("/get_api_source")
    .then((resp) => {
      return resp.text();
    })
    .then((MapJS) => {
      window.Function(MapJS)(); //ファイル実行
      initMap(); //APIリクエストのcallbackではなくここで実行
    })
    .catch(() => {
      alert("エラーが発生しました。");
    });
};

//地図上に路線図と道路状況を表示するlayerの表示制御
function setUpLayersListener(id,layerController, m) {
  return function () {
    const layerButton = document.getElementById(id);
    if (!layerButton.checked) {
      layerController.setMap(null);
    } else {
      layerController.setMap(m);
    }
  }
}

function initMap() {
  const map = new google.maps.Map(document.getElementById("map"), {
    mapTypeControl: false,
    center: { lat: 35.6816228, lng: 139.765199 },
    zoom: 10,
    //streetViewを無効化
    streetViewControl: false,
    options: {
      gestureHandling: "greedy", //地図埋め込み時Ctrボタン要求の無効化
    },
  });
  const placesService = new google.maps.places.PlacesService(map);
  const infowindow = new google.maps.InfoWindow();
  const infowindowContent = document.getElementById("infowindow-content");
  infowindow.setContent(infowindowContent);

  //地図上に路線図と道路状況を表示するlayer
  const transitLayer = new google.maps.TransitLayer();
  const trafficLayer = new google.maps.TrafficLayer();
  map.controls[google.maps.ControlPosition.TOP_LEFT].push(document.getElementById("add-layers"));
  document.getElementById("add-transport-layer").addEventListener("change", setUpLayersListener("add-transport-layer", transitLayer, map))
  document.getElementById("add-traffic-layer").addEventListener("change", setUpLayersListener("add-traffic-layer", trafficLayer, map))

  $("#add-route").attr("disabled", true);
  //１番目のルート要素をHTMLに追加
  $("#search-box").append(genSearchBox(routeID, colorMap[routeID]));
  document.getElementById("date" + String(routeID)).value = ymd;
  document.getElementById("date" + String(routeID)).min = ymd;
  document.getElementById("time" + String(routeID)).value = clock;
  //AutocompleteとDiretionsServiceのインスタンス化
  new AutocompleteDirectionsHandler(map, String(routeID),placesService,infowindow);
  $(".toggle-title").on("click", function () {
    $(this).toggleClass("active");
    $(this).next().slideToggle();
  });

  //ルートを決定するまで「次のルートを追加」ボタンが押せないメッセージを表示
  $("#add-route-panel").on("mouseover", function () {
    if (document.getElementById("add-route").disabled === true) {
      if (!$("#add-route").nextAll("small.error-info").length) {
        $("#add-route").after(
          '<br><small class="text-danger error-info">現在のルートを決定するまで次のルートの追加は出来ません。</small>'
        );
      }
    } else {
      if ($("#add-route").nextAll("small.error-info").length) {
        $("#add-route").nextAll("small.error-info").remove();
      }
    }
  });

  //ボタンが押されたら２番目以降のルート要素をHTMLに追加
  $("#add-route").on("click", function () {
    $("#add-route").attr("disabled", true);
    routeID++;
    if (routeID === 9) {
      document.getElementById("add-route").style.display = "none";
    }
    $("#search-box").append(genSearchBox(routeID, colorMap[routeID]));
    document.getElementById("date" + String(routeID)).value = ymd;
    document.getElementById("date" + String(routeID)).min = ymd;
    document.getElementById("time" + String(routeID)).value = clock;
    new AutocompleteDirectionsHandler(map, String(routeID), placesService,infowindow);
    $(".toggle-title").on("click", function () {
      $(".toggle-title").off("click");
      $(".toggle-title").on("click", function () {
        //activeの場合、CSSで下矢印を表示する
        $(this).toggleClass("active");
        $(this).next().slideToggle();
      });
    });
  });
}


class Elements {
    constructor(routeNum) {
      // 経路オプションのラジオボタン
      this.modeWalking = document.getElementById("changemode-walking" + routeNum);
      this.modeTransit = document.getElementById("changemode-transit" + routeNum);
      this.modeDriving = document.getElementById("changemode-driving" + routeNum);
      // 公共交通機関のオプション指定
      this.transitOption = document.getElementById("transit-time" + routeNum);
      this.departNow = document.getElementById("depart-now" + routeNum);
      this.specifyDeparture = document.getElementById("depart-time" + routeNum);
      this.specifyArrival = document.getElementById("arrival-time" + routeNum);

      //自動車のオプション指定
      this.drivingOption = document.getElementById("driving-option" + routeNum);
      this.avoidToll = document.getElementById("avoid-toll" + routeNum);
      this.avoidHighway = document.getElementById("avoid-highway" + routeNum);

      //出発地の入力
      this.originInput = document.getElementById("origin-input" + routeNum);
      //目的地の入力
      this.destinationInput = document.getElementById("destination-input" + routeNum);

      this.routeDecideButton = document.getElementById("route-decide" + routeNum);
    }
}


const MSG_CANNOT_USE_IN_JAPAN = "日本の公共交通機関情報はGoogleによる機能制限により、ご利用いただけません。" +
                                "海外の公共交通機関情報はご利用いただけます。"
const AUTO_COMPLETE_FIELDS = [
    "place_id",
    "geometry",
    "formatted_address",
    "utc_offset_minutes",
    ]

const defaultPolyLineOptions = {
  clickable: true,
  //選択したルート以外の場合、色を#808080に設定(選択されている場合、色付きだから、
  //元に戻すためには、全てのルートについて#808080に設定する必要あり。)
  strokeColor: "#808080",
  strokeOpacity: 0.7,
  strokeWeight: 7,
}

function optionsForAlternatives(zIndex) {
  let polyLineOptions = Object.assign({}, defaultPolyLineOptions);
  polyLineOptions.zIndex =  parseInt(zIndex);
  return polyLineOptions;
}

function optionsForSelected(strokeColor, zIndex) {
  let selectedPolyLine = Object.assign({}, defaultPolyLineOptions);
  //選択したルートの色をcolorCodeに従って変更
  selectedPolyLine.strokeColor = strokeColor;
  selectedPolyLine.strokeOpacity = 1.0;
  //色付きラインを一番上に表示するため、zIndexを大きくする
  selectedPolyLine.zIndex = parseInt(zIndex) + 1;
  return selectedPolyLine;
}



class AutocompleteDirectionsHandler {
  constructor(map, routeNum, placesService, infowindow) {
    /**
     * Assign the project to an employee.
     * @param {String} routeNum - ルートのIndex番号
     * @param {string} colorCode - ルートごとのカラーコード
     * @param {Object} map - google mapオブジェクト
     * @param {Object} directionRequest - directionServiceの引数に指定するオブジェクト
     * @param {String} originPlaceId - Autocomplete Serviceで地名から変換された地点ID
     * @param {Number} originLatitude - Autocomplete Serviceで地名から取得された緯度
     * @param {Number} originLongitude - Autocomplete Serviceで地名から取得された経度
     * @param {String} destinationPlaceId - Autocomplete Serviceで地名から変換された地点ID
     * @param {Object} elements - classに関連するDOM要素
     * @param {Number} timeDiffMin - TimeZone APIから取得された出発地のoffset
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
    this.originLatitude = 0;
    this.originLongitue = 0;
    this.destinationPlaceId = "";
    this.elements = new Elements(routeNum);
    this.timeDiffMin = 0;
    this.poly = [];
    this.focusedElementID = "";
    this.travelMode = google.maps.TravelMode.WALKING;
    this.directionsService = new google.maps.DirectionsService();
    this.directionsRenderer = new google.maps.DirectionsRenderer();
    //初期設定
    this.directionsRenderer.setMap(map);
    this.directionsRenderer.setPanel(
      document.getElementById("route-detail-panel" + this.routeNum)
    );

    const originAutocomplete = new google.maps.places.Autocomplete(this.elements.originInput);
    //Places detailは高額料金がかかるので、必要なフィールドを指定して、料金を下げる
    originAutocomplete.setFields(AUTO_COMPLETE_FIELDS);
    const destinationAutocomplete = new google.maps.places.Autocomplete(
      this.elements.destinationInput
    );
    //Places detailは高額料金がかかるので、必要なフィールドを指定して、料金を下げる
    destinationAutocomplete.setFields(AUTO_COMPLETE_FIELDS);

    //EventListenerの設定
    this.setupClickListener(
      "changemode-walking" + this.routeNum,
      google.maps.TravelMode.WALKING
    );
    this.setupClickListener(
      "changemode-transit" + this.routeNum,
      google.maps.TravelMode.TRANSIT
    );
    this.setupClickListener(
      "changemode-driving" + this.routeNum,
      google.maps.TravelMode.DRIVING
    );
    this.setupPlaceChangedListener(originAutocomplete, "ORIG", this);
    this.setupPlaceChangedListener(destinationAutocomplete, "DEST", this);
    this.setupOptionListener("date" + this.routeNum);
    this.setupOptionListener("time" + this.routeNum);
    this.setupTimeListener("depart-now" + this.routeNum, this.routeNum);
    this.setupOptionListener("avoid-toll" + this.routeNum);
    this.setupOptionListener("avoid-highway" + this.routeNum);
    this.setUpRouteSelectedListener(this, this.directionsRenderer);
    this.setUpDecideRouteListener(this, this.directionsRenderer);
    //マップ上のクリックに対する処理の設定
    this.placesService = placesService;
    this.infowindow = infowindow;
    this.infowindowContent = document.getElementById("infowindow-content");
    this.map.addListener("click", this.handleMapClick.bind(this));
    this.getFocusedElementID("origin-input" + this.routeNum, this);
    this.getFocusedElementID("destination-input" + this.routeNum, this);
    currRouteNum = routeNum;
  }

  //経路オプションのラジオボタンが押されたら発火
  setupClickListener(id, mode) {
    const radioButton = document.getElementById(id);
    radioButton.addEventListener("click", () => {
      if (id === this.elements.modeTransit.id) {
        this.elements.transitOption.style.display =
          "block";
      } else {
        this.elements.transitOption.style.display =
          "none";
      }
      if (id === this.elements.modeDriving.id) {
        this.elements.drivingOption.style.display = "block";
      } else if (id !== this.elements.modeDriving.id) {
        this.elements.drivingOption.style.display = "none";
      }
      this.travelMode = mode;
      this.route();
    });
  }

  //出発地と目的地の入力があった場合、発火
  setupPlaceChangedListener(autocomplete, mode, me) {
    autocomplete.bindTo("bounds", this.map);
    autocomplete.addListener("place_changed", () => {
      // me.infowindow.close(); //クリックした場所の詳細表示を削除
      const place = autocomplete.getPlace();
      if (!place.place_id) {
        window.alert("表示された選択肢の中から選んでください。");
        return;
      } else if (
        me.elements.modeTransit.checked &&
        place.formatted_address.indexOf("日本") !== -1
      ) {
        window.alert(MSG_CANNOT_USE_IN_JAPAN);
        return;
      }
      if (mode === "ORIG") {
        me.originPlaceId = place.place_id;
        //出発地にズームインする
        me.map.setCenter(place.geometry.location);
        me.map.setZoom(15);
      } else {
        me.destinationPlaceId = place.place_id;
      }

      //UTCとの時差をminnutes単位で取得
      me.timeDiffMin = place.utc_offset_minutes;
      //経度と緯度を設定
      me.originLatitude = place.geometry.location.lat();
      me.originLongitue = place.geometry.location.lng();
      me.getPlaceInfo(place.place_id, me);
      me.route();
    });
  }

  //経路オプションが設定された時発火
  setupOptionListener(id) {
    const optionChange = document.getElementById(id);
    optionChange.addEventListener("change", () => {
      this.route();
    });
  }

  //「すぐに出発」ボタンを有効化
  setupTimeListener(id, rNum) {
    const timeNow = document.getElementById(id);
    timeNow.addEventListener("click", () => {
      document.getElementById("date" + rNum).value = yyyy + "-" + mm + "-" + dd;
      document.getElementById("time" + rNum).value = clock;
      this.route();
    });
  }

  //複数ルートがある場合、パネルのルートを押したら発火
  setUpRouteSelectedListener(obj, directionsRenderer) {
    //documentに明記されていない
    directionsRenderer.addListener("routeindex_changed", function () {
      obj.elements.routeDecideButton.style.display = "block";
      const target = directionsRenderer.getRouteIndex();
      for (let i = 0; i < obj.poly.length; i++) {
        if (i === target) {
          obj.poly[i].setOptions(optionsForSelected(obj.colorCode, obj.routeNum));
        } else {
          obj.poly[i].setOptions(optionsForAlternatives(obj.routeNum));
        }
        obj.poly[i].setMap(obj.map);
      }
    });
  }

  setUpDecideRouteListener(obj, directionsRenderer) {
    obj.elements.routeDecideButton.addEventListener("click", function () {
        $("#add-route").attr("disabled", false);
        const target = directionsRenderer.getRouteIndex();
        //ルートを決定したら、toggleを閉じる
        $("#toggle-" + obj.routeNum).next().slideToggle();
        //directionsRendererから目的のルート情報を取得してrouteObjインスタンスを作成
        const ruoteOjb = {
          geocoded_waypoints: directionsRenderer.directions.geocoded_waypoints,
          request: directionsRenderer.directions.request,
          routes: [directionsRenderer.directions.routes[target]],
          status: directionsRenderer.directions.status,
          __proto__: directionsRenderer.directions.__proto__,
        };
        //選択したルートオブジェクトをmultiSearchReqに追加
        multiSearchReq["routes"][obj.routeNum] = ruoteOjb;
        for (let i = 0; i < obj.poly.length; i++) {
          if (i !== target) {
            obj.poly[i].setMap(null);
          }
        }
      });
  }

  //マップ上のクリックを扱うメソッド
  handleMapClick(clickedPlace) {
    if (this.routeNum !== currRouteNum) {
      return;
    }
    const me = this;
    // clickした場所の情報が入ったオブジェクトでplaceIdフィールドがある場合
    if ("placeId" in clickedPlace) {
      if (clickedPlace.placeId) {
        this.getPlaceInfoForClick(clickedPlace.placeId, me);
      }
    }
  }

  //クリックした場所の情報表示とルート検索を行うメソッド
  getPlaceInfo(placeId, me) {
    me.placesService.getDetails(
      {
        placeId: placeId,
        fields: [
          "icon",
          "name",
          "place_id",
          "formatted_address",
          "geometry",
          "utc_offset_minutes",
        ],
      },
      (place, status) => {
        if (
          status === "OK" &&
          place &&
          place.geometry &&
          place.geometry.location
        ) {
            me.infowindow.open(me.map);
            me.infowindow.setPosition(place.geometry.location);
            me.infowindowContent.children["place-icon"].src = place.icon;
            me.infowindowContent.children["place-name"].textContent = place.name;
            me.infowindowContent.children["place-address"].textContent =
                place.formatted_address;
        }
      }
    );
  }

  getPlaceInfoForClick(placeId, me) {
    me.placesService.getDetails(
        {
          placeId: placeId,
          fields: [
            "icon",
            "name",
            "place_id",
            "formatted_address",
            "geometry",
            "utc_offset_minutes",
          ],
        },
        (place, status) => {
          if (
              status === "OK" &&
              place &&
              place.geometry &&
              place.geometry.location
          ) {
            //入力が選択されていなければ、出発地として扱う
            if (!me.focusedElementID) {
              me.focusedElementID = "origin-input" + me.routeNum;
            }
            document.getElementById(me.focusedElementID).value =
                place.formatted_address;
            if (me.focusedElementID === "origin-input" + me.routeNum) {
              me.originPlaceId = place.place_id;
              //UTCとの時差をminutes単位で取得
              me.timeDiffMin = place.utc_offset_minutes;
            } else if (me.focusedElementID === "destination-input" + me.routeNum) {
              me.destinationPlaceId = place.place_id;
            }
            me.infowindow.open(me.map);
            me.infowindow.setPosition(place.geometry.location);
            me.infowindowContent.children["place-icon"].src = place.icon;
            me.infowindowContent.children["place-name"].textContent = place.name;
            me.infowindowContent.children["place-address"].textContent =
                place.formatted_address;

            me.route();
          }
        }
    );
  }

  //origin-input destination-inputどちらが選択されているか取得するメソッド
  getFocusedElementID(id, me) {
    const currentInput = document.getElementById(id);
    currentInput.addEventListener("focus", () => {
      currRouteNum = me.routeNum; //focusされたら、currRouteNumを選択されたルート番号に設定
      me.focusedElementID = id;
    });
  }

  //ルートをクリックした時のイベントリスナーに使用するcallback関数を返すメソッド
  polyLineListenerCallback(idx, me) {
    return function () {
      me.directionsRenderer.setRouteIndex(idx);
      for (let j = 0; j < me.poly.length; j++) {
        if (j === idx) {
          me.poly[j].setOptions(optionsForSelected(me.colorCode, me.routeNum));
        } else {
          me.poly[j].setOptions(optionsForAlternatives(me.routeNum));
        }
        me.poly[j].setMap(me.map);
      }
    };
  }

  //directions Serviceを使用し、ルート検索
  route() {
    if (!this.originPlaceId || !this.destinationPlaceId) {
      return;
    }
    const me = this;
    this.directionsRequest = {
      origin: { placeId: this.originPlaceId },
      destination: { placeId: this.destinationPlaceId },
      travelMode: this.travelMode,
      //↓複数ルートを返す場合、指定
      provideRouteAlternatives: true,
    };
    //公共交通機関を選択した場合
    if (this.elements.modeTransit.checked) {
      if (
          me.elements.originInput.value.indexOf("日本") !== -1
      ) {
        window.alert(MSG_CANNOT_USE_IN_JAPAN);
        return;
      }
      this.directionsRequest.transitOptions = {};
      //時間指定しない場合、現在時刻に設定
      me.directionsRequest.transitOptions.departureTime = new Date(
        ymd + "T" + clock
      );

      //「すぐに出発」以外のボタンが押されている場合
      if (!me.elements.departNow.checked) {
        //ブラウザのタイムゾーンでの指定時間
        let specTime = new Date(
          document.getElementById("date" + this.routeNum).value +
            "T" +
            document.getElementById("time" + this.routeNum).value
        );

        /*(ブラウザのタイムゾーンの時刻) ー (ブラウザのタイムゾーンのoffset) ー (入力地のタイムゾーンのoffset) = (入力地のタイムゾーンの時刻)
                例:ロサンゼルスの鉄道の3月1日,10:00出発を調べたい場合、
                (3月1日,10:00 Asia/Tokyo) -(-9 hors) - (-8 hours) = (3月2日 3:00 Asia/Tokyo) = (3月1日,10:00 America/Los_Angeles)
                (注意)Javascriptの場合、offsetはGMTより進んでいる場合、マイナスになり、TimeZone APIの場合、逆に進んでいる場合プラスになる*/
        specTime.setHours(
          specTime.getHours() -
            Math.round(tzoneOffsetminu / 60) -
            Math.round(me.timeDiffMin / 3600)
        );

        //出発時間を指定した場合
        if (me.elements.specifyDeparture.checked) {
          me.directionsRequest.transitOptions.departureTime = specTime;
        }
        //到着時間を指定した場合
        else if (me.elements.specifyArrival.checked) {
          me.directionsRequest.transitOptions.arrivalTime = specTime;
        }
      }
      //乗り換え回数が最小になるようセット
      me.directionsRequest.transitOptions.routingPreference = "FEWER_TRANSFERS";
    }

    //自動車ルートを指定した場合
    else if (this.elements.modeDriving.checked) {
      //有料道路不使用の場合
      if (me.elements.avoidToll.checked) {
        this.directionsRequest.avoidTolls = true;
      }
      //高速道路不使用の場合
      if (me.elements.avoidHighway.checked) {
        this.directionsRequest.avoidHighways = true;
      }
    }

    //Directions Serviceを使ったルート検索メソッド
    this.directionsService.route(this.directionsRequest, (response, status) => {
      if (status === "OK") {
        //検索結果表示前に、現在の表示を全て削除
        if (me.poly.length > 0) {
          for (let i = 0; i < me.poly.length; i++) {
            me.poly[i].setMap(null);
          }
          me.poly = [];
        }
        if (
          response.request.travelMode === "TRANSIT" &&
          response.routes[0].legs[0].start_address.match(/日本/)
        ) {
          me.elements.routeDecideButton.style.display =
            "none";
          alert(MSG_CANNOT_USE_IN_JAPAN);
          return;
        }
        //複数ルートが帰ってきた場合、それぞれについて、ラインを描画する
        for (let i = 0; i < response.routes.length; i++) {
          let routePolyline = new google.maps.Polyline();
          routePolyline.setPath(response.routes[i].overview_path);

          //
          routePolyline.setOptions(optionsForAlternatives(me.routeNum));
          routePolyline.setMap(me.map);
          me.poly.push(routePolyline);

          let callback = me.polyLineListenerCallback(i, me);
          me.poly[i].addListener("click", callback);
        }
        //インデックス番号0のルートに色をつける
        me.poly[0].setOptions(optionsForSelected(me.routeNum));
        me.directionsRenderer.setRouteIndex(0);

        //responseをRendererに渡して、パネルにルートを表示
        me.directionsRenderer.setOptions({
          suppressPolylines: true,
        });
        me.directionsRenderer.setDirections(response);

        //ルートが１つのみの場合、detail-panelが表示されないので、span要素で距離、所要時間を表示する
        if (response.routes.length === 1) {
          document.getElementById(
            "one-result-panel" + me.routeNum
          ).style.display = "block";
          document.getElementById("one-result-text" + me.routeNum).innerHTML =
            "<span>" +
            response.routes[0].summary +
            "</span>" +
            "<span class='ml-1'>" +
            response.routes[0].legs[0].distance.text +
            "</span>" +
            "<span>" +
            ".</span>" +
            "<span class='ml-1'>約" +
            "<span>" +
            response.routes[0].legs[0].duration.text;
          +"</span>" + "</span>";
        } else {
          //ルートが２つ以上の場合、必要ないので、表示しない
          document.getElementById(
            "one-result-panel" + me.routeNum
          ).style.display = "none";
        }
        //STATUS != OKの場合
      } else {
        me.elements.routeDecideButton.style.display =
          "none";
        window.alert(
          "ルートが見つかりませんでした。出発地と目的地の距離が遠すぎる場合、結果が表示されない場合があります。"
        );
      }
    });
  }
}

function genSearchBox(routeId, color) {
  const route_html_tpl = `<h2 class="toggle-title pl-3 pt-3 mb-0" id="toggle-${routeId}" style="background-color: ${color}">ルート${
    routeId + 1
  }</h2>
        <div class="search-fields">
            <hr color="white" class="mt-0">
            <div id="required-fields">
                <div style="width: 350px"><small>出発地と目的地を入力する場合、地名の一部を入力すると下に選択肢が表示されますので、その中からお選びください。</small></div>
                <input
                        id="origin-input${routeId}"
                        class="controls input-fields"
                        type="text"
                        placeholder="出発地の一部を入力または地名をクリック"
                />
                <br><br>
                <input
                        id="destination-input${routeId}"
                        class="controls input-fields"
                        type="text"
                        placeholder="目的地の一部を入力または地名をクリック"
                />
                <br>
                <div id="route-options${routeId}">
                    <span>移動方法：</span>
                    <br>
                    <input type="radio" name="type" id="changemode-walking${routeId}"  checked="checked"/>
                    <label class="mr-2" for="changemode-walking${routeId}">徒歩</label>

                    <input type="radio" name="type" id="changemode-transit${routeId}"/>
                    <label class="mr-2" for="changemode-transit${routeId}">公共交通機関</label>

                    <input type="radio" name="type" id="changemode-driving${routeId}"/>
                    <label for="changemode-driving${routeId}">自動車</label>
                </div>
            </div>
            <hr class="mt-2">
            <div id="option-fields">
                <div id="transit-time${routeId}" style="display: none">
                    <span>時間指定：</span>
                    <br>
                    <input type="radio" name="timespec${routeId}" id="depart-now${routeId}" checked="checked"/>
                    <label class="mr-2" for="depart-now${routeId}">すぐに出発</label>
                    <input type="radio" name="timespec${routeId}" id="depart-time${routeId}"/>
                    <label class="mr-2" for="depart-time${routeId}">出発時間</label>
                    <input type="radio" name="timespec${routeId}" id="arrival-time${routeId}"/>
                    <label for="arrival-time${routeId}">到着時間</label>
                    <br>
                    <input type="date" id="date${routeId}">
                    <input type="time" id="time${routeId}">
                </div>
                <div id="driving-option${routeId}" style="display: none"><span>運転時オプション：</span>
                    <br>
                    <input type="checkbox" name="driveing-option-select${routeId}" id="avoid-toll${routeId}"/>有料道路を使用しない
                    <br>
                    <input type="checkbox" name="driveing-option-select${routeId}" id="avoid-highway${routeId}"/>高速道路を使用しない
                </div>
            </div>
             <div id="route-detail-panel${routeId}" class="route-detail">
            </div>
            <div style="background-color: white; padding-bottom: 2px">
             <div class="ml-2 mb-2 border" id="one-result-panel${routeId}" style="color: black; display: none">
                <table>
                    <td>ルート:</td>
                    <tbody>
                    <tr>
                        <td>
                            <span id="one-result-text${routeId}" style="color: black"></span>
                        </td>
                    </tr>
                    </tbody>
                </table>
             </div>
             <button class="btn-primary mx-auto" id="route-decide${routeId}" style="display: none">このルートで決定</button>
             </div>
        </div>`;
  return route_html_tpl;
}

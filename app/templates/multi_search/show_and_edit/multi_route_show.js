//ルート番号
var routeID = Object.keys(routeInfo.routes).length - 1;
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
//routeInfoをコピー
const multiSearchUpdateReq = Object.assign({}, routeInfo);
delete multiSearchUpdateReq.route_count;
multiSearchUpdateReq["previous_title"] = routeInfo.title;

// 入力中のルート番号を入れる
var currRouteNum = "0";

//Ajax通信
$(function () {
  $("#save-route").click(function () {
    if (/[\.\$]/.test(document.getElementById("route-name").value)) {
      window.alert(".または$はルート名に使用できません。");
      return;
    }

    if (multiSearchUpdateReq.title === "") {
      window.alert("ルート名は１文字以上入力してください。");
      return;
    }
    multiSearchUpdateReq["title"] = document.getElementById("route-name").value;
    // 多重送信を防ぐため通信完了までボタンをdisableにする
    var button = $(this);
    button.attr("disabled", true);

    $.ajax({
      url: "/update_route", // 通信先のURL
      type: "POST", // 使用するHTTPメソッド
      data: JSON.stringify(multiSearchUpdateReq),
      contentType: "application/json",
      dataType: "json", // responseのデータの種類
      timespan: 1000, // 通信のタイムアウトの設定(ミリ秒)
    })
      //通信成功
      .done(function (data, textStatus, jqXHR) {
        window.location.href = "/mypage/show_routes";
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
var today = new Date();
var yyyy = today.getFullYear();
var mm = ("0" + (today.getMonth() + 1)).slice(-2); //getMonthは0 ~ 11
var dd = ("0" + today.getDate()).slice(-2);
var ymd = yyyy + "-" + mm + "-" + dd;
var hr = ("0" + today.getHours()).slice(-2);
var minu = ("0" + today.getMinutes()).slice(-2);
var clock = hr + ":" + minu + ":00";

//ブラウザのタイムゾーンのUTCからの時差をminutes単位で取得
var tzoneOffsetminu = today.getTimezoneOffset();

//Google Maps API実行ファイル読み込み
window.onload = function () {
  fetch("/get_apikey")
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

  //ルート名をDBから読み込んだルート名に設定
  document.getElementById("route-name").value = routeInfo.title;
  multiSearchUpdateReq.title = routeInfo.title;
  //保存されているルート要素をHTMLに追加
  for (var i = 0; i < Object.keys(routeInfo.routes).length; i++) {
    multiSearchUpdateReq.routes[String(i)] = routeInfo.routes[i];
    $("#search-box").append(genSearchBox(i, colorMap[i]));
    document.getElementById("origin-input" + String(i)).value =
      routeInfo.routes[i].routes[0].legs[0].start_address;
    document.getElementById("destination-input" + String(i)).value =
      routeInfo.routes[i].routes[0].legs[0].end_address;
    document.getElementById("date" + String(i)).value = ymd;
    document.getElementById("date" + String(i)).min = ymd;
    document.getElementById("time" + String(i)).value = clock;
    //AutocompleteとDiretionsServiceのインスタンス化
    var originID = routeInfo.routes[i].request.origin["placeId"];
    var destID = routeInfo.routes[i].request.destination["placeId"];

    new AutocompleteDirectionsHandler(map, String(i), originID, destID);
    //読み込んだルートのtoggleは閉じた状態で表示する
    $(".toggle-title").on("click", function () {
      $(".toggle-title").off("click");
      $(".toggle-title").on("click", function () {
        $(this).toggleClass("active");
        $(this).next().slideToggle();
      });
    });
    //toggleを閉じる
    $("#toggle-" + String(i))
      .next()
      .slideToggle();
  }

  //ルートを決定するまで「新しいルートを追加」ボタンが押せないメッセージを表示
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

  //ボタンが押されたら保存されたルートのインデックス以降のルート要素をHTMLに追加
  $("#add-route").on("click", function () {
    $("#add-route").attr("disabled", true);
    routeID++;
    currRouteNum = String(routeID);
    if (routeID === 9) {
      document.getElementById("add-route").style.display = "none";
    }
    $("#search-box").append(genSearchBox(routeID, colorMap[routeID]));
    document.getElementById("date" + String(routeID)).value = ymd;
    document.getElementById("date" + String(routeID)).min = ymd;
    document.getElementById("time" + String(i)).value = clock;
    new AutocompleteDirectionsHandler(map, String(routeID), "", "");
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

class AutocompleteDirectionsHandler {
  constructor(map, routeNum, originID, destID) {
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
    this.originPlaceId = originID;
    this.originLatitude = 0;
    this.originLongitue = 0;
    this.destinationPlaceId = destID;
    this.timeDiffMin = 0;
    this.poly = [];
    this.inputFieldID = "";
    this.travelMode = google.maps.TravelMode.WALKING;
    this.directionsService = new google.maps.DirectionsService();
    this.directionsRenderer = new google.maps.DirectionsRenderer();

    //初期設定
    this.directionsRenderer.setMap(map);
    this.directionsRenderer.setPanel(
      document.getElementById("route-detail-panel" + this.routeNum)
    );
    //保存されているルートのみ色付きラインを描画
    if (parseInt(routeNum) < Object.keys(routeInfo.routes).length) {
      this.directionsRenderer.setOptions({
        //Colorとopacity(不透明度)と太さを設定
        polylineOptions: {
          strokeColor: this.colorCode,
          strokeOpacity: 1.0,
          strokeWeight: 7,
        },
      });
      this.directionsRenderer.setDirections(routeInfo.routes[routeNum]);
      document.getElementById("one-result-panel" + routeNum).style.display =
        "block";
      document.getElementById("one-result-text" + routeNum).innerText =
        "ルート: " +
        routeInfo.routes[routeNum].routes[0].summary +
        " ," +
        routeInfo.routes[routeNum].routes[0].legs[0].distance.text +
        " ," +
        routeInfo.routes[routeNum].routes[0].legs[0].duration.text;
    }
    //出発地の入力値
    const originInput = document.getElementById("origin-input" + this.routeNum);
    //目的地の入力値
    const destinationInput = document.getElementById(
      "destination-input" + this.routeNum
    );
    const originAutocomplete = new google.maps.places.Autocomplete(originInput);
    //Places detailは高額料金がかかるので、必要なフィールドを指定して、料金を下げる
    originAutocomplete.setFields(["place_id", "geometry", "formatted_address", "utc_offset_minutes"]);
    const destinationAutocomplete = new google.maps.places.Autocomplete(
      destinationInput
    );
    //Places detailは高額料金がかかるので、必要なフィールドを指定して、料金を下げる
    destinationAutocomplete.setFields([
      "place_id",
      "geometry",
      "formatted_address",
      "utc_offset_minutes",
    ]);

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
    this.placesService = new google.maps.places.PlacesService(map);
    this.infowindow = new google.maps.InfoWindow();
    this.infowindowContent = document.getElementById("infowindow-content");
    this.infowindow.setContent(this.infowindowContent);
    this.map.addListener("click", this.handleMapClick.bind(this));
    this.getFocusedElementID("origin-input" + this.routeNum, this);
    this.getFocusedElementID("destination-input" + this.routeNum, this);
  }

  //経路オプションのラジオボタンが押されたら発火
  setupClickListener(id, mode) {
    const radioButton = document.getElementById(id);
    radioButton.addEventListener("click", () => {
      if (id === "changemode-transit" + this.routeNum) {
        document.getElementById("transit-time" + this.routeNum).style.display =
          "block";
      } else if (id !== "changemode-transit" + this.routeNum) {
        document.getElementById("transit-time" + this.routeNum).style.display =
          "none";
      }
      if (id === "changemode-driving" + this.routeNum) {
        document.getElementById(
          "driving-option" + this.routeNum
        ).style.display = "block";
      } else if (id !== "changemode-driving" + this.routeNum) {
        document.getElementById(
          "driving-option" + this.routeNum
        ).style.display = "none";
      }
      this.travelMode = mode;
      this.route();
    });
  }

  //出発地と目的地の入力があった場合、発火
  setupPlaceChangedListener(autocomplete, mode, me) {
    autocomplete.bindTo("bounds", this.map);
    autocomplete.addListener("place_changed", () => {
      me.infowindow.close(); //クリックした場所の詳細表示を削除
      const place = autocomplete.getPlace();
      if (!place.place_id) {
        window.alert("表示された選択肢の中から選んでください。");
        return;
      } else if (
        document.getElementById("changemode-transit" + me.routeNum).checked &&
        place.formatted_address.indexOf("日本") !== -1
      ) {
        window.alert(
          "日本の公共交通機関情報はGoogleによる機能制限により、ご利用いただけません。海外の公共交通機関情報はご利用いただけます。"
        );
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
      me.getPlaceInformation(place.place_id, me);
    });
  }

  //経路オプションが設定された時発火
  setupOptionListener(id) {
    const optionChange = document.getElementById(id);
    optionChange.addEventListener("change", () => {
      this.route();
    });
  }

  //すぐに出発ボタンを有効化
  setupTimeListener(id, rNum) {
    const timeNow = document.getElementById(id);
    timeNow.addEventListener("click", () => {
      document.getElementById("date" + rNum).value = ymd;
      document.getElementById("time" + rNum).value = clock;
      this.route();
    });
  }

  //複数ルートがある場合、パネルのルートを押したら発火
  setUpRouteSelectedListener(obj, directionsRenderer) {
    //documentに明記されていない
    google.maps.event.addListener(
      directionsRenderer,
      "routeindex_changed",
      function () {
        document.getElementById("route-decide" + obj.routeNum).style.display =
          "block";
        var target = directionsRenderer.getRouteIndex();
        for (var i = 0; i < obj.poly.length; i++) {
          if (i == target) {
            obj.poly[i].setOptions({
              //選択したルートの場合、色をcolorCodeに従って変更
              polylineOptions: {
                strokeColor: obj.colorCode,
                strokeOpacity: 1.0,
                strokeWeight: 7,
                //色付きラインを一番上に表示するため、zIndexを他のルートより大きくする。
                zIndex: parseInt(obj.routeNum) + 1,
              },
            });
          } else {
            obj.poly[i].setOptions({
              //選択したルート以外の場合、色を#808080に設定(選択されている場合、色付きだから、
              //元に戻すためには、全てのルートについて#808080に設定する必要あり。)
              polylineOptions: {
                strokeColor: "#808080",
                strokeOpacity: 0.7,
                strokeWeight: 7,
                //色付きラインを一番上に表示するため、zIndexを小さくする
                zIndex: parseInt(obj.routeNum),
              },
            });
          }
          obj.poly[i].setMap(obj.map);
        }
      }
    );
  }

  setUpDecideRouteListener(obj, directionsRenderer) {
    document
      .getElementById("route-decide" + obj.routeNum)
      .addEventListener("click", function () {
        $("#add-route").attr("disabled", false);
        var target = directionsRenderer.getRouteIndex();
        //ルートを決定したら、toggleを閉じる
        $("#toggle-" + obj.routeNum)
          .next()
          .slideToggle();
        //directionsRendererから目的のルート情報を取得してrouteObjインスタンスを作成
        var ruoteOjb = {
          geocoded_waypoints: directionsRenderer.directions.geocoded_waypoints,
          request: directionsRenderer.directions.request,
          routes: [directionsRenderer.directions.routes[target]],
          status: directionsRenderer.directions.status,
          __proto__: directionsRenderer.directions.__proto__,
        };
        //選択したルートオブジェクトをmultiSearchUpdateReqに追加
        multiSearchUpdateReq["routes"][obj.routeNum] = ruoteOjb;
        for (var i = 0; i < obj.poly.length; i++) {
          if (i != target) {
            obj.poly[i].setMap(null);
          }
        }
      });
  }

  //マップ上のクリックを扱うメソッド
  handleMapClick(clickedPlace) {
    if (this.routeNum !== currRouteNum) {
      return
    }
    const me = this;
    if ("placeId" in clickedPlace) {
      // デフォルトのinfo windowを無効化
      clickedPlace.stop();
      if (clickedPlace.placeId) {
        this.getPlaceInformation(clickedPlace.placeId, me);
      }
    }
  }

  //クリックした場所の情報表示とルート検索を行うメソッド
  getPlaceInformation(placeId, me) {
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
            if (!me.inputFieldID) {
              me.inputFieldID = "origin-input" + me.routeNum;
            }
            document.getElementById(me.inputFieldID).value =
                place.formatted_address;
            if (me.inputFieldID === "origin-input" + me.routeNum) {
              me.originPlaceId = place.place_id;
              //UTCとの時差をminutes単位で取得
              me.timeDiffMin = place.utc_offset_minutes;
            } else if (me.inputFieldID === "destination-input" + me.routeNum) {
              me.destinationPlaceId = place.place_id;
            }
            me.infowindow.close();
            me.infowindow.setPosition(place.geometry.location);
            me.infowindowContent.children["place-icon"].src = place.icon;
            me.infowindowContent.children["place-name"].textContent = place.name;
            me.infowindowContent.children["place-address"].textContent =
                place.formatted_address;
            me.infowindow.open(me.map);

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
      me.inputFieldID = id;
    });
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
    if (document.getElementById("changemode-transit" + this.routeNum).checked) {
      if (document.getElementById("origin-input" + me.routeNum).value.indexOf("日本") !== -1) {
        window.alert(
            "日本の公共交通機関情報はGoogleによる機能制限により、ご利用いただけません。海外の公共交通機関情報はご利用いただけます。"
        );
        return;
      }
      this.directionsRequest.transitOptions = {};
      //時間指定しない場合、現在時刻に設定
      me.directionsRequest.transitOptions.departureTime = new Date(
        ymd + "T" + clock
      );

      //「すぐに出発」以外のボタンが押されている場合
      if (!document.getElementById("depart-now" + me.routeNum).checked) {
        //ブラウザのタイムゾーンでの指定時間
        var specTime = new Date(
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
        if (document.getElementById("depart-time" + me.routeNum).checked) {
          me.directionsRequest.transitOptions.departureTime = specTime;
        }
        //到着時間を指定した場合
        else if (
          document.getElementById("arrival-time" + me.routeNum).checked
        ) {
          me.directionsRequest.transitOptions.arrivalTime = specTime;
        }
      }

      //乗り換え回数が最小になるようセット
      this.directionsRequest.transitOptions.routingPreference =
        "FEWER_TRANSFERS";
    }
    //自動車ルートを指定した場合
    else if (
      document.getElementById("changemode-driving" + this.routeNum).checked
    ) {
      //有料道路不使用の場合
      if (document.getElementById("avoid-toll" + this.routeNum).checked) {
        this.directionsRequest.avoidTolls = true;
      }
      //高速道路不使用の場合
      if (document.getElementById("avoid-highway" + this.routeNum).checked) {
        this.directionsRequest.avoidHighways = true;
      }
    }

    //Directions Serviceを使ったルート検索メソッド
    this.directionsService.route(this.directionsRequest, (response, status) => {
      if (status === "OK") {
        //検索結果表示前に、現在の表示を全て削除
        if (me.poly.length > 0) {
          for (var i = 0; i < me.poly.length; i++) {
            me.poly[i].setMap(null);
          }
          me.poly = [];
        }

        if (
          response.request.travelMode == "TRANSIT" &&
          response.routes[0].legs[0].start_address.match(/日本/)
        ) {
          document.getElementById("route-decide" + me.routeNum).style.display =
            "none";
          alert("日本の公共交通機関情報はGoogleによる機能制限により、ご利用いただけません。海外の公共交通機関情報はご利用いただけます。");
          return;
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
              strokeColor: "#808080",
              strokeOpacity: 0.5,
              strokeWeight: 7,
            },
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
      } else {
        document.getElementById("route-decide" + me.routeNum).style.display =
          "none";
        window.alert(
          "ルートが見つかりませんでした。出発地と目的地の距離が遠すぎる場合、結果が表示されない場合があります。"
        );
      }
    });
  }
}

function genSearchBox(routeId, color) {
  var route_html_tpl = `<h2 class="toggle-title pl-3 pt-3 mb-0" id="toggle-${routeId}" style="background-color: ${color}">ルート${
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
                    <input type="radio" name="timespec${routeId}" id="depart-now${routeId}"/>
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

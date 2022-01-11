var destinationFormattedAddress = {};
let originFormattedAddress = routeInfo.origin_address;
var searchResult = {};
const simulSearchUpdateReq = Object.assign({}, routeInfo);
simulSearchUpdateReq.previous_title = routeInfo.title;
let searchFlag = false;
var simulReq = {
    origin: "",
    destinations: {
    },
    mode: "",
    departure_time: "",
    latlng: {lat: "", lng: ""},
    avoid: [],
};
simulReq.origin = routeInfo.origin;
for(var i = 0; i < 10; i ++) {
    if (routeInfo.destinations.hasOwnProperty(i)){
    simulReq.destinations[i] = routeInfo.destinations[i].place_id;
}
}
simulReq.mode = routeInfo.mode;
simulReq.departure_time = routeInfo.departureTime;
simulReq.latlng = Object.assign({}, routeInfo.latlng);
simulReq.avoid = Object.assign([], routeInfo.avoid);

//検索結果保存
$(function () {
    $("#save-route").click(function () {
        var keys = Object.keys(simulSearchUpdateReq);
        if (keys.length === 0) {
            window.alert("ルートを１つ以上設定してください。");
            return;
        }
        if ($("#route-title").val() === "") {
            window.alert("保存名は１文字以上入力してください。");
            return;
        }
        if (/[\.\$]/.test(document.getElementById("route-title").value)) {
            window.alert(".または$はタイトル名に使用できません。");
            return;
        }
        simulSearchUpdateReq.title = document.getElementById("route-title").value;
        if (searchFlag) {
            simulSearchUpdateReq.origin = simulReq.origin;
            simulSearchUpdateReq.origin_address = originFormattedAddress;
            simulSearchUpdateReq.mode = simulReq.mode;
            simulSearchUpdateReq.departure_time = simulReq.departure_time;
            simulSearchUpdateReq.latlng = simulReq.latlng;
            simulSearchUpdateReq.avoid = simulReq.avoid;
            simulSearchUpdateReq.destinations = searchResult;
            for (var i = 1; i < 10; i++) {
                if (destinationFormattedAddress.hasOwnProperty(i)) {
                    simulSearchUpdateReq.destinations[i].address = destinationFormattedAddress[i];
                }
            }
        }
        // 多重送信を防ぐため通信完了までボタンをdisableにする
        var button = $(this);
        button.attr("disabled", true);

        $.ajax({
            url: "/simul_search/update_route", // 通信先のURL
            type: "POST", // 使用するHTTPメソッド
            data: JSON.stringify(simulSearchUpdateReq),
            contentType: "application/json",
            dataType: "json", // responseのデータの種類
            timespan: 1000, // 通信のタイムアウトの設定(ミリ秒)
        })
            //通信成功
            .done(function (data, textStatus, jqXHR) {
                window.location.href = "/simul_search/show_route/?route_title=" + simulSearchUpdateReq.title;
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


//Ajax通信
$(function () {
    $("#simul-search").click(function () {
        // 多重送信を防ぐため通信完了までボタンをdisableにする
        var button = $(this);
        button.attr("disabled", true);

        //公共交通機関選択の場合
        if (
            document.getElementById("transit").checked &&
            document.getElementById("set-future").checked
        ) {
            simulReq["mode"] = "transit";
            simulReq["departure_time"] =
                document.getElementById("date").value +
                "T" +
                document.getElementById("time").value;
        }
        //自動車選択の場合
        else if (document.getElementById("driving").checked) {
            var tmparr = [];
            if (document.getElementById("avoid-tolls").checked) {
                tmparr.push("tolls");
            }
            if (document.getElementById("avoid-highways").checked) {
                tmparr.push("highways");
            }
            simulReq["avoid"] = tmparr;
        }

        $.ajax({
            url: "/do_simul_search", // 通信先のURL
            type: "POST", // 使用するHTTPメソッド
            data: JSON.stringify(simulReq),
            contentType: "application/json",
            dataType: "json", // responseのデータの種類
            timespan: 2000, // 通信のタイムアウトの設定(ミリ秒)
            //通信成功
        })
            .done(function (simulSearchResult, textStatus, jqXHR) {
                searchFlag = true;
                searchResult = simulSearchResult;
                for (var i = 1; i < 10; i++) {
                    if (simulSearchResult[i]) {
                        document.getElementById("distance" + String(i)).innerText =
                            simulSearchResult[i].distance;
                        document.getElementById("duration" + String(i)).innerText =
                            simulSearchResult[i].duration;
                    }
                }
                //通信失敗
            })
            .fail(function (xhr, status, error) {
                // HTTPエラー時
                switch (xhr.status) {
                    case 400:
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


//Google Maps API実行ファイル読み込み
window.onload = function () {
    document.getElementById("origin-input").value = routeInfo.origin_address;
    for (var i = 1; i < 10; i++) {
        if (routeInfo.destinations.hasOwnProperty(i)) {
            destinationFormattedAddress[i] = routeInfo.destinations[i].address;
            var destinationInput = document.getElementById(
                "destination-input" + String(i)
            );
            destinationInput.value = routeInfo.destinations[i].address;
            simulReq["destinations"][i] = routeInfo.destinations[i].place_id;
            document.getElementById("distance" + String(i)).innerText =
                routeInfo.destinations[i].distance;
            document.getElementById("duration" + String(i)).innerText =
                routeInfo.destinations[i].duration;

        }
    }
    document.getElementById("route-title").value = routeInfo.title;
    fetch("/get_api_source")
        .then((resp) => {
            return resp.text();
        })
        .then((MapJS) => {
            window.Function(MapJS)(); //ファイル実行
            initAutocomplete(); //APIリクエストのcallbackではなくここで実行
        })
        .catch(() => {
            alert("エラーが発生しました。");
        });
};

//出発地と目的地の自動入力を設定
function initAutocomplete() {
    // 日付を設定
    document.getElementById("date").value = ymd;
    document.getElementById("date").min = ymd;
    document.getElementById("time").value = clock;

    var originInput = document.getElementById("origin-input");
    const originAutocomplete = new google.maps.places.Autocomplete(originInput);
    //Places detailは高額料金がかかるので、必要なフィールドを指定して、料金を下げる
    originAutocomplete.setFields(["place_id", "geometry", "formatted_address"]);
    originAutocomplete.addListener("place_changed", () => {
        const place = originAutocomplete.getPlace();
        if (!place.place_id) {
            window.alert("表示された選択肢の中から選んでください。");
        } else if (
            document.getElementById("transit").checked &&
            place.formatted_address.indexOf("日本") !== -1
        ) {
            window.alert(
                "日本の公共交通機関情報はGoogle Maps APIの仕様上、ご利用いただけません。"
            );
            return;
        }
        originFormattedAddress = place.formatted_address;
        simulReq["origin"] = place.place_id;
        simulReq["latlng"]["lat"] = String(place.geometry.location.lat());
        simulReq["latlng"]["lng"] = String(place.geometry.location.lng());
    });
    //目的地検索のAutocompleteを有効化
    for (var i = 1; i < 10; i++) {
        new AutocompleteHandler(String(i));
    }
    document
        .getElementById("walking")
        .addEventListener("click", setupClickListener);
    document
        .getElementById("transit")
        .addEventListener("click", setupClickListener);
    document
        .getElementById("driving")
        .addEventListener("click", setupClickListener);
}

//経路オプションのラジオボタンが押されたら発火
function setupClickListener() {
    if (this.id === "transit") {
        if (document.getElementById("origin-input").value.indexOf("日本") !== -1) {
            window.alert(
                "日本の公共交通機関情報はGoogle Maps APIの仕様上、ご利用いただけません。"
            );
            return;
        }
        document.getElementById("departure-time").style.display = "block";
    } else if (this.id !== "transit") {
        document.getElementById("departure-time").style.display = "none";
    }
    if (this.id === "driving") {
        document.getElementById("driving-option").style.display = "block";
    } else if (this.id !== "driving") {
        document.getElementById("driving-option").style.display = "none";
    }
    simulReq.mode = this.value;
}

//目的地検索のAutocompleteを設定
class AutocompleteHandler {
    constructor(routeNum) {
        this.routeNum = routeNum;
        const destinationInput = document.getElementById(
            "destination-input" + routeNum
        );
        const destinationAutocomplete = new google.maps.places.Autocomplete(
            destinationInput
        );
        //Places detailは高額料金がかかるので、必要なフィールドを指定して、料金を下げる
        destinationAutocomplete.setFields([
            "place_id",
            "geometry",
            "formatted_address",
        ]);
        //EventListenerの設定
        this.setupPlaceChangedListener(destinationAutocomplete);
    }

    //目的地の入力があった場合、発火
    setupPlaceChangedListener(autocomplete) {
        autocomplete.addListener("place_changed", () => {
            const place = autocomplete.getPlace();
            if (!place.place_id) {
                window.alert("表示された選択肢の中から選んでください。");
                return;
            } else if (
                document.getElementById("transit").checked &&
                place.formatted_address.indexOf("日本") !== -1
            ) {
                window.alert(
                    "日本の公共交通機関情報はGoogle Maps APIの仕様上、ご利用いただけません。"
                );
                return;
            }
            simulReq["destinations"][this.routeNum] = place.place_id;
            destinationFormattedAddress[this.routeNum] = place.formatted_address;
        });
    }
}

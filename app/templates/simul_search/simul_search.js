//Ajax通信
$(function () {
    $("#simul-search").click(function () {
        // 多重送信を防ぐため通信完了までボタンをdisableにする
        var button = $(this);
        button.attr("disabled", true);

        //公共交通機関選択の場合
        if (document.getElementById("transit").checked && document.getElementById("set-future").checked) {
            simulReq["departure_time"] = document.getElementById("date").value + "T" + document.getElementById("time").value + ":00Z";
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
            console.log(tmparr);
            simulReq["avoid"] = tmparr.join("|");
        }

        $.ajax({
            url: "/do_simul_search", // 通信先のURL
            type: "POST",		// 使用するHTTPメソッド
            data: JSON.stringify(simulReq),
            contentType: 'application/json',
            dataType: "json", // responseのデータの種類
            timespan: 2000,	// 通信のタイムアウトの設定(ミリ秒)
            //通信成功
        }).done(function (data, textStatus, jqXHR) {
            for (var i = 1; i < 10; i++) {
                if (data.resp[i]) {
                    document.getElementById("distance" + String(i)).innerText = data.resp[i][0];
                    document.getElementById("duration" + String(i)).innerText = data.resp[i][1];
                }
            }
            //通信失敗
        }).fail(function (xhr, status, error) {// HTTPエラー時
            alert("Server Error. Pleasy try again later.");
            //通信終了後
        }).always(function (arg1, status, arg2) {
            //status が "success" の場合は always(data, status, xhr) となるが
            //、"success" 以外の場合は always(xhr, status, error)となる。
            button.attr("disabled", false);  // ボタンを再び enableにする
        });
    });
});

var simulReq = {
    "origin": "",
    "destinations": {
        "1": "",
        "2": "",
        "3": "",
        "4": "",
        "5": "",
        "6": "",
        "7": "",
        "8": "",
        "9": "",
    },
    "mode": "",
    "departure_time": "",
    "avoid": "",
};

//出発地と目的地の自動入力を設定
function initAutocomplete() {
    var originInput = document.getElementById("origin-input");
    const originAutocomplete = new google.maps.places.Autocomplete(originInput);
    //Places detailは高額料金がかかるので、必要なフィールドを指定して、料金を下げる
    originAutocomplete.setFields(["place_id", "geometry", "formatted_address"]);
    originAutocomplete.addListener("place_changed", () => {
        const place = originAutocomplete.getPlace();
        if (!place.place_id) {
            window.alert("表示された選択肢の中から選んでください。");
        } else if (document.getElementById("transit").checked && place.formatted_address.indexOf("日本") !== -1) {
            window.alert("日本の公共交通機関情報はGoogle Maps APIの仕様上、ご利用いただけません。");
            return;
        }
        simulReq["origin"] = place.place_id;
    });
    //目的地検索のAutocompleteを有効化
    for (var i = 1; i < 10; i++) {
        new AutocompleteHandler(String(i));
    }
    document.getElementById("walking").addEventListener("click", setupClickListener);
    document.getElementById("transit").addEventListener("click", setupClickListener);
    document.getElementById("driving").addEventListener("click", setupClickListener);
}

//経路オプションのラジオボタンが押されたら発火
function setupClickListener() {
    if (this.id === "transit") {
        document.getElementById("departure-time").style.display = "block"
    } else if (this.id !== "transit") {
        document.getElementById("departure-time").style.display = "none"
    }
    if (this.id === "driving") {
        document.getElementById("driving-option").style.display = "block"
    } else if (this.id !== "driving") {
        document.getElementById("driving-option").style.display = "none"
    }
    simulReq.mode = this.value;
}

//目的地検索のAutocompleteを設定
class AutocompleteHandler {
    constructor(routeNum) {
        this.routeNum = routeNum;
        const destinationInput = document.getElementById("destination-input" + routeNum);
        const destinationAutocomplete = new google.maps.places.Autocomplete(destinationInput);
        //Places detailは高額料金がかかるので、必要なフィールドを指定して、料金を下げる
        destinationAutocomplete.setFields(["place_id", "geometry", "formatted_address"]);
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
            } else if (document.getElementById("transit").checked && place.formatted_address.indexOf("日本") !== -1) {
                window.alert("日本の公共交通機関情報はGoogle Maps APIの仕様上、ご利用いただけません。");
                return;
            }
            simulReq["destinations"][this.routeNum] = place.place_id;
        });
    }
}


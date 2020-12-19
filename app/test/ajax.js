var data = {
    origin: { placeId: this.originPlaceId },
    destination: { placeId: this.destinationPlaceId },
    travelMode: this.travelMode,
    //↓複数ルートを返す場合、指定
    // provideRouteAlternatives: true,
};
$.ajax({
    url:"http://localhost:8080/route_search", // 通信先のURL
    type:"POST",		// 使用するHTTPメソッド
    data:"origin=place_id:" + encodeURIComponent(this.originPlaceId) +
        "&destination=place_id:" + encodeURIComponent(this.destinationPlaceId) +
        "&travelMode=" + encodeURIComponent(this.travelMode), // 送信するデータ
    dataType:"json", // 応答のデータの種類
    // (xml/html/script/json/jsonp/text)
    timespan:1000 		// 通信のタイムアウトの設定(ミリ秒)

    // 2. doneは、通信に成功した時に実行される
    //  引数のdata1は、通信で取得したデータ
    //  引数のtextStatusは、通信結果のステータス
    //  引数のjqXHRは、XMLHttpRequestオブジェクト
}).done(function(data1,textStatus,jqXHR) {
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
    console.log(data1)
    me.directionsRenderer.setDirections(data1);
});
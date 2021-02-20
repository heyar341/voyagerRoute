function genSearchBox(routeId,color) {
    var route_html_tpl = `<h2 class="toggle-title pl-3 pt-3 mb-0" id="toggle-${routeId}" style="background-color: ${color}">ルート${routeId + 1}</h2>
        <div class="search-fields">
            <hr color="white" class="mt-0">
            <div id="required-fields">
                <input
                        id="origin-input${routeId}"
                        class="controls input-fields"
                        type="text"
                        placeholder="出発地を入力"
                />
                <br><br>
                <input
                        id="destination-input${routeId}"
                        class="controls input-fields"
                        type="text"
                        placeholder="目的地を入力"
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
                    <input type="radio" name="deparitime${routeId}" id="depart-time${routeId}"/>
                    <label class="mr-2" for="depart-time${routeId}">出発時間</label>
                    <input type="radio" name="arrivaltime${routeId}" id="arrival-time${routeId}"/>
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
             <div class="ml-2 mb-2" id="one-result-panel"><span id="one-result-text" style="color: black"></span></div>
             <button class="btn-primary mx-auto" id="route-decide${routeId}" style="display: none">このルートで決定</button>
             </div>
        </div>`
    return route_html_tpl
}
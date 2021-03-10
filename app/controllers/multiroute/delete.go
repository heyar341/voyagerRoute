package multiroute

import (
	"app/controllers"
	"app/cookiehandler"
	"app/customerr"
	"app/model"
	"fmt"
	"log"
	"net/http"
)

type deleteRoute struct {
	user       model.User
	routeTitle string
	err        error
}

//deleteRoute delete multi_route_title field from user document
func (d *deleteRoute) deleteRoute() {
	if d.err != nil {
		return
	}
	err := model.UpdateMultiRouteTitles(d.user.ID, d.routeTitle, "$unset", "")
	if err != nil {
		d.err = customerr.BaseErr{
			Op:  "Deleting route title",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while deleting %v from multi_route_titles: %w", d.routeTitle, err),
		}
		return
	}
}

func DeleteRoute(w http.ResponseWriter, req *http.Request) {
	var d deleteRoute
	d.err = controllers.CheckHTTPMethod(req)
	d.user, d.err = controllers.GetUserFromCtx(req)
	d.routeTitle = req.FormValue("title")
	d.deleteRoute()

	//レスポンス作成
	if d.err != nil {
		e := d.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage/show_routes")
		log.Printf("Error while deleting route title: %v", e.Err)
		return
	}

	successMsg := "ルート「" + d.routeTitle + "」を削除しました。"
	cookiehandler.MakeCookieAndRedirect(w, req, "success", successMsg, "/mypage/show_routes")
	log.Printf("User [%v] deleted route [%v]", d.user.UserName, d.routeTitle)
}

////routes collectionから削除 エラー解析に使うかもしれないので、rout自体は削除せずに残しておく

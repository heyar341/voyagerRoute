package multiroute

import (
	"app/customerr"
	"app/model"
	"fmt"
	"log"
	"net/http"
)

type deleteRoute struct {
	user       model.UserData
	routeTitle string
	err        error
}

func (d *deleteRoute) getUserFromCtx(req *http.Request) {
	if d.err != nil {
		return
	}
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		d.err = customerr.BaseErr{
			Op:  "Getting user from context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting user from reuest's context"),
		}
		return
	}
	d.user = user
}

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

func makeCookieAndRedirect(w http.ResponseWriter, req *http.Request, cName, cVal string) {
	c := &http.Cookie{
		Name:   cName,
		Value:  cVal,
		Path:   "/mypage/show_routes",
		MaxAge: 5,
	}
	http.SetCookie(w, c)
	http.Redirect(w, req, "/mypage/show_routes", http.StatusSeeOther)
}

func DeleteRoute(w http.ResponseWriter, req *http.Request) {
	var d = deleteRoute{}
	if req.Method != "POST" {
		d.err = customerr.BaseErr{
			Msg: "HTTPメソッドが不正です。",
			Err: fmt.Errorf("invalid HTTP method access"),
		}
	}
	//Auth middlewareからuserIDを取得
	d.getUserFromCtx(req)
	//requestから要挙するタイトルを取得
	d.routeTitle = req.FormValue("title")
	//「元のルート名をuser documentから削除」
	d.deleteRoute()

	//レスポンス作成
	if d.err != nil {
		e := d.err.(customerr.BaseErr)
		makeCookieAndRedirect(w, req, "msg", e.Msg)
		log.Printf("Error while deleting route title: %v", e.Err)
		return
	}

	makeCookieAndRedirect(w, req, "success", "ルート「"+d.routeTitle+"」を削除しました。")
	log.Printf("User [%v] deleted route [%v]", d.user.UserName, d.routeTitle)
}

////routes collectionから削除 エラー解析に使うかもしれないので、rout自体は削除せずに残しておく

package main

import (
	"app/controllers/auth"
	"app/controllers/middleware"
	"app/controllers/multiroute"
	"app/controllers/mypage"
	"app/controllers/profile"
	"app/controllers/simulsearch"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var home_tpl *template.Template

func init() {
	home_tpl = template.Must(template.Must(template.ParseGlob("templates/home/home.html")).ParseGlob("templates/includes/*.html"))
}

func main() {

	fmt.Println("App started")
	http.Handle("/favicon.ico", http.NotFoundHandler()) //favicon
	http.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("./templates"))))

	//「認証」
	http.HandleFunc("/register_form/", middleware.Auth(auth.RegisterForm))    //新規登録画面
	http.HandleFunc("/check_email", auth.EmailIsAvailable)                    //メールアドレスの可用確認APIのエドポイント
	http.HandleFunc("/register", middleware.RegisterValidator(auth.Register)) //仮登録実行用エンドポイント
	http.HandleFunc("/ask_confirm", middleware.Auth(auth.AskConfirmEmail))    //メールアドレス確認依頼画面
	http.HandleFunc("/login_form/", middleware.Auth(auth.LoginForm))          //ログイン画面
	http.HandleFunc("/login", middleware.LoginValidator(auth.Login))          //ログイン実行用エンドポイント
	http.HandleFunc("/confirm_register", auth.ConfirmRegister)                //本登録実行用エンドポイント
	http.HandleFunc("/logout", auth.Logout)                                   //ログアウト用エンドポイント

	//「まとめ検索」
	http.HandleFunc("/multi_search", middleware.Auth(multiroute.MultiSearchTpl))                               //検索画面
	http.HandleFunc("/get_timezone", middleware.Auth(multiroute.GetTimezone))                                  //タイムゾーン取得用エンドポイント
	http.HandleFunc("/routes_save", middleware.Auth(middleware.SaveRoutesValidator(multiroute.SaveRoutes)))    //保存用エンドポイント
	http.HandleFunc("/show_route/", middleware.Auth(multiroute.ShowAndEditRoutesTpl))                          //確認編集画面
	http.HandleFunc("/update_route", middleware.Auth(middleware.UpdateRouteValidator(multiroute.UpdateRoute))) //編集用エンドポイント
	http.HandleFunc("/delete_route", middleware.Auth(multiroute.DeleteRoute))                                  //削除用エンドポイント

	//「同時検索」
	http.HandleFunc("/simul_search", middleware.Auth(simulsearch.SimulSearchTpl))                   //検索画面
	http.HandleFunc("/do_simul_search", middleware.SimulSearchValidator(simulsearch.DoSimulSearch)) //検索実行用エンドポイント

	//「マイページ」
	http.HandleFunc("/mypage", middleware.Auth(mypage.ShowMypage))                 //マイページ表示
	http.HandleFunc("/mypage/show_routes/", middleware.Auth(mypage.ShowAllRoutes)) //保存したルート一覧
	http.HandleFunc("/mypage/delete_route", middleware.Auth(mypage.ConfirmDelete)) //削除確認

	//「プロフィール」
	http.HandleFunc("/profile", middleware.Auth(profile.ShowProfile)) //プロフィール画面

	//「ホーム」
	http.HandleFunc("/", middleware.Auth(home))

	http.ListenAndServe(":80", nil)
}

func home(w http.ResponseWriter, req *http.Request) {
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		log.Printf("Error whle gettibg data from context")
		data = map[string]interface{}{"isLoggedIn": false}
	}
	data["msg"] = req.URL.Query().Get("msg")
	data["success"] = req.URL.Query().Get("success")
	home_tpl.ExecuteTemplate(w, "home.html", data)
}

package main

import (
	"app/contexthandler"
	"app/controllers/api"
	"app/controllers/auth"
	"app/controllers/multiroute"
	"app/controllers/mypage"
	"app/controllers/profile"
	"app/controllers/simulsearch"
	"app/mailhandler"
	"app/middleware"
	"app/reqvalidator"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
)

var homeTpl *template.Template

func init() {
	homeTpl = template.Must(template.Must(template.ParseGlob("templates/home/home.html")).ParseGlob("templates/includes/*.html"))
}

func main() {

	fmt.Println("App started")
	http.Handle("/favicon.ico", http.NotFoundHandler()) //favicon
	http.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("./templates"))))

	//「認証」
	http.HandleFunc("/check_email", mailhandler.EmailIsAvailable)               //メールアドレスの可用確認APIのエドポイント
	http.HandleFunc("/register", reqvalidator.RegisterValidator(auth.Register)) //仮登録実行用画面とエンドポイント
	http.HandleFunc("/login", reqvalidator.LoginValidator(auth.Login))          //ログイン実行用画面とエンドポイント
	http.HandleFunc("/confirm_register/", auth.ConfirmRegister)                 //本登録実行用エンドポイント
	http.HandleFunc("/logout", auth.Logout)                                     //ログアウト用エンドポイント

	//「まとめ検索」
	http.HandleFunc("/multi_search", middleware.Auth(multiroute.MultiSearchTpl))                                 //検索画面
	http.HandleFunc("/get_api_source", api.GetApiSource)                                                         //Google Maps API Javascriptの実行に必要なJavascriptファイルを取得するためのエンドポイント
	http.HandleFunc("/routes_save", middleware.Auth(reqvalidator.SaveRoutesValidator(multiroute.SaveNewRoute)))  //保存用エンドポイント
	http.HandleFunc("/show_route/", middleware.Auth(multiroute.ShowAndEditRoutesTpl))                            //確認編集画面
	http.HandleFunc("/update_route", middleware.Auth(reqvalidator.UpdateRouteValidator(multiroute.UpdateRoute))) //編集用エンドポイント
	http.HandleFunc("/delete_route", middleware.Auth(multiroute.DeleteRoute))                                    //削除用エンドポイント

	//「同時検索」
	http.HandleFunc("/simul_search", middleware.Auth(simulsearch.SimulSearchTpl))                                                   //検索画面
	http.HandleFunc("/do_simul_search", reqvalidator.SimulSearchValidator(simulsearch.DoSimulSearch))                               //検索実行用エンドポイント
	http.HandleFunc("/simul_search/routes_save", middleware.Auth(reqvalidator.SaveSimulRouteValidator(simulsearch.SaveNewRoute)))   //保存用エンドポイント
	http.HandleFunc("/simul_search/show_route/", middleware.Auth(simulsearch.ShowAndEditSimulRoutesTpl))                            //確認編集画面
	http.HandleFunc("/simul_search/update_route", middleware.Auth(reqvalidator.UpdateSimulRouteValidator(simulsearch.UpdateRoute))) //編集用エンドポイント

	//「マイページ」
	http.HandleFunc("/mypage", middleware.Auth(mypage.ShowMypage))                                  //マイページ表示
	http.HandleFunc("/mypage/show_routes", middleware.Auth(mypage.ShowAllRoutes))                   //保存したルート一覧
	http.HandleFunc("/mypage/simul_search/show_routes", middleware.Auth(mypage.ShowAllSimulRoutes)) //保存した同時検索一覧
	http.HandleFunc("/mypage/delete_route", middleware.Auth(mypage.ConfirmDelete))                  //削除確認
	http.HandleFunc("/question_form", middleware.Auth(mypage.ShowQuestionForm))                     //お問い合わせ入力ページ
	http.HandleFunc("/send_question", middleware.Auth(mailhandler.SendQuestion))                    //お問い合わせ送信用エンドポイント

	//「プロフィール」
	http.HandleFunc("/profile/username_edit_form", middleware.Auth(profile.EditUserNameForm)) //プロフィール画面
	http.HandleFunc("/profile/username_edit", middleware.Auth(profile.UpdateUserName))        //ユーザー名編集画面
	http.HandleFunc("/profile/email_edit_form", middleware.Auth(profile.EditEmailForm))       //ユーザー名編集用エンドポインt
	http.HandleFunc("/profile/email_edit", middleware.Auth(profile.UpdateEmail))              //メールアドレス編集画面
	http.HandleFunc("/confirm_email/", middleware.Auth(profile.ConfirmUpdateEmail))           //メールアドレス編集用画面
	http.HandleFunc("/profile/password_edit_form", middleware.Auth(profile.EditPasswordForm)) //パスワード編集画面
	http.HandleFunc("/profile/password_edit", middleware.Auth(profile.UpdatePassword))        //パスワード編集用画面
	http.HandleFunc("/profile/", middleware.Auth(profile.ShowProfile))                        //プロフィール画面

	//「ホーム」
	http.HandleFunc("/", middleware.Auth(home))

	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, req *http.Request) {
	data := contexthandler.GetLoginStateFromCtx(req)
	//successメッセージがある場合
	c, err := req.Cookie("success")
	if err == nil {
		processCookie(w, c, data)
		return
	}
	//エラーメッセージがある場合
	c, err = req.Cookie("msg")
	if err == nil {
		processCookie(w, c, data)
		return
	}
	homeTpl.ExecuteTemplate(w, "home.html", data)
}

func processCookie(w http.ResponseWriter, c *http.Cookie, data map[string]interface{}) {
	b64Str, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		homeTpl.ExecuteTemplate(w, "home.html", data)
		return
	}
	data[c.Name] = string(b64Str)
	homeTpl.ExecuteTemplate(w, "home.html", data)
}

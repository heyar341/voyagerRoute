package main

import (
	"app/controllers/auth"
	"app/controllers/home"
	"app/controllers/multiroute"
	"app/controllers/mypage"
	"app/controllers/profile"
	"app/controllers/simulsearch"
	"app/internal/api"
	"app/internal/mailhandler"
	"app/middleware"
	"app/reqvalidator"
	"fmt"
	"net/http"
)

func main() {

	fmt.Println("App started")
	http.Handle("/favicon.ico", http.NotFoundHandler()) //favicon
	http.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("./templates"))))

	//「認証」
	http.HandleFunc("/check_email", mailhandler.EmailIsAvailable)               //メールアドレスの可用確認APIのエドポイント
	http.HandleFunc("/register", reqvalidator.RegisterValidator(auth.Register)) //仮登録実行用画面とエンドポイント
	http.HandleFunc("/login", reqvalidator.LoginValidator(auth.Login))          //ログイン実行用画面とエンドポイント
	http.HandleFunc("/confirm_register/", auth.ConfirmRegister)                 //本登録実行用エンドポイント
	http.HandleFunc("/logout", middleware.CheckHTTPMethod(auth.Logout))         //ログアウト用エンドポイント

	//「まとめ検索」
	http.HandleFunc("/multi_search", middleware.Auth(multiroute.Index))                                                                                                  //検索画面
	http.HandleFunc("/get_api_source", api.GetApiSource)                                                                                                                 //Google Maps API Javascriptの実行に必要なJavascriptファイルを取得するためのエンドポイント
	http.HandleFunc("/routes_save", middleware.Auth(reqvalidator.SaveRoutesValidator(multiroute.Save)))                                                                  //保存用エンドポイント
	http.HandleFunc("/multi_search/show_route/", middleware.Auth(multiroute.Show))                                                                                       //確認編集画面
	http.HandleFunc("/update_route", middleware.CheckHTTPMethod(middleware.CheckHTTPContentType(middleware.Auth(reqvalidator.UpdateRouteValidator(multiroute.Update))))) //編集用エンドポイント
	http.HandleFunc("/multi_search/delete_route", middleware.CheckHTTPMethod(middleware.Auth(multiroute.Delete)))                                                        //削除用エンドポイント

	//「同時検索」
	http.HandleFunc("/simul_search", middleware.Auth(simulsearch.Index))                                                                                                                    //検索画面
	http.HandleFunc("/do_simul_search", reqvalidator.SimulSearchValidator(simulsearch.Search))                                                                                              //検索実行用エンドポイント
	http.HandleFunc("/simul_search/routes_save", middleware.Auth(reqvalidator.SaveSimulRouteValidator(simulsearch.Save)))                                                                   //保存用エンドポイント
	http.HandleFunc("/simul_search/show_route/", middleware.Auth(simulsearch.Show))                                                                                                         //確認編集画面
	http.HandleFunc("/simul_search/update_route", middleware.CheckHTTPMethod(middleware.CheckHTTPContentType(middleware.Auth(reqvalidator.UpdateSimulRouteValidator(simulsearch.Update))))) //編集用エンドポイント
	http.HandleFunc("/simul_search/delete_route", middleware.CheckHTTPMethod(middleware.Auth(simulsearch.Delete)))                                                                          //編集用エンドポイント

	//「マイページ」
	http.HandleFunc("/mypage", middleware.Auth(mypage.Mypage))                      //マイページ表示
	http.HandleFunc("/mypage/show_routes/", middleware.Auth(mypage.ShowAllRoutes))  //保存したルート一覧
	http.HandleFunc("/mypage/delete_route/", middleware.Auth(mypage.ConfirmDelete)) //削除確認
	http.HandleFunc("/question_form", middleware.Auth(mypage.QuestionForm))         //お問い合わせ入力ページ
	http.HandleFunc("/send_question", middleware.Auth(mailhandler.SendQuestion))    //お問い合わせ送信用エンドポイント

	//「プロフィール」
	http.HandleFunc("/profile/username_edit", middleware.Auth(profile.EditUserName)) //ユーザー名編集画面
	http.HandleFunc("/profile/email_edit", middleware.Auth(profile.EditEmail))       //メールアドレス編集画面
	http.HandleFunc("/confirm_email/", middleware.Auth(profile.ConfirmUpdateEmail))  //メールアドレス編集用画面
	http.HandleFunc("/profile/password_edit", middleware.Auth(profile.EditPassword)) //パスワード編集用画面
	http.HandleFunc("/profile/", middleware.Auth(profile.ShowProfile))               //プロフィール画面

	//「ホーム」
	http.HandleFunc("/", middleware.Auth(home.Show))

	http.ListenAndServe(":8080", nil)
}

package profile

import (
	"app/contexthandler"
	"app/controllers"
	"app/cookiehandler"
	"app/customerr"
	"app/model"
	"fmt"
	"log"
	"net/http"
)

type updateUserName struct {
	user        model.User
	newUserName string
	err         error
}

const REDIRECT_URI_TO_UPDATE_USERNAME_FORM = "/profile/username_edit_form"

//getUserNameFromForm gets username from request form.
func (uU *updateUserName) getUserNameFromForm(req *http.Request) {
	if uU.err != nil {
		return
	}
	newUserName := req.FormValue("username")

	if newUserName == "" {
		uU.err = customerr.BaseErr{
			Op:  "get username from request form",
			Msg: "ユーザー名は１文字以上入力してください。",
			Err: fmt.Errorf("request's username was empty"),
		}
		return
	}
	uU.newUserName = newUserName
}

//updateUserName updates user document's username field.
func (uU *updateUserName) updateUserName() {
	if uU.err != nil {
		return
	}
	err := model.UpdateUser(uU.user.ID, "username", uU.newUserName)
	if err != nil {
		uU.err = customerr.BaseErr{
			Op:  "Saving editing email to DB",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while inserting editing email to editing_email collection %w", err),
		}
		return
	}
}

func UpdateUserName(w http.ResponseWriter, req *http.Request) {
	var uU updateUserName
	controllers.CheckHTTPMethod(req, &uU.err)
	contexthandler.GetUserFromCtx(req, &uU.user, &uU.err)
	uU.getUserNameFromForm(req)
	uU.updateUserName()

	if uU.err != nil {
		e := uU.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, REDIRECT_URI_TO_UPDATE_USERNAME_FORM)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	cookiehandler.MakeCookieAndRedirect(w, req, "success", "ユーザー名を変更しました。", "/profile")
}

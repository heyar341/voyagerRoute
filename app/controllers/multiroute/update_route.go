package multiroute

import (
	"app/customerr"
	"app/model"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strings"
	"time"
)

//ルート編集保存requestのフィールドを保存するstruct
type RouteUpdateRequest struct {
	ID            primitive.ObjectID     `json:"_id" bson:"_id"`
	Title         string                 `json:"title" bson:"title"`
	PreviousTitle string                 `json:"previous_title" bson:"previous_title"`
	Routes        map[string]model.Route `json:"routes" bson:"routes"`
}

//getUpdateRoutesInfo gets validated RouteUpdateRequest from context
func getUpdateRoutesInfo(req *http.Request) (*routesData, string) {
	reqFields, ok := req.Context().Value("reqFields").(RouteUpdateRequest)
	if !ok {
		return &routesData{
			err: customerr.BaseErr{
				Op:  "Get req routes info from context",
				Msg: "エラーが発生しました。",
				Err: fmt.Errorf("error while getting request fields from reuest's context"),
			},
		}, ""
	}
	return &routesData{
		routesInfo: model.MultiRoute{
			ID:     reqFields.ID,
			Title:  reqFields.Title,
			Routes: reqFields.Routes,
		},
		err: nil,
	}, reqFields.PreviousTitle
}

//updateRoute updates route document in routes collection
func (r *routesData) updateRoute() {
	if r.err != nil {
		return
	}
	err := r.routesInfo.UpdateRoute()
	if err != nil {
		e := customerr.BaseErr{
			Op:  "Save new multi route",
			Err: fmt.Errorf("error while updating multi route: %w", err),
		}

		if strings.Contains(err.Error(), "(BSONObjectTooLarge)") {
			e.Msg = "ルートのデータサイズが大きすぎるため、保存できません。\nルートの分割をお願いします。"
			r.err = e
		} else {
			e.Msg = "エラーが発生しました。"
			r.err = e
		}
	}
}

//updateRouteTitles updates timestamp of multi_route_titles field in users collection, and delete previous title if title was changed.
func (r *routesData) updateRouteTitles(pTitle string) {
	if r.err != nil {
		return
	}
	now := time.Now().UTC() //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
	if r.routesInfo.Title != pTitle {
		//「元のルート名をuser documentから削除」
		//documentではなく、document内のフィールドを削除する場合、Deleteではなく、Update operatorの$unsetを使って削除する
		//公式ドキュメントURL: https://docs.mongodb.com/manual/reference/operator/update/unset/
		err := model.UpdateMultiRouteTitles(r.userID, pTitle, "$unset", "")
		if err != nil {
			r.err = customerr.BaseErr{
				Op:  "Remove previous multi route title",
				Msg: "エラーが発生しました。",
				Err: fmt.Errorf("error while removing previous multi route title: %w", err),
			}
			return
		}
	}
	//「タイムスタンプを更新または追加」
	err := model.UpdateMultiRouteTitles(r.userID, r.routesInfo.Title, "$set", now)
	if err != nil {
		r.err = customerr.BaseErr{
			Op:  "Set new multi route title and timestamp",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while setting new multi route title and timestamp: %w", err),
		}
		return
	}
}

//ルートを更新保存するための関数
func UpdateRoute(w http.ResponseWriter, req *http.Request) {
	//バリデーション完了後のルート情報と変更前のタイトルを取得
	r, pTitle := getUpdateRoutesInfo(req)
	//Auth middlewareからuserIDを取得
	r.getUserID(req)
	//ルートの更新
	r.updateRoute()
	//ルート名とタイムスタンプの更新
	r.updateRouteTitles(pTitle)

	if r.err != nil {
		e := r.err.(customerr.BaseErr)
		http.Error(w, e.Msg, http.StatusInternalServerError)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	respMsg := map[string]string{"msg": "ok"}
	json.NewEncoder(w).Encode(&respMsg)
}

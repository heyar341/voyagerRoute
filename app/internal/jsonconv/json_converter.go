package jsonconv

import (
	"app/internal/customerr"
	"app/internal/errormsg"
	"encoding/json"
	"fmt"
)

//StructToJSON makes JSON object from a struct
func StructToJSON(s interface{}, err *error) string {
	//レスポンス作成
	jsonEnc, e := json.Marshal(s)
	if e != nil {
		*err = customerr.BaseErr{
			Op:  "json marshaling multiRoute struct",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while json marshaling: %w", err),
		}
		return ""
	}
	//JSONのバイナリ形式のままだとtemplateで読み込めないので、stringに変換
	return string(jsonEnc)
}

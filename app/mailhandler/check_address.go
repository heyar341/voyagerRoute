package mailhandler

import (
	"net"
	"regexp"
	"strings"
)

//メールアドレスの正規表現
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmailValid(email string) bool {
	//文字数チェック
	if len(email) < 3 && len(email) > 254 {
		return false
		//正規表現でチェック
	} else if !emailRegex.MatchString(email) {
		return false
	}
	domain := strings.Split(email, "@")[1]
	mx, err := net.LookupMX(domain)
	if err != nil {
		return false
	} else if len(mx) == 0 {
		return false
	}
	return true
}

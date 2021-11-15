package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (j JsonTime) MarshaJSON() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-04-02"))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32    `json:"id"`
	Mobile   string   `json:"mobile"`
	NickName string   `json:"nick_name"`
	Birthday JsonTime `json:"birthday"`
	Gender   string   `json:"gender"`
}

package types

import "fmt"

type Message struct {
	ErrCode  int         `json:"ErrCode" xml:"ErrCode"`
	Resource interface{} `json:"Resource" xml:"Resource"`
	Hint     string      `json:"Hint" xml:"Hint"`
}

func (m *Message)ToString()string {
	return fmt.Sprintf("{ErrCode:%d, Resource:%v, Hint:%s} \n", m.ErrCode, m.Resource, m.Hint)
}
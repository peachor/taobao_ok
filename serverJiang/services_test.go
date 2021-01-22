package serverJiang

import (
	"testing"
)

func TestServerJiang(t *testing.T) {
	var s ServerJiang
	ss := make(map[string]string)
	ss["text"] = "商品到货通知"
	ss["desp"] = "优衣库到货了！！！！"
	s.Data = ss
	s.SCKey = "xxx"
	s.Do()
}

package serverJiang

import (
	gg "github.com/lycblank/goprogressbar"
	"testing"
	"time"
)

func TestServerJiang(t *testing.T) {

	//var s ServerJiang
	//ss := make(map[string]string)
	//ss["text"] = "商品到货通知"
	//ss["desp"] = "优衣库到货了！！！！"
	//s.Data=ss
	//s.Do()

	//
	go func(num int64) {
		bar := gg.NewProgressBar(num * 60)
		for i := 1; i <= int(num*60); i++ {
			bar.Play(int64(i))
			time.Sleep(1 * time.Second)
		}
		bar.Finish()

	}(int64(1))

	time.Sleep(2*time.Minute)

}

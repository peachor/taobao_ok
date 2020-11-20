package serverJiang

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ServerJiang struct {
	Data map[string]string
}

func (s ServerJiang) Do() {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	DataUrlVal := url.Values{}
	for key, val := range s.Data {
		DataUrlVal.Add(key, val)
	}
	resp, err := client.Post("https://sc.ftqq.com/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.send", "application/x-www-form-urlencoded", strings.NewReader(DataUrlVal.Encode()))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(body)
}

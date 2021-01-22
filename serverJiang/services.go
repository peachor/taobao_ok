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
	SCKey string
	Data  map[string]string
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

	resp, err := client.Post(fmt.Sprintf("https://sc.ftqq.com/%s.send", s.SCKey), "application/x-www-form-urlencoded", strings.NewReader(DataUrlVal.Encode()))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

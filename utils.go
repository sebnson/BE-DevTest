package main

import (
	"io/ioutil"
	"net/http"
)

func makeHTTPGetRequest(url string) ([]byte, error) {
	//HTTP GET 요청 보내기
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	//응답 데이터 읽기
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

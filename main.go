package main

import (
	"encoding/json"
	"strconv"
	//"fmt"
	"log"
	"time"
)

func main() {
	// 스케줄러 시작
	go scheduler()

	// API 서버 시작
	startAPIServer()
}

func scheduler() {
	for {
		// Chainlink 컨트랙트, Bitfinex API에서 토큰 가격 정보를 조회, 데이터베이스에 저장
		err := fetchDataAndSave()
		if err != nil {
			log.Println("Error fetching data and saving to database:", err)
		}

		// 30초 간격으로 스케줄링
		time.Sleep(30 * time.Second)
	}
}

func fetchDataAndSave() error {
	// Chainlink 컨트랙트에서 토큰 가격 정보 조회
	tokenPrice, err := getLatestTokenData()
	if err != nil {
		return err
	}

	// Bitfinex API에서 토큰 가격 정보 조회
	tokenSymbol := "USTUSD" // USDT 토큰 가격 정보를 조회
	// tokenSymbol에 맞는 Bitfinex API URL 생성, HTTP 요청 전송
	bitfinexURL := "https://api.bitfinex.com/v1/pubticker/" + tokenSymbol
	response, err := makeHTTPGetRequest(bitfinexURL)
	if err != nil {
		return err
	}

	// 응답 데이터 파싱하여 토큰 가격 정보 가져오기
	var priceInfo struct {
		LastPrice string `json:"last_price"`
	}
	err = json.Unmarshal(response, &priceInfo)
	if err != nil {
		return err
	}

	// 문자열 형태의 토큰 가격을 실수로 변환
	tokenPrice, err := strconv.ParseFloat(priceInfo.LastPrice, 64)
	if err != nil {
		return err
	}

	// 데이터베이스에 토큰 가격 정보 저장
	err = saveTokenPrice(tokenSymbol, tokenPrice)
	if err != nil {
		return err
	}

	return nil
}

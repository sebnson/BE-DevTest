package main

import (
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

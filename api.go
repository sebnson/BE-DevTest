package main

import (
	"encoding/json"
	"log"
	"net/http"
)

//API 서버 시작
func startAPIServer() {
	http.HandleFunc("/token/price", handleTokenPriceRequest)
	http.HandleFunc("/averageTokenPrice", getAverageTokenPrice)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//토큰 가격 정보 반환
func handleTokenPriceRequest(w http.ResponseWriter, r *http.Request) {
	// 요청 파라미터에서 토큰 심볼을 읽기 (예: /token/price?symbol=USDT)
	tokenSymbol := r.URL.Query().Get("symbol")

	// 데이터베이스에서 토큰 가격 정보 조회
	priceInfo, err := getTokenPrice(tokenSymbol)
	if err != nil {
		http.Error(w, "Error getting token price", http.StatusInternalServerError)
		return
	}

	// JSON 응답 반환 //응답 데이터를 JSON 형식으로 변환하여 클라이언트에 반환
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(priceInfo)
	if err != nil {
		http.Error(w, "Error encoding JSON response.", http.StatusInternalServerError)
		return
	}
}

package main

import (
	"encoding/json"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Chainlink 컨트랙트 주소
var chainlinkAddress = map[string]string{
	"DAI": "0x0630521aC362bc7A19a4eE44b57cE72Ea34AD01c", //DAI 토큰피드
	"ETH": "0x143db3CEEfbdfe5631aDD3E50f7614B6ba708BA7", //ETH 토큰피드
}

var bitfinexToken = map[string]string{
	"USDT": "USTUSD",
	"ETH":  "ETHUSD",
}

var chainlinkABI abi.ABI
var err error

func init() {
	// ChainLink.json에서 ABI 읽어오기
	abiContent, err := os.ReadFile("ChainLink.json")
	if err != nil {
		log.Fatalf("Failed to read ABI file: %v", err)
	}

	// ABI 파싱
	chainlinkABI, err = abi.JSON(strings.NewReader(string(abiContent)))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}
}

func fetchDataAndSave() error {
	// Chainlink 컨트랙트에서 토큰 가격 정보 조회
	for token, chainlinkToken := range chainlinkAddress {
		// Chainlink 토큰 가격 정보 조회 및 저장
		chainlinkPrice, err := getChainlinkTokenData(chainlinkToken)
		if err != nil {
			log.Printf("Error fetching price from Chainlink: %v", err)
			continue
		}
		err = saveTokenPrice(token, chainlinkPrice, "chainlink")
		if err != nil {
			log.Printf("Error saving price from Chainlink: %v", err)
			continue
		}

	}
	//  bitfinex에서 토큰 가격 정보 조회
	for token, bitfinexToken := range bitfinexToken {
		bitfinexPrice, err := getBitfinexTokenData(bitfinexToken)
		if err != nil {
			log.Printf("Error fetching price from bitfinex: %v", err)
			continue
		}
		err = saveTokenPrice(token, bitfinexPrice, "bitfinex")
		if err != nil {
			log.Printf("Error saving price from bitfinex: %v", err)
			continue
		}
	}
	return nil
}

func getChainlinkTokenData(address string) (float64, error) {
	client, err := ethclient.Dial("https://data-seed-prebsc-1-s1.binance.org:8545/") // Ethereum 노드에 연결
	if err != nil {
		return 0, err
	}

	// latestRoundData() 함수 호출용 데이터 생성
	data, err := chainlinkABI.Pack("latestRoundData")
	if err != nil {
		return 0, err
	}

	addr := common.HexToAddress(address)
	callMsg := ethereum.CallMsg{
		To:   &addr,
		Data: data,
	}

	result, err := client.CallContract(nil, callMsg, nil)
	if err != nil {
		return 0, err
	}

	// 데이터 파싱 및 가격 리턴
	currentTokenPrice, _, err := parseTokenPriceData(result)
	if err != nil {
		return 0, err
	}

	return currentTokenPrice, nil
}

func getBitfinexTokenData(bitfinexToken string) (float64, error) {
	// tokenSymbol에 맞는 Bitfinex API URL 생성, HTTP 요청 전송
	bitfinexURL := "https://api.bitfinex.com/v1/pubticker/" + bitfinexToken
	response, err := makeHTTPGetRequest(bitfinexURL)
	if err != nil {
		return 0, err
	}

	// 응답 데이터 파싱, 토큰 가격 정보 가져오기
	var priceInfo struct {
		LastPrice string `json:"last_price"`
	}
	err = json.Unmarshal(response, &priceInfo)
	if err != nil {
		return 0, err
	}

	// 문자열 형태의 토큰 가격을 실수로 변환
	bitfinexPrice, err := strconv.ParseFloat(priceInfo.LastPrice, 64)
	if err != nil {
		return 0, err
	}

	return bitfinexPrice, nil
}

func parseTokenPriceData(data []byte) (float64, int64, error) {
	var roundData struct {
		Answer    *big.Int
		UpdatedAt *big.Int
	}

	err = chainlinkABI.UnpackIntoInterface(&roundData, "latestRoundData", data)
	if err != nil {
		return 0, 0, err
	}

	// uint80 타입인 answer를 float64로 변환하여 가격 리턴
	tokenPrice, _ := new(big.Float).SetInt(roundData.Answer).Float64()

	// updatedAt 값을 int64로 변환하여 타임스탬프 리턴
	updatedAt := roundData.UpdatedAt.Int64()

	return tokenPrice, updatedAt, nil
}

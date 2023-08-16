package main

import (
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var chainlinkABI abi.ABI
var err error

// Chainlink 컨트랙트 주소
var chainlinkAddress = []common.Address{
	common.HexToAddress("0x0630521aC362bc7A19a4eE44b57cE72Ea34AD01c"), //DAI 토큰피드
	common.HexToAddress("0x143db3CEEfbdfe5631aDD3E50f7614B6ba708BA7"), //ETH 토큰피드
}

func init() {
	chainlinkABI, err = abi.JSON(strings.NewReader("ChainLink.JSON"))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}
}

// 토큰 가격 정보 조회
func getLatestTokenData(address common.Address) (float64, error) {
	client, err := ethclient.Dial("https://data-seed-prebsc-1-s1.binance.org:8545/") // Ethereum 노드에 연결
	if err != nil {
		return 0, err
	}

	// latestRoundData() 함수 호출용 데이터 생성
	data, err := chainlinkABI.Pack("latestRoundData")
	if err != nil {
		return 0, err
	}

	var tokenPrice float64

	for _, address := range chainlinkAddress {
		callMsg := ethereum.CallMsg{
			To:   &address,
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
		tokenPrice += currentTokenPrice // Here, I'm just adding up the prices, but you might want to handle the results differently
	}

	return tokenPrice, nil
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

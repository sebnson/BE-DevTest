package main

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var chainlinkABI = []byte("...")                  // Chainlink 컨트랙트 ABI 바이트 코드
var chainlinkAddress = common.HexToAddress("...") // Chainlink 컨트랙트 주소

// Chainlink 컨트랙트로부터 토큰 가격 정보 조회
func getLatestTokenData() (float64, error) {
	client, err := ethclient.Dial("https://bsc-testnet-url") // Ethereum 노드에 연결
	if err != nil {
		return 0, err
	}

	parsedABI, err := abi.JSON(bytes.NewReader(chainlinkABI))
	if err != nil {
		return 0, err
	}

	data, err := client.CallContract(nil, ethereum.CallMsg{
		To:   &chainlinkAddress,
		Data: parsedABI.Pack("latestRoundData"), // latestRoundData() 함수 호출
	}, nil)
	if err != nil {
		return 0, err
	}

	// 데이터 파싱 및 가격 리턴
	price, _, err := parseTokenPriceData(data)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func parseTokenPriceData(data []byte) (float64, int64, error) {
	// Chainlink 컨트랙트 latestRoundData() 함수의 반환값 파싱
	parsedABI, err := abi.JSON(bytes.NewReader(chainlinkABI))
	if err != nil {
		return 0, 0, err
	}

	var roundData struct {
		Answer    *big.Int
		UpdatedAt *big.Int
	}

	err = parsedABI.UnpackIntoInterface(&roundData, "latestRoundData", data)
	if err != nil {
		return 0, 0, err
	}

	// uint80 타입인 answer를 float64로 변환하여 가격 리턴
	price := new(big.Float).SetInt(roundData.Answer).Float64()

	// updatedAt 값을 int64로 변환하여 타임스탬프 리턴
	updatedAt := roundData.UpdatedAt.Int64()

	return price, updatedAt, nil
}

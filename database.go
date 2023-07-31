package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	_ "github.com/go-sql-driver/mysql"

)

const (
	dbUsername = "mysql_username"
	dbPassword = "mysql_password"
	dbName     = "database_name"
)

func connectToDB() (*sql.DB, error) {
	// 데이터베이스 연결 설정
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName))
	if err != nil {
		return nil, err
	}

	return db, nil
}


func saveTokenPrice(tokenSymbol string, price float64) error {
	db, err := connectToDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 데이터베이스에 토큰 가격 정보 저장
	query := "INSERT INTO token_prices (symbol, price, timestamp) VALUES (?, ?, ?)"
	_, err = db.Exec(query, tokenSymbol, tokenPrice, time.Now())
	if err != nil {
		return err 
	}

	return nil
}

func getTokenPrice(tokenSymbol string) (float64, error) {
	db, err := connectToDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// 토큰 가격 정보 조회
	query := "SELECT price FROM token_prices WHERE symbol = ? ORDER BY timestamp DESC LIMIT 1"
	row := db.QueryRow(query, tokenSymbol)

	var price float64
	err = row.Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("No token price found for the given token symbol")
		}
		return 0, err
	}
	return 0, nil
}

// 특정 시간 구간 동안의 평균 토큰 가격 정보 조회
func getAverageTokenPrice(tokenSymbol, source string, startTime, endTime time.Time)(float64, error){
	db, err := connectToDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// 주어진 시간 구간에 해당하는 토큰 가격들 조회
	query := "SELECT price FROM token_prices WHERE symbol = ? AND timestamp BETWEEN ? AND ?"
	rows, err := db.Query(query, tokenSymbol, startTime, endTime)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var totalPrice float64
	var count int

	// 조회된 토큰 가격들의 합, 개수 계산
	for rows.Next() {
		var price float64
		err := rows.Scan(&price)
		if err != nil {
			return 0, err
		}
		totalPrice += price
		count++
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	// 평균 계산
	if count > 0 {
		avgPrice := totalPrice / float64(count)
		return avgPrice, nil
	}

	return 0, fmt.Errorf("No token prices found for given time period")
}

## BE-DevTest: DeFi서비스 백엔드 개발 인턴

사용언어: Golang

데이터베이스: MySQL

----

#### main.go
scheduler 함수를 호출하여 스케줄러를 시작하고, startAPIServer 함수를 호출하여 API 서버를 시작합니다. Chainlink 컨트랙트와 Bitfinex API에서 토큰 가격 정보를 조회하고 데이터베이스에 저장하는 fetchDataAndSave 함수를 포함합니다.

#### contract.go: 
Ethereum 블록체인에서 Chainlink 컨스랙트로부터 토큰 가격 정보를 조회합니다.
getLatestTokenData 함수 호출 시 Ethereum 노드에 연결하여 Chainlink 컨트랙트의 latestRoundData 함수를 호출하여 토큰 가격 정보를 가져온 후, 토큰의 가격과 업데이트 시간을 반환합니다.

#### api.go: 
HTTP API 서버를 구현하는 부분으로, 엔드포인트(/token/price)로부터 토큰 심볼을 받아와 데이터베이스에서 해당 토큰의 가격 정보를 조회하고, JSON 형식으로 클라이언트에게 응답합니다. 

#### database.go: 
MySQL 데이터베이스와 상호작용하는 코드가 들어 있는 파일입니다. connectToDB 함수를 사용하여 데이터베이스에 연결합니다. 데이터베이스에 토큰 가격 정보를 저장하는 saveTokenPrice 함수와 토큰 가격 정보를 조회하는 getTokenPrice 함수, 특정 시간 구간동안의 평균 토큰 가격 정보를 조회하는 getAverageTokenPrice 함수를 포함합니다.

#### utils.go: 
HTTP GET 요청을 보내고 응답 데이터를 읽어오는 기능을 수행합니다.


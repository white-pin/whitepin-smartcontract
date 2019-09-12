package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type EvaluationChaincode struct {
}

type RecordType int
const RecordTypeUser RecordType = 1
const RecordTypeTrade RecordType = 2
const RecordTypeScoreTemp RecordType = 3

// TODO data put, get 공통화
// TODO Error 발생 공통화
// TODO 체인코드 로그 구문 추가
// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(EvaluationChaincode))
	if err != nil {
		fmt.Printf("Error starting Evaluation Chaincode: %s", err)
	}
}

// Init initializes chaincode
// Properties 설정 (default)
func (t *EvaluationChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Printf("Init Evaluation Chaincode.")

	// properties 설정
	err := SetProperties(stub, defaultEvaluationLimit, defaultOpenScoreDuration)
	if err != nil {
		return shim.Error(err.Error())
	}

	// total data 설정
	if err := AddUser(stub, total_user); err != nil {
		shim.Error(err.Error())
	}

	return shim.Success(nil)
}


func (t *EvaluationChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response  {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	switch function {
	case "getProperties" : return t.getProperties(stub) // 프로퍼티 조회
	case "setProperties" : return t.setProperties(stub, args) // 프로퍼티 설정
	case "addUser": return t.addUser(stub, args) // 사용자 생성
	case "queryUser": return t.queryUser(stub, args) // 사용자 조회
	case "createTrade": return t.createTrade(stub, args) // 거래 생성
	case "queryTradeWithId": return t.queryTradeWithId(stub, args) // 거래 조회
	case "queryScoreTemp": return t.queryScoreTemp(stub, args) // 임시 평가정수 조회
	case "queryTradeWithCondition": return t.queryTradeWithCondition(stub, args) // 거래 조회 (query string 사용)
	case "queryScoreTempWithTradeId": return t.queryScoreTempWithTradeId(stub, args) // 임시 평가정수 조회 (query string 사용, tradeId로만 조회 가능)
	case "closeTrade": return t.closeTrade(stub, args) // 거래 완료 처리 (판매자 또는 구매자)
	case "enrollTempScore": return t.enrollTempScore(stub, args) // 임시 평가점수 등록 (판매자 또는 구매자)
	case "enrollScore": return t.enrollScore(stub, args) // 거래 점수 등록 (판매자, 구매자 동시에)
	default:
		err := errors.Errorf("No matched function. : %s", function)
		return shim.Error(err.Error())
	}
}

// 프로퍼티 조회
func (t *EvaluationChaincode) getProperties(stub shim.ChaincodeStubInterface) pb.Response {
	prpty, err := GetProperties(stub)
	if err != nil {
		shim.Error(err.Error())
	}

	byteData, err := json.Marshal(prpty)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(byteData)
}

// 프로퍼티 설정
// args[0] : 평가 입력 기다려주는 시간 (default 14일, 1,209,600 = 14 * 24 * 60 * 60) 이시간 이후에는 0점 처리
// args[1] : 거래 당사자들의 모든 평가 입력 후 공개하기 까지 걸리는 시간 (default 5일, 432,000 = 5 * 24 * 60 * 60)
func (t *EvaluationChaincode) setProperties(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	err := SetProperties(stub, args[0], args[1])
	if err != nil {
		shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// 사용자 생성
// args[0] : 사용자 토큰
func (t *EvaluationChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	err := AddUser(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("Add a new user [%s] finished.\n", args[0])
	return shim.Success(nil)
}


// 사용자 조회
// args[0] : 사용자 토큰
func (t *EvaluationChaincode) queryUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	byteData, err := GetUser(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Get user successfully.")
	return shim.Success(byteData)
}


// 거래 생성
// args[0] : 거래 ID
// args[1] : 서비스 코드
// args[2] : 판매자 토큰
// args[3] : 구매자 토큰
func (t *EvaluationChaincode) createTrade(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	err := CreateTrade(stub, args[0], args[1], args[2], args[3])
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("Add a new trade [%s] finished.\n", args[0])
	return shim.Success(nil)
}


// 거래 조회
// args[0] : 거래 ID
func (t *EvaluationChaincode) queryTradeWithId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	byteData, err := GetTradeWithId(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Get trade successfully.")
	return shim.Success(byteData)
}


// Temp 점수 조회
// args[0] : Temp 점수 키
func (t *EvaluationChaincode) queryScoreTemp(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	scoreTemp, err := GetScoreTempWithKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if scoreTemp == nil {
		err = errors.New("There is no matched data.")
		return shim.Error(err.Error())
	}

	fmt.Println("Get Temp score successfully.")
	return shim.Success(scoreTemp)
}


// 거래 조회 query 작성 후 추가
// args[0] : query string. (거래는 다양하게 불러올 필요가 있으므로 query string 자체를 변수로 받도록)
func (t *EvaluationChaincode) queryTradeWithCondition(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	byteData, err := GetTradeWithQueryString(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(byteData)
}


// 임시평가점수 조회 query 작성 후 추가
// args[0] : query string. (거래는 다양하게 불러올 필요가 있으므로 query string 자체를 변수로 받도록)
func (t *EvaluationChaincode) queryScoreTempWithTradeId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	byteData, err := GetScoreTempWithTradeId(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if byteData == nil {
		err = errors.New("There is no matched data.")
		return shim.Error(err.Error())
	}

	return shim.Success(byteData)
}


// 거래 완료 처리. 판매자, 구매자 둘다 완료처리해야 최종 완료됨. args[1]은 userTkn으로, 판매자, 구매자인지 판별하여 해당 대상에 대해서 완료처리.
// args[0] : 거래 ID
// args[1] : 사용자 토큰 (판매자든 구매자든)
func (t *EvaluationChaincode) closeTrade(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	err := CloseTrade(stub , args[0], args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}


// Temp 점수 등록(구매자 or 판매자) : 둘다 등록해야 최종 등록됨
//args[0] := tradeId
//args[1] := userTkn
//args[2] := scoreOrigin // "[3,4,5]" 의 format
//args[3] := key (encryption에 사용될)
func (t *EvaluationChaincode) enrollTempScore(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// AES_GCM 생성 및 키 설정
	aes_gcm, err := GCM_Key(args[3])
	if err != nil {
		return shim.Error(err.Error())
	}

	// nonce 설정. (거래 데이터셋에서 시간 가져오기)
	var trade Trade
	rsltTrade, err := GetTradeWithId(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	err = json.Unmarshal(rsltTrade, &trade)
	if err != nil {
		return shim.Error(err.Error())
	}

	date := trade.Date
	aes_gcm.nonce = fmt.Sprintf("%d%02d%02dA%02d%02d%02d%09d",
		date.Year(), date.Month(), date.Day(),
		date.Hour(), date.Minute(), date.Second(), date.Nanosecond())

	aes_gcm.plainTxt = args[2] // "[3,4,5]" 의 format

	// 판매자인지 구매자인지 판별
	var division string
	switch args[1] {
	case trade.SellerTkn: division = "buy" // 사용자의 토큰이 판매자 토큰과 일치한다는 것은, 점수를 매기는 주체 = 판매자, 점수가 매겨지는 대상 = 구매자. 즉 구매자가 받는 점수
	case trade.BuyerTkn: division = "sell" // 사용자의 토큰이 구매자 토큰과 일치한다는 것은, 점수를 매기는 주체 = 구매자, 점수가 매겨지는 대상 = 판매자. 즉 판매자가 받는 점수
	default:
		return shim.Error(errors.Errorf("User \"%s\" doesn't participated this trade.", args[1]).Error())
	}

	// 점수 암호화
	err = aes_gcm.GCM_encrypt()
	if err != nil {
		return shim.Error(err.Error())
	}

	score := aes_gcm.chipherTxt

	err = SetScoreTempWithTradeId(stub, args[0], score, division)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


// 거래 점수 등록
// args[0] : 거래 ID
// args[1] : 암호화 해제 키
func (t *EvaluationChaincode) enrollScore(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	byteData, err := GetScoreTempWithTradeId(stub, args[0])
	//byteData, err := GetScoreTempWithQueryString(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if byteData == nil {
		err := errors.New("There is no matched data.")
		return shim.Error(err.Error())
	}

	var scoreTemp ScoreTemp
	err = json.Unmarshal(byteData, &scoreTemp)
	if err != nil {
		return shim.Error(err.Error())
	}

	tradeId := scoreTemp.TradeId

	// AES_GCM 생성 및 키 설정
	aes_gcm, err := GCM_Key(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	// nonce 설정. (거래 데이터셋에서 시간 가져오기)
	var trade Trade
	rsltTrade, err := GetTradeWithId(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	err = json.Unmarshal(rsltTrade, &trade)
	if err != nil {
		return shim.Error(err.Error())
	}

	date := trade.Date
	aes_gcm.nonce = fmt.Sprintf("%d%02d%02dA%02d%02d%02d%09d",
		date.Year(), date.Month(), date.Day(),
		date.Hour(), date.Minute(), date.Second(), date.Nanosecond())

	// 복호화 (sell)
	aes_gcm.chipherTxt = scoreTemp.Score.SellScore // "[3,4,5]" 의 암호화된 format
	err = aes_gcm.GCM_decrypt()
	if err != nil {
		return shim.Error(err.Error())
	}
	sellScorePlainTxt := aes_gcm.plainTxt

	// 복호화 (buy)
	aes_gcm.chipherTxt = scoreTemp.Score.BuyScore // "[3,4,5]" 의 암호화된 format
	err = aes_gcm.GCM_decrypt()
	if err != nil {
		return shim.Error(err.Error())
	}
	buyScorePlainTxt := aes_gcm.plainTxt

	var sellScore []int // [3,4,5] format의 int array
	var buyScore []int // [3,4,5] format의 int array
	for _, val := range strings.Split(sellScorePlainTxt[1:len(sellScorePlainTxt)-1], ",") {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return shim.Error(err.Error())
		}
		sellScore = append(sellScore, intVal)
	}
	if len(sellScore) != 3 { // 평가 질문이 3개 이므로 점수도 3개
		err := errors.New("Seller score has error.")
		return shim.Error(err.Error())
	}
	for _, val := range strings.Split(buyScorePlainTxt[1:len(buyScorePlainTxt)-1], ",") {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return shim.Error(err.Error())
		}
		buyScore = append(buyScore, intVal)
	}
	if len(buyScore) != 3 { // 평가 질문이 3개 이므로 점수도 3개
		err := errors.New("Buyer score has error.")
		return shim.Error(err.Error())
	}

	err = EvaluateTrade(stub, tradeId, sellScore, buyScore)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Update User score 구매자, 판매자, 전체
	err = UpdateUserScore(stub, trade.SellerTkn , sellScore, "sell")
	if err != nil {
		return shim.Error(err.Error())
	}
	err = UpdateUserScore(stub, trade.BuyerTkn , buyScore, "buy")
	if err != nil {
		return shim.Error(err.Error())
	}
	err = UpdateTotalScore(stub, sellScore, buyScore)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
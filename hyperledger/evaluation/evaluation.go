package main

import (
	"fmt"
	"strconv"
	"strings"
	
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
)

type EvaluationChaincode struct {
}

type RecordType int
const RecordTypeUser RecordType = 1
const RecordTypeTrade RecordType = 2
const RecordTypeScoreMeta RecordType = 3

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
// ===========================
func (t *EvaluationChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Printf("Init Evaluation Chaincode.")
	return shim.Success(nil)
}


func (t *EvaluationChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response  {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	switch function {
	case "addUser": return t.addUser(stub, args) // 사용자 생성
	case "queryUser": return t.queryUser(stub, args) // 사용자 조회
	case "createTrade": return t.createTrade(stub, args) // 거래 생성
	case "queryTradeWithId": return t.queryTradeWithId(stub, args) // 거래 조회
	case "addScoreMeta": return t.addScoreMeta(stub, args) // 임시 평가점수 저장
	case "queryScoreMeta": return t.queryScoreMeta(stub, args) // 임시 평가정수 조회
	case "queryTradeWithCondition": return t.queryTradeWithCondition(stub, args) // 거래 조회 (query string 사용)
	case "closeTrade": return t.closeTrade(stub, args) // 거래 완료 처리 (판매자 또는 구마재)
	case "enrollMetaScore": return t.enrollMetaScore(stub, args) // 임시 평가점수 등록 (판매자 또는 구마재)
	case "queryMetaScoreWithCondition": return t.queryMetaScoreWithCondition(stub, args) // 임시 평가점수 조회 (query string 사용)
	case "enrollScore": return t.enrollScore(stub, args) // 거래 점수 등록 (판매자, 구매자 동시에)
	default:
		err := errors.Errorf("No matched function. : %s", function)
		return shim.Error(err.Error())
	}
}


// 사용자 생성 +
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


// 사용자 조회 +
func (t *EvaluationChaincode) queryUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	user, err := GetUser(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	byteData,err := json.Marshal(user)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Get user successfully.")
	return shim.Success(byteData)
}


// 거래 생성 +
func (t *EvaluationChaincode) createTrade(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	// TODO 서비스 코드, 판매자 토큰, 구매자 토큰 존재하는지 검증 로직 필요
	err := CreateTrade(stub, args[0], args[1], args[2], args[3])
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("Add a new trade [%s] finished.\n", args[0])
	return shim.Success(nil)
}


// 거래 조회 +
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


// meta 점수 생성 +
func (t *EvaluationChaincode) addScoreMeta(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	err := AddScoreMeta(stub, args[0], args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


// meta 점수 조회 +
func (t *EvaluationChaincode) queryScoreMeta(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	scoreMeta, err := GetScoreMetaWithKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if scoreMeta == nil {
		err = errors.New("There is no matched data.")
		return shim.Error(err.Error())
	}

	fmt.Println("Get meta score successfully.")
	return shim.Success(scoreMeta)
}


// 거래 조회 query 작성 후 추가 +
func (t *EvaluationChaincode) queryTradeWithCondition(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// TODO selector 유효성 검증

	byteData, err := GetTradeWithQueryString(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(byteData)
}


// 거래 완료 처리. 판매자, 구매자 둘다 완료처리해야 최종 완료됨. args[1]은 userTkn으로, 판매자, 구매자인지 판별하여 해당 대상에 대해서 완료처리. +
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


// meta 점수 등록(구매자 or 판매자) : 둘다 등록해야 최종 등록됨 +
func (t *EvaluationChaincode) enrollMetaScore(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	scoreKey := args[0]
	scoreOrigin := args[1] // "[3,4,5]" 의 format
	division := args[2]

	// TODO 암호화하는 부분 변경필요 (현재는 단순 AES Encryption)

	score := scoreOrigin
	//scoreOrigin = scoreOrigin[1:len(scoreOrigin)-1] // "3,4,5" 의 format

	err := SetScoreMetaWithKey(stub, scoreKey, score, division)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


// meta 점수 조회 (query) +
func (t *EvaluationChaincode) queryMetaScoreWithCondition(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	buffer, err := GetScoreMetaWithQueryString(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if buffer == nil {
		err = errors.New("There is no matched data.")
		return shim.Error(err.Error())
	}

	return shim.Success(buffer)
}


// 거래 점수 등록
func (t *EvaluationChaincode) enrollScore(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// TODO query 부분 작성되면 주석 변경

	byteData, err := GetScoreMetaWithQueryString(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if byteData == nil {
		err := errors.New("There is no matched data.")
		return shim.Error(err.Error())
	}

	var scoreMeta ScoreMeta
	err = json.Unmarshal(byteData, &scoreMeta)
	if err != nil {
		return shim.Error(err.Error())
	}

	tradeId := scoreMeta.TradeId
	sellScoreChiper := scoreMeta.Score.SellScore // "[3,4,5]" 의 format
	buyScoreChiper := scoreMeta.Score.BuyScore // "[3,4,5]" 의 format

	var sellScore []int // [3,4,5] format의 int array
	var buyScore []int // [3,4,5] format의 int array
	for _, val := range strings.Split(sellScoreChiper[1:len(sellScoreChiper)-1], ",") {
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
	for _, val := range strings.Split(buyScoreChiper[1:len(buyScoreChiper)-1], ",") {
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

	return shim.Success(nil)
}
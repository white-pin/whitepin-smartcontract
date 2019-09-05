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
	case "queryTradeWithUserTkn": return t.queryTradeWithUserTkn(stub, args) // 거래 조회 (query string 사용)
	default:
		err := errors.Errorf("No matched function. : %s", function)
		return shim.Error(err.Error())
	}
}


// 사용자 생성
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


// 거래 생성
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


// 거래 조회
func (t *EvaluationChaincode) queryTradeWithId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	trade, err := GetTradeWithId(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	byteData, err := json.Marshal(trade)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Get trade successfully.")
	return shim.Success(byteData)
}


// meta 점수 생성
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


// meta 점수 조회
func (t *EvaluationChaincode) queryScoreMeta(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	scoreMeta, err := GetScoreMetaWithKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	byteData, err := json.Marshal(scoreMeta)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Get meta score successfully.")
	return shim.Success(byteData)
}


// 거래 조회 query 작성 후 추가
func (t *EvaluationChaincode) queryTradeWithUserTkn(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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


// 거래 완료 처리(division, "sell" : 구매자 / "buy" : 판매자) : 둘다 완료처리해야 최종 완료됨
func (t *EvaluationChaincode) closeTrade(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	tradeId := args[0]
	userTkn := args[1]
	division := args[2]

	err := CloseTrade(stub , tradeId, userTkn, division)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}


// meta 점수 등록(구매자 or 판매자) : 둘다 등록해야 최종 등록됨
func (t *EvaluationChaincode) enrollMetaScore(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	scoreKey := args[0]
	scoreOrigin := args[1] // "[3,4,5]" 의 format
	division := args[2]
	//aeskey := args[3]

	// TODO 암호화하는 부분 변경필요 (현재는 단순 AES Encryption)
	//chiper, err := sw.AESCBCPKCS7Encrypt(bytes.NewBufferString(aeskey).Bytes(), bytes.NewBufferString(scoreOrigin).Bytes())
	//if err != nil {
	//	return shim.Error(err.Error())
	//}
	score := scoreOrigin
	//scoreOrigin = scoreOrigin[1:len(scoreOrigin)-1] // "3,4,5" 의 format

	err := SetScoreMetaWithKey(stub, scoreKey, score, division)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


// TODO meta 점수 조회


// 거래 점수 등록
func (t *EvaluationChaincode) enrollScore(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	tradeId := args[0]
	// TODO query 부분 작성되면 주석 변경
	// ------------------------------------------------------------
	//scoremeta.GetScoreMetaWithTradeId(stub, tradeId)
	scoreKey := args[1]
	scoreMeta, err := GetScoreMetaWithKey(stub, scoreKey)
	// ------------------------------------------------------------

	if err != nil {
		return shim.Error(err.Error())
	}

	sellScoreChiper := scoreMeta.Score.SellScore // "[3,4,5]" 의 format
	buyScoreChiper := scoreMeta.Score.BuyScore // "[3,4,5]" 의 format

	var sellScore []int // [3, 4, 5] format의 int array
	var buyScore []int // [3, 4, 5] format의 int array
	for _, val := range strings.Split(sellScoreChiper[1:len(sellScoreChiper)-1], ",") {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return shim.Error(err.Error())
		}
		sellScore = append(sellScore, intVal)
	}
	for _, val := range strings.Split(buyScoreChiper[1:len(buyScoreChiper)-1], ",") {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return shim.Error(err.Error())
		}
		buyScore = append(buyScore, intVal)
	}

	err = EvaluateTrade(stub, tradeId, sellScore, buyScore)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
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
const defaultOrder string = "desc"

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
	if err := AddUser(stub, TotalUser); err != nil {
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
	case "queryTradeWithQueryString" : return t.queryTradeWithQueryString(stub, args) // queryString으로 거래이력 조회. (메소드로 제공하지 않는 기능들에 대해서 조회 가능하도록)
	case "queryTradeWithUser" : return t.queryTradeWithUser(stub, args) // 사용자로 거래이력 조회. (판매, 구매, 모두 : sell, buy, all)
	case "queryTradeWithUserService" : return t.queryTradeWithUserService(stub, args) // 사용자, 서비스 코드로 거래이력 조회. (판매, 구매, 모두 : sell, buy, all)
	case "queryTradeWithService" : return t.queryTradeWithService(stub, args) // 서비스 코드로 거래이력 조회.
	case "queryTradeWithId": return t.queryTradeWithId(stub, args) // 거래 조회
	case "queryScoreTemp": return t.queryScoreTemp(stub, args) // 임시 평가정수 조회
	case "queryScoreTempWithTradeId": return t.queryScoreTempWithTradeId(stub, args) // 임시 평가정수 조회 (query string 사용, tradeId로만 조회 가능)
	case "closeTrade": return t.closeTrade(stub, args) // 거래 완료 처리 (판매자 또는 구매자)
	case "enrollTempScore": return t.enrollTempScore(stub, args) // 임시 평가점수 등록 (판매자 또는 구매자)
	case "enrollScore": return t.enrollScore(stub, args) // 거래 점수 등록 (판매자, 구매자 동시에)
	// TODO history
	default:
		return shim.Error(errors.Errorf("No matched function. : %s", function).Error())
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

	byteData, err := getDataWithKey(stub, args[0])
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


	// Update User score 구매자, 판매자, 전체
	err = InitUserScore(stub, args[2], "sell")
	if err != nil {
		return shim.Error(err.Error())
	}
	err = InitUserScore(stub, args[3], "buy")
	if err != nil {
		return shim.Error(err.Error())
	}
	err = InitUserScore(stub, TotalUser, "tot")
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

	scoreTemp, err := getDataWithKey(stub, args[0])
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


// 거래 조회 query string 사용. (거래는 다양하게 불러올 필요가 있으므로 query string 자체를 변수로 받도록)
// args[0] : query string.
func (t *EvaluationChaincode) queryTradeWithQueryString(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	byteData, err := getQueryString(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(byteData)
}
// 거래 조회 (사용자 토큰으로 조회)
// args[0] : 사용자 정보, UserTkn
// args[1] : 사용자의 조건 (판매 : "sell", 구매 : "buy", 둘다 : "all")
// args[2] : 시간에 대한 ordering. (default is desc.)
// args[3] : 일반조회, 페이지 조회(normal, page, default is normal)
// args[4] (optional) : 한 페이지 사이즈
// args[5] (optional) : 페이지 번호
// args[6] (optional) : 북마크 (첫 페이지 조회 또는 default는 "")
func (t *EvaluationChaincode) queryTradeWithUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 && len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 4 or 7")
	}
	var byteData []byte
	var err error
	userTkn := args[0]
	condition := args[1]
	ordering := defaultOrder
	if args[2] != "" {
		ordering = args[2]
	}
	reportType := args[3]
	pageSize := args[4]
	pageNum := args[5]
	bookmark := args[6]
	queryString := ""

	switch condition {
	case "sell" :
		queryString = "{\"selector\":{\"SellerTkn\":\""+userTkn+"\",\"RecType\":2},\"sort\":[{\"Date\":\""+ordering+"\"}]}"
	case "buy" :
		queryString = "{\"selector\":{\"BuyerTkn\":\""+userTkn+"\",\"RecType\":2},\"sort\":[{\"Date\":\""+ordering+"\"}]}"
	case "all" : // 전체 조회는 sort가 안되도록 설계가 됐기 때문에 sorting 불가능.
		queryString = "{\"selector\":{\"$or\":[{\"SellerTkn\":\""+userTkn+"\",\"RecType\":2},{\"BuyerTkn\":\""+userTkn+"\",\"RecType\":2}]}}"
	default :
		return shim.Error(errors.New("Unexpected user condition. It can be 'sell', 'buy' and 'all'").Error())
	}

	switch reportType {
	case "normal" :
		byteData, err = getQueryString(stub, queryString)
		if err != nil {
			return shim.Error(err.Error())
		}
	case "page" :
		byteData, err = getPageDataWithQueryString(stub, queryString, pageSize, pageNum, bookmark)
		if err != nil {
			return shim.Error(errors.New("Paging Query Failed").Error())
		}
	default:
		return shim.Error(errors.New("Unexpected reoport option. 'normal' and 'page' is available.").Error())
	}

	return shim.Success(byteData)
}
// 거래 조회 (사용자 토큰, 서비스 코드로 조회)
// args[0] : 사용자 정보, UserTkn
// args[1] : 서비스 코드, ServiceCode
// args[2] : 사용자의 조건 (판매 : "sell", 구매 : "buy", 둘다 : "all")
// args[3] : 시간에 대한 ordering. (default is desc.)
// args[4] : 일반조회, 페이지 조회(normal, page, default is normal)
// args[5] (optional) : 한 페이지 사이즈
// args[6] (optional) : 페이지 번호
// args[7] (optional) : 북마크 (첫 페이지 조회 또는 default는 "")
func (t *EvaluationChaincode) queryTradeWithUserService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 && len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 5 or 8")
	}
	var byteData []byte
	var err error
	userTkn := args[0]
	serviceCode := args[1]
	condition := args[2]
	ordering := defaultOrder
	if args[3] != "" {
		ordering = args[3]
	}
	reportType := args[4]
	pageSize := args[5]
	pageNum := args[6]
	bookmark := args[7]

	queryString := ""
	switch condition {
	case "sell" :
		queryString = "{\"selector\":{\"SellerTkn\":\""+userTkn+"\",\"ServiceCode\":\""+serviceCode+"\",\"RecType\":2},\"sort\":[{\"Date\":\""+ordering+"\"}]}"
	case "buy" :
		queryString = "{\"selector\":{\"BuyerTkn\":\""+userTkn+"\",\"ServiceCode\":\""+serviceCode+"\",\"RecType\":2},\"sort\":[{\"Date\":\""+ordering+"\"}]}"
	case "all" : // 전체 조회는 sort가 안되도록 설계가 됐기 때문에 sorting 불가능.
		queryString = "{\"selector\":{\"$or\":[{\"SellerTkn\":\""+userTkn+"\",\"ServiceCode\":\""+serviceCode+"\",\"RecType\":2},{\"BuyerTkn\":\""+userTkn+"\",\"ServiceCode\":\""+serviceCode+"\",\"RecType\":2}]}}"
	default :
		return shim.Error(errors.New("Unexpected user condition. It can be 'sell', 'buy' and 'all'").Error())
	}

	switch reportType {
	case "normal" :
		byteData, err = getQueryString(stub, queryString)
		if err != nil {
			return shim.Error(err.Error())
		}
	case "page" :
		byteData, err = getPageDataWithQueryString(stub, queryString, pageSize, pageNum, bookmark)
		if err != nil {
			return shim.Error(errors.New("Paging Query Failed").Error())
		}
	default:
		return shim.Error(errors.New("Unexpected reoport option. 'normal' and 'page' is available.").Error())
	}

	return shim.Success(byteData)
}
// 거래 조회 (서비스 코드로 조회)
// args[0] : 서비스 코드, ServiceCode
// args[1] : 시간에 대한 ordering. (default is desc.)
// args[2] : 일반조회, 페이지 조회(normal, page, default is normal)
// args[3] (optional) : 한 페이지 사이즈
// args[4] (optional) : 페이지 번호
// args[5] (optional) : 북마크 (첫 페이지 조회 또는 default는 "")
func (t *EvaluationChaincode) queryTradeWithService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 && len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 3 or 6")
	}
	var byteData []byte
	var err error
	serviceCode := args[0]
	ordering := defaultOrder
	if args[1] != "" {
		ordering = args[1]
	}
	reportType := args[2]
	pageSize := args[3]
	pageNum := args[4]
	bookmark := args[5]
	queryString := "{\"selector\":{\"ServiceCode\":\""+serviceCode+"\",\"RecType\":2},\"sort\":[{\"Date\":\""+ordering+"\"}]}"

	switch reportType {
	case "normal" :
		byteData, err = getQueryString(stub, queryString)
		if err != nil {
			return shim.Error(err.Error())
		}
	case "page" :
		byteData, err = getPageDataWithQueryString(stub, queryString, pageSize, pageNum, bookmark)
		if err != nil {
			return shim.Error(errors.New("Paging Query Failed").Error())
		}
	default:
		return shim.Error(errors.New("Unexpected reoport option. 'normal' and 'page' is available.").Error())
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
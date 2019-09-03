package main

import (
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
)

type EvaluationChaincode struct {
}

// User data-set
type user struct {
	UserTkn string `json:"UserTkn"`
	SellAmt int64 `json:"SellAmt"`
	BuyAmt int64 `json:"BuyAmt"`
	SellAvg struct {
		TotAvg    float32 `json:"TotAvg"`
		EvalAvg01 float32 `json:"EvalAvg01"`
		EvalAvg02 float32 `json:"EvalAvg02"`
		EvalAvg03 float32 `json:"EvalAvg03"`
	} `json:"SellAvg"`
	BuyAvg struct {
		TotAvg    float32 `json:"TotAvg"`
		EvalAvg01 float32 `json:"EvalAvg01"`
		EvalAvg02 float32 `json:"EvalAvg02"`
		EvalAvg03 float32 `json:"EvalAvg03"`
	} `json:"BuyAvg"`
	TradeAvg struct {
		TotAvg    float32 `json:"TotAvg"`
		EvalAvg01 float32 `json:"EvalAvg01"`
		EvalAvg02 float32 `json:"EvalAvg02"`
		EvalAvg03 float32 `json:"EvalAvg03"`
	} `json:"BuyAvg"`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(EvaluationChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *EvaluationChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}


func (t *EvaluationChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response  {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "createTrade" { //create a new Trade
		return t.createTrade(stub, args)
	} else if function == "closeTrade" { //change owner of a specific marble
		return t.transferMarble(stub, args)
	}

	return shim.Success(nil)
}


// TODO 거래 생성
func (t *EvaluationChaincode) createTrade(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

// TODO 거래 종료(구매자 or 판매자) : 둘다 종료해야 최종 종료처리됨
func (t *EvaluationChaincode) closeTrade(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

// TODO meta 점수 등록(구매자 or 판매자) : 둘다 등록해야 최종 등록됨
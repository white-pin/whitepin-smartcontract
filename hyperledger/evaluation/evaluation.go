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

// Trade data-set
type trade struct {
	TradeId string `json:"TradeId"`
	SellerTkn string `json:"SellerTkn"`
	BuyerTkn string `json:"BuyerTkn"`
	Date time.Time `json:"Date"`
	Close struct {
		SellDone bool `json:"SellDone"`
		BuyDone bool `json:"BuyDone"`
		SellDate time.Time `json:"SellDate"`
		BuyDate time.Time `json:"BuyDate"`
	}
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
	if function == "initMarble" { //create a new marble
		return t.initMarble(stub, args)
	} else if function == "transferMarble" { //change owner of a specific marble
		return t.transferMarble(stub, args)
	}

	return shim.Success(nil)
}


func (t *EvaluationChaincode) initMarble(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
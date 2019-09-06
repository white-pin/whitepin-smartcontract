package main

import (
	"bytes"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"time"
)

// Trade data-set
type Trade struct {
	TradeId string `json:"TradeId"`
	ServiceCode string `json:"ServiceCode"`
	SellerTkn string `json:"SellerTkn"`
	BuyerTkn string `json:"BuyerTkn"`
	Date time.Time `json:"Date"`
	Close struct {
		SellDone bool `json:"SellDone"`
		BuyDone bool `json:"BuyDone"`
		SellDate time.Time `json:"SellDate"`
		BuyDate time.Time `json:"BuyDate"`
	}
	Score struct {
		SellScore []int `json:"SellScore"`
		BuyScore []int `json:"BuyScore"`
	} `json:"Score"`
}


// 거래 등록하기 {TradeId, ServiceCode, SellerTkn, BuyerTkn}
func CreateTrade(stub shim.ChaincodeStubInterface, tradeId string, serviceCode string, sellerTkn string, buyerTkn string) error {
	var trade Trade
	trade.TradeId = tradeId
	trade.ServiceCode = serviceCode
	trade.SellerTkn = sellerTkn
	trade.BuyerTkn = buyerTkn
	trade.Date = time.Now()

	// TODO 기존 거래 있는지 검증

	err := putTrade(stub, trade, tradeId)
	if err != nil {
		return err
	}

	return nil
}


// 거래 가져오기 (TradeId : key)
func GetTradeWithId(stub shim.ChaincodeStubInterface, tradeId string) ([]byte, error) {
	trade, err := getTrade(stub, tradeId)
	if err != nil {
		return nil, err
	}
	if trade == nil {
		err := errors.New("There is no matched data.")
		return nil, err
	}

	return trade, nil
}


// 거래 가져오기 query 필요
func GetTradeWithQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}

	buffer := bytes.Buffer{}
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return buffer.Bytes(), nil
}


// 거래 완료하기 (구매자, 판매자 각각) {TradeId, UserTkn}
func CloseTrade(stub shim.ChaincodeStubInterface, tradeId string, userTkn string) error {
	var trade Trade
	bytesData, err := getTrade(stub, tradeId)
	if err != nil {
		return err
	}
	if bytesData == nil {
		err := errors.New("There is no matched data.")
		return err
	}

	err = json.Unmarshal(bytesData, &trade)
	if err != nil {
		return err
	}

	switch userTkn {
	case trade.SellerTkn: trade.Close.SellDone = true
	trade.Close.SellDate = time.Now()
	case trade.BuyerTkn: trade.Close.BuyDone = true
	trade.Close.BuyDate = time.Now()
	default:
		err := errors.Errorf("user %s is not both seller and buyer.", userTkn)
		return err
	}

	err = putTrade(stub, trade, tradeId)
	if err != nil {
		return err
	}

	return nil
}


// 거래 점수 등록 (meta score 로부터 구매자, 판매자 평가점수 동시에 등록)
func EvaluateTrade(stub shim.ChaincodeStubInterface, tradeId string, sellScore []int, buyScore []int) error {
	var trade Trade

	bytesData, err := getTrade(stub, tradeId)
	if err != nil {
		return err
	}
	if bytesData == nil {
		err := errors.New("There is no matched data.")
		return err
	}

	err = json.Unmarshal(bytesData, &trade)
	if err != nil {
		return err
	}

	trade.Score.SellScore = sellScore
	trade.Score.BuyScore = buyScore

	err = putTrade(stub, trade, tradeId)
	if err != nil {
		return err
	}

	return nil
}



func getTrade(stub shim.ChaincodeStubInterface, tradeId string) ([]byte, error) {
	byteData, err := stub.GetState(tradeId)
	if err != nil {
		err = errors.Errorf("Failed to get data. : TradeId is \"%s\"", tradeId)
		return nil, err
	}
	if byteData == nil {
		err := errors.Errorf("There is no trade : TradeId is \"%s\"", tradeId)
		return nil, err
	}

	return byteData, nil
}

func putTrade(stub shim.ChaincodeStubInterface, trade Trade, tradeId string) error {
	inputData, err := json.Marshal(trade)
	if err != nil {
		err := errors.New("Failed to json encoding.")
		return err
	}

	err = stub.PutState(tradeId, inputData)
	if err != nil {
		err := errors.New("Failed to store data.")
		return err
	}

	return nil
}
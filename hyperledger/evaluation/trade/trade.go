package trade

import (
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

// TODO 거래 등록하기 {TradeId, ServiceCode, SellerTkn, BuyerTkn}
func CreateTrade(stub shim.ChaincodeStubInterface, tradeId string, serviceCode string, sellerTkn string, buyerTkn string) error {
	var trade Trade

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

// TODO 거래 가져오기 {TradeId} {UserTkn, sell/buy} {ServiceCode}
func GetTradeWithId(stub shim.ChaincodeStubInterface, tradeId string, serviceCode string, sellerTkn string, buyerTkn string) error {

	return nil
}

// TODO 거래 가져오기 {UserTkn, sell/buy} {ServiceCode}
func GetTradeWithUserTkn(stub shim.ChaincodeStubInterface, userTkn string) error {

	return nil
}

// TODO 거래 가져오기 {ServiceCode}
func GetTradeWithServiceCode(stub shim.ChaincodeStubInterface, serviceCode string) error {

	return nil
}

// TODO 거래 완료하기 (구매자, 판매자 각각) {TradeId, UserTkn, sell/buy}
func CloseTrade(stub shim.ChaincodeStubInterface, tradeId string, userTkn string, division string) error {

	return nil
}
package scoremeta

import (
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ScoreMeta data는 점수가 실제 거래 데이터에 저장되기 전에 암호화 된 값으로 임시로 저장하고 있는다.
// 특정 조건이 달성되면 점수를 실제 거래 데이터에 저장하고 공개한다.
// TODO 위에서 언급한 특정 조건 지정.
type ScoreMeta struct {
	TradeId string `json:"TradeId"`
	Score struct {
		SellScore string `json:"SellScore"`
		BuyScore string `json:"BuyScore"`
	}
}

// 점수
func GetScoreData(stub shim.ChaincodeStubInterface, tradeId string) ([]byte, error) {
	byteData, err := stub.GetState(tradeId)
	if err != nil {
		return nil, err
	}
	return byteData, nil
}
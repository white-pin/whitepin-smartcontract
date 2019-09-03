package scoremeta

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
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


// 점수 가져오기
func GetScoreData(stub shim.ChaincodeStubInterface, tradeId string) (ScoreMeta, error) {
	var scoreMeta ScoreMeta

	byteData, err := stub.GetState(tradeId)
	if err != nil {
		err = errors.Errorf("Failed to get data : TradeId is \"%s\"", tradeId)
		return ScoreMeta{}, err
	}

	err = json.Unmarshal(byteData, &scoreMeta)
	if err != nil {
		err = errors.New("Failed to json decoding.")
		return ScoreMeta{}, err
	}

	return scoreMeta, nil
}


// 점수 설정 division : "sell", "buy". sell인 경우는 판매자의 점수이고(구매자가 매긴 점수), buy인 경우는 구매자의 점수이다.(판매자가 매긴 점수)
func SetScore(stub shim.ChaincodeStubInterface, tradeId string, division string, score string) error {
	var scoreMeta ScoreMeta

	byteData, err := stub.GetState(tradeId)
	if err != nil {
		err = errors.Errorf("Failed to get data : TradeId is \"%s\"", tradeId)
		return err
	}

	err = json.Unmarshal(byteData, &scoreMeta)
	if err != nil {
		err = errors.New("Failed to json decoding.")
		return err
	}

	switch division {
	case "sell": scoreMeta.Score.SellScore = score
	case "buy": scoreMeta.Score.BuyScore = score
	default:
		err := errors.New("Division is wrong. Available value is \"sell\" and \"buy\"")
		return err
	}

	inputData, err := json.Marshal(scoreMeta)
	if err != nil {
		err := errors.New("Failed to json encoding.")
		return err
	}

	err = stub.PutState(tradeId, inputData)
	if err != nil {
		err := errors.New("Failed to store data.")
		return err
	}
	fmt.Printf("Set \"%s\" score successfuly.", division)

	return nil
}
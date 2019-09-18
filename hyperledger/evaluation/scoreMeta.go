package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"time"
)

// ScoreTemp data는 점수가 실제 거래 데이터에 저장되기 전에 암호화 된 값으로 임시로 저장하고 있는다.
// 조건이 달성되면 점수를 실제 거래 데이터에 저장하고 공개한다.
// 조건 : 판매자와 구매자 모두 구매완료한 날로부터 14일후.
// 조건 : 판매자와 구매자가 서로 평가를 완료한 일로부터 5일후
type ScoreTemp struct {
	RecType RecordType `json:"RecType"` // ScoreTemp : 3
	ScoreKey string `json:"ScoreKey"`
	TradeId string `json:"TradeId"`
	UserTkn string `json:"UserTkn"`
	ExpiryDate time.Time `json:"ExpiryDate"`
	Score struct {
		SellScore string `json:"SellScore"`
		BuyScore string `json:"BuyScore"`
	} `json:"Score"`
	IsExpired bool `json:"IsExpired"`
	//Sync bool `json:"Sync"`
}


// 점수 생성
func AddScoreTemp(stub shim.ChaincodeStubInterface, scoreKey string, tradeId string) error {
	var scoreTemp ScoreTemp

	scoreTemp.RecType = RecordTypeScoreTemp
	scoreTemp.ScoreKey = scoreKey
	scoreTemp.TradeId = tradeId
	scoreTemp.IsExpired = false
	//scoreTemp.Sync = false

	// 같은 scoreKey로 이미 임시 점수를 생성했는지 확인
	_, err := getDataWithKey(stub, scoreKey)
	if err == nil {
		err := errors.Errorf("The score key already used. : %s", scoreKey)
		return err
	}

	byteData, err := json.Marshal(scoreTemp)
	if err != nil {
		err = errors.New("Failed to json encoding.")
		return err
	}

	err = stub.PutState(scoreKey, byteData)
	if err != nil {
		return err
	}

	return nil
}


// 점수 설정 (query) division : "sell", "buy". sell인 경우는 판매자의 점수이고(구매자가 매긴 점수), buy인 경우는 구매자의 점수이다.(판매자가 매긴 점수)
func SetScoreTempWithTradeId(stub shim.ChaincodeStubInterface, tradeId string, score string, division string) error {
	var scoreTemp ScoreTemp
	bothSetScoreFlag := false

	byteData, err := GetScoreTempWithTradeId(stub, tradeId)
	if err != nil {
		return err
	}
	if byteData == nil {
		err := errors.New("There is no matched data.")
		return err
	}

	err = json.Unmarshal(byteData, &scoreTemp)
	if err != nil {
		return errors.New("Failed to json decoding.")
	}

	switch division {
	case "sell":
		scoreTemp.Score.SellScore = score
		if scoreTemp.Score.BuyScore != "" {
			bothSetScoreFlag = true
		}
	case "buy":
		scoreTemp.Score.BuyScore = score
		if scoreTemp.Score.SellScore != "" {
			bothSetScoreFlag = true
		}
	default:
		err := errors.New("Division is wrong. Available value is \"sell\" and \"buy\"")
		return err
	}

	// 거래 당사자 모두 리뷰를 등록한 경우 공개일이 지나면 공개하도록 만료일을 변경한다.
	if bothSetScoreFlag {
		prpty, err := GetProperties(stub)
		if err != nil {
			return err
		}

		scoreTemp.ExpiryDate = time.Now().Add(prpty.OpenScoreDuration) // 지금으로부터 + 평가기간 limit
	}

	inputData, err := json.Marshal(scoreTemp)
	if err != nil {
		err := errors.New("Failed to json encoding.")
		return err
	}

	err = stub.PutState(scoreTemp.ScoreKey, inputData)
	if err != nil {
		err := errors.New("Failed to store data.")
		return err
	}
	fmt.Printf("Set \"%s\" score successfuly.", division)

	return nil
}


// 임시 평가점수 삭제
func DelScoreTemp(stub shim.ChaincodeStubInterface, tradeId string) error {
	var scoreTemp ScoreTemp

	byteData, err := GetScoreTempWithTradeId(stub, tradeId)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteData, &scoreTemp)
	if err != nil {
		err = errors.New("Failed to json decoding.")
		return err
	}

	// 임시 평가점수 삭제 (평가 종료 후)
	err = stub.DelState(scoreTemp.ScoreKey)
	if err != nil {
		return err
	}

	return nil
}

func SetScoreTempExpiryWithTradeId(stub shim.ChaincodeStubInterface, tradeId string, duration time.Duration) error {
	var scoreTemp ScoreTemp
	byteData, err := GetScoreTempWithTradeId(stub, tradeId)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteData, &scoreTemp)
	if err != nil {
		return err
	}

	scoreTemp.ExpiryDate = time.Now().Add(duration) // 지금으로부터 + 평가기간 limit

	inputData, err := json.Marshal(scoreTemp)
	if err != nil {
		return err
	}

	err = stub.PutState(scoreTemp.ScoreKey, inputData)
	if err != nil {
		return err
	}

	return nil
}


func GetScoreTempWithTradeId(stub shim.ChaincodeStubInterface, tradeId string) ([]byte, error) {
	var byteData []byte

	queryString := "{\"selector\":{\"TradeId\":\""+tradeId+"\",\"RecType\":3},\"use_index\":[\"_design/indexTradeDoc\",\"indexTrade\"]}"

	resultsIterators, err := stub.GetQueryResult(queryString)
	if err != nil {
		err = errors.Errorf("Failed to get Trade : query string is wrong : \"%s\"", queryString)
		return nil, err
	}

	if resultsIterators.HasNext() {
		response, err := resultsIterators.Next()
		if err != nil {
			return nil, err
		}
		byteData = response.Value
	}
	if resultsIterators.HasNext() {
		err := errors.New("Temp score must matched only 1 record.")
		return nil, err
	}
	defer resultsIterators.Close()

	return byteData, nil
}


// change sync state
//func UpdateSyncState(stub shim.ChaincodeStubInterface, tradeId string) error {
//	var scoreTemp ScoreTemp
//
//	byteData, err := GetScoreTempWithTradeId(stub, tradeId)
//	if err != nil {
//		return err
//	}
//
//	err = json.Unmarshal(byteData, &scoreTemp)
//	if err != nil {
//		return err
//	}
//
//	scoreTemp.Sync = true
//
//	inputData, err := json.Marshal(byteData)
//	if err != nil {
//		return err
//	}
//
//	err = stub.PutState(scoreTemp.ScoreKey, inputData)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}


// get function for batch
func GetScoreTempForBatch(stub shim.ChaincodeStubInterface) ([]byte, error) {
	buffer := bytes.Buffer{}
	queryString := "{\"selector\":{\"RecType\":3}}"

	resultsIterators, err := stub.GetQueryResult(queryString)
	if err != nil {
		err = errors.Errorf("Failed to get Trade : query string is wrong : \"%s\"", queryString)
		return nil, err
	}

	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	if resultsIterators.HasNext() {
		response, err := resultsIterators.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		byteData := response.Value
		var scoreTemp ScoreTemp
		err = json.Unmarshal(byteData, &scoreTemp)
		if err != nil {
			return nil, err
		}

		// 만료시간이 지난경우 isExpired true로 변경
		if scoreTemp.ExpiryDate.After(time.Now()) {
			scoreTemp.IsExpired = true
			inputData, err := json.Marshal(byteData)
			if err != nil {
				return nil, err
			}
			stub.PutState(scoreTemp.ScoreKey, inputData)
		}

		buffer.Write(response.Value)
		bArrayMemberAlreadyWritten = true
	}
	defer resultsIterators.Close()
	buffer.WriteString("]")

	return buffer.Bytes(), nil
}
package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
	"time"
)

type Properties struct {
	EvaluationLimit time.Duration `json:"EvaluationLimit"` // 평가 입력 기다려주는 시간 (default 14일, 1,209,600 = 14 * 24 * 60 * 60) 이시간 이후에는 0점 처리
	OpenScoreDuration time.Duration `json:"OpenScoreDuration"` // 거래 당사자들의 모든 평가 입력 후 공개하기 까지 걸리는 시간 (default 5일, 432,000 = 5 * 24 * 60 * 60)
}

const propertyKey string = "PROPERTIES"
const defaultEvaluationLimit string = "120s" // 2분, 초단위 (시연)
const defaultOpenScoreDuration string = "30s" // 10초, 초단위 (시연)
//const defaultEvaluationLimit string = "1209600s" // 14일, 초단위 (테스트)
//const defaultOpenScoreDuration string = "432000s" // 5일, 초단위 (테스트)

func GetProperties (stub shim.ChaincodeStubInterface) (Properties, error) {
	var prpty Properties

	byteData, err := stub.GetState(propertyKey)
	if err != nil {
		return Properties{}, err
	}
	if byteData == nil {
		return Properties{}, nil
	}

	err = json.Unmarshal(byteData, &prpty)
	if err != nil {
		return Properties{}, err
	}

	return prpty, nil
}

func SetProperties(stub shim.ChaincodeStubInterface, evaluationLimit string, openScoreDuration string) error {
	prpty, err := GetProperties(stub)
	if err != nil {
		return err
	}

	if !strings.Contains(evaluationLimit,"s") {
		evaluationLimit += "s"
	}
	if !strings.Contains(openScoreDuration,"s") {
		openScoreDuration += "s"
	}

	if evaluationLimit != "" {
		evaluationLimitNew, err := time.ParseDuration(evaluationLimit)
		if err != nil {
			return err
		}
		prpty.EvaluationLimit = evaluationLimitNew
	}
	if openScoreDuration != "" {
		openScoreDurationNew, err := time.ParseDuration(openScoreDuration)
		if err != nil {
			return err
		}
		prpty.OpenScoreDuration = openScoreDurationNew
	}

	inputData, err := json.Marshal(prpty)
	if err != nil {
		return err
	}

	err = stub.PutState(propertyKey, inputData)
	if err != nil {
		return err
	}

	return nil
}
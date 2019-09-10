package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"time"
)

type Properties struct {
	EvaluationLimit time.Duration `json:"EvaluationLimit"` // 평가 입력 기다려주는 시간 (default 14일, 1,209,600 = 14 * 24 * 60 * 60) 이시간 이후에는 0점 처리
	OpenScoreDuration time.Duration `json:"OpenScoreDuration"` // 거래 당사자들의 모든 평가 입력 후 공개하기 까지 걸리는 시간 (default 5일, 432,000 = 5 * 24 * 60 * 60)
}

const propertyKey string = "PROPERTIES"
const default_evaluationLimit time.Duration = 1209600 // 14일, 초단위
const default_openScoreDuration time.Duration = 432000 // 5일, 초단위

func InitProperties(stub shim.ChaincodeStubInterface) error {
	var prpty Properties

	prpty.EvaluationLimit = default_evaluationLimit
	prpty.OpenScoreDuration = default_openScoreDuration

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

func GetProperties (stub shim.ChaincodeStubInterface) (Properties, error) {
	var prpty Properties

	byteData, err := stub.GetState(propertyKey)
	if err != nil {
		return Properties{}, err
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

	if evaluationLimit != "" {
		evaluationLimit_new, err := time.ParseDuration(evaluationLimit)
		if err != nil {
			return err
		}
		prpty.EvaluationLimit = evaluationLimit_new
	}
	if openScoreDuration != "" {
		openScoreDuration_new, err := time.ParseDuration(openScoreDuration)
		if err != nil {
			return err
		}
		prpty.OpenScoreDuration = openScoreDuration_new
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
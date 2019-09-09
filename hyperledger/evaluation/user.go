package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"time"
)

// User data-set
type User struct {
	RecType RecordType `json:"RecType"` // User : 1
	UserTkn string `json:"UserTkn"`
	SellAmt int `json:"SellAmt"`
	BuyAmt int `json:"BuyAmt"`
	SellEx int `json:"SellEx"`
	BuyEx int `json:"BuyEx"`
	Date time.Time `json:"Date"`
	SellSum struct {
		TotSum    int `json:"TotSum"`
		EvalSum01 int `json:"EvalSum01"`
		EvalSum02 int `json:"EvalSum02"`
		EvalSum03 int `json:"EvalSum03"`
	} `json:"SellSum"`
	BuySum struct {
		TotSum    int `json:"TotSum"`
		EvalSum01 int `json:"EvalSum01"`
		EvalSum02 int `json:"EvalSum02"`
		EvalSum03 int `json:"EvalSum03"`
	} `json:"BuySum"`
	TradeSum struct {
		TotSum    int `json:"TotSum"`
		EvalSum01 int `json:"EvalSum01"`
		EvalSum02 int `json:"EvalSum02"`
		EvalSum03 int `json:"EvalSum03"`
	} `json:"TradeSum"`
}


// 사용자 추가
func AddUser(stub shim.ChaincodeStubInterface, userTkn string) error {
	var user User

	byteData, err := stub.GetState(userTkn)
	if err != nil {
		err = errors.Errorf("Failed to get User : UserTkn is \"%s\"", userTkn)
		return nil
	}
	if byteData != nil {
		err := errors.Errorf("User \"%s\" is already exist.", userTkn)
		return err
	}

	user.RecType = RecordTypeUser
	user.UserTkn = userTkn
	user.Date = time.Now()

	inputData, err := json.Marshal(user)
	if err != nil {
		err = errors.New("Failed to json encoding.")
		return err
	}

	err = stub.PutState(userTkn, inputData)
	if err != nil {
		err := errors.New("Failed to store data.")
		return err
	}

	return nil
}


// 사용자 점수 업데이트. division : "sell", "buy". sell인 경우는 판매자의 점수이고(구매자가 매긴 점수), buy인 경우는 구매자의 점수이다.(판매자가 매긴 점수)
func UpdateUserScore(stub shim.ChaincodeStubInterface, userTkn string, score []int, division string) error {
	var user User

	byteData, err := stub.GetState(userTkn)
	if err != nil {
		err = errors.Errorf("Failed to get User : UserTkn is \"%s\"", userTkn)
		return nil
	}
	if byteData == nil {
		err := errors.Errorf("There is no user : UserTkn \"%s\"", userTkn)
		return err
	}

	err = json.Unmarshal(byteData, &user)
	if err != nil {
		err := errors.New("Failed to json decoding.")
		return err
	}

	// 점수 합산
	switch division {
	case "sell": user.SellAmt++
		user.SellSum.EvalSum01 += score[0]
		user.SellSum.EvalSum02 += score[1]
		user.SellSum.EvalSum03 += score[2]

		user.TradeSum.EvalSum01 += score[0]
		user.TradeSum.EvalSum02 += score[1]
		user.TradeSum.EvalSum03 += score[2]
		for _, item := range score {
			user.SellSum.TotSum += item
			user.TradeSum.TotSum += item
		}

	case "buy": user.BuyAmt++
		user.BuySum.EvalSum01 += score[0]
		user.BuySum.EvalSum02 += score[1]
		user.BuySum.EvalSum03 += score[2]

		user.TradeSum.EvalSum01 += score[0]
		user.TradeSum.EvalSum02 += score[1]
		user.TradeSum.EvalSum03 += score[2]
		for _, item := range score {
			user.BuySum.TotSum += item
			user.TradeSum.TotSum += item
		}
	default:
		err := errors.New("Division is wrong. Available value is \"sell\" and \"buy\"")
		return err
	}

	inputData, err := json.Marshal(user)
	if err != nil {
		err = errors.New("Failed to json encoding.")
		return err
	}

	err = stub.PutState(userTkn, inputData)
	if err != nil {
		err := errors.New("Failed to store data.")
		return err
	}

	return nil
}


// 사용자 조회
func GetUser(stub shim.ChaincodeStubInterface, userTkn string) ([]byte, error){

	byteData, err := stub.GetState(userTkn)
	if err != nil {
		err = errors.Errorf("Failed to get User : UserTkn is \"%s\"", userTkn)
		return nil, err
	}
	if byteData == nil {
		err := errors.Errorf("There is no user : UserTkn \"%s\"", userTkn)
		return nil, err
	}

	return byteData, nil
}
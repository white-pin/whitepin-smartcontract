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

// UserTkn = "TOTAL_USER" 인 Record는 전체 사용자의 총 계를 나타낸다.
const TotalUser string = "TOTAL_USER"


// 거래 생성시 사용자 데이터 Init (SellAmt/BuyAmt +1)
func InitUserScore(stub shim.ChaincodeStubInterface, userTkn string, division string) error {
	var user User
	byteData, err := getDataWithKey(stub, userTkn)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteData, &user)
	if err != nil {
		return errors.New("Failed to json decoding.")
	}

	switch division {
	case "sell":
		user.SellAmt++
		user.SellEx++
	case "buy":
		user.BuyAmt++
		user.BuyEx++
	case "tot":
		user.SellAmt++
		user.SellEx++
		user.BuyAmt++
		user.BuyEx++
	default:
		return errors.New("Not allowed division detected.")
	}

	inputData, err := json.Marshal(user)
	if err != nil {
		return errors.New("Failed to json encoding.")
	}

	err = stub.PutState(userTkn, inputData)
	if err != nil {
		return err
	}

	return nil
}


// 사용자 추가
func AddUser(stub shim.ChaincodeStubInterface, userTkn string) error {
	var user User

	byteData, err := stub.GetState(userTkn)
	if err != nil {
		err = errors.Errorf("Failed to get User : UserTkn is \"%s\"", userTkn)
		return nil
	}
	if byteData != nil && userTkn != TotalUser {
		err := errors.Errorf("User \"%s\" is already exist.", userTkn)
		return err
	} else if byteData != nil && userTkn == TotalUser {
		return nil
	}

	user.RecType = RecordTypeUser
	user.UserTkn = userTkn
	user.Date = time.Now()

	inputData, err := json.Marshal(user)
	if err != nil {
		return errors.New("Failed to json encoding.")
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
	case "sell":
		//user.SellAmt++
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

		// 점수가 등록되지 않아서 0 점처리 된 경우
		if !(score[0] == 0 && score[1] == 0 && score[2] == 0) {
			user.SellEx--
		}

	case "buy":
		//user.BuyAmt++
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

		// 점수가 등록되지 않아서 0 점처리 된 경우
		if !(score[0] == 0 && score[1] == 0 && score[2] == 0) {
			user.BuyEx--
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


func UpdateTotalScore(stub shim.ChaincodeStubInterface, sellScore []int, buyScore []int) error {
	var tot User

	byteData, err := stub.GetState(TotalUser)
	if err != nil {
		err = errors.New("Failed to get Total")
		return nil
	}
	if byteData == nil {
		err := errors.New("There is no Total")
		return err
	}

	err = json.Unmarshal(byteData, &tot)
	if err != nil {
		err := errors.New("Failed to json decoding.")
		return err
	}

	//tot.SellAmt++
	tot.SellSum.EvalSum01 += sellScore[0]
	tot.SellSum.EvalSum02 += sellScore[1]
	tot.SellSum.EvalSum03 += sellScore[2]

	tot.TradeSum.EvalSum01 += sellScore[0]
	tot.TradeSum.EvalSum02 += sellScore[1]
	tot.TradeSum.EvalSum03 += sellScore[2]
	for _, item := range sellScore {
		tot.SellSum.TotSum += item
		tot.TradeSum.TotSum += item
	}
	// 점수가 등록되지 않아서 0 점처리 된 경우
	if !(sellScore[0] == 0 && sellScore[1] == 0 && sellScore[2] == 0) {
		tot.SellEx--
	}

	//tot.BuyAmt++
	tot.BuySum.EvalSum01 += buyScore[0]
	tot.BuySum.EvalSum02 += buyScore[1]
	tot.BuySum.EvalSum03 += buyScore[2]

	tot.TradeSum.EvalSum01 += buyScore[0]
	tot.TradeSum.EvalSum02 += buyScore[1]
	tot.TradeSum.EvalSum03 += buyScore[2]
	for _, item := range buyScore {
		tot.BuySum.TotSum += item
		tot.TradeSum.TotSum += item
	}
	// 점수가 등록되지 않아서 0 점처리 된 경우
	if !(buyScore[0] == 0 && buyScore[1] == 0 && buyScore[2] == 0) {
		tot.BuyEx--
	}

	inputData, err := json.Marshal(tot)
	if err != nil {
		err = errors.New("Failed to json encoding.")
		return err
	}

	err = stub.PutState(TotalUser, inputData)
	if err != nil {
		err := errors.New("Failed to store data.")
		return err
	}

	return nil
}


// 사용자 조회
//func GetUser(stub shim.ChaincodeStubInterface, userTkn string) ([]byte, error){
//
//	byteData, err := stub.GetState(userTkn)
//	if err != nil {
//		err = errors.Errorf("Failed to get User : UserTkn is \"%s\"", userTkn)
//		return nil, err
//	}
//	if byteData == nil {
//		err := errors.Errorf("There is no user : UserTkn \"%s\"", userTkn)
//		return nil, err
//	}
//
//	return byteData, nil
//}
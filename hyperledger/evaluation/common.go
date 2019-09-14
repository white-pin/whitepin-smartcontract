package main

import (
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"log"
)

func getQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	buffer := bytes.Buffer{}
	bArrayMemberAlreadyWritten := false

	log.Println("--- Getting query result INIT ---")

	resultIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}

	buffer.WriteString("[")
	if resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.Write(queryResponse.Value)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	log.Println("--- Getting query result FINISHED ---")

	return buffer.Bytes(), nil
}

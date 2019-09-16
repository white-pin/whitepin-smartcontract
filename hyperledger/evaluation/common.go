package main

import (
	"bytes"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func getDataWithKey(stub shim.ChaincodeStubInterface, key string) ([]byte, error) {
	log.Println("--- Getting data with key INIT ---")

	byteData, err := stub.GetState(key)
	if err != nil {
		err = errors.Errorf("Failed to get data with key : Key \"%s\"", key)
		return nil, err
	}
	if byteData == nil {
		err = errors.Errorf("There is no data with key : Key \"%s\"", key)
		return nil, err
	}

	log.Println("--- Getting data with key FINISHED ---")

	return byteData, nil
}


func getQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	log.Println("--- Getting query result INIT ---")

	resultIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultIterator)

	log.Println("--- Getting query result FINISHED ---")

	return buffer.Bytes(), nil
}


func getPageDataWithQueryString(stub shim.ChaincodeStubInterface, queryString string, inputPageSize string, inputPageNum string, inputBookmark string) ([]byte, error) {
	var resultsIterator shim.StateQueryIteratorInterface
	var buffer *bytes.Buffer
	var bookmark string
	//return type of ParseInt is int64
	pageSize, err := strconv.ParseInt(inputPageSize, 10, 32)
	if err != nil {
		return nil, err
	}
	pageNum, err := strconv.Atoi(inputPageNum)
	if err != nil {
		return nil, err
	}

	if inputBookmark != "" {
		bookmark = inputBookmark
		resultsIterator, _, err = stub.GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
		if err != nil {
			return nil, err
		}
		defer resultsIterator.Close()
	} else {
		if pageNum > 1 {
			for i := 0; i < pageNum-1 ; i++ {
				_, responseMetadata, err := stub.GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
				if err != nil {
					return nil, err
				}
				bookmark = responseMetadata.Bookmark
			}
		}
		resultsIterator, _, err = stub.GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
		if err != nil {
			return nil, err
		}
		defer resultsIterator.Close()
	}
	buffer, err = constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}


func constructQueryResponseFromIterator(resultIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error){
	buffer := bytes.Buffer{}
	bArrayMemberAlreadyWritten := false

	buffer.WriteString("[")
	for resultIterator.HasNext() {
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

	return &buffer, nil
}

package main

import (
	"bytes"
	"log"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

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

	return &buffer, nil
}

//func addPageMetaData(buffer *bytes.Buffer, responseMetadata *pb.QueryResponseMetadata) *bytes.Buffer {
//
//	buffer.WriteString("[{\"ResponseMetadata\":{\"RecordsCount\":")
//	buffer.WriteString("\"")
//	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
//	buffer.WriteString("\"")
//	buffer.WriteString(", \"Bookmark\":")
//	buffer.WriteString("\"")
//	buffer.WriteString(responseMetadata.Bookmark)
//	buffer.WriteString("\"}}]")
//
//	return buffer
//}

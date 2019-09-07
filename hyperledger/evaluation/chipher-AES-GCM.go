package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
)

type AES_GCM struct {
	key string // hex encoding 해서 사용. 32자리 : AES128, 64자리 AES256. (encoding 후 각각 16byte, 32byte)
	nonce string // hex encoding 해서 사용. 24자리. (encoding 후 12 byte)
	plainTxt string // 원문
	chipherTxt string // 암호문
}


func GCM_Key(keyString string) (AES_GCM, error) {
	var aes_gcm AES_GCM
	hashKey := sha256.Sum256([]byte(keyString))
	var keyBuffer []byte
	for _, item := range hashKey {
		keyBuffer = append(keyBuffer, item)
	}
	key := hex.EncodeToString(keyBuffer)
	//fmt.Printf("%s\n%x %d\n", key, hashKey, len(key))

	aes_gcm.key = key
	return aes_gcm, nil
}


func (aes_gcm *AES_GCM) GCM_encrypt() error {
	key, _ := hex.DecodeString(aes_gcm.key)
	plaintext := []byte(aes_gcm.plainTxt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	nonce, _ := hex.DecodeString(aes_gcm.nonce)

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	//fmt.Printf("%x\n", ciphertext)

	aes_gcm.chipherTxt = hex.EncodeToString(ciphertext)
	return nil
}


func (aes_gcm *AES_GCM) GCM_decrypt() error {
	key, _ := hex.DecodeString(aes_gcm.key)
	ciphertext, _ := hex.DecodeString(aes_gcm.chipherTxt)
	nonce, _ := hex.DecodeString(aes_gcm.nonce)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	//fmt.Printf("%s\n", plaintext)
	aes_gcm.plainTxt = string(plaintext)

	return nil
}
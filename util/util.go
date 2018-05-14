package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"elastic-search/global"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

var lock sync.Mutex

// WriteResults - the method to write results to file
func WriteResults(values []global.Result, file string) {
	lock.Lock()
	defer lock.Unlock()

	// f, err := os.Create(fmt.Sprintf("/home/eugen/go/src/elastic-search/config/results-%v.log", time.Now().Location()))
	f, err := os.Create(file)

	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, result := range values {
		_, err = f.WriteString(fmt.Sprintf("%s\n", result.String()))
		if err != nil {
			panic(err)
		}
	}
}

// DecryptString decript
func DecryptString(cryptoText string, keyString string) (plainTextString string, err error) {

	// Format the keyString so that it's 32 bytes.
	newKeyString, err := hashTo32Bytes(keyString)
	// Encode the cryptoText to base 64.
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher([]byte(newKeyString))

	if err != nil {
		panic(err)
	}

	if len(cipherText) < aes.BlockSize {
		panic("cipherText too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

// EncryptString is the en
// Takes two string, plainText and keyString.
// plainText is the text that needs to be encrypted by keyString.
// The function will output the resulting crypto text and an error variable.
func EncryptString(plainText string, keyString string) (cipherTextString string, err error) {

	// Format the keyString so that it's 32 bytes.
	newKeyString, err := hashTo32Bytes(keyString)

	if err != nil {
		return "", err
	}

	key := []byte(newKeyString)
	value := []byte(plainText)

	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	cipherText := make([]byte, aes.BlockSize+len(value))

	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], value)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

// As we cannot use a variable length key, we must cut the users key
// up to or down to 32 bytes. To do this the function takes a hash
// of the key and cuts it down to 32 bytes.
func hashTo32Bytes(input string) (output string, err error) {

	if len(input) == 0 {
		return "", errors.New("No input supplied")
	}

	hasher := sha256.New()
	hasher.Write([]byte(input))

	stringToSHA256 := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	// Cut the length down to 32 bytes and return.
	return stringToSHA256[:32], nil
}

func checkConnection(remoteServer string) bool {

	return false
}

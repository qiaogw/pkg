//aes加密解密
package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

//AesEncrypt aes加密
func AesEncrypt(plantText []byte, key []byte, iv []byte) (string, error) {

	block, err := aes.NewCipher(key)

	//选择加密算法

	if err != nil {

		return "", err

	}

	plantText = PKCS7Padding(plantText, block.BlockSize())

	if len(plantText)%aes.BlockSize != 0 { //块大小在aes.BlockSize中定义

		panic("plantText is not a multiple of the block size")

	}

	blockModel := cipher.NewCBCEncrypter(block, iv)

	ciphertext := make([]byte, len(plantText))

	blockModel.CryptBlocks(ciphertext, plantText)

	return base64.StdEncoding.EncodeToString(ciphertext), nil

}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {

	padding := blockSize - len(ciphertext)%blockSize

	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)

}

//AesDecrypt aes解密
func AesDecrypt(cs string, key []byte, iv []byte) ([]byte, error) {
	ciphertext, _ := base64.StdEncoding.DecodeString(cs)

	block, err := aes.NewCipher(key)

	//选择加密算法

	if err != nil {

		return nil, err

	}

	blockModel := cipher.NewCBCDecrypter(block, iv)

	plantText := make([]byte, len(ciphertext))

	blockModel.CryptBlocks(plantText, ciphertext)

	plantText = PKCS7UnPadding(plantText, block.BlockSize())

	return plantText, nil

}

func PKCS7UnPadding(plantText []byte, blockSize int) []byte {

	length := len(plantText)

	unpadding := int(plantText[length-1])

	return plantText[:(length - unpadding)]

}

package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// Encrypt 使用AES加密算法对输入数据进行加密
func Encrypt(plainText []byte, key []byte, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 对明文进行PKCS7填充
	plainText = PKCS7Padding(plainText, block.BlockSize())

	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)

	// 返回Base64编码后的密文
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// PKCS7Padding 对数据进行PKCS7填充
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// Decrypt 使用AES解密算法对输入数据进行解密
func Decrypt(cipherText string, key []byte, iv []byte) ([]byte, error) {
	cipherData, _ := base64.StdEncoding.DecodeString(cipherText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherData))
	blockMode.CryptBlocks(plainText, cipherData)

	// 去除PKCS7填充
	plainText = PKCS7UnPadding(plainText, block.BlockSize())

	return plainText, nil
}

// PKCS7UnPadding 去除PKCS7填充
func PKCS7UnPadding(data []byte, blockSize int) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

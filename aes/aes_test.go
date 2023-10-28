package aes

import (
	"bytes"
	"testing"
)

func TestAesEncryptionAndDecryption(t *testing.T) {
	key := []byte("0123456789012345")
	iv := []byte("1234567890123456")
	plainText := []byte("Hello, AES!")

	// 加密
	encrypted, err := Encrypt(plainText, key, iv)
	if err != nil {
		t.Errorf("Error encrypting: %v", err)
	}

	// 解密
	decrypted, err := Decrypt(encrypted, key, iv)
	if err != nil {
		t.Errorf("Error decrypting: %v", err)
	}

	// 比较解密后的明文与原始明文
	if !bytes.Equal(decrypted, plainText) {
		t.Errorf("Decrypted text does not match original plain text.")
	}
}

func TestPKCS7PaddingAndUnPadding(t *testing.T) {
	blockSize := 16
	data := []byte("Test Data")

	// 进行PKCS7填充
	paddedData := PKCS7Padding(data, blockSize)
	if len(paddedData)%blockSize != 0 {
		t.Errorf("PKCS7 padding error")
	}

	// 进行PKCS7去填充
	unpaddedData := PKCS7UnPadding(paddedData, blockSize)
	if !bytes.Equal(unpaddedData, data) {
		t.Errorf("PKCS7 unpadding error")
	}
}

package genrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/qiaogw/pkg/tools"

	"os"
)

// GenKey 密钥文件生成
func GenKey(bits int,jwtPublicPath,jwtPrivatePath string) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "私钥",
		Bytes: derStream,
	}
	err = tools.CheckPath(jwtPrivatePath)
	file, err := os.Create(jwtPrivatePath)
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "公钥",
		Bytes: derPkix,
	}
	err = tools.CheckPath(jwtPublicPath)
	file, err = os.Create(jwtPublicPath)
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

// RsaGenKey 生成RSA公钥和私钥并保存在对应的目录文件下
// 参数bits: 指定生成的秘钥的长度, 单位: bit
func RsaGenKey(bits int, privateFileName, publicFileName string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	x509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	privateFile, err := os.Create(privateFileName)
	if err != nil {
		return err
	}
	defer privateFile.Close()
	privateBlock := pem.Block{
		Type:  "私钥",
		Bytes: x509PrivateKey,
	}

	if err = pem.Encode(privateFile, &privateBlock); err != nil {
		return err
	}
	publicKey := privateKey.PublicKey
	x509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	publicFile, _ := os.Create(publicFileName)
	defer publicFile.Close()
	publicBlock := pem.Block{
		Type:  "公钥",
		Bytes: x509PublicKey,
	}
	if err = pem.Encode(publicFile, &publicBlock); err != nil {
		return err
	}
	return nil
}

// RSAEncrypt RSA公钥加密
func RSAEncrypt(src []byte, filename string) ([]byte, error) {
	// 根据文件名读出文件内容
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)

	// 从数据中找出pem格式的块
	block, _ := pem.Decode(buf)
	if block == nil {
		return nil, err
	}

	// 解析一个der编码的公钥//对应于生成秘钥的x509.MarshalPKIXPublicKey(&publicKey)
	//publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 公钥加密
	publicKey := pubKey.(*rsa.PublicKey)
	result, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, src)
	return result, err
}

//RSADecryptRSA私钥解密
func RSADecrypt(src []byte, filename string) ([]byte, error) {
	// 根据文件名读出内容
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)

	// 从数据中解析出pem块
	block, _ := pem.Decode(buf)
	if block == nil {
		return nil, err
	}

	// 解析出一个der编码的私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 私钥解密
	result, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, src)
	if err != nil {
		return nil, err
	}
	return result, nil
}

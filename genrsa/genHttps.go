//Package genrsa 生成https私有证书
package genrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"

	//	"encoding/asn1"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	rd "math/rand"
	"os"
	"time"
	//	"github.com/astaxie/beego"
)

func init() {
	rd.Seed(time.Now().UnixNano())
}

// CertInformation 证书信息
type CertInformation struct {
	Country            []string
	Organization       []string
	OrganizationalUnit []string
	EmailAddress       []string
	Province           []string
	Locality           []string
	CommonName         string
	CrtName, KeyName   string
	IsCA               bool
	Names              []pkix.AttributeTypeAndValue
}

// CreateCRT 创建证书
func CreateCRT(RootCa *x509.Certificate, RootKey *rsa.PrivateKey, info CertInformation) error {
	Crt := newCertificate(info)
	Key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	var buf []byte
	if RootCa == nil || RootKey == nil {
		//创建自签名证书
		buf, err = x509.CreateCertificate(rand.Reader, Crt, Crt, &Key.PublicKey, Key)
	} else {
		//使用根证书签名
		buf, err = x509.CreateCertificate(rand.Reader, Crt, RootCa, &Key.PublicKey, RootKey)
	}
	if err != nil {
		return err
	}

	err = write(info.CrtName, "CERTIFICATE", buf)
	if err != nil {
		return err
	}

	buf = x509.MarshalPKCS1PrivateKey(Key)
	return write(info.KeyName, "PRIVATE KEY", buf)
}

//编码写入文件
func write(filename, Type string, p []byte) error {
	File, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer File.Close()
	var b *pem.Block = &pem.Block{Bytes: p, Type: Type}
	return pem.Encode(File, b)

}

// Parse 创建所有证书
func Parse(crtPath, keyPath string) (rootcertificate *x509.Certificate, rootPrivateKey *rsa.PrivateKey, err error) {
	rootcertificate, err = ParseCrt(crtPath)
	if err != nil {
		return
	}
	rootPrivateKey, err = ParseKey(keyPath)
	return
}

//ParseCrt 创建证书
func ParseCrt(path string) (*x509.Certificate, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// p := &pem.Block{}
	p, _ := pem.Decode(buf)
	return x509.ParseCertificate(p.Bytes)
}

//ParseKey 创建key
func ParseKey(path string) (*rsa.PrivateKey, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p, _ := pem.Decode(buf)
	return x509.ParsePKCS1PrivateKey(p.Bytes)
}

func newCertificate(info CertInformation) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(rd.Int63()),
		Subject: pkix.Name{
			Country:            info.Country,
			Organization:       info.Organization,
			OrganizationalUnit: info.OrganizationalUnit,
			Province:           info.Province,
			CommonName:         info.CommonName,
			Locality:           info.Locality,
			ExtraNames:         info.Names,
		},
		NotBefore:             time.Now(),                                                                 //证书的开始时间
		NotAfter:              time.Now().AddDate(20, 0, 0),                                               //证书的结束时间
		BasicConstraintsValid: true,                                                                       //基本的有效性约束
		IsCA:                  info.IsCA,                                                                  //是否是根证书
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, //证书用途
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		EmailAddresses:        info.EmailAddress,
	}
}

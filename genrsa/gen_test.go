package genrsa

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestCrt(t *testing.T) {
	baseinfo := CertInformation{Country: []string{"CN"}, Organization: []string{"WS"}, IsCA: true,
		OrganizationalUnit: []string{"goadmin"}, EmailAddress: []string{"arzang@163.com"},
		Locality: []string{"LanZhou"}, Province: []string{"GansuSu"}, CommonName: "ilive-Stacks",
		CrtName: "test_root.crt", KeyName: "test_root.key"}

	err := CreateCRT(nil, nil, baseinfo)
	if err != nil {
		t.Log("Create crt error,Error info:", err)
		return
	}
	crtinfo := baseinfo
	crtinfo.IsCA = false
	crtinfo.CrtName = "test_server.crt"
	crtinfo.KeyName = "test_server.key"
	//添加扩展字段用来做自定义使用

	crtinfo.Names = []pkix.AttributeTypeAndValue{
		{Type: asn1.ObjectIdentifier{2, 1, 3}, Value: "MAC_ADDR"},
	}

	crt, pri, err := Parse(baseinfo.CrtName, baseinfo.KeyName)
	if err != nil {
		t.Log("Parse crt error,Error info:", err)
		return
	}
	err = CreateCRT(crt, pri, crtinfo)
	if err != nil {
		t.Log("Create crt error,Error info:", err)
	}
	os.Remove(baseinfo.CrtName)
	os.Remove(baseinfo.KeyName)
	os.Remove(crtinfo.CrtName)
	os.Remove(crtinfo.KeyName)
}

func TestRSA(t *testing.T) {
	private := "privateKey.pem"
	public := "pubulicKey.pem"
	RsaGenKey(2048, private, public)
	baseinfo := CertInformation{Country: []string{"CN"}, Organization: []string{"WS"}, IsCA: true,
		OrganizationalUnit: []string{"goadmin"}, EmailAddress: []string{"arzang@163.com"},
		Locality: []string{"LanZhou"}, Province: []string{"GansuSu"}, CommonName: "ilive-Stacks",
		CrtName: "test_root.crt", KeyName: "test_root.key"}
	msg, _ := json.Marshal(baseinfo)
	cipherText, err := RSAEncrypt(msg, public)
	fmt.Println(string(cipherText))
	if err != nil {
		t.Fatal(err)
	}
	plainText, err := RSADecrypt(cipherText, private)
	var v CertInformation
	json.Unmarshal(plainText, &v)
	fmt.Println(v)
	if err != nil {
		t.Fatal(err)
	}
}

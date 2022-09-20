package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

const (
	LicenseToolsPath    = "/Users/qilin.wang/IdeaProjects/goProject/"
	CAPath              = LicenseToolsPath + "ca.pem"
	RSAPriKeyPath       = LicenseToolsPath + "pri.key"
	ReqPath             = LicenseToolsPath + "req.txt"
	ApplicationCodePath = LicenseToolsPath + "application-code.txt"
	GenLicReqPath       = LicenseToolsPath + "genlicreq"
)

type LicenseInfo struct {
	Id            string `json:"licenseid"`
	User          string `json:"user"`
	LicenseType   int    `json:"type"`
	Hostnum       int    `json:"hostnum"`
	Cpunum        int    `json:"Cpunum"`
	IssueTime     string `json:"issuetime"`
	ExpiredTime   string `json:"expiretime"`
	Prodinfo      string `json:"prodinfo"`
	B64thumbprint string `json:"thumbprint"`
	VmNum         int    `json:"vmnum"`
}

type ApplicationCode struct {
	PrivateKey     string `json:"privateKey"`
	LicenseRequest string `json:"licenseRequest"`
}

type licenseBody struct {
	License string `json:"license"`
	Aeskey  string `json:"aeskey"`
}

func generateApplicationCode() error {
	//check files
	err := checkTools([]string{CAPath, GenLicReqPath})
	if err != nil {
		return err
	}
	//chmod + x
	err = exec.Command("sh", "-c", fmt.Sprintf("chmod +x %s", GenLicReqPath)).Run()
	if err != nil {
		return err
	}
	//generate req&pri.key
	err = exec.Command("sh", "-c", fmt.Sprintf("%s -privkey %s -reqfile %s %s", GenLicReqPath, RSAPriKeyPath, ReqPath, CAPath)).Run()
	if err != nil {
		return err
	}
	//generate applicationCode
	reqFile, err := ioutil.ReadFile(ReqPath)
	if err != nil {
		return err
	}
	req := string(reqFile)
	keyFile, err := ioutil.ReadFile(RSAPriKeyPath)
	if err != nil {
		return err
	}
	key := string(keyFile)
	marshal, err := json.Marshal(ApplicationCode{
		PrivateKey:     key,
		LicenseRequest: req,
	})
	err = ioutil.WriteFile(ApplicationCodePath, []byte(base64.StdEncoding.EncodeToString(marshal)), 0644)
	if err != nil {
		return err
	}
	return err
}

func checkTools(toolsPath []string) error {
	var err error
	for _, v := range toolsPath {
		_, err = os.Stat(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func DecryptLicenseTxt(licenseTxt string) (error, LicenseInfo) {
	var license LicenseInfo
	//check ca.pem/pri.key
	err := checkTools([]string{CAPath, RSAPriKeyPath})
	if err != nil {
		return err, license
	}
	//verify
	//ca, err := ioutil.ReadFile(CAPath)
	bytes, err := ioutil.ReadFile(RSAPriKeyPath)
	if err != nil {
		return err, license
	}
	//pem.Decode(ca)
	//

	//get lic json
	data, err := base64.StdEncoding.DecodeString(licenseTxt)
	if err != nil {
		return err, license
	}
	var lic licenseBody
	err = json.Unmarshal(data, &lic)
	if err != nil {
		return err, license
	}
	if lic.Aeskey == "" {
		decodeString, err := base64.StdEncoding.DecodeString(lic.License)
		if err != nil {
			return err, license
		}
		err = json.Unmarshal(decodeString, &license)
		if err != nil {
			return err, license
		}
		return nil, license
	}
	//get rsa key
	realRsaKey, err := RsaDecrypt([]byte(lic.Aeskey), bytes)
	if err != nil {
		return err, license
	}
	//aes decrypt
	decrypt, err := AesDecrypt([]byte(lic.License), realRsaKey)
	if err != nil {
		return err, license
	}
	err = json.Unmarshal(decrypt, &license)
	if err != nil {
		return err, license
	}
	return nil, license
}

func AesDecrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	blockMode := cipher.NewCFBDecrypter(block, key[:blockSize])
	crypted := make([]byte, len(data))
	blockMode.XORKeyStream(crypted, data)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

func RsaDecrypt(cipherText []byte, priKey []byte) ([]byte, error) {
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("private key error")
	}

	parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	//key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("parse error")
	}
	key := parseResult.(*rsa.PrivateKey)
	return rsa.DecryptPKCS1v15(rand.Reader, key, cipherText)
}

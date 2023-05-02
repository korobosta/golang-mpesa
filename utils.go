package mpesa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	randString "math/rand"
	"time"
)

func RandomString(length int) string {
    randString.Seed(time.Now().UnixNano())
    b := make([]byte, length+2)
    rand.Read(b)
    return fmt.Sprintf("%x", b)[2 : length+2]
}

func EncryptWithPublicKey( password string,env int) string {
	//byteMessage := []byte(message)
	var f = ""
	if env == 1{
		f = "ProductionCertificate.cer"
	}else{
		f = "SandboxCertificate.cer"
	}

	publicKey, err := ioutil.ReadFile(f)
	if err != nil {
		log.Println(err)
	}

	block, _ := pem.Decode(publicKey)

	var cert *x509.Certificate
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Println(err)
	}

	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	reader := rand.Reader
	signature, err := rsa.EncryptPKCS1v15(reader, rsaPublicKey, []byte(password))
	if err != nil {
		log.Println(err)
	}


	return base64.StdEncoding.EncodeToString(signature)
	
}
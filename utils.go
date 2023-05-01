package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"log"
)

func EncryptWithPublicKey( msg string,env int, password string) string {
	//byteMessage := []byte(message)
	var f = ""
	if env == 1{
		f = "certs/ProductionCertificate.cer"
	}else{
		f = "certs/SandboxCertificate.cer"
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
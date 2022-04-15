package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// create certs directory if not exists
	err := os.MkdirAll("certs", os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	// generate root private key
	rootPriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalln(err)
	}

	// convert root private key to DER
	rootPrivDER, err := x509.MarshalECPrivateKey(rootPriv)
	if err != nil {
		log.Fatalln(err)
	}

	// convert root private key to PEM
	rootPrivPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: rootPrivDER,
	})

	// save root private key to certs/ca.key
	err = os.WriteFile(filepath.Join("certs", "ca.key"), rootPrivPEM, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	// create root certificate
	rootTmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "root",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	rootCertDER, err := x509.CreateCertificate(rand.Reader, rootTmpl, rootTmpl, rootPriv.Public(), rootPriv)
	if err != nil {
		log.Fatalln(err)
	}
	rootCert, err := x509.ParseCertificate(rootCertDER)
	if err != nil {
		log.Fatalln(err)
	}

	// create root certificate pool for verification
	certPool := x509.NewCertPool()
	certPool.AddCert(rootCert)

	// convert root certificate to PEM
	rootCertPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: rootCertDER,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// save root certificate to certs/ca.crt
	err = os.WriteFile(filepath.Join("certs", "ca.crt"), rootCertPEM, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	// generate server private key
	servPriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalln(err)
	}

	// convert server private key to DER
	servPrivDER, err := x509.MarshalECPrivateKey(servPriv)
	if err != nil {
		log.Fatalln(err)
	}

	// convert server private key to PEM
	servPrivPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: servPrivDER,
	})

	// save server private key to certs/server.key
	err = os.WriteFile(filepath.Join("certs", "server.key"), servPrivPEM, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	// create server certificate
	servTmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "server",
		},
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(10, 0, 0),
		DNSNames:    []string{"127.0.0.1"}, // server hostname
	}
	servCertDER, err := x509.CreateCertificate(rand.Reader, servTmpl, rootCert, servPriv.Public(), rootPriv)
	if err != nil {
		return
	}
	servCert, err := x509.ParseCertificate(servCertDER)
	if err != nil {
		log.Fatalln(err)
	}

	// verify server certificate
	_, err = servCert.Verify(x509.VerifyOptions{
		Roots: certPool,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// convert server certificate to PEM
	servCertPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: servCertDER,
	})

	// save server certificate to certs/server.crt
	err = os.WriteFile(filepath.Join("certs", "server.crt"), servCertPEM, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	// generate client private key
	cliPriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalln(err)
	}

	// convert client private key to DER
	cliPrivDER, err := x509.MarshalECPrivateKey(cliPriv)
	if err != nil {
		log.Fatalln(err)
	}

	// convert client private key to PEM
	cliPrivPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: cliPrivDER,
	})

	// save client private key to certs/client.key
	err = os.WriteFile(filepath.Join("certs", "client.key"), cliPrivPEM, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	// convert client private key to PKCS8
	cliPrivPKCS8, err := x509.MarshalPKCS8PrivateKey(cliPriv)
	if err != nil {
		log.Fatalln(err)
	}

	// save client private key (PKCS8) to certs/client.pk8
	err = os.WriteFile(filepath.Join("certs", "client.pk8"), cliPrivPKCS8, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	// create client certificate
	cliTmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "postgres", // postgresql server user
		},
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(10, 0, 0),
	}
	cliCertDER, err := x509.CreateCertificate(rand.Reader, cliTmpl, rootCert, cliPriv.Public(), rootPriv)
	if err != nil {
		return
	}
	cliCert, err := x509.ParseCertificate(cliCertDER)
	if err != nil {
		log.Fatalln(err)
	}

	// verify client certificate
	_, err = cliCert.Verify(x509.VerifyOptions{
		Roots: certPool,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// convert client certificate to PEM
	cliCertPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cliCertDER,
	})

	// save client certificate to certs/client.crt
	err = os.WriteFile(filepath.Join("certs", "client.crt"), cliCertPEM, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
}

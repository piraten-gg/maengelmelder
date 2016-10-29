/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	pseudoRandom "math/rand"
	"os"
	"time"
)

// create a TLS certificate for the side if it does not exist
func CreateCertificateIfNotExist(certFile string, keyFile string, commonName string, bits int) {
	if _, err := os.Stat(certFile); err == nil {
		// certFile already exists
		return
	}
	if _, err := os.Stat(keyFile); err == nil {
		// keyFile already exists
		return
	}

	fmt.Println("Generating RSA key and certificate.")
	pseudoRandom.Seed(time.Now().UTC().UnixNano())

	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Println("Generate key:", err)
	}

	notBefore := time.Now()
	notAfter := time.Now().AddDate(10, 0, 0) // the certificate should expire in 10 years

	template := x509.Certificate{
		SerialNumber: big.NewInt(pseudoRandom.Int63()),
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		fmt.Println("create cert:", err)
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		fmt.Println("save cert:", err)
	}
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		fmt.Println("save cert:", err)
	}
	err = certOut.Close()
	if err != nil {
		fmt.Println("save cert:", err)
	}

	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println("save key:", err)
	}
	err = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	if err != nil {
		fmt.Println("save key:", err)
	}
	err = keyOut.Close()
	if err != nil {
		fmt.Println("save key:", err)
	}
}

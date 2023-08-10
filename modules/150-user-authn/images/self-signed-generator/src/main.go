/*
Copyright 2023 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

func main() {
	certHosts := os.Args[1:]

	// Generate a CA private key
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate CA private key: %v", err)
	}

	// Create a self-signed CA certificate
	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "CA"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // Valid for 10 years
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		log.Fatalf("Failed to create CA certificate: %v", err)
	}

	// Save the CA private key to a file
	caPrivateKeyFile, err := os.Create("/certs/ca.key")
	if err != nil {
		log.Fatalf("Failed to create CA private key file: %v", err)
	}
	err = pem.Encode(caPrivateKeyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivateKey),
	})
	if err != nil {
		log.Fatal(err)
	}
	caPrivateKeyFile.Close()

	// Save the CA certificate to a file
	caCertificateFile, err := os.Create("/certs/ca.crt")
	if err != nil {
		log.Fatalf("Failed to create CA certificate file: %v", err)
	}
	err = pem.Encode(caCertificateFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		log.Fatal(err)
	}
	caCertificateFile.Close()

	// Generate a private key for the server certificate
	serverPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate server private key: %v", err)
	}

	// Sign the server certificate with the CA certificate and private key
	serverCertificate, err := x509.CreateCertificate(rand.Reader, &x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{CommonName: "Self-signed"},
		DNSNames:              certHosts,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}, caTemplate, serverPrivateKey.Public(), caPrivateKey)
	if err != nil {
		log.Fatalf("Failed to sign server certificate: %v", err)
	}

	// Save the server private key to a file
	serverPrivateKeyFile, err := os.Create("/certs/tls.key")
	if err != nil {
		log.Fatalf("Failed to create server private key file: %v", err)
	}
	err = pem.Encode(serverPrivateKeyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(serverPrivateKey),
	})
	if err != nil {
		log.Fatal(err)
	}
	serverPrivateKeyFile.Close()

	// Save the server certificate to a file
	serverCertificateFile, err := os.Create("/certs/tls.crt")
	if err != nil {
		log.Fatalf("Failed to create server certificate file: %v", err)
	}
	err = pem.Encode(serverCertificateFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCertificate,
	})
	if err != nil {
		log.Fatal(err)
	}
	serverCertificateFile.Close()

	log.Print("Self-signed certificate, key, and CA certificate, key generated successfully.")
}

package generator

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"log"
	"os"
)

func GenerateRSAKey() {
	privateKey, generateRSAErr := rsa.GenerateKey(rand.Reader, 2048)
	if generateRSAErr != nil {
		log.Fatal("Failed to generate key: " + generateRSAErr.Error())
	}

	envFile, openEnvFileErr := os.OpenFile(".env", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if openEnvFileErr != nil {
		log.Fatal("Failed to locate env file: " + openEnvFileErr.Error())
	}

	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	}

	encodePrivateKey := pem.EncodeToMemory(privateKeyBlock)
	encodePublicKey := pem.EncodeToMemory(publicKeyBlock)

	base64PrivateKeyBlock := base64.StdEncoding.EncodeToString(encodePrivateKey)
	base64PublicKeyBlock := base64.StdEncoding.EncodeToString(encodePublicKey)

	_, writePrivateKeyToEnvErr := envFile.WriteString("\nPRIVATE_KEY=" + base64PrivateKeyBlock + "\n")
	if writePrivateKeyToEnvErr != nil {
		log.Fatal("Failed to write to env file: " + writePrivateKeyToEnvErr.Error())
	}

	_, writePublicKeyToEnvErr := envFile.WriteString("\nPUBLIC_KEY=" + base64PublicKeyBlock + "\n")
	if writePublicKeyToEnvErr != nil {
		log.Fatal("Failed to write to env file: " + writePublicKeyToEnvErr.Error())
	}
}

func GenerateCSRFKey() (string, error) {
	csrfKey := make([]byte, 32)
	_, genCsrfKeyErr := rand.Read(csrfKey)
	if genCsrfKeyErr != nil {
		return "", errors.New("failed to create CSRF key: " + genCsrfKeyErr.Error())
	}

	pureString := hex.EncodeToString(csrfKey)

	return pureString, nil
}

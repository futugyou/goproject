package extensions

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
)

// AesCTREncrypt encrypts plaintext using AES in CTR mode
// key 16, 24 or 32 bytes for AES-128, AES-192, or AES-256
func AesCTREncrypt(plaintext, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AesCTRDecrypt decrypts a base64 encoded ciphertext using AES in CTR mode
// key 16, 24 or 32 bytes for AES-128, AES-192, or AES-256
func AesCTRDecrypt(EncryptText, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(EncryptText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// AesCFBEncrypt encrypts plaintext using AES in CFB mode
// key 16, 24 or 32 bytes for AES-128, AES-192, or AES-256
func AesCFBEncrypt(plaintext, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AesCFBDecrypt decrypts a base64 encoded ciphertext using AES in CFB mode
// key 16, 24 or 32 bytes for AES-128, AES-192, or AES-256
func AesCFBDecrypt(EncryptText, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(EncryptText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// AesGCMEncrypt encrypts plaintext using AES in GCM mode
// key 16, 24 or 32 bytes for AES-128, AES-192, or AES-256
func AesGCMEncrypt(plaintext, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AesGCMDecrypt decrypts a base64 encoded ciphertext using AES in GCM mode
// key 16, 24 or 32 bytes for AES-128, AES-192, or AES-256
func AesGCMDecrypt(encodedCiphertext, key string) (string, error) {
	enc, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", err
	}

	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(enc) < aesGCM.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := enc[:aesGCM.NonceSize()], enc[aesGCM.NonceSize():]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func RsaEncryptMessage(message string, publicKeyStr string) (string, error) {
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return "", errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("not a valid RSA public key")
	}

	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(message), nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func RsaDecrypt(cipherText string, privateKeyStr string) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return "", errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, cipherTextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

// Deprecated: Sha256 is deprecated.
func Sha1(text string) string {
	h := sha1.New()
	h.Write([]byte(text))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Deprecated: Sha256 is deprecated.
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

func Sha256(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

func Sha3(text string) string {
	hash := sha3.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

func Hmac(text string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(text))
	return hex.EncodeToString(mac.Sum(nil))
}

func Blake2b(text string) string {
	hash := blake2b.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

func String(n int64) string {
	buf := [11]byte{}
	pos := len(buf)
	i := n
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

// set n = 32
func GenerateRandomKey(n int) (string, error) {
	key := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func GenerateKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// bits = 2048
func GenerateKeyPair(bits int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil
}

// Generate an ECDSA key pair and convert it to a PEM formatted string
func GenerateECDSAKeyPair() (string, string, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil
}

// Signing a message using an ECDSA private key
func SignMessage(privateKeyStr, message string) (string, string, error) {
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return "", "", fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", "", err
	}

	hash := sha256.Sum256([]byte(message))
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", "", err
	}

	rBytes, err := r.MarshalText()
	if err != nil {
		return "", "", err
	}
	sBytes, err := s.MarshalText()
	if err != nil {
		return "", "", err
	}

	return base64.StdEncoding.EncodeToString(rBytes), base64.StdEncoding.EncodeToString(sBytes), nil
}

// Verify signature using ECDSA public key
func VerifySignature(publicKeyStr, message, rStr, sStr string) (bool, error) {
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return false, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	publicKey, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return false, fmt.Errorf("not a valid ECDSA public key")
	}

	hash := sha256.Sum256([]byte(message))

	rBytes, err := base64.StdEncoding.DecodeString(rStr)
	if err != nil {
		return false, err
	}
	sBytes, err := base64.StdEncoding.DecodeString(sStr)
	if err != nil {
		return false, err
	}

	var r, s big.Int
	if err := r.UnmarshalText(rBytes); err != nil {
		return false, err
	}
	if err := s.UnmarshalText(sBytes); err != nil {
		return false, err
	}

	return ecdsa.Verify(publicKey, hash[:], &r, &s), nil
}

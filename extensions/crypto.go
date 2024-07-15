package extensions

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
)

// AesCTRCrypto encrypts plaintext using AES in CTR mode
// key 16, 24 or 32 bytes for AES-128, AES-192, or AES-256
func AesCTRCrypto(plaintext, key string) (string, error) {
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
func AesCTRDecrypt(cryptoText, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
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

// AesCFBCrypto encrypts plaintext using AES in CFB mode
// key 16, 24 or 32 bytes for AES-128, AES-192, or AES-256
func AesCFBCrypto(plaintext, key string) (string, error) {
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
func AesCFBDecrypt(cryptoText, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
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

// AesGCMCrypto encrypts plaintext using AES in GCM mode
// key 16, 24 or 32 bytes for AES-128, AES-192, or AES-256
func AesGCMCrypto(plaintext, key string) (string, error) {
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

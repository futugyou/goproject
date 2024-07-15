package extensions

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
)

// key 16, 24 or 32 bytes for AES-128, AES-192, or AES-256
func AesCTRCrypto(key string, plaintext string) (string, string, error) {
	keyString := []byte(key)
	plaintextString := []byte(plaintext)

	block, err := aes.NewCipher(keyString)
	if err != nil {
		return "", "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintextString))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", "", err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintextString)

	return hex.EncodeToString(ciphertext), hex.EncodeToString(iv), nil
}

// key 16, 24 or 32 bytes for AES-128, AES-192, or AES-256
func AesCTRDecrypt(key string, ciphertext string, iv string) (string, error) {
	keyString := []byte(key)
	ivString, err := hex.DecodeString(iv)
	if err != nil {
		return "", err
	}
	ciphertextString, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyString)
	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(ciphertextString))
	stream := cipher.NewCTR(block, ivString)
	stream.XORKeyStream(plaintext, ciphertextString)

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

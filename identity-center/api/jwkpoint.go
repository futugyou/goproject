package api

import (
	_ "github.com/joho/godotenv/autoload"

	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

func Jwks(w http.ResponseWriter, r *http.Request) {

	signed_key := os.Getenv("signed_key")

	raw, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("failed to generate new rsa private key: %s\n", err)
		return
	}

	key, err := jwk.FromRaw(raw)
	if err != nil {
		fmt.Printf("failed to create RSA key: %s\n", err)
		return
	}
	if _, ok := key.(jwk.RSAPrivateKey); !ok {
		fmt.Printf("expected jwk.RSAPrivateKey, got %T\n", key)
		return
	}

	key.Set(jwk.KeyIDKey, signed_key)

	buf, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		fmt.Printf("failed to marshal key into JSON: %s\n", err)
		return
	}
	result := "{ \"keys\": [ " + string(buf) + " ] }"
	w.Write([]byte(result))
	w.WriteHeader(200)
}

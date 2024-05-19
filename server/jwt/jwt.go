package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type JWTPayload struct {
	User string `json:"user"`
	Time string `json:"time"`
}

type JWTDecoded struct {
	header  JWTHeader
	payload JWTPayload
}

func EncodeBase64URL(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

func DecodeBase64(token string) (JWTDecoded, error) {
	encodedParts := strings.Split(token, ".")
	if len(encodedParts) != 3 {
		return JWTDecoded{}, fmt.Errorf("invalid JWT token format")
	}

	decodedHeader, err := base64.RawURLEncoding.DecodeString(encodedParts[0])
	if err != nil {
		return JWTDecoded{}, fmt.Errorf("error decoding token header:, %v", err)
	}

	decodedPayload, err := base64.RawURLEncoding.DecodeString(encodedParts[1])
	if err != nil {
		return JWTDecoded{}, fmt.Errorf("error decoding token payload:, %v", err)
	}

	var header JWTHeader
	err = json.Unmarshal(decodedHeader, &header)
	if err != nil {
		return JWTDecoded{}, fmt.Errorf("error unmarshaling header:, %v", err)
	}

	var payload JWTPayload
	err = json.Unmarshal(decodedPayload, &payload)
	if err != nil {
		return JWTDecoded{}, fmt.Errorf("error unmarshaling payload:, %v", err)
	}

	return JWTDecoded{
		header,
		payload,
	}, nil
}

func Construct(user string) (string, error) {
	header := JWTHeader{
		Alg: "HS256",
		Typ: "JWT",
	}

	payload := JWTPayload{
		User: user,
		Time: time.Now().String(),
	}

	secret := os.Getenv("jwtsecret")

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("json marshal error with header: %v", err)
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("json marshal error with payload: %v", err)
	}

	encodedHeader := EncodeBase64URL(headerJSON)
	encodedPayload := EncodeBase64URL(payloadJSON)

	message := fmt.Sprintf("%s.%s", encodedHeader, encodedPayload)

	signature := hmac.New(sha256.New, []byte(secret))
	signature.Write([]byte(message))
	encodedSignature := EncodeBase64URL(signature.Sum(nil))

	token := fmt.Sprintf("%s.%s", message, encodedSignature)

	return token, nil
}

// does this function even make sense i dont know
// i just decode it and encode it again and then check the signature
// hmm
func Verify(token string) (bool, error) {
	secret := os.Getenv("jwtsecret")
	decoded, err := DecodeBase64(token)
	if err != nil {
		return false, fmt.Errorf("token decoding error: %v", err)
	}

	headerJSON, err := json.Marshal(decoded.header)
	if err != nil {
		return false, fmt.Errorf("json marshal error with header: %v", err)
	}

	payloadJSON, err := json.Marshal(decoded.payload)
	if err != nil {
		return false, fmt.Errorf("json marshal error with payload: %v", err)
	}

	encodedHeader := EncodeBase64URL(headerJSON)
	encodedPayload := EncodeBase64URL(payloadJSON)

	message := fmt.Sprintf("%s.%s", encodedHeader, encodedPayload)

	signature := hmac.New(sha256.New, []byte(secret))
	signature.Write([]byte(message))
	encodedSignature := EncodeBase64URL(signature.Sum(nil))

	reconstructedToken := fmt.Sprintf("%s.%s", message, encodedSignature)

	return token == reconstructedToken, nil
}

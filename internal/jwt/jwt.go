package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Exp    int64  `json:"exp"`
}

func Generate(secret, userID, email string, ttl time.Duration) (string, error) {
	header := base64.RawURLEncoding.EncodeToString(mustJSON(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}))

	claims := Claims{
		UserID: userID,
		Email:  email,
		Exp:    time.Now().Add(ttl).Unix(),
	}
	payload := base64.RawURLEncoding.EncodeToString(mustJSON(claims))

	sig := sign(secret, header+"."+payload)
	return header + "." + payload + "." + sig, nil
}

func Verify(secret, token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	expected := sign(secret, parts[0]+"."+parts[1])
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return nil, ErrInvalidToken
	}

	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(raw, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if time.Now().Unix() > claims.Exp {
		return nil, ErrExpiredToken
	}

	return &claims, nil
}

func sign(secret, data string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

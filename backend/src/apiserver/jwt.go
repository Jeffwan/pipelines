package main

import (
	"crypto/rsa"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// jwtAuthenticator checks the jwt token.
type jwtAuthenticator struct {
	publicKey *rsa.PublicKey
}

// from https://bytedance.feishu.cn/docs/doccnY7xSMZEJ0nT5cSB99
type claims struct {
	jwt.StandardClaims
	Username     string `json:"username,omitempty"`
	Type         string `json:"type,omitempty"`
	Organization string `json:"organization,omitempty"`
	AvatarURL    string `json:"avatar_url,omitempty"`
	Email        string `json:"email,omitempty"`
	WorkCountry  string `json:"work_country,omitempty"`
}

func (c *claims) Valid() error {
	return c.StandardClaims.Valid()
}

func NewJwtAuthenticator(path string) (*jwtAuthenticator, error) {
	publicKey, err := getPublicKey(path)
	if err != nil {
		return nil, err
	}
	return &jwtAuthenticator{
		publicKey: publicKey,
	}, nil
}

func (j *jwtAuthenticator) AuthenticateToken(token string) (string, error) {
	c := &claims{}
	t, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (interface{}, error) {
		return j.publicKey, nil
	})

	if err != nil {
		return "", err
	}
	if !t.Valid {
		return "", t.Claims.Valid()
	}
	return c.Username, nil
}

// the path can be a url or a file path.
func getPublicKey(path string) (*rsa.PublicKey, error) {
	var data []byte
	var err error

	if strings.HasPrefix(path, "http") {
		var resp *http.Response
		resp, err = http.Get(path) // nolint
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		data, err = ioutil.ReadAll(resp.Body)
	} else {
		data, err = ioutil.ReadFile(path)
	}

	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(data)
}

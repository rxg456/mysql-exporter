package handler

import (
	"encoding/base64"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthSecrets map[string]string

func Auth(handler http.Handler, secrets AuthSecrets) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		secret := request.Header.Get("Authorization")
		if !isAuth(secret, secrets) {
			response.Header().Set("WWW-Authenticate", `Basic realm=""`)
			response.WriteHeader(401)
			return
		}
		handler.ServeHTTP(response, request)

	})
}

func isAuth(secret string, secrets AuthSecrets) bool {
	if secrets == nil {
		return true
	}
	nodes := strings.Fields(secret)
	if len(nodes) != 2 {
		return false
	}
	plaintext, err := base64.StdEncoding.DecodeString(nodes[1])
	if err != nil {
		return false
	}
	nodes = strings.SplitN(string(plaintext), ":", 2)
	if len(nodes) != 2 {
		return false
	}
	hasher, ok := secrets[nodes[0]]
	return ok && bcrypt.CompareHashAndPassword([]byte(hasher), []byte(nodes[1])) == nil
}

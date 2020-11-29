package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/geobuff/geobuff-api/config"
)

// Response holds the response message.
type Response struct {
	Message string `json:"message"`
}

// Jwks holds a list of json web keys.
type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

// JSONWebKeys holds fields related to the JSON Web Key Set for this API.
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// GetJwtMiddleware returns the Auth0 middleware instance we will use to handle authorized endpoints.
func GetJwtMiddleware() *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			aud := config.Values.Auth0.Audience
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("invalid audience")
			}

			iss := config.Values.Auth0.Issuer
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(fmt.Sprintf("%s.well-known/jwks.json", config.Values.Auth0.Issuer))

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

// HasUserReadPermissions confirms the user making a request has the correct permissions to read user data.
func HasUserReadPermissions(request *http.Request) (bool, error) {
	permissions, err := getPermissions(request)
	if err != nil {
		return false, err
	}

	return permissionPresent(permissions, "read:users"), nil
}

// HasUserWritePermissions confirms the user making a request has the correct permissions to update users.
func HasUserWritePermissions(request *http.Request) (bool, error) {
	permissions, err := getPermissions(request)
	if err != nil {
		return false, err
	}

	return permissionPresent(permissions, "write:users"), nil
}

func getPermissions(request *http.Request) ([]interface{}, error) {
	header := request.Header.Get("Authorization")
	tokenString := header[7:]
	parser := new(jwt.Parser)
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		permissions := claims["permissions"]
		if permissions == nil {
			return nil, nil
		}
		return permissions.([]interface{}), nil
	}
	return nil, errors.New("failed to parse claims from token")
}

func permissionPresent(permissions []interface{}, target string) bool {
	for _, permission := range permissions {
		if permission == target {
			return true
		}
	}
	return false
}

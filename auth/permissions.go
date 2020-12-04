package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/geobuff/geobuff-api/database"
)

const (
	// ReadUsers gives permission to read user data.
	ReadUsers = "read:users"

	// WriteUsers gives permission to update users.
	WriteUsers = "write:users"

	// ReadScores gives permission to read score data.
	ReadScores = "read:scores"

	// WriteScores gives permission to update scores.
	WriteScores = "write:scores"

	// WriteLeaderboard gives permission to update leaderboard entries.
	WriteLeaderboard = "write:leaderboard"
)

// HasPermission confirms the user making a request has the correct permissions to complete the action.
func HasPermission(request *http.Request, permission string) (bool, error) {
	permissions, err := getPermissions(request)
	if err != nil {
		return false, err
	}
	return permissionPresent(permissions, permission), nil
}

// ValidUser confirms the user making a request is either making changes to their own data or has the correct
// permissions to complete the action.
func ValidUser(request *http.Request, userID int, permission string) (int, error) {
	switch user, err := database.GetUser(userID); err {
	case sql.ErrNoRows:
		return http.StatusNotFound, err
	case nil:
		if matchingUser, err := matchingUser(request, user.Username); err != nil {
			return http.StatusInternalServerError, err
		} else if !matchingUser {
			if hasPermission, err := HasPermission(request, permission); err != nil {
				return http.StatusInternalServerError, err
			} else if !hasPermission {
				return http.StatusUnauthorized, errors.New("missing or invalid permissions")
			}
		}
	default:
		return http.StatusInternalServerError, err
	}
	return 200, nil
}

func getPermissions(request *http.Request) ([]interface{}, error) {
	token, _, err := parseToken(request)
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

func parseToken(request *http.Request) (*jwt.Token, []string, error) {
	header := request.Header.Get("Authorization")
	tokenString := header[7:]
	parser := new(jwt.Parser)
	return parser.ParseUnverified(tokenString, jwt.MapClaims{})
}

func permissionPresent(permissions []interface{}, target string) bool {
	for _, permission := range permissions {
		if permission == target {
			return true
		}
	}
	return false
}

func matchingUser(request *http.Request, target string) (bool, error) {
	username, err := getUsername(request)
	if err != nil {
		return false, err
	}
	return username == target, nil
}

func getUsername(request *http.Request) (string, error) {
	token, _, err := parseToken(request)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return fmt.Sprintf("%s", claims["http://geobuff.com/username"]), nil
	}
	return "", errors.New("failed to parse claims from token")
}

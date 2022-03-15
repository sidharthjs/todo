package jwt

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestJWTMethods(t *testing.T) {

	testCases := []struct {
		userID           string
		userName         string
		expectedUserID   string
		expectedUserName string
	}{{
		userID:           "1001",
		userName:         "john101",
		expectedUserID:   "1001",
		expectedUserName: "john101",
	},
	}

	assert := assert.New(t)
	for _, testCase := range testCases {
		token, err := CreateJWTToken(testCase.userID, testCase.userName)
		assert.Nil(err)

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		assert.Nil(err)

		userID, userName, err := GetUserFromJWTToken(parsedToken)
		assert.Nil(err)
		assert.Equal(testCase.expectedUserID, userID)
		assert.Equal(testCase.expectedUserName, userName)
	}
}

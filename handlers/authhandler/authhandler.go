package authhandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"local/sidharthjs/todo/jwt"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

const jwtSecret = "aJWTSecret"

//AuthHandler struct definition
type AuthHandler struct {
	ClientID       string
	ClientSecret   string
	LoginURI       string
	AccessTokenURI string
	RedirectURI    string
	ProfileURI     string
	UsersService   string
}

// New returns AuthHandler
func New(clientID, clientSecret, loginURI, accessTokenURI, redirectURI, profileURI, usersService string) *AuthHandler {
	return &AuthHandler{
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		LoginURI:       loginURI,
		AccessTokenURI: accessTokenURI,
		RedirectURI:    redirectURI,
		ProfileURI:     profileURI,
		UsersService:   usersService,
	}
}

//InitiateOAuth redirects to the github login page
func (ah *AuthHandler) InitiateOAuth(c *fiber.Ctx) error {
	redirectURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&state=state", ah.LoginURI, ah.ClientID, ah.RedirectURI)
	return c.Redirect(redirectURL)
}

// ProcessCallback executes the logic of callback page after OAuth initialization
func (ah *AuthHandler) ProcessCallback(c *fiber.Ctx) error {

	// Get github access token of the user
	data := url.Values{}
	data.Set("client_id", ah.ClientID)
	data.Set("client_secret", ah.ClientSecret)
	data.Set("code", c.Query("code"))
	data.Set("redirect_uri", ah.RedirectURI)

	req, err := http.NewRequest("POST", ah.AccessTokenURI, strings.NewReader(data.Encode()))
	if err != nil {
		log.Errorf("error while creating the token request: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error during authentication",
		})
	}
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("error while making token request: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error during authentication",
		})
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error while reading token response: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error during authentication",
		})

	}

	m := make(map[string]string)
	err = json.Unmarshal(b, &m)
	if err != nil {
		log.Errorf("error while unmarshaling token response: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error during authentication",
		})
	}
	log.Debugf("access_token: %s", m["access_token"])

	// Get profile id, username by hitting github user API
	req, err = http.NewRequest("GET", ah.ProfileURI+"/user", strings.NewReader(data.Encode()))
	if err != nil {
		log.Errorf("error while creating the token request: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error during authentication",
		})
	}
	req.Header.Add("Authorization", "Bearer "+m["access_token"])

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("error while making github user request: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error during authentication",
		})
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error while reading github user response: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error during authentication",
		})
	}

	user := struct {
		ID       int64  `json:"id"`
		Username string `json:"login"`
	}{}
	err = json.Unmarshal(b, &user)
	if err != nil {
		log.Errorf("error while unmarshaling github user response: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error during authentication",
		})
	}

	// Create JWT token
	token, err := jwt.CreateJWTToken(strconv.FormatInt(user.ID, 10), user.Username)
	if err != nil {
		log.Errorf("error while generating JWT token: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error during authentication",
		})
	}
	log.Debugf("JWT token: %s", token)

	// Create User
	postBody, _ := json.Marshal(map[string]string{
		"username": user.Username,
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err = http.Post(ah.UsersService+"/users/"+strconv.FormatInt(user.ID, 10), "application/json", responseBody)
	if err != nil {
		log.Errorf("error while creating user: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error while creating user",
		})
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error while reading create user response: %s", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error while reading create user response",
		})
	}
	log.Printf(string(body))
	log.Debugf("User %s successfully created", user.Username)

	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Welcome %s!\nYour JWT token: %s\n\nPlease refer to README.md for curl commands to test the app.", user.Username, token))
}

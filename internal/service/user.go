package service

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/CoRide-tw/backend/internal/config"
	"github.com/CoRide-tw/backend/internal/db"
	"github.com/CoRide-tw/backend/internal/model"
	"github.com/CoRide-tw/backend/internal/util"
	"github.com/gin-gonic/gin"
)

type userSvc struct{}

func (s *userSvc) OauthUrl(c *gin.Context) {
	baseUrl, err := url.Parse("https://accounts.google.com/o/oauth2/v2/auth")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	params := url.Values{}
	params.Add("client_id", config.Env.GoogleOAuthClientId)
	params.Add("redirect_uri", config.Env.GoogleOAuthRedirectUrl)
	params.Add("response_type", "code")
	params.Add("scope", config.Env.GoogleOauthScope)
	params.Add("access_type", "offline")
	params.Add("include_granted_scopes", "true")
	baseUrl.RawQuery = params.Encode()

	c.JSON(http.StatusOK, gin.H{"url": baseUrl.String()})
}

type oauthCode struct {
	Code string `json:"code"`
}

func (s *userSvc) OAuthLogin(c *gin.Context) {
	// get code
	var res oauthCode
	if err := c.BindJSON(&res); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := s.getAccessToken(res.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData, err := s.getUserData(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// upsert user
	userResp, upsertErr := db.UpsertUser(&model.User{
		Name:       userData.Name,
		Email:      userData.Email,
		GoogleId:   userData.GoogleId,
		PictureUrl: userData.PictureUrl,
	})
	if upsertErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": upsertErr.Error()})
		return
	}

	// sign JWT token
	jwtToken, tokenErr := util.GenerateJWT(userResp.Id, config.Env.CoRideJwtSecret)
	if tokenErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"user":  userResp,
		"token": jwtToken,
	})
}

func (s *userSvc) Get(c *gin.Context) {
	stringId := c.Param("id")
	userId, err := strconv.Atoi(stringId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be integer"})
		return
	}

	authUid, authUidExist := c.Get("userId")
	if !authUidExist {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}
	if authUid != int32(userId) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	// get user from db
	user, err := db.GetUser(int32(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *userSvc) Update(c *gin.Context) {
	stringId := c.Param("id")
	userId, err := strconv.Atoi(stringId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be integer"})
		return
	}
	authUid, authUidExist := c.Get("userId")
	if !authUidExist {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}
	if authUid != int32(userId) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	var user model.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := db.UpdateUser(int32(userId), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

type oauthExchangeRes struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func (s *userSvc) getAccessToken(code string) (*oauthExchangeRes, error) {
	baseUrl, err := url.Parse("https://oauth2.googleapis.com/token")
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("code", code)
	params.Add("client_id", config.Env.GoogleOAuthClientId)
	params.Add("client_secret", config.Env.GoogleOAuthClientSecret)
	params.Add("redirect_uri", config.Env.GoogleOAuthRedirectUrl)
	params.Add("grant_type", "authorization_code")
	baseUrl.RawQuery = params.Encode()

	req, err := http.NewRequest(
		"POST",
		baseUrl.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		var err error
		return nil, err
	}

	var res oauthExchangeRes
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

type userData struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	GoogleId   string `json:"sub"`
	PictureUrl string `json:"picture"`
}

func (s *userSvc) getUserData(accessToken string) (*userData, error) {
	baseUrl, err := url.Parse("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("access_token", accessToken)
	baseUrl.RawQuery = params.Encode()

	req, err := http.NewRequest(
		"GET",
		baseUrl.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Read the response as a byte slice
	var user userData
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

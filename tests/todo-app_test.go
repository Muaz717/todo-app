package tests

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	testmodels "github.com/Muaz717/todo-app/tests/models"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	host           = "0.0.0.0:8083"
	passDefaultLen = 10
)

func TestTodo_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	email := gofakeit.Email()
	pass := fakeRandomPassword()

	reqBody := testmodels.RegRequest{
		Email:    email,
		Password: pass,
	}

	e := httpexpect.Default(t, u.String())

	e.POST("/auth/sign-up").
		WithJSON(reqBody).
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("status").HasValue("status", "OK").
		ContainsKey("msg").HasValue("msg", "You successfully registered")

	var loginResp testmodels.LoginResp

	e.POST("/auth/sign-in").
		WithJSON(reqBody).
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("status").HasValue("status", "OK").
		ContainsKey("msg").HasValue("msg", "You got token").
		Decode(&loginResp)

	token := loginResp.Token
	require.NotEmpty(t, token)

	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte("Darm"), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, email, claims["email"].(string))

	const deltaSeconds = 1

	time, err := time.ParseDuration("12h")
	assert.NoError(t, err)
	// check if exp of token is in correct range, ttl get from st.Cfg.TokenTTL
	assert.InDelta(t, loginTime.Add(time).Unix(), claims["exp"].(float64), deltaSeconds)

	itemReqBody := testmodels.CreateReq{
		Title:       gofakeit.JobTitle(),
		Description: gofakeit.JobDescriptor(),
	}

	bearerToken := fmt.Sprintf("Bearer %s", token)
	e.POST("/api/items/").
		WithHeader("Authorization", bearerToken).
		WithJSON(itemReqBody).
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("status").HasValue("status", "OK").
		ContainsKey("msg").HasValue("msg", "Item successfully created")

	e.GET("/api/items/").
		WithHeader("Authorization", bearerToken).
		Expect().
		Status(200)
}

func fakeRandomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}

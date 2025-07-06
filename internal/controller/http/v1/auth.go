package v1

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/azoma13/auth-service/config"
	"github.com/azoma13/auth-service/internal/service"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var source = "logInWithId"

type authRoutes struct {
	authService service.Auth
}

func newAuthRoutes(g *echo.Group, authService service.Auth) {
	r := authRoutes{
		authService: authService,
	}

	g.POST("/sign-up", r.signUp)
	g.POST("/sign-in", r.signIn)
	g.POST("/log-in", r.logIn)
}

type signInput struct {
	Id       string
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

func (r *authRoutes) signUp(c echo.Context) error {
	var input signInput

	if err := c.Bind(&input); err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.authService.CreateUser(c.Request().Context(), service.AuthCreateUserInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Id string `json:"id"`
	}

	return c.JSON(http.StatusCreated, response{
		Id: id,
	})
}

func (r *authRoutes) signIn(c echo.Context) error {
	var input signInput

	if err := c.Bind(&input); err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	claimsAccess := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(config.Cfg.AccessTokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   "access_token",
	}

	claimsRefresh := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(config.Cfg.RefreshTokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   "refresh_token",
	}

	ctx := context.WithValue(c.Request().Context(), "source", "sign-in")
	c.SetRequest(c.Request().WithContext(ctx))
	accessToken, refreshToken, err := getAccessAndRefreshToken(c.Request().Context(), r.authService, input, claimsAccess, claimsRefresh)
	if err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	err = createAccount(c, r.authService, input.Username, refreshToken)
	if err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	setCookie(c, "accessToken", accessToken, config.Cfg.AccessTokenTTL, 0)
	setCookie(c, "refreshToken", base64.StdEncoding.EncodeToString([]byte(refreshToken)), config.Cfg.RefreshTokenTTL, 0)

	return c.JSON(http.StatusNoContent, nil)
}

func getAccessAndRefreshToken(ctx context.Context, authService service.Auth, input signInput, claimsAccess, claimsRefresh jwt.StandardClaims) (string, string, error) {
	src := ctx.Value("source").(string)
	if src != source {
		accessToken, err := authService.GenerateToken(ctx, service.AuthGenerateTokenInput{
			Username: input.Username,
			Password: input.Password,
			TokenClaims: service.TokenClaims{
				StandardClaims: claimsAccess,
			},
		})

		if err != nil {
			return "", "", err
		}

		refreshToken, err := authService.GenerateToken(ctx, service.AuthGenerateTokenInput{
			Username: input.Username,
			Password: input.Password,
			TokenClaims: service.TokenClaims{
				StandardClaims: claimsRefresh,
			},
		})
		if err != nil {
			return "", "", err
		}
		return accessToken, refreshToken, nil
	} else {

		accessToken, err := authService.GenerateToken(ctx, service.AuthGenerateTokenInput{
			Id: input.Id,
			TokenClaims: service.TokenClaims{
				StandardClaims: claimsAccess,
			},
		})

		if err != nil {
			return "", "", err
		}
		refreshToken, err := authService.GenerateToken(ctx, service.AuthGenerateTokenInput{
			Id: input.Id,
			TokenClaims: service.TokenClaims{
				StandardClaims: claimsRefresh,
			},
		})
		if err != nil {
			return "", "", err
		}
		return accessToken, refreshToken, nil
	}
}

func setCookie(c echo.Context, nameCookie, token string, tokenTTL time.Duration, maxAge int) {
	c.SetCookie(&http.Cookie{
		Name:     nameCookie,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(tokenTTL),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   maxAge,
	})
}

func createAccount(c echo.Context, authService service.Auth, username, refreshToken string) error {
	err := authService.CreateAccount(c.Request().Context(), service.AuthCreateAccountInput{
		Username:      username,
		RefreshToken:  refreshToken,
		UserAgent:     c.Request().Header.Get("User-Agent"),
		XForwardedFor: c.Request().Header.Get("X-Forwarded-For"),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *authRoutes) logIn(c echo.Context) error {
	id := c.Request().URL.Query().Get("id")
	if id == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}
	var input signInput
	input.Id = id

	claimsAccess := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(config.Cfg.AccessTokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   "access_token",
	}

	claimsRefresh := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(config.Cfg.RefreshTokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   "refresh_token",
	}

	ctx := context.WithValue(c.Request().Context(), "source", source)
	c.SetRequest(c.Request().WithContext(ctx))
	accessToken, refreshToken, err := getAccessAndRefreshToken(ctx, r.authService, input, claimsAccess, claimsRefresh)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	err = r.authService.CreateAccount(c.Request().Context(), service.AuthCreateAccountInput{
		UserId:        id,
		RefreshToken:  refreshToken,
		UserAgent:     c.Request().Header.Get("User-Agent"),
		XForwardedFor: c.Request().Header.Get("X-Forwarded-For"),
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	setCookie(c, "accessToken", accessToken, config.Cfg.AccessTokenTTL, 0)
	setCookie(c, "refreshToken", base64.StdEncoding.EncodeToString([]byte(refreshToken)), config.Cfg.RefreshTokenTTL, 0)

	return c.JSON(http.StatusNoContent, nil)
}

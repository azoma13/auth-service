package v1

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/azoma13/auth-service/config"
	"github.com/azoma13/auth-service/internal/entity"
	"github.com/azoma13/auth-service/internal/service"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type accountRoutes struct {
	accountService service.Account
}

func newAccountRoutes(g *echo.Group, accountService service.Account) {
	r := &accountRoutes{
		accountService: accountService,
	}

	g.GET("/guid", r.getGuid)
	g.PUT("/refresh", r.updateTokens)
	g.DELETE("/sign-out", r.signOut)
}

func (r *accountRoutes) getGuid(c echo.Context) error {
	id, err := getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid get userid")
		return fmt.Errorf("invalid get userid")
	}

	type response struct {
		Id string `json:"guid"`
	}

	return c.JSON(http.StatusOK, response{
		Id: id,
	})
}

func (r *accountRoutes) updateTokens(c echo.Context) error {
	id, err := getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid get userid")
		return fmt.Errorf("invalid get userid")
	}

	cookie, err := c.Cookie("refreshToken")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	oldRefreshToken, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	account, err := r.accountService.GetAccount(c.Request().Context(), service.AccountGetInput{
		UserId:        id,
		RefreshToken:  string(oldRefreshToken),
		UserAgent:     c.Request().Header.Get("User-Agent"),
		XForwardedFor: c.Request().Header.Get("X-Forwarded-For"),
	})

	switch err {
	case service.ErrDifferentUserAgent:
		c.Set("source", service.ErrDifferentUserAgent)
		err := r.signOut(c)
		return err
	case service.ErrDifferentXForwardedFor:
		log.Println("sdf")
		webhook()
	case nil:
	default:
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	accessToken, refreshToken, err := r.newTokens(c, account)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	err = r.accountService.UpdateRefreshToken(c.Request().Context(), service.AccountUpdateInput{
		Id:            account.Id,
		RefreshToken:  refreshToken,
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

func getId(c echo.Context) (string, error) {
	userId := c.Get(userIdCtx)
	if userId == nil {
		return "", fmt.Errorf("invalid get userid")
	}

	id, ok := userId.(string)
	if !ok {
		return "", fmt.Errorf("error convert userid")
	}

	return id, nil
}

func (r *accountRoutes) newTokens(c echo.Context, account entity.Account) (string, string, error) {
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

	accessToken, err := r.accountService.GenerateToken(c.Request().Context(), service.TokenClaims{
		StandardClaims: claimsAccess,
		UserId:         account.UserId,
	})
	if err != nil {
		return "", "", err
	}

	refreshToken, err := r.accountService.GenerateToken(c.Request().Context(), service.TokenClaims{
		StandardClaims: claimsRefresh,
		UserId:         account.UserId,
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func webhook() {

}

func (r *accountRoutes) signOut(c echo.Context) error {
	cookie, err := c.Cookie("refreshToken")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	refreshToken, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	userId := c.Get(userIdCtx)
	if userId == nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	id, ok := userId.(string)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	err = r.accountService.DeleteAccount(c.Request().Context(), service.AuthDeleteAccountInput{
		UserId:       id,
		RefreshToken: string(refreshToken),
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	setCookie(c, "accessToken", "", 0, -1)
	setCookie(c, "refreshToken", "", config.Cfg.RefreshTokenTTL, -1)

	source := c.Get("source")
	if source == service.ErrDifferentUserAgent {
		newErrorResponse(c, http.StatusForbidden, "please sign-in again")
		return service.ErrDifferentUserAgent
	}
	return c.JSON(http.StatusNoContent, nil)
}

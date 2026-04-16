package routes

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func sessionCookieName() string {
	if v := os.Getenv("SESSION_COOKIE_NAME"); v != "" {
		return v
	}
	return "session"
}

func sessionCookieSecure() bool {
	return os.Getenv("SESSION_COOKIE_SECURE") == "true"
}

func sessionCookieDomain() string {
	return os.Getenv("SESSION_COOKIE_DOMAIN")
}

func sessionCookieSameSite() http.SameSite {
	switch strings.ToLower(os.Getenv("SESSION_COOKIE_SAMESITE")) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

func setSessionCookie(c *gin.Context, sessionID string) {
	maxAge := int(sessionManager.IdleExpiration() / time.Second)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     sessionCookieName(),
		Value:    sessionID,
		Domain:   sessionCookieDomain(),
		Path:     "/",
		MaxAge:   maxAge,
		Expires:  time.Now().Add(sessionManager.IdleExpiration()),
		HttpOnly: true,
		Secure:   sessionCookieSecure(),
		SameSite: sessionCookieSameSite(),
	})
}

func clearSessionCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     sessionCookieName(),
		Value:    "",
		Domain:   sessionCookieDomain(),
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   sessionCookieSecure(),
		SameSite: sessionCookieSameSite(),
	})
}

func sessionIDFromRequest(c *gin.Context) string {
	if cookie, err := c.Cookie(sessionCookieName()); err == nil && cookie != "" {
		return cookie
	}
	return ""
}

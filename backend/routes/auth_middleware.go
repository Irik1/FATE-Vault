package routes

import (
	"context"
	"net/http"
	"time"

	"FATE-Vault/backend/db"
	"FATE-Vault/backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserFromSessionID loads the current user using opaque session data.
func UserFromSessionID(ctx context.Context, sessionID string) (*models.Users, *Session, error) {
	if db.Client == nil {
		return nil, nil, mongo.ErrClientDisconnected
	}
	coll := db.Client.Database("main").Collection("users")
	session := sessionManager.ReadValid(sessionID)
	if session == nil {
		return nil, nil, mongo.ErrNoDocuments
	}
	filter := bson.M{"_id": session.UserID}
	var user models.Users
	err := coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, nil, err
	}
	return &user, session, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Help caches keep responses separated per session cookie.
		c.Writer.Header().Add("Vary", "Cookie")
		c.Writer.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)

		sessionID := sessionIDFromRequest(c)
		if sessionID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user, session, err := UserFromSessionID(ctx, sessionID)
		if err != nil {
			clearSessionCookie(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired session"})
			c.Abort()
			return
		}

		c.Set("userId", user.ID)
		c.Set("username", user.Username)
		c.Set("role", user.Role)
		c.Set("user", *user)
		sessionManager.Touch(session)
		setSessionCookie(c, session.ID)

		c.Next()
	}
}

func GetUserFromContext(c *gin.Context) (*models.Users, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	userModel, ok := user.(models.Users)
	if !ok {
		return nil, false
	}
	return &userModel, true
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireAdminOrOwner(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		userId, _ := c.Get("userId")
		resourceId := c.Param(paramName)

		if role == "admin" || userId == resourceId {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		c.Abort()
	}
}

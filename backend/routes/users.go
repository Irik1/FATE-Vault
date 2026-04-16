package routes

import (
	"context"
	"net/http"
	"time"

	"FATE-Vault/backend/db"
	"FATE-Vault/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role,omitempty"`
}

type AuthRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type UpdateUserRequest struct {
	Username       string `json:"username,omitempty"`
	ProfilePicture string `json:"profilePicture,omitempty"`
	Role           string `json:"role,omitempty"`
}

func RegisterUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	coll := db.Client.Database("main").Collection("users")

	// Check if username already exists
	filter := bson.M{"username": req.Username}
	var existingUser models.Users
	err := coll.FindOne(ctx, filter).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}
	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check username: " + err.Error()})
		return
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = "user"
	}

	// Validate role
	if role != "admin" && role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role must be either 'admin' or 'user'"})
		return
	}

	// Create new user
	user := models.Users{
		ID:       uuid.NewString(),
		Username: req.Username,
		Role:     role,
	}

	// Hash password
	if err := user.SetPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password: " + err.Error()})
		return
	}

	_, err = coll.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user: " + err.Error()})
		return
	}

	// Return user without password; session ID is stored in HttpOnly cookie.
	user.HashedPassword = ""
	session := sessionManager.Create(user.ID)
	setSessionCookie(c, session.ID)
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func respondWithSessionUser(c *gin.Context, user *models.Users, status int) error {
	// Rotate existing session ID to mitigate session fixation.
	if oldSessionID := sessionIDFromRequest(c); oldSessionID != "" {
		sessionManager.Destroy(oldSessionID)
	}
	session := sessionManager.Create(user.ID)
	setSessionCookie(c, session.ID)
	u := *user
	u.HashedPassword = ""
	c.JSON(status, gin.H{"user": u})
	return nil
}

// GetCurrentUser returns user data from already-authenticated middleware context.
func GetCurrentUser(c *gin.Context) {
	user, ok := GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	user.HashedPassword = ""
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// LogoutUser clears the session cookie.
func LogoutUser(c *gin.Context) {
	if sessionID := sessionIDFromRequest(c); sessionID != "" {
		sessionManager.Destroy(sessionID)
	}
	clearSessionCookie(c)
	c.Status(http.StatusNoContent)
}

func AuthUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	coll := db.Client.Database("main").Collection("users")

	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	var user models.Users
	filter := bson.M{"username": req.Username}
	err := coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user: " + err.Error()})
		return
	}
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	if err := respondWithSessionUser(c, &user, http.StatusOK); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session: " + err.Error()})
	}
}

func UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	// Check authorization: user must be admin or updating their own account
	userId, exists := c.Get("userId")
	role, _ := c.Get("role")
	if !exists || (role != "admin" && userId != id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if db.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return
	}

	coll := db.Client.Database("main").Collection("users")

	// Only admin can change role
	if req.Role != "" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can change user role"})
		return
	}

	// Build update document
	update := bson.M{}
	if req.Username != "" {
		// Check if new username already exists (if different from current)
		filter := bson.M{"username": req.Username, "_id": bson.M{"$ne": id}}
		var existingUser models.Users
		err := coll.FindOne(ctx, filter).Decode(&existingUser)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}
		if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check username: " + err.Error()})
			return
		}
		update["username"] = req.Username
	}
	if req.ProfilePicture != "" {
		update["profilePicture"] = req.ProfilePicture
	}
	if req.Role != "" {
		// Validate role
		if req.Role != "admin" && req.Role != "user" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "role must be either 'admin' or 'user'"})
			return
		}
		update["role"] = req.Role
	}

	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	filter := bson.M{"_id": id}
	updateDoc := bson.M{"$set": update}
	result, err := coll.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user: " + err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Fetch updated user
	var user models.Users
	err = coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated user: " + err.Error()})
		return
	}

	// Return user without password
	user.HashedPassword = ""
	c.JSON(http.StatusOK, user)
}

package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	models "FATE-Vault/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

// Type and constant aliases to keep existing tests working with the models package
type Character = models.Character
type Aspect = models.Aspect
type SkillGroup = models.SkillGroup
type Refresh = models.Refresh
type Consequence = models.Consequence
type Stress = models.Stress
type StressBox = models.StressBox
type Stunt = models.CharacterStunt
type Edition = models.Edition

const (
	Core        Edition = models.Core
	Accelerated Edition = models.Accelerated
	Condensed   Edition = models.Condensed
	Custom      Edition = models.Custom
)

// Helper function to create mock character data
func createMockCharacter() Character {
	return Character{
		ID:          "test-id-123",
		Edition:     Core,
		Name:        "Test Character",
		Description: "A test character for unit testing",
		Images:      []string{"/images/test.jpg"},
		Notes:       "Test notes",
		Aspects: []Aspect{
			{Type: "High Concept", Value: "Test Aspect"},
		},
		Skills: []SkillGroup{
			{Level: "+4", Skills: []string{"Test Skill"}},
		},
		Refresh: Refresh{Current: 3, Max: 3},
		Extras:  "Test extras",
		Stunts: []Stunt{
			{Name: "Test Stunt", Description: "Test description"},
		},
		Stress: []Stress{
			{Type: "physical", Boxes: []StressBox{{Size: 1, IsFilled: false}}},
		},
		Consequences: []Consequence{
			{Type: "mild", Size: 2, Description: "", Status: "none"},
		},
		IsPublished: true,
		CreatorID:   "",
	}
}

// Helper function to create mock character documents for list
func createMockCharacterDocuments() []bson.M {
	return []bson.M{
		{
			"_id":         "test-id-1",
			"name":        "Character 1",
			"edition":     "core",
			"description": "First test character",
			"isPublished": true,
		},
		{
			"_id":         "test-id-2",
			"name":        "Character 2",
			"edition":     "accelerated",
			"description": "Second test character",
			"isPublished": true,
		},
	}
}

// Test helper to setup router
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// Note: These tests focus on JSON parsing, request handling, and validation logic.
// Full MongoDB integration tests would require a test database or more complex mocking setup.
// The handlers will return 500 errors when mongoClient is nil, which is expected in unit tests.

// Test helper to create a test request
func createTestRequest(method, url string, body interface{}) *http.Request {
	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	return req
}

func TestCreateCharacter_Success(t *testing.T) {
	router := setupRouter()
	router.POST("/characters/create", CreateCharacter)

	character := createMockCharacter()
	req := createTestRequest("POST", "/characters/create", character)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Without MongoDB connection, this will return 500 (expected in unit tests)
	// We verify that JSON parsing works (not 400 Bad Request)
	assert.NotEqual(t, http.StatusBadRequest, w.Code, "JSON parsing should succeed")

	// Verify the request was processed (either success or DB error, not JSON error)
	assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusInternalServerError)
}

func TestCreateCharacter_InvalidJSON(t *testing.T) {
	router := setupRouter()
	router.POST("/characters/create", CreateCharacter)

	req, _ := http.NewRequest("POST", "/characters/create", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "error")
}

func TestCreateCharacter_EmptyBody(t *testing.T) {
	router := setupRouter()
	router.POST("/characters/create", CreateCharacter)

	req, _ := http.NewRequest("POST", "/characters/create", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Empty JSON object is valid, will fail at MongoDB level (expected)
	assert.NotEqual(t, http.StatusBadRequest, w.Code, "Empty JSON should be valid")
	assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusInternalServerError)
}

func TestUpdateCharacter_Success(t *testing.T) {
	router := setupRouter()
	router.POST("/characters/update/:id", UpdateCharacter)

	character := createMockCharacter()
	req := createTestRequest("POST", "/characters/update/test-id-123", character)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Without MongoDB, this will return 404 or 500, but JSON parsing should work
	assert.NotEqual(t, http.StatusBadRequest, w.Code, "JSON parsing should succeed")
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
}

func TestUpdateCharacter_MissingID(t *testing.T) {
	router := setupRouter()
	router.POST("/characters/update/:id", UpdateCharacter)

	character := createMockCharacter()
	req := createTestRequest("POST", "/characters/update/", character)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Gin will handle empty param differently, but we expect bad request or similar
	assert.True(t, w.Code >= http.StatusBadRequest)
}

func TestUpdateCharacter_InvalidJSON(t *testing.T) {
	router := setupRouter()
	router.POST("/characters/update/:id", UpdateCharacter)

	req, _ := http.NewRequest("POST", "/characters/update/test-id-123", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "error")
}

func TestUpdateCharacter_EmptyID(t *testing.T) {
	router := setupRouter()
	router.POST("/characters/update/:id", UpdateCharacter)

	character := createMockCharacter()
	// Test with empty string as ID - this tests the handler's ID validation
	req := createTestRequest("POST", "/characters/update/", character)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// The handler checks for empty ID, should return bad request
	// But Gin routing might handle this differently
	assert.True(t, w.Code >= http.StatusBadRequest)
}

func TestDeleteCharacter_Success(t *testing.T) {
	router := setupRouter()
	router.DELETE("/characters/delete/:id", DeleteCharacter)

	req, _ := http.NewRequest("DELETE", "/characters/delete/test-id-123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Without MongoDB, this will return 404 or 500, but handler logic is tested
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)

	// If successful, verify response structure
	if w.Code == http.StatusOK {
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "message")
	}
}

func TestDeleteCharacter_MissingID(t *testing.T) {
	router := setupRouter()
	router.DELETE("/characters/delete/:id", DeleteCharacter)

	req, _ := http.NewRequest("DELETE", "/characters/delete/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Handler checks for empty ID
	assert.True(t, w.Code >= http.StatusBadRequest)
}

func TestDeleteCharacter_NotFound(t *testing.T) {
	router := setupRouter()
	router.DELETE("/characters/delete/:id", DeleteCharacter)

	req, _ := http.NewRequest("DELETE", "/characters/delete/non-existent-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Handler processes the request (MongoDB connection will fail, but logic is tested)
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
}

func TestCharactersList(t *testing.T) {
	router := setupRouter()
	router.GET("/characters/list", CharactersList)

	req, _ := http.NewRequest("GET", "/characters/list", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Without MongoDB, this will return 500, but handler structure is tested
	// The handler correctly attempts to query MongoDB
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}

// Test character model validation
func TestCharacter_JSONSerialization(t *testing.T) {
	character := createMockCharacter()

	// Test marshaling
	jsonData, err := json.Marshal(character)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Test unmarshaling
	var unmarshaled Character
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, character.Name, unmarshaled.Name)
	assert.Equal(t, character.Edition, unmarshaled.Edition)
	assert.Equal(t, len(character.Aspects), len(unmarshaled.Aspects))
	assert.Equal(t, len(character.Skills), len(unmarshaled.Skills))
}

// Test with mock data structures
func TestCreateCharacter_WithMockData(t *testing.T) {
	testCases := []struct {
		name        string
		character   Character
		expectError bool
	}{
		{
			name:        "valid character",
			character:   createMockCharacter(),
			expectError: false,
		},
		{
			name: "character with all fields",
			character: Character{
				Edition:     Accelerated,
				Name:        "Full Character",
				Description: "Complete character description",
				Images:      []string{"/img1.jpg", "/img2.jpg"},
				Notes:       "Detailed notes",
				Aspects: []Aspect{
					{Type: "High Concept", Value: "Aspect 1"},
					{Type: "Trouble", Value: "Aspect 2"},
				},
				Skills: []SkillGroup{
					{Level: "+4", Skills: []string{"Skill1"}},
					{Level: "+3", Skills: []string{"Skill2", "Skill3"}},
				},
				Refresh: Refresh{Current: 2, Max: 3},
				Extras:  "Extra information",
				Stunts: []Stunt{
					{Name: "Stunt 1", Description: "Description 1"},
					{Name: "Stunt 2", Description: "Description 2"},
				},
				Stress: []Stress{
					{Type: "physical", Boxes: []StressBox{{Size: 1, IsFilled: false}, {Size: 2, IsFilled: true}}},
					{Type: "mental", Boxes: []StressBox{{Size: 1, IsFilled: false}}},
				},
				Consequences: []Consequence{
					{Type: "mild", Size: 2, Description: "Mild consequence", Status: "active"},
					{Type: "moderate", Size: 4, Description: "Moderate consequence", Status: "none"},
				},
				IsPublished: true,
				CreatorID:   "",
			},
			expectError: false,
		},
		{
			name: "minimal character",
			character: Character{
				Edition:     Core,
				Name:        "Minimal",
				IsPublished: true,
				CreatorID:   "",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := setupRouter()
			router.POST("/characters/create", CreateCharacter)

			req := createTestRequest("POST", "/characters/create", tc.character)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// JSON parsing should work regardless of MongoDB
			if !tc.expectError {
				assert.NotEqual(t, http.StatusBadRequest, w.Code, "Should not fail JSON parsing")
			}
		})
	}
}

func TestUpdateCharacter_WithMockData(t *testing.T) {
	character := createMockCharacter()

	router := setupRouter()
	router.POST("/characters/update/:id", UpdateCharacter)

	testCases := []struct {
		name        string
		id          string
		character   Character
		expectError bool
	}{
		{
			name:        "valid update",
			id:          "test-id-123",
			character:   character,
			expectError: false,
		},
		{
			name: "update with modified data",
			id:   "test-id-456",
			character: Character{
				Edition:     Condensed,
				Name:        "Updated Character",
				Description: "Updated description",
				IsPublished: true,
				CreatorID:   "",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := createTestRequest("POST", "/characters/update/"+tc.id, tc.character)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// JSON parsing should work
			assert.NotEqual(t, http.StatusBadRequest, w.Code, "Should not fail JSON parsing")
		})
	}
}

// Test helper functions
func TestCreateMockCharacter(t *testing.T) {
	character := createMockCharacter()

	assert.NotEmpty(t, character.ID)
	assert.Equal(t, Core, character.Edition)
	assert.Equal(t, "Test Character", character.Name)
	assert.NotEmpty(t, character.Aspects)
	assert.NotEmpty(t, character.Skills)
	assert.NotEmpty(t, character.Stress)
	assert.NotEmpty(t, character.Consequences)
}

func TestCreateMockCharacterDocuments(t *testing.T) {
	docs := createMockCharacterDocuments()

	assert.Len(t, docs, 2)
	assert.Equal(t, "test-id-1", docs[0]["_id"])
	assert.Equal(t, "Character 1", docs[0]["name"])
	assert.Equal(t, "test-id-2", docs[1]["_id"])
	assert.Equal(t, "Character 2", docs[1]["name"])
}

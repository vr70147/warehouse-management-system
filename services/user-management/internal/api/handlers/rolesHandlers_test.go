package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"user-management/internal/api/handlers"
	"user-management/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateRole(t *testing.T) {
	// Set up the environment variable for JWT secret
	os.Setenv("JWT_SECRET", "test_secret")

	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_role.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Role{})

	r := SetupRouter()
	r.POST("/roles", func(c *gin.Context) {
		c.Set("account_id", uint(1)) // Set the account ID in the context
		handlers.CreateRole(db)(c)
	})

	// Define the test case for successful role creation
	t.Run("CreateRoleSuccess", func(t *testing.T) {
		role := model.Role{
			Role:         "Manager",
			Description:  "Manages the team",
			IsActive:     true,
			DepartmentID: 1,
		}
		jsonValue, _ := json.Marshal(role)
		req, _ := http.NewRequest("POST", "/roles", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Role created successfully", response["message"])
		assert.NotNil(t, response["data"])
	})

	// Define the test case for invalid request body
	t.Run("CreateRoleInvalidBody", func(t *testing.T) {
		invalidJson := `{"role": "Manager", "description": "Manages the team", "isActive": true "department_id": 1}` // Missing comma after true
		req, _ := http.NewRequest("POST", "/roles", bytes.NewBuffer([]byte(invalidJson)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		t.Log("Response Body:", w.Body.String()) // Add this line for debugging

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to read body", response["error"])
	})

	// Clean up the database
	db.Exec("DELETE FROM roles")
}

func TestUpdateRole(t *testing.T) {
	// Set up the environment variable for JWT secret
	os.Setenv("JWT_SECRET", "test_secret")

	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_role.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Role{})

	role := model.Role{
		Role:         "Manager",
		Description:  "Manages the team",
		IsActive:     true,
		DepartmentID: 1,
		AccountID:    1,
	}
	db.Create(&role)

	r := SetupRouter()
	r.PUT("/roles/:id", func(c *gin.Context) {
		c.Set("account_id", uint(1)) // Set the account ID in the context
		handlers.UpdateRole(db)(c)
	})

	// Define the test case for successful role update
	t.Run("UpdateRoleSuccess", func(t *testing.T) {
		updateData := struct {
			Role         string `json:"role"`
			Description  string `json:"description"`
			IsActive     bool   `json:"is_active"`
			DepartmentID uint   `json:"department_id"`
		}{
			Role:         "Updated Manager",
			Description:  "Updated description",
			IsActive:     false,
			DepartmentID: 2,
		}
		jsonValue, _ := json.Marshal(updateData)
		req, _ := http.NewRequest("PUT", "/roles/"+strconv.Itoa(int(role.ID)), bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Role updated successfully", response["message"])
		assert.NotNil(t, response["data"])
	})

	// Define the test case for invalid request body
	t.Run("UpdateRoleInvalidBody", func(t *testing.T) {
		invalidJson := `{"role": "Manager", "description": "Manages the team", "isActive": true "department_id": 1}` // Missing comma after true

		req, _ := http.NewRequest("PUT", "/roles/"+strconv.Itoa(int(role.ID)), bytes.NewBuffer([]byte(invalidJson)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request data", response["error"])
	})

	// Clean up the database
	db.Exec("DELETE FROM roles")
}

func TestGetRoles(t *testing.T) {
	// Set up the environment variable for JWT secret
	os.Setenv("JWT_SECRET", "test_secret")

	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_role.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Role{})

	role := model.Role{
		Role:         "Manager",
		Description:  "Manages the team",
		IsActive:     true,
		DepartmentID: 1,
		AccountID:    1,
	}
	db.Create(&role)

	r := SetupRouter()
	r.GET("/roles", func(c *gin.Context) {
		c.Set("account_id", uint(1)) // Set the account ID in the context
		handlers.GetRoles(db)(c)
	})

	// Define the test case for retrieving roles
	t.Run("GetRolesSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/roles", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Roles retrieved successfully", response["message"])
		assert.NotNil(t, response["roles"])
	})

	// Clean up the database
	db.Exec("DELETE FROM roles")
}

func TestSoftDeleteRole(t *testing.T) {
	// Set up the environment variable for JWT secret
	os.Setenv("JWT_SECRET", "test_secret")

	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_role.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Role{})

	role := model.Role{
		Role:         "Manager",
		Description:  "Manages the team",
		IsActive:     true,
		DepartmentID: 1,
		AccountID:    1,
	}
	db.Create(&role)

	r := SetupRouter()
	r.DELETE("/roles/:id", func(c *gin.Context) {
		c.Set("account_id", uint(1)) // Set the account ID in the context
		handlers.SoftDeleteRole(db)(c)
	})

	// Define the test case for successful role deletion
	t.Run("SoftDeleteRoleSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/roles/"+strconv.Itoa(int(role.ID)), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Role soft deleted successfully", response["message"])
	})

	// Define the test case for invalid role ID
	t.Run("SoftDeleteRoleInvalidID", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/roles/9999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Role not found", response["error"])
	})

	// Clean up the database
	db.Exec("DELETE FROM roles")
}
func TestHardDeleteRole(t *testing.T) {
	// Set up the environment variable for JWT secret
	os.Setenv("JWT_SECRET", "test_secret")

	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_role.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Role{})

	role := model.Role{
		Role:         "Manager",
		Description:  "Manages the team",
		IsActive:     true,
		DepartmentID: 1,
		AccountID:    1,
	}
	db.Create(&role)

	r := SetupRouter()
	r.DELETE("/roles/hard/:id", func(c *gin.Context) {
		c.Set("account_id", uint(1)) // Set the account ID in the context
		handlers.HardDeleteRole(db)(c)
	})

	// Define the test case for successful role hard deletion
	t.Run("HardDeleteRoleSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/roles/hard/"+strconv.Itoa(int(role.ID)), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Role hard deleted successfully", response["message"])
	})

	// Define the test case for invalid role ID
	t.Run("HardDeleteRoleInvalidID", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/roles/hard/9999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Role not found", response["error"])
	})

	// Clean up the database
	db.Exec("DELETE FROM roles")
}

func TestRecoverRole(t *testing.T) {
	// Set up the environment variable for JWT secret
	os.Setenv("JWT_SECRET", "test_secret")

	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_role.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Role{})

	role := model.Role{
		Role:         "Manager",
		Description:  "Manages the team",
		IsActive:     true,
		DepartmentID: 1,
		AccountID:    1,
	}
	db.Create(&role)
	db.Delete(&role) // Soft delete the role

	r := SetupRouter()
	r.POST("/roles/:id/recover", func(c *gin.Context) {
		c.Set("account_id", uint(1)) // Set the account ID in the context
		handlers.RecoverRole(db)(c)
	})

	// Define the test case for successful role recovery
	t.Run("RecoverRoleSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/roles/"+strconv.Itoa(int(role.ID))+"/recover", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Role recovered successfully", response["message"])
	})

	// Define the test case for invalid role ID
	t.Run("RecoverRoleInvalidID", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/roles/9999/recover", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Role not found", response["error"])
	})

	// Clean up the database
	db.Exec("DELETE FROM roles")
}

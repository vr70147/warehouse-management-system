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

// SetupTestDB initializes the database for testing
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test_department.db"), &gorm.Config{})
	assert.NoError(t, err)
	db.AutoMigrate(&model.Department{}, &model.Role{}, &model.User{})
	return db
}

func TestCreateDepartment(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	db := SetupTestDB(t)
	defer db.Exec("DELETE FROM departments")

	r := SetupRouter()
	r.POST("/departments", func(c *gin.Context) {
		c.Set("account_id", uint(1))
		handlers.CreateDepartment(db)(c)
	})

	t.Run("CreateDepartmentSuccess", func(t *testing.T) {
		department := model.Department{Name: "Engineering"}
		jsonValue, _ := json.Marshal(department)
		req, _ := http.NewRequest("POST", "/departments", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Department created successfully", response["message"])
		assert.NotNil(t, response["data"])
	})

	t.Run("CreateDepartmentInvalidBody", func(t *testing.T) {
		invalidJson := `{"name": "Engineering"`
		req, _ := http.NewRequest("POST", "/departments", bytes.NewBuffer([]byte(invalidJson)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request data", response["error"])
	})
}

func TestUpdateDepartment(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	db := SetupTestDB(t)
	defer db.Exec("DELETE FROM departments")

	department := model.Department{Name: "Engineering", AccountID: 1}
	db.Create(&department)

	r := SetupRouter()
	r.PUT("/departments/:id", func(c *gin.Context) {
		c.Set("account_id", uint(1))
		handlers.UpdateDepartment(db)(c)
	})

	t.Run("UpdateDepartmentSuccess", func(t *testing.T) {
		updateData := struct {
			Name string `json:"name"`
		}{
			Name: "Updated Engineering",
		}
		jsonValue, _ := json.Marshal(updateData)
		req, _ := http.NewRequest("PUT", "/departments/"+strconv.Itoa(int(department.ID)), bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Department updated successfully", response["message"])
		assert.NotNil(t, response["data"])
	})

	t.Run("UpdateDepartmentInvalidBody", func(t *testing.T) {
		invalidJson := `{"name": "Engineering"` // missing closing quote
		req, _ := http.NewRequest("PUT", "/departments/"+strconv.Itoa(int(department.ID)), bytes.NewBuffer([]byte(invalidJson)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request data", response["error"])
	})
}

func TestGetDepartments(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	db := SetupTestDB(t)
	defer db.Exec("DELETE FROM departments")

	department := model.Department{Name: "Engineering", AccountID: 1}
	result := db.Create(&department)
	if result.Error != nil {
		t.Fatalf("Failed to create department: %v", result.Error)
	}

	r := SetupRouter()
	r.GET("/departments", func(c *gin.Context) {
		c.Set("account_id", uint(1))
		handlers.GetDepartments(db)(c)
	})

	t.Run("GetDepartmentsSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/departments", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		t.Log("Response Body:", w.Body.String()) // Add this line for debugging
		assert.Equal(t, "Departments retrieved successfully", response["message"])
		assert.NotNil(t, response["department"])
	})
}

func TestSoftDeleteDepartment(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	db := SetupTestDB(t)
	defer db.Exec("DELETE FROM departments")

	department := model.Department{Name: "Engineering", AccountID: 1}
	db.Create(&department)

	r := SetupRouter()
	r.DELETE("/departments/:id", func(c *gin.Context) {
		c.Set("account_id", uint(1))
		handlers.SoftDeleteDepartment(db)(c)
	})

	t.Run("SoftDeleteDepartmentSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/departments/"+strconv.Itoa(int(department.ID)), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Department soft deleted successfully", response["message"])
	})

	t.Run("SoftDeleteDepartmentInvalidID", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/departments/9999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Department not found", response["error"])
	})
}

func TestHardDeleteDepartment(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	db := SetupTestDB(t)
	defer db.Exec("DELETE FROM departments")

	department := model.Department{Name: "Engineering", AccountID: 1}
	db.Create(&department)

	r := SetupRouter()
	r.DELETE("/departments/hard/:id", func(c *gin.Context) {
		c.Set("account_id", uint(1))
		handlers.HardDeleteDepartment(db)(c)
	})

	t.Run("HardDeleteDepartmentSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/departments/hard/"+strconv.Itoa(int(department.ID)), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Department hard deleted successfully", response["message"])
	})

	t.Run("HardDeleteDepartmentInvalidID", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/departments/hard/9999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Department not found", response["error"])
	})
}

func TestRecoverDepartment(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	db := SetupTestDB(t)
	defer db.Exec("DELETE FROM departments")

	department := model.Department{Name: "Engineering", AccountID: 1}
	db.Create(&department)
	db.Delete(&department)

	r := SetupRouter()
	r.POST("/departments/:id/recover", func(c *gin.Context) {
		c.Set("account_id", uint(1))
		handlers.RecoverDepartment(db)(c)
	})

	t.Run("RecoverDepartmentSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/departments/"+strconv.Itoa(int(department.ID))+"/recover", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Department recovered successfully", response["message"])
	})

	t.Run("RecoverDepartmentInvalidID", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/departments/9999/recover", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Department not found", response["error"])
	})
}

func TestGetUsersByDepartment(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	db := SetupTestDB(t)
	defer db.Exec("DELETE FROM users")
	defer db.Exec("DELETE FROM roles")
	defer db.Exec("DELETE FROM departments")

	department := model.Department{Name: "Engineering", AccountID: 1}
	db.Create(&department)

	role := model.Role{
		Role:         "Developer",
		DepartmentID: department.ID,
		AccountID:    1,
	}
	db.Create(&role)

	user := model.User{
		PersonalID: "12345",
		Name:       "Test User",
		Email:      "user@example.com",
		Age:        25,
		BirthDate:  "1999-01-01",
		RoleID:     role.ID,
		AccountID:  1,
		Permission: model.PermissionManager,
	}
	db.Create(&user)

	r := SetupRouter()
	r.GET("/departments/users", func(c *gin.Context) {
		c.Set("account_id", uint(1))
		handlers.GetUsersByDepartment(db)(c)
	})

	t.Run("GetUsersByDepartmentSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/departments/users?department=Engineering", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Users retrieved successfully", response["message"])
		assert.NotNil(t, response["users"])
	})

	t.Run("GetUsersByDepartmentInvalidDepartment", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/departments/users?department=InvalidDept", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "No users found in the specified department", response["error"])
	})

	t.Run("GetUsersByDepartmentMissingDepartment", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/departments/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Department name is required", response["error"])
	})
}

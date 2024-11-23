package tests_test

import (
	"bytes"
	"encoding/json"
	"inventory-management/internal/model"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSupplier(t *testing.T) {
	db, token, _ := setupTestEnvironment()
	r := SetupRouter(db)

	t.Run("CreateSupplierSuccess", func(t *testing.T) {
		supplier := model.Supplier{
			Name:        "Test Supplier",
			Description: "Test Description",
			Email:       "supplier@example.com",
			Contact:     "1234567890",
			AccountID:   1,
		}
		jsonValue, _ := json.Marshal(supplier)
		req, _ := http.NewRequest("POST", "/suppliers", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Supplier created successfully", response.Message)
	})

	// Clean up the database
	db.Exec("DELETE FROM suppliers")
}

func TestUpdateSupplier(t *testing.T) {
	db, token, testUser := setupTestEnvironment()
	r := SetupRouter(db)

	supplier := model.Supplier{
		Name:        "Test Supplier",
		Description: "Test Description",
		Email:       "supplier@example.com",
		Contact:     "1234567890",
		AccountID:   testUser.AccountID,
	}
	db.Create(&supplier)

	t.Run("UpdateSupplierSuccess", func(t *testing.T) {
		updatedSupplier := model.Supplier{
			Name:        "Updated Supplier",
			Description: "Updated Description",
			Email:       "updated_supplier@example.com",
			Contact:     "0987654321",
		}
		jsonValue, _ := json.Marshal(updatedSupplier)
		req, _ := http.NewRequest("PUT", "/suppliers/"+strconv.Itoa(int(supplier.ID)), bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Supplier updated successfully", response.Message)
	})

	// Clean up the database
	db.Exec("DELETE FROM suppliers")
}

func TestGetSuppliers(t *testing.T) {
	db, token, testUser := setupTestEnvironment()
	r := SetupRouter(db)

	// Create some test suppliers
	supplier1 := model.Supplier{
		Name:        "Test Supplier 1",
		Description: "Test Description 1",
		Email:       "supplier1@example.com",
		Contact:     "1234567890",
		AccountID:   testUser.AccountID,
	}
	supplier2 := model.Supplier{
		Name:        "Test Supplier 2",
		Description: "Test Description 2",
		Email:       "supplier2@example.com",
		Contact:     "0987654321",
		AccountID:   testUser.AccountID,
	}
	db.Create(&supplier1)
	db.Create(&supplier2)

	t.Run("GetAllSuppliersSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/suppliers", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuppliersResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Suppliers))
	})

	// Clean up the database
	db.Exec("DELETE FROM suppliers")
}

func TestSoftDeleteSupplier(t *testing.T) {
	db, token, testUser := setupTestEnvironment()
	r := SetupRouter(db)

	supplier := model.Supplier{
		Name:        "Test Supplier",
		Description: "Test Description",
		Email:       "supplier@example.com",
		Contact:     "1234567890",
		AccountID:   testUser.AccountID,
	}
	db.Create(&supplier)

	t.Run("SoftDeleteSupplierSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/suppliers/"+strconv.Itoa(int(supplier.ID)), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Supplier deleted successfully", response.Message)
	})

	// Clean up the database
	db.Exec("DELETE FROM suppliers")
}

func TestHardDeleteSupplier(t *testing.T) {
	db, token, testUser := setupTestEnvironment()
	r := SetupRouter(db)

	supplier := model.Supplier{
		Name:        "Test Supplier",
		Description: "Test Description",
		Email:       "supplier@example.com",
		Contact:     "1234567890",
		AccountID:   testUser.AccountID,
	}
	db.Create(&supplier)

	t.Run("HardDeleteSupplierSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/suppliers/hard/"+strconv.Itoa(int(supplier.ID)), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Supplier deleted permanently", response.Message)
	})

	// Clean up the database
	db.Exec("DELETE FROM suppliers")
}

func TestRecoverSupplier(t *testing.T) {
	db, token, testUser := setupTestEnvironment()
	r := SetupRouter(db)

	supplier := model.Supplier{
		Name:        "Test Supplier",
		Description: "Test Description",
		Email:       "supplier@example.com",
		Contact:     "1234567890",
		AccountID:   testUser.AccountID,
	}
	db.Create(&supplier)

	// Soft delete the supplier
	db.Delete(&supplier)

	t.Run("RecoverSupplierSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/suppliers/"+strconv.Itoa(int(supplier.ID))+"/recover", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Supplier recovered successfully", response.Message)
	})

	// Clean up the database
	db.Exec("DELETE FROM suppliers")
}

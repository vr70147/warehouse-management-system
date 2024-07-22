package handlers_test

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

func TestCreateProduct(t *testing.T) {
	db, token, _ := setupTestEnvironment()
	r := SetupRouter(db)

	t.Run("CreateProductSuccess", func(t *testing.T) {
		product := model.Product{
			Name:      "Test Product",
			AccountID: 1,
		}
		jsonValue, _ := json.Marshal(product)
		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.Product
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, product.Name, response.Name)
		assert.Equal(t, product.AccountID, response.AccountID)
	})

	// Clean up the database
	db.Exec("DELETE FROM products")
}

func TestGetProducts(t *testing.T) {
	db, token, testUser := setupTestEnvironment()
	r := SetupRouter(db)

	// Create some test products
	product1 := model.Product{
		Name:      "Test Product 1",
		AccountID: testUser.AccountID,
	}
	product2 := model.Product{
		Name:      "Test Product 2",
		AccountID: testUser.AccountID,
	}
	db.Create(&product1)
	db.Create(&product2)

	t.Run("GetAllProductsSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []model.Product
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response))
	})

	t.Run("GetProductByIDSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products?id="+strconv.Itoa(int(product1.ID)), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.Product
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, product1.Name, response.Name)
		assert.Equal(t, product1.AccountID, response.AccountID)
	})

	t.Run("GetProductByNameSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products?name=Test Product 1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []model.Product
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response))
		assert.Equal(t, product1.Name, response[0].Name)
	})

	// Clean up the database
	db.Exec("DELETE FROM products")
}

func TestUpdateProduct(t *testing.T) {
	db, token, testUser := setupTestEnvironment()
	r := SetupRouter(db)

	product := model.Product{
		Name:      "Test Product",
		AccountID: testUser.AccountID,
	}
	db.Create(&product)

	t.Run("UpdateProductSuccess", func(t *testing.T) {
		updatedProduct := model.Product{
			Name: "Updated Product",
		}
		jsonValue, _ := json.Marshal(updatedProduct)
		req, _ := http.NewRequest("PUT", "/products/"+strconv.Itoa(int(product.ID)), bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.Product
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, updatedProduct.Name, response.Name)
	})

	// Clean up the database
	db.Exec("DELETE FROM products")
}

func TestSoftDeleteProduct(t *testing.T) {
	db, token, testUser := setupTestEnvironment()
	r := SetupRouter(db)

	product := model.Product{
		Name:      "Test Product",
		AccountID: testUser.AccountID,
	}
	db.Create(&product)

	t.Run("SoftDeleteProductSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/products/"+strconv.Itoa(int(product.ID)), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Product soft deleted", response.Message)
	})

	// Clean up the database
	db.Exec("DELETE FROM products")
}

func TestHardDeleteProduct(t *testing.T) {
	db, token, testUser := setupTestEnvironment()
	r := SetupRouter(db)

	product := model.Product{
		Name:      "Test Product",
		AccountID: testUser.AccountID,
	}
	db.Create(&product)

	t.Run("HardDeleteProductSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/products/hard/"+strconv.Itoa(int(product.ID)), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Product and associated stocks hard deleted", response.Message)
	})

	// Clean up the database
	db.Exec("DELETE FROM products")
}

func TestRecoverProduct(t *testing.T) {
	db, token, testUser := setupTestEnvironment()
	r := SetupRouter(db)

	product := model.Product{
		Name:      "Test Product",
		AccountID: testUser.AccountID,
	}
	db.Create(&product)

	// Soft delete the product
	db.Delete(&product)

	t.Run("RecoverProductSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/products/"+strconv.Itoa(int(product.ID))+"/recover", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Product recovered", response.Message)
	})

	// Clean up the database
	db.Exec("DELETE FROM products")
}

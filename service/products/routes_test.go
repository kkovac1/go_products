package products

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kkovac1/products/types"
)

type mockProductsStore struct{}

func TestProductsServiceHandlers(t *testing.T) {
	productsStore := &mockProductsStore{}
	handler := NewHandler(productsStore)

	t.Run(("Should create product"), func(t *testing.T) {
		payload := types.CreateProductPayload{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       10.0,
			Quantity:    5,
		}

		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products", handler.handleCreateProduct)

		router.ServeHTTP(rr, req)

		if condition := rr.Code != http.StatusCreated; condition {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

func (m *mockProductsStore) GetProductById(productId int) (*types.Product, error) {
	return &types.Product{}, nil
}

func (m *mockProductsStore) GetProductsByIds(productIds []int) ([]types.Product, error) {
	return []types.Product{}, nil
}

func (m *mockProductsStore) GetAllProducts() ([]*types.Product, error) {
	return []*types.Product{}, nil
}

func (m *mockProductsStore) CreateProduct(product types.Product) error {
	return nil
}

func (m *mockProductsStore) UpdateProduct(product types.Product) error {
	return nil
}

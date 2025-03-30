package cart

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/kkovac1/products/service/auth"
	"github.com/kkovac1/products/types"
	"github.com/kkovac1/products/utils"
)

type Handler struct {
	store         types.OrderStore
	productsStore types.ProductsStore
	userStore     types.UserStore
}

func NewHandler(store types.OrderStore, productsStore types.ProductsStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, productsStore: productsStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	// Handle checkout logic here
	var payload types.CartCheckoutPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	productIds, err := getCartItemsIds(payload.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Check if the products exist and are in stock
	products, err := h.productsStore.GetProductsByIds(productIds)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	UserId := 0
	orderId, totalPrice, err := h.CreateOrder(products, payload.Items, UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"OrderId":    orderId,
		"TotalPrice": totalPrice,
	})

}

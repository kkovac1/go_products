package cart

import (
	"fmt"

	"github.com/kkovac1/products/types"
)

func getCartItemsIds(items []types.CartItem) ([]int, error) {
	productIds := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", item.ProductId)
		}

		productIds[i] = item.ProductId
	}

	return productIds, nil
}

func (h *Handler) CreateOrder(products []types.Product, cartItems []types.CartItem, userId int) (int, float64, error) {
	//Check if all products are in stock
	productsMap := make(map[int]types.Product)
	for _, product := range products {
		productsMap[product.ID] = product
	}

	err := CheckIfProductInStock(productsMap, cartItems)
	if err != nil {
		return 0, 0, err
	}

	// Calculate total price
	totalPrice := CalculateTotalPrice(productsMap, cartItems)

	// Reduce quntity of products in stock
	for _, item := range cartItems {
		product := productsMap[item.ProductId]
		product.Quantity -= item.Quantity
		h.productsStore.UpdateProduct(productsMap[item.ProductId])
	}

	// Create order
	orderId, err := h.store.CreateOrder(types.Order{
		UserID:  userId,
		Total:   totalPrice,
		Status:  "pending",
		Address: "123 Main St",
	})
	if err != nil {
		return 0, 0, err
	}

	// Create order items
	for _, item := range cartItems {
		err := h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderId,
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
			Price:     productsMap[item.ProductId].Price,
		})
		if err != nil {
			return 0, 0, err
		}
	}

	return orderId, totalPrice, nil
}

func CheckIfProductInStock(products map[int]types.Product, cartItems []types.CartItem) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductId]
		if !ok {
			return fmt.Errorf("product %d not found", item.ProductId)
		}
		if product.Quantity < item.Quantity {
			return fmt.Errorf("not enough stock for product %d", item.ProductId)
		}
	}

	return nil
}

func CalculateTotalPrice(products map[int]types.Product, cartItems []types.CartItem) float64 {
	var totalPrice float64
	for _, item := range cartItems {
		product := products[item.ProductId]
		price := product.Price * float64(item.Quantity)
		totalPrice += price
	}
	return totalPrice
}

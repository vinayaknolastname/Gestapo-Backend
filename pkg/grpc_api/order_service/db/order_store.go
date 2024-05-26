package db

import (
	"context"
	"fmt"
	"time"

	"github.com/akmal4410/gestapo/internal/database"
	"github.com/akmal4410/gestapo/pkg/grpc_api/order_service/db/entity"
	user_entity "github.com/akmal4410/gestapo/pkg/grpc_api/user_service/db/entity"
	"github.com/lib/pq"

	"github.com/akmal4410/gestapo/pkg/utils"
	"github.com/google/uuid"
)

type OrderStore struct {
	storage *database.Storage
}

func NewOrderStore(storage *database.Storage) *OrderStore {
	return &OrderStore{storage: storage}
}

// returns true if the user has order more than two time
func (store *OrderStore) CheckCODIsAvailable(UserID string) (bool, error) {
	selectQuery := `SELECT COUNT(user_id) FROM order_details WHERE user_id = $1;`
	var count int
	err := store.storage.DB.QueryRow(selectQuery, UserID).Scan(&count)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return false, err
	}
	return count > 2, nil
}

func (store *OrderStore) CreateOrder(req *entity.CreateOrderReq) error {
	createdAt := time.Now()
	updatedAt := time.Now()

	ctx := context.Background()
	tx, err := store.storage.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	paymentID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	var status string
	if req.PaymentMode == utils.COD {
		status = utils.PaymentPending
	} else {
		status = utils.PaymentCompleted
	}

	insertPaymentQuery := `
	INSERT INTO payment_details
	(id, amount, provider, status, transaction_id,  created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	_, err = tx.Exec(insertPaymentQuery, paymentID, req.Amount, req.PaymentMode, status, req.TransactionID, createdAt, updatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	orderDetailID, err := uuid.NewRandom()
	if err != nil {
		tx.Rollback()
		return err
	}

	insertOrderDetailQuery := `
	INSERT INTO order_details
	(id, user_id, payment_id, address_id, promo_id, amount, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	_, err = tx.Exec(insertOrderDetailQuery, orderDetailID, req.UserID, paymentID, req.AddressID, req.PromoID, req.Amount, createdAt, updatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	//select all the cart items
	selectOrderItemsQuery := `
	SELECT product_id, inventory_id, quantity, price 
	FROM cart_items 
	WHERE cart_id = $1;
	`
	rows, err := store.storage.DB.Query(selectOrderItemsQuery, req.CartID)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	var cartItems []*user_entity.CartItemRes
	for rows.Next() {
		var item user_entity.CartItemRes

		err := rows.Scan(
			&item.ProductID,
			&item.InventoryID,
			&item.Quantity,
			&item.Price,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
		cartItems = append(cartItems, &item)
	}

	err = rows.Err()
	if err != nil {
		tx.Rollback()
		return err
	}

	var discountedPercent *float64
	if req.PromoID != nil {
		selectQuery := `SELECT percent FROM promo_codes WHERE id = $1;`
		err = tx.QueryRow(selectQuery, req.PromoID).Scan(&discountedPercent)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, item := range cartItems {
		amount := item.Price
		if discountedPercent != nil {
			amount = amount * (1 - *discountedPercent/100)
		}

		var size float32
		selectSizeQuery := `SELECT size FROM inventories WHERE id = $1;`
		err = tx.QueryRow(selectSizeQuery, item.InventoryID).Scan(&size)
		if err != nil {
			tx.Rollback()
			return err
		}

		orderItemID, err := uuid.NewRandom()
		if err != nil {
			tx.Rollback()
			return err
		}

		insertOrderItemQuery := `
		INSERT INTO order_items
		(id, order_id, product_id, size, quantity, amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
		`

		_, err = tx.Exec(insertOrderItemQuery, orderItemID, orderDetailID, item.ProductID, size, item.Quantity, amount, utils.OrderActive, createdAt, updatedAt)
		if err != nil {
			tx.Rollback()
			return err
		}

		//Inserting into tracking_details table
		trackingID, err := uuid.NewRandom()
		if err != nil {
			tx.Rollback()
			return err
		}

		insertTrackingQuery := `
		INSERT INTO tracking_details
		(id, order_item_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5);
		`

		_, err = tx.Exec(insertTrackingQuery, trackingID, orderItemID, utils.TrackingStatus0, createdAt, updatedAt)
		if err != nil {
			tx.Rollback()
			return err
		}

		//Inserting into tracking_items table
		trackingItemID, err := uuid.NewRandom()
		if err != nil {
			tx.Rollback()
			return err
		}

		insertTrackingItmeQuery := `
		INSERT INTO tracking_items
		(id, tracking_id, title, summary, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6);
		`

		_, err = tx.Exec(insertTrackingItmeQuery, trackingItemID, trackingID, utils.TrackingTitles[utils.TrackingStatus0], utils.TrackingSummeries[utils.TrackingStatus0], createdAt, updatedAt)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Update quantity in inventories in table
		updateQuery := `
        UPDATE inventories
        SET quantity = quantity - $1, updated_at = $2
        WHERE id = $3;
    	`
		res, err := tx.Exec(updateQuery, item.Quantity, updatedAt, item.InventoryID)
		if err != nil {
			tx.Rollback()
			return err
		}
		n, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}
		if n == 0 {
			tx.Rollback()
			return fmt.Errorf("could update inventories")
		}
	}

	//Deleting the cart_items
	deleteCartItemsQuery := `DELETE FROM cart_items WHERE cart_id = $1;`

	res, err := store.storage.DB.Exec(deleteCartItemsQuery, req.CartID)
	if err != nil {
		tx.Rollback()
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if n == 0 {
		tx.Rollback()
		return fmt.Errorf("could not clear the cart items")
	}

	//Deleting the cart_items
	deleteCartQuery := `DELETE FROM carts WHERE id = $1;`

	res, err = store.storage.DB.Exec(deleteCartQuery, req.CartID)
	if err != nil {
		tx.Rollback()
		return err
	}
	n, err = res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if n == 0 {
		tx.Rollback()
		return fmt.Errorf("could not clear the cart")
	}

	tx.Commit()
	return nil
}

func (store *OrderStore) GetUserOrders(userID, status string) ([]*entity.UserOrderRes, error) {
	selectQuery := `
	SELECT
    oi.id AS id,
	p.id AS product_id,
    p.product_name AS product_name,
	p.images AS product_images,
    oi.size AS size,
    oi.amount AS price,
	oi.status AS status
	FROM
    order_items oi
	LEFT JOIN
    products p ON oi.product_id = p.id
	LEFT JOIN
    order_details o ON o.user_id = $1
	WHERE 
	oi.status = $2;
	`

	rows, err := store.storage.DB.Query(selectQuery, userID, status)
	if err != nil {
		return nil, err
	}

	var orders []*entity.UserOrderRes

	defer rows.Close()
	for rows.Next() {
		var order entity.UserOrderRes
		var images pq.StringArray

		err := rows.Scan(
			&order.ID,
			&order.ProductID,
			&order.ProductName,
			&images,
			&order.Size,
			&order.Price,
			&order.Status,
		)
		if err != nil {
			return nil, err
		}
		order.ProductImage = []string(images)[0]
		orders = append(orders, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil

}

func (store *OrderStore) GetMerchantOrders(merchantID, status string) ([]*entity.UserOrderRes, error) {
	selectQuery := `
	SELECT
    oi.id AS id,
    p.product_name AS product_name,
	p.images AS product_images,
    oi.size AS size,
    oi.amount AS price,
	oi.status AS status
	FROM
    order_items oi
	LEFT JOIN
    products p ON oi.product_id = p.id
	WHERE 
	p.merchent_id = $1 AND oi.status = $2;
	`

	rows, err := store.storage.DB.Query(selectQuery, merchantID, status)
	if err != nil {
		return nil, err
	}

	var orders []*entity.UserOrderRes

	defer rows.Close()
	for rows.Next() {
		var order entity.UserOrderRes
		var images pq.StringArray

		err := rows.Scan(
			&order.ID,
			&order.ProductName,
			&images,
			&order.Size,
			&order.Price,
			&order.Status,
		)
		if err != nil {
			return nil, err
		}
		order.ProductImage = []string(images)[0]
		orders = append(orders, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil

}

func (store *OrderStore) IsMerchantCanUpdate(orderItemID, merchantID string) (bool, error) {
	selectQuery := `SELECT COUNT(oi.id) 
	FROM order_items oi
	JOIN products p ON oi.product_id = p.id
	WHERE oi.id = $1 AND p.merchent_id = $2;
	`
	var count int
	err := store.storage.DB.QueryRow(selectQuery, orderItemID, merchantID).Scan(&count)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return false, err
	}
	return count > 0, nil
}

func (store *OrderStore) GetMerchantTrackingStatus(orderItemID string) (int, error) {
	selectQuery := `SELECT status FROM tracking_details WHERE order_item_id = $1;`
	var status int
	err := store.storage.DB.QueryRow(selectQuery, orderItemID).Scan(&status)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return 0, err
	}

	return status, nil
}

func (store *OrderStore) UpdateOrderStatus(orderItemID string) error {
	ctx := context.Background()
	tx, err := store.storage.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	createdAt := time.Now()
	updatedAt := time.Now()

	updateTrackingDetailsQuery := `
	UPDATE tracking_details
	SET status = LEAST(status + 1, 3), updated_at = $2
	WHERE order_item_id = $1;
	`
	res, err := tx.Exec(updateTrackingDetailsQuery, orderItemID, updatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if n == 0 {
		tx.Rollback()
		return fmt.Errorf("couldnot update the tracking_details")
	}

	selectQuery := `SELECT id, status FROM tracking_details WHERE order_item_id = $1;`
	var trackingID string
	var status int
	err = store.storage.DB.QueryRow(selectQuery, orderItemID).Scan(&trackingID, &status)
	if err != nil {
		tx.Rollback()
		return err
	}

	status = status + 1 //beacuse only after transaction is completed it will update the status

	// Inserting into tracking_items table
	trackingItemID, err := uuid.NewRandom()
	if err != nil {
		tx.Rollback()
		return err
	}

	insertTrackingItmeQuery := `
	INSERT INTO tracking_items
	(id, tracking_id, title, summary, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err = tx.Exec(insertTrackingItmeQuery, trackingItemID, trackingID, utils.TrackingTitles[status], utils.TrackingSummeries[status], createdAt, updatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	n, err = res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if n == 0 {
		tx.Rollback()
		return fmt.Errorf("couldnot insert into tracking_items")
	}

	if status >= 3 {
		updateOrderItemQuery := `
		UPDATE order_items
		SET status = $2, updated_at = $3
		WHERE id = $1;
		`
		res, err := tx.Exec(updateOrderItemQuery, orderItemID, utils.OrderCompleted, updatedAt)
		if err != nil {
			tx.Rollback()
			return err
		}

		n, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}
		if n == 0 {
			tx.Rollback()
			return fmt.Errorf("couldnot update the order_items")
		}
	}

	tx.Commit()
	return nil
}

func (store *OrderStore) GetOrderTrackingDetails(orderItemId string) ([]*entity.TrackingDetailsRes, error) {
	var details []*entity.TrackingDetailsRes
	selectQuery := `
	SELECT 
	t.status AS status,
	ti.title AS title,
	ti.summary AS summary,
	ti.updated_at AS time
	FROM tracking_details t
	LEFT JOIN 
	tracking_items ti ON t.id = ti.tracking_id
	WHERE t.order_item_id = $1;
	`
	rows, err := store.storage.DB.Query(selectQuery, orderItemId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var address entity.TrackingDetailsRes
		err := rows.Scan(
			&address.Status,
			&address.Title,
			&address.Summary,
			&address.Time,
		)
		if err != nil {
			return nil, err
		}

		details = append(details, &address)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return details, nil
}

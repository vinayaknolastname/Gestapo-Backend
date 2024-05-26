package db

import (
	"fmt"
	"time"

	"github.com/akmal4410/gestapo/internal/database"
	"github.com/akmal4410/gestapo/pkg/grpc_api/product_service/db/entity"
	"github.com/akmal4410/gestapo/pkg/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ProductStore struct {
	storage *database.Storage
}

func NewProductStore(storage *database.Storage) *ProductStore {
	return &ProductStore{storage: storage}

}
func (store ProductStore) CheckDataExist(table, column, value string) (bool, error) {
	checkQuery := fmt.Sprintf(`SELECT * FROM %s WHERE %s = $1;`, table, column)
	res, err := store.storage.DB.Exec(checkQuery, value)
	if err != nil {
		return false, err
	}

	result, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return result != 0, nil
}

func (store *ProductStore) GetProductsForUser(merchantId *string, userId string) ([]*entity.GetProductRes, error) {
	var products []*entity.GetProductRes
	selectQuery := `
	SELECT 
    p.id AS product_id, 
    p.product_name AS product_name, 
    p.images AS product_images, 
    p.price AS product_price,
    AVG(r.star) AS star,
    w.id AS wishlist_id
	FROM 
    products p
	LEFT JOIN 
    reviews r ON p.id = r.product_id
	LEFT JOIN 
    wishlists w ON p.id = w.product_id 
    WHERE p.merchent_id = COALESCE($1, p.merchent_id) AND (w.user_id = $2 OR w.user_id IS NULL)
	GROUP BY 
    p.id, p.product_name, p.images, p.price, w.id;
    `

	rows, err := store.storage.DB.Query(selectQuery, merchantId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.GetProductRes
		var images pq.StringArray

		err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&images,
			&product.Price,
			&product.ReviewStar,
			&product.WishlistID,
		)
		if err != nil {
			return nil, err
		}
		product.ProductImages = []string(images)
		products = append(products, &product)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return []*entity.GetProductRes{}, nil
	}

	return products, nil
}

func (store *ProductStore) GetProductsForMerchants(merchantId *string) ([]*entity.GetProductRes, error) {
	var products []*entity.GetProductRes
	selectQuery := `
	SELECT id, product_name, images, price
	FROM products
	WHERE merchent_id = COALESCE($1, merchent_id);
    `

	rows, err := store.storage.DB.Query(selectQuery, merchantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.GetProductRes
		var images pq.StringArray

		err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&images,
			&product.Price,
		)
		if err != nil {
			return nil, err
		}
		product.ProductImages = []string(images)
		products = append(products, &product)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return []*entity.GetProductRes{}, nil
	}

	return products, nil
}

func (store *ProductStore) GetProductByIdForUser(productId, userId string) (*entity.GetProductRes, error) {
	selectQuery := `
	SELECT
    p.id AS id,
	p.merchent_id AS merchent_id,
    p.product_name AS product_name,
    p.description AS description,
    c.category_name AS category_name,
    p.size AS size,
    p.price AS price,
    CASE
        WHEN d.end_time IS NOT NULL AND d.end_time > NOW()
		THEN p.price - (p.price * d.percent / 100) 
        ELSE NULL
    END AS discount_price,
    p.images AS product_images,
	AVG(r.star) AS star,
    w.id AS wishlist_id
	FROM
    products p
	LEFT JOIN
    categories c ON p.category_id = c.id
	LEFT JOIN
    discounts d ON p.discount_id = d.id
	LEFT JOIN 
    reviews r ON p.id = r.product_id
	LEFT JOIN 
    wishlists w ON p.id = w.product_id 
	WHERE 
	p.id = $1 AND (w.user_id = $2 OR w.user_id IS NULL)
	GROUP BY 
    p.id, p.merchent_id, p.product_name, p.description, c.category_name, p.size, p.price, p.images, w.id, d.end_time, d.percent;
	`
	rows := store.storage.DB.QueryRow(selectQuery, productId, userId)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	var product entity.GetProductRes

	var images pq.StringArray
	var sizes pq.Float64Array

	err := rows.Scan(
		&product.ID,
		&product.MerchantID,
		&product.ProductName,
		&product.Description,
		&product.CategoryName,
		&sizes,
		&product.Price,
		&product.DiscountPrice,
		&images,
		&product.ReviewStar,
		&product.WishlistID,
	)
	product.ProductImages = []string(images)
	// Convert pq.Float64Array to []float64
	var sizeList []float64
	for _, v := range sizes {
		sizeList = append(sizeList, float64(v))
	}
	product.Size = &sizeList

	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (store *ProductStore) GetProductByIdForMerchant(productId string) (*entity.GetProductRes, error) {
	selectQuery := `
	SELECT
    p.id AS id,
	p.merchent_id AS merchent_id,
    p.product_name AS product_name,
    p.description AS description,
    c.category_name AS category_name,
    p.size AS size,
    p.price AS price,
    CASE
        WHEN d.end_time IS NOT NULL AND d.end_time > NOW()
		THEN p.price - (p.price * d.percent / 100) 
        ELSE NULL
    END AS discount_price,
    p.images AS product_images
	FROM
    products p
	LEFT JOIN
    categories c ON p.category_id = c.id
	LEFT JOIN
    discounts d ON p.discount_id = d.id
	WHERE 
	p.id = $1;
	`
	rows := store.storage.DB.QueryRow(selectQuery, productId)
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	var product entity.GetProductRes

	var images pq.StringArray
	var sizes pq.Float64Array

	err := rows.Scan(
		&product.ID,
		&product.MerchantID,
		&product.ProductName,
		&product.Description,
		&product.CategoryName,
		&sizes,
		&product.Price,
		&product.DiscountPrice,
		&images,
	)
	product.ProductImages = []string(images)
	// Convert pq.Float64Array to []float64
	var sizeList []float64
	for _, v := range sizes {
		sizeList = append(sizeList, float64(v))
	}
	product.Size = &sizeList

	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (store *ProductStore) IsUserCanAddReview(orderItemID, userID string) (bool, error) {
	selectQuery := `SELECT COUNT(oi.id) 
	FROM order_items oi
	JOIN products p ON oi.product_id = p.id
	JOIN order_details od ON oi.order_id = od.id
	WHERE oi.id = $1 AND od.user_id = $2 AND oi.status = $3;
	`
	var count int
	err := store.storage.DB.QueryRow(selectQuery, orderItemID, userID, utils.OrderCompleted).Scan(&count)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return false, err
	}
	return count > 0, nil
}

func (store *ProductStore) IsUserAlreadyAddedReview(productID, userID string) (bool, error) {
	selectQuery := `SELECT COUNT(*) 
	FROM reviews
	WHERE product_id = $1 AND user_id = $2;
	`
	var count int
	err := store.storage.DB.QueryRow(selectQuery, productID, userID).Scan(&count)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return false, err
	}
	return count > 0, nil
}

func (store *ProductStore) AddProductReview(req *entity.AddReviewReq) error {
	createdAt := time.Now()
	updatedAt := time.Now()

	uuId, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	insertQuery := `
	INSERT INTO reviews (id, product_id, user_id, star, review, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	_, err = store.storage.DB.Exec(insertQuery, uuId, req.ProductID, req.UserID, req.Star, req.Review, createdAt, updatedAt)
	if err != nil {
		return err
	}
	return nil
}

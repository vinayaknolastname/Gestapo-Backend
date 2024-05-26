package db

import (
	"context"
	"fmt"
	"time"

	"github.com/akmal4410/gestapo/internal/database"
	"github.com/akmal4410/gestapo/pkg/grpc_api/merchant_service/db/entity"
	product_entity "github.com/akmal4410/gestapo/pkg/grpc_api/product_service/db/entity"
	user_entity "github.com/akmal4410/gestapo/pkg/grpc_api/user_service/db/entity"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type MerchantStore struct {
	storage *database.Storage
}

func NewMerchantStore(storage *database.Storage) *MerchantStore {
	return &MerchantStore{storage: storage}

}

func (store MerchantStore) CheckDataExist(table, column, value string) (bool, error) {
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

func (store *MerchantStore) GetProfile(userId string) (*entity.GetMerchantRes, error) {
	selectQuery := `
	SELECT id, profile_image, full_name, user_name, phone, email, dob, gender, user_type 
	FROM user_data WHERE id = $1;
	`
	rows := store.storage.DB.QueryRow(selectQuery, userId)
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	var user entity.GetMerchantRes
	err := rows.Scan(
		&user.ID,
		&user.ProfileImage,
		&user.FullName,
		&user.UserName,
		&user.Phone,
		&user.Email,
		&user.DOB,
		&user.Gender,
		&user.UserType,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (store *MerchantStore) UpdateProfile(id string, req *entity.EditMerchantReq) error {
	updatedAt := time.Now()

	updateQuery := `UPDATE user_data
	SET profile_image = $2, full_name = $3, dob = $4, gender = $5, updated_at = $6
	WHERE id = $1;`

	res, err := store.storage.DB.Exec(updateQuery, id, req.ProfileImage, req.FullName, req.DOB, req.Gender, updatedAt)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("couldnot update the user")
	}
	return nil
}

func (store *MerchantStore) InsertProduct(userId, productId string, req *entity.AddProductReq) error {
	createdAt := time.Now()
	updatedAt := time.Now()

	ctx := context.Background()

	tx, err := store.storage.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	insertProductQuery := `
		INSERT INTO products
		(id, merchent_id, category_id, product_name, description, images, size, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
		`
	productImages := pq.StringArray(req.ProductImages)
	productSizes := pq.Float64Array(req.Sizes)

	_, err = tx.Exec(insertProductQuery, productId, userId, req.CategoryId, req.ProductName, req.Description, productImages, productSizes, req.Price, createdAt, updatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, size := range req.Sizes {
		inventoryId, err := uuid.NewRandom()
		if err != nil {
			tx.Rollback()
			return err
		}
		insertInventoryQuery := `
		INSERT INTO inventories
		(id, product_id, size, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6);
		`
		_, err = tx.Exec(insertInventoryQuery, inventoryId.String(), productId, size, req.Quantity, createdAt, updatedAt)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (store *MerchantStore) UpdateProduct(id string, req *entity.EditProductReq) error {
	updateQuery := `
	UPDATE products
	SET product_name = $2, description = $3, images = $4, price = $5, updated_at = $6
	WHERE id = $1;
	`
	updatedAt := time.Now()
	productImages := pq.StringArray(req.ProductImages)

	res, err := store.storage.DB.Exec(updateQuery, id, req.ProductName, req.Description, productImages, req.Price, updatedAt)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("couldnot update the products")
	}
	return nil
}

func (store *MerchantStore) GetProductById(productId string) (*product_entity.GetProductRes, error) {
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

	var product product_entity.GetProductRes

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

func (store *MerchantStore) DeleteProduct(productId string) error {
	deleteQuery := `
        DELETE FROM products
        USING inventories
        WHERE 
		products.id = $1
        AND 
		products.inventory_id = inventories.id;
    `

	res, err := store.storage.DB.Exec(deleteQuery, productId)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("could not delete the product")
	}
	return nil
}

func (store *MerchantStore) AddProductDiscount(req *entity.AddDiscountReq) error {
	ctx := context.Background()
	tx, err := store.storage.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	discountId, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	createdAt := time.Now()
	updatedAt := time.Now()

	if req.CardColor == "" {
		req.CardColor = "0xFF808080"
	}

	insertQuery := `
	INSERT INTO discounts
	(id, merchent_id, name, description, percent, card_color, start_time, end_time, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`

	res, err := tx.Exec(insertQuery, discountId.String(), req.MerchantId, req.DiscountName, req.Description, req.Percentage, req.CardColor, req.StartTime, req.EndTime, createdAt, updatedAt)
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
		return fmt.Errorf("couldnot insert the discounts")
	}

	updateQuery := `UPDATE products
	SET discount_id = $2, updated_at = $3 
	WHERE id = $1;
	`
	res, err = tx.Exec(updateQuery, req.ProductId, discountId.String(), updatedAt)
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
		return fmt.Errorf("couldnot update the products")
	}
	tx.Commit()
	return nil
}

func (store *MerchantStore) EditProductDiscount(discountId string, req *entity.EditDiscountReq) error {
	updateQuery := `
	UPDATE discounts
	SET name = COALESCE($2, name),
		description = COALESCE($3, description),
		percent = COALESCE($4, percent),
		card_color = COALESCE($5, card_color),
		start_time = COALESCE($6, start_time),
		end_time = COALESCE($7, end_time),
		updated_at = $8
	WHERE id = $1;
	`
	updatedAt := time.Now()
	res, err := store.storage.DB.Exec(updateQuery, discountId, req.DiscountName, req.Description, req.Percentage, req.CardColor, req.StartTime, req.EndTime, updatedAt)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("couldnot update the discount")
	}
	return nil
}

func (store *MerchantStore) GetAllDiscount(merchantId *string) ([]*user_entity.DiscountRes, error) {
	selectQuery := `
	SELECT 
    p.id AS product_id,
    d.name AS name,
	d.description AS description,
    d.percent AS percent,
    p.images[1] AS image
	FROM products p
	JOIN discounts d ON p.discount_id = d.id
	WHERE  d.end_time > NOW() AND d.merchent_id = COALESCE($1, d.merchent_id)
	ORDER BY d.percent DESC;
	`

	rows, err := store.storage.DB.Query(selectQuery, merchantId)
	if err != nil {
		return nil, err
	}
	var discounts []*user_entity.DiscountRes
	defer rows.Close()
	for rows.Next() {
		var discount user_entity.DiscountRes

		err := rows.Scan(
			&discount.ProductID,
			&discount.Name,
			&discount.Description,
			&discount.Percentage,
			&discount.ProductImage,
		)
		if err != nil {
			return nil, err
		}
		discounts = append(discounts, &discount)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return discounts, nil
}

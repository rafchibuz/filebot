package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"wildberries-order-service/internal/models"
)

// PostgresOrderRepository реализует хранение заказов в PostgreSQL
type PostgresOrderRepository struct {
	db *pgxpool.Pool
}

// NewPostgresOrderRepository создает новый PostgreSQL репозиторий
func NewPostgresOrderRepository(connectionString string) (*PostgresOrderRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Настройки пула соединений
	config.MaxConns = 30
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Successfully connected to PostgreSQL")

	return &PostgresOrderRepository{
		db: db,
	}, nil
}

// Close закрывает соединение с базой данных
func (r *PostgresOrderRepository) Close() {
	if r.db != nil {
		r.db.Close()
	}
}

// GetByID возвращает заказ по ID
func (r *PostgresOrderRepository) GetByID(orderID string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Получаем основную информацию о заказе
	orderQuery := `
		SELECT order_uid, track_number, entry, locale, internal_signature, 
		       customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders 
		WHERE order_uid = $1`

	var order models.Order
	err = tx.QueryRow(ctx, orderQuery, orderID).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Получаем информацию о доставке
	deliveryQuery := `
		SELECT name, phone, zip, city, address, region, email 
		FROM deliveries 
		WHERE order_uid = $1`

	err = tx.QueryRow(ctx, deliveryQuery, orderID).Scan(
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
		&order.Delivery.Email,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to get delivery info: %w", err)
	}

	// Получаем информацию о платеже
	paymentQuery := `
		SELECT transaction, request_id, currency, provider, amount, payment_dt, 
		       bank, delivery_cost, goods_total, custom_fee 
		FROM payments 
		WHERE order_uid = $1`

	err = tx.QueryRow(ctx, paymentQuery, orderID).Scan(
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to get payment info: %w", err)
	}

	// Получаем товары
	itemsQuery := `
		SELECT chrt_id, track_number, price, rid, name, sale, size, 
		       total_price, nm_id, brand, status 
		FROM items 
		WHERE order_uid = $1 
		ORDER BY id`

	rows, err := tx.Query(ctx, itemsQuery, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID,
			&item.Name, &item.Sale, &item.Size, &item.TotalPrice,
			&item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error reading items: %w", rows.Err())
	}

	order.Items = items

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &order, nil
}

// Save сохраняет заказ в базу данных
func (r *PostgresOrderRepository) Save(order *models.Order) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	if order.OrderUID == "" {
		return errors.New("order_uid is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Вставляем основную информацию о заказе
	orderQuery := `
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature, 
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO UPDATE SET
			track_number = EXCLUDED.track_number,
			entry = EXCLUDED.entry,
			locale = EXCLUDED.locale,
			internal_signature = EXCLUDED.internal_signature,
			customer_id = EXCLUDED.customer_id,
			delivery_service = EXCLUDED.delivery_service,
			shardkey = EXCLUDED.shardkey,
			sm_id = EXCLUDED.sm_id,
			date_created = EXCLUDED.date_created,
			oof_shard = EXCLUDED.oof_shard`

	_, err = tx.Exec(ctx, orderQuery,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	// Сохраняем информацию о доставке
	deliveryQuery := `
		INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (order_uid) DO UPDATE SET
			name = EXCLUDED.name,
			phone = EXCLUDED.phone,
			zip = EXCLUDED.zip,
			city = EXCLUDED.city,
			address = EXCLUDED.address,
			region = EXCLUDED.region,
			email = EXCLUDED.email`

	_, err = tx.Exec(ctx, deliveryQuery,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone,
		order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to save delivery: %w", err)
	}

	// Сохраняем информацию о платеже
	paymentQuery := `
		INSERT INTO payments (
			order_uid, transaction, request_id, currency, provider, amount, 
			payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO UPDATE SET
			transaction = EXCLUDED.transaction,
			request_id = EXCLUDED.request_id,
			currency = EXCLUDED.currency,
			provider = EXCLUDED.provider,
			amount = EXCLUDED.amount,
			payment_dt = EXCLUDED.payment_dt,
			bank = EXCLUDED.bank,
			delivery_cost = EXCLUDED.delivery_cost,
			goods_total = EXCLUDED.goods_total,
			custom_fee = EXCLUDED.custom_fee`

	_, err = tx.Exec(ctx, paymentQuery,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee,
	)
	if err != nil {
		return fmt.Errorf("failed to save payment: %w", err)
	}

	// Удаляем старые товары и вставляем новые
	_, err = tx.Exec(ctx, "DELETE FROM items WHERE order_uid = $1", order.OrderUID)
	if err != nil {
		return fmt.Errorf("failed to delete old items: %w", err)
	}

	// Вставляем товары
	for _, item := range order.Items {
		itemQuery := `
			INSERT INTO items (
				order_uid, chrt_id, track_number, price, rid, name, sale, 
				size, total_price, nm_id, brand, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

		_, err = tx.Exec(ctx, itemQuery,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price,
			item.RID, item.Name, item.Sale, item.Size, item.TotalPrice,
			item.NmID, item.Brand, item.Status,
		)
		if err != nil {
			return fmt.Errorf("failed to save item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Order %s saved to PostgreSQL", order.OrderUID)
	return nil
}

// GetAll возвращает все заказы
func (r *PostgresOrderRepository) GetAll() ([]*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Получаем список всех order_uid
	orderUIDs := []string{}
	rows, err := r.db.Query(ctx, "SELECT order_uid FROM orders ORDER BY date_created DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to get order UIDs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var orderUID string
		if err := rows.Scan(&orderUID); err != nil {
			return nil, fmt.Errorf("failed to scan order UID: %w", err)
		}
		orderUIDs = append(orderUIDs, orderUID)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error reading order UIDs: %w", rows.Err())
	}

	// Получаем полную информацию о каждом заказе
	orders := make([]*models.Order, 0, len(orderUIDs))
	for _, orderUID := range orderUIDs {
		order, err := r.GetByID(orderUID)
		if err != nil {
			log.Printf("Warning: failed to get order %s: %v", orderUID, err)
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}
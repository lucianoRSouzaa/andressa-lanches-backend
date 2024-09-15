package repository

import (
	"andressa-lanches/internal/domain/product"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	Pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{Pool: pool}
}

func (r *ProductRepository) Create(ctx context.Context, p *product.Product) error {
	query := `
        INSERT INTO products (name, price, description, category_id)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	err := r.Pool.QueryRow(ctx, query, p.Name, p.Price, p.Description, p.CategoryID).Scan(&p.ID)
	return err
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	query := `
		SELECT id, name, price, description, category_id
		FROM products
		WHERE id = $1
	`
	row := r.Pool.QueryRow(ctx, query, id)

	var p product.Product
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.CategoryID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *product.Product) error {
	query := `
		UPDATE products
		SET name = $1, price = $2, description = $3, category_id = $4
		WHERE id = $5
	`
	_, err := r.Pool.Exec(ctx, query, product.Name, product.Price, product.Description, product.CategoryID, product.ID)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
        DELETE FROM products
        WHERE id = $1
    `
	_, err := r.Pool.Exec(ctx, query, id)
	return err
}

func (r *ProductRepository) List(ctx context.Context) ([]*product.Product, error) {
	query := `
		SELECT id, name, price, description, category_id
		FROM products
	`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*product.Product
	for rows.Next() {
		var p product.Product
		err = rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.CategoryID)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

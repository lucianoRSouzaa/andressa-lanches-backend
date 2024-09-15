package repository

import (
	"andressa-lanches/internal/domain/category"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	Pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{Pool: pool}
}

func (r *CategoryRepository) Create(ctx context.Context, c *category.Category) error {
	query := `
        INSERT INTO categories (name, description)
        VALUES ($1, $2)
        RETURNING id
    `
	err := r.Pool.QueryRow(ctx, query, c.Name, c.Description).Scan(&c.ID)
	return err
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*category.Category, error) {
	query := `
        SELECT id, name, description
        FROM categories
        WHERE id = $1
    `
	row := r.Pool.QueryRow(ctx, query, id)

	var c category.Category
	err := row.Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Update(ctx context.Context, c *category.Category) error {
	query := `
        UPDATE categories
        SET name = $1, description = $2
        WHERE id = $3
    `
	_, err := r.Pool.Exec(ctx, query, c.Name, c.Description, c.ID)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
        DELETE FROM categories
        WHERE id = $1
    `
	_, err := r.Pool.Exec(ctx, query, id)
	return err
}

func (r *CategoryRepository) List(ctx context.Context) ([]*category.Category, error) {
	query := `
        SELECT id, name, description
        FROM categories
    `
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*category.Category
	for rows.Next() {
		var c category.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

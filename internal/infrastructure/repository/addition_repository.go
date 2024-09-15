package repository

import (
	"andressa-lanches/internal/domain/addition"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdditionRepository struct {
	Pool *pgxpool.Pool
}

func NewAdditionRepository(pool *pgxpool.Pool) *AdditionRepository {
	return &AdditionRepository{Pool: pool}
}

func (r *AdditionRepository) Create(ctx context.Context, a *addition.Addition) error {
	query := `
        INSERT INTO additions (name, price)
        VALUES ($1, $2)
        RETURNING id
    `
	err := r.Pool.QueryRow(ctx, query, a.Name, a.Price).Scan(&a.ID)
	return err
}

func (r *AdditionRepository) GetByID(ctx context.Context, id uuid.UUID) (*addition.Addition, error) {
	query := `
        SELECT id, name, price
        FROM additions
        WHERE id = $1
    `
	a := &addition.Addition{}
	err := r.Pool.QueryRow(ctx, query, id).Scan(&a.ID, &a.Name, &a.Price)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}

func (r *AdditionRepository) Update(ctx context.Context, a *addition.Addition) error {
	query := `
        UPDATE additions
        SET name = $1, price = $2
        WHERE id = $3
    `
	_, err := r.Pool.Exec(ctx, query, a.Name, a.Price, a.ID)
	return err
}

func (r *AdditionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
        DELETE FROM additions
        WHERE id = $1
    `
	_, err := r.Pool.Exec(ctx, query, id)
	return err
}

func (r *AdditionRepository) List(ctx context.Context) ([]*addition.Addition, error) {
	query := `
        SELECT id, name, price
        FROM additions
    `
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var additions []*addition.Addition
	for rows.Next() {
		a := &addition.Addition{}
		err := rows.Scan(&a.ID, &a.Name, &a.Price)
		if err != nil {
			return nil, err
		}
		additions = append(additions, a)
	}
	return additions, nil
}

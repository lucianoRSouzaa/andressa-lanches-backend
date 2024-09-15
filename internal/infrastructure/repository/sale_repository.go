package repository

import (
	"andressa-lanches/internal/domain/addition"
	"andressa-lanches/internal/domain/sale"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SaleRepository struct {
	Pool *pgxpool.Pool
}

func NewSaleRepository(pool *pgxpool.Pool) *SaleRepository {
	return &SaleRepository{Pool: pool}
}

func (r *SaleRepository) Create(ctx context.Context, s *sale.Sale) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	saleQuery := `
        INSERT INTO sales (date, total_amount, discount, additional_charges)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	err = tx.QueryRow(ctx, saleQuery, s.Date, s.TotalAmount, s.Discount, s.AdditionalCharges).Scan(&s.ID)
	if err != nil {
		return err
	}

	saleItemQuery := `
        INSERT INTO sale_items (sale_id, product_id, quantity, unit_price, total_price)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING item_id
    `

	saleItemAdditionQuery := `
        INSERT INTO sale_item_additions (sale_id, item_id, addition_id)
        VALUES ($1, $2, $3)
    `

	for i := range s.Items {
		item := &s.Items[i]
		err = tx.QueryRow(ctx, saleItemQuery, s.ID, item.ProductID, item.Quantity, item.UnitPrice, item.TotalPrice).Scan(&item.ItemID)
		if err != nil {
			return err
		}

		if len(item.Additions) > 0 {
			batch := &pgx.Batch{}
			for _, addition := range item.Additions {
				batch.Queue(saleItemAdditionQuery, s.ID, item.ItemID, addition.ID)
			}
			results := tx.SendBatch(ctx, batch)
			for range item.Additions {
				_, err := results.Exec()
				if err != nil {
					_ = results.Close()
					return err
				}
			}
			err = results.Close()
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit(ctx)
	return err
}

func (r *SaleRepository) GetByID(ctx context.Context, id uuid.UUID) (*sale.Sale, error) {
	saleQuery := `
        SELECT id, date, total_amount, discount, additional_charges
        FROM sales
        WHERE id = $1
    `
	row := r.Pool.QueryRow(ctx, saleQuery, id)

	var s sale.Sale
	err := row.Scan(&s.ID, &s.Date, &s.TotalAmount, &s.Discount, &s.AdditionalCharges)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	saleItemsQuery := `
        SELECT sale_id, item_id, product_id, quantity, unit_price, total_price
        FROM sale_items
        WHERE sale_id = $1
    `
	rows, err := r.Pool.Query(ctx, saleItemsQuery, s.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []sale.SaleItem
	for rows.Next() {
		var item sale.SaleItem
		err := rows.Scan(&item.SaleID, &item.ItemID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.TotalPrice)
		if err != nil {
			return nil, err
		}

		additionsQuery := `
            SELECT a.id, a.name, a.price
            FROM sale_item_additions sia
            INNER JOIN additions a ON sia.addition_id = a.id
            WHERE sia.sale_id = $1 AND sia.item_id = $2
        `
		additionRows, err := r.Pool.Query(ctx, additionsQuery, item.SaleID, item.ItemID)
		if err != nil {
			return nil, err
		}

		var additions []addition.Addition
		for additionRows.Next() {
			var add addition.Addition
			err := additionRows.Scan(&add.ID, &add.Name, &add.Price)
			if err != nil {
				additionRows.Close()
				return nil, err
			}
			additions = append(additions, add)
		}
		additionRows.Close()

		item.Additions = additions
		items = append(items, item)
	}

	s.Items = items
	return &s, nil
}

func (r *SaleRepository) List(ctx context.Context) ([]*sale.Sale, error) {
	salesQuery := `
        SELECT id, date, total_amount, discount, additional_charges
        FROM sales
        ORDER BY date DESC
    `
	salesRows, err := r.Pool.Query(ctx, salesQuery)
	if err != nil {
		return nil, err
	}
	defer salesRows.Close()

	var salesList []*sale.Sale
	for salesRows.Next() {
		var s sale.Sale
		err := salesRows.Scan(&s.ID, &s.Date, &s.TotalAmount, &s.Discount, &s.AdditionalCharges)
		if err != nil {
			return nil, err
		}

		saleItemsQuery := `
            SELECT sale_id, item_id, product_id, quantity, unit_price, total_price
            FROM sale_items
            WHERE sale_id = $1
        `
		itemsRows, err := r.Pool.Query(ctx, saleItemsQuery, s.ID)
		if err != nil {
			return nil, err
		}

		var items []sale.SaleItem
		for itemsRows.Next() {
			var item sale.SaleItem
			err := itemsRows.Scan(&item.SaleID, &item.ItemID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.TotalPrice)
			if err != nil {
				itemsRows.Close()
				return nil, err
			}

			additionsQuery := `
                SELECT a.id, a.name, a.price
                FROM sale_item_additions sia
                INNER JOIN additions a ON sia.addition_id = a.id
                WHERE sia.sale_id = $1 AND sia.item_id = $2
            `
			additionRows, err := r.Pool.Query(ctx, additionsQuery, item.SaleID, item.ItemID)
			if err != nil {
				itemsRows.Close()
				return nil, err
			}

			var additions []addition.Addition
			for additionRows.Next() {
				var add addition.Addition
				err := additionRows.Scan(&add.ID, &add.Name, &add.Price)
				if err != nil {
					additionRows.Close()
					itemsRows.Close()
					return nil, err
				}
				additions = append(additions, add)
			}
			additionRows.Close()

			item.Additions = additions
			items = append(items, item)
		}
		itemsRows.Close()

		s.Items = items
		salesList = append(salesList, &s)
	}

	return salesList, nil
}

func (r *SaleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	deleteAdditionsQuery := `
        DELETE FROM sale_item_additions
        WHERE sale_id = $1
    `
	_, err = tx.Exec(ctx, deleteAdditionsQuery, id)
	if err != nil {
		return err
	}

	deleteItemsQuery := `
        DELETE FROM sale_items
        WHERE sale_id = $1
    `
	_, err = tx.Exec(ctx, deleteItemsQuery, id)
	if err != nil {
		return err
	}

	deleteSaleQuery := `
        DELETE FROM sales
        WHERE id = $1
    `
	result, err := tx.Exec(ctx, deleteSaleQuery, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("venda não encontrada")
	}

	err = tx.Commit(ctx)
	return err
}

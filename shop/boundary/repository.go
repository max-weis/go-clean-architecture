package boundary

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
	"webshop/shop/control"
	"webshop/shop/entity"
)

type (
	sqlxRepository struct {
		db *sqlx.DB
	}

	productEntity struct {
		ID          string    `db:"id"`
		Title       string    `db:"title"`
		Description string    `db:"description"`
		Price       uint64    `db:"price"`
		CreatedAt   time.Time `db:"created_at"`
		ModifiedAt  time.Time `db:"modified_at"`
	}
)

func ProvideSqlxRepository(db *sqlx.DB) control.ProductRepository {
	return sqlxRepository{db: db}
}

const save = `
insert into products(id, title, description, price, created_at, modified_at)
values($1, $2, $3, $4, $5, $6);
`

func (repository sqlxRepository) Save(ctx context.Context, product entity.Product) error {
	tx := repository.db.MustBeginTx(ctx, &sql.TxOptions{})
	p := mapToEntity(product)
	tx.MustExecContext(ctx, save, p.ID, p.Title, p.Description, p.Price, p.CreatedAt, p.ModifiedAt)
	return tx.Commit()
}

func mapToEntity(p entity.Product) *productEntity {
	return &productEntity{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		ModifiedAt:  p.ModifiedAt,
	}
}

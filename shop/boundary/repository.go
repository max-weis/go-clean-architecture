package boundary

import (
	"context"
	"database/sql"
	"fmt"
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

const findPaginated = `
select *
from products
where case when not $1 then price > 0 end
order by $2
limit $3 offset $4;
`

func (repository sqlxRepository) FindPaginated(ctx context.Context, filterObject entity.FilterObject) ([]entity.Product, error) {
	tx := repository.db.MustBeginTx(ctx, &sql.TxOptions{})

	var products []productEntity
	rows, err := tx.Queryx(findPaginated, filterObject.Free, mapOrder(filterObject.Sort), filterObject.Limit, filterObject.Offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var product productEntity
		if err := rows.StructScan(&product); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return mapFromEntities(products), nil
}

const findByID = `
select *
from products
where id = $1;
`

func (repository sqlxRepository) FindByID(ctx context.Context, id string) (entity.Product, error) {
	tx := repository.db.MustBeginTx(ctx, &sql.TxOptions{})

	rows, err := tx.Queryx(findByID, id)
	if err != nil {
		return entity.Product{}, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return entity.Product{}, err
	}

	for rows.Next() {
		var product productEntity
		if err := rows.StructScan(&product); err != nil {
			return entity.Product{}, err
		}

		return mapFromEntity(product), nil
	}

	return entity.Product{}, nil
}

func mapOrder(sort entity.Sorting) string {
	format := "%s %s"

	var attribute string
	switch sort {
	case entity.IdAsc, entity.IdDesc:
		attribute = "id"
	case entity.TitleAsc, entity.TitleDesc:
		attribute = "title"
	case entity.PriceAsc, entity.PriceDesc:
		attribute = "price"
	case entity.CreatedAtAsc, entity.CreatedAtDesc:
		attribute = "created_at"
	case entity.ModifiedAtAsc, entity.ModifiedAtDesc:
		attribute = "modified_at"
	default:
		attribute = "id"
	}

	var order string
	switch sort {
	case entity.IdAsc, entity.TitleAsc, entity.PriceAsc, entity.CreatedAtAsc, entity.ModifiedAtAsc:
		order = "asc"
	case entity.IdDesc, entity.TitleDesc, entity.PriceDesc, entity.CreatedAtDesc, entity.ModifiedAtDesc:
		order = "desc"
	default:
		order = "asc"
	}

	return fmt.Sprintf(format, attribute, order)
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

func mapFromEntities(productEntities []productEntity) []entity.Product {
	var products []entity.Product
	for _, productEntity := range productEntities {
		products = append(products, mapFromEntity(productEntity))
	}

	return products
}

func mapFromEntity(p productEntity) entity.Product {
	return entity.Product{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		ModifiedAt:  p.ModifiedAt,
	}
}

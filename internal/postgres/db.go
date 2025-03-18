package postgres

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"project_sem/internal/model"
)

type Postgres interface {
	Get() ([]*model.Product, error)
	Create(*sql.Tx, *model.Product) error
	Begin() (*sql.Tx, error)
	Close()
	GetUnicCategory(tx *sql.Tx) (int, error)
	GetTotalPrice(tx *sql.Tx) (float64, error)
}

type postgess struct {
	db *sql.DB
}

func NewDB() (Postgres, error) {
	ps := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	db, err := sql.Open("pgx", ps)
	if err != nil {
		return nil, fmt.Errorf("error postgres %w", err)
	}
	return &postgess{db: db}, nil
}
func (p *postgess) Close() {
	p.db.Close()
}

const getSQL = `SELECT id, name,category,price,create_date FROM prices`

func (p *postgess) Get() ([]*model.Product, error) {
	rows, err := p.db.Query(getSQL)
	if err != nil {
		return nil, fmt.Errorf("error postgres %w", err)
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		var product model.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Category, &product.Price, &product.Data)
		if err != nil {
			return nil, fmt.Errorf("error scan postgres %w", err)
		}
		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error rows scan %w", err)
	}
	return products, nil
}

const insertSQL = `INSERT INTO prices ( id, name,category,price,create_date ) 
VALUES ($1,$2,$3,$4,$5 )`

func (p *postgess) Create(tx *sql.Tx, pr *model.Product) error {
	_, err := tx.Exec(insertSQL, pr.ID, pr.Name, pr.Category, pr.Price, pr.Data)
	if err != nil {
		return fmt.Errorf("error postgres %w", err)
	}
	return nil
}

const totalPriceSQL = `SELECT SUM(price) AS TotalPrice FROM prices`

func (p *postgess) GetTotalPrice(tx *sql.Tx) (float64, error) {
	var totalPriceVal float64
	if err := tx.QueryRow(totalPriceSQL).Scan(&totalPriceVal); err != nil {
		return 0, fmt.Errorf("error postgres %w", err)
	}
	return totalPriceVal, nil
}

const totalCountCategorySQL = `SELECT COUNT(DISTINCT category) FROM prices `

func (p *postgess) GetUnicCategory(tx *sql.Tx) (int, error) {
	var totalCountCategory int
	if err := tx.QueryRow(totalCountCategorySQL).Scan(&totalCountCategory); err != nil {
		return 0, fmt.Errorf("error postgres %w", err)
	}
	return totalCountCategory, nil
}

func (p *postgess) Begin() (*sql.Tx, error) {
	return p.db.Begin()
}

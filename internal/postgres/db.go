package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
	"project_sem/internal/model"
)

type Postgres interface {
	Migrate() error
	Get() ([]*model.Product, error)
	Create(*model.Product) error
	Close()
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

const migrateSQL = `CREATE TABLE IF NOT EXISTS prices (id INTEGER,
 name varchar(30),
    category varchar(30),
    price double precision,
    create_date date)`

func (p *postgess) Migrate() error {
	_, err := p.db.Exec(migrateSQL)
	if err != nil {
		return fmt.Errorf("error postgres %w", err)
	}
	return nil
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
	return products, nil
}

const insertSQL = `INSERT INTO prices ( id, name,category,price,create_date ) 
VALUES ($1,$2,$3,$4,$5 )`

func (p *postgess) Create(pr *model.Product) error {
	_, err := p.db.Exec(insertSQL, pr.ID, pr.Name, pr.Category, pr.Price, pr.Data)
	if err != nil {
		return fmt.Errorf("error postgres %w", err)
	}
	return nil
}

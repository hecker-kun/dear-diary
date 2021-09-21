package diary

import (
	"database/sql"
	"log"
	"time"

	"github.com/baryon-m/dear-diary/domain/entity"
)

type PostgreSQLRepo struct {
	db *sql.DB
}

func NewPostgreSQLRepoRepository(db *sql.DB) *PostgreSQLRepo {
	return &PostgreSQLRepo{db: db}
}

func (m *PostgreSQLRepo) Create(e *Entry) (entity.ID, error) {
	query, err := m.db.Prepare("INSERT INTO entries(id, author, content, created_at) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return e.ID, err
	}

	_, err = query.Exec(e.ID, e.Author, e.Content, time.Now().Format(time.RFC3339))
	if err != nil {
		return e.ID, err
	}

	err = query.Close()
	if err != nil {
		return e.ID, err
	}

	return e.ID, nil
}

func (m *PostgreSQLRepo) Delete(id entity.ID) error {
	query, err := m.db.Prepare("DELETE FROM entries WHERE id=$1")
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = query.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgreSQLRepo) FetchOne(id entity.ID) (*Entry, error) {
	return fetchOne(id, m.db)
}

func fetchOne(id entity.ID, db *sql.DB) (*Entry, error) {
	query, err := db.Prepare("SELECT * FROM entries WHERE id = $1")
	if err != nil {
		return nil, err
	}

	var e Entry

	rows, err := query.Query(id)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&e.ID, &e.Author, &e.Content, &e.CreatedAt)
	}
	return &e, nil
}

func (m *PostgreSQLRepo) FetchAll() ([]*Entry, error) {
	return fetchAll(m.db)
}

func fetchAll(db *sql.DB) ([]*Entry, error) {
	query, err := db.Prepare("SELECT id, author, content, created_at FROM entries")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var entries []*Entry

	rows, err := query.Query()
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for rows.Next() {
		var e Entry
		err = rows.Scan(&e.ID, &e.Author, &e.Content, &e.CreatedAt)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		entries = append(entries, &e)
	}

	return entries, nil
}

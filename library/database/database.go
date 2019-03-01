package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"vib-library/library/structs"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "viberate"
	password = "viberate"
	dbName   = "viberate"
)

type Database struct {
	conn *sql.DB
}

//  NewDatabase initializes and returns new database instance
func NewDatabase() (*Database, error) {
	s := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	conn, err := sql.Open("postgres", s)
	if err != nil {
		return nil, err
	}

	return &Database{
		conn: conn,
	}, nil
}

func (db Database) Ping() bool {
	if err := db.conn.Ping(); err != nil {
		return false
	}
	return true
}

func (db Database) AddMember(firstName, lastName string) error {
	stmt := `INSERT INTO member (first_name, last_name) VALUES ($1, $2)`
	_, err := db.conn.Exec(stmt, firstName, lastName)
	return err
}

func (db Database) GetMembers() ([]structs.Member, error) {
	members := []structs.Member{}

	stmt := `SELECT id, first_name, last_name FROM member`
	rows, err := db.conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := structs.Member{}
		if err := rows.Scan(&m.Id, &m.FirstName, &m.LastName); err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func (db Database) RentBook(memberId, bookId string) error {
	// check if book available
	var available bool
	stmt := `SELECT amount > 0 FROM availability WHERE book_id=$1`
	err := db.conn.QueryRow(stmt, bookId).Scan(&available)
	if err != nil {
		return err
	}

	if available == false {
		return errors.New("book unavailable")
	}

	// begin transaction
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	// rent book
	stmt = `INSERT INTO rent (member_id, book_id, time) VALUES ($1, $2, $3)`
	_, err = tx.Exec(stmt, memberId, bookId, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	// change availabile amount
	stmt = `UPDATE availability SET amount=amount-1 WHERE book_id=$1`
	_, err = tx.Exec(stmt, bookId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// commit
	return tx.Commit()
}

func (db Database) ReturnBook(rentId string) error {
	// get book id
	var bookId string
	stmt := `SELECT book_id FROM rent WHERE id=$1`
	err := db.conn.QueryRow(stmt, rentId).Scan(&bookId)
	if err != nil {
		return err
	}

	// begin transaction
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	// mark book as returned
	stmt = `UPDATE rent SET time_return=$2 WHERE id=$1 AND time_return IS NULL`
	_, err = tx.Exec(stmt, rentId, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	// change availabile amount
	stmt = `UPDATE availability SET amount=amount+1 WHERE book_id=$1`
	_, err = tx.Exec(stmt, bookId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// commit
	return tx.Commit()
}

func (db Database) GetAvailableBooks() ([]structs.Book, error) {
	books := []structs.Book{}

	stmt := `SELECT b.id, b.title, b.author, b.year, a.amount FROM book b JOIN availability a ON b.id = a.book_id WHERE a.amount > 0`
	rows, err := db.conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		b := structs.Book{}
		if err := rows.Scan(&b.Id, &b.Title, &b.Author, &b.Year, &b.AvailableAmount); err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

package vivom

import (
	"database/sql"
	"errors"
)

type Row interface {
	TableName() string
	Columns() []string
}

type SelectableRow interface {
	Row

	ScanValues() []interface{}
}

type InsertableRow interface {
	Row

	GetID() int
	SetID(int)
	Validate() error
	Values() []interface{}
}

type SelectableRows interface {
	Row

	Next() SelectableRow
}

func csv(values []string) string {
	list := ""
	for i, value := range values {
		if i != (len(values) - 1) {
			list += value + ", "
		} else {
			list += value
		}
	}
	return list
}

func csQ(columns []string) string {
	questions := make([]string, len(columns))
	for i, _ := range columns {
		questions[i] = "?"
	}
	return csv(questions)
}

func Insert(r InsertableRow, db *sql.DB) error {
	if r.GetID() != 0 {
		return errors.New("can't insert tag with id")
	}

	err := r.Validate()
	if err != nil {
		return err
	}

	query := "INSERT INTO " + r.TableName() + " (" + csv(r.Columns()) + ") values (" + csQ(r.Columns()) + ")"
	res, err := db.Exec(query, r.Values()...)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	r.SetID(int(id))

	return nil
}

func Select(r SelectableRow, id string, db *sql.DB) error {
	row := db.QueryRow("SELECT id, "+csv(r.Columns())+" FROM "+r.TableName()+" WHERE id=?", id)
	err := row.Scan(r.ScanValues()...)
	if err != nil {
		return err
	}

	return nil
}

func SelectAll(rs SelectableRows, db *sql.DB) error {
	rows, err := db.Query("SELECT id, " + csv(rs.Columns()) + " FROM " + rs.TableName())
	if err != nil {
		return err
	}

	for rows.Next() {
		r := rs.Next()
		err := rows.Scan(r.ScanValues()...)
		if err != nil {
			return err
		}
	}

	return nil
}

func SelectAllBy(rs SelectableRows, column string, value string, db *sql.DB) error {
	query := "SELECT id, " + csv(rs.Columns()) + " FROM " + rs.TableName() + " WHERE "
	query += column + "=?"
	rows, err := db.Query(query, value)
	if err != nil {
		return err
	}

	for rows.Next() {
		r := rs.Next()
		err := rows.Scan(r.ScanValues()...)
		if err != nil {
			return err
		}
	}

	return nil
}

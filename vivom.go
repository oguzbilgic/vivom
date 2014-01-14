// Package vivom is a dead simple, tiny, but powerful ORM library.
package vivom

import (
	"database/sql"
	"errors"
)

type Row interface {
	Table() string

	// Columns returns column names of the table, starting with the primary key.
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

func csQ(n int) string {
	questions := make([]string, n)
	for i := 0; i < n; i++ {
		questions[i] = "?"
	}
	return csv(questions)
}

type Vivom struct {
	db *sql.DB
}

func New(db *sql.DB) *Vivom {
	return &Vivom{db}
}

func (v *Vivom) Insert(r InsertableRow) error {
	if r.GetID() != 0 {
		return errors.New("can't insert tag with id")
	}

	err := r.Validate()
	if err != nil {
		return err
	}

	query := "INSERT INTO " + r.Table() + " (" + csv(r.Columns()[1:]) + ") values (" + csQ(len(r.Columns())-1) + ")"
	res, err := v.db.Exec(query, r.Values()...)
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

func (v *Vivom) Update(r InsertableRow) error {
	if r.GetID() == 0 {
		return errors.New("doesn't have an ID")
	}

	err := r.Validate()
	if err != nil {
		return err
	}

	columns := r.Columns()[1:]
	for i, column := range columns {
		columns[i] = column + "=?"
	}

	query := "UPDATE " + r.Table() + " SET " + csv(columns) + " WHERE " + r.Columns()[0] + "=?"
	_, err = v.db.Exec(query, append(r.Values(), r.GetID())...)

	return err
}

func (v *Vivom) Select(r SelectableRow, id string) error {
	row := v.db.QueryRow("SELECT "+csv(r.Columns())+" FROM "+r.Table()+" WHERE "+r.Columns()[0]+"=?", id)
	return row.Scan(r.ScanValues()...)
}

func (v *Vivom) SelectAll(rs SelectableRows) error {
	return v.SelectAllBy(rs, "", "")
}

func (v *Vivom) SelectAllBy(rs SelectableRows, column string, value string) error {
	query := "SELECT " + csv(rs.Columns()) + " FROM " + rs.Table()

	if column != "" && value != "" {
		// TODO Make sure this line is safe
		query += " WHERE " + column + "=" + value
	}

	rows, err := v.db.Query(query)
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

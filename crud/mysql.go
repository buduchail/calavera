package crud

import (
	"fmt"
	"sync"
	"errors"
	"strings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/buduchail/catrina"
)

type (
	MySqlCRUD struct {
		db      *sql.DB
		table   string
		id      string
		fields  []string
		hydrate MySqlHydrateFunc

		stmt struct {
			lock             sync.RWMutex
			insertStatement  *sql.Stmt
			selectStatements map[string]*sql.Stmt
			updateStatement  *sql.Stmt
			deleteStatement  *sql.Stmt
		}
	}

	MySqlHydrateFunc func(rows sql.Rows) (interface{}, error)
)

func NewMySqlCRUD(dsn string, table string, fields []string, hydrate MySqlHydrateFunc) (*MySqlCRUD, error) {

	if len(fields) == 0 {
		return nil, errors.New("At least one field must be defined")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	mysql := MySqlCRUD{
		db:      db,
		table:   table,
		id:      fields[0],
		fields:  fields,
		hydrate: hydrate,
	}

	mysql.stmt.lock = sync.RWMutex{}
	mysql.stmt.selectStatements = make(map[string]*sql.Stmt, 0)

	return &mysql, nil
}

// Helper methods

func (r *MySqlCRUD) getInsertStatement() (*sql.Stmt, error) {

	if r.stmt.insertStatement == nil {
		r.stmt.lock.Lock()
		defer r.stmt.lock.Unlock()
		// first field is ID field
		fields := strings.Join(r.fields[1:], ",")
		values := strings.Repeat(",?", len(r.fields)-1)[1:]
		stmt, err := r.db.Prepare(
			fmt.Sprintf(
				"INSERT INTO %s (%s) VALUES (%s)",
				r.table,
				fields,
				values,
			),
		)
		if err != nil {
			return nil, err
		}
		r.stmt.insertStatement = stmt
	}

	return r.stmt.insertStatement, nil
}

func (r *MySqlCRUD) getSelectStatement(where string) (*sql.Stmt, error) {

	_, prepared := r.stmt.selectStatements[where]
	if !prepared {
		r.stmt.lock.Lock()
		defer r.stmt.lock.Unlock()
		stmt, err := r.db.Prepare(
			fmt.Sprintf(
				"SELECT %s FROM %s WHERE %s",
				strings.Join(r.fields, ","),
				r.table,
				where,
			),
		)
		if err != nil {
			return nil, err
		}
		r.stmt.selectStatements[where] = stmt
	}

	return r.stmt.selectStatements[where], nil
}

func (r *MySqlCRUD) getUpdateStatement() (*sql.Stmt, error) {

	if r.stmt.updateStatement == nil {
		r.stmt.lock.Lock()
		defer r.stmt.lock.Unlock()
		fields := ""
		// first field is ID field
		for _, f := range r.fields[1:] {
			fields += ", " + f + " = ?"
		}
		stmt, err := r.db.Prepare(
			fmt.Sprintf(
				"UPDATE %s SET %s WHERE %s = ?",
				r.table,
				fields[2:],
				r.id,
			),
		)
		if err != nil {
			return nil, err
		}
		r.stmt.updateStatement = stmt
	}

	return r.stmt.updateStatement, nil
}

func (r *MySqlCRUD) getDeleteStatement() (*sql.Stmt, error) {

	if r.stmt.deleteStatement == nil {
		r.stmt.lock.Lock()
		defer r.stmt.lock.Unlock()
		stmt, err := r.db.Prepare(
			fmt.Sprintf(
				"DELETE FROM %s WHERE %s = ?",
				r.table,
				r.id,
			),
		)
		if err != nil {
			return nil, err
		}
		r.stmt.deleteStatement = stmt
	}

	return r.stmt.deleteStatement, nil
}

func (r *MySqlCRUD) selectMany(where string, values []interface{}) (<-chan catrina.Row, error) {

	result := make(chan catrina.Row)

	stmt, err := r.getSelectStatement(where)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(result)
		defer rows.Close()

		for rows.Next() {
			obj, err := r.hydrate(*rows)
			if err != nil {
				result <- catrina.Row{nil, err}
			} else {
				result <- catrina.Row{obj, nil}
			}
		}

		err = rows.Err()
		if err != nil {
			// TODO: should we panic here?
			result <- catrina.Row{nil, err}
		}
	}()

	return result, nil
}

func (r *MySqlCRUD) castValues(values []catrina.Value) []interface{} {

	interfaces := make([]interface{}, len(values))
	for i, v := range values {
		interfaces[i] = v
	}

	return interfaces
}

// Public interface

func (r *MySqlCRUD) Insert(values []catrina.Value) (id catrina.Value, e error) {

	if len(r.fields)-1 != len(values) {
		return nil, errors.New("Value count does not match field count")
	}

	stmt, err := r.getInsertStatement()
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(r.castValues(values)...)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return lastID, nil
}

func (r *MySqlCRUD) Select(id catrina.Value) (catrina.Object, error) {

	stmt, err := r.getSelectStatement(r.id + " = ?")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("No rows found")
	}

	return r.hydrate(*rows)
}

func (r *MySqlCRUD) SelectWhereFields(fields []string, values []catrina.Value) (<-chan catrina.Row, error) {

	if len(fields) != len(values) {
		return nil, errors.New("Fields and values do not match")
	}

	where := ""
	for _, f := range fields {
		where += " AND " + f + " = ?"
	}

	return r.selectMany(where[4:], r.castValues(values))
}

func (r *MySqlCRUD) SelectWhereRange(field string, min, max catrina.Value) (<-chan catrina.Row, error) {

	return r.selectMany(
		field+" BETWEEN ? AND ?",
		r.castValues([]catrina.Value{min, max}),
	)
}

func (r *MySqlCRUD) SelectWhereExpression(where string, values []catrina.Value) (<-chan catrina.Row, error) {

	return r.selectMany(where, r.castValues(values))
}

func (r *MySqlCRUD) Update(id catrina.Value, values []catrina.Value) error {

	if len(r.fields)-1 != len(values) {
		return errors.New("Value count does not match field count")
	}

	stmt, err := r.getUpdateStatement()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(r.castValues(append(values, id))...)
	if err != nil {
		// TODO: don't treat warning as errors (e.g. trimmed data)
		return err
	}

	return nil
}

func (r *MySqlCRUD) Delete(id catrina.Value) error {

	stmt, err := r.getDeleteStatement()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		// TODO: don't treat warning as errors (e.g. trimmed data)
		return err
	}

	return nil
}

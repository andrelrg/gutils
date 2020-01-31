package db

import (
	"database/sql"
	"fmt"
)

type Database struct{
	ConnectionLine string
	Conn *sql.DB
	Driver string
}


//Ping Tests the connection.
func (d *Database) Ping() error {
	return d.Conn.Ping()
}

func (d *Database) Connect() error {
	conn, err := sql.Open("postgres", d.ConnectionLine)
	if err != nil {
		return err
	}
	d.Conn = conn
	return d.Ping()
}

// Execute executes the query received with the given parameters.
func (d *Database) Execute(query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result

	stmtIns, err := d.Conn.Prepare(query)
	if err != nil {
		return result, err
	}
	defer stmtIns.Close()

	result, err = stmtIns.Exec(args...)
	if err != nil {
		return result, err
	}

	return result, nil
}

//MapScan Get the result of a query in this format: map[string]interface{}
func (d *Database) MapScan(query string, args ...interface{}) (map[string]interface{}, error) {
	stmt, err := d.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	entry := make(map[string]interface{})
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry = make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
	}
	return entry, nil
}

//SliceMapScan Fetch all lines of given select
func (d *Database) SliceMapScan(query string, args ...interface{}) ([]map[string]interface{}, error) {
	stmt, err := d.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil
}

//QueryRow Get the next Query Row.
func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.Conn.QueryRow(query, args...)
}

//GetJSON Get first row in JSON format.
func (d *Database) GetJSON(sqlString string) (map[string]interface{}, error) {
	rows, err := d.Conn.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	entry := make(map[string]interface{})
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry = make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
	}
	return entry, nil
}

// GetJSONList Fetch all lines of given select
func (d *Database) GetJSONList(sqlString string) ([]map[string]interface{}, error) {
	rows, err := d.Conn.Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil
}

// Close is responsible for closing database connection
func (d *Database) Close() {
	err := d.Conn.Close()
	fmt.Println(err)
	panic(0)
}

package db

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astrolink/gutils/cache"
	"strconv"
	"strings"
	"time"
)

type Database struct {
	ConnectionLine string
	Conn           *sql.DB
	Driver         string
	CacheConfig    cache.Config
}

//Ping Tests the connection.
func (d *Database) Ping() error {
	return d.Conn.Ping()
}

func (d *Database) Connect() error {
	conn, err := sql.Open(d.Driver, d.ConnectionLine)
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


// ExecuteWithTx executes the received query with the parameters provided within a transaction.
func (d *Database) ExecuteWithTx(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result

	stmtIns, err := tx.Prepare(query)
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

			switch values[i].(type) {
			case []uint8:
				// converting everything to an array of interface
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}

				in := strings.Index(v.(string), "{")
				v = strings.Replace(v.(string), "{", "", -1)
				v = strings.Replace(v.(string), "}", "", -1)
				items := strings.Split(v.(string), ",")
				if v.(string) == "" {
					v = make([]string, 0)
				} else if len(items) == 1 && in == -1 {
					v = string(b)
				} else {
					v = items
				}
			default:
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}
			}

			entry[col] = v
		}
	}
	return entry, nil
}

//  BasicMapScan Get the result of a query in this format: map[string]interface{}
//  Individual types:
//	bool, for booleans
//	float64, for numbers
//	string, for strings
//	[]interface{}, for arrays
func (d *Database) BasicMapScan(query string, args ...interface{}) (map[string]interface{}, error) {
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

			switch values[i].(type) {
			case string:
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}
			case bool:
				v = val
			case int, int8, int16, int32, int64, float32:
				value := fmt.Sprintf("%v", val)
				parseFloat, parseErr := strconv.ParseFloat(value, 64)

				if parseErr != nil {
					v = val
				} else {
					v = parseFloat
				}
			case float64:
				v = val
			case []uint8:
				// converting everything to an array of interface
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}

				in := strings.Index(v.(string), "{")
				v = strings.Replace(v.(string), "{", "", -1)
				v = strings.Replace(v.(string), "}", "", -1)
				items := strings.Split(v.(string), ",")
				if v.(string) == "" {
					v = make([]string, 0)
				} else if len(items) == 1 && in == -1 {
					v = string(b)
				} else {
					v = items
				}
			default:
				if val == nil{
					v = nil
				} else {
					jsonType, marshalErr := json.Marshal(val)

					if marshalErr != nil {
						v = val
					} else {
						v = string(jsonType)
					}
				}
			}

			entry[col] = v
		}
	}
	return entry, nil
}

func (d *Database) MapScanRedis(query string, duration time.Duration, args ...interface{}) (map[string]interface{}, error) {
	if d.CacheConfig == nil {
		return d.BasicMapScan(query, args...)
	}

	redis, err := cache.NewRedis(d.CacheConfig)

	if err != nil {
		fmt.Println("fail to connect to the cache: ", err.Error())

		return d.BasicMapScan(query, args...)
	}

	defer redis.Close()

	var hashString string

	hashString = generateKeyByQueryAndArgs(query, args...)

	cacheResult, err := redis.Get(hashString)

	if err != nil {
		return d.setRedisByMapScan(query, redis, hashString, duration, args...)
	}

	var result map[string]interface{}

	err = json.Unmarshal([]byte(cacheResult), &result)

	if err != nil {
		fmt.Println("unmarshal error: " + err.Error())

		return d.setRedisByMapScan(query, redis, hashString, duration, args...)
	}

	return result, nil
}

func (d *Database) SliceMapScanRedis(query string, duration time.Duration, args ...interface{}) ([]map[string]interface{}, error) {
	if d.CacheConfig == nil {
		return d.BasicSliceMapScan(query, args...)
	}

	redis, err := cache.NewRedis(d.CacheConfig)

	if err != nil {
		fmt.Println("fail to connect to the cache: ", err.Error())

		return d.BasicSliceMapScan(query, args...)
	}

	defer redis.Close()

	var hashString string

	hashString = generateKeyByQueryAndArgs(query, args...)

	cacheResult, err := redis.Get(hashString)

	if err != nil {
		return d.setRedisBySliceMapScan(query, redis, hashString, duration, args...)
	}

	var result []map[string]interface{}

	err = json.Unmarshal([]byte(cacheResult), &result)

	if err != nil {
		fmt.Println("unmarshal error: " + err.Error())

		return d.setRedisBySliceMapScan(query, redis, hashString, duration, args...)
	}

	return result, nil
}

func generateKeyByQueryAndArgs(query string, args ...interface{}) string {
	for _, s := range args {
		query += fmt.Sprintf("%v", s)
	}

	data := []byte(query)
	hash := sha256.Sum256(data)

	var hashString string

	hashString = fmt.Sprintf("%x", hash[:])

	return hashString
}

func (d *Database) setRedisByMapScan(query string, redis *cache.Redis, redisKey string, duration time.Duration, args ...interface{}) (map[string]interface{}, error) {
	data, mapScanError := d.BasicMapScan(query, args...)

	if data == nil {
		return data, mapScanError
	}

	jsonBytes, marshalError := json.Marshal(data)

	if marshalError == nil {
		_ = redis.Set(redisKey, string(jsonBytes), duration)
	}

	return data, mapScanError
}

func (d *Database) setRedisBySliceMapScan(query string, redis *cache.Redis, redisKey string, duration time.Duration, args ...interface{}) ([]map[string]interface{}, error) {
	data, mapScanError := d.BasicSliceMapScan(query, args...)

	if data == nil {
		return data, mapScanError
	}

	jsonBytes, marshalError := json.Marshal(data)

	if marshalError == nil {
		_ = redis.Set(redisKey, string(jsonBytes), duration)
	}

	return data, mapScanError
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

			switch values[i].(type) {
			case []uint8:
				// converting everything to an array of interface
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}

				in := strings.Index(v.(string), "{")
				v = strings.Replace(v.(string), "{", "", -1)
				v = strings.Replace(v.(string), "}", "", -1)
				items := strings.Split(v.(string), ",")
				if v.(string) == "" {
					v = make([]string, 0)
				} else if len(items) == 1 && in == -1 {
					v = string(b)
				} else {
					v = items
				}
			default:
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}
			}

			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil
}

//  BasicSliceMapScan Fetch all lines of given select
//  Individual types:
//	bool, for booleans
//	float64, for numbers
//	string, for strings
//	[]interface{}, for arrays
func (d *Database) BasicSliceMapScan(query string, args ...interface{}) ([]map[string]interface{}, error) {
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

			switch values[i].(type) {
			case string:
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}
			case bool:
				v = val
			case int, int8, int16, int32, int64, float32:
				value := fmt.Sprintf("%v", val)
				parseFloat, parseErr := strconv.ParseFloat(value, 64)

				if parseErr != nil {
					v = val
				} else {
					v = parseFloat
				}
			case float64:
				v = val
			case []uint8:
				// converting everything to an array of interface
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}

				in := strings.Index(v.(string), "{")
				v = strings.Replace(v.(string), "{", "", -1)
				v = strings.Replace(v.(string), "}", "", -1)
				items := strings.Split(v.(string), ",")
				if v.(string) == "" {
					v = make([]string, 0)
				} else if len(items) == 1 && in == -1 {
					v = string(b)
				} else {
					v = items
				}
			default:
				if val == nil{
					v = nil
				} else {
					jsonType, marshalErr := json.Marshal(val)

					if marshalErr != nil {
						v = val
					} else {
						v = string(jsonType)
					}
				}
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

func (d *Database) StartTransaction() (*sql.Tx, error) {
	return d.Conn.Begin()
}


// Close is responsible for closing database connection
func (d *Database) Close() {
	err := d.Conn.Close()
	if err != nil {
		fmt.Println(err)
	}
}

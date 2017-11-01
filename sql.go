package sqlext

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

var ErrNotSupport = errors.New("only support insert []T or []*T")

// BatchInsert, rows should be []T or []*T
func BatchInsert(db *sql.DB, rows interface{}) (result sql.Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	kind := reflect.TypeOf(rows).Kind()
	if kind != reflect.Slice {
		return nil, ErrNotSupport
	}

	sqlStr, values := genInsertSql(rows)
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(values...)
}

func genInsertSql(rows interface{}) (string, []interface{}) {

	var (
		column     string
		needColumn = true
		table      string
		values     = []interface{}{}
		raw        = reflect.ValueOf(rows)
		sql        = "INSERT INTO %s ( %s ) VALUES "
	)

	for i := 0; i < raw.Len(); i++ {
		val := reflect.ValueOf(raw.Index(i).Interface())
		for val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		line := "("
		tp := reflect.Indirect(val).Type()
		if table == "" {
			table = snake(tp.Name())
		}
		for i := 0; i < val.NumField(); i++ {
			if needColumn {
				column += snake(tp.Field(i).Name) + ","
			}
			line += "?,"
			values = append(values, val.Field(i).Interface())
		}
		line = strings.TrimSuffix(line, ",")
		line += "),"
		sql += line
		needColumn = false
	}

	sql = strings.TrimSuffix(sql, ",")
	column = strings.TrimSuffix(column, ",")
	return fmt.Sprintf(sql, table, column), values
}

// MapToStruct fill struct with map value
func MapToStruct(m map[string]interface{}, s interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	kind := reflect.TypeOf(s).Kind()
	if kind != reflect.Ptr {
		return errors.New("second param should be a pointer")
	}

	rs := reflect.ValueOf(s).Elem()
	tp := reflect.Indirect(rs).Type()

	for i := 0; i < rs.NumField(); i++ {
		if v, ok := m[snake(tp.Field(i).Name)]; ok {
			val := reflect.ValueOf(v)
			rs.Field(i).Set(val)
		}
	}
	return nil
}

func snake(in string) string {
	runes := []rune(in)
	length := len(runes)
	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 &&
			unicode.IsUpper(runes[i]) &&
			((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}
	return string(out)
}

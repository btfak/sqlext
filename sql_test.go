package sqlext

import (
	"fmt"
	"testing"
	"time"
)

type Group struct {
	ID         int64
	Name       string
	CreateTime time.Time
}

func TestGenInsertSql(t *testing.T) {
	t1 := time.Date(2016, 1, 1, 0, 0, 0, 0, time.Local)
	row1 := Group{ID: 1, Name: "row1", CreateTime: t1}
	t2 := time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)
	row2 := Group{ID: 2, Name: "row2", CreateTime: t2}

	wantSql := "INSERT INTO group ( id,name,create_time ) VALUES (?,?,?),(?,?,?)"
	wantValues := []interface{}{1, "row1", t1, 2, "row2", t2}
	sql := ""
	values := []interface{}{}

	check := func() {
		if sql != wantSql {
			t.Error("sql error", sql)
		}
		if len(values) != 6 {
			t.Error("value error", values)
		}
		for k, v := range values {
			if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", wantValues[k]) {
				t.Error("value error", v, wantValues[k])
			}
		}
	}

	rows := []Group{row1, row2}
	sql, values = genInsertSql(rows)
	check()

	ptrRows := []*Group{&row1, &row2}
	sql, values = genInsertSql(ptrRows)
	check()

}

func BenchmarkGenInsertSql(b *testing.B) {
	t1 := time.Date(2016, 1, 1, 0, 0, 0, 0, time.Local)
	row1 := Group{ID: 1, Name: "row1", CreateTime: t1}
	rows := []interface{}{}
	for i := 0; i < 100; i++ {
		rows = append(rows, row1)
	}
	for i := 0; i < b.N; i++ {
		genInsertSql(rows)
	}
}

type User struct {
	ID    int64
	Year  int32
	Name  string
	Blank string
	Time  time.Time
}

var now = time.Now()
var m = map[string]interface{}{
	"id":     int64(100),
	"year":   int32(24),
	"name":   "gocql",
	"extend": true,
	"time":   now,
}

func TestMapToStruct(t *testing.T) {
	var user User
	err := MapToStruct(m, &user)
	if err != nil {
		t.Fatal(err)
	}
	if user.ID != 100 || user.Year != 24 || user.Name != "gocql" ||
		user.Blank != "" || user.Time != now {
		t.Fatal("map to struct fail")
	}
}

func BenchmarkMapToStruct(b *testing.B) {
	var user User
	for i := 0; i < b.N; i++ {
		err := MapToStruct(m, &user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

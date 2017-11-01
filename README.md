# sqlext

sqlext use to extend [gorm](https://github.com/jinzhu/gorm) and [gocql](https://github.com/gocql/gocql)

- gorm not support batch insert，`BatchInsert` is a universal batch insert API.
- gocql not support bind querying data to struct，`MapToStruct` fill struct fields with map value


##example

- BatchInsert  

```go
type GroupMember struct {
	ID        int64
	GroupID   int64
	Type      int8
	Extra     []byte
	JoinTime  time.Time
}

var mydb *gorm.DB
var members = []GroupMember{}
sqlext.BatchInsert(mydb.DB(),members)

//生成的SQL: INSERT INTO group_member (id, group_id, type, extra, join_time) VALUES (?,?,?,?,?), (?,?,?,?,?) ...
```

- MapToStruct

```go
// MapToStruct
func GetUsersByIDs(userIDs []int64) (users []User, err error) {
	iter := Cassandra.Query("SELECT * FROM user WHERE id IN ?", userIDs).Iter()
	for {
		row := make(map[string]interface{})
		if !iter.MapScan(row) {
			break
		}
		user := User{}
		err = sqlext.MapToStruct(row, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
```
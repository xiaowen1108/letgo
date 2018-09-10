package letgo

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

type Db struct {
	coon           *sql.DB
	sql            string //本次sql
	tableName      string
	tableNameAlias string
	field          []string
	join           string
	set            map[string]string
	where          map[string]interface{}
	whereIn        map[string][]interface{}
	whereBetween   map[string][2]interface{}
	order          map[string]string
	limitNum       int
}

var (
	TableNameError = errors.New("table name is null")
)

const (
	SELECT = iota
	UPDATE
	DELETE
	INSERT
)

func (db *Db) GetSQl(op int) (string, error) {
	switch op {
	case SELECT:
		return db.GetSelectSQl()
	case UPDATE:
		return db.GetSelectSQl()
	case DELETE:
		return db.GetSelectSQl()
	case INSERT:
		return db.GetSelectSQl()
	default:
		return db.GetSelectSQl()
	}
}
func (db *Db) GetSelectSQl() (string, error) {
	var sql bytes.Buffer
	sql.WriteString("SELECT ")
	//field
	db.writeField(&sql)
	//table name
	sql.WriteString(" FROM `")
	sql.WriteString(db.tableName)
	sql.WriteString("`")
	//alias
	db.writeAlias(&sql)
	//JOIN
	db.writeJoin(&sql)
	//where
	db.writeWhere(&sql)
	//order
	db.writeOrder(&sql)
	//limit
	db.writeLimit(&sql)
	db.sql = sql.String()
	return db.sql, nil
}
func (db *Db) writeField(sql *bytes.Buffer) {
	if db.field != nil {
		for k, v := range db.field {
			if k != 0 {
				sql.WriteString(", ")
			}
			sql.WriteString("`")
			sql.WriteString(v)
			sql.WriteString("`")
		}
	} else {
		sql.WriteString("*")
	}
}
func (db *Db) writeAlias(sql *bytes.Buffer) {
	if db.tableNameAlias != "" {
		sql.WriteString(" AS `")
		sql.WriteString(db.tableNameAlias)
		sql.WriteString("`")
	}
}
func (db *Db) writeJoin(sql *bytes.Buffer) {
	if db.join != "" {
		sql.WriteString(" ")
		sql.WriteString(db.join)
		sql.WriteString(" ")
	}
}
func (db *Db) writeWhere(sql *bytes.Buffer) {
	var ix int
	if db.where != nil || db.whereIn != nil || db.whereBetween != nil {
		sql.WriteString(" WHERE ")
	}
	if db.where != nil {
		for k, v := range db.where {
			if ix != 0 {
				sql.WriteString(" AND ")
			}
			if !strings.Contains(k, "?") {
				sql.WriteString("`")
				sql.WriteString(k)
				sql.WriteString("` = '")
				sql.WriteString(v.(string))
				sql.WriteString("'")
			} else {
				k = strings.Replace(k, "?", "'%v'", 1)
				sql.WriteString(fmt.Sprintf(k, v))
			}
			ix++
		}
	}
	if db.whereIn != nil {
		for k, v := range db.whereIn {
			if ix != 0 {
				sql.WriteString(" AND ")
			}
			sql.WriteString("`")
			sql.WriteString(k)
			sql.WriteString("` IN (")
			for k1, v1 := range v {
				if k1 != 0 {
					sql.WriteString(", ")
				}
				sql.WriteString("'")
				in, ok := v1.(string)
				if !ok {
					continue
				}
				sql.WriteString(in)
				sql.WriteString("'")
			}
			sql.WriteString(")")
			ix++
		}
	}
	if db.whereBetween != nil {
		for k, v := range db.whereBetween {
			if ix != 0 {
				sql.WriteString(" AND ")
			}
			start, ok := v[0].(string)
			if !ok {
				continue
			}
			end, ok := v[1].(string)
			if !ok {
				continue
			}
			sql.WriteString("`")
			sql.WriteString(k)
			sql.WriteString("` BETWEEN (")
			sql.WriteString("'")
			sql.WriteString(start)
			sql.WriteString("', ")
			sql.WriteString("'")
			sql.WriteString(end)
			sql.WriteString("'")
			sql.WriteString(")")
			ix++
		}
	}
}
func (db *Db) writeOrder(sql *bytes.Buffer) {
	if db.order != nil {
		var ix int
		sql.WriteString(" ORDER BY ")
		for k, v := range db.order {
			if ix != 0 {
				sql.WriteString(", ")
				ix++
			}
			sql.WriteString(k)
			sql.WriteString(" ")
			sql.WriteString(v)
		}
	}
}

func (db *Db) writeLimit(sql *bytes.Buffer) {
	if db.limitNum > 0 {
		sql.WriteString(" LIMIT ")
		sql.WriteString(strconv.Itoa(db.limitNum))
	}
}

//table
func (db *Db) Table(name string) *Db {
	db.tableName = name
	return db
}

//table
func (db *Db) Alias(name string) *Db {
	db.tableNameAlias = name
	return db
}

//Feild
func (db *Db) Field(field ...string) *Db {
	db.field = field
	return db
}

//JOIN
func (db *Db) Join(join string) *Db {
	db.join = join
	return db
}

//Set
func (db *Db) Set(key string, value string) *Db {
	if db.set == nil {
		db.set = make(map[string]string)
	}
	db.set[key] = value
	return db
}

//Where  aaa,bbb    aaa = ?, bbb
func (db *Db) Where(key string, value interface{}) *Db {
	if db.where == nil {
		db.where = make(map[string]interface{})
	}
	switch value.(type) {
	case int:
		db.where[key] = strconv.Itoa(value.(int))
	case string:
		db.where[key] = value.(string)
	}
	return db
}
func (db *Db) WhereIn(key string, in []interface{}) *Db {
	if db.whereIn == nil {
		db.whereIn = make(map[string][]interface{})
	}
	db.whereIn[key] = in
	return db
}
func (db *Db) WhereBetween(key string, between [2]interface{}) *Db {
	if db.whereBetween == nil {
		db.whereBetween = make(map[string][2]interface{})
	}
	db.whereBetween[key] = between
	return db
}

//order
func (db *Db) Order(key string, order string) *Db {
	if db.order == nil {
		db.order = make(map[string]string)
	}
	db.order[key] = order
	return db
}

//limit
func (db *Db) Limit(num int) *Db {
	db.limitNum = num
	return db
}

//select
func (db *Db) Select() (string, error) {
	return db.GetSQl(SELECT)
}

//update
func (db *Db) Update() (string, error) {
	return db.GetSQl(UPDATE)
}

//delete
func (db *Db) Delete() (string, error) {
	return db.GetSQl(DELETE)
}

//insert
func (db *Db) Insert() (string, error) {
	return db.GetSQl(INSERT)
}

//清空
func (db *Db) Clear() {

}

func (db *Db) GetLastSQl() string {
	return db.sql
}

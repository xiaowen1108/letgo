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
	coon         *sql.DB
	sql          string //本次sql
	tableName    string
	field        []string
	set          map[string]string
	where        map[string]interface{}
	whereIn      map[string][]interface{}
	whereBetween map[string][2]interface{}
	order        map[string]string
	limitNum     int
}

var (
	TableNameError = errors.New("table name is null")
)

func (db *Db) GetSQl() (string, error) {
	var field, set, where, whereIn, whereBetween, order, sql bytes.Buffer
	var ix int
	if db.tableName == "" {
		return "", TableNameError
	}
	if db.field != nil {
		for k, v := range db.field {
			if k != 0 {
				field.WriteString(", ")
			}
			field.WriteString("`")
			field.WriteString(v)
			field.WriteString("`")
		}
	}
	if db.set != nil {
		for k, v := range db.set {
			if ix == 0 {
				set.WriteString(", ")
			}
			set.WriteString("`")
			set.WriteString(k)
			set.WriteString("` = ")
			set.WriteString(v)
			ix++
		}
		ix = 0
	}
	if db.where != nil {
		for k, v := range db.where {
			if ix != 0 {
				where.WriteString(" AND ")
			}
			if !strings.Contains(k, "?") {
				where.WriteString("`")
				where.WriteString(k)
				where.WriteString("` = ")
				where.WriteString(v.(string))
			} else {
				k = strings.Replace(k, "?", "%v", 1)
				where.WriteString(fmt.Sprintf(k, v))
			}
			ix++
		}
		ix = 0
	}
	if db.whereIn != nil {

	}
	if db.whereBetween != nil {

	}
	if db.order != nil {
		for k, v := range db.order {
			if ix == 0 {
				set.WriteString(", ")
				ix++
			}
			order.WriteString(k)
			order.WriteString(" ")
			order.WriteString(v)
		}
		ix = 0
	}
	//组装sql  SELECT
	sql.WriteString("SELECT ")
	if field.Len() != 0 {
		sql.WriteString(field.String())
	} else {
		sql.WriteString("*")
	}
	sql.WriteString(" FROM `")
	sql.WriteString(db.tableName)
	sql.WriteString("`")
	if where.Len() != 0 {
		sql.WriteString(" WHERE ")
		sql.WriteString(where.String())
	}
	if whereIn.Len() != 0 && whereBetween.Len() != 0 {

	}
	if order.Len() != 0 {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(order.String())
	}
	if db.limitNum > 0 {
		sql.WriteString(" LIMIT ")
		sql.WriteString(strconv.Itoa(db.limitNum))
	}
	sql.WriteString(";")
	return sql.String(), nil
}

//table
func (db *Db) Table(name string) *Db {
	db.tableName = name
	return db
}

//Feild
func (db *Db) Field(field ...string) *Db {
	db.field = field
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
func (db *Db) Select() {

}

//update
func (db *Db) Update() {

}

//delete
func (db *Db) Delete() {

}

//清空
func (db *Db) Clear() {

}

func (db *Db) GetLastSQl() string {
	return db.sql
}

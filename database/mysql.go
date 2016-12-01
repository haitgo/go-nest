package database

import (
	"go-nest/conv"
	"strings"

	"github.com/go-xorm/xorm"
)

type H map[string]interface{}

/*
表类
*/
type Mysql struct {
	Engine *xorm.Engine
	arg    []interface{} //where使用
	table  string        //表名
	field  string        //查询字段
	sql    string        //sql语句
}

func (obj Mysql) New() *Mysql {
	return &obj
}

/*
通过mysql对象实例化表
*/
func (this *Mysql) Table(table string) *Mysql {
	this.table = table
	return this
}

/*
需要查询的字段
*/
func (this *Mysql) Field(cols string) *Mysql {
	this.field = cols
	return this
}

/*
查询条件，请不要输入where，
*/
func (this *Mysql) Where(sql string, arg ...interface{}) *Mysql {
	if sql != "" {
		this.sql = " WHERE " + sql
	}
	this.arg = arg
	return this
}

/*
插入数据，返回自增长id，可以输入一维map也可以输入2维map
*/
func (this *Mysql) Insert(data interface{}) (insertId int64, err error) {
	sqlaS, argA := this.arr2insert(data)
	sqlS := "INSERT INTO " + this.table + sqlaS
	argA = append(argA, this.arg...)
	result, err := this.Engine.Exec(sqlS, argA...)
	if err != nil {
		return
	}
	insertId, _ = result.LastInsertId()
	return
}

/*
修改数据，返回受影响的行数
*/
func (this *Mysql) Update(data H) (aff int64, err error) {
	sqlaS, argA := this.arr2update(data)
	sqlS := "UPDATE " + this.table + " SET " + sqlaS + this.sql
	argA = append(argA, this.arg...)
	result, err := this.Engine.Exec(sqlS, argA...)
	aff, _ = result.RowsAffected()
	return
}

/*
修改单条数据，返回受影响的行数
*/
func (this *Mysql) UpdateOne(data H) (rowsAffected int64, err error) {
	this.sql = this.sql + " LIMIT 1"
	return this.Update(data)
}

/*
删除数据，返回受影响的行数
*/
func (this *Mysql) Delete() (aff int64, err error) {
	sqlS := "DELETE  FROM " + this.table + this.sql
	result, err := this.Engine.Exec(sqlS, this.arg...)
	aff, _ = result.RowsAffected()
	return
}

/*
删除单条数据，返回受影响的行数
*/
func (this *Mysql) DeleteOne() (affectedNums int64, err error) {
	this.sql = this.sql + " LIMIT 1"
	return this.Delete()
}

/*
数组转换为修改sql语句
*/
func (this *Mysql) arr2update(arr H) (sql string, args []interface{}) {
	var sqlArr []string
	for k, v := range arr {
		if v == "" {
			continue
		}
		sqlArr = append(sqlArr, "`"+k+"`=?")
		args = append(args, v)
	}
	sql = strings.Join(sqlArr, ",")
	return
}

/*
数组转换为插入sql语句
*/
func (this *Mysql) arr2insert(arr interface{}) (sql string, args []interface{}) {
	var colsArr, valueArr []string
	var sqlCols, sqlValue string
	if arrData, ok := arr.(H); ok { //一维数组
		for k, v := range arrData {
			v = conv.MustString(v, "")
			if v == "" {
				continue
			}
			colsArr = append(colsArr, "`"+k+"`")
			valueArr = append(valueArr, "?")
			args = append(args, v)
		}
		sqlValue = "(" + strings.Join(valueArr, ",") + ")"
	} else if arrData, ok := arr.([]H); ok { //2维数组
		var cols []string //map数据的列
		for index, brr := range arrData {
			if index == 0 {
				for k, _ := range brr {
					colsArr = append(colsArr, "`"+k+"`")
					cols = append(cols, k)
				}
			}
			var valueBrr []string
			for _, k := range cols {
				v := conv.MustString(brr[k], "")
				valueBrr = append(valueBrr, "?")
				args = append(args, v)
			}
			value := "(" + strings.Join(valueBrr, ",") + ")"
			valueArr = append(valueArr, value)
		}
		sqlValue = strings.Join(valueArr, ",")
	}
	sqlCols = strings.Join(colsArr, ",")
	sql = "(" + sqlCols + ") values  " + sqlValue
	return
}

/**
数据库操作：所有的模型层都应该继承Model，方便使用数据库，默认自带add,modify,delete,get
          方法，可以使用Transaction添加事务处理
*/
package database

import (
	"errors"
	"math"

	"github.com/go-xorm/xorm"
)

var (
	Engine *xorm.Engine //数据库引擎,引入该包的时候来初始化
)

type Model struct {
	obj           interface{}  `xorm:"="` //继承者使用该对象
	id            *int32       `xorm:"="` //主键编号
	session       *Session     `xorm:"="` //数据库会话
	isRootSession bool         `xorm:"="` //是否为根会话
	engine        *xorm.Engine `xorm:"="` //数据库引擎，默认为配置的数据库引擎
}

//数据库引擎初始化
//obj=继承者实例化的对象，id=继承者的主键编号
func (this *Model) Construct(obj interface{}, id *int32) {
	this.obj = obj
	this.id = id
	this.engine = Engine
}

//数据库引擎,如果给了参数就设置引擎
func (this *Model) Engine(arg ...*xorm.Engine) *xorm.Engine {
	if len(arg) == 1 {
		this.engine = arg[0]
	}
	return this.engine
}

//事务引擎,如果给了arg参数，则为设置否则new一个session
func (this *Model) TransactionEngine(arg ...*Session) *Session {
	if this.session != nil {
		return this.session
	}
	if len(arg) == 1 {
		this.session = arg[0]
	} else if len(arg) == 0 {
		this.session = newSession()
		this.session.Begin() //开启事务
		this.isRootSession = true
	}
	return this.session
}

//自动提交session,如果处理成功则commit，否则rollback
func (this *Model) AutoCommit() error {
	return this.session.autoCommit(this.isRootSession)
}

//获取模型的数据，如果给了参数则将查询结果返回该参数，否则将查询的结果映射到默认的obj
func (this *Model) Info(arg ...interface{}) (bool, error) {
	obj := this.obj
	if len(arg) == 1 {
		obj = arg[0]
	}
	return this.Engine().Id(this.id).Get(obj)
}

//添加数据，如果使用了TransactionEngine,则使用session操作，默认使用Engine()
func (this *Model) Add() (int64, error) {
	if this.session != nil {
		return this.session.InsertOne(this.obj)
	}
	return this.Engine().InsertOne(this.obj)
}

//修改数据，如果使用了TransactionEngine,则使用session操作，默认使用Engine()
func (this *Model) Modify() (err error) {
	var aff int64
	if this.session != nil {
		aff, err = this.session.Id(this.id).Update(this.obj)
	} else {
		aff, err = this.Engine().Id(this.id).Update(this.obj)
	}
	if aff != 1 {
		err = errors.New("修改失败")
	}
	return err
}

//删除数据，如果使用了TransactionEngine,则使用session操作，默认使用Engine()
func (this *Model) Delete() (err error) {
	var aff int64
	if this.session != nil {
		aff, err = this.session.Id(this.id).Delete(this.obj)
	} else {
		aff, err = this.Engine().Id(this.id).Delete(this.obj)
	}
	if aff != 1 {
		err = errors.New("删除失败")
	}
	return err
}

//计算分页
//参数：totle=总记录数量，current=当前页码，limit=每一页数量
//返回：p=分页数据，start=sql查询limit起点
func (this *Model) Page(totle, current, limit int) (p *Page, start int) {
	p = new(Page)
	p.Pages = int(math.Ceil(float64(totle) / float64(limit))) //计算出页码
	if current <= 1 {                                         //页首不越界
		current = 1
	}
	if current >= p.Pages { //页末不越界
		current = p.Pages
	}
	p.Totle = totle
	p.Limit = limit
	p.Current = current
	start = (current - 1) * limit
	return
}

//------------------------------------------------------------------------------

//数据库会话，当使用事务的时候比用
type Session struct {
	xorm.Session      //继承xorm的session
	isHandleOK   bool //是否提交成功
}

func newSession() *Session {
	o := new(Session)
	return o
}

//自动判断是否事务提交成功，必须是在root下生成的事务才有效
func (this *Session) autoCommit(isRoot bool) error {
	if isRoot && this.isHandleOK {
		return this.Commit()
	}
	return nil
}

//会话处理成功
func (this *Session) HandleOK() {
	this.isHandleOK = true
}

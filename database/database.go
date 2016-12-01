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
	DO_INVALID = errors.New("数据操作无效")
	DO_FAIL    = errors.New("数据查询失败")
)

var (
	Engine *xorm.Engine //数据库引擎,引入该包的时候来初始化
)

type Model struct {
	omit          []string     //去除
	cols          []string     //需要的列
	mustCols      []string     //必须要的列
	obj           interface{}  `xorm:"-"` //继承者使用该对象
	id            *int32       `xorm:"-"` //主键编号
	session       *Session     `xorm:"-"` //数据库会话
	isRootSession bool         `xorm:"-"` //是否为根会话
	engine        *xorm.Engine `xorm:"-"` //数据库引擎，默认为配置的数据库引擎
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
func (this *Model) SessionEngine(arg ...*Session) *Session {
	if this.session != nil {
		return this.session
	}
	if len(arg) == 1 {
		this.session = arg[0]
	} else {
		this.session = NewSession(this.engine)
		this.session.Begin() //开启事务
		this.isRootSession = true
	}
	return this.session
}

//自动提交session,如果处理成功则commit，否则rollback
func (this *Model) AutoCommit() error {
	return this.session.autoCommit(this.isRootSession)
}

func (this *Model) Omit(omit ...string) *Model {
	this.omit = omit
	return this
}
func (this *Model) Cols(cols ...string) *Model {
	this.cols = cols
	return this
}
func (this *Model) MustCols(cols ...string) *Model {
	this.mustCols = cols
	return this
}
func (this *Model) initSqlArr() {
	this.mustCols = []string{}
	this.cols = []string{}
	this.omit = []string{}
}

//获取模型的数据，如果给了参数则将查询结果返回该参数，否则将查询的结果映射到默认的obj
func (this *Model) Info(arg ...interface{}) (err error) {
	obj := this.obj
	if len(arg) == 1 {
		obj = arg[0]
	}
	var b bool
	b, err = this.Engine().Id(this.id).Get(obj)
	if !b && err == nil {
		err = DO_FAIL
	}
	return
}

//获取模型的数据，如果给了参数则将查询结果返回该参数，否则将查询的结果映射到默认的obj
//不指定id查询，可以多个条件
func (this *Model) Infoi(arg ...interface{}) (err error) {
	obj := this.obj
	if len(arg) == 1 {
		obj = arg[0]
	}
	var b bool
	b, err = this.Engine().Get(obj)
	if !b && err == nil {
		err = DO_FAIL
	}
	return
}

//根据当前环境获取数据库session
func (this *Model) getsession() *xorm.Session {
	var sess *xorm.Session
	if this.session != nil {
		sess = this.session.Cols()
	} else {
		sess = this.Engine().Cols()
	}
	if len(this.mustCols) > 0 {
		sess.MustCols(this.mustCols...)
	}
	if len(this.cols) > 0 {
		sess.Cols(this.cols...)
	}
	if len(this.omit) > 0 {
		sess.Omit(this.omit...)
	}
	return sess
}

//添加数据，如果使用了TransactionEngine,则使用session操作，默认使用Engine()
//可以在add前调用回调修改参数
func (this *Model) Add() (err error) {
	var aff int64
	sess := this.getsession()
	aff, err = sess.InsertOne(this.obj)
	this.initSqlArr()
	if aff == 1 && err == nil {
		return nil
	}
	if aff == 0 && err == nil {
		err = DO_INVALID
	}
	return err
}

//修改数据，如果使用了TransactionEngine,则使用session操作，默认使用Engine()
func (this *Model) Mod() (err error) {
	var aff int64
	sess := this.getsession()
	aff, err = sess.Id(this.id).Update(this.obj)
	this.initSqlArr()
	if aff == 1 && err == nil {
		return nil
	}
	if aff == 0 && err == nil {
		err = DO_INVALID
	}
	return err
}

//删除数据，如果使用了TransactionEngine,则使用session操作，默认使用Engine()
func (this *Model) Del() (err error) {
	var aff int64
	sess := this.getsession()
	aff, err = sess.Id(this.id).Delete(this.obj)
	this.initSqlArr()
	if aff == 1 && err == nil {
		return nil
	}
	if aff == 0 && err == nil {
		err = DO_INVALID
	}
	return err
}

//计算分页
//参数：query=查询参数，obj=查询对象，current=当前页码，limit=每一页数量
//返回：p=分页数据，start=sql查询limit起点
func (this *Model) Page(query *xorm.Session, obj interface{}, current, limit int) (p *Page, start int) {
	totle, err := query.Clone().Count(obj)
	if err != nil {
		return
	}
	p = new(Page)
	p.Pages = int(math.Ceil(float64(totle) / float64(limit))) //计算出页码
	if current <= 1 {                                         //页首不越界
		current = 1
	}
	if current >= p.Pages { //页末不越界
		current = p.Pages
	}
	p.Totle = int(totle)
	p.Limit = limit
	p.Current = current
	start = (current - 1) * limit
	return
}

//实例化使用map保存数据的engine
func (this *Model) Mysql(engine *xorm.Engine, table string) *Mysql {
	return Mysql{Engine: engine, table: table}.New()
}

//------------------------------------------------------------------------------

//数据库会话，当使用事务的时候比用
type Session struct {
	xorm.Session      //继承xorm的session
	isSuccess    bool //是否提交成功
}

func NewSession(eg *xorm.Engine) *Session {
	o := new(Session)
	o.Engine = eg
	o.Init()
	return o
}

//自动判断是否事务提交成功，必须是在root下生成的事务才有效
func (this *Session) autoCommit(isRoot bool) error {
	if isRoot && this.isSuccess {
		return this.Commit()
	}
	return nil
}

//会话处理成功
func (this *Session) Success() {
	this.isSuccess = true
}

package xorm

import (
	"context"
	"database/sql"
	"io"
	"reflect"
	"time"

	"github.com/go-xorm/builder"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

// ORM is wrapper interface for wizard.Xorm
type ORM interface {
	ReadOnly(Identifier, bool)
	IsReadOnly(Identifier) bool
	SetAutoTransaction(Identifier, bool)
	IsAutoTransaction(Identifier) bool

	Master(interface{}) Engine
	MasterByKey(interface{}, interface{}) Engine
	Masters(interface{}) []Engine
	Slave(interface{}) Engine
	SlaveByKey(interface{}, interface{}) Engine
	Slaves(interface{}) []Engine

	Get(interface{}, func(Session) (bool, error)) (bool, error)
	Find(interface{}, func(Session) error) error
	Count(interface{}, func(Session) (int64, error)) (int64, error)
	Insert(Identifier, interface{}, func(Session) (int64, error)) (int64, error)
	Update(Identifier, interface{}, func(Session) (int64, error)) (int64, error)
	FindParallel(interface{}, interface{}, string, ...interface{}) error
	FindParallelByCondition(interface{}, FindCondition) error
	CountParallelByCondition(interface{}, FindCondition) ([]int64, error)
	UpdateParallelByCondition(interface{}, UpdateCondition) (int64, error)
	GetUsingMaster(Identifier, interface{}, func(Session) (bool, error)) (bool, error)
	FindUsingMaster(Identifier, interface{}, func(Session) error) error
	CountUsingMaster(Identifier, interface{}, func(Session) (int64, error)) (int64, error)

	NewMasterSession(interface{}) (Session, error)

	UseMasterSession(Identifier, interface{}) (Session, error)
	UseMasterSessionByKey(Identifier, interface{}, interface{}) (Session, error)
	UseSlaveSession(Identifier, interface{}) (Session, error)
	UseSlaveSessionByKey(Identifier, interface{}, interface{}) (Session, error)
	UseAllMasterSessions(Identifier, interface{}) ([]Session, error)

	ForceNewTransaction(interface{}) (Session, error)
	Transaction(Identifier, interface{}) (Session, error)
	TransactionByKey(Identifier, interface{}, interface{}) (Session, error)
	AutoTransaction(Identifier, interface{}, Session) error
	CommitAll(Identifier) error
	RollbackAll(Identifier) error
	CloseAll(Identifier)
}

// Session is interface for xorm.Session.
// (xorm v0.7.1)
type Session interface {
	Delete(bean interface{}) (int64, error)

	Find(rowsSlicePtr interface{}, condiBean ...interface{}) error
	FindAndCount(rowsSlicePtr interface{}, condiBean ...interface{}) (int64, error)

	Get(bean interface{}) (bool, error)

	Update(bean interface{}, condiBean ...interface{}) (int64, error)

	Sql(query string, args ...interface{}) *xorm.Session
	SQL(query interface{}, args ...interface{}) *xorm.Session
	Where(query interface{}, args ...interface{}) *xorm.Session
	And(query interface{}, args ...interface{}) *xorm.Session
	Or(query interface{}, args ...interface{}) *xorm.Session
	Id(id interface{}) *xorm.Session
	ID(id interface{}) *xorm.Session
	In(column string, args ...interface{}) *xorm.Session
	NotIn(column string, args ...interface{}) *xorm.Session
	Conds() builder.Cond

	Incr(column string, arg ...interface{}) *xorm.Session
	Decr(column string, arg ...interface{}) *xorm.Session
	SetExpr(column string, expression string) *xorm.Session
	Select(str string) *xorm.Session
	Cols(columns ...string) *xorm.Session
	AllCols() *xorm.Session
	MustCols(columns ...string) *xorm.Session
	UseBool(columns ...string) *xorm.Session
	Distinct(columns ...string) *xorm.Session
	Omit(columns ...string) *xorm.Session
	Nullable(columns ...string) *xorm.Session
	NoAutoTime() *xorm.Session

	Count(bean ...interface{}) (int64, error)
	Sum(bean interface{}, columnName string) (res float64, err error)
	SumInt(bean interface{}, columnName string) (res int64, err error)
	Sums(bean interface{}, columnNames ...string) ([]float64, error)
	SumsInt(bean interface{}, columnNames ...string) ([]int64, error)

	Clone() *xorm.Session
	Init()
	Close()
	ContextCache(context xorm.ContextCache) *xorm.Session
	IsClosed() bool
	Prepare() *xorm.Session
	Before(closures func(interface{})) *xorm.Session
	After(closures func(interface{})) *xorm.Session
	Table(tableNameOrBean interface{}) *xorm.Session
	Alias(alias string) *xorm.Session
	NoCascade() *xorm.Session
	ForUpdate() *xorm.Session
	NoAutoCondition(no ...bool) *xorm.Session
	Limit(limit int, start ...int) *xorm.Session
	OrderBy(order string) *xorm.Session
	Desc(colNames ...string) *xorm.Session
	Asc(colNames ...string) *xorm.Session
	StoreEngine(storeEngine string) *xorm.Session
	Charset(charset string) *xorm.Session
	Cascade(trueOrFalse ...bool) *xorm.Session
	NoCache() *xorm.Session
	Join(joinOperator string, tablename interface{}, condition string, args ...interface{}) *xorm.Session
	GroupBy(keys string) *xorm.Session
	Having(conditions string) *xorm.Session
	DB() *core.DB
	LastSQL() (string, []interface{})
	Unscoped() *xorm.Session

	Rows(bean interface{}) (*xorm.Rows, error)
	Iterate(bean interface{}, fun xorm.IterFunc) error
	BufferSize(size int) *xorm.Session

	PingContext(ctx context.Context) error

	Insert(beans ...interface{}) (int64, error)
	InsertMulti(rowsSlicePtr interface{}) (int64, error)
	InsertOne(bean interface{}) (int64, error)

	Query(sqlorArgs ...interface{}) ([]map[string][]byte, error)
	QueryString(sqlorArgs ...interface{}) ([]map[string]string, error)
	QuerySliceString(sqlorArgs ...interface{}) ([][]string, error)
	QueryInterface(sqlorArgs ...interface{}) ([]map[string]interface{}, error)

	Exist(bean ...interface{}) (bool, error)

	Begin() error
	Rollback() error
	Commit() error

	Exec(sqlorArgs ...interface{}) (sql.Result, error)

	Ping() error
	CreateTable(bean interface{}) error
	CreateIndexes(bean interface{}) error
	CreateUniques(bean interface{}) error
	DropIndexes(bean interface{}) error
	DropTable(beanOrTableName interface{}) error
	IsTableExist(beanOrTableName interface{}) (bool, error)
	IsTableEmpty(bean interface{}) (bool, error)
	Sync2(beans ...interface{}) error
}

// Engine is interface for xorm.Engine
type Engine interface {
	Transaction(f func(*xorm.Session) (interface{}, error)) (interface{}, error)

	PingContext(ctx context.Context) error

	SetCacher(tableName string, cacher core.Cacher)
	GetCacher(tableName string) core.Cacher
	BufferSize(size int) *xorm.Session
	CondDeleted(colName string) builder.Cond
	ShowSQL(show ...bool)
	ShowExecTime(show ...bool)
	Logger() core.ILogger
	SetLogger(logger core.ILogger)
	SetLogLevel(level core.LogLevel)
	SetDisableGlobalCache(disable bool)
	DriverName() string
	DataSourceName() string
	SetMapper(mapper core.IMapper)
	SetTableMapper(mapper core.IMapper)
	SetColumnMapper(mapper core.IMapper)
	SupportInsertMany() bool
	QuoteStr() string
	Quote(value string) string
	QuoteTo(buf *builder.StringBuilder, value string)
	SqlType(c *core.Column) string
	SQLType(c *core.Column) string
	AutoIncrStr() string
	SetConnMaxLifetime(d time.Duration)
	SetMaxOpenConns(conns int)
	SetMaxIdleConns(conns int)
	SetDefaultCacher(cacher core.Cacher)
	GetDefaultCacher() core.Cacher
	NoCache() *xorm.Session
	NoCascade() *xorm.Session
	MapCacher(bean interface{}, cacher core.Cacher) error
	NewDB() (*core.DB, error)
	DB() *core.DB
	Dialect() core.Dialect
	NewSession() *xorm.Session
	Close() error
	Ping() error
	Sql(querystring string, args ...interface{}) *xorm.Session
	SQL(query interface{}, args ...interface{}) *xorm.Session
	NoAutoTime() *xorm.Session
	NoAutoCondition(no ...bool) *xorm.Session
	DBMetas() ([]*core.Table, error)
	DumpAllToFile(fp string, tp ...core.DbType) error
	DumpAll(w io.Writer, tp ...core.DbType) error
	DumpTablesToFile(tables []*core.Table, fp string, tp ...core.DbType) error
	DumpTables(tables []*core.Table, w io.Writer, tp ...core.DbType) error
	Cascade(trueOrFalse ...bool) *xorm.Session
	Where(query interface{}, args ...interface{}) *xorm.Session
	Id(id interface{}) *xorm.Session
	ID(id interface{}) *xorm.Session
	Before(closures func(interface{})) *xorm.Session
	After(closures func(interface{})) *xorm.Session
	Charset(charset string) *xorm.Session
	StoreEngine(storeEngine string) *xorm.Session
	Distinct(columns ...string) *xorm.Session
	Select(str string) *xorm.Session
	Cols(columns ...string) *xorm.Session
	AllCols() *xorm.Session
	MustCols(columns ...string) *xorm.Session
	UseBool(columns ...string) *xorm.Session
	Omit(columns ...string) *xorm.Session
	Nullable(columns ...string) *xorm.Session
	In(column string, args ...interface{}) *xorm.Session
	NotIn(column string, args ...interface{}) *xorm.Session
	Incr(column string, arg ...interface{}) *xorm.Session
	Decr(column string, arg ...interface{}) *xorm.Session
	SetExpr(column string, expression string) *xorm.Session
	Table(tableNameOrBean interface{}) *xorm.Session
	Alias(alias string) *xorm.Session
	Limit(limit int, start ...int) *xorm.Session
	Desc(colNames ...string) *xorm.Session
	Asc(colNames ...string) *xorm.Session
	OrderBy(order string) *xorm.Session
	Prepare() *xorm.Session
	Join(joinOperator string, tablename interface{}, condition string, args ...interface{}) *xorm.Session
	GroupBy(keys string) *xorm.Session
	Having(conditions string) *xorm.Session
	UnMapType(t reflect.Type)
	GobRegister(v interface{}) *xorm.Engine
	TableInfo(bean interface{}) *xorm.Table
	IsTableEmpty(bean interface{}) (bool, error)
	IsTableExist(beanOrTableName interface{}) (bool, error)
	IdOf(bean interface{}) core.PK
	IDOf(bean interface{}) core.PK
	IdOfV(rv reflect.Value) core.PK
	IDOfV(rv reflect.Value) core.PK
	CreateIndexes(bean interface{}) error
	CreateUniques(bean interface{}) error
	ClearCacheBean(bean interface{}, id string) error
	ClearCache(beans ...interface{}) error
	Sync(beans ...interface{}) error
	Sync2(beans ...interface{}) error
	CreateTables(beans ...interface{}) error
	DropTables(beans ...interface{}) error
	DropIndexes(bean interface{}) error
	Exec(sqlorArgs ...interface{}) (sql.Result, error)
	Query(sqlorArgs ...interface{}) (resultsSlice []map[string][]byte, err error)
	QueryString(sqlorArgs ...interface{}) ([]map[string]string, error)
	QueryInterface(sqlorArgs ...interface{}) ([]map[string]interface{}, error)
	Insert(beans ...interface{}) (int64, error)
	InsertOne(bean interface{}) (int64, error)
	Update(bean interface{}, condiBeans ...interface{}) (int64, error)
	Delete(bean interface{}) (int64, error)
	Get(bean interface{}) (bool, error)
	Exist(bean ...interface{}) (bool, error)
	Find(beans interface{}, condiBeans ...interface{}) error
	FindAndCount(rowsSlicePtr interface{}, condiBean ...interface{}) (int64, error)
	Iterate(bean interface{}, fun xorm.IterFunc) error
	Rows(bean interface{}) (*xorm.Rows, error)
	Count(bean ...interface{}) (int64, error)
	Sum(bean interface{}, colName string) (float64, error)
	SumInt(bean interface{}, colName string) (int64, error)
	Sums(bean interface{}, colNames ...string) ([]float64, error)
	SumsInt(bean interface{}, colNames ...string) ([]int64, error)
	ImportFile(ddlPath string) ([]sql.Result, error)
	Import(r io.Reader) ([]sql.Result, error)
	GetColumnMapper() core.IMapper
	GetTableMapper() core.IMapper
	GetTZLocation() *time.Location
	SetTZLocation(tz *time.Location)
	GetTZDatabase() *time.Location
	SetTZDatabase(tz *time.Location)
	SetSchema(schema string)
	Unscoped() *xorm.Session

	Clone() (*xorm.Engine, error)

	TableName(bean interface{}, includeSchema ...bool) string
}

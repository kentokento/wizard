package xorm

import (
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

// Session is interface for xorm.Interface (xorm.Session)
type Session interface {
	xorm.Interface

	Init()
	Close()
	Begin() error
	Rollback() error
	Commit() error

	And(query interface{}, args ...interface{}) *xorm.Session
}

// Engine is interface for xorm.EngineInterface (xorm.Engine)
type Engine interface {
	xorm.EngineInterface
}

package geeorm

import (
	"errors"
	"geeorm/session"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "gee.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

func TestNewEngine(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
}

//func OpenDB(t *testing.T) *Engine {
//	t.Helper()
//	engine, err := NewEngine("sqlite3", "gee.db")
//	if err != nil {
//		t.Fatal("无法连接")
//	}
//	return engine
//}

type User struct {
	Name string `geerom:"PRIMARY KEY"`
	Age  int
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactioncommit(t)
	})
}
func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()

	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result []interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{
			Name: "Tom",
			Age:  18,
		})
		//err != nil触发回滚
		return nil, errors.New("error")
	})
	if err == nil || s.HasTable() {
		t.Fatal("不能回滚")
	}
}
func transactioncommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result []interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{
			Name: "Tom",
			Age:  18,
		})
		//err != nil触发回滚
		return
	})
	u := &User{}
	_ = s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("事务不能提交")
	}
}

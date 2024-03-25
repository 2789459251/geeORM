package geeorm

import (
	"database/sql"
	"fmt"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/session"
	"strings"
)

// Engine is the main struct of geeorm, manages all db sessions and transactions.
type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// NewEngine create a instance of Engine
// connect database and ping it to test whether it's alive
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	// Send a ping to make sure the database connection is alive.
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	// make sure the specific dialect exists
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

// Close database connection
func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// NewSession creates a new session for next operations
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}

type TxFunc func(s *session.Session) ([]interface{}, error)

// 事务执行函数
func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		//需要修护，事务出问题了
		if p := recover(); p != nil {
			//p表示的是捕获的panic的值
			_ = s.Rollback()
			//业务出问题了，panic一下
			panic(p)
		} else if err != nil {
			//事务有问题，回滚
			_ = s.Rollback()
		} else {
			err = s.Commit()
		}
	}()
	//f就是事务
	return f(s)
}
func difference(a []string, b []string) (diff []string) {
	//a有b没有的
	mapb := make(map[string]bool)
	for _, v := range b {
		mapb[v] = true
	}
	for _, v := range a {
		if _, ok := mapb[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}
func (engine Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(s *session.Session) (result []interface{}, err error) {
		if !s.Model(value).HasTable() {
			log.Infof("%s表还未创建", s.RefTable().Name)
			return nil, s.CreateTable()
		}
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("added cols %v,delete cols %v", addCols, delCols)
		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %v ADD COLUMN %s %s", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}
		if len(delCols) == 0 {
			return
		}
		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ",")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s ;", tmp, fieldStr, s.RefTable().Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}

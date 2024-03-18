package geeorm

/*用户交互*/
import (
	"database/sql"
	"geeORM/geeorm/mylog"
	"geeORM/geeorm/session"
)

type Engine struct {
	db *sql.DB
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		mylog.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		mylog.Error(err)
		return
	}
	e = &Engine{db: db}
	mylog.Info("数据库连接成功")
	return
}
func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		mylog.Error("关闭数据库失败")
	}
	mylog.Info("关闭数据库成功")
}
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db)
}

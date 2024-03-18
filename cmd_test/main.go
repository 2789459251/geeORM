package main

import (
	"geeORM/geeorm"
	"geeORM/geeorm/mylog"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	engine, _ := geeorm.NewEngine("sqlite3", "gee.db")
	defer engine.Close()

	s := engine.NewSession()

	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User (Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE User (`Name text`);").Exec()
	result, _ := s.Raw("INSERT INTO User (`Name`) values (?),(?)", "tom", "sim").Exec()
	count, _ := result.RowsAffected()
	mylog.Infof("操作成功，有%d行受影响", count)
}

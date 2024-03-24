package session

import (
	"fmt"
	"geeorm/log"
	"testing"
)

type Account struct {
	Id       int `geeorm:"primary key"`
	Password string
}

func (account *Account) BeforeInsert(s *Session) error {
	log.Info("before insert", account)
	account.Id += 1000
	return nil
}
func (account *Account) AfterQuery(s *Session) error {
	log.Info("after query", account)
	account.Password = "*********"
	return nil
}
func TestSession_CallMethod(t *testing.T) {
	s := NewSession().Model(&Account{})
	_ = s.DropTable()
	_ = s.CreateTable()
	_, _ = s.Insert(&Account{
		Id:       1,
		Password: "123456",
	}, &Account{
		Id:       2,
		Password: "qwerty",
	})
	u := &Account{}
	err := s.First(u)
	fmt.Println(u)
	if err != nil || u.Id != 1001 || u.Password != "*********" {
		t.Fatal("Failed to call hooks after query, got", u)
	}
}

package session

import "testing"

var (
	user1 = &User{"Tom", 19}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewSession().Model(&User{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("record测试初始化失败")
	}
	return s
}
func TestSession_Find(t *testing.T) {
	s := testRecordInit(t)
	var Users []User
	if err := s.Find(&Users); err != nil || len(Users) != 2 {
		t.Fatal("查询失败")
	}
}
func TestSession_Insert(t *testing.T) {
	s := testRecordInit(t)
	affect, err := s.Insert(user3)
	if err != nil || affect != 1 {
		t.Fatal("插入失败")
	}
}

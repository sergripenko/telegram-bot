package models

import "github.com/astaxie/beego/orm"

type Users struct {
	Id        int    `orm:"column(id);pk;auto" json:"id"`
	UserName  string `orm:"column(user_name)" json:"user_name"`
	FirstName string `orm:"column(first_name)" json:"first_name"`
	LastName  string `orm:"column(last_name)" json:"last_name"`
	ChatID    int    `orm:"column(chat_id)" json:"chat_id"`
	Base
}

func (t *Users) TableName() string {
	return "users"
}

func init() {
	orm.RegisterModel(new(Users))
}

// AddUsers insert a new User into database and returns
// last inserted Id on success.
func AddUser(m *Users) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

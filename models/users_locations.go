package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/labstack/gommon/log"
)

type UsersLocations struct {
	Id        int    `orm:"column(id);pk;auto" json:"id"`
	User      *Users `orm:"column(user_id);rel(fk)" json:"user"`
	Latitude  int    `orm:"column(latitude)" json:"latitude"`
	Longitude int    `orm:"column(longitude)" json:"longitude"`
	City      string `orm:"column(city)" json:"city"`
	Base
}

func (t *UsersLocations) TableName() string {
	return "users_locations"
}

func init() {
	orm.RegisterModel(new(UsersLocations))
}

func AddUsersLocation(lat, long, chatId int) {
	log.Info("AddUsersLocation")
	o := orm.NewOrm()
	var user Users
	o.QueryTable(new(Users)).Filter("chat_id", chatId).One(&user)
	usersLoc := &UsersLocations{
		User:      &user,
		Latitude:  lat,
		Longitude: long,
	}
	_, err := o.Insert(usersLoc)
	log.Error(err)
}

func IfLocationExist(chatId int) bool {
	log.Info("IfLocationExist")
	o := orm.NewOrm()
	var user Users
	o.QueryTable(new(Users)).Filter("chat_id", chatId).One(&user)
	exist := o.QueryTable(new(UsersLocations)).Filter("user_id", user.Id).Exist()
	return exist
}

func AddUserLocationCity(chatId int, city string) {
	log.Info("AddUserLocationCity")
	o := orm.NewOrm()
	var user Users
	o.QueryTable(new(Users)).Filter("chat_id", chatId).One(&user)
	var loc UsersLocations
	o.QueryTable(new(UsersLocations)).Filter("user_id", user.Id).One(&loc)
	loc.City = city
	o.Update(&loc)
}

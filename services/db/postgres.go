package db

import (
	"os"

	"github.com/labstack/gommon/log"

	"github.com/astaxie/beego/orm"
	"gopkg.in/yaml.v2"
)

type DBconf struct {
	Dialect  string `yaml:"dialect"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	DBName   string `yaml:"dbname"`
	Password string `yaml:"password"`
	SSlMode  string `yaml:"sslmode"`
}

func getDBConfig() (*DBconf, error) {
	f, err := os.Open("db_conf.yml")

	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg DBconf
	decoder := yaml.NewDecoder(f)

	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func InitDB() {
	conf, err := getDBConfig()
	if err != nil {
		log.Error(err)
	}
	orm.RegisterDriver(conf.Dialect, orm.DRPostgres)
	dbparams := "user=" + conf.User +
		" password=" + conf.Password +
		" host=" + conf.Host +
		" port=" + conf.Port +
		" dbname=" + conf.DBName +
		" sslmode=" + conf.SSlMode
	orm.RegisterDataBase("default", "postgres", dbparams)
}

package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"path"
	"runtime"
)

type Database interface {
	GetRecord(instance interface{}, query string, params ...interface{}) error
	GetRecords(instances interface{}, query string, params ...interface{}) error
	Exec(query string, params ...interface{}) error
	Close()
}

type dbase struct {
	connection *sqlx.DB // handle database
}

func getAbsPath() string {
	_, filename, _, _ := runtime.Caller(1)
	pwd := path.Join(path.Dir(filename))
	return pwd
}

func initEnvConfigs() {
	dir := getAbsPath()
	viper.SetConfigName("config")
	viper.AddConfigPath(dir + "/../../config")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func getDBConfigs() string {
	cmd := "host=" + viper.GetString("db.host")
	cmd += " user=" + viper.GetString("db.user")
	cmd += " password=" + viper.GetString("db.password")
	cmd += " dbname=" + viper.GetString("db.dbname")
	cmd += " sslmode=disable"
	return cmd
}

// open new connect to Postgres database
func Open() Database {
	initEnvConfigs()
	dbConfigs := getDBConfigs()

	db, err := sqlx.Connect("postgres", dbConfigs)
	if err != nil {
		log.Fatal(err)
	}

	return &dbase{db}
}

// GetRecord from the database
func (s *dbase) GetRecord(instance interface{}, query string, params ...interface{}) error {
	err := s.connection.Get(instance, query, params...)

	return err
}

// GetRecord from the database
func (s *dbase) Exec(query string, params ...interface{}) error {
	_, err := s.connection.Exec(query, params...)

	return err
}

func (s *dbase) GetRecords(instances interface{}, query string, params ...interface{}) error {
	return s.connection.Select(instances, query, params...)
}

func (s *dbase) Close() {
	s.connection.Close()
}

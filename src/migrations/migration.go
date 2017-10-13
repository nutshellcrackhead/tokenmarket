package main

import (
	"fmt"
	"log"
	"path"
	"runtime"
	// "path/filepath"
	"database/sql"
	"github.com/tanel/dbmigrate"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func initEnvConfigs() {
	dir := getAbsPath()
	viper.SetConfigName("config")
	viper.AddConfigPath(dir + "/../config/")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func getAbsPath() string {
	_, filename, _, _ := runtime.Caller(1)
	pwd := path.Join(path.Dir(filename))
	return pwd
}

func getDBConfigs() string {
	cmd := "host=" + viper.GetString("db.host")
	cmd += " user=" + viper.GetString("db.user")
	cmd += " password=" + viper.GetString("db.password")
	cmd += " dbname=" + viper.GetString("db.dbname")
	cmd += " sslmode=disable"
	//fmt.Println(cmd)
	return cmd
}

func main() {
	initEnvConfigs()
	dbConfigs := getDBConfigs()

	db, err := sql.Open("postgres", dbConfigs)
	if err != nil {
		log.Fatal(err)
	}

	migrationsPath := getMigrationsPath()
	//fmt.Println(migrationsPath)
	if err := dbmigrate.Run(db, migrationsPath); err != nil {
		log.Println(err)
	}
}

func getMigrationsPath() string {
	return path.Join(getAbsPath(), "db", "migrate")
}

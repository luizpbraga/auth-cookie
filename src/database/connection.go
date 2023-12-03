package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/luizpbraga/auth-cookie/src/database/models"
)

var Db *sql.DB

func Connect() (*sql.DB, error) {
	err := godotenv.Load("./.env")

	if err != nil {
		log.Fatal("load end error: ", err)
		return nil, err
	}

	user := os.Getenv("DBUSER")
	dbName := os.Getenv("DBNAME")
	address := os.Getenv("DBADDRESS")
	password := os.Getenv("DBPASSWORD")

	fmt.Print(user+dbName+address+password, "\n")

	cfg := mysql.Config{
		AllowNativePasswords: true,
		Net:                  "tcp",
		User:                 user,
		Addr:                 address,
		Passwd:               password,
		DBName:               dbName,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal("open err:", err)
		return db, err
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Not Connected: ", err)
		return db, err
	}

	if err := CreateUserTable(db); err != nil {
		log.Fatal("Can not create User Table: ", err)
		return db, err
	}

	Db = db

	return db, nil
}

func FindUserByEmail(email string) (*models.User, error) {
	query := "select id, name, email, password from User where email = (?) limit 1"
	row := Db.QueryRow(query, email)
	user := new(models.User)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	return user, err
}

func FindUserById(id string) (*models.User, error) {
	query := "select id, name, email, password from User where id = (?) limit 1"
	row := Db.QueryRow(query, id)
	user := new(models.User)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	return user, err
}

func InserUserTable(user *models.User) error {
	query := "insert into User (name, email, password) values (?, ?, ?)"
	_, err := Db.Exec(query, user.Name, user.Email, user.Password)
	return err
}

func CreateUserTable(db *sql.DB) error {
	query := `
create table if not exists User(
  id          bigint unsigned not null auto_increment,
  name        longtext, 
  email       longtext unique, 
  password    longtext,
  primary key (id)
);
  `
	_, err := db.Exec(query)
	return err
}

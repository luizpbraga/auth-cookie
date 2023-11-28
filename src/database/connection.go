package database

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/luizpbraga/auth-cookie/src/database/models"
)

var Db *sql.DB

func Connect() (*sql.DB, error) {
	cfg := mysql.Config{
		Net:                  "tcp",
		User:                 "test",
		Addr:                 "localhost",
		Passwd:               "test",
		DBName:               "Auth",
		AllowNativePasswords: true,
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

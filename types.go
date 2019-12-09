package main

import "github.com/go-sql-driver/mysql"

// User database object
type User struct {
	userid             int            `db:"userid"`
	handle             string         `db:"handle"`
	discord            string         `db:"discord"`
	password           string         `db:"password"`
	email              string         `db:"email"`
	img                string         `db:"img"`
	shortBio           string         `db:"shortBio"`
	bio                string         `db:"bio"`
	firstName          string         `db:"firstName"`
	lastName           string         `db:"lastName"`
	accountCreation    mysql.NullTime `db:"accountCreation"`
	bday               mysql.NullTime `db:"bday"`
	position           int            `db:"position"`
	referredby         int            `db:"referredby"`
	rank               int            `db:"rank"`
	status             int            `db:"status"`
	message            string         `db:"message"`
	approvedby         int            `db:"approvedby"`
	approveddate       mysql.NullTime `db:"approveddate"`
	newsletter         int            `db:"newsletter"`
	terms              int            `db:"terms"`
	blogupdatesemail   int            `db:"blogupdatesemail"`
	leadermessageemail int            `db:"leadermessageemail"`
	newfeatureemail    int            `db:"newfeatureemail"`
	newmission         int            `db:"newmission"`
}

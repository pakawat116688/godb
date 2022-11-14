package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type employee struct {
	id int
	name string
	salary int
	tel string
	status int
}

func main() {
	os.Remove("./mydata.db")
	println("Hello")
	mydb, err := sql.Open("sqlite3", "./mydata.db")
	if err != nil {
		println("Error Can not Open database....")
		panic(err)
	}
	defer mydb.Close()

	err = create_table(mydb)
	if err != nil {
		panic(err)
	}

	err = insert_data(mydb,"Doraemon", 80000, "06x-xxx-xxxx", 0)
	if err != nil {
		panic(err)
	}

	err = insert_data(mydb,"Nobita", 90000, "061-xxx-xxxx", 0)
	if err != nil {
		panic(err)
	}

	err = insert_data(mydb,"Sisuga", 50000, "062-xxx-xxxx", 0)
	if err != nil {
		panic(err)
	}

	err = insert_data(mydb,"Giant", 70000, "063-xxx-xxxx", 0)
	if err != nil {
		panic(err)
	}

	err = insert_data(mydb,"Sunio", 200000, "099-xxx-xxxx", 1)
	if err != nil {
		panic(err)
	}

	data, err := getData(mydb)
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		fmt.Printf("%#v \n",value)
	}

	err = update_table(mydb, "Nobi Nobita", 2)
	if err != nil {
		panic(err)
	}

	getid, err := getdataById(mydb, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println("Get ID 2 --> ",*getid)


	// delete one column
	err = del_dataById(mydb, 2)
	if err != nil {
		panic(err)
	}

	data, err = getData(mydb)
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		fmt.Printf("%#v \n",value)
	}

	// delete all
	err = del_data(mydb)
	if err != nil {
		panic(err)
	}

	// after delete
	println("After Delete......")
	data, err = getData(mydb)
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		fmt.Printf("%#v \n",value)
	}

}

// Create
func create_table(db *sql.DB) error {

	tx, err := db.Begin()

	if err != nil {
		return err
	}
	create_table := `CREATE TABLE "Employee" (
		"id"	INTEGER,
		"name"	TEXT,
		"salary"	NUMERIC,
		"tel"	TEXT,
		"status"	NUMERIC,
		PRIMARY KEY("id" AUTOINCREMENT)
	);`

	statement, err := tx.Prepare(create_table)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	
	println("Employee Table created")

	return nil
}

// Update
func insert_data(db *sql.DB, name string, salary int, tel string, status int) error {

	println("Insert Employee Record.....")
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := `insert into employee(name, salary, tel, status)
		values (?, ?, ?, ?)`

	state, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	_, err = state.Exec(name, salary, tel, status)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	println("Insert Data Success.....")

	return nil
}


// Read All
func getData(db *sql.DB) ([]employee, error) {

	println("Get Data From Employee Database....")
	if err := db.Ping(); err != nil {
		return nil, err
	}
	
	query := "select * from Employee"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	empployee := []employee{}

	for rows.Next() {
		emp := employee{}
		err = rows.Scan(&emp.id, &emp.name, &emp.salary, &emp.tel, &emp.status)
		if err != nil {
			return nil, err
		}
		empployee = append(empployee, emp)
	}

	return empployee, nil
}

// Read By Id
func getdataById(db *sql.DB, id int) (*employee, error)  {
	
	err := db.Ping()
	if err != nil{
		return nil, err
	}

	query := "select * from Employee where id=?"
	state ,err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	emp := employee{}
	err = state.QueryRow(id).Scan(&emp.id, &emp.name, &emp.salary, &emp.tel, &emp.status)
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func update_table(db *sql.DB, name string, id int) error {
	
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := "update Employee set name=? where id=?"
	state, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	_, err = state.Exec(name, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	fmt.Printf("Update column id %v success...",id)

	return nil
}

func del_data(db *sql.DB) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := "delete from Employee"
	state, err := tx.Exec(query)
	if err != nil {
		return err
	}

	affected, err := state.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil
	}

	if affected <= 0 {
		tx.Rollback()
		return errors.New("cannot delete")
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	fmt.Println("delete all success...")

	return nil
}

func del_dataById(db *sql.DB, id int) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := "delete from Employee where id=?"
	state, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	_, err = state.Exec(id)
	if err != nil {
		tx.Rollback()
		return nil
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	fmt.Printf("delete id %v success...\n",id)

	return nil
}
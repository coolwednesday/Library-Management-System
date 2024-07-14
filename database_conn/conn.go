package database_conn

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func NewConnection() (*sql.DB, error) {
	//Connecting to the mysql Database
	s, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/library")
	if err != nil {
		return s, err
	}

	//var title, author, name string
	//var isbn,id int

	//Creating Books Table
	/*
		_, err = s.Exec(`create table books(isbn int primary key,title text,author text)`)
		if err != nil {
			fmt.Println(err)
		}else{
			fmt.Println("Table created")
		}
	*/

	//Creating Users Table
	/*
		_, err = s.Exec(`create table users(id int primary key,name text)`)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Table users created")
		}
	*/

	//Creating lendingRecords tables
	/*
		_, err = s.Exec(`CREATE TABLE lendingRecords(actionid int auto_increment primary key ,userid int not null,bookid int not null, foreign key (userid) references users(id), foreign key (bookid) references books(isbn))`)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Lending records table created")
		}
	*/

	//Inserting Users in the Database in users table
	/*
		var u = []User{
			{"User1", 1567},
			{"User2", 8902},
			{"User3", 7903},
			{"user4", 9056},
		}
		for _, v := range u {
			_, err = s.Exec(`INSERT INTO users(id, name) VALUES (?,?)`, v.Id, v.Name)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("User" + v.Name + "inserted successfully")
			}

		}
	*/

	//Inserting Books in the Database in books table
	/*
		var b = []Book{
			{"Book1", "Author1", 12345},
			{"Book2", "Author2", 12905},
			{"Book3", "Author3", 12785},
			{"Book4", "Author4", 19905},
		}
		for _, v := range b {
			_, err := s.Exec(`INSERT INTO books (title,author,isbn) VALUES (?,?,?)`, v.title, v.author, v.isbn)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Inserted", v.title)
			}
		}
	*/

	//Inserting Lending Records into the table
	/*
		records := make(map[int]int)
		records[12905] = 1567
		records[19905] = 9056
		for key, value := range records {
			_, err = s.Exec(`INSERT INTO lendingRecords(userid,bookid)VALUES((Select id from users where id=?),(Select isbn from books where isbn=?))`, value, key)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Lending records added to table")
			}
		}
	*/

	/*Fetching All the books in Database

	list := []Book{}
	rows, err := s.Query(`SELECT title,author,isbn FROM books`)
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		err = rows.Scan(&title, &author, &isbn)
		if err != nil {
			fmt.Println(err)
		} else {
			list = append(list, Book{
				isbn,
				title,
				author,
			})
		}
	}

	for _, v := range list {
		fmt.Println(v)
	}
	*/

	//Updating an existing records
	/*
		record := Book{
			12345,
			"Book1",
			"Author5",
		}
	*/

	/*
		_, err = s.Exec(`UPDATE books SET author="Author5" where isbn=` + strconv.Itoa(record.ISBN))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Book updated")
		}
	*/
	//Deleting an existing record
	/*
		_, err = s.Exec(`DELETE FROM books WHERE ISBN=` + strconv.Itoa(record.ISBN))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Book deleted")
		}

	*/
	return s, nil
}

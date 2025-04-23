package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return
	}
}
func dbConnectInteractive() *sql.DB {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter DB Driver (e.g. mysql): ")
	driver, _ := reader.ReadString('\n')
	driver = strings.TrimSpace(driver)

	fmt.Print("Enter Username: ")
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)

	fmt.Print("Enter Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	//Initial connection without specifying a database
	db, err := sql.Open(driver, user+":"+password+"@/")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to DB server:", err)
	}

	//Fetching Databases
	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		log.Fatal("Cannot list databases:", err)
	}
	defer rows.Close()

	var databases []string
	var dbName string
	i := 1
	fmt.Println("\nAvailable Databases:")
	for rows.Next() {
		rows.Scan(&dbName)
		fmt.Printf("%d. %s\n", i, dbName)
		databases = append(databases, dbName)
		i++
	}

	fmt.Print("\nEnter the number of the database you want to use: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(databases) {
		log.Fatal("Invalid database selection.")
	}
	selectedDB := databases[index-1]

	//Reconnect to the selected database
	finalDB, err := sql.Open(driver, user+":"+password+"@/"+selectedDB)
	if err != nil {
		log.Fatal(err)
	}
	err = finalDB.Ping()
	if err != nil {
		log.Fatal("Failed to connect to selected DB:", err)
	}

	fmt.Printf("\nConnected to database: %s\n", selectedDB)
	return finalDB
}

//func readCities(db *sql.DB) {
//	fmt.Println("===== Cities =====")
//	rows, err := db.Query("SELECT ID, Name, CountryCode, Population FROM city")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var id int
//		var name, code string
//		var population int
//		err := rows.Scan(&id, &name, &code, &population)
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Printf("%d | %-20s | %s | %d\n", id, name, code, population)
//	}
//}

//func readCountries(db *sql.DB) {
//	fmt.Println("\n===== Countries =====")
//	rows, err := db.Query("SELECT Code, Name, Continent, Population FROM country")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var code, name, continent string
//		var population int
//		err := rows.Scan(&code, &name, &continent, &population)
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Printf("%s | %-30s | %-10s | %d\n", code, name, continent, population)
//	}
//}

//	func readLanguages(db *sql.DB) {
//		fmt.Println("\n===== Languages =====")
//		rows, err := db.Query("SELECT CountryCode, Language, IsOfficial FROM countrylanguage")
//		if err != nil {
//			log.Fatal(err)
//		}
//		defer rows.Close()
//
//		for rows.Next() {
//			var code, language, official string
//			err := rows.Scan(&code, &language, &official)
//			if err != nil {
//				log.Fatal(err)
//			}
//			fmt.Printf("%s | %-20s | Official: %s\n", code, language, official)
//		}
//	}
func listTables(db *sql.DB) []string {
	fmt.Println("Available tables:")
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tables []string
	var tableName string
	i := 1
	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d. %s\n", i, tableName)
		tables = append(tables, tableName)
		i++
	}
	fmt.Println(strings.Repeat("-", 30))
	return tables
}
func readTableByName(db *sql.DB, tableName string) {
	// Displays the content of any table from the database in a dynamic way.
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error reading table '%s': %v\n", tableName, err)
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of rows is : ", len(columns))
	printRows(tableName, columns, rows)

}
func printRows(tableName string, columns []string, rows *sql.Rows) {

	values := make([]interface{}, len(columns))    // It will contain the values that will be returned from the table.
	valuePtrs := make([]interface{}, len(columns)) //These are indicators of these values so that we can pass them
	// to rows.Scan.

	fmt.Println("\n===== " + strings.ToUpper(tableName) + " =====")

	for i := range columns { // Formating
		fmt.Printf("%-20s", columns[i])
	}

	fmt.Println("\n" + strings.Repeat("-", 20*len(columns)))

	for rows.Next() {
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal(err)
		}

		for _, val := range values {
			var v string

			if b, ok := val.([]byte); ok {
				v = string(b)
			} else if val == nil {
				v = "NULL"
			} else {
				v = fmt.Sprintf("%v", val)
			}
			fmt.Printf("%-20s", v)
		}
		fmt.Println()
	}
}
func updateRowInteractive(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter table name: ")
	table, _ := reader.ReadString('\n')
	table = strings.TrimSpace(table)

	fmt.Print("Enter column to identify row (e.g. id): ")
	idColumn, _ := reader.ReadString('\n')
	idColumn = strings.TrimSpace(idColumn)

	fmt.Print("Enter value of that column (e.g. 5): ")
	idValue, _ := reader.ReadString('\n')
	idValue = strings.TrimSpace(idValue)

	fmt.Print("Enter column to update: ")
	updateColumn, _ := reader.ReadString('\n')
	updateColumn = strings.TrimSpace(updateColumn)

	fmt.Print("Enter new value: ")
	newValue, _ := reader.ReadString('\n')
	newValue = strings.TrimSpace(newValue)

	query := fmt.Sprintf("UPDATE %s SET %s = ? WHERE %s = ?", table, updateColumn, idColumn)
	result, err := db.Exec(query, newValue, idValue)
	if err != nil {
		fmt.Println("Error updating row:", err)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Rows affected:", rows)
}
func deleteRowInteractive(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter table name: ")
	table, _ := reader.ReadString('\n')
	table = strings.TrimSpace(table)

	fmt.Print("Enter column to identify row (e.g. id): ")
	idColumn, _ := reader.ReadString('\n')
	idColumn = strings.TrimSpace(idColumn)

	fmt.Print("Enter value to delete: ")
	idValue, _ := reader.ReadString('\n')
	idValue = strings.TrimSpace(idValue)

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", table, idColumn)
	result, err := db.Exec(query, idValue)
	if err != nil {
		fmt.Println("Error deleting row:", err)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Rows affected:", rows)
}
func selectQueryManual(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter SELECT query: ")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	printRows("Query Result", columns, rows)
}

// INSERT INTO city (Name, CountryCode, District, Population) VALUES ('NewCity', 'USA', 'TestState', 12345);
func insertQueryManual(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter INSERT query: ")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	result, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error executing INSERT:", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Insert executed, but couldn't retrieve rows affected.")
		return
	}

	fmt.Printf("Insert successful. Rows affected: %d\n", rowsAffected)
}

func main() {
	db := dbConnectInteractive()
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		fmt.Println("\n===== MENU =====")
		fmt.Println("1. Show Tables")
		fmt.Println("2. Show Table by Number")
		fmt.Println("3. Update Row")
		fmt.Println("4. Delete Row")
		fmt.Println("5. Run SELECT Query")
		fmt.Println("6. Run INSERT Query")
		fmt.Println("0. Exit")
		fmt.Print("Enter your choice: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid choice, please try again")
			continue
		}

		switch choice {
		case 0:
			fmt.Println("Exiting...")
			return
		case 1:
			listTables(db)
			fmt.Println("Press ENTER to return...")
			reader.ReadString('\n')
		case 2:
			tables := listTables(db)
			if len(tables) == 0 {
				fmt.Println("No tables found.")
				break
			}
			fmt.Print("Enter table number: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			index, err := strconv.Atoi(input)
			if err != nil || index < 1 || index > len(tables) {
				fmt.Println("Invalid table number.")
				break
			}
			tableName := tables[index-1]
			readTableByName(db, tableName)
			fmt.Println("\nPress ENTER to return to menu...")
			reader.ReadString('\n')
		case 3:
			updateRowInteractive(db)
			fmt.Println("\nPress ENTER to return to menu...")
			reader.ReadString('\n')
		case 4:
			deleteRowInteractive(db)
			fmt.Println("\nPress ENTER to return to menu...")
			reader.ReadString('\n')
		case 5:
			selectQueryManual(db)
			fmt.Println("\nPress ENTER to return to menu...")
			reader.ReadString('\n')
		case 6:
			insertQueryManual(db)
			fmt.Println("\nPress ENTER to return to menu...")
			reader.ReadString('\n')
		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}

}

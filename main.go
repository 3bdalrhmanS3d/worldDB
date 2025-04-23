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

var currentDatabase string

func clearScreen() {
	var cmd *exec.Cmd
	var err error
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	err = cmd.Run()

	if err != nil {
		fmt.Println(err)
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
	db, err := sql.Open(driver, fmt.Sprintf("%s:%s@/", user, password))
	if err != nil {
		log.Fatal("Connection failed:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to DB server:", err)
	}

	fmt.Println("\nConnected to DB server successfully.")

	return db
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

func createDatabase(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the name of the database to create: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	query := fmt.Sprintf("CREATE DATABASE %s", name)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Created database %s\n", name)
}
func dropDatabase(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the name of the database to drop: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	query := fmt.Sprintf("Drop DATABASE %s", name)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Dropped database %s\n", name)
}
func createTable(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the name of the table to create: ")
	tableName, _ := reader.ReadString('\n')
	tableName = strings.TrimSpace(tableName)

	columns := getColumns()
	if columns == "" {
		fmt.Println("Failed to create table due to invalid columns.")
		return
	}

	// CREATE TABLE
	query := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, columns)
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	fmt.Printf("Table '%s' created successfully.\n", tableName)
}
func dropTable(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the name of the table to drop: ")
	tableName, _ := reader.ReadString('\n')
	tableName = strings.TrimSpace(tableName)

	query := fmt.Sprintf("DROP TABLE %s", tableName)
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error dropping table:", err)
		return
	}

	fmt.Printf("Table '%s' dropped successfully.\n", tableName)
}
func useDatabase(db *sql.DB) {

	databases, err := showDatabases(db)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter the number of the database you want to use: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(databases) {
		fmt.Println("Invalid selection. Please choose a valid number.")
		return
	}

	selectedDB := databases[index-1]

	_, err = db.Exec("USE " + selectedDB)
	if err != nil {
		fmt.Printf("Error selecting database '%s': %v\n", selectedDB, err)
		return
	}

	currentDatabase = selectedDB
	fmt.Printf("Now using database: %s\n", selectedDB)
}
func getColumns() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the number of columns: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	numColumns, err := strconv.Atoi(input)
	if err != nil || numColumns <= 0 {
		fmt.Println("Invalid number of columns")
		return ""
	}

	var columns []string
	for i := 1; i <= numColumns; i++ {
		fmt.Printf("Enter column %d name: ", i)
		columnName, _ := reader.ReadString('\n')
		columnName = strings.TrimSpace(columnName)

		fmt.Printf("Enter type for column %d (e.g., INT, VARCHAR(50)): ", i)
		columnType, _ := reader.ReadString('\n')
		columnType = strings.TrimSpace(columnType)

		columnDefinition := fmt.Sprintf("%s %s", columnName, columnType)
		columns = append(columns, columnDefinition)
	}

	return strings.Join(columns, ", ")
}
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

	var count int
	queryCheck := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", table, idColumn)
	err := db.QueryRow(queryCheck, idValue).Scan(&count)
	if err != nil {
		fmt.Println("Error checking row existence:", err)
		return
	}

	if count == 0 {
		fmt.Println("No matching row found to delete.")
		return
	}

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
func showDatabases(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		return nil, fmt.Errorf("cannot list databases: %v", err)
	}
	defer rows.Close()

	var databases []string
	var dbName string
	i := 1
	fmt.Println("\nAvailable Databases:")
	for rows.Next() {
		err = rows.Scan(&dbName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d. %s\n", i, dbName)
		databases = append(databases, dbName)
		i++
	}

	return databases, nil
}

func main() {
	db := dbConnectInteractive()
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		fmt.Println("\n===== MENU =====")

		if currentDatabase == "" {
			fmt.Println("1. Create Database")
			fmt.Println("2. Delete Database")
			fmt.Println("3. Show Databases")
			fmt.Println("4. Use Database")
		} else {
			fmt.Printf("Current Database: %s\n", currentDatabase)
			fmt.Println("1. Show Tables")
			fmt.Println("2. Select Table")
			fmt.Println("3. Create Table")
			fmt.Println("4. Drop Table")
			fmt.Println("5. Update Row")
			fmt.Println("6. Delete Row")
			fmt.Println("7. Run SELECT Query")
			fmt.Println("8. Run INSERT Query")
			fmt.Println("9. Change Database")
		}
		fmt.Println("0. Exit")

		fmt.Print("\nEnter your choice: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid choice, please try again.")
			continue
		}

		switch choice {
		case 0:
			fmt.Println("Exiting...")
			return

		case 1:
			if currentDatabase == "" {
				createDatabase(db)
			} else {
				listTables(db)
			}
		case 2:
			if currentDatabase == "" {
				dropDatabase(db)
			} else {
				fmt.Print("Enter table name: ")
				tableName, _ := reader.ReadString('\n')
				tableName = strings.TrimSpace(tableName)
				readTableByName(db, tableName)
			}
		case 3:
			if currentDatabase == "" {
				showDatabases(db)
			} else {
				createTable(db)
			}
		case 4:
			if currentDatabase == "" {
				useDatabase(db)
			} else {
				dropTable(db)
			}
		case 5:
			if currentDatabase == "" {
				fmt.Println("Please select a database first.")
			} else {
				updateRowInteractive(db)
			}
		case 6:
			if currentDatabase == "" {
				fmt.Println("Please select a database first.")
			} else {
				deleteRowInteractive(db)
			}
		case 7:
			if currentDatabase == "" {
				fmt.Println("Please select a database first.")
			} else {
				selectQueryManual(db)
			}
		case 8:
			if currentDatabase == "" {
				fmt.Println("Please select a database first.")
			} else {
				insertQueryManual(db)
			}
		case 9:
			if currentDatabase == "" {
				fmt.Println("Invalid choice, please try again.")
			} else {
				useDatabase(db)
			}
		default:
			fmt.Println("Invalid choice, please try again.")
		}

		fmt.Println("\nPress Enter to continue...")
		reader.ReadString('\n')

	}
}

// Nils Elde
// https://gitlab.com/nilsanderselde

package dbconnect

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	// DBInfo stores the info used to connect to the database
	DBInfo string

	// DB is database connection
	DB *sql.DB
)

// GetDBInfo gets the info to connect to database from external file
func GetDBInfo() string {
	file, err := os.Open("env/dbinfo")

	if err != nil {
		fmt.Println("DB connection information not found.")
		return ""
	}
	defer file.Close()

	var login []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		login = strings.Split(scanner.Text(), ",")
	}

	if len(login) == 5 {
		host := login[0]
		port, _ := strconv.Atoi(login[1])
		user := login[2]
		password := login[3]
		dbname := login[4]
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	}
	fmt.Println("Invalid database connection information provided.")
	return ""

}

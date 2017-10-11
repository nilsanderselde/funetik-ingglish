// Nils Elde
// https://gitlab.com/nilsanderselde

package dbconnect

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// DBInfo stores the info used to connect to the database
var DBInfo string

// GetDBInfo gets the info to connect to database from external file
func GetDBInfo() string {
	file, err := os.Open("db/dbinfo")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var login []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		login = strings.Split(scanner.Text(), ",")
	}

	host := login[0]
	port, _ := strconv.Atoi(login[1])
	user := login[2]
	password := login[3]
	dbname := login[4]

	fmt.Println("Ready to connect.")
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

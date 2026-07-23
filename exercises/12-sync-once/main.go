package main

import (
	"fmt"
	"sync"
	"time"
)

type DBConnection struct {
	DSN string
}

var (
	once sync.Once
	db   *DBConnection
)

func connect(dsn string) *DBConnection {
	fmt.Println("→ connecting to database...")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("→ connection established")
	return &DBConnection{DSN: dsn}
}

func getDBConnection() *DBConnection {
	once.Do(func() {
		db = connect("postgres://localhost/mydb")
	})
	return db
}

func main() {
	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			conn := getDBConnection()
			fmt.Printf("goroutine %2d using connection: %s\n", id, conn.DSN)
		}(i)
	}

	wg.Wait()

	fmt.Println("\ncalling getDBConnection() again from main:")
	conn := getDBConnection()
	fmt.Printf("got connection: %s\n", conn.DSN)
}

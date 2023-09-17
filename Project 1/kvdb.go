/**
 * @file kvdb.go
 * @author Joshua Calzadillas (jmc1241@usnh.edu)
 * @brief Project 1 - KVStore
 * @date 2023-09-10
 */
package main

import "fmt"

// Database test
func DBTest(db *map[string]string) {
	// Adding item to db
	DBAddItem(db, "h1", "Hello World!")
	DBAddItem(db, "h2", "Hello Josh!")
	DBAddItem(db, "h3", "Poop!")

	// Adding item to db
	DBRemoveItem(db, "h3")
	DBRemoveItem(db, "h4")

	// Database Printout
	fmt.Println(*db)

	// Get items
	fmt.Println(DBGet(db, "h2"))
	fmt.Println(DBGet(db, "h4"))
}

// DB add item
func DBAddItem(db *map[string]string, Key string, Value string) {
	// Adding item to db
	(*db)[Key] = Value
}

// DB Get item based on key
func DBGet(db *map[string]string, Key string) (string, string) {
	return Key, (*db)[Key]
}

// DB remove item
func DBRemoveItem(db *map[string]string, Key string) {
	if (*db)[Key] != "" {
		delete((*db), Key)
	} else {
		return
	}
}

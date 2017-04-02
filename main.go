package main

import (
	"flag"
	"fmt"
)

/* main program */
func main() {
	// parse command line flags
	flag.Parse()
	// check if positional arguments are provided
	if flag.NArg() != 1 {
		fmt.Println("Usage: backup <file.json>")
		return
	}
	// read and parse backup configuration file
	rsync, err := ParseConfigFile(flag.Arg(0))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// make sure backup directories exist
	if err := rsync.Init(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	// perform backup
	if allow, err := rsync.AllowBackup(); err != nil {
		fmt.Println("AllowBackup():", err)
		return
	} else if allow {
		if err := DoLock(rsync); err != nil {
			fmt.Println("DoLock():", err)
			return
		}
	}
}

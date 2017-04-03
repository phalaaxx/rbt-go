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
		fmt.Println("Usage: backup file.json [file1.json ...]")
		return
	}
	for idx := 0; idx < flag.NArg(); idx++ {
		// read and parse backup configuration file
		rsync, err := ParseConfigFile(flag.Arg(idx))
		if err != nil {
			fmt.Printf("Error: %s: %v\n", flag.Arg(idx), err)
			continue
		}
		// make sure backup directories exist
		if err := rsync.Init(); err != nil {
			fmt.Printf("Init(): %s: %v\n", rsync.Name, err)
			continue
		}
		// perform backup
		if allow, err := rsync.AllowBackup(); err != nil {
			fmt.Printf("AllowBackup(): %s: %v\n", rsync.Name, err)
			continue
		} else if allow {
			if err := DoLock(rsync); err != nil {
				fmt.Printf("DoLock(): %s: %v\n", rsync.Name, err)
				continue
			}
		}
	}
}

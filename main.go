package main

import (
	"flag"
	"fmt"
	"strings"
)

// ConfigFiles defines a list of strings
type ConfigFiles []string

// String returns a representation of the object
func (f *ConfigFiles) String() string {
	return strings.Join(*f, ",")
}

// Set adds a new value to the list
func (f *ConfigFiles) Set(value string) error {
	*f = append(*f, value)
	return nil
}

/* main program */
func main() {
	// parse command line flags
	var files ConfigFiles
	flag.Var(&files, "f", "Backup configuration file.")
	flag.Parse()
	// make sure there are file names provided
	if len(files) == 0 {
		fmt.Println("For usage use -h")
		return
	}
	// run backup process for all provided arguments
	for _, file := range files {
		// read and parse backup configuration file
		rsync, err := ParseConfigFile(file)
		if err != nil {
			fmt.Printf("Error: %s: %v\n", file, err)
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

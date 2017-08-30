package main

import (
	"flag"
	"fmt"
	"os"
	"path"
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
	// make sure file name ends in .json
	if !strings.HasSuffix(value, ".json") {
		value = fmt.Sprintf("%s.json", value)
	}
	// check if file exists
	if _, err := os.Stat(value); err == nil {
		*f = append(*f, value)
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	// does not exist, try in /etc/rbt
	value = path.Clean(path.Join("/etc/rbt", value))
	if _, err := os.Stat(value); err == nil {
		*f = append(*f, value)
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	return os.ErrNotExist
}

/* main program */
func main() {
	// parse command line flags
	var files ConfigFiles
	var verbose bool
	flag.Var(&files, "f", "Backup configuration file.")
	flag.BoolVar(&verbose, "v", false, "Verbose output.")
	flag.Parse()
	// make sure there are file names provided
	if len(files) == 0 {
		// if program name is not rbt this may indicate a symlink
		if strings.HasSuffix(os.Args[0], "rbt") {
			fmt.Println("For usage use -h")
			return
		}
		files.Set(path.Base(os.Args[0]))
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
			if err := DoLock(rsync, verbose); err != nil {
				fmt.Printf("DoLock(): %s: %v\n", rsync.Name, err)
				continue
			}
		}
	}
}

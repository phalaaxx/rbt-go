package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// RsyncOptions contains configuration options for rsync
type RsyncOptions struct {
	Name    string   `json:"name"`
	Backups int      `json:"backups"`
	Target  string   `json:"target"`
	Files   []string `json:"files"`
	Exclude []string `json:"exclude"`
}

// GetTarget parses target string value
func (r RsyncOptions) GetTarget() string {
	return strings.Replace(r.Target, "$name", r.Name, -1)
}

// GetTargetBackup renders the destination directory for rsync
func (r RsyncOptions) GetTargetBackup() string {
	return fmt.Sprintf("%s/backup.%d", r.GetTarget(), r.Backups)
}

// GetLastBackup returns directory name of last completed backup
func (r RsyncOptions) GetLastBackup() string {
	return fmt.Sprintf("%s/backup.0", r.GetTarget())
}

// GetLastBackupTime returns the time when the last backup was completed
func (r RsyncOptions) GetLastBackupTime() (*time.Time, error) {
	// create completed file
	completed := fmt.Sprintf("%s/backup.0/completed", r.GetTarget())
	File, err := os.Open(completed)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return new(time.Time), nil
	}
	defer File.Close()
	// get time from completed file
	var stats CompletedStats
	decoder := json.NewDecoder(File)
	if err := decoder.Decode(&stats); err != nil {
		return nil, err
	}
	return &stats.Timestamp, nil
}

// AllowBackup returns true if enough time has passed since the last backup
func (r RsyncOptions) AllowBackup() (bool, error) {
	// get finish time of latest backup
	last, err := r.GetLastBackupTime()
	if err != nil {
		return false, err
	}
	// allow next backup to start at least one hour after the last one
	return last.Add(time.Hour).Before(time.Now()), nil
}

// Options generates options for running rsync
func (r *RsyncOptions) Options() []string {
	options := []string{
		"-aR",
		"--delete",
		"--stats",
		fmt.Sprintf("--link-dest=%s", r.GetLastBackup()),
	}
	// add source files / directories
	for _, dir := range r.Files {
		options = append(options, fmt.Sprintf("%s:%s", r.Name, dir))
	}
	// add excludes list
	for _, dir := range r.Exclude {
		options = append(options, "--exclude", dir)
	}
	// add target directory
	options = append(options, r.GetTargetBackup())
	return options
}

// Rotate backup directories by one
func (r *RsyncOptions) Rotate() error {
	target := r.GetTarget()
	// rename last backup directory to a temporary name
	TmpDir := fmt.Sprintf("%s/backup.tmp", target)
	if err := os.Rename(r.GetTargetBackup(), TmpDir); err != nil {
		return err
	}
	// rotate remaining directories
	for idx := r.Backups; idx > 0; idx-- {
		SrcDir := fmt.Sprintf("%s/backup.%d", target, idx-1)
		DstDir := fmt.Sprintf("%s/backup.%d", target, idx)
		if err := os.Rename(SrcDir, DstDir); err != nil {
			return err
		}
	}
	// move temp directory to backup.0
	if err := os.Rename(TmpDir, r.GetLastBackup()); err != nil {
		return err
	}
	return nil
}

// SoveCompleted saves backup stats in the completed file
func (r *RsyncOptions) SaveCompleted(duration int64) error {
	// create completed file
	completed := fmt.Sprintf("%s/backup.%d/completed", r.GetTarget(), r.Backups)
	File, err := os.Create(completed)
	if err != nil {
		return err
	}
	defer File.Close()
	// prepare completed data
	stats := CompletedStats{
		time.Now(),
		duration,
	}
	// encode json data
	encoder := json.NewEncoder(File)
	if err := encoder.Encode(&stats); err != nil {
		return err
	}
	// set file permissions
	if err := File.Chmod(0644); err != nil {
		return err
	}
	return nil
}

// Run rsync data transfer and rotate directories on success
func (r *RsyncOptions) Run() error {
	// prepare rsync command
	cmd := exec.Command(
		"/usr/bin/rsync",
		r.Options()...)
	// get start time
	StartTime := time.Now().Unix()
	// start rsync to transfer data
	if err := cmd.Start(); err != nil {
		return err
	}
	// wait for command to exit
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				// exit status 24 means files changed during transfer
				// this is a non-fatal error
				if status.ExitStatus() != 24 {
					return err
				}
			}
		} else {
			return err
		}
	}
	// create completed file
	duration := time.Now().Unix() - StartTime
	if err := r.SaveCompleted(duration); err != nil {
		return err
	}
	// rotate backups and return
	if err := r.Rotate(); err != nil {
		return err
	}
	return nil
}

// Init makes sure backup directories exist prior to backup procedure
func (r *RsyncOptions) Init() error {
	target := r.GetTarget()
	// make sure root backup directory exists
	if _, err := os.Stat(target); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// create root backup directory
		if err := os.MkdirAll(target, 0755); err != nil {
			return err
		}
	}
	// make sure all backup directories exist
	for idx := 0; idx < r.Backups; idx++ {
		BackupDir := fmt.Sprintf("%s/backup.%d", target, idx)
		if _, err := os.Stat(BackupDir); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			// make backup directory
			if err := os.Mkdir(BackupDir, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}

// ParseConfigFile reads and parses a json configuration file
func ParseConfigFile(Name string) (*RsyncOptions, error) {
	// open configuration file for reading
	File, err := os.Open(Name)
	if err != nil {
		return nil, err
	}
	defer File.Close()
	decoder := json.NewDecoder(File)

	// decode json data
	rsync := new(RsyncOptions)
	if err := decoder.Decode(rsync); err != nil {
		fmt.Println("DECODER:", err)
		return nil, err
	}

	return rsync, nil
}

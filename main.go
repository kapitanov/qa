package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"
)

const (
	versionNumber = "vUnknown"
)

func (t *commandConfig) execute() error {
	command, err := exec.LookPath(t.Command)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to locate executable \"%s\" in PATH", t.Command)
		return fmt.Errorf(errorMessage)
	}

	os.Stderr.Sync()
	cmd := &exec.Cmd{
		Path: command,
		Args: append([]string{command}, t.Arguments...),
	}

	log.Printf("+ '%s' %s", cmd.Path, cmd.Args)
	err = cmd.Run()

	if err == nil {
		log.Printf("%s exited with %d", cmd.Path, 0)
		return nil
	}

	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			exitCode := status.ExitStatus()
			log.Printf("%s exited with %d", cmd.Path, exitCode)
			os.Exit(exitCode)
			return nil
		}
	}

	log.Printf("< %s", err)
	return err
}

func processFlags() bool {
	verbose := flag.Bool("v", false, "Enable verbose logging")
	verbose2 := flag.Bool("verbose", false, "Enable verbose logging")
	version := flag.Bool("version", false, "Print version and exit")
	customPath := flag.String("c", "", "Custom config file path")

	flag.Parse()

	log.SetFlags(0)
	if !*verbose && !*verbose2 {
		log.SetOutput(ioutil.Discard)
	}

	customConfigPath = *customPath

	if *version {
		fmt.Fprintf(os.Stderr, "%s\n", versionNumber)
		return false
	}

	return true
}

func main() {
	// Process command line flags
	if !processFlags() {
		return
	}

	// Load config file
	items, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
		return
	}

	if items == nil {
		return
	}

	// Run UI
	selectedItem, err := runUI(items)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
		return
	}

	// Run selected command
	if selectedItem != nil {
		err := selectedItem.execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
			return
		}
	}
}

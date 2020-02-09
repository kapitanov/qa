package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func (t *commandConfig) execute() error {
	cmd := &exec.Cmd{
		Path: t.Command,
		Args: append([]string{t.Command}, t.Arguments...),
	}
	log.Printf("+ '%s' %s", cmd.Path, cmd.Args)
	err := cmd.Run()
	return err
}

func processFlags() {
	verbose := flag.Bool("v", false, "Enable verbose logging")
	customPath := flag.String("c", "", "Custom config file path")

	flag.Parse()

	log.SetFlags(0)
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	customConfigPath = *customPath
}

func main() {
	// Process command line flags
	processFlags()

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

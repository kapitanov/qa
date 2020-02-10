package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
)

var (
	homeDir            string
	customConfigPath   string
	defaultConfigPath  string
	fallbackConfigPath string
)

type configRoot struct {
	Commands []*commandConfig `json:"commands"`
}

type commandConfig struct {
	Name      string   `json:"name"`
	Command   string   `json:"cmd"`
	Arguments []string `json:"args"`
}

func (cmd *commandConfig) prepare() {
	cmd.Command = substEnv(cmd.Command)
	if cmd.Arguments == nil {
		cmd.Arguments = make([]string, 0)
	}

	for i := range cmd.Arguments {
		cmd.Arguments[i] = substEnv(cmd.Arguments[i])
	}
}

func (cmd *commandConfig) validate() error {
	if cmd.Name == "" {
		return fmt.Errorf("No name specified")
	}

	if cmd.Command == "" {
		return fmt.Errorf("No command specified")
	}

	return nil
}

func init() {
	env := "HOME"
	switch runtime.GOOS {
	case "windows":
		env = "USERPROFILE"
	}
	homeDir := os.Getenv(env)

	defaultConfigPath = path.Join(homeDir, ".qa")

	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}

	fallbackConfigPath = path.Join(path.Dir(executable), "config.json")
}

func loadConfig() ([]*commandConfig, error) {
	if customConfigPath != "" {
		// Try to load custom config
		commands, err := tryLoadConfig(customConfigPath)
		if err != nil {
			return nil, err
		}

		if commands == nil {
			fmt.Fprintf(os.Stderr, "Config file \"%s\" doesn't exist\n", customConfigPath)
			os.Exit(1)
			return nil, nil
		}
	}

	// Try to load user config
	commands, err := tryLoadConfig(defaultConfigPath)
	if err != nil {
		return nil, err
	}

	if commands == nil {
		fmt.Fprintf(os.Stderr, "Config file \"%s\" doesn't exist\n", defaultConfigPath)
		fmt.Fprintf(os.Stderr, "Create config file and add commands you need\n")
		fmt.Fprintf(os.Stderr, "Here is an example of config file content:\n%s\n", configTemplate)

		return nil, nil
	}

	if len(commands) == 0 {
		fmt.Fprintf(os.Stderr, "Config file \"%s\" is empty - there're no command defined.", defaultConfigPath)
		fmt.Fprintf(os.Stderr, "Edit config file and add commands you need\n")
		fmt.Fprintf(os.Stderr, "Here is an example of config file content:\n%s\n", configTemplate)

		return nil, nil
	}

	return commands, nil
}

func substEnv(str string) string {
	str = strings.Replace(str, "~", homeDir, -1)
	str = strings.Replace(str, "$HOME", homeDir, -1)

	cwd, err := os.Getwd()
	if err != nil {
		str = strings.Replace(str, "$(pwd)", cwd, -1)
	}

	str = os.ExpandEnv(str)
	return str
}

func tryLoadConfig(path string) ([]*commandConfig, error) {
	if path == "" {
		return nil, nil
	}

	log.Printf("Trying to load config from \"%s\"", path)
	file, err := ioutil.ReadFile(path)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		log.Printf("Trying to load config from \"%s\": %s", path, err)
		return nil, err
	}

	data := configRoot{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	if data.Commands == nil {
		return nil, fmt.Errorf("Maformed config file \"%s\"", path)
	}

	for i, cmd := range data.Commands {
		cmd.prepare()
		log.Printf("Found command: \"%s\"", cmd.Name)
		err = cmd.validate()
		if err != nil {
			log.Printf("Command #%d is not valid: %s", i, err)
			return nil, err
		}
	}

	return data.Commands, nil
}

const configTemplate = `{
    "commands" :[
        {
            "name":"Command display name",
            "cmd" :"command_to_execute"
        },
        {
            "name":"Command display name",
            "cmd" :"command_to_execute",
            "args" : ["list","of","arguments"]
        }
    ]
}
`

func createEmptyConfig(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(configTemplate)
	if err != nil {
		return err
	}
	return nil
}

# QuickAccess

![demo](demo.gif)

QuickAccess (QA) is a tiny app written in Go.
It serves one purpose: to create a quick access menu in your terminal.
Menu items are basically shortcuts to shell command.

You can use this app to:

* quickly connect to remote servers by SSH
* simplify various daily tasks

## Build

To build this app for your current platform, type:

```shell
go get && go build
```

To build it for every supported platform (linux/x64 and windows/x64), type:

```shell
make
```

## Installing

### Installing from source

```shell
go get -u github.com/kapitanov/qa
```

## Configuring

You need to create a JSON file `.qa` in your home directory with the following content:

```json
{
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
```

Root element `commands` is required and it has to be an array.
Its items map are menu items. Their properties are:

* `name` - a menu item title (required element)
* `cmd` - a command to run (required element)
* `args` - a list of command's arguments (optional)

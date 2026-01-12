package main

import(
	"errors"	
	f "fmt"
	"os"
	) 

func CommandInterpreter(input string){

	commands := returnCurrentCommands()

	if _, ok := commands[input]; ok {
		commands[input].callback()
	} else {
		f.Println("Uknown command")
		return
	}
}

func displayHelp() error {
	f.Println("Welcome to the Pokedex!")
	f.Println("Usage: \n")
	commands := returnCurrentCommands()
	for _, value := range commands {
		f.Printf("%v: %v\n", value.name, value.description)
	}
	return errors.New("Blalba")
}

func commandExit() error {
    f.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return errors.New("Blabla")
}

type cliCommand struct {
		name        string
		description string
		callback    func() error
	}

func returnCurrentCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:		 "help",
			description: "Displays a help message",
			callback:     displayHelp,
		},
	}
	return commands
}
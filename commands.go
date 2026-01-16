package main

import(
	"errors"	
	f "fmt"
	"os"
	api "github.com/FinFreezer/go-pokedex/apihandler"
	) 

func CommandInterpreter(config *api.Config, input []string) (*api.Config) {

	commands := returnCurrentCommands()

	if _, ok := commands[input[0]]; ok {
		config, err := commands[input[0]].callback(config, input[1:])
		if err != nil {
			f.Println("Error reaching command, %w", err)
			return config
		}
		return config
	} else {
		f.Println("Unknown command")
		return config
	}
}

func displayHelp(config *api.Config, addParams []string) (*api.Config, error) {
	f.Println("Welcome to the Pokedex!")
	f.Println("Usage: ")
	commands := returnCurrentCommands()
	for _, value := range commands {
		f.Printf("%v: %v\n", value.name, value.description)
	}
	return config, nil
}

func commandExit(config *api.Config, addParams []string) (*api.Config, error) {
    f.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return config, errors.New("Blabla")
}

func getNextLoc(config *api.Config, addParams []string) (*api.Config, error) {
	if config == nil {
		f.Println("Calling API with Init")
		config = api.Start(config, "Init", addParams)
		return config, nil
	
	} else {
		f.Println("Calling API with Next")
		config = api.Start(config, "Next", addParams)
		return config, nil
	}
}

func getPreviousLoc(config *api.Config, addParams []string) (*api.Config, error) {
	f.Println("Calling API with Previous")
	api.Start(config, "Previous", addParams)
	return config, nil
}

func exploreLoc(config *api.Config, addParams []string) (*api.Config, error) {
	f.Println("Starting API for exploration.")
	api.Start(config, "Explore", addParams)
	return config, nil
}

type cliCommand struct {
		name        string
		description string
		callback    func(config *api.Config, addParams []string) (*api.Config, error)
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
		"map": {
			name: 		 "map",
			description: "Display 20 locations from the world map",
			callback: 	 getNextLoc,
		},
		"mapb": {
			name: 		 "mapb",
			description: "Display the previous 20 locations from the world map",
			callback: 	  getPreviousLoc,
		},
		"explore": {
			name:		 "explore",
			description: "List all Pokemon found at current location.",
			callback: 	 exploreLoc,
		},
	}
	return commands
}
package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
    api "github.com/FinFreezer/go-pokedex/apihandler"
    //cache "github.com/FinFreezer/go-pokedex/internal/pokecache"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    var config *api.Config
    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        input := cleanInput(scanner.Text())
        config = CommandInterpreter(input[0], config)
        
    }
}

func cleanInput (text string) []string {
    return strings.Fields(strings.ToLower(text))
}
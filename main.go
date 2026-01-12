package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        input := cleanInput(scanner.Text())
        CommandInterpreter(input[0])
        
    }
}

func cleanInput (text string) []string {
    //separatedStrings := []string

    return strings.Fields(strings.ToLower(text))
}
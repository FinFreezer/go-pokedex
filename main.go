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
        inputFirst := input[0]
        fmt.Printf("Your command was: %s\n", inputFirst)
    }
}

func cleanInput (text string) []string {
    //separatedStrings := []string

    return strings.Fields(strings.ToLower(text))
}
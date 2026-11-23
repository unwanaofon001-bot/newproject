package main

import (
    "fmt"
    "os"
    "strings"
)

// holds one color+substring pair
type ColorRule struct {
    color     string
    substring string
}

func main() {
    rules, text := parseInput()
    result := printAscii(rules, text)
    fmt.Print(result)
}

func parseInput() ([]ColorRule, string) {
    input := os.Args[1:]
    if len(input) < 1 {
        fmt.Println("Usage: go run . --color=COLOR [SUBSTRING] --color=COLOR [SUBSTRING] TEXT")
        os.Exit(1)
    }

    var rules []ColorRule
    var text string

    // walk through all arguments
    i := 0
    for i < len(input) {
        if strings.HasPrefix(input[i], "--color=") {
            // extract the color name after "="
            color := strings.Split(input[i], "=")[1]

            // the next argument is the substring for this color
            substring := ""
            if i+1 < len(input) && !strings.HasPrefix(input[i+1], "--color=") {
                // check if next arg is the final text or a substring
                // the last non-flag argument is always the full text
                // so we peek ahead — if there's still more args after, this is a substring
                if i+2 < len(input) {
                    substring = input[i+1]
                    i++ // skip the substring argument
                }
            }

            rules = append(rules, ColorRule{color: color, substring: substring})
        } else {
            // last argument with no --color= prefix is the full text
            text = input[i]
        }
        i++
    }

    return rules, text
}

func printAscii(rules []ColorRule, text string) string {
    data, err := os.ReadFile("standard.txt")
    if err != nil {
        fmt.Println("error found")
        return ""
    }

    content := strings.ReplaceAll(string(data), "\r\n", "\n")
    lines := strings.Split(content, "\n")

    colorCodes := map[string]string{
        "black":  "\033[30m",
        "red":    "\033[31m",
        "green":  "\033[32m",
        "yellow": "\033[33m",
        "blue":   "\033[34m",
        "purple": "\033[35m",
        "cyan":   "\033[36m",
        "white":  "\033[37m",
        "orange": "\033[38;5;208m",
        "pink":   "\033[38;2;255;105;180m",
    }

    reset := "\033[0m"

    // colorMap stores which color code to use for each character position
    // empty string means no color for that position
    colorMap := make([]string, len(text))

    // go through each rule and mark positions in colorMap
    for _, rule := range rules {
        code := colorCodes[rule.color]
        if rule.substring == "" {
            // no substring means color the whole text
            for i := range text {
                colorMap[i] = code
            }
        } else {
            // find where the substring appears and mark those positions
            for i := 0; i <= len(text)-len(rule.substring); i++ {
                if text[i:i+len(rule.substring)] == rule.substring {
                    for k := 0; k < len(rule.substring); k++ {
                        colorMap[i+k] = code
                    }
                }
            }
        }
    }

    var output strings.Builder
    for row := 0; row < 9; row++ {
        for i, char := range text {
            ascii := int(char)
            if ascii < 32 || ascii > 126 {
                continue
            }
            value := (ascii - 32) * 9
            if colorMap[i] != "" {
                output.WriteString(colorMap[i] + lines[value+row] + reset)
            } else {
                output.WriteString(lines[value+row])
            }
        }
        output.WriteString("\n")
    }

    return output.String()
}
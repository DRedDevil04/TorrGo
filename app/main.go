package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func decodeBencode(bencodedString string) (interface{}, error) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		var firstColonIndex int

		for i := 0; i < len(bencodedString); i++ {
			if bencodedString[i] == ':' {
				firstColonIndex = i
				break
			}
		}

		lengthStr := bencodedString[:firstColonIndex]

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", err
		}

		return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
	} else if unicode.IsLetter(rune(bencodedString[0])) {
		if bencodedString[0] == 'i' {
			// Handle integer
			var end int
			for i := 1; i < len(bencodedString); i++ {
				if bencodedString[i] == 'e' {
					end = i
					break
				}
			}
			intValue, err := strconv.Atoi(bencodedString[1:end])
			if err != nil {
				return "", err
			}
			return intValue, nil
		} else if bencodedString[0] == 'l' {
			// Handle list
			var list []interface{}
			i := 1
			for i < len(bencodedString)-1 {
				value, err := decodeBencode(bencodedString[i:])
				if err != nil {
					return nil, err
				}
				list = append(list, value)

				// Update `i` to move past the decoded value
				// Assuming `value` is a string, integer, or list, calculate its length
				switch v := value.(type) {
				case string:
					i += len(v) + len(strconv.Itoa(len(v))) + 1 // length prefix + colon
				case int:
					i += len(strconv.Itoa(v)) + 2 // 'i' + number + 'e'
				default:
					return "", fmt.Errorf("Invalid type in list")
				}
				fmt.Println("List detected", value, "i ", i)
			}
			return list, nil
		} else {
			return "", fmt.Errorf("Invalid bencoded string")
		}
	} else {
		return "", fmt.Errorf("Only strings are supported at the moment")
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		// Uncomment this block to pass the first stage

		bencodedValue := os.Args[2]

		decoded, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}

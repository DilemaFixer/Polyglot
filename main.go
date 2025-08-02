package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Dictionaries struct {
	Dicts map[string]*Dictionary `json:"dicts"`
}

type Dictionary struct {
	From  string            `json:"from"`
	To    string            `json:"to"`
	Words map[string]string `json:"words"`
}

func NewDictionary(from string, to string) *Dictionary {
	return &Dictionary{
		From:  from,
		To:    to,
		Words: make(map[string]string),
	}
}

func NewDictionaries() *Dictionaries {
	return &Dictionaries{
		Dicts: make(map[string]*Dictionary),
	}
}

func (dicts *Dictionaries) AddDictionary(key string, dict *Dictionary) {
	dicts.Dicts[key] = dict
}

func (dicts *Dictionaries) GetDictionary(key string) *Dictionary {
	return dicts.Dicts[key]
}

func (d *Dictionary) AddWord(word, translation string) {
	d.Words[word] = translation
}

func (d *Dictionary) GetTranslation(word string) (string, bool) {
	translation, exists := d.Words[word]
	return translation, exists
}

func (dicts *Dictionaries) SaveToJson(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(dicts); err != nil {
		return fmt.Errorf("failed to encode JSON: %v", err)
	}

	return nil
}

func LoadDictsFromJson(filename string) (*Dictionaries, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var dicts Dictionaries
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&dicts); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v", err)
	}

	return &dicts, nil
}

func initializeDictionaries() *Dictionaries {
	dictionaries := NewDictionaries()

	enRuDict := NewDictionary("English", "Russian")
	enRuDict.AddWord("hello", "привет")
	enRuDict.AddWord("world", "мир")
	enRuDict.AddWord("computer", "компьютер")

	deRuDict := NewDictionary("German", "Russian")
	deRuDict.AddWord("hallo", "привет")
	deRuDict.AddWord("welt", "мир")
	deRuDict.AddWord("computer", "компьютер")

	frRuDict := NewDictionary("French", "Russian")
	frRuDict.AddWord("bonjour", "привет")
	frRuDict.AddWord("monde", "мир")
	frRuDict.AddWord("ordinateur", "компьютер")

	dictionaries.AddDictionary("en-ru", enRuDict)
	dictionaries.AddDictionary("de-ru", deRuDict)
	dictionaries.AddDictionary("fr-ru", frRuDict)

	return dictionaries
}

func findMode(dict *Dictionary, scanner *bufio.Scanner) {
	fmt.Printf("Find mode for %s->%s dictionary\n", dict.From, dict.To)
	fmt.Println("Type 'exit' to return to main menu")

	for {
		fmt.Print("Enter word to translate: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "exit" {
			break
		}

		if input == "" {
			continue
		}

		if translation, exists := dict.GetTranslation(input); exists {
			fmt.Printf("Translation: %s = %s\n", input, translation)
		} else {
			fmt.Printf("Word '%s' not found in dictionary\n", input)
		}
	}
}

func writeMode(dict *Dictionary, scanner *bufio.Scanner) {
	fmt.Printf("Write mode for %s->%s dictionary\n", dict.From, dict.To)
	fmt.Println("Type 'exit' to return to main menu")

	for {
		fmt.Print("Enter word: ")
		if !scanner.Scan() {
			break
		}

		word := strings.TrimSpace(scanner.Text())
		if word == "exit" {
			break
		}

		if word == "" {
			continue
		}

		fmt.Print("Enter translation: ")
		if !scanner.Scan() {
			break
		}

		translation := strings.TrimSpace(scanner.Text())
		if translation == "" {
			fmt.Println("Translation cannot be empty")
			continue
		}

		dict.AddWord(word, translation)
		fmt.Printf("Added: %s = %s\n", word, translation)
	}
}

func printUsage() {
	fmt.Println("Usage: program <dictionary> <mode>")
	fmt.Println("Dictionaries: en-ru, de-ru, fr-ru")
	fmt.Println("Modes: -find, -write")
	fmt.Println("Example: program en-ru -find")
}

func main() {
	if len(os.Args) != 3 {
		printUsage()
		return
	}

	dictKey := os.Args[1]
	mode := os.Args[2]

	filename := "dictionaries.json"
	var dictionaries *Dictionaries
	var err error

	dictionaries, err = LoadDictsFromJson(filename)
	if err != nil {
		fmt.Printf("Failed to load dictionaries, creating new ones: %v\n", err)
		dictionaries = initializeDictionaries()
	}

	dict := dictionaries.GetDictionary(dictKey)
	if dict == nil {
		fmt.Printf("Dictionary '%s' not found\n", dictKey)
		fmt.Println("Available dictionaries:")
		for key, d := range dictionaries.Dicts {
			fmt.Printf("  %s (%s->%s)\n", key, d.From, d.To)
		}
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		switch mode {
		case "-find":
			findMode(dict, scanner)
		case "-write":
			writeMode(dict, scanner)
		default:
			fmt.Printf("Unknown mode '%s'\n", mode)
			printUsage()
			return
		}

		fmt.Print("Save changes and continue? (y/n/q): ")
		if !scanner.Scan() {
			break
		}

		choice := strings.ToLower(strings.TrimSpace(scanner.Text()))

		if choice == "q" {
			break
		}

		if choice == "y" {
			if err := dictionaries.SaveToJson(filename); err != nil {
				fmt.Printf("Failed to save dictionaries: %v\n", err)
			} else {
				fmt.Println("Dictionaries saved successfully")
			}
		}

		fmt.Print("Switch mode? (find/write/exit): ")
		if !scanner.Scan() {
			break
		}

		newMode := strings.TrimSpace(scanner.Text())
		switch newMode {
		case "find":
			mode = "-find"
		case "write":
			mode = "-write"
		case "exit":
			return
		default:
			fmt.Printf("Unknown mode '%s', keeping current mode\n", newMode)
		}
	}

	if err := dictionaries.SaveToJson(filename); err != nil {
		fmt.Printf("Failed to save dictionaries: %v\n", err)
	} else {
		fmt.Println("Dictionaries saved successfully")
	}
}

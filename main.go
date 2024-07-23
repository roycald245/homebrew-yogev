package main

import (
	"bufio"
	"embed"
	"fmt"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"strings"
)

//go:embed resources
var f embed.FS

func main() {
	var rootCmd = &cobra.Command{
		Use:   "yogev",
		Short: "Yogev is a basic command line tool for education and facts",
	}

	// Define a command to fetch a fact
	var factCmd = &cobra.Command{
		Use:   "fact",
		Short: "Get a random fact",
		Run: func(cmd *cobra.Command, args []string) {
			generateFact()
		},
	}

	rootCmd.AddCommand(factCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateFact() {
	template := randomLine("templates")
	words := strings.Fields(template)

	// Replace placeholders with random words

	for i, word := range words {
		trimmedWord := strings.Trim(word, ".,;:!?)}")
		switch trimmedWord {
		case "%n":
			words[i] = randomLine("nouns")
		case "%ns":
			words[i] = pluralizeNoun(randomLine("nouns"))
		case "%vs":
			words[i] = randomLine("verbs")
		case "%v":
			words[i] = thridPersoniseVerb(randomLine("verbs"))
		case "%a":
			words[i] = randomLine("adjectives")
		case "%num":
			words[i] = fmt.Sprintf("%d", rand.Intn(1000000))
		case "%per":
			words[i] = fmt.Sprintf("%d%%", rand.Intn(100))
		}
	}

	fmt.Println(strings.Join(words, " "))
}

func randomLine(wordType string) string {
	file, err := f.ReadFile(fmt.Sprintf("resources/%s/%s.txt", wordType, "english"))
	if err != nil {
		panic("Got an error. Very weird..." + err.Error())
	}

	lines, err := readLines(file)
	if err != nil {
		panic("Got an error. Very weird..." + err.Error())
	}
	return lines[rand.Intn(len(lines))]
}

func readLines(file []byte) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(string(file)))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func pluralizeNoun(noun string) string {
	if strings.HasSuffix(noun, "s") || strings.HasSuffix(noun, "sh") || strings.HasSuffix(noun, "ch") || strings.HasSuffix(noun, "x") || strings.HasSuffix(noun, "z") {
		return noun + "es"
	}
	if strings.HasSuffix(noun, "y") && !isVowelBeforeY(noun) {
		return noun[:len(noun)-1] + "ies"
	}
	if strings.HasSuffix(noun, "f") {
		return noun[:len(noun)-1] + "ves"
	}
	if strings.HasSuffix(noun, "fe") {
		return noun[:len(noun)-2] + "ves"
	}
	switch noun {
	case "child":
		return "children"
	case "mouse":
		return "mice"
	case "sheep":
		return "sheep"
	default:
		return noun + "s"
	}
}

func thridPersoniseVerb(verb string) string {
	if strings.HasSuffix(verb, "s") || strings.HasSuffix(verb, "sh") || strings.HasSuffix(verb, "ch") || strings.HasSuffix(verb, "x") || strings.HasSuffix(verb, "z") {
		return verb + "es"
	}
	if strings.HasSuffix(verb, "y") && !isVowelBeforeY(verb) {
		return verb[:len(verb)-1] + "ies"
	}
	return verb + "s"
}

func isVowelBeforeY(noun string) bool {
	if len(noun) < 2 {
		return false
	}
	switch noun[len(noun)-2] {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	default:
		return false
	}
}

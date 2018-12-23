package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type Link struct {
	Title string `json:"title"`
	Count int    `json:"count"`
}

type Note struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
	Links []Link   `json:"links"`
}

func NewNote(r io.Reader) *Note {
	var title string
	links := make([]string, 0)
	tags := make([]string, 0)

	tagsRegexp := regexp.MustCompile(`(^|\s)#[^\s]+`)
	linksRegexp := regexp.MustCompile(`\[\[([^]]+)]]`)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// Empty line
		if len(scanner.Text()) == 0 {
			continue
		}

		// Title line
		if strings.HasPrefix(scanner.Text(), "# ") {
			title = strings.TrimLeft(scanner.Text(), "# ")
			continue
		}

		// Tags line
		tokens := tagsRegexp.FindAllString(scanner.Text(), -1)
		if tokens != nil {
			tmp := tokens[:0]
			for i, t := range tokens {
				tokens[i] = strings.Replace(strings.TrimSpace(t), "#", "", -1)
				if len(tokens[i]) > 0 {
					tmp = append(tmp, tokens[i])
				}
			}
			tags = append(tags, tmp...)
			continue
		}

		// Content line with links optionally
		tokens = linksRegexp.FindAllString(scanner.Text(), -1)
		if tokens != nil {
			for i, t := range tokens {
				tokens[i] = strings.TrimFunc(strings.TrimSpace(t), func(r rune) bool {
					return r == '[' || r == ']'
				})
			}
			links = append(links, tokens...)
			continue
		}
	}

	return &Note{
		Title: title,
		Tags:  UniqStrings(tags),
		Links: ToLinks(FreqStrings(links)),
	}
}

func (n *Note) Marshal() (string, error) {
	b, err := json.Marshal(n)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// UniqStrings returns a unique subset of the string slice provided
func UniqStrings(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)
	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

// FreqStrings returns a frequency map of the string slice provided
func FreqStrings(input []string) map[string]int {
	m := make(map[string]int)
	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = 1
		} else {
			m[val] = m[val] + 1
		}
	}
	return m
}

// ToLinks returns a slice of Links from frequency map
func ToLinks(input map[string]int) []Link {
	arr := make([]Link, 0, len(input))
	for key, val := range input {
		arr = append(arr, Link{
			Title: key,
			Count: val,
		})
	}
	return arr
}

func main() {
	app := cli.NewApp()
	app.Name = "bear-tools"

	jsonCmd := cli.Command{
		Name:  "json",
		Usage: "Convert Bear App note from Markdown to JSON",
		Action: func(c *cli.Context) error {
			s, err := NewNote(os.Stdin).Marshal()
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(s)
			return nil
		},
	}

	app.Commands = []cli.Command{
		jsonCmd,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

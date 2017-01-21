package markov

import (
	"math/rand"
	"strings"
)

type Markov struct {
	entries map[string][]string
}

const (
	start_of_string = "START_OF_STRING_START_OF_STRING"
	end_of_string   = "END_OF_STRING_END_OF_STRING"
)

func New() *Markov {
	return &Markov{make(map[string][]string)}
}

func (markov *Markov) Update(s string) {
	var prev string = start_of_string
	for _, word := range strings.Split(s, " ") {
		markov.entries[prev] = append(markov.entries[prev], word)
		prev = word
	}
	markov.entries[prev] = append(markov.entries[prev], end_of_string)
}

func (markov *Markov) Generate(r *rand.Rand) string {
	var output = []string{}
	var word = start_of_string

	for {
		entry := markov.entries[word]
		word = entry[r.Intn(len(entry))]

		if word == end_of_string {
			break
		}

		output = append(output, word)
	}
	return strings.Join(output, " ")
}

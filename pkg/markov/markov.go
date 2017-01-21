package markov

import (
	"math/rand"
	"strings"
)

type MarkovEntry struct {
	next []string
}

type Markov struct {
	entries map[string]*MarkovEntry
}

const (
	start_of_string = "START_OF_STRING_START_OF_STRING"
	end_of_string   = "END_OF_STRING_END_OF_STRING"
)

func New() *Markov {
	return &Markov{make(map[string]*MarkovEntry)}
}

func insert_into_markov(markov *Markov, prev string, word string) {
	entry, ok := markov.entries[prev]

	if ok {
		entry.next = append(entry.next, word)
	} else {
		entry = &MarkovEntry{[]string{word}}
		markov.entries[prev] = entry
	}
}

func (markov *Markov) Update(s string) {
	var prev string = start_of_string
	for _, word := range strings.Split(s, " ") {
		insert_into_markov(markov, prev, word)
		prev = word
	}
	insert_into_markov(markov, prev, end_of_string)
}

func (markov *Markov) Generate(r *rand.Rand) string {
	var output = []string{}
	var word = start_of_string

	for {
		entry, _ := markov.entries[word]
		word = entry.next[r.Intn(len(entry.next))]

		if word == end_of_string {
			break
		}

		output = append(output, word)
	}
	return strings.Join(output, " ")
}

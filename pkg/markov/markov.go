package markov

import (
	"math/rand"
	"strings"
)

type Markov struct {
	entries map[string][]string
}

type Prefix struct {
	words [2]string
}

func (p *Prefix) key() string {
	if len(p.words[0]) > 0 && len(p.words[1]) > 0 {
		return p.words[0] + " " + p.words[1]
	} else if len(p.words[0]) > 0 {
		return p.words[0]
	} else {
		return ""
	}
}

func (p *Prefix) push(word string) {
	if len(p.words[0]) > 0 && len(p.words[1]) > 0 {
		p.words[0] = p.words[1]
		p.words[1] = word
	} else if len(p.words[0]) > 0 {
		p.words[1] = word
	} else {
		p.words[0] = word
	}
}

const (
	start_of_string = ""
	end_of_string = "\n"
)

func New() *Markov {
	return &Markov{make(map[string][]string)}
}

func (markov *Markov) Update(s string) {
	var prefix Prefix
	for _, word := range strings.Split(s, " ") {
		key := prefix.key()
		markov.entries[key] = append(markov.entries[key], word)
		prefix.push(word)
	}

	key := prefix.key()
	markov.entries[key] = append(markov.entries[key], end_of_string)
}

func (markov *Markov) Generate(r *rand.Rand) string {
	var output = []string{}
	var word = start_of_string

	var prefix Prefix
	for {
		key := prefix.key()
		entry := markov.entries[key]
		word = entry[r.Intn(len(entry))]

		if word == end_of_string {
			break
		}

		output = append(output, word)
		prefix.push(word)
	}
	return strings.Join(output, " ")
}

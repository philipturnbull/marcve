package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

type CVE struct {
	XMLName xml.Name `xml:"cve"`
	Items   []Item   `xml:"item"`
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Name        string   `xml:"name,attr"`
	Description string   `xml:"desc"`
}

func parse_xml(filename string) (*CVE, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cve_data, _ := ioutil.ReadAll(file)

	var cve *CVE
	xml.Unmarshal(cve_data, &cve)
	return cve, nil
}

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

func insert_into_markov(markov *Markov, prev string, word string) {
	entry, ok := markov.entries[prev]

	if ok {
		entry.next = append(entry.next, word)
	} else {
		entry = &MarkovEntry{[]string{word}}
		markov.entries[prev] = entry
	}
}

func (markov *Markov) update(s string) {
	var prev string = start_of_string
	for _, word := range strings.Split(s, " ") {
		insert_into_markov(markov, prev, word)
		prev = word
	}
	insert_into_markov(markov, prev, end_of_string)
}

func (markov *Markov) generate(r *rand.Rand) string {
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

func rand_from_string(s string) *rand.Rand {
	var seed int64

	s_hash := md5.Sum([]byte(s))
	buf := bytes.NewReader(s_hash[:8])
	binary.Read(buf, binary.LittleEndian, &seed)

	return rand.New(rand.NewSource(seed))
}

func main() {
	var cve_id string
	var xml_filename string

	flag.StringVar(&cve_id, "id", "CVE-2017-0001", "CVE identifier")
	flag.StringVar(&xml_filename, "filename", "allitems.xml", "CVE XML filename")

	flag.Parse()

	cve, err := parse_xml(xml_filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	var markov = Markov{make(map[string]*MarkovEntry)}
	for _, item := range cve.Items {
		markov.update(item.Description)
	}

	fmt.Println(markov.generate(rand_from_string(cve_id)))
}

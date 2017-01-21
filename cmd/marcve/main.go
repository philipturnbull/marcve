package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/philipturnbull/marcve/pkg/markov"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
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

func rand_from_string(s string) *rand.Rand {
	var seed int64

	s_hash := md5.Sum([]byte(s))
	buf := bytes.NewReader(s_hash[:8])
	binary.Read(buf, binary.LittleEndian, &seed)

	return rand.New(rand.NewSource(seed))
}

type CVEHandler struct {
	chain *markov.Markov
}

func (h *CVEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, h.chain.Generate(rand_from_string("CVE-"+vars["year"]+"-"+vars["id"])))
}

func run_server(chain *markov.Markov) {
	handler := &CVEHandler{chain}

	r := mux.NewRouter()
	r.Handle("/cve/{year}/{id}", handler)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func main() {
	var xml_filename string
	flag.StringVar(&xml_filename, "filename", "allitems.xml", "CVE XML filename")
	flag.Parse()

	cve, err := parse_xml(xml_filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	chain := markov.New()
	for _, item := range cve.Items {
		chain.Update(item.Description)
	}

	run_server(chain)
}

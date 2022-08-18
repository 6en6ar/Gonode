package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var torlistUrl string = "https://www.dan.me.uk/torlist/"
var tf string = "torlist.txt"

func CountLines(r io.Reader) (int, error) {
	buff := make([]byte, 32*1024)
	count := 0
	sep := []byte{'\n'}

	for {
		c, err := r.Read(buff)
		count += bytes.Count(buff[:c], sep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
func CountBytes(s []byte) int {
	nl := []byte{'\n'}
	n := bytes.Count(s, nl)
	if len(s) > 0 && !bytes.HasSuffix(s, nl) {
		n++
	}
	return n
}
func UpdateFile(body []byte) {
	err := os.WriteFile(tf, body, 0644)
	if err != nil {
		fmt.Println("Error updating file")
	}
}
func ReadListAndCompare() {
	res, err := http.Get(torlistUrl)
	if err != nil {
		fmt.Println("Can't get list, try later")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	f := OpenFile()
	defer f.Close()
	// Funkcija
	cBody := CountBytes(body)
	cFile, _ := CountLines(f)
	if cBody == cFile {
		fmt.Println("List of tor addresses is up to date ...")
		os.Exit(1)
		//Do nothing
	} else {
		UpdateFile(body)
		fmt.Printf("Total number of ip addresses --> %d", cBody)
		fmt.Println()
		ct := time.Now()
		fmt.Println("Last updated : ", ct.Format("2006.01.02 15:04:05"))

	}
}
func OpenFile() *os.File {
	f, err := os.Open(tf)
	if err != nil {
		fmt.Println("Error ocurred opening file!")
	}
	return f
}

var banner = `
   ______                      __   
  / ____/___  ____  ____  ____/ /__ 
 / / __/ __ \/ __ \/ __ \/ __  / _ \
/ /_/ / /_/ / / / / /_/ / /_/ /  __/
\____/\____/_/ /_/\____/\__,_/\___/ 

Coded by 6en6ar 3:)
								  

`

func main() {
	fmt.Println(banner)
	// can only be updated every 30 minutes
	ReadListAndCompare()
}

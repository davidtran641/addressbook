package main

import (
	"io/ioutil"
	"log"
	"os"
	"fmt"
	"io"
	addressbook "github.com/davidtran641/addressbook/pb"
	"github.com/golang/protobuf/proto"
)

func listPeople(w io.Writer, book *addressbook.AddressBook) {
	for _, p := range book.People {
		writePerson(w, p)
	}
}

func writePerson(w io.Writer, p *addressbook.Person) {
	fmt.Fprintln(w, "Person ID: ", p.Id)
	fmt.Fprintln(w, "	Name: ", p.Name)
	if p.Email != "" {
		fmt.Fprintln(w, "	Email: ", p.Email)
	}
	for _, pn := range p.Phones {
		switch pn.Type {
		case addressbook.Person_MOBILE:
			fmt.Fprintf(w, "	Mobile phone #: ")
		case addressbook.Person_HOME:
			fmt.Fprintf(w, "	Home phone #: ")
		case addressbook.Person_WORK:
			fmt.Fprintf(w, "	Work phone #: ")
		}
		fmt.Fprintln(w, pn.Number)
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <address book>\n", os.Args[0])
	}
	fname := os.Args[1]

	in, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalln("Error reading file: ", err)
	}

	book := &addressbook.AddressBook{}

	if err := proto.Unmarshal(in, book); err != nil {
		log.Fatalln("Failed to parse address book: ", err)
	}

	listPeople(os.Stdout, book)
}

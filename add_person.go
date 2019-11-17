package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"github.com/davidtran641/addressbook/pb"
	proto "github.com/golang/protobuf/proto"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <address book file>", os.Args[0])
	}
	fname := os.Args[1]

	in, err := ioutil.ReadFile(fname)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s: File Not found. Creating a new file", fname)
		} else {
			log.Fatalln("Error reading file: ", err)
		}
	}

	book := &addressbook.AddressBook{}
	if err := proto.Unmarshal(in, book); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	addr, err := promptForAddress(os.Stdin)
	if err != nil {
		log.Fatalln("Error with address: ", err)
	}
	book.People = append(book.People, addr)

	out, err := proto.Marshal(book)
	if err != nil {
		log.Fatalln("Failed to encode address book: ", err)
	}
	if err := ioutil.WriteFile(fname, out, 0644); err != nil {
		log.Fatalln("Failed to write address book: ", err)
	}
}

func promptForAddress(r io.Reader) (*addressbook.Person, error) {
	p := &addressbook.Person{}

	reader := bufio.NewReader(r)
	fmt.Print("Enter new person ID number: ")

	if _, err := fmt.Fscanf(reader, "%d\n", &p.Id); err != nil {
		return p, err
	}

	fmt.Print("Enter name:")
	name, err := reader.ReadString('\n')
	if err != nil {
		return p, err
	}

	p.Name = strings.TrimSpace(name)

	fmt.Print("Enter email address: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return p, err
	}
	p.Email = strings.TrimSpace(email)

	for {
		fmt.Print("Enter phone number Or leave blank to finish: ")
		phone, err := reader.ReadString('\n')
		if err != nil {
			return p, err
		}
		phone = strings.TrimSpace(phone)
		if phone == "" {
			break
		}

		pn := &addressbook.Person_PhoneNumber{Number: phone}

		fmt.Print("Is this a mobile, home or work phone? ")
		ptype, err := reader.ReadString('\n')
		if err != nil {
			return p, err
		}
		ptype = strings.TrimSpace(ptype)
		switch ptype {
		case "mobile":
			pn.Type = addressbook.Person_MOBILE
		case "home":
			pn.Type = addressbook.Person_HOME
		case "work":
			pn.Type = addressbook.Person_WORK
		default:
			fmt.Printf("Unknow phone type %v. Using default\n", ptype)
		}

		p.Phones = append(p.Phones, pn)

	}
	return p, nil
}

package charset_test

import (
	"bytes"
	"fmt"
	"github.com/rogpeppe/go-charset/charset"
	_ "github.com/rogpeppe/go-charset/data"
	"io/ioutil"
	"log"
	"strings"
)

func ExampleNewReader() {
	r, err := charset.NewReader("latin1", strings.NewReader("\xa35 for Pepp\xe9"))
	if err != nil {
		log.Fatal(err)
	}
	result, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", result)
	// Output: £5 for Peppé
}

func ExampleNewWriter() {
	buf := new(bytes.Buffer)
	w, err := charset.NewWriter("latin1", buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "£5 for Peppé")
	w.Close()
	fmt.Printf("%q\n", buf.Bytes())
	// Output: "\xa35 for Pepp\xe9"
}

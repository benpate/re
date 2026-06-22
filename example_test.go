package re_test

import (
	"fmt"
	"io"
	"net/http/httptest"
	"strings"

	"github.com/benpate/re"
)

// Example_reReader shows the core feature: a re.Reader can be read to completion
// more than once. After each read reaches EOF, the cursor rewinds so the next
// read starts again from the beginning.
func Example_reReader() {

	// A standard io.Reader can only be read once.
	singleReader := strings.NewReader("Hello World.")

	// Wrap it so it can be read repeatedly.
	multipleReader, err := re.NewReader(singleReader)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		body, _ := io.ReadAll(multipleReader)
		fmt.Println(string(body))
	}

	// Output:
	// Hello World.
	// Hello World.
	// Hello World.
}

// ExampleReader_Reset shows how Reset rewinds the cursor mid-stream, so a Reader
// that has only been partially consumed can be read again from the start.
func ExampleReader_Reset() {

	reader := re.NewReaderFromBytes([]byte("Hello World."))

	// Consume the first five bytes.
	first := make([]byte, 5)
	if _, err := reader.Read(first); err != nil {
		panic(err)
	}
	fmt.Printf("first read: %s\n", first)

	// Rewind, then read the whole thing.
	reader.Reset()
	all, _ := io.ReadAll(reader)
	fmt.Printf("after reset: %s\n", all)

	// Output:
	// first read: Hello
	// after reset: Hello World.
}

// ExampleReader_Read shows that draining a Reader through a buffer smaller than
// its contents returns the full content, in order, across multiple Read calls.
func ExampleReader_Read() {

	reader := re.NewReaderFromBytes([]byte("Hello World."))

	// Read through a small (4-byte) buffer until EOF.
	buffer := make([]byte, 4)
	for {
		n, err := reader.Read(buffer)
		fmt.Printf("%q\n", buffer[:n])
		if err == io.EOF {
			break
		}
	}

	// Output:
	// "Hell"
	// "o Wo"
	// "rld."
	// ""
}

// ExampleReadRequestBody shows that after reading a request body, the body is
// still readable by a later handler, because ReadRequestBody replaces it with a
// fresh reader.
func ExampleReadRequestBody() {

	request := httptest.NewRequest("POST", "https://example.com", strings.NewReader("request payload"))

	// One layer reads the body...
	body, _ := re.ReadRequestBody(request)
	fmt.Printf("first reader: %s\n", body)

	// ...and a later layer can still read the same body.
	again, _ := io.ReadAll(request.Body)
	fmt.Printf("second reader: %s\n", again)

	// Output:
	// first reader: request payload
	// second reader: request payload
}

package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("http://google.com")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Response:", resp) // Taking a naive common sense approach (that isn't actually right) to log out the response for now
}

// In the next section, we're going to start really digging into the documentation around this resp object right here and figure out how we can
// actually really print the body of the response back to the user

// Currently only printing out:
// Response: &{200 OK 200 HTTP/1.1 1 1 map[Cache-Control:[private, max-age=0] Content-Type:[text/html; charset=ISO-8859-1] Date:[Tue, 31 Mar 2020 07:12:33 GMT]
// Expires:[-1] P3p:[CP="This is not a P3P policy! See g.co/p3phelp for more info."] Server:[gws] Set-Cookie:[1P_JAR=2020-03-31-07; expires=Thu, 30-Apr-2020
// 07:12:33 GMT; path=/; domain=.google.com; Secure NID=201=VN3yh2oxloA6xuUstop5ftZ2v_YXge1K7p-YAsmkujdmbsik23v-jcWCPIjKPM8mEjYbSYbpt9fkM7ivEMXMM3hUGrLcxUMFk1Zl
// _TdLX1JU0lREMxnxbtsyW9uj1B5IJfXinQh_j6CzN9BBUZi6uqVZt6InXqnoj6Rrmtw7YHo; expires=Wed, 30-Sep-2020 07:12:32 GMT; path=/; domain=.google.com; HttpOnly]
// X-Frame-Options:[SAMEORIGIN] X-Xss-Protection:[0]] 0xc0001842a0 -1 [] false true map[] 0xc000192000 <nil>}

// https://golang.org/pkg/net/http/#Response shows that we're clearing logging out a struct right now. So looking further into the struct, we find something interesting:
// the Body property and it has a type of io.ReadCloser
// If we were working with a lot of other languages, definitely like Ruby and Javascript comes to mind, the body of the response is usually just a blob of text or
// maybe some Json or whatever it might be, just something that very clearly represents the body of the response. But here with Go, it's clearly that's not the case.
// https://golang.org/pkg/io/#ReadCloser which we then click then tells us very quickly that this ReadCloser thing is actually an interface
//    ReadCloser is the interface that groups the basic Read and Close methods.
//       type ReadCloser interface {
//          Reader
//          Closer
//       }
// https://golang.org/pkg/io/#Reader then shows something we're more familiar to, i.e.
//    type Reader interface {
//       Read(p []byte) (n int, err error)
//    }
// Also, check out https://golang.org/pkg/io/#Closer
//    type Closer interface {
//       Close() error
//    }

// In summary, we're currently working out:
// A Response Struct that has...
// - Status of type string
// - StatusCode of type int
// - Body of type io.ReadCloser
//
// io.ReadCloser Interface groups...
// - Reader
// - Closer
//
// io.Reader Interface has...
// - Read([]byte) (int, error)
//
// io.Closer Interface has...
// - Close() (error)
//
// Notes:
// 1. Firstly, why was a interface used as the type inside of a struct?
//    Basically, if we specify a interface as a value inside of a struct, we're saying that the body field right here can have any value assigned to it, so long as it fulfills this interface.
//    So this is kind of like a free lease on us, to provide us a little bit of flexibility. It's like saying you can put any type in here as long as it satisfies the ReadCloser interface.
//    Thus, in this case, all we really have to do is look at that ReadCloser interface, drill through the documentation and then eventually come to find that we need to define a function
//    called Read and one called Close.
//    - So if we sat down and made some type of struct that had a function called Read and one called Close that obeyed all these other types in here we could then freely make a response
//      struct and assign that type to this body field. That's why we are seeing an interface in the place of an actual type inside of a struct.
//
// 2. Next, why we saw the kind of funny syntax around the ReadCloser?
//    So specifically we were talking about interfaces just a little bit ago and we were a little bit more used to seeing syntax that we list out a function name the set of parentheses and
//    then the return type. However, it's clearly not the case here.
//    - Well, in Go, we can take multiple interfaces so different interfaces and assemble them together to form another interface. I.e. both Reader and Closer are interfaces.
//      The ReadCloser interface says if you want to fulfill to satisfy the requirements of this interface, you have to satisfy the requirements of both the Reader and Closer interfaces.
//      So all we're really doing here is defining a new interface by putting together pieces of other ones inside of our application so we can freely embed one interface into another
//      as much as we please, as long as it actually serves the purpose of building our application. So in reality, what really matters to us is what the Reader interface and Closer interface
//      are requiring of us.

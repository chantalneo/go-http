package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type logWriter struct{}

func main() {
	resp, err := http.Get("http://google.com")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// bs := make([]byte, 99999) // using a different syntax from []byte{}. Check notes' point 5
	// resp.Body.Read(bs)        // in reality we don't actually always make some bite slice like this and pass it off to the read function whenever we want to read data out of response
	// fmt.Println(string(bs))

	lw := logWriter{}

	io.Copy(lw, resp.Body)
}

func (logWriter) Write(bs []byte) (int, error) {
	fmt.Println(string(bs))
	fmt.Println("Just wrote this many bytes", len(bs))
	return len(bs), nil
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
//
// 3. Instructor threw a possibility (just a possibility, though neither is this how the world works nor how Go is written):
//    Source of Input                              | Returns?!? | To Print It... (Imagining a world without interfaces)
//    HTTP Request Body                            | []flargen  | func printHTTP([]flargen)
//    Text file on hard drive                      | []string   | func printFile([]string)
//    Image file on hard drive                     | jpegne     | func printImage(jpegne)
//    User entering text into command line         | []byte     | func printText([]byte)
//    Data from analog sensor plugged into machine | []float    | func printData([]float)
//
//    The solution to all this stuff is a Reader interface
//    No matter what the source of input is, we're going to see this reader interface all over the place
//    Now the purpose of the reader interface is to say hey we understand that there is a lot of different sources of data coming into your program,
//    but for each of these different sources if they implement the Reader interface then we get some very common interface or some common point of contact that we can
//    use to take that input and then pipe it off to different places inside of our codebase, without writing a bunch of custom functions to work with each of these individual return types
//    TL;DR
//    Reader interface can be considered as being literally an interface or an adapter to take these radically different sources of input and translate them into some common medium,
//    that all these other functions that we might have can easily work with
//
//    The Reader expects us to spit out a byte slice from the different sources, so that it becomes an output data that anyone can work with
//
//    A more realistic version would be:
//    Source of Input
//    HTTP Request Body                             [implements Reader] ----------\
//    Text file on hard drive                       [implements Reader] -----------\
//    Image file on hard drive                      [implements Reader] ------------>------------ []byte, the output data that anyone can work with
//    User entering text into command line          [implements Reader] -----------/
//    Data from analog sensor plugged into machine  [implements Reader] ----------/
//
// 4. What does the Read function do?
//    The entire point of all this interface stuff is to say that the request body has already taken care of this work for us...
//
//    Things that wants
//    to read the body
//    (something that
//    wants to see the
//    Reader interface)               Thing that implements Reader
//                                    Read([]byte) (int, err)
//    Byte Slice ----------------->           Byte Slice
//                                                /|\
//                                                 |
//                                       Raw body of response
//                         (the function takes the data from the raw body of
//                        response and injects/pushes it into that byte slice)
//
// Remember all that discussion around pointer, if we pass a slice off to a function, that function can freely modify the slice, and then the original slice
// back in our code world (before we pass it off to the function), which is the original slice gets modified. So that's what's really going on behind the scenes
// and that kind of explains the Read function signature right there.
//    Instructor: so even though we are passing in this byte size and I don't know about you, but if I was going to read for something in like a classic programming language,
//    I would expect to call Read and then be returned a byte slice, you know something that says hey we just read from the data source. Here is a bite slice that's coming back
//    that's full of all the data you want to care or you actually care about. But with this function right here the signature is just a little bit different to kind of work with
//    that whole concept of pointers and whatnot we spoke about earlier
//
// Putting myself in the equation, I want to read data out of the request body. So if I'm going to make direct access to this Read function, then I'd create a byte slice
// So that I pass it to the thing that is implementing the reader which for our application is the request body, i.e. the request body is implementing the Reader interface
//
// A little about the about two return arguments on this thing are the two return values:
// - The int is referring to the number of bytes that was read into that slice, which we can use for error checking or a little bit of something that says hey here's
// how much data we just shoved into the slice
// - The err is basically an error object that says okay well hey you know here's maybe something went wrong something didn't go quite right
//
// 5. There we're making a slice of type byte and make sure that there are 99999 elements available inside of it
//    Yes, a bite slice can grow and shrink, but the read function is not set up to automatically resize the slice, if the slices are already full
//    So instead we take this approach of just making an arbitrarily large bite slice that's basically big enough for all this data to fit into
//
// 6. Source of data -> Reader -> []byte (Output data that anyone can work with)
//    []byte -> Writer -> Some form of output
//
//    We can kind of think of this writer interface as doing something like this. We take our bite slice we pass it to some value that implements the writer interface. The writer interface
//    is essentially describing something that can take some info inside of our program and send it outside of our program
//
//    []byte -> Writer -> Source of output
//                        - Outgoing HTTP Request
//                        - Text file on hard drive
//                        - Image file on hard drive
//                        - Terminal
//
//    So we need to find something in the standard library that implements the Writer interface, and use that to log out all the data that we're receiving from the Reader!
//
// 7. In relation to the conclusion of point 6
//    https://golang.org/pkg/io/#Writer
//    https://golang.org/pkg/io/#Copy
//
//									(First Argument)                                         (Second Argument)
//    io.Copy        Something that implements the writer interface           Something that implements the Reader interface
//                                          V                                                         V
//                                      os.Stdout                                                  resp.Body
//                                          V
//                                 value of type File
//                                          V
//                         File has a function called 'Write'
//                                          V
//                    Therefore, it implements the Writer interface
//
// That's why we can pass os.Stdout to io.Copy
//
// 7. Actual Source Code of io.Copy
//    func Copy(dst Writer, src Reader) (written int64, err error) {
// 	     return copyBuffer(dst, src, nil)
//    }
//
//    func CopyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {
// 	     if buf != nil && len(buf) == 0 {
// 		    panic("empty buffer in io.CopyBuffer")
// 	     }
// 	     return copyBuffer(dst, src, buf)
//    }
//
// 	  func copyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {
//       // If the reader has a WriteTo method, use it to do the copy.
// 	     // Avoids an allocation and a copy.
// 	     if wt, ok := src.(WriterTo); ok {
// 		    return wt.WriteTo(dst)
// 	     }
// 	     // Similarly, if the writer has a ReadFrom method, use it to do the copy.
// 	     if rt, ok := dst.(ReaderFrom); ok {
// 		    return rt.ReadFrom(src)
// 	     }
// 	     if buf == nil {
// 		    size := 32 * 1024
// 		    if l, ok := src.(*LimitedReader); ok && int64(size) > l.N {
// 			   if l.N < 1 {
// 				  size = 1
// 			   } else {
// 			      size = int(l.N)
// 			   }
// 		    }
// 		    buf = make([]byte, size)
// 	     }
// 	     for {
// 		    nr, er := src.Read(buf)
// 		    if nr > 0 {
// 			   nw, ew := dst.Write(buf[0:nr])
// 			      if nw > 0 {
// 				     written += int64(nw)
// 			      }
// 			      if ew != nil {
// 				     err = ew
// 				     break
// 			      }
// 			      if nr != nw {
// 				     err = ErrShortWrite
// 				     break
// 			      }
// 		    }
// 		    if er != nil {
// 			   if er != EOF {
// 			      err = er
// 			   }
// 			   break
// 		    }
// 	      }
// 	      return written, err
//     }

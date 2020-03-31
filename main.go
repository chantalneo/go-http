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

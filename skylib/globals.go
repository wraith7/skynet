//Copyright (c) 2011 Brian Ketelsen

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package skylib

import (
	"path"
	"flag"
	"os"
	"expvar"
)



var Port *int = flag.Int("port", 0, "tcp port to listen")
var Name *string = flag.String("name", os.Args[0], "name of this server")
var BindIP *string = flag.String("bindaddress", "127.0.0.1", "address to bind")
var LogFileName *string = flag.String("logFileName", "myservice." + path.Base(os.Args[0]) + ".log", "name of logfile")
var LogLevel *int = flag.Int("logLevel", 1, "log level (1-5)")
var Protocol *string = flag.String("protocol", "http+gob", "RPC message transport protocol (default is http+gob; try json")
var Requests *expvar.Int
var Errors *expvar.Int
var Goroutines *expvar.Int

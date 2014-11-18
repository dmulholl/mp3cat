package main

import "fmt"
import "os"
import "io"


// Flag controlling the display of debugging information.
var debugMode = false


// debug prints debugging information to stderr.
func debug(message string) {
    if debugMode {
        fmt.Fprintln(os.Stderr, message)
    }
}


// fillBuffer attemtps to read len(buffer) bytes from the input stream.
// Returns a boolean indicating success.
func fillBuffer(stream io.Reader, buffer []byte) bool {
    n, _ := io.ReadFull(stream, buffer)
    if n < len(buffer) {
        return false
    }
    return true
}

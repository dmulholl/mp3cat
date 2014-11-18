package mp3lib


import (
    "fmt"
    "os"
    "io"
)


// Flag controlling the display of debugging information.
var DebugMode = false


// debug prints debugging information to stderr.
func debug(message string) {
    if DebugMode {
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

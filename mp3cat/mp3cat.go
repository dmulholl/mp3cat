/*
    Command line utility for concatenating MP3 files without re-encoding.
*/
package main


import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "github.com/dmulholland/mp3cat/mp3lib"
    "github.com/dmulholland/clio/go/clio"
)


// Application version number.
const version = "2.1.0"


// Command line help text.
var helptext = fmt.Sprintf(`
Usage: %s [FLAGS] [OPTIONS] ARGUMENTS

  Concatenates MP3 files without re-encoding. Supports both constant bit rate
  (CBR) and variable bit rate (VBR) files. Strips ID3 tags and garbage data
  from the output.

Arguments:
  <files>           List of input files to merge.

Options:
  -o, --out <file>  Output filename. Defaults to 'output.mp3'.

Flags:
  -f, --force       Overwrite an existing output file.
      --help        Display this help text and exit.
  -v, --verbose     Report progress.
      --version     Display the application's version number and exit.
`, filepath.Base(os.Args[0]))


// Application entry point.
func main() {

    // Initialize an argument parser.
    parser := clio.NewParser(helptext, version)

    // Register flags.
    parser.AddFlag("force", 'f')
    parser.AddFlag("verbose", 'v')
    parser.AddFlag("debug")

    // Register options.
    parser.AddStrOpt("out", "output.mp3", 'o')

    // Parse the command line arguments.
    parser.Parse()

    // Make sure we have a list of input files.
    if !parser.HasArgs() {
        fmt.Fprintln(os.Stderr, "Error: you must supply a list of files to merge.")
        os.Exit(1)
    }

    // Set debug mode if the user supplied a --debug flag.
    if parser.GetFlag("debug") {
        mp3lib.DebugMode = true
    }

    // Merge the input files.
    mergeFiles(
        parser.GetStrOpt("out"),
        parser.GetArgs(),
        parser.GetFlag("force"),
        parser.GetFlag("verbose"))
}


// Create a new file at the specified output path containing the merged
// contents of the list of input files.
func mergeFiles(outputPath string, inputPaths []string, overwrite bool, verbose bool) {

    var totalFrames uint32
    var totalBytes uint32
    var totalFiles int
    var firstBitRate int
    var isVBR bool

    // Only overwrite an existing file if the --force flag has been used.
    if _, err := os.Stat(outputPath); err == nil {
        if !overwrite {
            fmt.Fprintf(os.Stderr, "Error: the file '%v' already exists. ", outputPath)
            fmt.Fprintf(os.Stderr, "Use --force to overwrite it.\n")
            os.Exit(1)
        }
    }

    // If the list of input files includes the output file we'll end up
    // in an infinite loop.
    for _, filepath := range inputPaths {
        if filepath == outputPath {
            fmt.Fprintln(os.Stderr, "Error: the list of input files includes the output file.")
            os.Exit(1)
        }
    }

    // Create the output file.
    outputFile, err := os.Create(outputPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    // Loop over the input files and append their MP3 frames to the
    // output file.
    for _, filepath := range inputPaths {

        if verbose {
            fmt.Println("Merging:", filepath)
        }

        inputFile, err := os.Open(filepath)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }

        isFirstFrame := true

        for {

            // Read the next frame from the input file.
            frame := mp3lib.NextFrame(inputFile)
            if frame == nil {
                break
            }

            // Skip the first frame if it's a VBR header.
            if isFirstFrame {
                isFirstFrame = false
                if mp3lib.IsXingHeader(frame) || mp3lib.IsVbriHeader(frame) {
                    continue
                }
            }

            // If we detect more than one bitrate we'll need to add
            // a VBR header to the output file.
            if firstBitRate == 0 {
                firstBitRate = frame.BitRate
            } else if frame.BitRate != firstBitRate {
                isVBR = true
            }

            // Write the frame to the output file.
            _, err := outputFile.Write(frame.RawBytes)
            if err != nil {
                fmt.Fprintln(os.Stderr, err)
                os.Exit(1)
            }

            totalFrames += 1
            totalBytes += uint32(len(frame.RawBytes))
        }

        inputFile.Close()
        totalFiles += 1
    }

    outputFile.Close()

    // If we detected multiple bitrates, prepend a VBR header to the file.
    if isVBR {
        if verbose {
            fmt.Println("VBR data detected. Adding Xing header.")
        }
        addXingHeader(outputPath, totalFrames, totalBytes)
    }

    if verbose {
        fmt.Printf("%v files merged.\n", totalFiles)
    }
}


// Prepend an Xing VBR header to the specified MP3 file.
func addXingHeader(filepath string, totalFrames, totalBytes uint32) {

    outputFile, err := os.Create(filepath + ".mp3cat.tmp")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    inputFile, err := os.Open(filepath)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    templateFrame := mp3lib.NextFrame(inputFile)
    inputFile.Seek(0, 0)

    xingHeader := mp3lib.NewXingHeader(templateFrame, totalFrames, totalBytes)

    _, err = outputFile.Write(xingHeader.RawBytes)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    _, err = io.Copy(outputFile, inputFile)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    outputFile.Close()
    inputFile.Close()

    err = os.Remove(filepath)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    err = os.Rename(filepath + ".mp3cat.tmp", filepath)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

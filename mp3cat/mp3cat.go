/*
    Command line utility for concatenating MP3 files without re-encoding.

    Supports both constant bit rate (CBR) and variable bit rate (VBR) files.
    Strips ID3 tags and garbage data from the output.

      * Author: Darren Mulholland <dmulholland@outlook.ie>
      * License: Public Domain

*/
package main


import (
    "fmt"
    "io"
    "os"
    "flag"
    "github.com/dmulholland/mp3cat/mp3lib"
)


const version = "1.0.0"


const usage = `Usage: mp3cat [FLAGS] ARGUMENTS

Concatenates MP3 files without re-encoding. Supports both CBR and VBR files.
Strips ID3 tags and garbage data from the output.

Arguments:

  <outfile>        Output file.
  <infiles>        List of input files to concatenate.

Flags:

  -f, --force      Force overwriting of existing output files.
  -v, --verbose    Report progress.
  --help           Display this help text and exit.
  --version        Display version number and exit.`


var helpFlag = flag.Bool("help", false, "print help text and exit")
var versionFlag = flag.Bool("version", false, "print version and exit")
var debugFlag = flag.Bool("debug", false, "print debug information")
var forceFlag = flag.Bool("force", false, "force overwriting of existing output files")
var verboseFlag = flag.Bool("verbose", false, "increase verbosity")


func init() {
    flag.BoolVar(forceFlag, "f", false, "force overwriting of existing output files")
    flag.BoolVar(verboseFlag, "v", false, "increase verbosity")
}


func main() {

    flag.Usage = func() {
        fmt.Println()
        fmt.Println(usage)
    }

    flag.Parse()

    if *helpFlag {
        fmt.Println(usage)
        os.Exit(0)
    }

    if *versionFlag {
        fmt.Println(version)
        os.Exit(0)
    }

    if *debugFlag {
        mp3lib.DebugMode = true
    }

    if flag.NArg() < 2 {
        fmt.Fprintln(os.Stderr, "Error: you must supply at least two arguments.\n")
        fmt.Fprintln(os.Stderr, usage)
        os.Exit(1)
    }

    mergeFiles(flag.Arg(0), flag.Args()[1:])
}


func mergeFiles(outputPath string, inputPaths []string) {

    var totalFrames uint32
    var totalBytes uint32
    var totalFiles int
    var firstBitRate int
    var isVBR bool

    if _, err := os.Stat(outputPath); err == nil {
        if !(*forceFlag) {
            fmt.Fprintf(os.Stderr, "Error: \"%v\" already exists.\n", outputPath)
            fmt.Fprintf(os.Stderr, "Use the --force flag to overwrite it.\n")
            os.Exit(1)
        }
    }

    for _, filepath := range inputPaths {
        if filepath == outputPath {
            fmt.Fprintln(os.Stderr, "Error: the list of input files includes the output file.")
            os.Exit(1)
        }
    }

    outputFile, err := os.Create(outputPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    for _, filepath := range inputPaths {

        if *verboseFlag {
            fmt.Println("Merging:", filepath)
        }

        inputFile, err := os.Open(filepath)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }

        isFirstFrame := true

        for {
            frame := mp3lib.NextFrame(inputFile)
            if frame == nil {
                break
            }

            if isFirstFrame {
                isFirstFrame = false
                if mp3lib.IsXingHeader(frame) || mp3lib.IsVbriHeader(frame) {
                    continue
                }
            }

            if firstBitRate == 0 {
                firstBitRate = frame.BitRate
            } else if frame.BitRate != firstBitRate {
                isVBR = true
            }

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

    if isVBR {
        if *verboseFlag {
            fmt.Println("VBR data detected. Adding Xing header.")
        }
        addXingHeader(outputPath, totalFrames, totalBytes)
    }

    if *verboseFlag {
        fmt.Printf("%v files merged.\n", totalFiles)
    }
}


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

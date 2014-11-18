/*
    Command line utility for concatenating MP3 files without re-encoding.

    Author: Darren Mulholland <dmulholland@outlook.ie>
    License: Public Domain
*/

package main


import (
    "fmt"
    "os"
    "flag"
)


var version = "0.1.0"


var usage = `Usage: mp3cat [FLAGS] ARGUMENTS

Arguments:

  <output-file>     Output filename.
  <input-files>     List of input files to concatenate.

Flags:

  --help            Display this help text and exit.
  --version         Display version number and exit.`


func main() {

    var helpFlag = flag.Bool("help", false, "print help text and exit")
    var versionFlag = flag.Bool("version", false, "print version and exit")
    var debugFlag = flag.Bool("debug", false, "print debug information")

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
        debugMode = true
    }

    if flag.NArg() < 2 {
        fmt.Fprintln(os.Stderr, "error: too few arguments\n")
        fmt.Fprintln(os.Stderr, usage)
        os.Exit(1)
    }

    outputPath := flag.Arg(0)
    inputPaths := flag.Args()[1:]

    outputFile, err := os.Create(outputPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer outputFile.Close()

    totalFrames := 0
    firstBitRate := 0
    isVBR := false

    for _, filepath := range inputPaths {

        inputFile, err := os.Open(filepath)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        defer inputFile.Close()

        isFirstFrame := true

        for {
            frame := nextFrame(inputFile)
            if frame == nil {
                break
            }

            if isFirstFrame {
                isFirstFrame = false
                if isXingHeader(frame) || isVbriHeader(frame) {
                    debug("skipping vbr header")
                    continue
                }
            }

            if firstBitRate == 0 {
                firstBitRate = frame.BitRate
            } else if firstBitRate != frame.BitRate {
                isVBR = true
            }

            _, err := outputFile.Write(frame.RawBytes)
            if err != nil {
                fmt.Fprintln(os.Stderr, err)
                os.Exit(1)
            }

            totalFrames += 1
        }

        inputFile.Close()
    }

    outputFile.Close()

    if isVBR {
        // We need to rewrite the file, adding an Xing header at the front.
        debug("vbr file")
    }

    debug(fmt.Sprintf("total frames: %v", totalFrames))
}

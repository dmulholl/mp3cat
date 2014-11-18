/*
    Command line utility for concatenating MP3 files without re-encoding.

    Author: Darren Mulholland <dmulholland@outlook.ie>
    License: Public Domain
*/

package main


import (
    "fmt"
    "io"
    "os"
    "flag"
)


var version = "0.2.0"


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

    var totalFrames uint32
    var totalBytes uint32

    firstBitRate := 0
    isVBR := false

    for _, filepath := range inputPaths {

        inputFile, err := os.Open(filepath)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }

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
            totalBytes += uint32(len(frame.RawBytes))
        }

        inputFile.Close()
    }

    outputFile.Close()

    if isVBR {

        // We need to rewrite the file, adding an Xing header at the front.
        debug("vbr file")

        outputFile, err := os.Create(outputPath + ".tmp")
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }

        inputFile, err := os.Open(outputPath)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }

        firstFrame := nextFrame(inputFile)
        inputFile.Seek(0, 0)

        xingHeader := newXingHeader(firstFrame, totalFrames, totalBytes)

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

        inputFile.Close()
        outputFile.Close()

        err = os.Remove(outputPath)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }

        err = os.Rename(outputPath + ".tmp", outputPath)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
    }

    debug(fmt.Sprintf("total frames: %v", totalFrames))
}

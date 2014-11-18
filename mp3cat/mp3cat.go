/*
    Command line utility for concatenating MP3 files without re-encoding.

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


const version = "0.3.0"


const usage = `Usage: mp3cat [FLAGS] ARGUMENTS

Arguments:

  <outfile>    Output filename.
  <infiles>    List of input files to concatenate.

Flags:

  --help       Display this help text and exit.
  --version    Display version number and exit.`


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
        mp3lib.DebugMode = true
    }

    if flag.NArg() < 2 {
        fmt.Fprintln(os.Stderr, "error: too few arguments\n")
        fmt.Fprintln(os.Stderr, usage)
        os.Exit(1)
    }

    mergeFiles(flag.Arg(0), flag.Args()[1:])
}


func mergeFiles(outputPath string, inputPaths []string) {

    var totalFrames uint32
    var totalBytes uint32
    var firstBitRate int
    var isVBR bool

    outputFile, err := os.Create(outputPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    for _, filepath := range inputPaths {

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
    }

    outputFile.Close()

    if isVBR {
        addXingHeader(outputPath, totalFrames, totalBytes)
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

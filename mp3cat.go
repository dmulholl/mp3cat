package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dmulholl/argo/v4"
	"github.com/dmulholl/mp3lib"
	"golang.org/x/term"
)

const version = "4.3.0"

var helptext = fmt.Sprintf(`
Usage: %s [files]

  This tool concatenates MP3 files without re-encoding. Input files can be
  specified as a list of filenames:

    $ mp3cat one.mp3 two.mp3 three.mp3

  Alternatively, an entire directory of .mp3 files can be concatenated:

    $ mp3cat --dir /path/to/directory

Arguments:
  [files]                 List of files to merge.

Options:
  -d, --dir <path>        Directory of files to merge.
  -m, --meta <n>          Copy ID3 metadata from the n-th input file.
  -o, --out <path>        Output filepath. Defaults to 'output.mp3'.

Flags:
  -f, --force             Overwrite an existing output file.
  -h, --help              Display this help text and exit.
  -q, --quiet             Quiet mode. Only output error messages.
  -v, --version           Display the version number and exit.
`, filepath.Base(os.Args[0]))

func main() {
	// Parse the command line arguments.
	parser := argo.NewParser()
	parser.Helptext = helptext
	parser.Version = version
	parser.NewFlag("force f")
	parser.NewFlag("quiet q")
	parser.NewFlag("debug")
	parser.NewStringOption("out o", "output.mp3")
	parser.NewStringOption("dir d", "")
	parser.NewStringOption("interlace i", "")
	parser.NewIntOption("meta m", 0)

	if err := parser.ParseOsArgs(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s.\n", err)
		os.Exit(1)
	}

	// Make sure we have a list of files to merge.
	var files []string
	if parser.Found("dir") {
		err := filepath.Walk(
			parser.StringValue("dir"),
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if strings.ToLower(filepath.Ext(info.Name())) == ".mp3" {
					files = append(files, path)
				}
				return nil
			})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		if len(files) == 0 {
			fmt.Fprintln(os.Stderr, "Error: no files found.")
			os.Exit(1)
		}
	} else if len(parser.Args) > 0 {
		files = parser.Args
	} else {
		fmt.Fprintln(os.Stderr, "Error: you must specify files to merge.")
		os.Exit(1)
	}

	// Are we copying the ID3 tag from the n-th input file?
	var tagpath string
	if parser.Found("meta") {
		tagindex := parser.IntValue("meta") - 1
		if tagindex < 0 || tagindex > len(files)-1 {
			fmt.Fprintln(os.Stderr, "Error: --meta argument is out of range.")
			os.Exit(1)
		}
		tagpath = files[tagindex]
	}

	// Are we interlacing a spacer file?
	if parser.Found("interlace") {
		files = interlace(files, parser.StringValue("interlace"))
	}

	// Make sure all the files in the list actually exist.
	validateFiles(files)

	// Set debug mode if the user supplied a --debug flag.
	if parser.Found("debug") {
		mp3lib.DebugMode = true
	}

	// Merge the input files.
	merge(
		parser.StringValue("out"),
		tagpath,
		files,
		parser.Found("force"),
		parser.Found("quiet"))
}

// Check that all the files in the list exist.
func validateFiles(files []string) {
	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			fmt.Fprintf(os.Stderr, "Error: the file '%v' does not exist.\n", file)
			os.Exit(1)
		}
	}
}

// Interlace a spacer file between each file in the list.
func interlace(files []string, spacer string) []string {
	var interlaced []string
	for _, file := range files {
		interlaced = append(interlaced, file)
		interlaced = append(interlaced, spacer)
	}
	return interlaced[:len(interlaced)-1]
}

// Create a new file at [outpath] containing the merged contents of the list of input files.
func merge(outpath, tagpath string, inpaths []string, force, quiet bool) {
	var totalFrames uint32
	var totalBytes uint32
	var totalFiles int
	var firstBitRate int
	var isVBR bool

	// Only overwrite an existing file if the --force flag has been used.
	if _, err := os.Stat(outpath); err == nil {
		if !force {
			fmt.Fprintf(os.Stderr, "Error: the file '%v' already exists.\n", outpath)
			os.Exit(1)
		}
	}

	// If the list of input files includes the output file we'll end up in an infinite loop.
	for _, filepath := range inpaths {
		if filepath == outpath {
			fmt.Fprintln(os.Stderr, "Error: the list of input files includes the output file.")
			os.Exit(1)
		}
	}

	// Create the output file.
	outfile, err := os.Create(outpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if !quiet {
		printLine()
	}

	// Loop over the input files and append their MP3 frames to the output file.
	for _, inpath := range inpaths {
		if !quiet {
			fmt.Println("+", inpath)
		}

		infile, err := os.Open(inpath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		isFirstFrame := true

		for {
			// Read the next frame from the input file.
			frame := mp3lib.NextFrame(infile)
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

			// If we detect more than one bitrate we'll need to add a VBR header to the output file.
			if firstBitRate == 0 {
				firstBitRate = frame.BitRate
			} else if frame.BitRate != firstBitRate {
				isVBR = true
			}

			// Write the frame to the output file.
			_, err := outfile.Write(frame.RawBytes)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}

			totalFrames += 1
			totalBytes += uint32(len(frame.RawBytes))
		}

		infile.Close()
		totalFiles += 1
	}

	outfile.Close()
	if !quiet {
		printLine()
	}

	// If we detected multiple bitrates, prepend a VBR header to the file.
	if isVBR {
		if !quiet {
			fmt.Println("• Multiple bitrates detected. Adding VBR header.")
		}
		addXingHeader(outpath, totalFrames, totalBytes)
	}

	// Copy the ID3v2 tag from the n-th input file if requested. Order of operations is important
	// here. The ID3 tag must be the first item in the file - in particular, it must come *before*
	// any VBR header.
	if tagpath != "" {
		if !quiet {
			fmt.Printf("• Copying ID3 tag from: %s\n", tagpath)
		}
		addID3v2Tag(outpath, tagpath)
	}

	// Print a count of the number of files merged.
	if !quiet {
		fmt.Printf("• %v files merged.\n", totalFiles)
		printLine()
	}
}

// Prepend an Xing VBR header to the specified MP3 file.
func addXingHeader(filepath string, totalFrames, totalBytes uint32) {

	outputFile, err := os.Create(filepath + ".mp3cat.tmp")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	inputFile, err := os.Open(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	xingHeader := mp3lib.NewXingHeader(totalFrames, totalBytes)

	_, err = outputFile.Write(xingHeader.RawBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	outputFile.Close()
	inputFile.Close()

	err = os.Remove(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	err = os.Rename(filepath+".mp3cat.tmp", filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

// Prepend an ID3v2 tag to the MP3 file at mp3Path, copying from tagPath.
func addID3v2Tag(mp3Path, tagPath string) {

	tagFile, err := os.Open(tagPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	id3tag := mp3lib.NextID3v2Tag(tagFile)
	tagFile.Close()

	if id3tag != nil {
		outputFile, err := os.Create(mp3Path + ".mp3cat.tmp")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		inputFile, err := os.Open(mp3Path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		_, err = outputFile.Write(id3tag.RawBytes)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		_, err = io.Copy(outputFile, inputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		outputFile.Close()
		inputFile.Close()

		err = os.Remove(mp3Path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		err = os.Rename(mp3Path+".mp3cat.tmp", mp3Path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	}
}

// Print a line to stdout if we're running in a terminal.
func printLine() {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		width, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err == nil {
			if runtime.GOOS == "windows" {
				for i := 0; i < width; i++ {
					fmt.Print("-")
				}
				fmt.Println()
			} else {
				fmt.Print("\u001B[90m")
				for i := 0; i < width; i++ {
					fmt.Print("─")
				}
				fmt.Println("\u001B[0m")
			}
		}
	}
}

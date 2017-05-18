/*
   MP3Cat is a fast command line utility for concatenating MP3 files
   without re-encoding. It supports both constant bit rate (CBR) and
   variable bit rate (VBR) files.
*/
package main

import (
	"fmt"
	"github.com/dmulholland/clio/go/clio"
	"github.com/dmulholland/mp3lib"
	id3 "github.com/mikkyang/id3-go"
	v2 "github.com/mikkyang/id3-go/v2"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Application version number.
const version = "2.2.1"

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
  -i, --id3         Copy ID3 tags from the first file.
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
	parser.AddFlag("force f")
	parser.AddFlag("verbose v")
	parser.AddFlag("debug d")
	parser.AddFlag("id3 i")

	// Register options.
	parser.AddStr("out o", "output.mp3")

	// Parse the command line arguments.
	parser.Parse()

	// Make sure we have a list of input files.
	if !parser.HasArgs() {
		fmt.Fprintln(
			os.Stderr,
			"Error: you must supply a list of files to merge.")
		os.Exit(1)
	}

	// Set debug mode if the user supplied a --debug flag.
	if parser.GetFlag("debug") {
		mp3lib.DebugMode = true
	}

	// Merge the input files.
	mergeFiles(
		parser.GetStr("out"),
		parser.GetArgs(),
		parser.GetFlag("force"),
		parser.GetFlag("verbose"),
		parser.GetFlag("id3"))
}

// Create a new file at the specified output path containing the merged
// contents of the list of input files.
func mergeFiles(outputPath string, inputPaths []string, force, verbose bool, addId3 bool) {

	var totalFrames uint32
	var totalBytes uint32
	var totalFiles int
	var firstBitRate int
	var isVBR bool

	// Only overwrite an existing file if the --force flag has been used.
	if _, err := os.Stat(outputPath); err == nil {
		if !force {
			fmt.Fprintf(
				os.Stderr,
				"Error: the file '%v' already exists.\n", outputPath)
			os.Exit(1)
		}
	}

	// If the list of input files includes the output file we'll end up in an
	// infinite loop.
	for _, filepath := range inputPaths {
		if filepath == outputPath {
			fmt.Fprintln(
				os.Stderr,
				"Error: the list of input files includes the output file.")
			os.Exit(1)
		}
	}

	// Create the output file.
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Loop over the input files and append their MP3 frames to the output file.
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

			// If we detect more than one bitrate we'll need to add a VBR
			// header to the output file.
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

	title, artist, album, year, genre, sets, comments /*, image*/ := getId3Info(inputPaths[0])

	if addId3 {
		mp3File, err := id3.Open(outputPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer mp3File.Close()
		mp3File.SetTitle(title)
		mp3File.SetArtist(artist)
		mp3File.SetAlbum(album)
		mp3File.SetYear(year)
		mp3File.SetGenre(genre)

		if sets != "" {
			TPOS := v2.NewTextFrame(v2.V23FrameTypeMap["TPOS"], sets)
			mp3File.AddFrames(TPOS)
		}

		TRCK := v2.NewTextFrame(v2.V23FrameTypeMap["TRCK"], "1/1")
		mp3File.AddFrames(TRCK)

		if comments != nil {
			COMM := v2.NewUnsynchTextFrame(v2.V23FrameTypeMap["COM"], "Comments", strings.Join(comments, ""))
			mp3File.AddFrames(COMM)
		}

		/* This code isn't present and I coundn't find any workarounds using NewDataFrame to ge tit to work.  I even tried editing the code in github.com/dmulholland/mp3lib/v2/frame.go without much luck.
		if image != nil {
			APIC := v2.NewImageFrame(v2.V23FrameTypeMap["APIC"], image)
			mp3File.AddFrames(APIC)
		}*/
	}

	if verbose {
		fmt.Printf("%v files merged.\n", totalFiles)
		if addId3 {
			fmt.Printf("\t TIT2=%v\n", title)
			fmt.Printf("\t TPE1=%v\n", artist)
			fmt.Printf("\t TALB=%v\n", album)
			fmt.Printf("\t TYER=%v\n", year)
			fmt.Printf("\t TCON=%v\n", genre)
			fmt.Printf("\t TPOS=%v\n", sets)
			fmt.Printf("\t TRCK=%v\n", "1/1")
			//fmt.Printf("\t APIC=%v\n", len(image.Bytes()))
		}
	}
}

func getId3Info(filepath string) (title string, artist string, album string, year string, genre string, sets string, comments []string /*, image *v2.ImageFrame*/) {
	mp3File, err := id3.Open(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer mp3File.Close()
	title = mp3File.Title()
	artist = mp3File.Artist()
	album = mp3File.Album()
	year = mp3File.Year()
	genre = mp3File.Genre()
	comments = mp3File.Comments()
	TPOS := mp3File.Frame("TPOS") // v2
	if TPOS != nil {
		sets = TPOS.String()
	} else {
		TPOS = mp3File.Frame("TPA") //v1
		if TPOS != nil {
			sets = TPOS.String()
		}
	}
	/* No reason to read the image and I can't seam to save it.
	APIC := mp3File.Frame("APIC")
	if APIC != nil {
		image = APIC.(*v2.ImageFrame)
	} else {
		APIC = mp3File.Frame("PIC")
		if APIC != nil {
			image = APIC.(*v2.ImageFrame)
		}
	}*/

	return
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

	err = os.Rename(filepath+".mp3cat.tmp", filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}


# mp3cat

A command line utility for joining MP3 files without re-encoding. Supports both constant bit rate (CBR) and variable bit rate (VBR) files.


## Usage

    Usage: mp3cat [FLAGS] [OPTIONS] ARGUMENTS

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

Example:

    $ mp3cat -o merged.mp3 one.mp3 two.mp3 three.mp3


## Installation

If you have Go installed you can run:

    $ go get github.com/dmulholland/mp3cat/mp3cat

This will download, compile, and install the latest version of the application to your `$GOPATH/bin` directory.

Alternatively, you can download a precompiled binary from the [releases](https://github.com/dmulholland/mp3cat/releases) page.


## License

This work has been placed in the public domain.

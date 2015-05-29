
# mp3cat

A command line utility for joining MP3 files without re-encoding. Supports both constant bit rate (CBR) and variable bit rate (VBR) files.


## Usage

    Usage: mp3cat [FLAGS] ARGUMENTS

      Concatenates MP3 files without re-encoding. Supports both constant bit rate
      (CBR) and variable bit rate (VBR) files. Strips ID3 tags and garbage data
      from the output.

    Arguments:
      <outfile>        Output file.
      <infiles>        List of input files to merge.

    Flags:
      -f, --force      Force overwriting of existing output files.
      -v, --verbose    Report progress.
      --help           Display this help text and exit.
      --version        Display version number and exit.

Example:

    $ mp3cat out.mp3 1.mp3 2.mp3 3.mp3


## Installation

If you have Go installed you can run:

    $ go get github.com/dmulholland/mp3cat/mp3cat

This will download, compile, and install the latest version of the application to your `$GOPATH/bin` directory.

Alternatively, you can download a precompiled binary from the [releases](https://github.com/dmulholland/mp3cat/releases) page.


## License

This work has been placed in the public domain.

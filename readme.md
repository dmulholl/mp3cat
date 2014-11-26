
mp3cat
======

Command line utility for joining MP3 files without re-encoding. Supports both constant
bit rate (CBR) and variable bit rate (VBR) files.


### Usage

    Usage: mp3cat [FLAGS] ARGUMENTS

    Concatenates MP3 files without re-encoding. Supports both CBR and VBR files.
    Strips ID3 tags and garbage data from the output.

    Arguments:

      <outfile>        Output file.
      <infiles>        List of input files to concatenate.

    Flags:

      -f, --force      Force overwriting of existing output files.
      -v, --verbose    Report progress.
      --help           Display this help text and exit.
      --version        Display version number and exit.

Example:

    $ mp3cat output.mp3 input-1.mp3 input-2.mp3 input-3.mp3


### Installation

    $ go get github.com/dmulholland/mp3cat/mp3cat


### License

This work has been placed in the public domain.

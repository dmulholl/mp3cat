
# MP3Cat

MP3Cat is a fast command-line utility for merging MP3 files without re-encoding. It supports both constant bit rate (CBR) and variable bit rate (VBR) MP3 files.

    $ mp3cat --help

    Usage: mp3cat [FLAGS] [OPTIONS] ARGUMENTS

      Concatenates MP3 files without re-encoding. Supports both constant bit
      rate (CBR) and variable bit rate (VBR) files. Strips ID3 tags and
      garbage data from the output.

    Arguments:
      <files>           List of input files to merge.

    Options:
      -o, --out <file>  Output filename. Defaults to 'output.mp3'.

    Flags:
      -f, --force       Overwrite an existing output file.
          --help        Display this help text and exit.
      -t, --tag         Copy the ID3 tag from the first input file.
      -v, --verbose     Report progress.
          --version     Display the application's version number and exit.

See the project's [homepage][] for details and binary downloads.

[homepage]: http://mulholland.xyz/dev/mp3cat/

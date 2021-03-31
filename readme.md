# MP3Cat

[documentation]: http://www.dmulholl.com/dev/mp3cat.html
[releases]: https://github.com/dmulholl/mp3cat/releases
[mp3binder]: https://github.com/crra/mp3binder


MP3Cat is a simple command line utility for concatenating MP3 files without re-encoding.

<p align="center">
    <img src="mp3cat.png" width="600px">
</p>

* [Documentation][]
* [Pre-built Binaries][releases]

Run `mp3cat --help` to view the command line help:

    Usage: mp3cat [files]

      This tool concatenates MP3 files without re-encoding.

    Arguments:
      [files]                 List of files to merge.

    Options:
      -d, --dir <path>        Directory of files to merge.
      -o, --out <path>        Output filepath.

    Flags:
      -f, --force             Overwrite an existing output file.
      -h, --help              Display this help text and exit.
      -q, --quiet             Run in quiet mode.
      -v, --version           Display the version number.

Files to be joined can be specified as a list of filenames:

    $ mp3cat one.mp3 two.mp3 three.mp3

Alternatively, an entire directory of `.mp3` files can be concatenated:

    $ mp3cat --dir /path/to/directory

Note that this application is in *maintenance mode* &mdash; it's intentionally simple and I'm not planning to add any extra features.
If you need something more complex check out [mp3binder][] which is an active fork with additional functionality.

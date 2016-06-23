
# Go Version

Install:

    go get github.com/dmulholland/clio/go/clio

Import:

    import "github.com/dmulholland/clio/go/clio"


## Usage

Initialize an argument parser, optionally specifying help text and a version
string:

    func NewParser(helptext string, version string) *ArgParser

Supplying help text activates the automatic `--help` flag; supplying a version string activates the automatic `--version` flag. An empty string `""` can be passed for either parameter.

You can now register your application's options and commands on the parser instance as explained below.

Once the required options and commands have been registered, call the parser's `Parse()` method to process the application's command line arguments.

    func (parser *ArgParser) Parse()

Parsed option values can be retrieved from the parser instance itself.


### Options

Clio supports long-form options (`--foo`) with single-character aliases (`-f`). Note that when registering an option you should omit the leading dashes, i.e. you should register the option name as `"foo"` rather than `"--foo"`.

Registering options:

*   `func (parser *ArgParser) AddFlag(name string, alias ...rune)`

    Register a flag, optionally specifying a single-character alias. A flag is
    a boolean option - it takes no argument but is either present (true) or
    absent (false). The alias parameter may be omitted.

*   `func (parser *ArgParser) AddStrOpt(name string, defVal string, alias ...rune)`

    Register a string option and its default value, optionally specifying a
    single-character alias. The alias parameter may be omitted.

*   `func (parser *ArgParser) AddIntOpt(name string, defVal int, alias ...rune)`

    Register an integer option and its default value, optionally specifying a
    single-character alias. The alias parameter may be omitted.

*   `func (parser *ArgParser) AddFloatOpt(name string, defVal float64, alias ...rune)`

    Register a floating-point option and its default value, optionally
    specifying a single-character alias. The alias parameter may be omitted.

Retrieving values:

*   `func (parser *ArgParser) GetFlag(name string) bool`

*   `func (parser *ArgParser) GetStrOpt(name string) string`

*   `func (parser *ArgParser) GetIntOpt(name string) int`

*   `func (parser *ArgParser) GetFloatOpt(name string) float64`

All options have default values which are used when the option is omitted from the command line arguments.

Note that Clio supports the standard `--` option-parsing switch. All command line arguments following a `--` will be treated as positional arguments rather than options, even if they begin with a single or double dash.


### Positional Arguments

The following methods provide access to positional arguments:

*   `func (parser *ArgParser) HasArgs() bool`

    Returns true if at least one positional argument has been found.

*   `func (parser *ArgParser) NumArgs() int`

    Returns the number of positional arguments.

*   `func (parser *ArgParser) GetArg(index int) string`

    Returns the positional argument at the specified index.

*   `func (parser *ArgParser) GetArgs() []string`

    Returns the positional arguments as a slice of strings.

*   `func (parser *ArgParser) GetArgsAsInts() []int`

    Attempts to parse and return the positional arguments as a slice of
    integers. Exits with an error message on failure.

*   `func (parser *ArgParser) GetArgsAsFloats() []float64`

    Attempts to parse and return the positional arguments as a slice of floats.
    Exits with an error message on failure.


### Commands

Clio supports git-style command interfaces with arbitrarily-nested commands. Register a command on a parser using the `AddCmd()` method:

    func (parser *ArgParser) AddCmd(command string, callback Callback, helptext string) *ArgParser

This method returns the `ArgParser` instance associated with the new command. You can register flags and options on this sub-parser using the methods listed above. (Note that you do not need to call `Parse()` on the command parser instance - calling `Parse()` on the root parser is sufficient.)

Commands support an automatic `--help` flag and an automatic `help <cmd>` command.

The supplied callback function will be called if the command is found. This callback should accept the command's sub-parser instance as its sole argument.

Other command-related methods are:

*   `func (parser *ArgParser) HasCmd() bool`

    Returns true if the parser has identified a command.

*   `func (parser *ArgParser) GetCmd() string`

    Returns the command name, if a command was identified.

*   `func (parser *ArgParser) GetCmdParser() *ArgParser`

    Returns the command parser, if a command was identified.

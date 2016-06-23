/*
    Package clio is a minimalist argument-parsing library for creating elegant
    command-line interfaces.
*/
package clio


import (
    "fmt"
    "os"
    "strings"
    "strconv"
    "unicode"
    "sort"
)


// Package version number.
const Version = "1.0.0"


// Enum for classifying option types. We use 'flag' as a synonym for boolean
// options, i.e. options that are either present (true) or absent (false). All
// other option types require an argument.
const (
    flagType = iota
    strType
    intType
    floatType
)


// An option can have a boolean, string, integer, or floating point value.
type option struct {
    optType int
    boolVal bool
    strVal string
    intVal int
    floatVal float64
}


// String returns a string representation of the option's value.
func (opt *option) String() string {
    var str string
    switch opt.optType {
    case flagType:
        str = fmt.Sprintf("%v", opt.boolVal)
    case strType:
        str = opt.strVal
    case intType:
        str = fmt.Sprintf("%v", opt.intVal)
    case floatType:
        str = fmt.Sprintf("%v", opt.floatVal)
    }
    return str
}


// Callback function for processing commands.
type Callback func(*ArgParser)


// Makes a slice of string arguments available as a stream.
type argStream struct {
    args []string
    index int
    length int
}


// Initializes a new argStream instance.
func newArgStream(args []string) *argStream {
    return &argStream{
        args: args,
        index: 0,
        length: len(args),
    }
}


// Returns true if the stream contains at least one more argument.
func (stream *argStream) hasNext() bool {
    return stream.index < stream.length
}


// Returns the next argument from the stream.
func (stream *argStream) next() string {
    stream.index += 1
    return stream.args[stream.index - 1]
}


// Returns the next argument from the stream without consuming it.
func (stream *argStream) peek() string {
    return stream.args[stream.index]
}


// Returns a slice containing all the remaining arguments from the stream.
func (stream *argStream) remainder() []string {
    return stream.args[stream.index:]
}


// An ArgParser instance stores registered options and parsed command line
// arguments.
//
// Note that every registered command recursively receives an ArgParser instance
// of its own. In theory commands can be stacked to any depth, although in
// practice even two levels is confusing for users and best avoided.
type ArgParser struct {

    // Help text for the application or command.
    helptext string

    // Application version number.
    version string

    // Stores option objects indexed by option name.
    options map[string]*option

    // Stores option objects indexed by single-letter shortcut.
    shortcuts map[rune]*option

    // Stores command sub-parser instances indexed by command.
    commands map[string]*ArgParser

    // Stores command callbacks indexed by command.
    callbacks map[string]Callback

    // Stores positional arguments parsed from the input array.
    arguments []string

    // Stores the command string, if a command is found.
    command string

    // Stores the command's parser instance, if a command is found.
    commandParser *ArgParser
}


// NewParser initializes a new ArgParser instance.
func NewParser(helptext string, version string) *ArgParser {
    return &ArgParser {
        helptext: strings.TrimSpace(helptext),
        version: strings.TrimSpace(version),
        options: make(map[string]*option),
        shortcuts: make(map[rune]*option),
        commands: make(map[string]*ArgParser),
        callbacks: make(map[string]Callback),
        arguments: make([]string, 0, 10),
    }
}


// AddFlag registers a flag (a boolean option) on a parser instance.
// The caller can optionally specify a single-letter shortcut alias.
func (parser *ArgParser) AddFlag(name string, alias ...rune) {
    opt := option{
        optType: flagType,
        boolVal: false,
    }
    parser.options[name] = &opt
    for _, c := range alias {
        parser.shortcuts[c] = &opt
    }
}


// AddStrOpt registers a string option on a parser instance.
// The caller can optionally specify a single-letter shortcut alias.
func (parser *ArgParser) AddStrOpt(name string, defVal string, alias ...rune) {
    opt := option{
        optType: strType,
        strVal: defVal,
    }
    parser.options[name] = &opt
    for _, c := range alias {
        parser.shortcuts[c] = &opt
    }
}


// AddIntOpt registers an integer option on a parser instance.
// The caller can optionally specify a single-letter shortcut alias.
func (parser *ArgParser) AddIntOpt(name string, defVal int, alias ...rune) {
    opt := option{
        optType: intType,
        intVal: defVal,
    }
    parser.options[name] = &opt
    for _, c := range alias {
        parser.shortcuts[c] = &opt
    }
}


// AddFloatOpt registers a float option on a parser instance.
// The caller can optionally specify a single-letter shortcut alias.
func (parser *ArgParser) AddFloatOpt(name string, defVal float64, alias ...rune) {
    opt := option{
        optType: floatType,
        floatVal: defVal,
    }
    parser.options[name] = &opt
    for _, c := range alias {
        parser.shortcuts[c] = &opt
    }
}


// AddCmd registers a command on a parser instance.
func (parser *ArgParser) AddCmd(command string, callback Callback, helptext string) *ArgParser {
    cmdParser := NewParser(helptext, "")
    parser.commands[command] = cmdParser
    parser.callbacks[command] = callback
    return cmdParser
}


// Help prints the parser's help text, then exits.
func (parser *ArgParser) Help() {
    fmt.Println(parser.helptext)
    os.Exit(0)
}


// ParseArgs parses the specified slice of string arguments.
func (parser *ArgParser) ParseArgs(args []string) {

    // Switch to turn off parsing if we encounter a -- argument.
    // Everything following the -- will be treated as a positional argument.
    parsing := true

    // Convert the input slice into a stream.
    stream := newArgStream(args)

    // Loop while we have arguments to process.
    for stream.hasNext() {

        // Fetch the next argument from the stream.
        arg := stream.next()

        // If parsing has been turned off, simply add the argument to the
        // list of positionals.
        if !parsing {
            parser.arguments = append(parser.arguments, arg)
            continue
        }

        // If we encounter a -- argument, turn off parsing.
        if arg == "--" {
            parsing = false
            continue
        }

        // Is the argument a long-form option or flag?
        if strings.HasPrefix(arg, "--") {

            // Strip the -- prefix.
            arg = arg[2:]

            // Is the argument a registered option name?
            if opt, ok := parser.options[arg]; ok {

                // If the option is a flag, store the boolean true.
                if opt.optType == flagType {
                    opt.boolVal = true
                    continue
                }

                // Not a flag, so check for a following argument.
                if !stream.hasNext() {
                    fmt.Fprintf(os.Stderr, "Error: missing argument for the --%v option.\n", arg)
                    os.Exit(1)
                }

                // Fetch the argument from the stream and attempt to parse it.
                nextarg := stream.next()

                switch opt.optType {

                case strType:
                    opt.strVal = nextarg

                case intType:
                    intVal, err := strconv.ParseInt(nextarg, 0, 0)
                    if err != nil {
                        fmt.Fprintf(os.Stderr, "Error: cannot parse '%v' as an integer.\n", nextarg)
                        os.Exit(1)
                    }
                    opt.intVal = int(intVal)

                case floatType:
                    floatVal, err := strconv.ParseFloat(nextarg, 64)
                    if err != nil {
                        fmt.Fprintf(os.Stderr, "Error: cannot parse '%v' as a float.\n", nextarg)
                        os.Exit(1)
                    }
                    opt.floatVal = floatVal
                }

                // We have successfully parsed a long-form option with an
                // argument. Move on to the next argument in the stream.
                continue
            }

            // Is the argument the automatic --help command?
            if arg == "help" && parser.helptext != "" {
                fmt.Println(parser.helptext)
                os.Exit(0)
            }

            // Is the argument the automatic --version command.
            if arg == "version" && parser.version != "" {
                fmt.Println(parser.version)
                os.Exit(0)
            }

            // The argument is not a registered or automatic option.
            // Print an error message and exit.
            fmt.Fprintf(os.Stderr, "Error: --%v is not a recognised option.\n", arg)
            os.Exit(1)
        }

        // Is the argument a short-form option or flag?
        if strings.HasPrefix(arg, "-"){

            // If the argument consists of a sigle dash or a dash followed by
            // a digit, treat it as a positional argument.
            if arg == "-" || unicode.IsDigit([]rune(arg)[1]) {
                parser.arguments = append(parser.arguments, arg)
                continue
            }

            // Examine each character individually to allow for condensed
            // short-form arguments, i.e.
            //     -a -b foo -c bar
            // is equivalent to:
            //     -abc foo bar
            for _, c := range arg[1:] {

                // Is the character a registered shortcut?
                if opt, ok := parser.shortcuts[c]; ok {

                    // If the option is a flag, store the boolean true.
                    if opt.optType == flagType {
                        opt.boolVal = true
                        continue
                    }

                    // Not a flag, so check for a following argument.
                    if !stream.hasNext() {
                        fmt.Fprintf(os.Stderr, "Error: missing argument for the -%c option.\n", c)
                        os.Exit(1)
                    }

                    // Fetch the argument from the stream and attempt to parse it.
                    nextarg := stream.next()

                    switch opt.optType {

                    case strType:
                        opt.strVal = nextarg

                    case intType:
                        intVal, err := strconv.ParseInt(nextarg, 0, 0)
                        if err != nil {
                            fmt.Fprintf(os.Stderr, "Error: cannot parse '%v' as an integer.\n", nextarg)
                            os.Exit(1)
                        }
                        opt.intVal = int(intVal)

                    case floatType:
                        floatVal, err := strconv.ParseFloat(nextarg, 64)
                        if err != nil {
                            fmt.Fprintf(os.Stderr, "Error: cannot parse '%v' as a float.\n", nextarg)
                            os.Exit(1)
                        }
                        opt.floatVal = floatVal
                    }

                    // We have successfully parsed a single short-form option.
                    // Move on to the next short-form option in the block.
                    continue
                }

                // Not a registered shortcut. Print an error and exit.
                fmt.Fprintf(os.Stderr, "Error: -%c is not a recognised option.\n", c)
                os.Exit(1)
            }

            // We have successfully parsed a block of short-form options.
            // Move on to the next argument in the stream.
            continue
        }

        // Is the argument a registered command?
        if cmdParser, ok := parser.commands[arg]; ok {
            cmdParser.ParseArgs(stream.remainder())
            parser.callbacks[arg](cmdParser)
            parser.command = arg
            parser.commandParser = cmdParser
            break
        }

        // Is the argument the automatic 'help' command?
        if arg == "help"{
            if stream.hasNext() {
                command := stream.next()
                if cmdParser, ok := parser.commands[command]; ok {
                    fmt.Println(cmdParser.helptext)
                    os.Exit(0)
                } else {
                    fmt.Fprintf(os.Stderr, "Error: '%v' is not a recognised command.\n", command)
                    os.Exit(1)
                }
            } else {
                fmt.Fprintf(os.Stderr, "Error: the help command requires an argument.\n")
                os.Exit(1)
            }
        }

        // If we get here, we have a positional argument.
        parser.arguments = append(parser.arguments, arg)
    }
}


// Parse parses the application's command line arguments.
func (parser *ArgParser) Parse() {
    parser.ParseArgs(os.Args[1:])
}


// GetFlag returns true if the named flag was found.
func (parser *ArgParser) GetFlag(name string) bool {
    return parser.options[name].boolVal
}


// GetStrOpt returns the value of the named option.
func (parser *ArgParser) GetStrOpt(name string) string {
    return parser.options[name].strVal
}


// GetIntOpt returns the value of the named option.
func (parser *ArgParser) GetIntOpt(name string) int {
    return parser.options[name].intVal
}


// GetFloatOpt returns the value of the named option.
func (parser *ArgParser) GetFloatOpt(name string) float64 {
    return parser.options[name].floatVal
}


// HasArgs returns true if the parser has identified one or more positional
// arguments.
func (parser *ArgParser) HasArgs() bool {
    return len(parser.arguments) > 0
}


// NumArgs returns the number of positional arguments.
func (parser *ArgParser) NumArgs() int {
    return len(parser.arguments)
}


// GetArg returns the positional argument at the specified index.
func (parser *ArgParser) GetArg(index int) string {
    return parser.arguments[index]
}


// GetArgs returns the parser's positional arguments as a slice of strings.
func (parser *ArgParser) GetArgs() []string {
    return parser.arguments
}


// GetArgsAsInts attempts to parse and return the positional arguments as a
// slice of integers. The application will exit with an error message if any
// of the arguments cannot be parsed as an integer.
func (parser *ArgParser) GetArgsAsInts() []int {
    intList := make([]int, 0, 10)
    for _, strArg := range parser.arguments {
        intArg, err := strconv.ParseInt(strArg, 0, 0)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: cannot parse '%v' as an integer.\n", strArg)
            os.Exit(1)
        }
        intList = append(intList, int(intArg))
    }
    return intList
}


// GetArgsAsFloats attempts to parse and return the positional arguments as a
// slice of floats. The application will exit with an error message if any
// of the arguments cannot be parsed as a float.
func (parser *ArgParser) GetArgsAsFloats() []float64 {
    floatList := make([]float64, 0, 10)
    for _, strArg := range parser.arguments {
        floatArg, err := strconv.ParseFloat(strArg, 64)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: cannot parse '%v' as a float.\n", strArg)
            os.Exit(1)
        }
        floatList = append(floatList, floatArg)
    }
    return floatList
}


// HasCmd returns true if the parser has identified a command.
func (parser *ArgParser) HasCmd() bool {
    return parser.command != ""
}


// GetCmd returns the command string, if a command was found.
func (parser *ArgParser) GetCmd() string {
    return parser.command
}


// GetCmdParser returns the command's parser instance, if a command was found.
func (parser *ArgParser) GetCmdParser() *ArgParser {
    return parser.commandParser
}


// String returns a string representation of the parser instance.
func (parser *ArgParser) String() string {
    lines := make([]string, 0, 10)

    lines = append(lines, "Options:")
    if len(parser.options) > 0 {
        names := make([]string, 0, len(parser.options))
        for name := range parser.options {
            names = append(names, name)
        }
        sort.Strings(names)
        for _, name := range names {
            lines = append(lines, fmt.Sprintf("  %v: %v", name, parser.options[name]))
        }
    } else {
        lines = append(lines, "  [none]")
    }

    lines = append(lines, "\nArguments:")
    if len(parser.arguments) > 0 {
        for _, arg := range parser.arguments {
            lines = append(lines, fmt.Sprintf("  %v", arg))
        }
    } else {
        lines = append(lines, "  [none]")
    }

    lines = append(lines, "\nCommand:")
    if parser.HasCmd() {
        lines = append(lines, fmt.Sprintf("  %v", parser.GetCmd()))
    } else {
        lines = append(lines, "  [none]")
    }

    return strings.Join(lines, "\n")
}

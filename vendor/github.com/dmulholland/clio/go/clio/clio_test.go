package clio


import (
    "testing"
)


// -------------------------------------------------------------------------
// Boolean options.
// -------------------------------------------------------------------------


func TestBoolOptionEmpty(t *testing.T) {
    parser := NewParser("", "")
    parser.AddFlag("bool")
    parser.ParseArgs([]string{})
    if parser.GetFlag("bool") != false {
        t.Fail()
    }
}


func TestBoolOptionMissing(t *testing.T) {
    parser := NewParser("", "")
    parser.AddFlag("bool")
    parser.ParseArgs([]string{"foo", "bar"})
    if parser.GetFlag("bool") != false {
        t.Fail()
    }
}


func TestBoolOptionLongform(t *testing.T) {
    parser := NewParser("", "")
    parser.AddFlag("bool")
    parser.ParseArgs([]string{"--bool"})
    if parser.GetFlag("bool") != true {
        t.Fail()
    }
}


func TestBoolOptionShortform(t *testing.T) {
    parser := NewParser("", "")
    parser.AddFlag("bool", 'b')
    parser.ParseArgs([]string{"-b"})
    if parser.GetFlag("bool") != true {
        t.Fail()
    }
}


// -------------------------------------------------------------------------
// String options.
// -------------------------------------------------------------------------


func TestStringOptionEmpty(t *testing.T) {
    parser := NewParser("", "")
    parser.AddStrOpt("string", "default")
    parser.ParseArgs([]string{})
    if parser.GetStrOpt("string") != "default" {
        t.Fail()
    }
}


func TestStringOptionMissing(t *testing.T) {
    parser := NewParser("", "")
    parser.AddStrOpt("string", "default")
    parser.ParseArgs([]string{"foo", "bar"})
    if parser.GetStrOpt("string") != "default" {
        t.Fail()
    }
}


func TestStringOptionLongform(t *testing.T) {
    parser := NewParser("", "")
    parser.AddStrOpt("string", "default")
    parser.ParseArgs([]string{"--string", "value"})
    if parser.GetStrOpt("string") != "value" {
        t.Fail()
    }
}


func TestStringOptionShortform(t *testing.T) {
    parser := NewParser("", "")
    parser.AddStrOpt("string", "default", 's')
    parser.ParseArgs([]string{"-s", "value"})
    if parser.GetStrOpt("string") != "value" {
        t.Fail()
    }
}


// -------------------------------------------------------------------------
// Integer options.
// -------------------------------------------------------------------------


func TestIntOptionEmpty(t *testing.T) {
    parser := NewParser("", "")
    parser.AddIntOpt("int", 101)
    parser.ParseArgs([]string{})
    if parser.GetIntOpt("int") != 101 {
        t.Fail()
    }
}


func TestIntOptionMissing(t *testing.T) {
    parser := NewParser("", "")
    parser.AddIntOpt("int", 101)
    parser.ParseArgs([]string{"foo", "bar"})
    if parser.GetIntOpt("int") != 101 {
        t.Fail()
    }
}


func TestIntOptionLongform(t *testing.T) {
    parser := NewParser("", "")
    parser.AddIntOpt("int", 101)
    parser.ParseArgs([]string{"--int", "202"})
    if parser.GetIntOpt("int") != 202 {
        t.Fail()
    }
}


func TestIntOptionShortform(t *testing.T) {
    parser := NewParser("", "")
    parser.AddIntOpt("int", 101, 'i')
    parser.ParseArgs([]string{"-i", "202"})
    if parser.GetIntOpt("int") != 202 {
        t.Fail()
    }
}


func TestIntOptionNegative(t *testing.T) {
    parser := NewParser("", "")
    parser.AddIntOpt("int", 101)
    parser.ParseArgs([]string{"--int", "-202"})
    if parser.GetIntOpt("int") != -202 {
        t.Fail()
    }
}


// -------------------------------------------------------------------------
// Float options.
// -------------------------------------------------------------------------


func TestFloatOptionEmpty(t *testing.T) {
    parser := NewParser("", "")
    parser.AddFloatOpt("float", 1.1)
    parser.ParseArgs([]string{})
    if parser.GetFloatOpt("float") != 1.1 {
        t.Fail()
    }
}


func TestFloatOptionMissing(t *testing.T) {
    parser := NewParser("", "")
    parser.AddFloatOpt("float", 1.1)
    parser.ParseArgs([]string{"foo", "bar"})
    if parser.GetFloatOpt("float") != 1.1 {
        t.Fail()
    }
}


func TestFloatOptionLongform(t *testing.T) {
    parser := NewParser("", "")
    parser.AddFloatOpt("float", 1.1)
    parser.ParseArgs([]string{"--float", "2.2"})
    if parser.GetFloatOpt("float") != 2.2 {
        t.Fail()
    }
}


func TestFloatOptionShortform(t *testing.T) {
    parser := NewParser("", "")
    parser.AddFloatOpt("float", 1.1, 'f')
    parser.ParseArgs([]string{"-f", "2.2"})
    if parser.GetFloatOpt("float") != 2.2 {
        t.Fail()
    }
}


func TestFloatOptionNegative(t *testing.T) {
    parser := NewParser("", "")
    parser.AddFloatOpt("float", 1.1)
    parser.ParseArgs([]string{"--float", "-2.2"})
    if parser.GetFloatOpt("float") != -2.2 {
        t.Fail()
    }
}


// -------------------------------------------------------------------------
// Multiple options.
// -------------------------------------------------------------------------


func TestMultiOptionsEmpty(t *testing.T) {
    parser := NewParser("", "")
    parser.AddFlag("bool1")
    parser.AddFlag("bool2", 'b')
    parser.AddStrOpt("string1", "default1")
    parser.AddStrOpt("string2", "default2", 's')
    parser.AddIntOpt("int1", 101)
    parser.AddIntOpt("int2", 202, 'i')
    parser.AddFloatOpt("float1", 1.1)
    parser.AddFloatOpt("float2", 2.2, 'f')
    parser.ParseArgs([]string{})
    if parser.GetFlag("bool1") != false {
        t.Fail()
    }
    if parser.GetFlag("bool2") != false {
        t.Fail()
    }
    if parser.GetStrOpt("string1") != "default1" {
        t.Fail()
    }
    if parser.GetStrOpt("string2") != "default2" {
        t.Fail()
    }
    if parser.GetIntOpt("int1") != 101 {
        t.Fail()
    }
    if parser.GetIntOpt("int2") != 202 {
        t.Fail()
    }
    if parser.GetFloatOpt("float1") != 1.1 {
        t.Fail()
    }
    if parser.GetFloatOpt("float2") != 2.2 {
        t.Fail()
    }
}


func TestMultiOptionsLongform(t *testing.T) {
    parser := NewParser("", "")
    parser.AddFlag("bool1")
    parser.AddFlag("bool2", 'b')
    parser.AddStrOpt("string1", "default1")
    parser.AddStrOpt("string2", "default2", 's')
    parser.AddIntOpt("int1", 101)
    parser.AddIntOpt("int2", 202, 'i')
    parser.AddFloatOpt("float1", 1.1)
    parser.AddFloatOpt("float2", 2.2, 'f')
    parser.ParseArgs([]string{
        "--bool1",
        "--bool2",
        "--string1", "value1",
        "--string2", "value2",
        "--int1", "303",
        "--int2", "404",
        "--float1", "3.3",
        "--float2", "4.4",
    })
    if parser.GetFlag("bool1") != true {
        t.Fail()
    }
    if parser.GetFlag("bool2") != true {
        t.Fail()
    }
    if parser.GetStrOpt("string1") != "value1" {
        t.Fail()
    }
    if parser.GetStrOpt("string2") != "value2" {
        t.Fail()
    }
    if parser.GetIntOpt("int1") != 303 {
        t.Fail()
    }
    if parser.GetIntOpt("int2") != 404 {
        t.Fail()
    }
    if parser.GetFloatOpt("float1") != 3.3 {
        t.Fail()
    }
    if parser.GetFloatOpt("float2") != 4.4 {
        t.Fail()
    }
}


func TestMultiOptionsShortform(t *testing.T) {
    parser := NewParser("", "")
    parser.AddFlag("bool1")
    parser.AddFlag("bool2", 'b')
    parser.AddStrOpt("string1", "default1")
    parser.AddStrOpt("string2", "default2", 's')
    parser.AddIntOpt("int1", 101)
    parser.AddIntOpt("int2", 202, 'i')
    parser.AddFloatOpt("float1", 1.1)
    parser.AddFloatOpt("float2", 2.2, 'f')
    parser.ParseArgs([]string{
        "--bool1",
        "-b",
        "--string1", "value1",
        "-s", "value2",
        "--int1", "303",
        "-i", "404",
        "--float1", "3.3",
        "-f", "4.4",
    })
    if parser.GetFlag("bool1") != true {
        t.Fail()
    }
    if parser.GetFlag("bool2") != true {
        t.Fail()
    }
    if parser.GetStrOpt("string1") != "value1" {
        t.Fail()
    }
    if parser.GetStrOpt("string2") != "value2" {
        t.Fail()
    }
    if parser.GetIntOpt("int1") != 303 {
        t.Fail()
    }
    if parser.GetIntOpt("int2") != 404 {
        t.Fail()
    }
    if parser.GetFloatOpt("float1") != 3.3 {
        t.Fail()
    }
    if parser.GetFloatOpt("float2") != 4.4 {
        t.Fail()
    }
}


// -------------------------------------------------------------------------
// Condensed short-form options.
// -------------------------------------------------------------------------


func TestCondensedOptions(t *testing.T) {
    parser := NewParser("", "")
    parser.AddFlag("bool", 'b')
    parser.AddStrOpt("string", "default", 's')
    parser.AddIntOpt("int", 101, 'i')
    parser.AddFloatOpt("float", 1.1, 'f')
    parser.ParseArgs([]string{"-bsif", "value", "202", "2.2"})
    if parser.GetFlag("bool") != true {
        t.Fail()
    }
    if parser.GetStrOpt("string") != "value" {
        t.Fail()
    }
    if parser.GetIntOpt("int") != 202 {
        t.Fail()
    }
    if parser.GetFloatOpt("float") != 2.2 {
        t.Fail()
    }
}


// -------------------------------------------------------------------------
// Positional arguments.
// -------------------------------------------------------------------------


func TestPositionalArgsEmpty(t *testing.T) {
    parser := NewParser("", "")
    parser.ParseArgs([]string{})
    if parser.HasArgs() != false {
        t.Fail()
    }
}


func TestPositionalArgs(t *testing.T) {
    parser := NewParser("", "")
    parser.ParseArgs([]string{"foo", "bar"})
    if parser.HasArgs() != true {
        t.Fail()
    }
    if parser.NumArgs() != 2 {
        t.Fail()
    }
    if parser.GetArg(0) != "foo" {
        t.Fail()
    }
    if parser.GetArg(1) != "bar" {
        t.Fail()
    }
    if parser.GetArgs()[0] != "foo" {
        t.Fail()
    }
    if parser.GetArgs()[1] != "bar" {
        t.Fail()
    }
}


func TestPositionalArgsAsInts(t *testing.T) {
    parser := NewParser("", "")
    parser.ParseArgs([]string{"1", "11"})
    if parser.GetArgsAsInts()[0] != 1 {
        t.Fail()
    }
    if parser.GetArgsAsInts()[1] != 11 {
        t.Fail()
    }
}


func TestPositionalArgsAsFloats(t *testing.T) {
    parser := NewParser("", "")
    parser.ParseArgs([]string{"1.1", "11.1"})
    if parser.GetArgsAsFloats()[0] != 1.1 {
        t.Fail()
    }
    if parser.GetArgsAsFloats()[1] != 11.1 {
        t.Fail()
    }
}


// -------------------------------------------------------------------------
// Commands
// -------------------------------------------------------------------------


func callback(parser *ArgParser) {}


func TestCommandAbsent(t *testing.T) {
    parser := NewParser("", "")
    parser.AddCmd("cmd", callback, "helptext")
    parser.ParseArgs([]string{})
    if parser.HasCmd() != false {
        t.Fail()
    }
}


func TestCommandPresent(t *testing.T) {
    parser := NewParser("", "")
    cmdParser := parser.AddCmd("cmd", callback, "helptext")
    parser.ParseArgs([]string{"cmd"})
    if parser.HasCmd() != true {
        t.Fail()
    }
    if parser.GetCmd() != "cmd" {
        t.Fail()
    }
    if parser.GetCmdParser() != cmdParser {
        t.Fail()
    }
}


func TestCommandWithOptions(t *testing.T) {
    parser := NewParser("", "")
    cmdParser := parser.AddCmd("cmd", callback, "helptext")
    cmdParser.AddFlag("bool")
    cmdParser.AddStrOpt("string", "default")
    cmdParser.AddIntOpt("int", 101)
    cmdParser.AddFloatOpt("float", 1.1)
    parser.ParseArgs([]string{
        "cmd",
        "foo", "bar",
        "--string", "value",
        "--int", "202",
        "--float", "2.2",
    })
    if parser.HasCmd() != true {
        t.Fail()
    }
    if parser.GetCmd() != "cmd" {
        t.Fail()
    }
    if parser.GetCmdParser() != cmdParser {
        t.Fail()
    }
    if cmdParser.HasArgs() != true {
        t.Fail()
    }
    if cmdParser.NumArgs() != 2 {
        t.Fail()
    }
    if cmdParser.GetStrOpt("string") != "value" {
        t.Fail()
    }
    if cmdParser.GetIntOpt("int") != 202 {
        t.Fail()
    }
    if cmdParser.GetFloatOpt("float") != 2.2 {
        t.Fail()
    }
}

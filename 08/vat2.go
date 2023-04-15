package main

import (
	"bufio"
	"bytes"
	"fmt"

	"path/filepath"
	"strconv"
	"strings"

	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	local = iota + 1
	argument
	this
	that

	functionCallPrep = "(functionCallPrep)\n// push return address \n@R13\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n// push lcl\n@1\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n// push args\n@2\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n// push this\n@3\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n// push that\n@4\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n//+ARG = (*sp) - 5 - arg num\n@5\nD=A\n@R14\nD=M+D\n@SP\nD=M-D\n@ARG\nM=D\n@SP\nD=M\n@LCL\nM=D\n// R15 stores the call address \n@R15\nA=M\n0;JMP\n"
	returnPrep       = "(return-prep)\n// pop last in stack to arg\n@SP\nAM=M-1\nD=M\n@ARG\nA=M\nM=D\nD=A\n@R13\nM=D\n// pop that\n@LCL\nD=M\n@SP\nAM=D\nD=M\n@THAT\nM=D\n// pop this\n@SP\nAM=M-1\nD=M\n@THIS\nM=D\n// pop to arg\n@SP\nAM=M-1\nD=M\n@ARG\nM=D\n// pop to lcl\n@SP\nAM=M-1\nD=M\n@LCL\nM=D\n// return address\n@SP\nAM=M-1\nD=M\n@R15\nM=D\n@R13\nD=M\n@SP\nM=D\n@R15\nA=M\n0;JMP\n"
)

var (
	app = cli.NewApp()
	seg = map[string]int{
		"local":    local,
		"argument": argument,
		"this":     this,
		"that":     that,
	}
	ops = map[string]string{
		"add": "+",
		"sub": "-",
		"and": "&",
		"or":  "|",

		"not": "!",
		"neg": "-",
	}
)

type vm struct {
	reader   *bufio.Scanner
	path     string
	currfunc string
	curr     []string
	skip     int
	count    int
}

func new(file *os.File, path string) *vm {
	v := &vm{
		reader: bufio.NewScanner(file),
		path:   path,
	}
	v.next()
	return v
}

func (v *vm) next() {
	token := v.reader.Text()
	val := strings.TrimSpace(string(token))
	if com := strings.Index(val, "//"); com != -1 {
		val = val[:com]
	}
	v.curr = strings.Split(val, " ")
}
func (v *vm) skipEmpty() {
	for len(v.curr) == 0 {
		v.next()
	}
}
func (v *vm) Compare() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nD=M-D\nM=-1\n@skip$%d$\nD;", v.skip))
	switch v.curr[0] {
	case "lt":
		buffer.WriteString("JLT\n")
	case "gt":
		buffer.WriteString("JGT\n")
	case "eq":
		buffer.WriteString("JEQ\n")
	}
	buffer.WriteString(fmt.Sprintf("@SP\nA=M\nM=0\n(skip$%d$)\n@SP\nM=M+1\n", v.skip))
	v.skip++
	return buffer.String()
}
func (v *vm) Pop() string {
	var buffer string
	switch v.curr[1] {
	case "local", "argument", "this", "that":
		buffer = fmt.Sprintf("@%s\nD=A\n@%d\nD=M+D\n@R13\nM=D\n@SP\nAM=M-1\nD=M\n@R13\nA=M\nM=D\n", v.curr[2], seg[v.curr[1]])
	default:
		buffer = v.helper("@SP\nAM=M-1\nD=M\n@%s\nM=D\n")
	}
	return buffer
}
func (v *vm) Push() string {
	switch v.curr[1] {
	case "constant":
		return fmt.Sprintf("@%s\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n", v.curr[2])
	case "local", "argument", "this", "that":
		return fmt.Sprintf("@%s\nD=A\n@%d\nA=M+D\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n", v.curr[2], seg[v.curr[1]])
	default:
		return v.helper("@%s\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")
	}
}
func (v *vm) helper(dummy string) string {
	var buffer string
	switch v.curr[1] {
	case "static": //> 16 - 255
		buffer = filepath.Base(v.path) + "$" + v.curr[2]
	case "pointer": //> pointer 0/1
		if v.curr[2] == "0" {
			buffer = "3"
		} else if v.curr[2] == "1" {
			buffer = "4"
		}
	case "temp": //> 5 - 12
		num, _ := strconv.Atoi(v.curr[2])
		buffer = strconv.Itoa(5 + num)
	default:
		buffer = v.curr[2]
	}
	return fmt.Sprintf(dummy, buffer)
}
func (v *vm) Label() string {
	if len(v.currfunc) != 0 {
		return fmt.Sprintf("(%s$%s$%s)\n", v.path, v.currfunc, v.curr[1])
	}
	return fmt.Sprintf("(%s$%s)\n", v.path, v.curr[1])
}
func (v *vm) Goto() string {
	if len(v.currfunc) != 0 {
		return fmt.Sprintf("@%s$%s$%s\n0;JMP\n", v.path, v.currfunc, v.curr[1])
	}
	return fmt.Sprintf("@%s$%s\n0;JMP\n", v.path, v.curr[1])
}
func (v *vm) IfGoto() string {
	if len(v.currfunc) != 0 {
		return fmt.Sprintf("@SP\nAM=M-1\nD=M\n@%s$%s$%s\nD;JNE\n", v.path, v.currfunc, v.curr[1])
	}
	return fmt.Sprintf("@SP\nAM=M-1\nD=M\n@%s$%s$\nD;JNE\n", v.path, v.curr[1])
}
func (v *vm) Call() string {
	temp := fmt.Sprintf("address$%s$%d", v.curr[1], v.count)
	v.count++
	return fmt.Sprintf(
		"@%s\nD=A\n@R13\nM=D\n@%s$%s\nD=A\n@R15\nM=D\n@%s\nD=A\n@R14\nM=D\n@functionCallPrep\n0;JMP\n(%s)\n",
		temp, v.path, v.curr[1], v.curr[2], temp,
	)
}
func (v *vm) Func() string {
	v.currfunc = v.curr[1]
	return fmt.Sprintf("(%s$%s)\n@%s\nD=A\n@SP\nM=D+M\n", v.path, v.curr[1], v.curr[2])
}
func (v *vm) Return() string {
	v.currfunc = ""
	return "@return-prep\n0;JMP\n"
}
func (v *vm) Ops() string {
	var buffer bytes.Buffer
	buffer.WriteString("@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nM=M")
	buffer.WriteString(ops[v.curr[0]])
	buffer.WriteString("D\n@SP\nM=M+1\n")
	return buffer.String()
}
func (v *vm) Negate() string {
	var buffer bytes.Buffer
	buffer.WriteString("@SP\nA=M-1\nM=")
	buffer.WriteString(ops[v.curr[0]])
	buffer.WriteString("M\n")
	return buffer.String()
}

func (v *vm) code() string {
	var output bytes.Buffer
	v.skipEmpty()

	output.WriteString("// " + strings.Join(v.curr, " ") + "\n")

	switch v.curr[0] {
	case "pop":
		output.WriteString(v.Pop())
	case "push":
		output.WriteString(v.Push())
	case "label":
		output.WriteString(v.Label())
	case "goto":
		output.WriteString(v.Goto())
	case "if-goto":
		output.WriteString(v.IfGoto())
	case "call":
		output.WriteString(v.Call())
	case "function":
		output.WriteString(v.Func())
	case "return":
		output.WriteString(v.Return())
	case "lt", "gt", "eq":
		output.WriteString(v.Compare())
	case "not", "neg":
		output.WriteString(v.Negate())
	default:
		output.WriteString(v.Ops())
	}
	v.next()
	return output.String()
}

func codeWriter(file *os.File, path string) chan string {
	c := make(chan string)
	v := new(file, path)
	go func() {
		defer close(c)
		for v.reader.Scan() {
			c <- v.code()
		}
	}()
	return c
}

type builder struct {
	path    string
	name    string
	isDir   bool
	content bytes.Buffer
	file    *os.File
}

func (b *builder) addPrep() {
	b.content.
		WriteString(returnPrep)
	b.content.
		WriteString(functionCallPrep)
}
func check(err error, message string) {
	if err != nil {
		fmt.Printf(message, err.Error())
		os.Exit(1)
	}
}
func (b *builder) dumcContent() {
	dir := filepath.Dir(b.path)
	if b.isDir {
		dir = b.path
	}
	file, err := os.Create(filepath.Join(dir, b.name+".asm"))
	b.addPrep()
	defer file.Close()
	_, err1 := file.WriteString(b.content.String())
	check(err, "bad path:")
	check(err1, "bad path:")

	file.Sync()
}
func newBuilder(path string) *builder {
	fi, err := os.Open(path)
	check(err, "bad path:")
	defer fi.Close()
	fifo, _ := fi.Stat()
	return &builder{
		path:  filepath.Clean(path),
		name:  filepath.Base(path),
		isDir: fifo.IsDir(),
		file:  fi,
	}
}
func (b *builder) readFile(fileName string) {
	fi, _ := os.Open(filepath.Join(b.path, fileName))
	defer fi.Close()
	for line := range codeWriter(fi, fmt.Sprintf("%s$%s", b.name, fileName[:len(fileName)-3])) {
		b.content.WriteString(line)
	}
}
func VM(path string) {
	b := newBuilder(path)
	if !b.isDir {
		b.readFile(b.file.Name())
		b.dumcContent()
		return
	}

	files, err := b.file.ReadDir(0)
	check(err, "unexpected error:")
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".vm") {
			continue
		}
		b.readFile(file.Name())
	}
	b.dumcContent()
}

func clibuild() {
	info()
	app.Commands = []*cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "compiles a .vm file into asm",
			Action: func(c *cli.Context) error {
				path := c.Args().Get(0)
				VM(path)
				return nil
			},
		},
	}
}

func info() {
	app.Name = "vm compiler"
	app.Usage = "compiles the vm"
	app.Authors = []*cli.Author{{Name: "encleine", Email: ""}}
	app.Version = "0.0.2"
}

func main() {
	clibuild()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

type strChan chan []string

func (f strChan) Next() *[]string {
	c, ok := <-f
	if !ok {
		return nil
	}
	return &c
}

func parse(data []byte) strChan {
	c := make(chan []string)
	str := strings.NewReader(string(data))
	r := bufio.NewReader(str)
	go func() {
		defer close(c)
		for {
			token, _, err := r.ReadLine()

			val := strings.TrimSpace(string(token))
			if com := strings.Index(val, "//"); com != -1 {
				val = val[:com]
			}
			switch {
			case len(val) > 0:
				c <- strings.Split(val, " ")
			case err != nil:
				return
			}
		}
	}()
	return c
}

var obs = map[string]string{
	"add": "@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nM=M+D\n@SP\nM=M+1\n",
	"sub": "@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nM=M-D\n@SP\nM=M+1\n",
	"and": "@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nM=M&D\n@SP\nM=M+1\n",
	"or":  "@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nM=M|D\n@SP\nM=M+1\n",

	"gt": "@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nD=D-M\nM=1\n@skip%d\nD;JGT\n@SP\nA=M\nM=0\n(skip%d)\n@SP\nM=M+1\n",
	"lt": "@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nD=D-M\nM=1\n@skip%d\nD;JLT\n@SP\nA=M\nM=0\n(skip%d)\n@SP\nM=M+1\n",
	"eq": "@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nD=D-M\nM=1\n@skip%d\nD;JEQ\n@SP\nA=M\nM=0\n(skip%d)\n@SP\nM=M+1\n",

	"not": "@SP\nA=M-1\nM=!M",
	"neg": "@SP\nA=M-1\nM=-M",
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

const (
	local int = iota + 1
	argument
	this
	that
)

var seg = map[string]int{
	"local":    local,
	"argument": argument,
	"this":     this,
	"that":     that,
}

func ternar(v bool) func(v1 any, v2 any) any {
	return func(v1 any, v2 any) any {
		if v {
			return v1
		}
		return v2
	}
}

func codeWriter(path string) (dummy string) {
	data, err := os.ReadFile(path)
	check(err)
	labCount := 0
	// dummy = "@256\nD = A\n@SP\nM = D\n"

	for lst := range parse(data) {
		temp := ""
		dummy += "// " + strings.Join(lst[:], " ") + "\n"
		switch op := lst[0]; op {
		case "pop":
			switch lst[1] {
			case "local", "argument", "this", "that":
				temp = "@%s\nD=A\n@%d\nD=M+D\n@R13\nM=D\n@SP\nAM=M-1\nD=M\n@R13\nA=M\nM=D\n"
				dummy += fmt.Sprintf(temp, lst[2], seg[op])
				continue
			default:
				temp = "@SP\nAM=M-1\nD=M\n@%s\nM=D\n"
			}
		case "push":
			switch lst[1] {
			case "constant":
				temp = "@%s\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
				dummy += fmt.Sprintf(temp, lst[2])
				continue
			case "local", "argument", "this", "that":
				temp = "@%s\nD=A\n@%d\nA=M+D\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
				dummy += fmt.Sprintf(temp, lst[2], seg[op])
				continue
			default:
				temp = "@%s\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
			}
		case "lt", "gt", "eq":
			dummy += fmt.Sprintf(obs[lst[0]], labCount, labCount)
			labCount++
			continue
		default:
			dummy += obs[lst[0]]
			continue
		}

		switch op := lst[1]; op {
		case "static": //> 16 - 255
			dummy += fmt.Sprintf(temp, filepath.Base(path)+lst[2])
		case "pointer": //> pointer 0/1
			dummy += fmt.Sprintf(temp, ternar(lst[2] == "0")("3", "4"))
		case "temp": //> 5 - 12
			num, err := strconv.Atoi(lst[2])
			check(err)
			dummy += fmt.Sprintf(temp, strconv.Itoa(5+num))
		}

	}
	return
}

func build(c *cli.Context) error {
	path := c.Args().Get(0)
	if !strings.HasSuffix(path, ".vm") {
		fmt.Println("not a supported file format")
		return nil
	}
	file, e := os.Create(fmt.Sprintf("%s.%s", path[:len(path)-3], "asm"))
	check(e)
	dummy := codeWriter(path)
	defer file.Close()
	_, err := file.WriteString(dummy)
	file.Sync()
	check(err)
	return nil
}

var app = cli.NewApp()

func clibuild() {
	info()
	app.Commands = []*cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "vm compiler",
			Action:  build,
		},
	}
}

func info() {
	app.Name = "vm compiler"
	app.Usage = "compiles the vm"
	app.Authors = []*cli.Author{{Name: "encleine", Email: ""}}
	app.Version = "0.0.1"
}

func main() {
	clibuild()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

var symbolTable = map[string]uint16{
	"R1": 1, "R0": 0, "R2": 2, "R3": 3, "R4": 4, "R5": 5, "R6": 6, "R7": 7,
	"R8": 8, "R9": 9, "R10": 10, "R11": 11, "R12": 12, "R13": 13, "R14": 14, "R15": 15,
	"SP": 0, "LCL": 1, "ARG": 2, "THIS": 3, "THAT": 4,
	"SCREEN": 16384, "KBD": 24576,
}
var labelTable = map[string]uint16{}

var compTable = map[string]string{
	"0": "0101010", "1": "0111111", "-1": "0111010",
	"D": "0001100", "A": "0110000", "M": "1110000",
	"!D": "0001101", "!A": "0110001", "!M": "1110001",
	"-D": "0001111", "-A": "0110011", "-M": "1110011",
	"D+1": "0011111", "A+1": "0110111", "M+1": "1110111",
	"D-1": "0001110", "A-1": "0110010", "M-1": "1110010",
	"D+A": "0000010", "D+M": "1000010",
	"D-A": "0010011", "D-M": "1010011",
	"A-D": "0000111", "M-D": "1000111",
	"D&A": "0000000", "D&M": "1000000",
	"D|A": "0010101", "D|M": "1010101",
}
var jumpTable = map[string]string{
	"JMP": "111",
	"JNE": "101",
	"JEQ": "010",
	"JGT": "001",
	"JLT": "100",
	"JGE": "011",
	"JLE": "110",
}
var distTable = map[string]string{
	"M":   "001",
	"D":   "010",
	"DM":  "011",
	"MD":  "011",
	"A":   "100",
	"AM":  "101",
	"MA":  "101",
	"AD":  "110",
	"DA":  "110",
	"ADM": "111",
}

func getlines(data string) chan string {
	c := make(chan string)
	str := strings.NewReader(data)
	r := bufio.NewReader(str)
	go func() {
		defer close(c)
		for {
			token, _, err := r.ReadLine()

			val := strings.ReplaceAll(strings.TrimSpace(string(token)), " ", "")
			if com := strings.Index(val, "//"); com != -1 {
				val = val[:com]
			}

			switch {
			case len(val) > 0:
				c <- val
			case err != nil:
				return
			}
		}
	}()
	return c
}

type token map[string][]string

func firstPass(filePath string) (items []token) {
	var lineCount uint16 = 0
	dat, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	for line := range getlines(string(dat)) {
		le := string(line[0])
		switch le {
		case "@":
			lineCount++
			items = append(items, token{"C": []string{line[1:]}})
			continue
		case "(":
			labelTable[strings.Trim(line, "()")] = lineCount
			continue
		default:
			lineCount++

			re := regexp.MustCompile("[^A-Za-z0-9+-/&/|!]+")
			split := re.Split(line, -1)

			if strings.Contains(line, ";") {
				items = append(items, token{"Aj": split})
			} else {
				items = append(items, token{"A": split})
			}
		}
	}
	return
}

func firstKey(mymap token) (key string) {
	for k := range mymap {
		return k
	}
	return
}

func secondPass(items []token) []token {
	var r = regexp.MustCompile("^[a-zA-z]")
	var varCount uint16 = 16
	for _, v := range items {
		switch firstKey(v) {
		case "A", "Aj":
			continue
		}

		key := v["C"][0]
		if val, ok := symbolTable[key]; ok {
			key = strconv.FormatInt(int64(val), 10)
		} else if val, ok := labelTable[key]; ok {
			key = strconv.FormatInt(int64(val), 10)
		} else if r.MatchString(key) {
			labelTable[key] = varCount
			key = strconv.FormatInt(int64(varCount), 10)
			varCount++
		}
		v["C"] = []string{key}
	}
	return items
}

func lastPass(items []token) (prog string) {
	line := []byte("0000000000000000")
	for _, v := range items {
		key := firstKey(v)
		switch key {
		case "C":
			val, err := strconv.ParseInt(v["C"][0], 0, 64)
			if err != nil {
				panic(err)
			}
			line = []byte("0" + fmt.Sprintf("%015b", val))
		case "A", "Aj":
			copy(line[:3], []byte("111"))
			list := v["A"]
			if key == "Aj" {
				list = v["Aj"]
			}
			for i, v := range list {
				switch i {
				case 0:
					if val, ok := distTable[v]; ok && (key != "Aj" && len(v) <= 2) {
						copy(line[10:13], []byte(val))
					} else if val, ok := compTable[v]; ok {
						copy(line[3:10], []byte(val))
					}
				case 1:
					if val, ok := compTable[v]; ok {
						copy(line[3:10], []byte(val))
					} else if val, ok := jumpTable[v]; ok {
						copy(line[13:16], []byte(val))
					}
				case 2:
					if val, ok := jumpTable[v]; ok {
						copy(line[13:16], []byte(val))
					}
				}
			}
		}
		prog += string(line) + "\n"
		line = []byte("0000000000000000")
	}
	return
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

var app = cli.NewApp()

func build() {
	app.Commands = []*cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "assembles hack .asm files down to hack binary",
			Action: func(c *cli.Context) error {
				path := c.Args().Get(0)
				if !strings.HasSuffix(path, ".asm") {
					fmt.Println("not a supported file format")
					return nil
				}
				mymap := firstPass(path)
				nxtmap := secondPass(mymap)
				lstmap := lastPass(nxtmap)

				file, e := os.Create(path[:len(path)-3] + "hack")
				check(e)
				defer file.Close()
				_, err := file.WriteString(lstmap)
				file.Sync()
				check(err)
				return nil
			},
		},
	}
}

func info() {
	app.Name = "hack assembler"
	app.Usage = "assembles the hack"
	app.Authors = []*cli.Author{{Name: "encleine", Email: ""}}
	app.Version = "0.0.1"
}

// * Hack asm so hasm
func main() {
	info()
	build()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

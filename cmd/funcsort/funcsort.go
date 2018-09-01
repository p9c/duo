package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
)

type method struct {
	name   string
	proto  string
	method string
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("a filename is required")
		os.Exit(1)
	}
	if args[1] == "-h" || args[1] == "--help" {
		fmt.Println("go type formatter for pipeline pattern libraries")
		fmt.Println()
		fmt.Println("usage:", os.Args[0], "<libraryfile.go>")
		fmt.Println()
		fmt.Println("can only work with a file with only one struct and one interface defined, which is anyway the pipeline pattern")
		os.Exit(0)
	}
	filename := os.Args[1]
	fh, _ := os.Open(filename)
	out, _ := ioutil.ReadAll(fh)
	s := string(out)
	splitted := strings.Split(s, "\n\n")
	var iface []string
	var methnames []method
	var funcs []string
	var typename string
	var Interface string
	var Struct string
	for i := range splitted {
		words := strings.Fields(splitted[i])
		lines := strings.Split(splitted[i], "\n")
		for j := range lines {
			if lines[j] == "" {
				break
			}
			if len(lines[j]) > 0 && lines[j][0] != '/' {
				words1 := strings.Split(lines[j], " ")
				if words1[0] == "type" && words1[2] == "struct" {
					typename = words1[1]
					Struct = strings.TrimSpace(splitted[i])
					splitted[i] = ""
				}
			}
		}
		if len(words) < 1 {

		} else if words[0] == "type" && words[2] == "interface" {
			head := lines[0]
			tail := lines[len(lines)-1]
			iface = append(iface, head)
			sort.Strings(lines)
			for j := range lines {
				if len(lines[j]) < 1 {
					break
				}
				if lines[j][0] == '\t' {
					iface = append(iface, lines[j])
					methre, _ := regexp.Compile(`[a-zA-Z_][a-zA-Z0-9_]*`)
					mn := methre.FindAll([]byte(lines[j]), 1)
					methnames = append(methnames, method{name: string(mn[0]), proto: lines[j]})
				}
			}
			iface = append(iface, tail)
			Interface = strings.Join(iface, "\n")
			splitted[i] = ""
		}
		for j := range lines {
			if len(lines[j]) > 1 {
				if lines[j][0] != '/' {
					if lines[j][0] == 'f' && lines[j][1] == 'u' && lines[j][2] == 'n' && lines[j][3] == 'c' {
						if lines[j][5] == '(' {
							words := strings.Split(lines[j], " ")
							name := strings.Split(words[3], "(")[0]
							for k := range methnames {
								if methnames[k].name == name {
									methnames[k].method = strings.TrimSpace(splitted[i])
									splitted[i] = ""
								}
							}
						} else {
							funcs = append(funcs, strings.TrimSpace(splitted[i]))
							splitted[i] = ""
						}
						break
					}
				}
			}
		}
	}
	for i := range methnames {
		if methnames[i].method == "" {
			methnames[i].method = "// " + methnames[i].name + " does...\n" +
				"func (r *" + typename + ") " + methnames[i].proto[1:] + " {\n\treturn r\n}"
		}
	}
	var output string
	for i := range splitted {
		if splitted[i] != "" {
			output += splitted[i] + "\n\n"
		}
	}
	output += Struct + "\n\n"
	for i := range funcs {
		output += funcs[i] + "\n\n"
	}
	output += Interface + "\n\n"
	for i := range methnames {
		output += methnames[i].method
		if i == len(methnames)-1 {
			output += "\n"
		} else {
			output += "\n\n"
		}
	}
	fh, err := os.OpenFile(os.Args[1], os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("unable to write to file")
		os.Exit(1)
	}
	fh.WriteString(output)
}

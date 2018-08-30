package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var S skeleton

type skeleton struct {
	pkg   string
	imp   []string
	name  string
	iface string
	vals  [][]string
	meths [][]string
}

func main() {
	text, _ := ioutil.ReadAll(os.Stdin)
	words := strings.Fields(string(text))
	var i int
	for {
		if words[i] == "package" {
			i++
			S.pkg = words[i]
			i++
		}
		if words[i] == "import" {
			i++
			if words[i] == "(" {
				i++
				for words[i] != ")" {
					S.imp = append(S.imp, words[i])
					i++
				}
				i++
			}
			S.imp = append(S.imp, words[i])
		}
		if words[i] == "type" {
			i++
			if words[i+1] == "struct" {
				S.name = words[i]
				i += 3
				for words[i] != "}" {
					S.vals = append(S.vals, []string{words[i], words[i+1]})
					i += 2
				}
				i += 3
			}
			S.iface = words[i]
			i += 2
			for words[i] != "}" {
				splits := strings.Split(words[i], "(")
				name := splits[0]
				fmt.Println(name)
				if len(splits) > 1 {
					if splits[1] == ")" {
						i++
					} else {
						splitten := strings.Split(splits[1], ",")
						fmt.Println(splitten)
						fmt.Println(splits[1])
					}
				}
				fmt.Println(words[i])
				i++
			}
		}
		break
	}
	fmt.Println(S)
}

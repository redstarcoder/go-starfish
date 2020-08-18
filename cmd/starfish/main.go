package main

import (
	"flag"
	"fmt"
	"github.com/redstarcoder/go-starfish/starfish"
	"io/ioutil"
	"os"
	"time"
)

var (
	showcodebox        = flag.Bool("c", false, "output the codebox each tick")
	flagscript         = flag.String("code", "", "execute the script supplied in 'code'")
	showstack          = flag.Bool("s", false, "output the stack each tick")
	help         *bool = flag.Bool("h", false, "display this help message")
	delay              = flag.Duration("t", 0, "time to sleep between ticks (ex: 100ms)")
	compmode           = flag.Bool("m", false, "run like the fishlanguage.com interpreter")
	initialstack       = &stack{[]float64{}}
	fName              = "fish"
)

func Error() {
	fmt.Println("Usage:", fName, "[args] <file>")
	flag.PrintDefaults()
}

func loadScript(fName string) string {
	file, err := os.Open(fName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func init() {
	fName = os.Args[0]
	flag.Var(initialstack, "i", "set the initial stack (ex: '\"Example\" 10 \"stack\"')")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if *help || (*flagscript == "" && len(args) == 0) {
		Error()
		return
	}
	var script string
	if script = *flagscript; script == "" {
		script = loadScript(args[0])
	}

	cB := starfish.NewCodeBox(script, initialstack.s, *compmode)
	if !*showcodebox && !*showstack && *delay == 0 {
		var (
			end    bool
			output string
		)
		for ; !end; output, end = cB.Swim() {
			if output != "" {
				fmt.Print(output)
			}
		}
		return
	}
	if *showcodebox {
		cB.PrintBox()
	}
	if *showstack && cB.StackLength() > 0 {
		fmt.Println("Stack:", cB.Stack())
	}
	time.Sleep(*delay)
	var (
		end    bool
		output string
	)
	for ; !end; output, end = cB.Swim() {
		if output != "" {
			fmt.Print(output)
		}
		if *showcodebox {
			cB.PrintBox()
		}
		if *showstack && cB.StackLength() > 0 {
			fmt.Println("Stack:", cB.Stack())
		}
		time.Sleep(*delay)
	}
}

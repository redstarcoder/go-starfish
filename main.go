package main

import (
	"flag"
	"fmt"
	"github.com/redstarcoder/go-fish/fish"
	"io/ioutil"
	"os"
	"time"
)

var (
	showcodebox = flag.Bool("c", false, "outputs the codebox each tick")
	showstack = flag.Bool("s", false, "outputs the stack each tick")
	help *bool = flag.Bool("h", false, "displays this help message")
	delay = flag.Duration("t", 0, "time to sleep between ticks (ex: 100ms)")
	initialstack = &stack{[]float64{}}
	fName = "fish"
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
	flag.Var(initialstack, "i", "sets the initial stack (ex: '\"Example\" 10 \"stack\"')")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if *help || len(args) == 0 {
		Error()
		return
	}
	script := loadScript(args[0])

	cB := fish.NewCodeBox(script, initialstack.s)
	if !*showcodebox && !*showstack && *delay == 0 {
		for !cB.Swim() {}
		return
	}
	if *showcodebox {
		cB.PrintBox()
	}
	if *showstack && cB.StackLength() > 0 {
		fmt.Println("Stack:", cB.Stack())
	}
	time.Sleep(*delay)
	for !cB.Swim() {
		if *showcodebox {
			cB.PrintBox()
		}
		if *showstack && cB.StackLength() > 0 {
			fmt.Println("Stack:", cB.Stack())
		}
		time.Sleep(*delay)
	}
}

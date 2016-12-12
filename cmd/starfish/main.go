package main

import (
	//"flag"
	"fmt"
	"github.com/redstarcoder/go-starfish/starfish"
	//"io/ioutil"
	"github.com/gopherjs/gopherjs/js"
	//"os"
	"time"
)

var (
	/*showcodebox = flag.Bool("c", false, "output the codebox each tick")
	flagscript = flag.String("code", "", "execute the script supplied in 'code'")
	showstack = flag.Bool("s", false, "output the stack each tick")
	help *bool = flag.Bool("h", false, "display this help message")
	delay = flag.Duration("t", 0, "time to sleep between ticks (ex: 100ms)")
	compmode = flag.Bool("m", false, "run like the fishlanguage.com interpreter")*/
	initialstack = &stack{[]float64{}}
)

func init() {
}

func main() {
	stop := false
	pause := false
	output := js.Global.Get("output")
	input := js.Global.Get("input")
	stack := js.Global.Get("stack")
	codebox := js.Global.Get("codebox")
	script := js.Global.Get("script")
	delay := js.Global.Get("delay")
	run := js.Global.Get("run")
	istack := js.Global.Get("initialstack")
	inputfield := js.Global.Get("inputfield")
	give := js.Global.Get("give")
	
	give.Call("addEventListener", "click", func() {
		input.Set("innerHTML", input.Get("innerHTML").String() + inputfield.Get("value").String())
		inputfield.Set("value", "")
	})
	
	js.Global.Get("end").Call("addEventListener", "click", func() {
		stop = true
		pause = false
	})
	
	js.Global.Get("pause").Call("addEventListener", "click", func() {
		pause = !pause
	})
	
	run.Call("addEventListener", "click", func() {
		run.Set("disabled", true)
		go func () {
			output.Set("innerHTML", "")
			stop = false
			pause = false
			delayms := delay.Get("value").Int()
			initialstack.Set(istack.Get("value").String())
			cB := starfish.NewCodeBox(script.Get("value").String(), initialstack.s, false, starfish.JSObjects{output, input, stack, codebox})
			cB.PrintBox()
			stack.Set("innerHTML", fmt.Sprintln(cB.Stack()))
			time.Sleep(time.Millisecond*time.Duration(delayms))
			for !cB.Swim() && !stop {
				cB.PrintBox()
				stack.Set("innerHTML", fmt.Sprintln(cB.Stack()))
				time.Sleep(time.Millisecond*time.Duration(delayms))
				for pause {
					time.Sleep(time.Millisecond*200)
				}
			}
			run.Set("disabled", false)
		}()
	})
	/*if *showcodebox {
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
	}*/
	
	//output.Set("innerHTML", output.Get("innerHTML").String() + "<br />Test")
}

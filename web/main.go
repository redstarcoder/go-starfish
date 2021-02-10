package main

import (
	"fmt"
	"github.com/redstarcoder/go-starfish/starfish"
	"strconv"
	"strings"
	"syscall/js"
)

var (
	globalCB       *starfish.CodeBox
	ExecutionTimer js.Value
	Paused         bool
	Run            bool // true while a script is running, can be changed to false to stop execution
)

func CreateBox() string {
	width, height := globalCB.Size()
	output := make([]byte, 0, width*7*height+(height+2)*4+10)
	output = append(output, []byte("<br>")...)
	box := globalCB.Box()
	for y, line := range box {
		for x, r := range line {
			if x != 0 || y != 0 {
				output = append(output, []byte(fmt.Sprintf(`<c id="%dx%d">`, x, y))...)
			} else {
				output = append(output, []byte(fmt.Sprintf(`<c id="%dx%d" class="s">`, x, y))...)
			}
			output = append(output, byte(r))
			output = append(output, []byte("</c>")...)
		}
		output = append(output, []byte("<br>")...)
	}
	return "<pre>" + string(output) + "<br></pre>"
}

func shareScript(this js.Value, args []js.Value) interface{} {
	window := js.Global()
	document := window.Get("document")

	url := window.Get("location").Get("href").String()
	if strings.Contains(url, "?") {
		url = strings.Split(url, "?")[0]
	}

	shareField := document.Call("getElementById", "sharefield")
	shareField.Set("value", url+"?script="+window.Get("LZString").Call("compressToEncodedURIComponent", document.Call("getElementById", "script").Get("value").String()).String())

	shareBox := document.Call("getElementById", "sharebox")
	shareBox.Call("setAttribute", "style", "")
	return nil
}

func stopscript(this js.Value, args []js.Value) interface{} {
	window := js.Global()
	window.Call("clearTimeout", ExecutionTimer)
	globalCB = nil
	Run = false
	Paused = false
	return nil
}

func pausescript(this js.Value, args []js.Value) interface{} {
	if globalCB == nil || !Run {
		return nil
	}
	window := js.Global()
	if !Paused {
		window.Call("clearTimeout", ExecutionTimer)
		Paused = true
	} else {
		document := window.Get("document")
		delay, err := strconv.Atoi(document.Call("getElementById", "delay").Get("value").String())
		if err != nil {
			panic("Delay must be a number, no decimals.")
		}
		stepRun(delay)
		Paused = false
	}
	return nil
}

func stepRun(delay int) {
	window := js.Global()
	document := window.Get("document")

	stepRunJS := func(this js.Value, args []js.Value) interface{} {
		go func() { // a bit hacky, but this throws everything into a goroutine so a sleep for example, won't block.
			outBox := document.Call("getElementById", "output")
			outString := outBox.Get("innerText").String()
			stackBox := document.Call("getElementById", "stack")

			x, y := globalCB.Loc()
			width, height := globalCB.Size()

			if x < width && y < height {
				document.Call("getElementById", fmt.Sprintf("%dx%d", x, y)).Call("setAttribute", "class", "")
			}

			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
					rString, _ := r.(string)
					outString += "\n" + rString
					outBox.Set("innerText", outString)
					if globalCB != nil {
						stackBox.Set("innerText", fmt.Sprintln(globalCB.Stack()))

						if x < width && y < height {
							if !globalCB.DeepSea() {
								document.Call("getElementById", fmt.Sprintf("%dx%d", x, y)).Call("setAttribute", "class", "s")
							} else {
								document.Call("getElementById", fmt.Sprintf("%dx%d", x, y)).Call("setAttribute", "class", "u")
							}
						}
					}
					stopscript(this, args)
				}
			}()

			shareScript(this, args)

			output, end := globalCB.Swim()
			if output != "" {
				if output[0] == 13 { // if carriage return...
					outString = ""
				} else {
					outString += output
				}
				outBox.Set("innerText", outString)
				// print(output)
			}
			if globalCB != nil {
				stackBox.Set("innerText", fmt.Sprintln(globalCB.Stack()))
			} else {
				end = true
			}
			if !end {
				x, y = globalCB.Loc()
				ExecutionTimer = window.Call("setTimeout", "stepRun();", delay).JSValue()
			}

			if !globalCB.DeepSea() {
				document.Call("getElementById", fmt.Sprintf("%dx%d", x, y)).Call("setAttribute", "class", "s")
			} else {
				document.Call("getElementById", fmt.Sprintf("%dx%d", x, y)).Call("setAttribute", "class", "u")
			}

			if end {
				print("\n")
				globalCB = nil
				Run = false
			}
		}()
		return nil
	}
	window.Set("stepRun", js.FuncOf(stepRunJS))
	window.Call("stepRun")
}

func runscript(this js.Value, args []js.Value) interface{} {
	if Run {
		return nil
	} else {
		Run = true
	}

	window := js.Global()
	document := window.Get("document")

	//output
	outBox := document.Call("getElementById", "output")
	outString := ""

	stackBox := document.Call("getElementById", "stack")

	defer func() {
		if r := recover(); r != nil {
			println(r.(string))
			outString += "\n" + r.(string)
			outBox.Set("innerText", outString)
			//stackBox.Set("innerText", fmt.Sprintln(cB.Stack()))
		}
	}()

	// initial stack
	stack := &stack{[]float64{}}
	stack.Set(document.Call("getElementById", "initialstack").Get("value").String())

	// codeBox
	cB := starfish.NewCodeBox(document.Call("getElementById", "script").Get("value").String(), stack.s, false)

	// delay
	delay, err := strconv.Atoi(document.Call("getElementById", "delay").Get("value").String())
	if err != nil {
		panic("Delay must be a number, no decimals.")
	}

	if delay <= 0 {
		go func() { // a bit hacky, but this throws everything into a goroutine so a sleep for example, won't block.
			var (
				end    bool
				output string
			)
			for ; !end && Run; output, end = cB.Swim() {
				if output != "" {
					if output[0] == 13 { // if carriage return...
						outString = ""
					} else {
						outString += output
					}
					outBox.Set("innerText", outString)
				}
			}
			stackBox.Set("innerText", fmt.Sprintln(cB.Stack()))
			print("\n")
			Run = false
		}()
	} else {
		globalCB = cB
		document.Call("getElementById", "output").Set("innerText", "")
		document.Call("getElementById", "codebox").Set("innerHTML", CreateBox())
		stepRun(delay)
	}
	return nil
}

func setup() {
	window := js.Global()
	document := window.Get("document")

	// callbacks
	window.Set("runscript", js.FuncOf(runscript))
	document.Call("getElementById", "run").Call("setAttribute", "onClick", "runscript();")
	window.Set("stopscript", js.FuncOf(stopscript))
	document.Call("getElementById", "end").Call("setAttribute", "onClick", "stopscript();")
	window.Set("pausescript", js.FuncOf(pausescript))
	document.Call("getElementById", "pause").Call("setAttribute", "onClick", "pausescript();")
	window.Set("sharescript", js.FuncOf(shareScript))
	document.Call("getElementById", "share").Call("setAttribute", "onClick", "sharescript();")

	// load script (if any)
	if s := window.Call("getUrlVars").Get("script").String(); s != "<undefined>" {
		document.Call("getElementById", "script").Set("value", window.Get("LZString").Call("decompressFromEncodedURIComponent", s).String())
	}
}

func main() {
	println("👍")
	// register functions
	setup()

	<-make(chan bool)
}

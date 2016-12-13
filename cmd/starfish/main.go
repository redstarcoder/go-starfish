package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/redstarcoder/go-starfish/starfish"
	"strings"
	"time"
)

var (
	initialstack = &stack{[]float64{}}
)

func init() {
}

func main() {
	stop := false
	pause := false
	delayms := 0
	showcodebox := false
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
	sharefield := js.Global.Get("sharefield")
	share := js.Global.Get("share")
	sharebox := js.Global.Get("sharebox")
	showhide := js.Global.Get("showhide")

	if s := js.Global.Call("getUrlVars").Get("script").String(); s != "undefined" {
		script.Set("value", js.Global.Get("LZString").Call("decompressFromEncodedURIComponent", s).String())
	}

	url := js.Global.Get("window").Get("location").Get("href").String()
	sharefield.Set("value", url)
	if strings.Contains(url, "?") {
		url = strings.Split(url, "?")[0]
	}

	showhide.Call("addEventListener", "click", func() {
		if showcodebox {
			showcodebox = false
			showhide.Set("innerHTML", "Show CodeBox")
			codebox.Get("style").Set("display", "none")
		} else {
			showcodebox = true
			showhide.Set("innerHTML", "Hide CodeBox")
			codebox.Get("style").Set("display", "inline-block")
		}
		sharefield.Set("value", url+"?script="+js.Global.Get("LZString").Call("compressToEncodedURIComponent", script.Get("value").String()).String())
	})

	share.Call("addEventListener", "click", func() {
		sharefield.Set("value", url+"?script="+js.Global.Get("LZString").Call("compressToEncodedURIComponent", script.Get("value").String()).String())
		sharebox.Get("style").Set("display", "block")
		sharefield.Call("select")
	})

	give.Call("addEventListener", "click", func() {
		input.Set("innerHTML", input.Get("innerHTML").String()+inputfield.Get("value").String())
		inputfield.Set("value", "")
	})

	js.Global.Get("end").Call("addEventListener", "click", func() {
		stop = true
		pause = false
	})

	js.Global.Get("pause").Call("addEventListener", "click", func() {
		pause = !pause
		delayms = delay.Get("value").Int()
	})

	run.Call("addEventListener", "click", func() {
		run.Set("disabled", true)
		output.Set("innerHTML", "")
		sharefield.Set("value", url+"?script="+js.Global.Get("LZString").Call("compressToEncodedURIComponent", script.Get("value").String()).String())
		go func() {
			stop = false
			pause = false
			delayms = delay.Get("value").Int()
			initialstack.Set(istack.Get("value").String())
			cB := starfish.NewCodeBox(script.Get("value").String(), initialstack.s, false, starfish.JSObjects{output, input, stack, codebox})
			cB.PrintBox()
			stack.Set("innerHTML", fmt.Sprintln(cB.Stack()))
			time.Sleep(time.Millisecond * time.Duration(delayms))
			for !cB.Swim() && !stop {
				if showcodebox {
					cB.PrintBox()
				}
				stack.Set("innerHTML", fmt.Sprintln(cB.Stack()))
				time.Sleep(time.Millisecond * time.Duration(delayms))
				for pause {
					if !showcodebox {
						cB.PrintBox()
					}
					time.Sleep(time.Millisecond * 200)
				}
			}
			run.Set("disabled", false)
		}()
	})
}

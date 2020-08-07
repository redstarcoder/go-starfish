// See https://esolangs.org/wiki/Starfish for more info.
package starfish

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Direction is a value representing the direction a ><> is swimming.
type Direction byte

const (
	Right Direction = iota
	Down
	Left
	Up
)

var reader chan byte

// Stack is a type representing a stack in ><>. It holds the stack values in S, as well as a register. The
// register may contain data, but will only be considered filled if filledRegister is also true.
type Stack struct {
	S              []float64
	register       float64
	filledRegister bool
}

// NewStack returns a pointer to a Stack populated with s.
func NewStack(s []float64) *Stack {
	newS := make([]float64, len(s))
	copy(newS, s)
	return &Stack{S: newS}
}

// Register implements "&".
func (s *Stack) Register() {
	if s.filledRegister {
		s.Push(s.register)
		s.filledRegister = false
	} else {
		s.register = s.Pop()
		s.filledRegister = true
	}
}

// Extend implements ":".
func (s *Stack) Extend() {
	s.Push(s.S[len(s.S)-1])
}

// Reverse implements "r".
func (s *Stack) Reverse() {
	newS := make([]float64, len(s.S))
	for i, ii := 0, len(s.S)-1; ii >= 0; i, ii = i+1, ii-1 {
		newS[i] = s.S[ii]
	}
	s.S = newS
}

// SwapTwo implements "$".
func (s *Stack) SwapTwo() {
	s.S[len(s.S)-1], s.S[len(s.S)-2] = s.S[len(s.S)-2], s.S[len(s.S)-1]
}

// SwapThree implements "@": with [1,2,3,4], calling "@" results in [1,4,2,3].
func (s *Stack) SwapThree() {
	s.S[len(s.S)-1], s.S[len(s.S)-2], s.S[len(s.S)-3] = s.S[len(s.S)-2], s.S[len(s.S)-3], s.S[len(s.S)-1]
}

// ShiftRight implements "}".
func (s *Stack) ShiftRight() {
	newS := make([]float64, 1, len(s.S))
	newS[0] = s.Pop()
	s.S = append(newS, s.S...)
}

// ShiftLeft implements "{".
func (s *Stack) ShiftLeft() {
	r := s.S[0]
	s.S = s.S[1:]
	s.Push(r)
}

// Push appends r to the end of the stack.
func (s *Stack) Push(r float64) {
	s.S = append(s.S, r)
}

// Pop removes the value on the end of the stack and returns it.
func (s *Stack) Pop() (r float64) {
	if len(s.S) > 0 {
		r = s.S[len(s.S)-1]
		s.S = s.S[:len(s.S)-1]
	} else {
		panic("Stack is empty!")
	}
	return
}

// getBytes removes c values from the stack, then returns them as a byte slice.
func (s *Stack) getBytes(c int) []byte {
	sData := s.S[len(s.S)-c:]
	s.S = s.S[:len(s.S)-c]
	bData := make([]byte, c)
	for i, v := range sData {
		bData[i] = byte(v)
	}
	return bData
}

func longestLineLength(lines []string) (l int) {
	for _, s := range lines {
		if len(s) > l {
			l = len(s)
		}
	}
	return
}

// CodeBox is an object usually created with NewCodeBox. It contains a ><> program complete with a stack,
// and is typically run in steps via CodeBox.Swim.
type CodeBox struct {
	fX, fY        int
	fDir          Direction
	wasLeft       bool
	escapedHook   bool
	width, height int
	box           [][]byte
	stacks        []*Stack
	p             int // Used to keep track of the current stack
	stringMode    byte
	compMode      bool
	deepSea       bool
	file          *os.File
}

// NewCodeBox returns a pointer to a new CodeBox. "script" should be a complete ><> script, "stack" should
// be the initial stack, and compatibilityMode should be set if fishinterpreter.com behaviour is needed.
func NewCodeBox(script string, stack []float64, compatibilityMode bool) *CodeBox {
	cB := new(CodeBox)

	script = strings.Replace(script, "\r", "", -1)
	if len(script) == 0 || script == "\n" {
		panic("Cannot accept script of length 0 (No room for the fish to survive).")
	}

	lines := strings.Split(script, "\n")
	cB.width = longestLineLength(lines)
	cB.height = len(lines)

	cB.box = make([][]byte, cB.height)
	for i, s := range lines {
		cB.box[i] = make([]byte, cB.width)
		for ii, r := 0, byte(0); ii < cB.width; ii++ {
			if ii < len(s) {
				r = byte(s[ii])
			} else {
				r = ' '
			}
			cB.box[i][ii] = byte(r)
		}
	}

	cB.stacks = []*Stack{NewStack(stack)}
	cB.compMode = compatibilityMode

	return cB
}

// Exe executes the instruction the ><> is currently on top of. It returns the string it intends to output (nil if none) and true when it executes ";".
func (cB *CodeBox) Exe(r byte) (string, bool) {
	switch r {
	case ' ':
		return "", false
	case '>':
		cB.fDir = Right
		cB.wasLeft = false
		return "", false
	case 'v':
		cB.fDir = Down
		return "", false
	case '<':
		cB.fDir = Left
		cB.wasLeft = true
		return "", false
	case '^':
		cB.fDir = Up
		return "", false
	case '|':
		if cB.fDir == Right {
			cB.fDir = Left
			cB.wasLeft = true
		} else if cB.fDir == Left {
			cB.fDir = Right
			cB.wasLeft = false
		}
		return "", false
	case '_':
		if cB.fDir == Down {
			cB.fDir = Up
		} else if cB.fDir == Up {
			cB.fDir = Down
		}
		return "", false
	case '#':
		switch cB.fDir {
		case Right:
			cB.fDir = Left
			cB.wasLeft = true
		case Down:
			cB.fDir = Up
		case Left:
			cB.fDir = Right
			cB.wasLeft = false
		case Up:
			cB.fDir = Down
		}
		return "", false
	case '/':
		switch cB.fDir {
		case Right:
			cB.fDir = Up
		case Down:
			cB.fDir = Left
			cB.wasLeft = true
		case Left:
			cB.fDir = Down
		case Up:
			cB.fDir = Right
			cB.wasLeft = false
		}
		return "", false
	case '\\':
		switch cB.fDir {
		case Right:
			cB.fDir = Down
		case Down:
			cB.fDir = Right
			cB.wasLeft = false
		case Left:
			cB.fDir = Up
		case Up:
			cB.fDir = Left
			cB.wasLeft = true
		}
		return "", false
	case 'x':
		cB.fDir = Direction(rand.Int31n(4))
		switch cB.fDir {
		case Right:
			cB.wasLeft = false
		case Left:
			cB.wasLeft = true
		}
		return "", false
	// *><> commands
	case 'O':
		cB.deepSea = false
		return "", false
	case '`':
		if cB.fDir == Down || cB.fDir == Up {
			if cB.wasLeft {
				cB.fDir = Left
			} else {
				cB.fDir = Right
			}
		} else {
			if cB.escapedHook {
				cB.fDir = Up
				cB.escapedHook = false
			} else {
				cB.fDir = Down
				cB.escapedHook = true
			}
		}
		return "", false
	}

	if cB.deepSea {
		return "", false
	}

	var output string

	switch r {
	default:
		panic(string(r))
	case ';':
		return "", true
	case '"', '\'':
		if cB.stringMode == 0 {
			cB.stringMode = r
		} else if r == cB.stringMode {
			cB.stringMode = 0
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		cB.Push(float64(r - '0'))
	case 'a', 'b', 'c', 'd', 'e', 'f':
		cB.Push(float64(r - 'a' + 10))
	case '&':
		cB.Register()
	case 'o':
		output = string(rune(cB.Pop()))
	case 'n':
		output = fmt.Sprintf("%v", cB.Pop())
	case 'r':
		cB.ReverseStack()
	case '+':
		cB.Push(cB.Pop() + cB.Pop())
	case '-':
		x := cB.Pop()
		y := cB.Pop()
		cB.Push(y - x)
	case '*':
		cB.Push(cB.Pop() * cB.Pop())
	case ',':
		x := cB.Pop()
		y := cB.Pop()
		cB.Push(y / x)
	case '%':
		x := cB.Pop()
		y := cB.Pop()
		cB.Push(float64(int64(y) % int64(x)))
	case '=':
		if cB.Pop() == cB.Pop() {
			cB.Push(1)
		} else {
			cB.Push(0)
		}
	case ')':
		x := cB.Pop()
		y := cB.Pop()
		if y > x {
			cB.Push(1)
		} else {
			cB.Push(0)
		}
	case '(':
		x := cB.Pop()
		y := cB.Pop()
		if y < x {
			cB.Push(1)
		} else {
			cB.Push(0)
		}
	case '!':
		cB.Move()
	case '?':
		if cB.Pop() == 0 {
			cB.Move()
		}
	case '.':
		cB.fY = int(cB.Pop())
		cB.fX = int(cB.Pop())
	case ':':
		cB.ExtendStack()
	case '~':
		cB.Pop()
	case '$':
		cB.StackSwapTwo()
	case '@':
		cB.StackSwapThree()
	case '}':
		cB.StackShiftRight()
	case '{':
		cB.StackShiftLeft()
	case ']':
		cB.CloseStack()
	case '[':
		cB.NewStack(int(cB.Pop()))
	case 'l':
		cB.Push(cB.StackLength())
	case 'g':
		cB.Push(float64(cB.box[int(cB.Pop())][int(cB.Pop())]))
	case 'p':
		cB.box[int(cB.Pop())][int(cB.Pop())] = byte(cB.Pop())
	case 'i':
		r := float64(-1)
		if cB.file == nil {
			b := byte(0)
			select {
			case b = <-reader:
				r = float64(b)
			default:
			}
		} else {
			bs := []byte{0}
			n, _ := cB.file.Read(bs)
			if n > 0 {
				r = float64(bs[0])
			}
		}
		cB.Push(r)
	// *><> commands
	case 'h':
		cB.Push(float64(time.Now().Hour()))
	case 'm':
		cB.Push(float64(time.Now().Minute()))
	case 's':
		cB.Push(float64(time.Now().Second()))
	case 'S':
		time.Sleep(time.Millisecond * 100 * time.Duration(cB.Pop()))
	case 'u':
		cB.deepSea = true
	case 'F':
		var err error
		count := int(cB.Pop())
		bData := cB.stacks[cB.p].getBytes(count)
		if cB.file != nil {
			cB.file.Close()
			err = ioutil.WriteFile(cB.file.Name(), bData, os.ModePerm)
			if err != nil {
				panic(err)
			}
			cB.file = nil
		} else {
			fName := string(bData)
			cB.file, err = os.Open(fName)
			if err != nil {
				cB.file, err = os.Create(fName)
				if err != nil {
					panic(err)
				}
			}
		}
	case 'C':
		cB.Call()
	case 'R':
		cB.Ret()
	case 'I':
		cB.p++
	case 'D':
		cB.p--
	}
	return output, false
}

// Move changes the fish's x/y coordinates based on CodeBox.fDir.
func (cB *CodeBox) Move() {
	switch cB.fDir {
	case Right:
		cB.fX++
		if cB.fX >= cB.width {
			cB.fX = 0
		}
	case Down:
		cB.fY++
		if cB.fY >= cB.height {
			cB.fY = 0
		}
	case Left:
		cB.fX--
		if cB.fX < 0 {
			cB.fX = cB.width - 1
		}
	case Up:
		cB.fY--
		if cB.fY < 0 {
			cB.fY = cB.height - 1
		}
	}
}

// Swim causes the ><> to execute an instruction, then move. It returns a string of non-zero length when it has output and true when it encounters ";".
func (cB *CodeBox) Swim() (string, bool) {
	defer func() {
		if r := recover(); r != nil {
			cB.PrintBox()
			fmt.Println("Stack:", cB.Stack())
			fmt.Println("something smells fishy...")
			os.Exit(1)
		}
	}()

	var (
		output string
		end    bool
	)

	if r := cB.box[cB.fY][cB.fX]; cB.stringMode != 0 && r != cB.stringMode {
		cB.Push(float64(r))
	} else {
		output, end = cB.Exe(r)
	}
	cB.Move()
	return output, end
}

// Stack returns the underlying Stack slice.
func (cB *CodeBox) Stack() []float64 {
	if cB.p >= 0 && cB.p < len(cB.stacks) {
		return cB.stacks[cB.p].S
	} else {
		return []float64{float64(cB.p)}
	}
}

// Push appends r to the end of the current stack.
func (cB *CodeBox) Push(r float64) {
	cB.stacks[cB.p].Push(r)
}

// Pop removes the value on the end of the current stack and returns it.
func (cB *CodeBox) Pop() float64 {
	return cB.stacks[cB.p].Pop()
}

// StackLength implements "l" on the current stack.
func (cB *CodeBox) StackLength() float64 {
	return float64(len(cB.stacks[cB.p].S))
}

// Register implements "&" on the current stack.
func (cB *CodeBox) Register() {
	cB.stacks[cB.p].Register()
}

// ReverseStack implements "r" on the current stack.
func (cB *CodeBox) ReverseStack() {
	cB.stacks[cB.p].Reverse()
}

// ExtendStack implements ":" on the current stack.
func (cB *CodeBox) ExtendStack() {
	cB.stacks[cB.p].Extend()
}

// StackSwapTwo implements "$" on the current stack.
func (cB *CodeBox) StackSwapTwo() {
	cB.stacks[cB.p].SwapTwo()
}

// StackSwapThree implements "@" on the current stack.
func (cB *CodeBox) StackSwapThree() {
	cB.stacks[cB.p].SwapThree()
}

// StackShiftRight implements "}" on the current stack.
func (cB *CodeBox) StackShiftRight() {
	cB.stacks[cB.p].ShiftRight()
}

// StackShiftLeft implements "{" on the current stack.
func (cB *CodeBox) StackShiftLeft() {
	cB.stacks[cB.p].ShiftLeft()
}

// CloseStack implements "]".
func (cB *CodeBox) CloseStack() {
	cB.p--
	if cB.compMode {
		cB.stacks[cB.p+1].Reverse() // This is done to match the fishlanguage.com interpreter...
	}
	cB.stacks[cB.p].S = append(cB.stacks[cB.p].S, cB.stacks[cB.p+1].S...)
	if cB.p+2 == len(cB.stacks) {
		cB.stacks = cB.stacks[:cB.p+1]
	} else {
		cB.stacks = append(cB.stacks[:cB.p+1], cB.stacks[cB.p+2:]...)
	}
}

// NewStack implements "[".
func (cB *CodeBox) NewStack(n int) {
	cB.p++
	if cB.p == len(cB.stacks) {
		cB.stacks = append(cB.stacks, NewStack(cB.stacks[cB.p-1].S[len(cB.stacks[cB.p-1].S)-n:]))
	} else {
		tstacks := make([]*Stack, cB.p+1, len(cB.stacks)+1)
		copy(tstacks, cB.stacks[:cB.p])
		tstacks[cB.p] = NewStack(cB.stacks[cB.p-1].S[len(cB.stacks[cB.p-1].S)-n:])
		tstacks = append(tstacks, cB.stacks[cB.p:]...)
		cB.stacks = tstacks
	}
	cB.stacks[cB.p-1].S = cB.stacks[cB.p-1].S[:len(cB.stacks[cB.p-1].S)-n]
	if cB.compMode {
		cB.stacks[cB.p].Reverse() // This is done to match the fishlanguage.com interpreter...
	}
}

// Call implements "C".
func (cB *CodeBox) Call() {
	cB.p++
	if cB.p == len(cB.stacks) {
		cB.stacks = append(cB.stacks, NewStack([]float64{float64(cB.fX), float64(cB.fY)}))
		cB.stacks[cB.p], cB.stacks[cB.p-1] = cB.stacks[cB.p-1], cB.stacks[cB.p]
	} else {
		tstacks := make([]*Stack, cB.p+1, len(cB.stacks)+1)
		copy(tstacks, cB.stacks[:cB.p])
		tstacks[cB.p] = tstacks[cB.p-1]
		tstacks[cB.p-1] = NewStack([]float64{float64(cB.fX), float64(cB.fY)})
		tstacks = append(tstacks, cB.stacks[cB.p:]...)
		cB.stacks = tstacks
	}
	cB.fY = int(cB.Pop())
	cB.fX = int(cB.Pop())
}

// Ret implements "R".
func (cB *CodeBox) Ret() {
	cB.p--
	cB.fY = int(cB.Pop())
	cB.fX = int(cB.Pop())
	cB.stacks[cB.p] = cB.stacks[cB.p+1]
	if cB.p+2 == len(cB.stacks) {
		cB.stacks = cB.stacks[:cB.p+1]
	} else {
		cB.stacks = append(cB.stacks[:cB.p+1], cB.stacks[cB.p+2:]...)
	}
}

// PrintBox outputs the codebox to stdout.
func (cB *CodeBox) PrintBox() {
	fmt.Println()
	for y, line := range cB.box {
		for x, r := range line {
			if x != cB.fX || y != cB.fY {
				fmt.Print(" " + string(rune(r)) + " ")
			} else {
				fmt.Print("*" + string(rune(r)) + "*")
			}
		}
		fmt.Println()
	}
}

// Size returns the CodeBox's width/height
func (cB *CodeBox) Size() (int, int) {
	return cB.width, cB.height
}

// Loc returns the CodeBox's x/y
func (cB *CodeBox) Loc() (int, int) {
	return cB.fX, cB.fY
}

// Box returns a copy of the 2D slice containing the *><> script
func (cB *CodeBox) Box() [][]byte {
	outBox := make([][]byte, len(cB.box))
	for i := range cB.box {
		outBox[i] = make([]byte, len(cB.box[i]))
		copy(outBox[i], cB.box[i])
	}
	return outBox
}

// DeepSea returns whether the codebox is in DeepSea mode
func (cB *CodeBox) DeepSea() bool {
	return cB.deepSea
}

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
	reader = make(chan byte, 1024)
	go func() {
		var err error
		b := make([]byte, 1024)
		for err == nil {
			n, err := os.Stdin.Read(b)
			if err == nil {
				for i := 0; i < n; i++ {
					reader <- b[i]
				}
			} else {
				fmt.Println(err)
				return
			}
		}
	}()
}

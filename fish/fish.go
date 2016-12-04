package fish

import (
	"fmt"
	"os"
	"math/rand"
	"strings"
	"time"
)

type Direction byte

const (
	Right Direction = iota
	Down
	Left
	Up
)

var reader chan byte

type Stack struct {
	S              []float64
	register       float64
	filledRegister bool
}

func NewStack(s []float64) *Stack {
	return &Stack{S: s}
}

func (s *Stack) Register() {
	if s.filledRegister {
		s.Push(s.register)
		s.filledRegister = false
	} else {
		s.register = s.Pop()
		s.filledRegister = true
	}
}

func (s *Stack) Extend() {
	s.Push(s.S[len(s.S)-1])
}

// r
func (s *Stack) Reverse() {
	newS := make([]float64, len(s.S))
	for i, ii := 0, len(s.S)-1; ii >= 0; i, ii = i+1, ii-1 {
		newS[i] = s.S[ii]
	}
	s.S = newS
}

func (s *Stack) SwapTwo() {
	x := s.S[len(s.S)-1]
	s.S[len(s.S)-1] = s.S[len(s.S)-2]
	s.S[len(s.S)-2] = x
}

// 1,2,3,4, calling @ results in 1,4,2,3
func (s *Stack) SwapThree() {
	x := s.S[len(s.S)-1]
	y := s.S[len(s.S)-2]
	s.S[len(s.S)-1] = y
	s.S[len(s.S)-2] = s.S[len(s.S)-3]
	s.S[len(s.S)-3] = x
}

// }
func (s *Stack) ShiftRight() {
	newS := make([]float64, 1, len(s.S))
	newS[0] = s.Pop()
	s.S = append(newS, s.S...)
}

// {
func (s *Stack) ShiftLeft() {
	r := s.S[0]
	s.S = s.S[1:]
	s.Push(r)
}

func (s *Stack) Push(r float64) {
	s.S = append(s.S, float64(r))
}

func (s *Stack) Pop() (r float64) {
	if len(s.S) > 0 {
		r = s.S[len(s.S)-1]
		s.S = s.S[:len(s.S)-1]
	} else {
		panic("Stack is empty!")
	}
	return
}

func longestLineLength(lines []string) (l int) {
	for _, s := range lines {
		if len(s) > l {
			l = len(s)
		}
	}
	return
}

type CodeBox struct {
	Fx, Fy        int
	FDir          Direction
	Width, Height int
	Box           [][]byte
	stacks        []*Stack
	p             int // Used to keep track of the current stack
	StringMode    byte
	compMode      bool
}

func NewCodeBox(script string, stack []float64, compatibilityMode bool) *CodeBox {
	cB := new(CodeBox)

	script = strings.Replace(script, "\r", "", -1)
	if len(script) == 0 || script == "\n" {
		panic("Cannot accept script of length 0 (No room for the fish to survive).")
	}

	lines := strings.Split(script, "\n")
	cB.Width = longestLineLength(lines)
	cB.Height = len(lines)

	cB.Box = make([][]byte, cB.Height)
	for i, s := range lines {
		cB.Box[i] = make([]byte, cB.Width)
		for ii, r := 0, byte(0); ii < cB.Width; ii++ {
			if ii < len(s) {
				r = byte(s[ii])
			} else {
				r = ' '
			}
			cB.Box[i][ii] = byte(r)
		}
	}

	cB.stacks = []*Stack{NewStack(stack)}
	cB.compMode = compatibilityMode

	return cB
}

func (cB *CodeBox) Exe(r byte) bool {
	switch r {
	default:
		panic(r)
	case ' ':
	case ';':
		return true
	case '>':
		cB.FDir = Right
	case 'v':
		cB.FDir = Down
	case '<':
		cB.FDir = Left
	case '^':
		cB.FDir = Up
	case '|':
		if cB.FDir == Right {
			cB.FDir = Left
		} else if cB.FDir == Left {
			cB.FDir = Right
		}
	case '_':
		if cB.FDir == Down {
			cB.FDir = Up
		} else if cB.FDir == Up {
			cB.FDir = Down
		}
	case '#':
		switch cB.FDir {
		case Right:
			cB.FDir = Left
		case Down:
			cB.FDir = Up
		case Left:
			cB.FDir = Right
		case Up:
			cB.FDir = Down
		}
	case '/':
		switch cB.FDir {
		case Right:
			cB.FDir = Up
		case Down:
			cB.FDir = Left
		case Left:
			cB.FDir = Down
		case Up:
			cB.FDir = Right
		}
	case '\\':
		switch cB.FDir {
		case Right:
			cB.FDir = Down
		case Down:
			cB.FDir = Right
		case Left:
			cB.FDir = Up
		case Up:
			cB.FDir = Left
		}
	case 'x':
		cB.FDir = Direction(rand.Int31n(4))
	case '"', '\'':
		if cB.StringMode == 0 {
			cB.StringMode = r
		} else if r == cB.StringMode {
			cB.StringMode = 0
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		cB.Push(float64(r - '0'))
	case 'a', 'b', 'c', 'd', 'e', 'f':
		cB.Push(float64(r - 'a' + 10))
	case '&':
		cB.Register()
	case 'o':
		print(string(byte(cB.Pop())))
	case 'n':
		fmt.Printf("%v", cB.Pop())
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
		cB.Fy = int(cB.Pop())
		cB.Fx = int(cB.Pop())
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
		cB.Push(float64(cB.Box[int(cB.Pop())][int(cB.Pop())]))
	case 'p':
		cB.Box[int(cB.Pop())][int(cB.Pop())] = byte(cB.Pop())
	case 'i':
		r := float64(-1)
		b := byte(0)
		select {
		case b = <-reader:
			r = float64(b)
		default:
		}
		cB.Push(r)
	}
	return false
}

func (cB *CodeBox) Move() {
	switch cB.FDir {
	case Right:
		cB.Fx++
		if cB.Fx >= cB.Width {
			cB.Fx = 0
		}
	case Down:
		cB.Fy++
		if cB.Fy >= cB.Height {
			cB.Fy = 0
		}
	case Left:
		cB.Fx--
		if cB.Fx < 0 {
			cB.Fx = cB.Width - 1
		}
	case Up:
		cB.Fy--
		if cB.Fy < 0 {
			cB.Fy = cB.Height - 1
		}
	}
}

func (cB *CodeBox) Swim() bool {
	defer func() {
		if r := recover(); r != nil {
			cB.PrintBox()
			println("something smells fishy...")
			os.Exit(1)
		}
	}()

	if r := cB.Box[cB.Fy][cB.Fx]; cB.StringMode != 0 && r != cB.StringMode {
		cB.Push(float64(r))
	} else if cB.Exe(r) {
		return true
	}
	cB.Move()
	return false
}

func (cB *CodeBox) Stack() []float64 {
	return cB.stacks[cB.p].S
}

func (cB *CodeBox) Push(r float64) {
	cB.stacks[cB.p].Push(r)
}

func (cB *CodeBox) Pop() float64 {
	return cB.stacks[cB.p].Pop()
}

// l
func (cB *CodeBox) StackLength() float64 {
	return float64(len(cB.stacks[cB.p].S))
}

// &
func (cB *CodeBox) Register() {
	cB.stacks[cB.p].Register()
}

// r
func (cB *CodeBox) ReverseStack() {
	cB.stacks[cB.p].Reverse()
}

// :
func (cB *CodeBox) ExtendStack() {
	cB.stacks[cB.p].Extend()
}

// $
func (cB *CodeBox) StackSwapTwo() {
	cB.stacks[cB.p].SwapTwo()
}

// @
func (cB *CodeBox) StackSwapThree() {
	cB.stacks[cB.p].SwapThree()
}

// }
func (cB *CodeBox) StackShiftRight() {
	cB.stacks[cB.p].ShiftRight()
}

// {
func (cB *CodeBox) StackShiftLeft() {
	cB.stacks[cB.p].ShiftLeft()
}

// ]
func (cB *CodeBox) CloseStack() {
	cB.p--
	if cB.compMode {
		cB.stacks[cB.p+1].Reverse() // This is done to match the fishlanguage.com interpreter...
	}
	cB.stacks[cB.p].S = append(cB.stacks[cB.p].S, cB.stacks[cB.p+1].S...)
}

// [
func (cB *CodeBox) NewStack(n int) {
	cB.p++
	if cB.p == len(cB.stacks) {
		cB.stacks = append(cB.stacks, NewStack(cB.stacks[cB.p-1].S[len(cB.stacks[cB.p-1].S)-n:]))
		cB.stacks[cB.p-1].S = cB.stacks[cB.p-1].S[:len(cB.stacks[cB.p-1].S)-n]
	} else {
		cB.stacks[cB.p].S = cB.stacks[cB.p-1].S[len(cB.stacks[cB.p-1].S)-n:]
		cB.stacks[cB.p].filledRegister = false
	}
	if cB.compMode {
		cB.stacks[cB.p].Reverse() // This is done to match the fishlanguage.com interpreter...
	}
}

func (cB *CodeBox) PrintBox() {
	println()
	for y, line := range cB.Box {
		for x, r := range line {
			if x != cB.Fx || y != cB.Fy {
				print(" " + string(rune(r)) + " ")
			} else {
				print("*" + string(rune(r)) + "*")
			}
		}
		println()
	}
}

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
	reader = make(chan byte)
	go func() {
		var err error
		b := make([]byte, 1)
		for err == nil {
			_, err = os.Stdin.Read(b)
			if err == nil {
				reader <-b[0]
			}
		}
	}()
}

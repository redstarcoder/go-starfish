package starfish

import (
	"log"
	"testing"
	"time"
)

const (
	TESTVALUE1 = 1
	TESTVALUE2 = 2
	TESTVALUE3 = 3
	TESTVALUE4 = 4
	SCRIPT     = `r>l5(?v~~~/:!|Ou+1Ox:@=?~~~~~~~!
~~l5(?v" "/
 ~;!?l<`  // Script used in "BenchmarkScript"
)

var (
	INITIALSTACK = []float64{float64('h'), float64('e'), float64('l'), float64('l'), float64('o'),
		float64(' '), float64('w'), float64('o'), float64('r'), float64('l'), float64('d')} // Stack used in "BenchmarkScript"
)

func runscript(script string, initialstack []float64, compMode bool) *CodeBox {
	cB := NewCodeBox(script, initialstack, compMode)
	now := time.Now()
	for !cB.Swim() {
		if time.Since(now) >= time.Second {
			log.Fatalln("script taking too long...")
		}
	}
	return cB
}

func BenchmarkScript(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		stack := make([]float64, len(INITIALSTACK))
		copy(stack, INITIALSTACK)
		cB := NewCodeBox(SCRIPT, stack, false)
		b.StartTimer()
		for !cB.Swim() {
		}
	}
	log.Println(b.N)
}

func TestStackRegister(t *testing.T) {
	cB := runscript("&;", []float64{TESTVALUE1, TESTVALUE2, TESTVALUE3}, false)
	s := cB.stacks[0]
	if len(s.S) != 2 || s.register != TESTVALUE3 || s.S[0] != TESTVALUE1 || !s.filledRegister {
		t.FailNow()
	}
	s.Register()
	if len(s.S) != 3 || s.S[0] != TESTVALUE1 || s.S[2] != TESTVALUE3 || s.filledRegister {
		t.FailNow()
	}
}

func TestStackExtend(t *testing.T) {
	cB := runscript(":;", []float64{TESTVALUE1, TESTVALUE2}, false)
	s := cB.stacks[0]
	if len(s.S) != 3 || s.S[2] != TESTVALUE2 {
		t.FailNow()
	}
}

func TestStackReverse(t *testing.T) {
	cB := runscript("r;", []float64{TESTVALUE1, TESTVALUE2, TESTVALUE3}, false)
	s := cB.stacks[0]
	if s.S[0] != TESTVALUE3 || s.S[1] != TESTVALUE2 || s.S[2] != TESTVALUE1 {
		t.FailNow()
	}
}

func TestStackSwapTwo(t *testing.T) {
	cB := runscript("$;", []float64{TESTVALUE1, TESTVALUE2, TESTVALUE3}, false)
	s := cB.stacks[0]
	if s.S[0] != TESTVALUE1 || s.S[1] != TESTVALUE3 || s.S[2] != TESTVALUE2 {
		t.FailNow()
	}
}

func TestStackSwapThree(t *testing.T) {
	cB := runscript("@;", []float64{TESTVALUE1, TESTVALUE2, TESTVALUE3, TESTVALUE4}, false)
	s := cB.stacks[0]
	if s.S[0] != TESTVALUE1 || s.S[1] != TESTVALUE4 || s.S[2] != TESTVALUE2 || s.S[3] != TESTVALUE3 {
		t.FailNow()
	}
}

func TestStackShiftLeft(t *testing.T) {
	cB := runscript("{;", []float64{TESTVALUE1, TESTVALUE2, TESTVALUE3, TESTVALUE4}, false)
	s := cB.stacks[0]
	if s.S[0] != TESTVALUE2 || s.S[1] != TESTVALUE3 || s.S[2] != TESTVALUE4 || s.S[3] != TESTVALUE1 {
		t.FailNow()
	}
}

func TestStackShiftRight(t *testing.T) {
	cB := runscript("};", []float64{TESTVALUE1, TESTVALUE2, TESTVALUE3, TESTVALUE4}, false)
	s := cB.stacks[0]
	if s.S[0] != TESTVALUE4 || s.S[1] != TESTVALUE1 || s.S[2] != TESTVALUE2 || s.S[3] != TESTVALUE3 {
		t.FailNow()
	}
}

func TestNewStackCloseStack(t *testing.T) {
	cB := NewCodeBox("[]", []float64{TESTVALUE1, TESTVALUE2, TESTVALUE3, TESTVALUE4, 2}, false)
	cB.Swim()
	s := cB.stacks[0]
	s2 := cB.stacks[1]
	if s.S[0] != TESTVALUE1 || s.S[1] != TESTVALUE2 || s2.S[0] != TESTVALUE3 || s2.S[1] != TESTVALUE4 || len(s.S) != 2 || len(s2.S) != 2 {
		t.FailNow()
	}

	cB.Swim()
	s = cB.stacks[0]
	if s.S[0] != TESTVALUE1 || s.S[1] != TESTVALUE2 || s.S[2] != TESTVALUE3 || s.S[3] != TESTVALUE4 || len(s.S) != 4 {
		t.FailNow()
	}
}

func TestNewStackCloseStackCompatibility(t *testing.T) {
	cB := NewCodeBox("[]", []float64{TESTVALUE1, TESTVALUE2, TESTVALUE3, TESTVALUE4, 2}, true)
	cB.Swim()
	s := cB.stacks[0]
	s2 := cB.stacks[1]
	if s.S[0] != TESTVALUE1 || s.S[1] != TESTVALUE2 || s2.S[1] != TESTVALUE3 || s2.S[0] != TESTVALUE4 || len(s.S) != 2 || len(s2.S) != 2 {
		t.FailNow()
	}

	cB.Swim()
	s = cB.stacks[0]
	if s.S[0] != TESTVALUE1 || s.S[1] != TESTVALUE2 || s.S[2] != TESTVALUE3 || s.S[3] != TESTVALUE4 || len(s.S) != 4 {
		t.FailNow()
	}
}

func TestPrintBox(t *testing.T) {
	cB := NewCodeBox(`"Hello test!";`, []float64{}, false)
	cB.PrintBox()
}

func TestStackLength(t *testing.T) {
	cB := NewCodeBox("l;", []float64{TESTVALUE1, TESTVALUE2, TESTVALUE3}, false)
	cB.Swim()
	if cB.Stack()[3] != 3 {
		t.Fail()
	}
}

func TestStackReturn(t *testing.T) {
	cB := NewCodeBox(";", []float64{TESTVALUE1, TESTVALUE3}, false)
	s := cB.Stack()
	if s[0] != TESTVALUE1 || s[1] != TESTVALUE3 {
		t.Fail()
	}
}

func TestMovement(t *testing.T) {
	cB := NewCodeBox(">;", []float64{}, false)
	cB.Swim()
	if !cB.Swim() {
		t.Fail()
	}

	cB = NewCodeBox("<;", []float64{}, false)
	cB.Swim()
	if !cB.Swim() {
		t.Fail()
	}

	cB = NewCodeBox("^\n;", []float64{}, false)
	cB.Swim()
	if !cB.Swim() {
		t.Fail()
	}

	cB = NewCodeBox("v\n;", []float64{}, false)
	cB.Swim()
	if !cB.Swim() {
		t.Fail()
	}

	cB = NewCodeBox("`;\n`", []float64{}, false)
	for i := 0; i < 5; i++ {
		if cB.Swim() {
			t.Fail()
		}
	}
	if !cB.Swim() {
		t.Fail()
	}
}

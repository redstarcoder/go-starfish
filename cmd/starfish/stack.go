package main

import (
	"errors"
	"strconv"
)

type stack struct {
	s []float64
}

func (s *stack) String() string {
	return ""
}

func (s *stack) Set(str string) error {
	var strMode byte
	runes := make([]rune, 0, 32)
	s.s = make([]float64, 0, 32)
	for _, r := range str {
		if strMode != 0 && byte(r) != strMode {
			s.s = append(s.s, float64(r))
			continue
		}
		switch r {
		default:
			return errors.New("Invalid initial stack")
		case ' ':
			if len(runes) > 0 {
				if f, err := strconv.ParseFloat(string(runes), 64); err == nil {
					s.s = append(s.s, f)
					runes = make([]rune, 0, 32)
				} else {
					return err
				}
			}
		case '\'', '"':
			if strMode == 0 {
				strMode = byte(r)
			} else {
				strMode = 0
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
			runes = append(runes, r)
		}
	}
	if f, err := strconv.ParseFloat(string(runes), 64); err == nil {
		s.s = append(s.s, f)
		runes = make([]rune, 0, 32)
	} else if len(runes) > 0 {
		return err
	}
	return nil
}

func (s *stack) Get() interface{} {
	return s.s
}

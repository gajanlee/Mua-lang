package vm

import (
	"muc/code"
	"muc/object"
)

type Frame struct {
	// fn *object.CompiledFunction
	cl *object.Closure
	ip int
	basePointer int
}

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	f := &Frame{
		cl: cl,
		ip: -1,
		basePointer: basePointer,	// the bottom of current stack
	}

	return f
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
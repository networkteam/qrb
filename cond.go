package jrm

import (
	"strings"
)

type binOp string

const (
	opEqual binOp = "="
)

type binaryExp struct {
	lft Exp
	rgt Exp
	op  binOp
}

func (b binaryExp) isExp() {}

func (b binaryExp) WriteSQL(sb *strings.Builder) (args []any) {
	args = append(args, b.lft.WriteSQL(sb)...)
	sb.WriteRune(' ')
	sb.WriteString(string(b.op))
	sb.WriteRune(' ')
	args = append(args, b.rgt.WriteSQL(sb)...)
	return args
}

func Eq(lft Exp, rgt Exp) Exp {
	return &binaryExp{
		lft: lft,
		rgt: rgt,
		op:  opEqual,
	}
}

type conjunctionExp struct {
	exps []Exp
}

func (c conjunctionExp) isExp() {}

func (c conjunctionExp) WriteSQL(sb *strings.Builder) (args []any) {
	for i, exp := range c.exps {
		if i > 0 {
			sb.WriteString(" AND ")
		}
		sb.WriteRune('(')
		expArgs := exp.WriteSQL(sb)
		sb.WriteRune(')')
		args = append(args, expArgs...)
	}
	return args
}

func And(exps ...Exp) Exp {
	return &conjunctionExp{
		exps: exps,
	}
}

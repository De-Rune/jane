// Copyright (c) 2024 - DeRuneLabs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package transpiler

import (
	"github.com/DeRuneLabs/jane/lexer"
	"github.com/DeRuneLabs/jane/lexer/tokens"
	"github.com/DeRuneLabs/jane/package/jnbits"
	"github.com/DeRuneLabs/jane/package/jntype"
)

func setshift(v *value, right uint64) {
	switch {
	case right <= 6:
		v.data.Type.Id = jntype.I8
	case right <= 7:
		v.data.Type.Id = jntype.U8
	case right <= 14:
		v.data.Type.Id = jntype.I16
	case right <= 15:
		v.data.Type.Id = jntype.U16
	case right <= 30:
		v.data.Type.Id = jntype.I32
	case right <= 31:
		v.data.Type.Id = jntype.U32
	case right <= 62:
		v.data.Type.Id = jntype.I64
	case right <= 63:
		v.data.Type.Id = jntype.U64
	case right <= 127:
		v.data.Type.Id = jntype.F32
	default:
		v.data.Type.Id = jntype.F64
	}
}

func bitize(v *value) {
	switch t := v.expr.(type) {
	case float64:
		v.data.Type.Id = jntype.FloatFromBits(jnbits.BitsizeFloat(t))
	case int64:
		v.data.Type.Id = jntype.IntFromBits(jnbits.BitsizeInt(t))
	case uint64:
		v.data.Type.Id = jntype.UIntFromBits(jnbits.BitsizeUInt(t))
	default:
		return
	}
	v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
}

func tonumf(expr any) float64 {
	switch t := expr.(type) {
	case float64:
		return t
	case int64:
		return float64(t)
	case uint64:
		return float64(t)
	}
	return 0
}

func tonumu(expr any) uint64 {
	switch t := expr.(type) {
	case float64:
		return uint64(t)
	case int64:
		return uint64(t)
	case uint64:
		return t
	}
	return 0
}

func tonums(expr any) int64 {
	switch t := expr.(type) {
	case float64:
		return int64(t)
	case int64:
		return t
	case uint64:
		return int64(t)
	}
	return 0
}

type solver struct {
	t  *Transpiler
	l  value
	r  value
	op lexer.Token
}

func (s *solver) eq(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case bool:
		v.expr = left == s.r.expr.(bool)
	case string:
		v.expr = left == s.r.expr.(string)
	case float64:
		v.expr = left == tonumf(s.r.expr)
	case int64:
		v.expr = left == tonums(s.r.expr)
	case uint64:
		v.expr = left == tonumu(s.r.expr)
	}
}

func (s *solver) noteq(v *value) {
	if !s.isConstExpr() {
		return
	}
	s.eq(v)
	v.expr = !v.expr.(bool)
}

func (s *solver) lt(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case float64:
		v.expr = left < tonumf(s.r.expr)
	case int64:
		v.expr = left < tonums(s.r.expr)
	case uint64:
		v.expr = left < tonumu(s.r.expr)
	}
}

func (s *solver) gt(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case float64:
		v.expr = left > tonumf(s.r.expr)
	case int64:
		v.expr = left > tonums(s.r.expr)
	case uint64:
		v.expr = left > tonumu(s.r.expr)
	}
}

func (s *solver) lteq(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case float64:
		v.expr = left <= tonumf(s.r.expr)
	case int64:
		v.expr = left <= tonums(s.r.expr)
	case uint64:
		v.expr = left <= tonumu(s.r.expr)
	}
}

func (s *solver) gteq(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case float64:
		v.expr = left >= tonumf(s.r.expr)
	case int64:
		v.expr = left >= tonums(s.r.expr)
	case uint64:
		v.expr = left >= tonumu(s.r.expr)
	}
}

func (s *solver) add(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case string:
		v.expr = left + s.r.expr.(string)
	case float64:
		v.expr = left + tonumf(s.r.expr)
	case int64:
		v.expr = int64(left + tonums(s.r.expr))
	case uint64:
		v.expr = uint64(left + tonumu(s.r.expr))
	}
}

func (s *solver) sub(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case float64:
		v.expr = left - tonumf(s.r.expr)
	case int64:
		v.expr = int64(left - tonums(s.r.expr))
	case uint64:
		v.expr = uint64(left - tonumu(s.r.expr))
	}
}

func (s *solver) mul(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case float64:
		v.expr = left * tonumf(s.r.expr)
	case int64:
		v.expr = int64(left * tonums(s.r.expr))
	case uint64:
		v.expr = uint64(left * tonumu(s.r.expr))
	}
}

func (s *solver) div(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case float64:
		right := tonumf(s.r.expr)
		if right != 0 {
			v.expr = left / right
		} else {
			s.t.pusherrtok(s.op, "divide_by_zero")
			v.expr = float64(0)
		}
	case int64:
		right := tonumf(s.r.expr)
		if right != 0 {
			v.expr = float64(left) / right
		} else {
			s.t.pusherrtok(s.op, "divide_by_zero")
			v.expr = int64(0)
		}
	case uint64:
		right := tonumf(s.r.expr)
		if right != 0 {
			v.expr = float64(left) / right
		} else {
			s.t.pusherrtok(s.op, "divide_by_zero")
			v.expr = float64(0)
		}
	}
}

func (s *solver) mod(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case int64:
		right := tonums(s.r.expr)
		if right != 0 {
			v.expr = left % right
		} else {
			s.t.pusherrtok(s.op, "divide_by_zero")
			v.expr = int64(0)
		}
	case uint64:
		right := tonumu(s.r.expr)
		if right != 0 {
			v.expr = left % right
		} else {
			s.t.pusherrtok(s.op, "divide_by_zero")
			v.expr = uint64(0)
		}
	}
}

func (s *solver) bitwiseAnd(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case int64:
		v.expr = left & tonums(s.r.expr)
	case uint64:
		v.expr = left & tonumu(s.r.expr)
	}
}

func (s *solver) bitwiseOr(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case int64:
		v.expr = left | tonums(s.r.expr)
	case uint64:
		v.expr = left | tonumu(s.r.expr)
	}
}

func (s *solver) bitwiseXor(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case int64:
		v.expr = left ^ tonums(s.r.expr)
	case uint64:
		v.expr = left ^ tonumu(s.r.expr)
	}
}

func (s *solver) urshift(v *value) {
	left := tonumu(s.l.expr)
	right := tonumu(s.r.expr)
	v.expr = left >> right
	setshift(v, right)
}

func (s *solver) rshift(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case int64:
		if left < 0 {
			right := tonumu(s.r.expr)
			v.expr = left >> right
			setshift(v, right)
		} else {
			s.urshift(v)
		}
	case uint64:
		s.urshift(v)
	}
}

func (s *solver) ulshift(v *value) {
	left := tonumu(s.l.expr)
	right := tonumu(s.r.expr)
	v.expr = left << right
	setshift(v, right)
}

func (s *solver) lshift(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case int64:
		if left < 0 {
			right := tonumu(s.r.expr)
			v.expr = left << right
			setshift(v, right)
		} else {
			s.ulshift(v)
		}
	case uint64:
		s.ulshift(v)
	}
}

func (s *solver) and(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case bool:
		v.expr = left && s.r.expr.(bool)
	}
}

func (s *solver) or(v *value) {
	if !s.isConstExpr() {
		return
	}
	switch left := s.l.expr.(type) {
	case bool:
		v.expr = left || s.r.expr.(bool)
	}
}

func (s *solver) ptr() (v value) {
	v.data.Token = s.op
	if !s.types_are_compatible(true) {
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "incompatible_types",
			s.r.data.Type.Kind, s.l.data.Type.Kind)
		return
	}
	if !typeIsPtr(s.l.data.Type) {
		s.l, s.r = s.r, s.l
	}
	switch s.op.Kind {
	case tokens.PLUS, tokens.MINUS:
		v.data.Type = s.l.data.Type
	case tokens.EQUALS, tokens.NOT_EQUALS, tokens.LESS, tokens.GREAT,
		tokens.GREAT_EQUAL, tokens.LESS_EQUAL:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
	default:
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "operator_not_for_jntype", s.op.Kind, "pointer")
	}
	return
}

func (s *solver) enum() (v value) {
	if typeIsEnum(s.l.data.Type) {
		s.l.data.Type = s.l.data.Type.Tag.(*Enum).Type
	}
	if typeIsEnum(s.r.data.Type) {
		s.r.data.Type = s.r.data.Type.Tag.(*Enum).Type
	}
	return s.solve()
}

func (s *solver) str() (v value) {
	v.data.Token = s.op
	if s.l.data.Type.Id != s.r.data.Type.Id {
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "incompatible_types",
			s.l.data.Type.Kind, s.r.data.Type.Kind)
		return
	}
	switch s.op.Kind {
	case tokens.PLUS:
		v.data.Type.Id = jntype.Str
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.add(&v)
	case tokens.EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.eq(&v)
	case tokens.NOT_EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.noteq(&v)
	default:
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "operator_not_for_jntype",
			s.op.Kind, tokens.STR)
	}
	return
}

func (s *solver) any() (v value) {
	v.data.Token = s.op
	switch s.op.Kind {
	case tokens.EQUALS, tokens.NOT_EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
	default:
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "operator_not_for_jntype", s.op.Kind, tokens.ANY)
	}
	return
}

func (s *solver) bool() (v value) {
	v.data.Token = s.op
	if !s.types_are_compatible(true) {
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "incompatible_types",
			s.r.data.Type.Kind, s.l.data.Type.Kind)
		return
	}
	switch s.op.Kind {
	case tokens.EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.eq(&v)
	case tokens.NOT_EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.noteq(&v)
	default:
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "operator_not_for_jntype",
			s.op.Kind, tokens.BOOL)
	}
	return
}

func (s *solver) floatMod() (v value, ok bool) {
	if !jntype.IsInteger(s.l.data.Type.Id) {
		if !jntype.IsInteger(s.r.data.Type.Id) {
			return
		}
		s.l, s.r = s.r, s.l
	}
	switch {
	case jntype.IsSignedInteger(s.l.data.Type.Id):
		switch {
		case integerAssignable(jntype.I64, s.r):
			return s.signed(), true
		case integerAssignable(jntype.U64, s.r):
			return s.unsigned(), true
		}
	case jntype.IsUnsignedInteger(s.l.data.Type.Id):
		if integerAssignable(jntype.I64, s.r) ||
			integerAssignable(jntype.U64, s.r) {
			return s.unsigned(), true
		}
	}
	return
}

func (s *solver) float() (v value) {
	v.data.Token = s.op
	if !jntype.IsNumeric(s.l.data.Type.Id) ||
		!jntype.IsNumeric(s.r.data.Type.Id) {
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "incompatible_types",
			s.r.data.Type.Kind, s.l.data.Type.Kind)
		return
	}
	switch s.op.Kind {
	case tokens.EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.eq(&v)
	case tokens.NOT_EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.noteq(&v)
	case tokens.LESS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.lt(&v)
	case tokens.GREAT:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.gt(&v)
	case tokens.GREAT_EQUAL:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.gteq(&v)
	case tokens.LESS_EQUAL:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.lteq(&v)
	case tokens.PLUS:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.add(&v)
	case tokens.MINUS:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.sub(&v)
	case tokens.STAR:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.mul(&v)
	case tokens.SOLIDUS:
		if jntype.IsInteger(s.l.data.Type.Id) && jntype.IsInteger(s.r.data.Type.Id) {
		} else if jntype.IsInteger(s.l.data.Type.Id) {
			s.r.data.Type = s.l.data.Type
		} else if jntype.IsInteger(s.r.data.Type.Id) {
			s.l.data.Type = s.r.data.Type
		}
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.div(&v)
	case tokens.PERCENT:
		var ok bool
		v, ok = s.floatMod()
		if ok {
			break
		}
		fallthrough
	default:
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "operator_not_for_float", s.op.Kind)
	}
	return
}

func (s *solver) signed() (v value) {
	v.data.Token = s.op
	if !jntype.IsNumeric(s.l.data.Type.Id) ||
		!jntype.IsNumeric(s.r.data.Type.Id) {
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "incompatible_types",
			s.r.data.Type.Kind, s.l.data.Type.Kind)
		return
	}
	switch s.op.Kind {
	case tokens.EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.eq(&v)
	case tokens.NOT_EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.noteq(&v)
	case tokens.LESS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.lt(&v)
	case tokens.GREAT:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.gt(&v)
	case tokens.GREAT_EQUAL:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.gteq(&v)
	case tokens.LESS_EQUAL:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.lteq(&v)
	case tokens.PLUS:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.add(&v)
	case tokens.MINUS:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.sub(&v)
	case tokens.STAR:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.mul(&v)
	case tokens.SOLIDUS:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.div(&v)
	case tokens.PERCENT:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.mod(&v)
	case tokens.AMPER:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.bitwiseAnd(&v)
	case tokens.VLINE:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.bitwiseOr(&v)
	case tokens.CARET:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.bitwiseXor(&v)
	case tokens.RSHIFT:
		v.data.Type.Id = jntype.U64
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		if !okForShifting(s.r) {
			s.t.pusherrtok(s.op, "bitshift_must_unsigned")
		}
		s.rshift(&v)
	case tokens.LSHIFT:
		v.data.Type.Id = jntype.U64
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		if !okForShifting(s.r) {
			s.t.pusherrtok(s.op, "bitshift_must_unsigned")
		}
		s.lshift(&v)
	default:
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "operator_not_for_int", s.op.Kind)
	}
	return
}

func (s *solver) unsigned() (v value) {
	v.data.Token = s.op
	if !jntype.IsNumeric(s.l.data.Type.Id) ||
		!jntype.IsNumeric(s.r.data.Type.Id) {
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "incompatible_types",
			s.r.data.Type.Kind, s.l.data.Type.Kind)
		return
	}
	switch s.op.Kind {
	case tokens.EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.eq(&v)
	case tokens.NOT_EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.noteq(&v)
	case tokens.LESS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.lt(&v)
	case tokens.GREAT:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.gt(&v)
	case tokens.GREAT_EQUAL:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.gteq(&v)
	case tokens.LESS_EQUAL:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		s.lteq(&v)
	case tokens.PLUS:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.add(&v)
	case tokens.MINUS:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.sub(&v)
	case tokens.STAR:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.mul(&v)
	case tokens.SOLIDUS:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.div(&v)
	case tokens.PERCENT:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.mod(&v)
	case tokens.AMPER:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.bitwiseAnd(&v)
	case tokens.VLINE:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.bitwiseOr(&v)
	case tokens.CARET:
		v.data.Type = s.l.data.Type
		if jntype.TypeGreaterThan(s.r.data.Type.Id, v.data.Type.Id) {
			v.data.Type = s.r.data.Type
		}
		s.bitwiseXor(&v)
	case tokens.RSHIFT:
		v.data.Type.Id = jntype.U64
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		if !okForShifting(s.r) {
			s.t.pusherrtok(s.op, "bitshift_must_unsigned")
		}
		s.rshift(&v)
	case tokens.LSHIFT:
		v.data.Type.Id = jntype.U64
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		if !okForShifting(s.r) {
			s.t.pusherrtok(s.op, "bitshift_must_unsigned")
		}
		s.lshift(&v)
	default:
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "operator_not_for_uint", s.op.Kind)
	}
	return
}

func (s *solver) logical() (v value) {
	if s.l.data.Type.Id != jntype.Bool ||
		s.r.data.Type.Id != jntype.Bool {
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "logical_not_bool")
		return
	}
	v.data.Token = s.op
	v.data.Type.Id = jntype.Bool
	v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
	if !s.isConstExpr() {
		return
	}
	switch s.op.Kind {
	case tokens.DOUBLE_AMPER:
		s.and(&v)
	case tokens.DOUBLE_VLINE:
		s.or(&v)
	}
	return
}

func (s *solver) array() (v value) {
	v.data.Token = s.op
	if !s.types_are_compatible(true) {
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "incompatible_types",
			s.r.data.Type.Kind, s.l.data.Type.Kind)
		return
	}
	switch s.op.Kind {
	case tokens.EQUALS, tokens.NOT_EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
	default:
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "operator_not_for_jntype", s.op.Kind, s.l.data.Type.Kind)
	}
	return
}

func (s *solver) slice() (v value) {
	v.data.Token = s.op
	if !s.types_are_compatible(true) {
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "incompatible_types",
			s.r.data.Type.Kind, s.l.data.Type.Kind)
		return
	}
	switch s.op.Kind {
	case tokens.EQUALS, tokens.NOT_EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
	default:
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "operator_not_for_jntype",
			s.op.Kind, s.l.data.Type.Kind)
	}
	return
}

func (s *solver) nil() (v value) {
	v.data.Token = s.op
	if !s.types_are_compatible(false) {
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "incompatible_types",
			s.r.data.Type.Kind, s.l.data.Type.Kind)
		return
	}
	switch s.op.Kind {
	case tokens.NOT_EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		if s.isConstExpr() {
			v.expr = s.l.expr != nil && s.r.expr != nil
		}
	case tokens.EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		if s.isConstExpr() {
			v.expr = s.l.expr == nil && s.r.expr == nil
		}
	default:
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "operator_not_for_jntype",
			s.op.Kind, tokens.NIL)
	}
	return
}

func (s *solver) structure() (v value) {
	v.data.Token = s.op
	if s.l.data.Type.Kind != s.r.data.Type.Kind {
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "incompatible_types",
			s.r.data.Type.Kind, s.l.data.Type.Kind)
		return
	}
	switch s.op.Kind {
	case tokens.NOT_EQUALS, tokens.EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
	default:
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "operator_not_for_jntype",
			s.op.Kind, tokens.STRUCT)
	}
	return
}

func (s *solver) jntrait() (v value) {
	v.data.Token = s.op
	if !s.types_are_compatible(true) {
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "incompatible_types",
			s.r.data.Type.Kind, s.l.data.Type.Kind)
		return
	}
	switch s.op.Kind {
	case tokens.NOT_EQUALS, tokens.EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
	default:
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "operator_not_for_jntype",
			s.op.Kind, tokens.TRAIT)
	}
	return
}

func (s *solver) function() (v value) {
	v.data.Token = s.op
	if (!typeIsPure(s.l.data.Type) || s.l.data.Type.Id != jntype.Nil) &&
		(!typeIsPure(s.r.data.Type) || s.r.data.Type.Id != jntype.Nil) {
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "incompatible_types",
			s.r.data.Type.Kind, s.l.data.Type.Kind)
		return
	}
	switch s.op.Kind {
	case tokens.NOT_EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
	case tokens.EQUALS:
		v.data.Type.Id = jntype.Bool
		v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
	default:
		s.t.eval.has_error = true
		s.t.pusherrtok(s.op, "operator_not_for_jntype",
			s.op.Kind, tokens.NIL)
	}
	return
}

func (s *solver) types_are_compatible(ignore_any bool) bool {
	checker := type_checker{
		left:         s.l.data.Type,
		right:        s.r.data.Type,
		ignore_any:   ignore_any,
		allow_assign: true,
	}
	ok := checker.check()
	return ok
}

func (s *solver) isConstExpr() bool {
	return s.l.constExpr && s.r.constExpr
}

func (s *solver) solve() (v value) {
	defer func() {
		if typeIsVoid(v.data.Type) {
			v.data.Type.Kind = jntype.TypeMap[v.data.Type.Id]
		} else {
			v.constExpr = s.isConstExpr()
			if v.constExpr {
				bitize(&v)
				v.model = getModel(v)
			}
		}
	}()
	switch s.op.Kind {
	case tokens.DOUBLE_AMPER, tokens.DOUBLE_VLINE:
		return s.logical()
	}
	switch {
	case typeIsFunc(s.l.data.Type), typeIsFunc(s.r.data.Type):
		return s.function()
	case typeIsArray(s.l.data.Type), typeIsArray(s.r.data.Type):
		return s.array()
	case typeIsSlice(s.l.data.Type), typeIsSlice(s.r.data.Type):
		return s.slice()
	case typeIsPtr(s.l.data.Type), typeIsPtr(s.r.data.Type):
		return s.ptr()
	case typeIsEnum(s.l.data.Type), typeIsEnum(s.r.data.Type):
		return s.enum()
	case typeIsStruct(s.l.data.Type), typeIsStruct(s.r.data.Type):
		return s.structure()
	case typeIsTrait(s.l.data.Type), typeIsTrait(s.r.data.Type):
		return s.jntrait()
	case s.l.data.Type.Id == jntype.Nil, s.r.data.Type.Id == jntype.Nil:
		return s.nil()
	case s.l.data.Type.Id == jntype.Any, s.r.data.Type.Id == jntype.Any:
		return s.any()
	case s.l.data.Type.Id == jntype.Bool, s.r.data.Type.Id == jntype.Bool:
		return s.bool()
	case s.l.data.Type.Id == jntype.Str, s.r.data.Type.Id == jntype.Str:
		return s.str()
	case jntype.IsFloat(s.l.data.Type.Id),
		jntype.IsFloat(s.r.data.Type.Id):
		return s.float()
	case jntype.IsUnsignedInteger(s.l.data.Type.Id),
		jntype.IsUnsignedInteger(s.r.data.Type.Id):
		return s.unsigned()
	case jntype.IsSignedNumeric(s.l.data.Type.Id),
		jntype.IsSignedNumeric(s.r.data.Type.Id):
		return s.signed()
	}
	return
}

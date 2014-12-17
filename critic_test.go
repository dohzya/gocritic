package gocritic

import (
	"bytes"
	"testing"
)

//
// {++  ++}
//

func TestIns(t *testing.T) {
	msg := "{++a++}"
	expt := "<ins>a</ins>"
	check(t, msg, expt)
}

func TestInsBegin(t *testing.T) {
	msg := "{++a++}b"
	expt := "<ins>a</ins>b"
	check(t, msg, expt)
}

func TestInsEnd(t *testing.T) {
	msg := "a{++b++}"
	expt := "a<ins>b</ins>"
	check(t, msg, expt)
}

func TestInsMiddle(t *testing.T) {
	msg := "a{++b++}c"
	expt := "a<ins>b</ins>c"
	check(t, msg, expt)
}

func TestInsBreak(t *testing.T) {
	msg := "a{++\nb++}c"
	expt := "a\n<ins class=\"break\">b</ins>c"
	check(t, msg, expt)
}

func TestInsNewParagraph(t *testing.T) {
	msg := "a{++\n++}b"
	expt := "a\n<ins class=\"break\">&nbsp;</ins>\nb"
	check(t, msg, expt)
}

func TestNoIns2(t *testing.T) {
	msg := "a++}b"
	expt := "a++}b"
	check(t, msg, expt)
}

//
// {--  --}
//

func TestDel(t *testing.T) {
	msg := "{--a--}"
	expt := "<del>a</del>"
	check(t, msg, expt)
}

func TestDelBegin(t *testing.T) {
	msg := "{--a--}b"
	expt := "<del>a</del>b"
	check(t, msg, expt)
}

func TestDelEnd(t *testing.T) {
	msg := "a{--b--}"
	expt := "a<del>b</del>"
	check(t, msg, expt)
}

func TestDelMiddle(t *testing.T) {
	msg := "a{--b--}c"
	expt := "a<del>b</del>c"
	check(t, msg, expt)
}

func TestNoDel1(t *testing.T) {
	msg := "{-a-}"
	expt := "{-a-}"
	check(t, msg, expt)
}

func TestNoDel2(t *testing.T) {
	msg := "a--}b"
	expt := "a--}b"
	check(t, msg, expt)
}

//
// {~~  ~>  ~~}
//

func TestSub(t *testing.T) {
	msg := "{~~a~>b~~}"
	expt := "<del>a</del><ins>b</ins>"
	check(t, msg, expt)
}

func TestSubBegin(t *testing.T) {
	msg := "{~~a~>b~~}c"
	expt := "<del>a</del><ins>b</ins>c"
	check(t, msg, expt)
}

func TestSubEnd(t *testing.T) {
	msg := "a{~~b~>c~~}"
	expt := "a<del>b</del><ins>c</ins>"
	check(t, msg, expt)
}

func TestSubMiddle(t *testing.T) {
	msg := "a{~~b~>c~~}d"
	expt := "a<del>b</del><ins>c</ins>d"
	check(t, msg, expt)
}

func TestSubBreak(t *testing.T) {
	msg := "a{~~b~>\nc~~}d"
	expt := "a<del>b</del>\n<ins class=\"break\">c</ins>d"
	check(t, msg, expt)
}

func TestSubNewParagraph(t *testing.T) {
	msg := "a{~~b~>\n~~}c"
	expt := "a<del>b</del>\n<ins class=\"break\">&nbsp;</ins>\nc"
	check(t, msg, expt)
}

func TestNoSub1(t *testing.T) {
	msg := "{~a~}"
	expt := "{~a~}"
	check(t, msg, expt)
}

func TestNoSub2(t *testing.T) {
	msg := "a~>b~~}c"
	expt := "a~>b~~}c"
	check(t, msg, expt)
}

func TestNoSub3(t *testing.T) {
	msg := "a~~}b"
	expt := "a~~}b"
	check(t, msg, expt)
}

//
// {==  ==}
//

func TestHighlight(t *testing.T) {
	msg := "{==a==}"
	expt := "<mark>a</mark>"
	check(t, msg, expt)
}

func TestHighlightBegin(t *testing.T) {
	msg := "{==a==}b"
	expt := "<mark>a</mark>b"
	check(t, msg, expt)
}

func TestHighlightEnd(t *testing.T) {
	msg := "a{==b==}"
	expt := "a<mark>b</mark>"
	check(t, msg, expt)
}

func TestHighlightMiddle(t *testing.T) {
	msg := "a{==b==}c"
	expt := "a<mark>b</mark>c"
	check(t, msg, expt)
}

func TestNoHighlight1(t *testing.T) {
	msg := "{=a=}"
	expt := "{=a=}"
	check(t, msg, expt)
}

func TestNoHighlight2(t *testing.T) {
	msg := "a==}b"
	expt := "a==}b"
	check(t, msg, expt)
}

//
// {>>  <<}
//

func TestComment(t *testing.T) {
	msg := "{>>a<<}"
	expt := "<span class=\"critic comment\">a</span>"
	check(t, msg, expt)
}

func TestCommentBegin(t *testing.T) {
	msg := "{>>a<<}b"
	expt := "<span class=\"critic comment\">a</span>b"
	check(t, msg, expt)
}

func TestCommentEnd(t *testing.T) {
	msg := "a{>>b<<}"
	expt := "a<span class=\"critic comment\">b</span>"
	check(t, msg, expt)
}

func TestCommentMiddle(t *testing.T) {
	msg := "a{>>b<<}c"
	expt := "a<span class=\"critic comment\">b</span>c"
	check(t, msg, expt)
}

func TestNoComment1(t *testing.T) {
	msg := "{>a<}"
	expt := "{>a<}"
	check(t, msg, expt)
}

func TestNoComment2(t *testing.T) {
	msg := "{<<a>>}"
	expt := "{<<a>>}"
	check(t, msg, expt)
}

func TestNoComment3(t *testing.T) {
	msg := "a<<}b"
	expt := "a<<}b"
	check(t, msg, expt)
}

//
// Helpers
//

type res struct {
	msg  string
	expt string
	real string
	err  error
}

func (r *res) isOk() bool {
	return r.err == nil && r.expt == r.real
}
func (r *res) isErr() bool {
	return r.err != nil
}
func (r *res) check(t *testing.T) {
	if r.isErr() {
		t.Errorf("Error: %v", r.err)
	} else if !r.isOk() {
		t.Errorf(`Failed: expected "%v" but had "%v"`, r.expt, r.real)
	}
}

func run(msg string, fopts ...func(*Options)) (string, error) {
	var res string
	in := bytes.NewBufferString(msg)
	out := bytes.NewBuffer([]byte{})
	if _, err := Critic(out, in, fopts...); err != nil {
		return res, err
	}
	return out.String(), nil
}

func test(msg, expt string, fopts ...func(*Options)) (r res) {
	r.msg = msg
	r.expt = expt
	real, err := run(msg, fopts...)
	if err != nil {
		r.err = err
		return
	}
	r.real = real
	return
}

func check(t *testing.T, msg, expt string, fopts ...func(*Options)) {
	r := test(msg, expt, fopts...)
	r.check(t)
}

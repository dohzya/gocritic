package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/russross/blackfriday"
)

const (
	critic_ins = iota
	critic_del
	critic_sub
	critic_comment
	critic_highligh
)

/*--------------------*/
/*--------------------*/
/*--------------------*/
func Critic(out *bytes.Buffer, data []byte) {
	offset := 0
	for offset < len(data) {

		// copy non-special chars
		for offset < len(data) && data[offset] != '{' {
			out.WriteByte(data[offset])
			offset++
		}

		// prevents index errors
		if offset > len(data)-6 {
			continue
		}

		sopen := data[offset : offset+3] // all delim are 3 chars length
		var oend int                     // offset of the actual content end
		var sclose string
		var tbegin, tclose string
		var kind int
		switch string(sopen) {
		case "{++":
			kind = critic_ins
			sclose = "++}"
			if data[offset+3] == '\n' {
				tbegin = "\n<ins class=\"break\">"
				tclose = "</ins>\n"
			} else {
				tbegin = "<ins>"
				tclose = "</ins>"
			}
		case "{--":
			kind = critic_del
			sclose = "--}"
			tbegin = "<del>"
			tclose = "</del>"
		case "{~~":
			kind = critic_sub
			sclose = "~~}"
		case "{==":
			kind = critic_comment
			sclose = "==}"
			tbegin = "<mark>"
			tclose = "</mark>"
		case "{>>":
			kind = critic_highligh
			sclose = "<<}"
			tbegin = `<span class="critic comment">`
			tclose = "</span>"
		default:
			out.WriteByte(data[offset])
			offset++
			continue
		}
		ostart := offset + 3 // offset of the actual content start
		oend = strings.Index(string(data[ostart:]), sclose)
		if oend < 1 {
			out.WriteByte(data[offset])
			offset++
			continue
		}
		oend += ostart

		if kind == critic_sub {
			odelstart := ostart
			oinsend := oend
			odelend := strings.Index(string(data[odelstart:]), "~>")
			if odelend < 1 {
				out.WriteByte(data[offset])
				offset++
				continue
			}
			odelend += odelstart
			oinsstart := odelend + 2 // len("~>")
			out.WriteString("<del>")
			out.Write(data[odelstart:odelend])
			out.WriteString("</del><ins>")
			out.Write(data[oinsstart:oinsend])
			out.WriteString("</ins>")
		} else {
			out.WriteString(tbegin)
			if kind == critic_ins {
				if data[ostart] == '\n' {
					if (oend - ostart) == 1 { // {++\n++}
						out.WriteString("&nbsp;")
					} else {
						out.Write(data[ostart+1 : oend])
					}
					out.WriteString(tclose)
				} else {
					out.Write(data[ostart:oend])
					out.WriteString(tclose)
				}
			} else {
				out.Write(data[ostart:oend])
				out.WriteString(tclose)
			}
		}

		offset = oend + len(sclose)
	}
}

func ex1(ext int) {
	critic := `lacus{++ est++} Pra{e}sent.`
	exp := `<p>lacus<ins> est</ins> Pra{e}sent.</p>
`
	md := bytes.NewBuffer(make([]byte, 0))
	Critic(md, []byte(critic))
	readb := blackfriday.MarkdownHtml(md.Bytes(), ext)
	real := string(readb)
	// fmt.Printf("critic  : ---%s---\n", critic)
	// fmt.Printf("md      : ---%s---\n", md)
	// fmt.Printf("real    : ---%v---\n", real[:len(real)-1])
	// fmt.Printf("expected: ---%v---\n", exp[:len(exp)-1])
	fmt.Printf("\n%v\n", real == exp)
}

func ex2(ext int) {
	critic := `Don't go around saying{-- to people that--} the world owes you
a living. The world owes you nothing. It was here first. {~~One~>Only one~~}
thing is impossible for God: To find {++any++} sense in any copyright law
on the planet. {==Truth is stranger than fiction==}{>>strange but true<<},
but it is because Fiction is obliged to stick to possibilities; Truth isn't.`
	exp := `<p>Don't go around saying<del> to people that</del> the world owes you
a living. The world owes you nothing. It was here first. <del>One</del><ins>Only one</ins>
thing is impossible for God: To find <ins>any</ins> sense in any copyright law
on the planet. <mark>Truth is stranger than fiction</mark><span class="critic comment">strange but true</span>,
but it is because Fiction is obliged to stick to possibilities; Truth isn't.</p>
`
	md := bytes.NewBuffer(make([]byte, 0))
	Critic(md, []byte(critic))
	readb := blackfriday.MarkdownHtml(md.Bytes(), ext)
	real := string(readb)
	// fmt.Printf("critic  : ---%s---\n", critic)
	// fmt.Printf("md      : ---%s---\n", md)
	// fmt.Printf("real    : ---%v---\n", real[:len(real)-1])
	// fmt.Printf("expected: ---%v---\n", exp[:len(exp)-1])
	fmt.Printf("\n%v\n", real == exp)
}

func ex3(ext int) {
	critic := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum at orci magna. Phasellus augue justo, sodales eu pulvinar ac, vulputate eget nulla. Mauris massa sem, tempor sed cursus et, semper tincidunt lacus.{++
++}Praesent sagittis, quam id egestas consequat, nisl orci vehicula libero, quis ultricies nulla magna interdum sem. Maecenas eget orci vitae eros accumsan mollis. Cras mi mi, rutrum id aliquam in, aliquet vitae tellus. Sed neque justo, cursus in commodo eget, facilisis eget nunc. Cras tincidunt auctor varius.`
	exp := `<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum at orci magna. Phasellus augue justo, sodales eu pulvinar ac, vulputate eget nulla. Mauris massa sem, tempor sed cursus et, semper tincidunt lacus.
<ins class="break">&nbsp;</ins>
Praesent sagittis, quam id egestas consequat, nisl orci vehicula libero, quis ultricies nulla magna interdum sem. Maecenas eget orci vitae eros accumsan mollis. Cras mi mi, rutrum id aliquam in, aliquet vitae tellus. Sed neque justo, cursus in commodo eget, facilisis eget nunc. Cras tincidunt auctor varius.</p>
`
	md := bytes.NewBuffer(make([]byte, 0))
	Critic(md, []byte(critic))
	readb := blackfriday.MarkdownHtml(md.Bytes(), ext)
	real := string(readb)
	// fmt.Printf("critic  : ---%s---\n", critic)
	// fmt.Printf("md      : ---%s---\n", md)
	// fmt.Printf("real    : ---%v---\n", real[:len(real)-1])
	// fmt.Printf("expected: ---%v---\n", exp[:len(exp)-1])
	fmt.Printf("\n%v\n", real == exp)
}

func main() {
	ext := blackfriday.CommonExtensions // | blackfriday.EXTENSION_CRITIC
	ex1(ext)
	ex2(ext)
	ex3(ext)
}

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
func Critic(out *bytes.Buffer, data []byte) int {
	sopen := data[:3] // all delim are 3 chars length
	ostart := 3       // offset of the actual content start
	var oend int      // offset of the actual content end
	var sclose string
	var tbegin, tclose string
	var kind int
	switch string(sopen) {
	case "{++":
		kind = critic_ins
		sclose = "++}"
		fmt.Printf("%s\n", data)
		if data[3] == '\n' {
			tbegin = `<ins class="break">`
		} else {
			tbegin = "<ins>"

		}
		tclose = "</ins>"
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
		return 0
	}
	oend = strings.Index(string(data[ostart:]), sclose)
	if oend < 1 {
		return 0
	}
	oend += ostart

	if kind == critic_sub {
		odelstart := ostart
		oinsend := oend
		odelend := strings.Index(string(data[odelstart:]), "~>")
		if odelend < 1 {
			return 0
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
		out.Write(data[ostart:oend])
		out.WriteString(tclose)
	}

	return oend + len(sclose)
}

func main() {
	ext := blackfriday.CommonExtensions // | blackfriday.EXTENSION_CRITIC

	critic := `lacus{++ est++} Praesent.`
	exp := `<p>lacus<ins> est</ins> Praesent</p>
`
	md := bytes.NewBuffer(make([]byte, len(critic)))
	Critic(md, []byte(critic))
	readb := blackfriday.MarkdownHtml(md.Bytes(), ext)
	real := string(readb)
	fmt.Printf("critic  : ---%s---\n", critic)
	fmt.Printf("md      : ---%s---\n", md)
	fmt.Printf("real    : ---%s---\n", real[:len(real)-1])
	fmt.Printf("expected: ---%s---\n", exp[:len(exp)-1])
	fmt.Printf("\n%v\n", real == exp)

	// 	input := `Don't go around saying{-- to people that--} the world owes you
	//  a living. The world owes you nothing. It was here first. {~~One~>Only one~~}
	//  thing is impossible for God: To find {++any++} sense in any copyright law
	//  on the planet. {==Truth is stranger than fiction==}{>>strange but true<<},
	//  but it is because Fiction is obliged to stick to possibilities; Truth isn't.`
	// 	exp := `<p>Don't go around saying<del> to people that</del> the world owes you
	//  a living. The world owes you nothing. It was here first. <del>One</del><ins>Only one</ins>
	//  thing is impossible for God: To find <ins>any</ins> sense in any copyright law
	//  on the planet. <mark>Truth is stranger than fiction</mark><span class="critic comment">strange but true</span>,
	//  but it is because Fiction is obliged to stick to possibilities; Truth isn't.</p>
	// `
	// 	readb := blackfriday.MarkdownHtml([]byte(input), ext)
	// 	real := string(readb)
	// 	fmt.Printf("raw:     ---%s---\n\n", input)
	// 	fmt.Printf("real:     ---%s---\n\n", real)
	// 	fmt.Printf("expected: ---%s---\n\n", exp)
	// 	fmt.Printf("%v\n\n", real == exp)

	// 	critic2 := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum at orci magna. Phasellus augue justo, sodales eu pulvinar ac, vulputate eget nulla. Mauris massa sem, tempor sed cursus et, semper tincidunt lacus.{++

	// ++}Praesent sagittis, quam id egestas consequat, nisl orci vehicula libero, quis ultricies nulla magna interdum sem. Maecenas eget orci vitae eros accumsan mollis. Cras mi mi, rutrum id aliquam in, aliquet vitae tellus. Sed neque justo, cursus in commodo eget, facilisis eget nunc. Cras tincidunt auctor varius.`
	// 	exp2 := `<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum at orci magna. Phasellus augue justo, sodales eu pulvinar ac, vulputate eget nulla. Mauris massa sem, tempor sed cursus et, semper tincidunt lacus.

	// <ins class=”break”>&nbsp;</ins>

	// Praesent sagittis, quam id egestas consequat, nisl orci vehicula libero, quis ultricies nulla magna interdum sem. Maecenas eget orci vitae eros accumsan mollis. Cras mi mi, rutrum id aliquam in, aliquet vitae tellus. Sed neque justo, cursus in commodo eget, facilisis eget nunc. Cras tincidunt auctor varius.</p>`
	// 	md2 := bytes.NewBuffer(make([]byte, len(critic2)))
	// 	Critic(md2, []byte(critic2))
	// 	readb2 := blackfriday.MarkdownHtml(md2.Bytes(), ext)
	// 	real2 := string(readb2)
	// 	fmt.Printf("critic2:     ---%s---\n\n", critic2)
	// 	fmt.Printf("md2:     ---%s---\n\n", md2)
	// 	fmt.Printf("real2:     ---%s---\n\n", real2)
	// 	fmt.Printf("expected2: ---%s---\n\n", exp2)
	// 	fmt.Printf("%v\n\n", real2 == exp2)

}

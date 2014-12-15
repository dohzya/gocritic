package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/russross/blackfriday"
)

// --------------------
// --------------------
// --------------------

type mode int

const (
	mNormal mode = iota
	mIns
	mDel
	mSub
	mComment
	mHighligh
)

var ops = map[mode]byte{
	mNormal:   '{',
	mIns:      '+',
	mDel:      '-',
	mSub:      '~',
	mComment:  '<',
	mHighligh: '=',
}

type context struct {
	mode         mode
	insTagged    bool
	insMultiline bool
}

func isOp(ctx context, c byte) bool {
	if c == ops[ctx.mode] {
		return true
	}
	if ctx.mode == mIns && !ctx.insTagged {
		return c != '\n' && c != '\r'
	}
	return false
}

// Critic converts critic markup into HTML
func Critic(w io.Writer, r io.Reader) (int, error) {
	rbuf := make([]byte, 32) // actual buffer
	buf := rbuf[2:]          // buf used for reading
	read := 0                // total bytes read
	bi := 2                  // index of 1st byte of rbuf which is a data
	ctx := context{mNormal, false, false}

	// bi allows to keep some bytes from an iteration to an other
main: // main iteration (1 loop = 1 read)
	for {
		ri, errr := r.Read(buf)
		read += ri
		if ri == 0 && errr != nil {
			if bi < 2 {
				// there are some bytes saved from the last iteration
				if _, err := w.Write(rbuf[bi:2]); err != nil {
					return read, err
				}
			}
			if errr != io.EOF {
				return read, errr
			}
			return read, nil
		}
		data := rbuf[bi : 2+ri]
		offset := 0

	sub: // iteration on the data read
		for offset < len(data) {
			i := offset
			// copy non-special chars
			for offset < len(data) && !isOp(ctx, data[offset]) {
				offset++
			}
			if _, err := w.Write(data[i:offset]); err != nil {
				return read, err
			}
			if ctx.mode == mIns && offset > i {
				ctx.insMultiline = true
			}
			if offset >= len(data) {
				bi = 2
				continue main
			}
			// if there are not enough chars to make op, save them for later
			// (actually there is an op of 2 chars only (`~>`) but it can't
			// be used at the EOF because it needs to be followed by `~~}`,
			// so we can store it for the next iteration and risk to not
			// handle it as an op if reaching EOF on the next read)
			if offset > len(data)-2 {
				rbuf[1] = data[offset]
				bi = 1
				continue main
			}
			if offset > len(data)-3 {
				rbuf[0] = data[offset]
				rbuf[1] = data[offset+1]
				bi = 0
				continue main
			}
			// there are more than 3 chars and it could be an op
			switch string(data[offset : offset+3]) {
			case "{++":
				ctx.mode = mIns
				ctx.insTagged = false
				ctx.insMultiline = false
				// the <ins> tag will be writen after having read all
				// `\n` following the `{++` tag.
				offset += 3
				bi = 2
			case "++}":
				var s string
				if !ctx.insTagged {
					if ctx.insMultiline {
						s = "<ins class=\"break\">&nbsp;</ins>\n"
					} else {
						s = "<ins>&nbsp;</ins>"
					}
				} else {
					s = "</ins>"
				}
				ctx.mode = mNormal
				ctx.insTagged = false
				ctx.insMultiline = false
				if _, err := w.Write([]byte(s)); err != nil {
					return read, err
				}
				offset += 3
				bi = 2
			case "{--":
				ctx.mode = mDel
				if _, err := w.Write([]byte("<del>")); err != nil {
					return read, err
				}
				offset += 3
				bi = 2
			case "--}":
				ctx.mode = mNormal
				if _, err := w.Write([]byte("</del>")); err != nil {
					return read, err
				}
				offset += 3
				bi = 2
			case "{~~":
				ctx.mode = mSub
				if _, err := w.Write([]byte("<del>")); err != nil {
					return read, err
				}
				offset += 3
				bi = 2
			case "~~}":
				ctx.mode = mNormal
				if _, err := w.Write([]byte("</ins>")); err != nil {
					return read, err
				}
				offset += 3
				bi = 2
			case "{==":
				ctx.mode = mHighligh
				if _, err := w.Write([]byte("<mark>")); err != nil {
					return read, err
				}
				offset += 3
				bi = 2
			case "==}":
				ctx.mode = mNormal
				if _, err := w.Write([]byte("</mark>")); err != nil {
					return read, err
				}
				offset += 3
				bi = 2
			case "{>>":
				ctx.mode = mComment
				if _, err := w.Write([]byte(`<span class="critic comment">`)); err != nil {
					return read, err
				}
				offset += 3
				bi = 2
			case "<<}":
				ctx.mode = mNormal
				if _, err := w.Write([]byte(`</span>`)); err != nil {
					return read, err
				}
				offset += 3
				bi = 2
			default:
				if ctx.mode == mIns && !ctx.insTagged {
					if _, err := w.Write([]byte("<ins>")); err != nil {
						return read, err
					}
					ctx.insTagged = true
				}
				if ctx.mode == mSub && string(data[offset:offset+2]) == "~>" {
					if _, err := w.Write([]byte(`</del><ins>`)); err != nil {
						return read, err
					}
					offset += 2
					bi = 2
					continue sub
				}
				if _, err := w.Write(data[offset : offset+1]); err != nil {
					return read, err
				}
				offset++
				bi = 2
				continue sub
			}
		}
	}
}

func ex(num int, critic string, exp string) {
	md := bytes.NewBuffer(make([]byte, 0))
	_, err := Critic(md, bytes.NewBufferString(critic))
	if err != nil {
		fmt.Printf("failed: %s\n", err.Error())
		return
	}
	readb := blackfriday.MarkdownHtml(md.Bytes(), blackfriday.CommonExtensions)
	real := string(readb)
	ok := real == exp
	if !ok {
		fmt.Printf("[%d] critic  : ---%s---\n", num, critic)
		fmt.Printf("[%d] md      : ---%s---\n", num, md)
		fmt.Printf("[%d] real    : ---%v---\n", num, real[:len(real)-1])
		fmt.Printf("[%d] expected: ---%v---\n", num, exp[:len(exp)-1])
	}
	fmt.Printf("[%d] %v\n", num, ok)
}

func main() {
	ex(
		1,
		`lacus{++ est++} Pra{e}sent.`,
		`<p>lacus<ins> est</ins> Pra{e}sent.</p>
`,
	)

	ex(
		2,
		`Don't go around saying{-- to people that--} the world owes you
a living. The world owes you nothing. It was here first. {~~One~>Only one~~}
thing is impossible for God: To find {++any++} sense in any copyright law
on the planet. {==Truth is stranger than fiction==}{>>strange but true<<},
but it is because Fiction is obliged to stick to possibilities; Truth isn't.`,
		`<p>Don't go around saying<del> to people that</del> the world owes you
a living. The world owes you nothing. It was here first. <del>One</del><ins>Only one</ins>
thing is impossible for God: To find <ins>any</ins> sense in any copyright law
on the planet. <mark>Truth is stranger than fiction</mark><span class="critic comment">strange but true</span>,
but it is because Fiction is obliged to stick to possibilities; Truth isn't.</p>
`,
	)

	ex(
		3,
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit.{++
++}Vestibulum at orci magna. Phasellus augue justo, sodales eu pulvinar ac,
vulputate eget nulla. Mauris massa sem, tempor sed cursus et, semper tincidunt
lacus. Praesent sagittis, quam id egestas consequat, nisl orci vehicula
libero, quis ultricies nulla magna interdum sem. Maecenas eget orci vitae
eros accumsan mollis. Cras mi mi, rutrum id aliquam in, {~~aliquet vitae~>
~~}tellus. Sed neque justo, cursus in commodo eget, facilisis eget nunc.
Cras tincidunt auctor varius.`,
		`<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
<ins class="break">&nbsp;</ins>
Vestibulum at orci magna. Phasellus augue justo, sodales eu pulvinar ac,
vulputate eget nulla. Mauris massa sem, tempor sed cursus et, semper tincidunt
lacus. Praesent sagittis, quam id egestas consequat, nisl orci vehicula
libero, quis ultricies nulla magna interdum sem. Maecenas eget orci vitae
eros accumsan mollis. Cras mi mi, rutrum id aliquam in, <del>aliquet vitae<del><ins class="break">&nbsp;</ins>
tellus. Sed neque justo, cursus in commodo eget, facilisis eget nunc.
Cras tincidunt auctor varius.</p>
`,
	)
}

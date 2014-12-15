package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

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
	mSub1
	mSub2
	mComment
	mHighligh
)

var ops = map[mode]byte{
	mNormal:   '{',
	mIns:      '+',
	mDel:      '-',
	mSub1:     '~',
	mSub2:     '~',
	mComment:  '<',
	mHighligh: '=',
}

type context struct {
	mode      mode
	insTagged bool
	multiline bool
}

func isOp(ctx context, c byte) bool {
	if c == ops[ctx.mode] {
		return true
	}
	if (ctx.mode == mIns || ctx.mode == mSub2) && !ctx.insTagged {
		return c != '\n' && c != '\r'
	}
	return false
}

// Critic converts critic markup into HTML
func Critic(w io.Writer, r io.Reader) (int, error) {
	rbuf := make([]byte, 3) // actual buffer
	buf := rbuf[2:]         // buf used for reading
	read := 0               // total bytes read
	bi := 2                 // index of 1st byte of rbuf which is a data
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
			if (ctx.mode == mIns || ctx.mode == mSub2) && offset > i {
				ctx.multiline = true
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
				ctx.multiline = false
				// the <ins> tag will be writen after having read all
				// `\n` following the `{++` tag.
				offset += 3
				bi = 2
			case "++}":
				var s string
				if !ctx.insTagged {
					if ctx.multiline {
						s = "<ins class=\"break\">&nbsp;</ins>\n"
					} else {
						s = "<ins>&nbsp;</ins>"
					}
				} else {
					s = "</ins>"
				}
				ctx.mode = mNormal
				ctx.insTagged = false
				ctx.multiline = false
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
				ctx.mode = mSub1
				if _, err := w.Write([]byte("<del>")); err != nil {
					return read, err
				}
				offset += 3
				bi = 2
			case "~~}":
				var s string
				if !ctx.insTagged {
					if ctx.multiline {
						s = "<ins class=\"break\">&nbsp;</ins>\n"
					} else {
						s = "<ins>&nbsp;</ins>"
					}
				} else {
					s = "</ins>"
				}
				ctx.mode = mNormal
				if _, err := w.Write([]byte(s)); err != nil {
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
				if (ctx.mode == mIns || ctx.mode == mSub2) && !ctx.insTagged {
					if _, err := w.Write([]byte("<ins>")); err != nil {
						return read, err
					}
					ctx.insTagged = true
				}
				if ctx.mode == mSub1 && string(data[offset:offset+2]) == "~>" {
					if _, err := w.Write([]byte(`</del>`)); err != nil {
						return read, err
					}
					ctx.mode = mSub2
					ctx.insTagged = false
					ctx.multiline = false
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

func main() {
	md := flag.Bool("md", false, "Use markdown parser")
	flag.Parse()

	var input io.Reader = os.Stdin
	var output io.Writer = os.Stdout

	if *md {
		bMd := bytes.NewBuffer(make([]byte, 0))
		if _, err := Critic(bMd, input); err != nil {
			fmt.Fprintf(os.Stderr, "[gocritic] Error during critic parsing: %s\n", err.Error())
			return
		}
		bHTML := blackfriday.MarkdownHtml(bMd.Bytes(), blackfriday.CommonExtensions)
		if _, err := output.Write(bHTML); err != nil {
			fmt.Fprintf(os.Stderr, "[gocritic] Error while writing result: %s\n", err.Error())
			return
		}
	} else {
		if _, err := Critic(output, input); err != nil {
			fmt.Fprintf(os.Stderr, "[gocritic] Error during critic parsing: %s\n", err.Error())
			return
		}
	}
}

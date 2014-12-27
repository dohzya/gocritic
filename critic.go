package gocritic

import "io"

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

type kind int

const (
	kdText kind = iota
	kdTag
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

func write(w io.Writer, b []byte, k kind, m mode, opts *Options) (int, error) {
	if opts.display(k, m) {
		return w.Write(b)
	}
	return 0, nil
}

func writeTxt(w io.Writer, b []byte, c context, o *Options) (int, error) {
	return write(w, b, kdText, c.mode, o)
}

func writeTg(w io.Writer, b []byte, c context, o *Options) (int, error) {
	return write(w, b, kdTag, c.mode, o)
}

// Critic converts critic markup into HTML
func Critic(w io.Writer, r io.Reader, fopts ...func(*Options)) (int, error) {
	opts := createOptions(fopts...)
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
				if _, err := writeTxt(w, rbuf[bi:2], ctx, opts); err != nil {
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
			if _, err := writeTxt(w, data[i:offset], ctx, opts); err != nil {
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
				ctx.mode = mNormal
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
				if _, err := writeTg(w, []byte(s), ctx, opts); err != nil {
					return read, err
				}
				ctx.insTagged = false
				ctx.multiline = false
				offset += 3
				bi = 2
			case "{--":
				ctx.mode = mDel
				if _, err := writeTg(w, []byte("<del>"), ctx, opts); err != nil {
					return read, err
				}
				offset += 3
				bi = 2
			case "--}":
				if _, err := writeTg(w, []byte("</del>"), ctx, opts); err != nil {
					return read, err
				}
				ctx.mode = mNormal
				offset += 3
				bi = 2
			case "{~~":
				ctx.mode = mSub1
				if _, err := writeTg(w, []byte("<del>"), ctx, opts); err != nil {
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
				if _, err := writeTg(w, []byte(s), ctx, opts); err != nil {
					return read, err
				}
				ctx.mode = mNormal
				offset += 3
				bi = 2
			case "{==":
				ctx.mode = mHighligh
				if _, err := writeTg(w, []byte("<mark>"), ctx, opts); err != nil {
					return read, err
				}
				offset += 3
				bi = 2
			case "==}":
				if _, err := writeTg(w, []byte("</mark>"), ctx, opts); err != nil {
					return read, err
				}
				ctx.mode = mNormal
				offset += 3
				bi = 2
			case "{>>":
				ctx.mode = mComment
				if _, err := writeTg(w, []byte(`<span class="critic comment">`), ctx, opts); err != nil {
					return read, err
				}
				offset += 3
				bi = 2
			case "<<}":
				if _, err := writeTg(w, []byte("</span>"), ctx, opts); err != nil {
					return read, err
				}
				ctx.mode = mNormal
				offset += 3
				bi = 2
			default:
				if (ctx.mode == mIns || ctx.mode == mSub2) && !ctx.insTagged {
					var s string
					if ctx.multiline {
						s = "<ins class=\"break\">"
					} else {
						s = "<ins>"
					}
					if _, err := writeTg(w, []byte(s), ctx, opts); err != nil {
						return read, err
					}
					ctx.insTagged = true
				}
				if ctx.mode == mSub1 && string(data[offset:offset+2]) == "~>" {
					if _, err := writeTg(w, []byte("</del>"), ctx, opts); err != nil {
						return read, err
					}
					ctx.mode = mSub2
					ctx.insTagged = false
					ctx.multiline = false
					offset += 2
					bi = 2
					continue sub
				}
				if _, err := writeTxt(w, data[offset:offset+1], ctx, opts); err != nil {
					return read, err
				}
				offset++
				bi = 2
				continue sub
			}
		}
	}
}

package gocritic

// Options allows client code to change the Critic's behavior.
//
// For now it allows to filter output.
type Options struct {
	filter filter
}

func (o *Options) hasFilter() bool {
	return o.filter != fNone
}
func (o *Options) display(m mode) bool {
	return o.filter == fNone || m == mNormal || m == mHighligh ||
		o.filter == fBefore && (m == mDel || m == mSub1) ||
		o.filter == fAfter && (m == mIns || m == mSub2)
}

func createOptions(fopts ...func(*Options)) *Options {
	opts := &Options{fNone}
	for _, f := range fopts {
		f(opts)
	}
	return opts
}

// FilterNone specifies that no filter is active
func FilterNone(o *Options) {
	o.filter = fNone
}

// FilterBefore specifies that only original source is returned
//
// Removes contents of {++++} and keeps content of {----}
func FilterBefore(o *Options) {
	o.filter = fBefore
}

// FilterAfter specifies that only reviewed source is returned
//
// Removes contents of {----} and keeps content of {++++}
func FilterAfter(o *Options) {
	o.filter = fAfter
}

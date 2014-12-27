package gocritic

// Options allows client code to change the Critic's behavior.
//
// For now it allows to filter output.
type Options struct {
	kinds map[kind]bool
	modes map[mode]bool
}

func (o *Options) display(k kind, m mode) bool {
	return o.kinds[k] && o.modes[m]
}

func createOptions(fopts ...func(*Options)) *Options {
	opts := &Options{make(map[kind]bool), make(map[mode]bool)}
	FilterShowAll(opts)
	for _, f := range fopts {
		f(opts)
	}
	return opts
}

// FilterShowAll specifies that no filter is active
func FilterShowAll(o *Options) {
	o.kinds[kdText] = true
	o.kinds[kdTag] = true
	o.modes[mNormal] = true
	o.modes[mIns] = true
	o.modes[mDel] = true
	o.modes[mSub1] = true
	o.modes[mSub2] = true
	o.modes[mComment] = true
	o.modes[mHighligh] = true
}

// FilterHideBefore specifies that original sources are not rendered
func FilterHideBefore(o *Options) {
	o.modes[mDel] = false
	o.modes[mSub1] = false
}

// FilterShowBefore specifies that original sources are rendered
func FilterShowBefore(o *Options) {
	o.modes[mDel] = true
	o.modes[mSub1] = true
}

// FilterHideAfter specifies that reviewed sources are not rendered
func FilterHideAfter(o *Options) {
	o.modes[mIns] = false
	o.modes[mSub2] = false
}

// FilterShowAfter specifies that reviewed sources not rendered
func FilterShowAfter(o *Options) {
	o.modes[mIns] = true
	o.modes[mSub2] = true
}

// FilterHideComments specifies that comments are not rendered
//
// This filter also deactivate tags rendering
func FilterHideComments(o *Options) {
	o.modes[mComment] = false
	o.kinds[kdTag] = false
}

// FilterShowComments specifies that comments not rendered
//
// This filter also activate tags rendering
func FilterShowComments(o *Options) {
	o.modes[mComment] = true
	o.kinds[kdTag] = true
}

// FilterHideTags specifies that tags are not rendered
func FilterHideTags(o *Options) {
	o.kinds[kdTag] = false
}

// FilterShowTags specifies that tags are rendered
func FilterShowTags(o *Options) {
	o.kinds[kdTag] = true
}

//
// High level
//

// FilterOnlyBefore specifies that only original sources is rendered
func FilterOnlyBefore(o *Options) {
	FilterHideAfter(o)
}

// FilterOnlyRawBefore specifies that only original sources is rendered
func FilterOnlyRawBefore(o *Options) {
	FilterHideAfter(o)
	FilterHideComments(o)
	FilterHideTags(o)
}

// FilterOnlyAfter specifies that only reviewed sources is rendered
func FilterOnlyAfter(o *Options) {
	FilterHideBefore(o)
}

// FilterOnlyRawAfter specifies that only reviewed sources is rendered
func FilterOnlyRawAfter(o *Options) {
	FilterHideBefore(o)
	FilterHideComments(o)
	FilterHideTags(o)
}

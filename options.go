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

// FilterHideOriginal specifies that original sources are not rendered
func FilterHideOriginal(o *Options) {
	o.modes[mDel] = false
	o.modes[mSub1] = false
}

// FilterShowOriginal specifies that original sources are rendered
func FilterShowOriginal(o *Options) {
	o.modes[mDel] = true
	o.modes[mSub1] = true
}

// FilterHideEdited specifies that edited sources are not rendered
func FilterHideEdited(o *Options) {
	o.modes[mIns] = false
	o.modes[mSub2] = false
}

// FilterShowEdited specifies that edited sources not rendered
func FilterShowEdited(o *Options) {
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

// FilterOnlyOriginal specifies that only original sources are rendered
func FilterOnlyOriginal(o *Options) {
	FilterHideEdited(o)
}

// FilterOnlyRawOriginal specifies that only original sources are rendered
func FilterOnlyRawOriginal(o *Options) {
	FilterHideEdited(o)
	FilterHideComments(o)
	FilterHideTags(o)
}

// FilterOnlyEdited specifies that only edited sources are rendered
func FilterOnlyEdited(o *Options) {
	FilterHideOriginal(o)
}

// FilterOnlyRawEdited specifies that only edited sources are rendered
func FilterOnlyRawEdited(o *Options) {
	FilterHideOriginal(o)
	FilterHideComments(o)
	FilterHideTags(o)
}

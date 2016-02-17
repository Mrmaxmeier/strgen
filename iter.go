package strgen

type Iterator interface {
	push(item)
	configure() error
	cycle()
	get() string
	length() int
	finished() bool
	setCyclePos(int)
	cleanup()
}

type TextIterator struct {
	text string
}

func (i *TextIterator) push(it item)         { i.text = it.val }
func (i *TextIterator) get() string          { return i.text }
func (i *TextIterator) cycle()               {}
func (i *TextIterator) configure() (e error) { return }
func (i *TextIterator) length() int          { return 1 }
func (i *TextIterator) finished() bool       { return true }
func (i *TextIterator) setCyclePos(int)      {}
func (i *TextIterator) cleanup()             {}

type IteratorsByLength []Iterator

func (s IteratorsByLength) Len() int {
	return len(s)
}
func (s IteratorsByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s IteratorsByLength) Less(i, j int) bool {
	return s[i].length() < s[j].length()
}

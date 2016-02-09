package strgen

import (
	"fmt"

	"github.com/bradfitz/slice"
)

type Generator struct {
	source    string
	err       error
	amount    int64
	left      int64
	current   int64
	infinite  bool
	doneCh    chan interface{}
	results   chan string
	iterators []Iterator
}

func (g *Generator) configure() error {
	if g.results == nil {
		g.results = make(chan string)
	}
	if g.doneCh == nil {
		g.doneCh = make(chan interface{})
	}
	_, items := lex(g.source)
	var currIter Iterator
	emit := func() {
		currIter.configure()
		g.iterators = append(g.iterators, currIter)
		currIter = nil
	}
	for item := range items {
		switch item.typ {
		case itemRange:
			currIter = &RangeIterator{}
		case itemChoice:
			currIter = &ChoiceIterator{}
		case itemText:
			if currIter == nil {
				currIter = &TextIterator{text: item.val}
				emit()
			} else {
				currIter.push(item)
			}
		case itemIterEnd:
			emit()
		case itemEOF:
		case itemError:
			g.err = fmt.Errorf(item.val)
			return g.err
		default:
			currIter.push(item)
		}
	}

	sortable := make([]Iterator, len(g.iterators))
	for i := 0; i < len(g.iterators); i++ {
		err := g.iterators[i].configure()
		if err != nil {
			g.err = err
			return err
		}
		sortable[i] = g.iterators[i]
	}
	slice.Sort(sortable[:], func(i, j int) bool {
		return sortable[i].length() < sortable[j].length() || sortable[j].length() == -1
	})
	cyclepos := 1
	for i := 0; i < len(sortable); i++ {
		sortable[i].setCyclePos(cyclepos)
		//fmt.Printf("%+v\n", sortable[i])
		cyclepos *= sortable[i].length()
		if sortable[i].length() == -1 {
			g.infinite = true
		}
	}
	if g.infinite {
		g.amount = -1
	} else {
		g.amount = int64(cyclepos)
	}
	return nil
}

func (g *Generator) generate() {
	g.left = g.amount
	for g.alive() {
		s := ""
		for i := 0; i < len(g.iterators); i++ {
			s += g.iterators[i].get()
			g.iterators[i].cycle()
		}
		select {
		case g.results <- s:
			g.current++
			if !g.infinite {
				g.left--
				if g.iterators[len(g.iterators)-1].finished() {
					g.kill()
					break
				}
			}
		case <-g.doneCh:
			break
		}
	}
	close(g.results)
}

func (g *Generator) alive() bool {
	select {
	case <-g.doneCh:
		return false
	default:
		return true
	}
}

func (g *Generator) kill() {
	if g.alive() {
		close(g.doneCh)
	}
}

func (g *Generator) next() (s string, err error) {
	select {
	case <-g.doneCh:
		err = fmt.Errorf("channel closed")
		return
	case s = <-g.results:
		if !g.alive() {
			err = fmt.Errorf("channel closed")
		}
		return
	}
}

func GenerateStrings(optionGen string) (<-chan string, int64, error) {
	g := Generator{source: optionGen}
	if g.configure() != nil {
		return nil, 0, g.err
	}
	go g.generate()
	return g.results, g.amount, g.err
}

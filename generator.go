package strgen

import (
	"fmt"
	"runtime"
	"sort"
)

type Generator struct {
	Source    string
	Err       error
	Amount    int64
	Left      int64
	Current   int64
	Infinite  bool
	DoneCh    chan interface{}
	Results   chan string
	Iterators []Iterator
}

func (g *Generator) Configure() error {
	if g.Results == nil {
		g.Results = make(chan string)
	}
	if g.DoneCh == nil {
		g.DoneCh = make(chan interface{})
	}
	_, items := lex(g.Source)
	var currIter Iterator
	emit := func() {
		g.Iterators = append(g.Iterators, currIter)
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
		case itemExpr:
			currIter = &ExpressionIterator{code: item.val}
		case itemIterEnd:
			emit()
		case itemEOF:
		case itemError:
			g.Err = fmt.Errorf(item.val)
			return g.Err
		default:
			currIter.push(item)
		}
	}

	sortable := make(IteratorsByLength, len(g.Iterators))
	for i := 0; i < len(g.Iterators); i++ {
		err := g.Iterators[i].configure()
		if err != nil {
			g.Err = err
			return err
		}
		sortable[i] = g.Iterators[i]
	}
	sort.Sort(sortable)
	cyclepos := 1
	for i := 0; i < len(sortable); i++ {
		sortable[i].setCyclePos(cyclepos)
		//fmt.Printf("%+v\n", sortable[i])
		cyclepos *= sortable[i].length()
		if sortable[i].length() == -1 {
			g.Infinite = true
		}
	}
	if g.Infinite {
		g.Amount = -1
	} else {
		g.Amount = int64(cyclepos)
	}
	return nil
}

func (g *Generator) Generate() {
	g.Left = g.Amount
	for g.Alive() {
		s := ""
		for i := 0; i < len(g.Iterators); i++ {
			s += g.Iterators[i].get()
			g.Iterators[i].cycle()
		}
		select {
		case g.Results <- s:
			g.Current++
			if !g.Infinite {
				g.Left--
				if g.Iterators[len(g.Iterators)-1].finished() {
					g.Close()
					break
				}
			}
		case <-g.DoneCh:
			break
		}
	}
	close(g.Results)

	for i := 0; i < len(g.Iterators); i++ {
		g.Iterators[i].cleanup()
	}
}

func (g *Generator) Alive() bool {
	if g.Err != nil {
		return false
	}
	select {
	case <-g.DoneCh:
		return false
	default:
		return true
	}
}

func (g *Generator) Close() {
	if g.Alive() {
		close(g.DoneCh)
		runtime.Gosched()
	}
}

func (g *Generator) Next() (s string, err error) {
	if !g.Alive() {
		err = fmt.Errorf("channel closed")
	}
	select {
	case <-g.DoneCh:
		err = fmt.Errorf("channel closed")
	case s = <-g.Results:
	}
	return
}

func GenerateStrings(optionGen string) (<-chan string, int64, error) {
	g := Generator{Source: optionGen}
	if g.Configure() != nil {
		return nil, 0, g.Err
	}
	go g.Generate()
	return g.Results, g.Amount, g.Err
}

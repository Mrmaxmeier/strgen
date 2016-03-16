package strgen

import (
	"strconv"

	"github.com/yuin/gopher-lua"
)

type ExpressionIterator struct {
	cycleLength  int
	tmpCycle     int
	cyclepos     int
	currentCycle int
	lastComputed *struct {
		cycle int
		value string
	}
	code         string
	state        *lua.LState
	valueChannel chan lua.LValue
	err          error
}

func (i *ExpressionIterator) push(it item) {
	if it.typ == itemText {
		if i.cycleLength == 0 {
			cycleLength, err := strconv.Atoi(it.val)
			if err != nil {
				i.err = err
			}
			i.cycleLength = cycleLength
		} else {
			i.code = "return " + it.val
		}
	}
}

func (i *ExpressionIterator) cycle() {
	i.tmpCycle = (i.tmpCycle + 1) % i.cyclepos
	if i.tmpCycle == 0 {
		i.currentCycle++
		if i.length() > 0 {
			i.currentCycle = i.currentCycle % i.length()
		}
	}
}

func (i *ExpressionIterator) get() string {
	if i.lastComputed != nil && i.currentCycle == i.lastComputed.cycle {
		return i.lastComputed.value
	}
	if i.lastComputed == nil {
		i.lastComputed = &struct {
			cycle int
			value string
		}{}
	}
	i.lastComputed.cycle = i.currentCycle
	i.state.SetGlobal("x", lua.LNumber(i.currentCycle))
	if err := i.state.DoString(i.code); err != nil {
		i.err = err
		i.lastComputed.value = i.err.Error()
		return i.lastComputed.value
	}
	i.lastComputed.value = i.state.Get(-1).String()
	i.state.Pop(1)
	return i.lastComputed.value
}

func (i *ExpressionIterator) configure() error {
	if i.err != nil {
		return i.err
	}

	i.state = lua.NewState(lua.Options{
		CallStackSize: 120,
		RegistrySize:  120 * 20,
		SkipOpenLibs:  true,
	})

	return nil
}
func (i *ExpressionIterator) length() int         { return i.cycleLength }
func (i *ExpressionIterator) finished() bool      { return i.currentCycle == 0 && i.tmpCycle == 0 }
func (i *ExpressionIterator) setCyclePos(pos int) { i.cyclepos = pos }
func (i *ExpressionIterator) cleanup() {
	i.state.Close()
}

package utils

const NoIdCreated = -1

//自动生成自动增加的ID，类型为int
type AutoIncIDGen struct {
	start, step, currentID int
}

type AutoIncIDOption struct {
	f func(*AutoIncIDGen)
}

func newAutoIncIDOption(funcOption func(id *AutoIncIDGen)) *AutoIncIDOption {
	return &AutoIncIDOption{
		f: funcOption,
	}
}

func (opt *AutoIncIDOption) apply(id *AutoIncIDGen) {
	opt.f(id)
}

//创建自增id生成器时的默认设置，起始id为1，步长为1
func defaultOption(id *AutoIncIDGen) {
	id.start = 1
	id.step = 1
	id.currentID = NoIdCreated
}

//创建id生成器时设置起始ID
func WithStartID(start int) *AutoIncIDOption {
	return newAutoIncIDOption(func(id *AutoIncIDGen) {
		id.start = start
	})
}

//创建id生成器时设置id自增步长
func WithStep(step int) *AutoIncIDOption {
	return newAutoIncIDOption(func(id *AutoIncIDGen) {
		id.step = step
	})
}

//创建自增id生成器
func NewAutoIncID(opts ...*AutoIncIDOption) *AutoIncIDGen {
	aID := AutoIncIDGen{}
	defaultOption(&aID)

	for _, opt := range opts {
		opt.apply(&aID)
	}

	return &aID
}

func (gen *AutoIncIDGen) GenID() int {
	if gen.currentID == NoIdCreated {
		gen.currentID = gen.start
	} else {
		gen.currentID += gen.step
	}

	return gen.currentID
}

func (gen *AutoIncIDGen) GetStart() int {
	return gen.start
}

func (gen *AutoIncIDGen) GetStep() int {
	return gen.step
}

func (gen *AutoIncIDGen) GetCurrentID() int {
	return gen.currentID
}

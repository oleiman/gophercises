package transform

// A bit of structure for working with templates in Exercise 17 (image transform service)

type transform struct {
	Mode int
	N    int
}

type TransformTable struct {
	Name       string
	Transforms [][]transform
	Final      bool
}

type TableOpt func(t *TransformTable)

func NewTransformTable(fname string, m, n int, opts ...TableOpt) TransformTable {
	t := TransformTable{
		Name:       fname,
		Transforms: make([][]transform, m),
		Final:      false,
	}
	for i := range t.Transforms {
		t.Transforms[i] = make([]transform, n)
	}
	TableDefaults(&t)
	for _, opt := range opts {
		opt(&t)
	}
	return t
}

func (t *TransformTable) map2D(f func(e *transform, i, j int)) {
	for i := range t.Transforms {
		for j := range t.Transforms[i] {
			f(&t.Transforms[i][j], i, j)
		}
	}
}

func SingleMode(mode int) TableOpt {
	return func(t *TransformTable) {
		t.map2D(func(e *transform, i, j int) {
			e.Mode = mode
		})
	}
}

func IncMode(base, inc int) TableOpt {
	return func(t *TransformTable) {
		t.map2D(func(e *transform, i, j int) {
			e.Mode = base + inc*(i*len(t.Transforms)+j)
		})
	}
}

func SingleN(n int) TableOpt {
	return func(t *TransformTable) {
		t.map2D(func(e *transform, i, j int) {
			e.N = n
		})
	}
}

func IncN(base, inc int) TableOpt {
	return func(t *TransformTable) {
		t.map2D(func(e *transform, i, j int) {
			e.N = base + inc*(len(t.Transforms)*i+j)
		})
	}
}

func Downloadable() TableOpt {
	return func(t *TransformTable) {
		t.Final = true
	}
}

func TableDefaults(t *TransformTable) {
	// IncMode(0, 1)(t)
	SingleMode(0)(t)
	SingleN(1)(t)
}

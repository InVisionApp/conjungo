package merge

import "reflect"

type funcSelector struct {
	typeFuncs   map[reflect.Type]MergeFunc
	kindFuncs   map[reflect.Kind]MergeFunc
	defaultFunc MergeFunc
}

func newFuncSelector() *funcSelector {
	return &funcSelector{
		typeFuncs: map[reflect.Type]MergeFunc{
			reflect.TypeOf(map[string]interface{}{}): mergeMap, //recursion becomes less obvious but allows custom handler
			reflect.TypeOf([]interface{}{}):          mergeSlice,
		},
		kindFuncs:   map[reflect.Kind]MergeFunc{},
		defaultFunc: defaultMergeFunc,
	}
}

func (f *funcSelector) SetTypeMergeFunc(t reflect.Type, mf MergeFunc) {
	f.typeFuncs[t] = mf
}

func (f *funcSelector) SetKindMergeFunc(k reflect.Kind, mf MergeFunc) {
	f.kindFuncs[k] = mf
}

func (f *funcSelector) SetDefaultMergeFunc(mf MergeFunc) {
	f.defaultFunc = mf
}

// Get func must always return a function.
func (f *funcSelector) GetFunc(i interface{}) MergeFunc {
	// prioritize a specific 'type' definition
	ti := reflect.TypeOf(i)
	if fx, ok := f.typeFuncs[ti]; ok {
		return fx
	}

	// then look for a more general 'kind'.
	if fx, ok := f.kindFuncs[ti.Kind()]; ok {
		return fx
	}

	return f.defaultFunc
}

// A function which defines how two items of the same type are merged together.
// Options are also passed in and it is the responsibility of the merge function to handle
// any variations in behavior that should occur. The value returned from the function will be
// written to directly to the target map, as long as there is no error.
type MergeFunc func(interface{}, interface{}, *Options) (interface{}, error)

func defaultMergeFunc(t, s interface{}, o *Options) (interface{}, error) {
	if o.Overwrite {
		return s, nil
	}

	return t, nil
}

func mergeMap(t, s interface{}, o *Options) (interface{}, error) {
	mapT, _ := t.(map[string]interface{})
	mapS, _ := s.(map[string]interface{})

	if err := merge(&mapT, mapS, o); err != nil {
		return nil, err
	}

	return mapT, nil
}

func mergeSlice(t, s interface{}, o *Options) (interface{}, error) {
	sliceT, _ := t.([]interface{})
	sliceS, _ := s.([]interface{})
	return append(sliceT, sliceS...), nil
}

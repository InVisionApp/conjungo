package merge

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"reflect"
)

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
	if nil == f.typeFuncs {
		f.typeFuncs = map[reflect.Type]MergeFunc{}
	}
	f.typeFuncs[t] = mf
}

func (f *funcSelector) SetKindMergeFunc(k reflect.Kind, mf MergeFunc) {
	if nil == f.kindFuncs {
		f.kindFuncs = map[reflect.Kind]MergeFunc{}
	}
	f.kindFuncs[k] = mf
}

func (f *funcSelector) SetDefaultMergeFunc(mf MergeFunc) {
	f.defaultFunc = mf
}

// Get func must always return a function.
// First looks for a merge func defined for its type. Type is the most specific way to categorize something,
// for example, struct type foo of package bar or map[string]string. Next it looks for a merge func defined for its
// kind, for example, struct or map. At this point, if nothing matches, it will fall back to the default merge definition.
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

	if f.defaultFunc != nil {
		return f.defaultFunc
	}

	return defaultMergeFunc
}

// A function which defines how two items of the same type are merged together.
// Options are also passed in and it is the responsibility of the merge function to handle
// any variations in behavior that should occur. The value returned from the function will be
// written to directly to the target map, as long as there is no error.
type MergeFunc func(interface{}, interface{}, *Options) (interface{}, error)

// The most basic merge function to be used as default behavior. In overwrite mode, it returns the source. Otherwise,
// it returns the target.
func defaultMergeFunc(t, s interface{}, o *Options) (interface{}, error) {
	if o.Overwrite {
		return s, nil
	}

	return t, nil
}

func mergeMap(t, s interface{}, o *Options) (interface{}, error) {
	mapT, _ := t.(map[string]interface{})
	mapS, _ := s.(map[string]interface{})

	// if empty, use the source
	if len(mapT) < 1 {
		return mapS, nil
	}

	for k, valS := range mapS {
		logrus.Debugf("MERGE T<>S '%s' :: %v <> %v", k, mapT[k], valS)
		val, err := merge(mapT[k], valS, o)
		if err != nil {
			return nil, fmt.Errorf("key '%s': %v", k, err)
		}
		mapT[k] = val
	}

	return mapT, nil
}

func mergeSlice(t, s interface{}, o *Options) (interface{}, error) {
	sliceT, _ := t.([]interface{})
	sliceS, _ := s.([]interface{})
	return append(sliceT, sliceS...), nil
}

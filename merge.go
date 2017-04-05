package merge

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Sirupsen/logrus"
)

type Options struct {
	Overwrite  bool
	MergeFuncs *funcSelector

	// to be used by merge functions to pass values down into recursive calls freely
	Context context.Context
}

func NewOptions() *Options {
	return &Options{
		Overwrite:  true,
		MergeFuncs: newFuncSelector(),
	}
}

// helper to wrap type assertion
func MergeMapStrIface(target, src map[string]interface{}, opt *Options) (map[string]interface{}, error) {
	val, err := Merge(target, src, opt)
	if err != nil {
		return nil, err
	}

	valMap, ok := val.(map[string]interface{})
		if ok {
		return valMap, nil
	}

	return nil, fmt.Errorf("Merge failed. Expected map[string]interface{} but got %v", reflect.TypeOf(val))
}

// public wrapper
func Merge(target, source interface{}, opt *Options) (interface{}, error) {
	// use defaults if none are provided
	if opt == nil {
		opt = NewOptions()
	}
	logrus.Debugf("OPT: %v", opt)

	return merge(target, source, opt)
}

func merge(target, src interface{}, opt *Options) (interface{}, error) {
	valS := reflect.ValueOf(src)
	valT := reflect.ValueOf(target)

	// if source is nil, skip
	if src == nil ||
		valS.Kind() == reflect.Ptr && valS.IsNil() {
		return target, nil
	}

	// if target is nil write to it
	if target == nil ||
		valT.Kind() == reflect.Ptr && valT.IsNil() {
		return src, nil
	}

	typeS := reflect.TypeOf(src)
	typeT := reflect.TypeOf(target)

	// if types do not match, bail
	if typeT != typeS {
		return nil, fmt.Errorf("Types do not match: %v, %v", typeT, typeS)
	}

	// look for a merge function
	f := opt.MergeFuncs.GetFunc(target)
	val, err := f(target, src, opt)
	if err != nil {
		return nil, err
	}

	return val, nil
}

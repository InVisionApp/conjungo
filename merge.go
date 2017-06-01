package conjungo

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

	merged, err := merge(reflect.ValueOf(target), reflect.ValueOf(source), opt)
	if err != nil {
		return nil, err
	}

	if !merged.IsValid() {
		return nil, nil
	}

	return merged.Interface(), nil
}

func merge(valT, valS reflect.Value, opt *Options) (reflect.Value, error) {
	// if source is nil, skip
	if !valS.IsValid() ||
		valS.Kind() == reflect.Ptr && valS.IsNil() {
		return valT, nil
	}

	// if target is nil write to it
	if !valT.IsValid() ||
		valT.Kind() == reflect.Ptr && valT.IsNil() {
		return valS, nil
	}

	// if types do not match, bail
	if valT.Type() != valS.Type() {
		return reflect.Value{}, fmt.Errorf("Types do not match: %v, %v", valT.Type(), valS.Type())
	}

	// look for a merge function
	f := opt.MergeFuncs.GetFunc(valT)
	val, err := f(valT, valS, opt)
	if err != nil {
		return reflect.Value{}, err
	}

	return val, nil
}

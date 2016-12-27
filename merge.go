package merge

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"reflect"
)

type Options struct {
	Overwrite  bool //TODO: actually use this
	mergeFuncs map[reflect.Type]MergeFunc
}

func NewOptions() *Options {
	return &Options{
		Overwrite: true,
		mergeFuncs: map[reflect.Type]MergeFunc{
			reflect.TypeOf(map[string]interface{}{}): mergeMap, //recursion becomes less obvious but allows custom handler
			reflect.TypeOf([]interface{}{}):          mergeSlice,
		},
	}
}

func (o *Options) SetMergeFunc(t reflect.Type, f MergeFunc) {
	o.mergeFuncs[t] = f
}

func Merge(target, src map[string]interface{}, opt *Options) (map[string]interface{}, error) {
	// use defaults if none are provided
	if opt == nil {
		opt = NewOptions()
	}
	logrus.Debugf("OPT: %v", opt)

	targetCopy := copyMap(target)
	if err := merge(&targetCopy, src, opt); err != nil {
		return nil, err
	}
	return targetCopy, nil
}

func merge(target *map[string]interface{}, src map[string]interface{}, opt *Options) error {
	logrus.Debugf("TARGET: %v", target)
	logrus.Debugf("SOURCE: %v", src)
	for k, v := range src {
		typeS := reflect.TypeOf(v)
		typeT := reflect.TypeOf((*target)[k])
		logrus.Debugf("TYPE T<>S '%s': %v <> %v", k, typeT, typeS)

		// a new value to insert
		if typeT == nil || typeS == nil {
			(*target)[k] = v
			continue
		}
		if typeT != typeS {
			return fmt.Errorf("Types do not match for key '%s': %v, %v", k, typeT, typeS)
		}

		// otherwise look for a merge function
		f, ok := opt.mergeFuncs[typeS]
		if ok { // if a custom merge is defined, use it
			if val, err := f((*target)[k], v, opt); err != nil {
				return err
			} else {
				(*target)[k] = val
			}
			continue
		}

		// otherwise just overwrite or insert new
		(*target)[k] = v

	}
	return nil
}

type MergeFunc func(interface{}, interface{}, *Options) (interface{}, error)

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

func copyMap(m map[string]interface{}) map[string]interface{} {
	newMap := map[string]interface{}{}
	for k, v := range m {
		newMap[k] = v
	}

	return newMap
}

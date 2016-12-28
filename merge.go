package merge

import (
	"fmt"
	"github.com/InVisionApp/go-merge/util"
	"github.com/Sirupsen/logrus"
	"reflect"
)

type Options struct {
	Overwrite  bool
	mergeFuncs map[reflect.Type]MergeFunc
	//TODO: also do one indexed by refelct.Kind to allow broader merge definitions
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

	targetCopy := util.CopyMap(target)
	if err := merge(&targetCopy, src, opt); err != nil {
		return nil, err
	}
	return targetCopy, nil
}

func merge(target *map[string]interface{}, src map[string]interface{}, opt *Options) error {
	for k, v := range src {
		typeS := reflect.TypeOf(v)

		origT, okT := (*target)[k]
		typeT := reflect.TypeOf(origT)

		logrus.Debugf("MERGE T<>S '%s' :: %v <> %v :: %v <> %v", k, origT, v, typeT, typeS)

		// if source is nil, skip
		if v == nil {
			continue
		}

		// insert if it does not exist in target
		if !okT {
			(*target)[k] = v
			continue
		}

		// if types do not match, bail
		if typeT != typeS {
			return fmt.Errorf("Types do not match for key '%s': %v, %v", k, typeT, typeS)
		}

		// otherwise look for a merge function
		f, ok := opt.mergeFuncs[typeS]
		if ok { // if a custom merge is defined, use it (and catch errors)
			if val, err := f((*target)[k], v, opt); err != nil {
				return err
			} else if opt.Overwrite {
				(*target)[k] = val
			}
			continue
		}

		// otherwise just overwrite
		if opt.Overwrite {
			(*target)[k] = v
		}
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

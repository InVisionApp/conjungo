package merge

import (
	"fmt"
	"github.com/InVisionApp/go-merge/util"
	"github.com/Sirupsen/logrus"
	"reflect"
)

type Options struct {
	Overwrite  bool
	MergeFuncs *funcSelector
}

func NewOptions() *Options {
	return &Options{
		Overwrite:  true,
		MergeFuncs: newFuncSelector(),
	}
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
	for k, valS := range src {
		typeS := reflect.TypeOf(valS)

		valT, okT := (*target)[k]
		typeT := reflect.TypeOf(valT)

		logrus.Debugf("MERGE T<>S '%s' :: %v <> %v :: %v <> %v", k, valT, valS, typeT, typeS)

		// if source is nil, skip
		if valS == nil {
			continue
		}

		// insert if it does not exist in target
		if !okT {
			(*target)[k] = valS
			continue
		}

		// if types do not match, bail
		if typeT != typeS {
			return fmt.Errorf("Types do not match for key '%s': %v, %v", k, typeT, typeS)
		}

		// look for a merge function
		f := opt.MergeFuncs.GetFunc(valT)
		val, err := f((*target)[k], valS, opt)
		if err != nil {
			return err
		}

		(*target)[k] = val
	}

	return nil
}

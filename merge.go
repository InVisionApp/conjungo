package merge

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"reflect"
)

type Options struct {
	Overwrite bool //TODO: actually use this
}

func Merge(target, src map[string]interface{}, opt *Options) (map[string]interface{}, error) {
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
		logrus.Debugf("TYPE T<>S: %v <> %v", typeT, typeS)

		// a new value to insert
		if typeT == nil { //TODO: do this better
			(*target)[k] = v
			continue
		}
		if typeS != typeT {
			return fmt.Errorf("Types do not match for key '%s'", k)
		}

		// if the value is a sub map handle it
		if typeS == reflect.TypeOf(map[string]interface{}{}) { // might be a better way to determine this
			mapT, _ := (*target)[k].(map[string]interface{})
			mapS, _ := v.(map[string]interface{})
			if err := merge(&mapT, mapS, opt); err != nil {
				return err
			}
			continue
		}

		// otherwise look for a merge function
		f, ok := mergeFuncs[typeS]
		if ok { // if a custom merge is defined, use it
			(*target)[k] = f((*target)[k], v, opt)
		} else { // otherwise just overwrite
			(*target)[k] = v
		}
	}
	return nil
}

type MergeFunc func(interface{}, interface{}, *Options) interface{}

var mergeFuncs = map[reflect.Type]MergeFunc{
	reflect.TypeOf(map[string]interface{}{}): func(t, s interface{}, o *Options) interface{} {
		return t
	},
}

func copyMap(m map[string]interface{}) map[string]interface{} {
	newMap := map[string]interface{}{}
	for k, v := range m {
		newMap[k] = v
	}

	return newMap
}

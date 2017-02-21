package merge

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("newFuncSelector", func() {
	var (
		fs *funcSelector
	)

	BeforeEach(func() {
		fs = newFuncSelector()
	})

	It("has the correct type mergefuncs", func() {
		mapMerge, mapOk := fs.typeFuncs[reflect.TypeOf(map[string]interface{}{})]
		Expect(mapOk).To(BeTrue())
		Expect(mapMerge).ToNot(BeNil())

		sliceMerge, sliceOK := fs.typeFuncs[reflect.TypeOf([]interface{}{})]
		Expect(sliceOK).To(BeTrue())
		Expect(sliceMerge).ToNot(BeNil())
	})

	It("has kind mergeFuncs map", func() {
		Expect(fs.kindFuncs).ToNot(BeNil())
	})

	It("has default mergeFunc", func() {
		Expect(fs.defaultFunc).ToNot(BeNil())
	})
})

var _ = Describe("Set Merge Func", func() {
	type TestKey struct{}

	var (
		fs *funcSelector
	)

	BeforeEach(func() {
		fs = newFuncSelector()
	})

	Context("Type Func", func() {
		It("adds the func correctly", func() {
			stubReturns := "uniqe string"
			t := reflect.TypeOf(TestKey{})
			fs.SetTypeMergeFunc(t, newMergeFuncStub(stubReturns))
			returned, _ := fs.typeFuncs[t](nil, nil, NewOptions())
			Expect(returned).To(Equal(stubReturns))
		})
	})

	Context("Kind Func", func() {
		It("adds the func correctly", func() {
			stubReturns := "uniqe string"
			k := reflect.TypeOf(TestKey{}).Kind()
			fs.SetKindMergeFunc(k, newMergeFuncStub(stubReturns))
			returned, _ := fs.kindFuncs[k](nil, nil, NewOptions())

			Expect(returned).To(Equal(stubReturns))
		})
	})

	Context("Default Func", func() {
		It("adds the func correctly", func() {
			stubReturns := "uniqe string"
			fs.SetDefaultMergeFunc(newMergeFuncStub(stubReturns))
			returned, _ := fs.defaultFunc(nil, nil, NewOptions())

			Expect(returned).To(Equal(stubReturns))
		})
	})

	Context("nil merge func maps", func() {
		It("does not panic", func() {
			fs = &funcSelector{}
			stubReturns := "uniqe string"
			f := newMergeFuncStub(stubReturns)
			t := reflect.TypeOf(TestKey{})
			k := reflect.TypeOf(TestKey{}).Kind()

			fs.SetTypeMergeFunc(t, f)
			fs.SetKindMergeFunc(k, f)
			fs.SetDefaultMergeFunc(f)
		})
	})
})

var _ = Describe("GetFunc", func() {
	type TestKey struct{}

	var (
		fs  *funcSelector
		key *TestKey
	)

	const (
		typeStubReturns    = "type"
		kindStubReturns    = "kind"
		defaultStubReturns = "default"
	)

	BeforeEach(func() {
		fs = &funcSelector{}
		key = &TestKey{}
	})

	Context("Type Func is defined", func() {
		BeforeEach(func() {
			fs.SetTypeMergeFunc(reflect.TypeOf(key), newMergeFuncStub(typeStubReturns))
		})

		It("gets the func", func() {
			f := fs.GetFunc(key)
			returned, _ := f(nil, nil, NewOptions())
			Expect(returned).To(Equal(typeStubReturns))
		})

		Context("kind func is also defined", func() {
			It("choses the type func", func() {
				fs.SetKindMergeFunc(reflect.TypeOf(key).Kind(), newMergeFuncStub(kindStubReturns))
				f := fs.GetFunc(key)
				returned, _ := f(nil, nil, NewOptions())
				Expect(returned).To(Equal(typeStubReturns))
			})
		})
	})

	Context("no type func defined", func() {
		Context("kind func is defined", func() {
			It("choses the kind func", func() {
				fs.SetKindMergeFunc(reflect.TypeOf(key).Kind(), newMergeFuncStub(kindStubReturns))
				f := fs.GetFunc(key)
				returned, _ := f(nil, nil, NewOptions())
				Expect(returned).To(Equal(kindStubReturns))
			})
		})

		Context("no kind func defined", func() {
			Context("default func defined", func() {
				It("choses the default func", func() {
					fs.SetDefaultMergeFunc(newMergeFuncStub(defaultStubReturns))
					f := fs.GetFunc(key)
					returned, _ := f(nil, nil, NewOptions())
					Expect(returned).To(Equal(defaultStubReturns))
				})
			})

			Context("no default func defined", func() {
				It("choses the global default func", func() {
					f := fs.GetFunc(key)
					returned, _ := f("a", "b", NewOptions())
					Expect(returned).To(Equal("b"))
				})
			})
		})
	})

	Context("no merge funcs defined", func() {
		It("returns defaultMergeFunc", func() {
			f := fs.GetFunc(key)
			Expect(f).ToNot(BeNil())
			merged, _ := f("a", "b", NewOptions())
			Expect(merged).To(Equal("b"))
		})
	})
})

func newMergeFuncStub(s string) MergeFunc {
	return func(interface{}, interface{}, *Options) (interface{}, error) {
		return s, nil
	}
}

var _ = Describe("defaultMergeFunc", func() {
	var opts *Options

	BeforeEach(func() {
		opts = &Options{}
	})

	Context("overwrite true", func() {
		It("returns source", func() {
			opts.Overwrite = true
			merged, err := defaultMergeFunc("a", "b", opts)
			Expect(err).ToNot(HaveOccurred())
			Expect(merged).To(Equal("b"))
		})
	})

	Context("overwrite false", func() {
		It("returns target", func() {
			opts.Overwrite = false
			merged, err := defaultMergeFunc("a", "b", opts)
			Expect(err).ToNot(HaveOccurred())
			Expect(merged).To(Equal("a"))
		})
	})
})

var _ = Describe("mergeMap", func() {
	var (
		targetMap, sourceMap map[string]interface{}
	)

	BeforeEach(func() {
		targetMap = map[string]interface{}{
			"A": "wrong",
			"B": 1,
			"C": map[string]interface{}{"foo": "unchanged", "bar": "orig"},
			"D": []interface{}{"unchanged", 0},
		}

		sourceMap = map[string]interface{}{
			"A": "correct",
			"B": 2,
			"C": map[string]interface{}{"bar": "newVal", "baz": "added"},
			"D": []interface{}{"added", 1},
		}
	})

	Context("happy path smoke test", func() {
		It("merges correctly", func() {
			merged, err := mergeMap(targetMap, sourceMap, NewOptions())

			Expect(err).ToNot(HaveOccurred())

			mergedMap, ok := merged.(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(mergedMap["A"]).To(Equal("correct"))
			Expect(mergedMap["B"]).To(Equal(2))

			subMap, ok := mergedMap["C"].(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(subMap["foo"]).To(Equal("unchanged"))
			Expect(subMap["bar"]).To(Equal("newVal"))
			Expect(subMap["baz"]).To(Equal("added"))

			subSlice, ok := mergedMap["D"]
			Expect(ok).To(BeTrue())
			Expect(subSlice).To(And(
				ContainElement("unchanged"),
				ContainElement(0),
				ContainElement("added"),
				ContainElement(1),
			))
		})
	})

	Context("overwrite is true", func() {
		It("overwrites", func() {
			opt := NewOptions()
			opt.Overwrite = true
			merged, err := mergeMap(targetMap, sourceMap, opt)

			Expect(err).ToNot(HaveOccurred())

			mergedMap, ok := merged.(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(mergedMap["A"]).To(Equal("correct"))
			Expect(mergedMap["B"]).To(Equal(2))

			subMap, ok := mergedMap["C"].(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(subMap["foo"]).To(Equal("unchanged"))
			Expect(subMap["bar"]).To(Equal("newVal"))
			Expect(subMap["baz"]).To(Equal("added"))
		})
	})

	Context("overwrite is false", func() {
		It("doesnt overwrite", func() {
			opt := NewOptions()
			opt.Overwrite = false
			merged, err := mergeMap(targetMap, sourceMap, opt)

			Expect(err).ToNot(HaveOccurred())

			mergedMap, ok := merged.(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(mergedMap["A"]).To(Equal("wrong"))
			Expect(mergedMap["B"]).To(Equal(1))

			subMap, ok := mergedMap["C"].(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(subMap["foo"]).To(Equal("unchanged"))
			Expect(subMap["bar"]).To(Equal("orig"))
			Expect(subMap["baz"]).To(Equal("added"))
		})
	})

	Context("empty target", func() {
		It("equals source", func() {
			merged, err := mergeMap(map[string]interface{}{}, sourceMap, NewOptions())

			Expect(err).ToNot(HaveOccurred())
			Expect(merged).To(Equal(sourceMap))
		})
	})

	Context("empty source", func() {
		It("equals target", func() {
			merged, err := mergeMap(targetMap, map[string]interface{}{}, NewOptions())

			Expect(err).ToNot(HaveOccurred())
			Expect(merged).To(Equal(targetMap))
		})
	})

	Context("nil target", func() {
		It("equals source", func() {
			merged, err := mergeMap(nil, sourceMap, NewOptions())

			Expect(err).ToNot(HaveOccurred())
			Expect(merged).To(Equal(sourceMap))
		})
	})

	Context("nil source", func() {
		It("equals target", func() {
			merged, err := mergeMap(targetMap, nil, NewOptions())

			Expect(err).ToNot(HaveOccurred())
			Expect(merged).To(Equal(targetMap))
		})
	})

	Context("mismatched field types", func() {
		It("errors", func() {
			targetMap["A"] = 0
			_, err := mergeMap(targetMap, sourceMap, NewOptions())

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("key 'A': Types do not match: int, string"))
		})

		Context("submap mismatch", func() {
			It("errors", func() {
				targetMap["C"] = map[string]interface{}{"bar": 0, "baz": "added"}
				_, err := mergeMap(targetMap, sourceMap, NewOptions())

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("key 'C': key 'bar': Types do not match: int, string"))
			})
		})
	})
})

var _ = Describe("mergeSlice", func() {
	var (
		targetSlice, sourceSlice []interface{}
	)

	BeforeEach(func() {
		targetSlice = []interface{}{3.6, "unchanged", 0}
		sourceSlice = []interface{}{1, "added", true}
	})

	Context("two populated slices", func() {
		It("merges them", func() {
			merged, err := mergeSlice(targetSlice, sourceSlice, NewOptions())
			Expect(err).ToNot(HaveOccurred())

			mergedSlice, ok := merged.([]interface{})
			Expect(ok).To(BeTrue())

			Expect(mergedSlice).To(And(
				ContainElement("unchanged"),
				ContainElement(0),
				ContainElement("added"),
				ContainElement(1),
				ContainElement(3.6),
				ContainElement(true),
			))
		})
	})

	Context("target slice is empty", func() {
		It("equals source", func() {
			merged, err := mergeSlice([]interface{}{}, sourceSlice, NewOptions())
			Expect(err).ToNot(HaveOccurred())

			mergedSlice, ok := merged.([]interface{})
			Expect(ok).To(BeTrue())

			Expect(len(mergedSlice)).To(Equal(len(sourceSlice)))
			Expect(mergedSlice).To(And(
				ContainElement("added"),
				ContainElement(1),
				ContainElement(true),
			))
			Expect(mergedSlice).ToNot(And(
				ContainElement("unchanged"),
				ContainElement(0),
				ContainElement(3.6),
			))
		})
	})

	Context("source slice is empty", func() {
		It("equals target", func() {
			merged, err := mergeSlice(targetSlice, []interface{}{}, NewOptions())
			Expect(err).ToNot(HaveOccurred())

			mergedSlice, ok := merged.([]interface{})
			Expect(ok).To(BeTrue())

			Expect(len(mergedSlice)).To(Equal(len(targetSlice)))
			Expect(mergedSlice).To(And(
				ContainElement("unchanged"),
				ContainElement(0),
				ContainElement(3.6),
			))
			Expect(mergedSlice).ToNot(And(
				ContainElement("added"),
				ContainElement(1),
				ContainElement(true),
			))
		})
	})

	Context("target slice is nil", func() {
		It("equals source", func() {
			merged, err := mergeSlice(nil, sourceSlice, NewOptions())
			Expect(err).ToNot(HaveOccurred())

			mergedSlice, ok := merged.([]interface{})
			Expect(ok).To(BeTrue())

			Expect(len(mergedSlice)).To(Equal(len(sourceSlice)))
			Expect(mergedSlice).To(And(
				ContainElement("added"),
				ContainElement(1),
				ContainElement(true),
			))
			Expect(mergedSlice).ToNot(And(
				ContainElement("unchanged"),
				ContainElement(0),
				ContainElement(3.6),
			))
		})
	})

	Context("source slice is nil", func() {
		It("equals target", func() {
			merged, err := mergeSlice(targetSlice, nil, NewOptions())
			Expect(err).ToNot(HaveOccurred())

			mergedSlice, ok := merged.([]interface{})
			Expect(ok).To(BeTrue())

			Expect(len(mergedSlice)).To(Equal(len(targetSlice)))
			Expect(mergedSlice).To(And(
				ContainElement("unchanged"),
				ContainElement(0),
				ContainElement(3.6),
			))
			Expect(mergedSlice).ToNot(And(
				ContainElement("added"),
				ContainElement(1),
				ContainElement(true),
			))
		})
	})

	Context("both slices are empty", func() {
		It("returns empty slice", func() {
			merged, err := mergeSlice([]interface{}{}, []interface{}{}, NewOptions())
			Expect(err).ToNot(HaveOccurred())

			mergedSlice, ok := merged.([]interface{})
			Expect(ok).To(BeTrue())

			Expect(mergedSlice).To(BeEmpty())
		})

		It("returns empty slice", func() {
			merged, err := mergeSlice(nil, nil, NewOptions())
			Expect(err).ToNot(HaveOccurred())

			mergedSlice, ok := merged.([]interface{})
			Expect(ok).To(BeTrue())

			Expect(mergedSlice).To(BeEmpty())
		})

	})
})

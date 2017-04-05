package merge

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
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

		structMerge, structOK := fs.kindFuncs[reflect.Struct]
		Expect(structOK).To(BeTrue())
		Expect(structMerge).ToNot(BeNil())
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

var _ = Describe("mergeStruct", func() {
	type Foo struct {
		Name    string
		Size    int
		Special bool
	}

	var (
		targetStruct, sourceStruct Foo
	)

	Context("two populated structs", func() {
		BeforeEach(func() {
			targetStruct = Foo{
				Name:    "target",
				Size:    2,
				Special: false,
			}

			sourceStruct = Foo{
				Name:    "source",
				Size:    4,
				Special: true,
			}
		})

		It("merges correctly", func() {
			merged, err := mergeStruct(&targetStruct, &sourceStruct, NewOptions())
			Expect(err).ToNot(HaveOccurred())
			Expect(merged).ToNot(BeNil())
			mergedStruct, ok := merged.(Foo)
			Expect(ok).To(BeTrue())
			Expect(mergedStruct.Name).To(Equal(sourceStruct.Name))
			Expect(mergedStruct.Size).To(Equal(sourceStruct.Size))
			Expect(mergedStruct.Special).To(Equal(sourceStruct.Special))
			Expect(targetStruct.Name).To(Equal("target"))
		})

		DescribeTable("pointer combinations",
			func(f func() (interface{}, error)) {
				merged, err := f()
				Expect(err).ToNot(HaveOccurred())
				Expect(merged).ToNot(BeNil())
				mergedStruct, ok := merged.(Foo)
				Expect(ok).To(BeTrue())
				Expect(mergedStruct.Name).To(Equal(sourceStruct.Name))
				Expect(mergedStruct.Size).To(Equal(sourceStruct.Size))
				Expect(mergedStruct.Special).To(Equal(sourceStruct.Special))
				Expect(targetStruct.Name).To(Equal("target"))
			},
			Entry(
				"pointer: T:n S:n",
				func() (interface{}, error) {
					return mergeStruct(targetStruct, sourceStruct, NewOptions())
				},
			),
			Entry(
				"pointer: T:y S:y",
				func() (interface{}, error) {
					return mergeStruct(&targetStruct, &sourceStruct, NewOptions())
				},
			),
			Entry(
				"pointer: T:y S:n",
				func() (interface{}, error) {
					return mergeStruct(&targetStruct, sourceStruct, NewOptions())
				},
			),
			Entry(
				"pointer: T:n S:y",
				func() (interface{}, error) {
					return mergeStruct(targetStruct, &sourceStruct, NewOptions())
				},
			),
		)
	})

	Context("partially populated", func() {
		Context("target is empty", func() {
			It("returns source", func() {
				merged, err := mergeStruct(Foo{}, sourceStruct, NewOptions())
				Expect(err).ToNot(HaveOccurred())
				Expect(merged).ToNot(BeNil())
				mergedStruct, ok := merged.(Foo)
				Expect(ok).To(BeTrue())
				Expect(mergedStruct.Name).To(Equal(sourceStruct.Name))
				Expect(mergedStruct.Size).To(Equal(sourceStruct.Size))
				Expect(mergedStruct.Special).To(Equal(sourceStruct.Special))
			})
		})

		Context("source is empty", func() {
			It("returns empty", func() {
				emptyFoo := Foo{}
				merged, err := mergeStruct(targetStruct, emptyFoo, NewOptions())
				Expect(err).ToNot(HaveOccurred())
				Expect(merged).ToNot(BeNil())
				mergedStruct, ok := merged.(Foo)
				Expect(ok).To(BeTrue())
				Expect(mergedStruct.Name).To(Equal(emptyFoo.Name))
				Expect(mergedStruct.Size).To(Equal(emptyFoo.Size))
				Expect(mergedStruct.Special).To(Equal(emptyFoo.Special))
			})
		})
	})

	Context("fields that need to be merged", func() {
		Context("slice fields", func() {
			type Baz struct {
				Slice []interface{}
			}
			var targetBaz, sourceBaz Baz

			BeforeEach(func() {
				targetBaz = Baz{
					Slice: []interface{}{"unchanged", 0},
				}
				sourceBaz = Baz{
					Slice: []interface{}{"added", 1},
				}
			})

			It("merges them", func() {
				merged, err := mergeStruct(targetBaz, sourceBaz, NewOptions())
				Expect(err).ToNot(HaveOccurred())
				mergedStruct, ok := merged.(Baz)
				Expect(ok).To(BeTrue())
				Expect(mergedStruct.Slice).To(And(
					ContainElement("unchanged"),
					ContainElement("added"),
					ContainElement(0),
					ContainElement(1),
				))
			})
		})

		Context("pointer fields", func() {
			type Baz struct {
				Ptr   *string
				Slice *[]interface{}
			}
			var targetBaz, sourceBaz Baz

			BeforeEach(func() {
				t := "target"
				s := "source"
				targetBaz = Baz{
					Ptr:   &t,
					Slice: &[]interface{}{"unchanged", 0},
				}
				sourceBaz = Baz{
					Ptr:   &s,
					Slice: &[]interface{}{"added", 1},
				}
			})

			It("handles them properly", func() {
				merged, err := mergeStruct(targetBaz, sourceBaz, NewOptions())
				Expect(err).ToNot(HaveOccurred())
				mergedStruct, ok := merged.(Baz)
				Expect(ok).To(BeTrue())
				Expect(*mergedStruct.Ptr).To(Equal("source"))
				Expect(*mergedStruct.Slice).To(And(
					ContainElement("added"),
					ContainElement(1),
				))
			})

			Context("source field is nil", func() {
				BeforeEach(func() {
					sourceBaz.Ptr = nil
					sourceBaz.Slice = nil
				})

				It("handles the nil", func() {
					By("with overwrite")

					merged, err := mergeStruct(targetBaz, sourceBaz, NewOptions())
					Expect(err).ToNot(HaveOccurred())
					mergedStruct, ok := merged.(Baz)
					Expect(ok).To(BeTrue())
					Expect(*mergedStruct.Ptr).To(Equal("target"))
					Expect(*mergedStruct.Slice).To(And(
						ContainElement("unchanged"),
						ContainElement(0),
					))
				})
			})
		})
	})

	// These are tested through the merge() func because that is what protects against panics
	Context("invalid entries", func() {
		var opt *Options

		BeforeEach(func() {
			opt = NewOptions()
			opt.MergeFuncs.SetKindMergeFunc(reflect.Struct, mergeStruct)
		})

		Context("nils", func() {
			Context("target nil", func() {
				It("return source", func() {
					merged, err := merge(nil, sourceStruct, opt)
					Expect(err).ToNot(HaveOccurred())
					Expect(merged).ToNot(BeNil())
					mergedStruct, ok := merged.(Foo)
					Expect(ok).To(BeTrue())
					Expect(mergedStruct.Name).To(Equal(sourceStruct.Name))
					Expect(mergedStruct.Size).To(Equal(sourceStruct.Size))
					Expect(mergedStruct.Special).To(Equal(sourceStruct.Special))
				})
			})

			Context("source nil", func() {
				It("return target", func() {
					merged, err := merge(targetStruct, nil, opt)
					Expect(err).ToNot(HaveOccurred())
					Expect(merged).ToNot(BeNil())
					mergedStruct, ok := merged.(Foo)
					Expect(ok).To(BeTrue())
					Expect(mergedStruct.Name).To(Equal(targetStruct.Name))
					Expect(mergedStruct.Size).To(Equal(targetStruct.Size))
					Expect(mergedStruct.Special).To(Equal(targetStruct.Special))
				})

			})
		})

		Context("struct types do not match", func() {
			It("errors", func() {
				type Bar struct {
					Baz float64
				}
				merged, err := merge(targetStruct, Bar{}, opt)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Types do not match: merge.Foo, merge.Bar"))
				Expect(merged).To(BeNil())
			})
		})
	})

	Context("non-struct type", func() {
		It("errors", func() {
			merged, err := mergeStruct(targetStruct, "a string", NewOptions())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("got non-struct kind (tagret: struct; source: string)"))
			Expect(merged).To(BeNil())
		})
	})

	Context("private fields on struct", func() {
		type Baz struct {
			Public  string
			private string
		}
		var targetBaz, sourceBaz Baz

		BeforeEach(func() {
			targetBaz = Baz{
				Public:  "target",
				private: "target",
			}
			sourceBaz = Baz{
				Public:  "source",
				private: "source",
			}
		})

		It("errors", func() {
			merged, err := mergeStruct(targetBaz, sourceBaz, NewOptions())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("problem with field(private) valid: true; can set: false"))
			Expect(merged).To(BeNil())
		})
	})

	Context("merge error on field", func() {
		type Baz struct {
			Foo interface{}
		}
		var targetBaz, sourceBaz Baz

		BeforeEach(func() {
			targetBaz = Baz{
				Foo: "string",
			}
			sourceBaz = Baz{
				Foo: 1,
			}
		})

		It("returns error", func() {
			merged, err := mergeStruct(targetBaz, sourceBaz, NewOptions())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to merge field `Baz.Foo`: Types do not match: string, int"))
			Expect(merged).To(BeNil())
		})
	})

	Context("receives wrong type on merge", func() {
		type Baz struct {
			Bar string
		}
		var (
			targetBaz, sourceBaz Baz
			opt                  *Options
		)

		BeforeEach(func() {
			opt = NewOptions()
			// merge func for string returns an int. wont be able to set the field to an int
			opt.MergeFuncs.SetKindMergeFunc(reflect.String, func(t, s interface{}, o *Options) (interface{}, error) {
				return 0, nil
			})

			targetBaz = Baz{
				Bar: "target",
			}
			sourceBaz = Baz{
				Bar: "source",
			}
		})

		It("errors", func() {
			merged, err := mergeStruct(targetBaz, sourceBaz, opt)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("types dont match string <> int"))
			Expect(merged).To(BeNil())
		})
	})
})

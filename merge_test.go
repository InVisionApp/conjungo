package merge

import (
	"encoding/json"
	"errors"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Options", func() {
	var (
		testOpts *Options
	)
	Context("happy path", func() {
		BeforeEach(func() {
			testOpts = NewOptions()
		})

		It("not nil", func() {
			Expect(testOpts).ToNot(BeNil())
		})

		It("defailts to overwrite true", func() {
			Expect(testOpts.Overwrite).To(BeTrue())
		})

		It("sets up mergefuncs", func() {
			Expect(testOpts.MergeFuncs).ToNot(BeNil())
		})
	})
})

var _ = Describe("Merge", func() {
	var (
		targetMap, sourceMap map[string]interface{}
	)

	Context("happy path smoke test", func() {
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
				"C": map[string]interface{}{"bar": "newVal", "safe": "added"},
				"D": []interface{}{"added", 1},
			}
		})

		It("merges correctly", func() {
			merged, err := Merge(targetMap, sourceMap, NewOptions())

			Expect(err).ToNot(HaveOccurred())
			newMap, ok := merged.(map[string]interface{})
			Expect(ok).To(BeTrue())

			jsonB, errJson := json.Marshal(newMap)
			Expect(errJson).ToNot(HaveOccurred())

			expectedJSON := `{
			  "A": "correct",
			  "B": 2,
			  "C": {
				"bar": "newVal",
				"foo": "unchanged",
				"safe": "added"
			  },
			  "D": [
				"unchanged",
				0,
				"added",
				1
			  ]
			}`
			Expect(jsonB).To(MatchJSON(expectedJSON))
		})

		It("accepts nil options", func() {
			_, err := Merge(targetMap, sourceMap, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("happy path - overwrite is false", func() {
		var (
			newMap map[string]interface{}
			err    error
		)

		BeforeEach(func() {
			targetMap = map[string]interface{}{
				"A": "original",
				"B": 1,
				"C": map[string]interface{}{"foo": "unchanged", "bar": "orig"},
				"D": []interface{}{"unchanged", 0},
			}

			sourceMap = map[string]interface{}{
				"A": "overwritten",
				"B": 2,
				"C": map[string]interface{}{"bar": "newVal", "baz": "added"},
				"D": []interface{}{"added", 1},
				"E": "inserted",
			}

			opts := NewOptions()
			opts.Overwrite = false

			var (
				merged interface{}
				ok     bool
			)

			merged, err = Merge(targetMap, sourceMap, opts)
			newMap, ok = merged.(map[string]interface{})
			Expect(ok).To(BeTrue())
		})

		It("does not error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("does not overwrite a top level string", func() {
			Expect(newMap["A"]).To(Equal("original"))
		})

		It("does not overwrite a top level int", func() {
			Expect(newMap["B"]).To(Equal(1))
		})

		It("inserts a new top level string", func() {
			Expect(newMap["E"]).To(Equal("inserted"))
		})

		Context("sub map", func() {
			var (
				newSubMap map[string]interface{}
				ok        bool
			)

			JustBeforeEach(func() {
				newSubMap, ok = newMap["C"].(map[string]interface{})
				Expect(ok).To(BeTrue())
			})

			It("does not overwrite a sub map value", func() {
				Expect(newSubMap["bar"]).To(Equal("orig"))
			})

			It("inserts a new sub map value", func() {
				Expect(newSubMap["baz"]).To(Equal("added"))
			})

			It("maintains unduplicated values", func() {
				Expect(newSubMap["foo"]).To(Equal("unchanged"))
			})
		})

		Context("sub slice", func() {
			It("merges properly", func() {
				newSubSlice, ok := newMap["D"].([]interface{})
				Expect(ok).To(BeTrue())

				Expect(len(newSubSlice)).To(Equal(4))
				Expect(newSubSlice).To(ContainElement("unchanged"))
				Expect(newSubSlice).To(ContainElement(0))
				Expect(newSubSlice).To(ContainElement("added"))
				Expect(newSubSlice).To(ContainElement(1))
			})
		})
	})

	Context("happy path specific types", func() {
		DescribeTable("basic types",
			func(target, source interface{}) {
				merged, err := Merge(target, source, NewOptions())

				Expect(err).ToNot(HaveOccurred())
				Expect(merged).To(Equal(source))
			},
			Entry("overwrites string",
				"wrong",
				"correct",
			),
			Entry("overwrites int",
				0,
				1,
			),
			Entry("overwrites float",
				0.0,
				1.0,
			),
			Entry("overwrites bool",
				true,
				false,
			),
		)

		Context("merge map", func() {
			const (
				testKey = "theKey"
			)

			It("merges correctly", func() {
				targetMap = map[string]interface{}{
					testKey: map[string]interface{}{"foo": "unchanged", "bar": "orig"},
				}

				sourceMap = map[string]interface{}{
					testKey: map[string]interface{}{"bar": "newVal", "baz": "added"},
				}

				merged, err := Merge(targetMap, sourceMap, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				mergedMap, ok := merged.(map[string]interface{})
				Expect(ok).To(BeTrue())

				dataMap, ok := mergedMap[testKey].(map[string]interface{})
				Expect(ok).To(BeTrue())

				Expect(dataMap["foo"]).To(Equal("unchanged"))
				Expect(dataMap["bar"]).To(Equal("newVal"))
				Expect(dataMap["baz"]).To(Equal("added"))
			})
		})

		Context("nil target", func() {
			It("merges correctly", func() {
				source := "bar"

				merged, err := Merge(nil, source, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				Expect(merged).To(Equal("bar"))
			})
		})

		Context("map within a map", func() {
			//TODO
		})

		Context("merge slice", func() {
			It("merges correctly", func() {
				target := []interface{}{"unchanged", 0}

				source := []interface{}{"added", 1}

				merged, err := Merge(target, source, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				dataSlice, ok := merged.([]interface{})
				Expect(ok).To(BeTrue())

				Expect(len(dataSlice)).To(Equal(4))
				Expect(dataSlice).To(ContainElement("unchanged"))
				Expect(dataSlice).To(ContainElement(0))
				Expect(dataSlice).To(ContainElement("added"))
				Expect(dataSlice).To(ContainElement(1))
			})
		})

		Context("nil source value", func() {
			It("doesnt error", func() {
				target := "foo"

				merged, err := Merge(target, nil, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				origVal, ok := merged.(string)
				Expect(ok).To(BeTrue())
				Expect(origVal).To(Equal("foo"))
			})
		})
	})

	Context("failure modes", func() {
		Context("merge func returns error", func() {
			It("returns an error", func() {
				target := errors.New("some err")
				source := errors.New("other err")

				opts := NewOptions()
				// define a merge func that always errors for the error type
				opts.MergeFuncs.SetTypeMergeFunc(
					reflect.TypeOf(errors.New("")),
					func(t, s interface{}, o *Options) (interface{}, error) {
						return nil, errors.New("returns error")
					},
				)
				_, err := Merge(target, source, opts)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("returns error"))
			})
		})

		Context("type mismatch", func() {
			It("returns an error", func() {
				target := 0
				source := ""

				_, err := Merge(target, source, NewOptions())

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Types do not match"))
			})
		})

		Context("type mismatch deeper in the tree", func() {
			It("returns an error", func() {
				testKey := "theKey"

				targetMap = map[string]interface{}{
					testKey: 0,
				}

				sourceMap = map[string]interface{}{
					testKey: "",
				}

				_, err := Merge(targetMap, sourceMap, NewOptions())

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Types do not match"))
			})
		})
	})
})

var _ = Describe("MergeMapStrIFace", func() {
	var (
		targetMap, sourceMap map[string]interface{}
		newMap               map[string]interface{}
		err                  error
	)

	BeforeEach(func() {
		targetMap = map[string]interface{}{
			"A": "original",
			"B": 1,
			"C": map[string]interface{}{"foo": "unchanged", "bar": "orig"},
			"D": []interface{}{"unchanged", 0},
		}

		sourceMap = map[string]interface{}{
			"A": "overwritten",
			"B": 2,
			"C": map[string]interface{}{"bar": "newVal", "baz": "added"},
			"D": []interface{}{"added", 1},
			"E": "inserted",
		}
	})

	Context("happy path", func() {
		JustBeforeEach(func() {
			opts := NewOptions()
			opts.Overwrite = false

			newMap, err = MergeMapStrIface(targetMap, sourceMap, opts)
		})

		It("does not error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("does not overwrite a top level string", func() {
			Expect(newMap["A"]).To(Equal("original"))
		})

		It("does not overwrite a top level int", func() {
			Expect(newMap["B"]).To(Equal(1))
		})

		It("inserts a new top level string", func() {
			Expect(newMap["E"]).To(Equal("inserted"))
		})
	})

	Context("Merge fails", func() {
		BeforeEach(func() {
			targetMap["F"] = errors.New("some err")
			sourceMap["F"] = errors.New("other err")

			opts := NewOptions()
			// define a merge func that always errors for the error type
			opts.MergeFuncs.SetTypeMergeFunc(
				reflect.TypeOf(errors.New("")),
				func(t, s interface{}, o *Options) (interface{}, error) {
					return nil, errors.New("returns error")
				},
			)

			newMap, err = MergeMapStrIface(targetMap, sourceMap, opts)
		})

		It("returns an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("returns error"))
		})
	})

	Context("merge returns non-map", func() {
		BeforeEach(func() {
			targetMap["F"] = errors.New("some err")
			sourceMap["F"] = errors.New("other err")

			opts := NewOptions()
			// define a merge func that returns wrong type
			opts.MergeFuncs.SetTypeMergeFunc(
				reflect.TypeOf(map[string]interface{}{}),
				func(t, s interface{}, o *Options) (interface{}, error) {
					return "a string", nil
				},
			)

			_, err = MergeMapStrIface(targetMap, sourceMap, opts)
		})

		It("returns an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Merge failed. Expected map[string]interface{} but got string"))
		})
	})
})

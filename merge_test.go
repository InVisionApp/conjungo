package conjungo

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

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
	type Foo struct {
		Name    string
		Size    int
		Special bool
		Tags    []interface{}
	}

	var (
		targetMap, sourceMap map[string]interface{}
	)

	Context("happy path smoke test", func() {
		Context("maps", func() {
			BeforeEach(func() {
				targetMap = map[string]interface{}{
					"A": "wrong",
					"B": 1,
					"C": map[string]interface{}{"foo": "unchanged", "bar": "orig"},
					"D": []interface{}{"unchanged", 0},
					"E": Foo{Name: "target", Size: 1, Special: false, Tags: []interface{}{"unchanged", 0}},
				}

				sourceMap = map[string]interface{}{
					"A": "correct",
					"B": 2,
					"C": map[string]interface{}{"bar": "newVal", "safe": "added"},
					"D": []interface{}{"added", 1},
					"E": Foo{Name: "source", Size: 3, Special: true, Tags: []interface{}{"added", 1}},
				}
			})

			It("merges correctly", func() {
				err := Merge(&targetMap, sourceMap, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				jsonB, errJson := json.Marshal(targetMap)
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
			  ],
			  "E": {
				"Name": "source",
				"Size": 3,
				"Special": true,
				"Tags": [
				  "unchanged",
				  0,
				  "added",
				  1
				]
			  }
			}`
				Expect(jsonB).To(MatchJSON(expectedJSON))
			})

			It("accepts nil options", func() {
				err := Merge(&targetMap, sourceMap, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("structs", func() {
			var (
				targetStruct, sourceStruct Foo
			)

			BeforeEach(func() {
				targetStruct = Foo{
					Name:    "target",
					Size:    1,
					Special: false,
					Tags:    []interface{}{"unchanged", 0},
				}

				sourceStruct = Foo{
					Name:    "source",
					Size:    3,
					Special: true,
					Tags:    []interface{}{"added", 1},
				}
			})

			It("merges correctly", func() {
				err := Merge(&targetStruct, sourceStruct, NewOptions())

				Expect(err).ToNot(HaveOccurred())
				//newMap, ok := merged.(map[string]interface{})
				//Expect(ok).To(BeTrue())

				jsonB, errJson := json.Marshal(targetStruct)
				Expect(errJson).ToNot(HaveOccurred())

				expectedJSON := `{
				"Name": "source",
				"Size": 3,
				"Special": true,
				"Tags": [
				  "unchanged",
				  0,
				  "added",
				  1
				]
			}`
				Expect(jsonB).To(MatchJSON(expectedJSON))
			})
		})
	})

	Context("happy path - overwrite is false", func() {
		var (
			err error
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

			err = Merge(&targetMap, sourceMap, opts)
		})

		It("does not error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("does not overwrite a top level string", func() {
			Expect(targetMap["A"]).To(Equal("original"))
		})

		It("does not overwrite a top level int", func() {
			Expect(targetMap["B"]).To(Equal(1))
		})

		It("inserts a new top level string", func() {
			Expect(targetMap["E"]).To(Equal("inserted"))
		})

		Context("sub map", func() {
			var (
				newSubMap map[string]interface{}
				ok        bool
			)

			JustBeforeEach(func() {
				newSubMap, ok = targetMap["C"].(map[string]interface{})
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
				newSubSlice, ok := targetMap["D"].([]interface{})
				Expect(ok).To(BeTrue())

				Expect(len(newSubSlice)).To(Equal(4))
				Expect(newSubSlice).To(ContainElement("unchanged"))
				Expect(newSubSlice).To(ContainElement(0))
				Expect(newSubSlice).To(ContainElement("added"))
				Expect(newSubSlice).To(ContainElement(1))
			})
		})

		Context("pointers", func() {
			Context("target is already pointer", func() {
				//TODO: implement
				//interface of pointer value
			})

			Context("pointer and value", func() {
				Context("target is pointer source is value", func() {
					//TODO: implement
				})

				Context("target is value source is pointer", func() {
					//TODO: implement
				})
			})
		})
	})

	Context("happy path specific types", func() {
		DescribeTable("basic types",
			func(target, source interface{}) {
				err := Merge(&target, source, NewOptions())

				Expect(err).ToNot(HaveOccurred())
				Expect(target).To(Equal(source))
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
			Entry("time",
				time.Now(),
				time.Now().Add(time.Hour),
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

				err := Merge(&targetMap, sourceMap, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				dataMap, ok := targetMap[testKey].(map[string]interface{})
				Expect(ok).To(BeTrue())

				Expect(dataMap["foo"]).To(Equal("unchanged"))
				Expect(dataMap["bar"]).To(Equal("newVal"))
				Expect(dataMap["baz"]).To(Equal("added"))
			})
		})

		Context("nil target", func() {
			It("merges correctly", func() {
				source := "bar"
				var target string

				err := Merge(&target, source, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				Expect(target).To(Equal("bar"))
			})
		})

		Context("nil source value", func() {
			It("doesnt error", func() {
				target := "foo"

				err := Merge(&target, nil, NewOptions())
				Expect(err).ToNot(HaveOccurred())
				Expect(target).To(Equal("foo"))
			})
		})

		Context("nil source and target value", func() {
			It("doesnt error", func() {
				var target int
				err := Merge(&target, nil, NewOptions())
				Expect(err).ToNot(HaveOccurred())
				Expect(target).To(BeZero())
			})
		})

		Context("nil pointer value", func() {
			Context("source nil", func() {
				It("doesnt error", func() {
					target := "foo"
					var source *string

					err := Merge(&target, source, NewOptions())
					Expect(err).ToNot(HaveOccurred())
					Expect(target).To(Equal("foo")) // unchanged
				})
			})

			Context("target nil", func() {
				It("doesnt error", func() {
					source := "foo"
					var target *string

					err := Merge(target, &source, NewOptions())
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("target can not be zero value"))
				})
			})
		})

		Context("merge slice", func() {
			It("merges correctly", func() {
				target := []interface{}{"unchanged", 0}

				source := []interface{}{"added", 1}

				err := Merge(&target, source, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				Expect(len(target)).To(Equal(4))
				Expect(target).To(ContainElement("unchanged"))
				Expect(target).To(ContainElement(0))
				Expect(target).To(ContainElement("added"))
				Expect(target).To(ContainElement(1))
			})
		})

		Context("merge struct", func() {
			type Thing struct {
				Foo string
			}

			It("merges correctly", func() {
				target := Thing{Foo: "bar"}
				source := Thing{Foo: "baz"}

				err := Merge(&target, source, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				Expect(target.Foo).To(Equal("baz"))
			})

			It("target struct and source value merges correctly", func() {
				target := Thing{Foo: "bar"}
				source := reflect.ValueOf(Thing{Foo: "baz"})

				err := Merge(&target, source, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				Expect(target.Foo).To(Equal("baz"))
			})

			It("target value and source struct merges correctly", func() {
				target := Thing{Foo: "bar"}
				source := Thing{Foo: "baz"}

				tVal := reflect.ValueOf(&target)
				err := Merge(tVal, source, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				Expect(target.Foo).To(Equal("baz"))
			})

			It("both values merges correctly", func() {
				target := Thing{Foo: "bar"}
				source := reflect.ValueOf(Thing{Foo: "baz"})

				err := Merge(reflect.ValueOf(&target), source, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				Expect(target.Foo).To(Equal("baz"))
			})

			It("pointer to target value errors", func() {
				target := Thing{Foo: "bar"}
				source := reflect.ValueOf(Thing{Foo: "baz"})

				tVal := reflect.ValueOf(&target)
				err := Merge(tVal, source, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				Expect(target.Foo).To(Equal("baz"))
			})
		})
	})

	Context("failure modes", func() {
		Context("target is not a pointer", func() {
			It("returns error", func() {
				err := Merge("foo", "bar", NewOptions())

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("target must be a pointer"))

			})
		})
		Context("merge func returns error", func() {
			It("returns an error", func() {
				target := errors.New("some err")
				source := errors.New("other err")

				opts := NewOptions()
				// define a merge func that always errors for the error type
				opts.MergeFuncs.SetTypeMergeFunc(
					reflect.TypeOf(errors.New("")),
					erroringMergeFunc,
				)
				err := Merge(&target, source, opts)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("returns error"))
			})
		})

		Context("type mismatch", func() {
			It("returns an error", func() {
				target := 0
				source := ""

				err := Merge(&target, source, NewOptions())

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

				err := Merge(&targetMap, sourceMap, NewOptions())

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Types do not match"))
			})
		})

		Context("merge returns wrong type", func() {
			type Bar struct{}

			It("returns an error", func() {
				target := Bar{}
				source := Bar{}

				opts := NewOptions()
				// define a merge func that returns wrong type
				opts.MergeFuncs.SetTypeMergeFunc(
					reflect.TypeOf(Bar{}),
					func(t, s reflect.Value, o *Options) (reflect.Value, error) {
						return reflect.ValueOf("a string"), nil
					},
				)

				err := Merge(&target, source, opts)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Merge failed: expected merged result to be conjungo.Bar but got string"))
			})
		})

	})
})

//TODO: maybe some of these tests are duplicates. Dedup them sometime.
var _ = Describe("MergeMapStrIFace", func() {
	var (
		targetMap, sourceMap map[string]interface{}
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

			err = Merge(&targetMap, sourceMap, opts)
		})

		It("does not error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("does not overwrite a top level string", func() {
			Expect(targetMap["A"]).To(Equal("original"))
		})

		It("does not overwrite a top level int", func() {
			Expect(targetMap["B"]).To(Equal(1))
		})

		It("inserts a new top level string", func() {
			Expect(targetMap["E"]).To(Equal("inserted"))
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
				erroringMergeFunc,
			)

			err = Merge(&targetMap, sourceMap, opts)
		})

		It("returns an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("returns error"))
		})
	})
})

func erroringMergeFunc(t, s reflect.Value, o *Options) (reflect.Value, error) {
	return reflect.Value{}, errors.New("returns error")
}

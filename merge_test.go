package merge

import (
	"encoding/json"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"reflect"
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

		It("has the correct mergefuncs", func() {
			Expect(len(testOpts.mergeFuncs)).To(Equal(2))

			mapMerge, mapOk := testOpts.mergeFuncs[reflect.TypeOf(map[string]interface{}{})]
			Expect(mapOk).To(BeTrue())
			Expect(mapMerge).ToNot(BeNil())

			sliceMerge, sliceOK := testOpts.mergeFuncs[reflect.TypeOf([]interface{}{})]
			Expect(sliceOK).To(BeTrue())
			Expect(sliceMerge).ToNot(BeNil())
		})
	})
})

var _ = Describe("merge", func() {
	const (
		testKey = "theKey"
	)

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
			newMap, err := Merge(targetMap, sourceMap, NewOptions())

			Expect(err).ToNot(HaveOccurred())

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
	})

	Context("happy path specific types", func() {
		DescribeTable("basic types",
			func(target, source interface{}) {
				targetMap = map[string]interface{}{
					testKey: target,
				}

				sMap := map[string]interface{}{
					testKey: source,
				}

				newMap, err := Merge(targetMap, sMap, NewOptions())

				Expect(err).ToNot(HaveOccurred())
				Expect(newMap[testKey]).To(Equal(sMap[testKey]))
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
			It("merges correctly", func() {
				targetMap = map[string]interface{}{
					testKey: map[string]interface{}{"foo": "unchanged", "bar": "orig"},
				}

				sourceMap = map[string]interface{}{
					testKey: map[string]interface{}{"bar": "newVal", "baz": "added"},
				}

				mergedMap, err := Merge(targetMap, sourceMap, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				dataMap, ok := mergedMap[testKey].(map[string]interface{})
				Expect(ok).To(BeTrue())

				Expect(dataMap["foo"]).To(Equal("unchanged"))
				Expect(dataMap["bar"]).To(Equal("newVal"))
				Expect(dataMap["baz"]).To(Equal("added"))
			})
		})

		Context("merge slice", func() {
			It("merges correctly", func() {
				targetMap = map[string]interface{}{
					testKey: []interface{}{"unchanged", 0},
				}

				sourceMap = map[string]interface{}{
					testKey: []interface{}{"added", 1},
				}

				mergedMap, err := Merge(targetMap, sourceMap, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				dataSlice, ok := mergedMap[testKey].([]interface{})
				Expect(ok).To(BeTrue())

				Expect(dataSlice).To(ContainElement("unchanged"))
				Expect(dataSlice).To(ContainElement(0))
				Expect(dataSlice).To(ContainElement("added"))
				Expect(dataSlice).To(ContainElement(1))
			})
		})

		Context("nil source value", func() {
			It("doesnt error", func() {
				targetMap = map[string]interface{}{
					testKey: 0,
				}

				sourceMap = map[string]interface{}{
					testKey: nil,
				}

				mergedMap, err := Merge(targetMap, sourceMap, NewOptions())
				Expect(err).ToNot(HaveOccurred())

				origVal, ok := mergedMap[testKey].(int)
				Expect(ok).To(BeTrue())
				Expect(origVal).To(Equal(0))
			})
		})
	})

	Context("failure modes", func() {
		Context("merge func returns error", func() {
			It("returns an error", func() {
				targetMap = map[string]interface{}{
					testKey: errors.New("some err"),
				}

				sourceMap = map[string]interface{}{
					testKey: errors.New("other err"),
				}

				opts := NewOptions()
				// define a merge func that always errors
				opts.SetMergeFunc(
					reflect.TypeOf(errors.New("")),
					func(t, s interface{}, o *Options) (interface{}, error) {
						return nil, errors.New("returns error")
					},
				)
				_, err := Merge(targetMap, sourceMap, opts)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("returns error"))
			})
		})

		Context("type mismatch returns error", func() {
			It("returns an error", func() {
				targetMap = map[string]interface{}{
					testKey: 0,
				}

				sourceMap = map[string]interface{}{
					testKey: "",
				}

				_, err := Merge(targetMap, sourceMap, NewOptions())

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Types do not match for key"))
			})
		})
	})
})

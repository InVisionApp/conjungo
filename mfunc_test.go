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

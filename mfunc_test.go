package merge

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("funcSelector", func() {
	Context("newFuncSelector", func() {
		It("has the correct mergefuncs", func() {
			fs := newFuncSelector()

			mapMerge, mapOk := fs.typeFuncs[reflect.TypeOf(map[string]interface{}{})]
			Expect(mapOk).To(BeTrue())
			Expect(mapMerge).ToNot(BeNil())

			sliceMerge, sliceOK := fs.typeFuncs[reflect.TypeOf([]interface{}{})]
			Expect(sliceOK).To(BeTrue())
			Expect(sliceMerge).ToNot(BeNil())
		})
	})
})

package logger_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mortedecai/rivulets-of-go/logger"
)

var _ = Describe("Logger", func() {
	It("should return a Logger", func() {
		const name = "foo"
		l, err := logger.New(name, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(l).ToNot(BeNil())
	})
	It("should return a new Logger when WithName is used", func() {
		const name = "foo"
		const name2 = "bar"
		l, err := logger.New(name, true)
		Expect(err).ToNot(HaveOccurred())
		l2 := l.WithName(name2)
		Expect(l).ToNot(Equal(l2))
	})
})

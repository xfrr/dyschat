package mongo_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xfrr/dyschat/pkg/testing"
)

var _ = Describe("Mongo Pkg Suite", func() {
	var (
		ctx context.Context
	)
	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("Connect", func() {
		It("should connect to the database", func() {
			dc := testing.NewMongoContainer()
			database, err := dc.Start(ctx, "mongo-test")
			Expect(err).ToNot(HaveOccurred())
			Expect(database).ToNot(BeNil())
			Eventually(database.Client().Ping(ctx, nil)).Should(Succeed())
			dc.Purge()
		})
	})
})

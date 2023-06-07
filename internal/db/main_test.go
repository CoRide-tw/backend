package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var pgPool *pgxpool.Pool

var _ = BeforeSuite(func() {
	var err error
	dbUrl := os.Getenv("POSTGRES_DATABASE_URL")
	pgPool, err = pgxpool.New(context.Background(), dbUrl)
	Expect(err).NotTo(HaveOccurred())
	Expect(InitDBClient(pgPool)).To(Succeed())
})

var _ = AfterSuite(func() {
	pgPool.Close()
})

//func TestDB(t *testing.T) {
//	RegisterFailHandler(Fail)
//	RunSpecs(t, "DB Suite")
//}

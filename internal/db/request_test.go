package db

import (
	"context"
	"fmt"
	"github.com/CoRide-tw/backend/internal/config"
	"github.com/CoRide-tw/backend/internal/constants"
	. "github.com/CoRide-tw/backend/internal/errors/generated/dberr"
	"github.com/CoRide-tw/backend/internal/model"
	. "github.com/DenChenn/blunder/pkg/blunder"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"log"
)

var _ = Describe("Request", func() {
	// initialize db connection and env
	config.Env = config.LoadEnv()
	pgPool, err := pgxpool.New(context.Background(), config.Env.PostgresDatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer pgPool.Close()

	// test data
	existedRequests := []model.Request{
		{RiderId: 1, RouteId: 1},
		{RiderId: 1, RouteId: 2},
		{RiderId: 2, RouteId: 1},
	}

	BeforeEach(func() {
		for i, request := range existedRequests {
			err := pgPool.QueryRow(context.Background(), fmt.Sprintf(
				`INSERT INTO requests (rider_id, route_id) VALUES (%d, %d) RETURNING id;`,
				request.RiderId,
				request.RouteId,
			)).Scan(&existedRequests[i].Id)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	AfterEach(func() {
		for _, request := range existedRequests {
			_, err := pgPool.Exec(context.Background(), fmt.Sprintf(
				`DELETE FROM requests WHERE id = %d;`,
				request.Id,
			))
			Expect(err).NotTo(HaveOccurred())
		}
	})

	Describe("GetRequest", func() {
		var (
			request *model.Request
			id      int32
			err     error
		)

		JustBeforeEach(func() {
			request, err = GetRequest(id)
		})

		When("request exists in database", func() {
			BeforeEach(func() {
				id = existedRequests[0].Id
			})

			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(request.Id).To(Equal(int32(1)))
				Expect(request.RiderId).To(Equal(int32(1)))
				Expect(request.RouteId).To(Equal(int32(1)))
			})
		})

		When("request does not exist in database", func() {
			BeforeEach(func() {
				id = 0
			})

			It("fails", func() {
				Expect(err).To(MatchError(ErrRequestNotFound))
				Expect(request).To(BeNil())
			})
		})
	})

	Describe("ListRequestsByRiderId", func() {
		var (
			requests []*model.Request
			riderId  int32
			err      error
		)

		JustBeforeEach(func() {
			requests, err = ListRequestsByRiderId(riderId)
		})

		When("requests exist in database", func() {
			BeforeEach(func() {
				riderId = existedRequests[0].RiderId
			})

			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(len(requests)).To(Equal(2))
			})
		})

		When("requests do not exist in database", func() {
			BeforeEach(func() {
				riderId = 0
			})

			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(len(requests)).To(Equal(0))
			})
		})
	})

	Describe("ListRequestsByRouteId", func() {
		var (
			requests []*model.Request
			routeId  int32
			err      error
		)

		JustBeforeEach(func() {
			requests, err = ListRequestsByRouteId(routeId)
		})

		When("requests exist in database", func() {
			BeforeEach(func() {
				routeId = existedRequests[0].RouteId
			})

			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(len(requests)).To(Equal(2))
			})
		})

		When("requests do not exist in database", func() {
			BeforeEach(func() {
				routeId = 0
			})

			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(len(requests)).To(Equal(0))
			})
		})
	})

	Describe("CreateRequest", func() {
		var (
			request *model.Request
			err     error
		)

		newRequest := model.Request{
			RiderId: 3,
			RouteId: 3,
		}

		JustBeforeEach(func() {
			request, err = CreateRequest(&newRequest)
		})

		AfterEach(func() {
			_, err := pgPool.Exec(context.Background(), fmt.Sprintf(
				`DELETE FROM requests WHERE rider_id = %d AND route_id = %d;`,
				newRequest.RiderId,
				newRequest.RouteId,
			))
			Expect(err).NotTo(HaveOccurred())
		})

		When("request does not exist in database", func() {
			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(request.RiderId).To(Equal(int32(3)))
				Expect(request.RouteId).To(Equal(int32(3)))
			})
		})
	})

	Describe("UpdateRequestStatus", func() {
		var (
			request *model.Request
			err     error
		)

		JustBeforeEach(func() {
			err = UpdateRequestStatus(existedRequests[0].Id, constants.RequestStatusCompleted)
		})

		When("request exists in database", func() {
			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(request.Id).To(Equal(existedRequests[0].Id))
				Expect(request.RiderId).To(Equal(existedRequests[0].RiderId))
				Expect(request.RouteId).To(Equal(existedRequests[0].RouteId))
				Expect(request.Status).To(Equal(constants.RequestStatusAccepted))
			})
		})

		When("request does not exist in database", func() {
			It("fails", func() {
				Expect(err).To(MatchError(ErrUndefined))
				Expect(request).To(BeNil())
			})
		})
	})
})

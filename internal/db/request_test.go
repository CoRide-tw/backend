package db

import (
	"context"
	"fmt"
	"github.com/CoRide-tw/backend/internal/constants"
	. "github.com/CoRide-tw/backend/internal/errors/generated/dberr"
	"github.com/CoRide-tw/backend/internal/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

const testCreateRequestSQL = `
	INSERT INTO requests (rider_id, route_id, pickup_location, dropoff_location, pickup_start_time, pickup_end_time, tips, status)
	VALUES (
		$1,
		$2,
		ST_SetSRID(ST_MakePoint($3, $4), 4326),
		ST_SetSRID(ST_MakePoint($5, $6), 4326),
		$7,
		$8,
		$9,
		$10
	)
	RETURNING id;
`

var _ = Describe("DBRequest", func() {
	// test data
	existedRequests := []model.Request{
		{
			RiderId:         -1,
			RouteId:         -1,
			PickupLong:      121.0134308229882,
			PickupLat:       24.79100321524295,
			DropoffLong:     121.01444872393937,
			DropoffLat:      24.79071289283521,
			PickupStartTime: time.Now(),
			PickupEndTime:   time.Now(),
			Tips:            100,
			Status:          constants.RequestStatusPending,
		},
		{
			RiderId:         -1,
			RouteId:         -2,
			PickupLong:      121.0134308229882,
			PickupLat:       24.79100321524295,
			DropoffLong:     121.01444872393937,
			DropoffLat:      24.79071289283521,
			PickupStartTime: time.Now(),
			PickupEndTime:   time.Now(),
			Tips:            100,
			Status:          constants.RequestStatusPending,
		},
		{
			RiderId:         -2,
			RouteId:         -1,
			PickupLong:      121.0134308229882,
			PickupLat:       24.79100321524295,
			DropoffLong:     121.01444872393937,
			DropoffLat:      24.79071289283521,
			PickupStartTime: time.Now(),
			PickupEndTime:   time.Now(),
			Tips:            100,
			Status:          constants.RequestStatusPending,
		},
	}

	BeforeEach(func() {
		for i, request := range existedRequests {
			err := pgPool.QueryRow(context.Background(), testCreateRequestSQL,
				request.RiderId,
				request.RouteId,
				request.PickupLong,
				request.PickupLat,
				request.DropoffLong,
				request.DropoffLat,
				request.PickupStartTime,
				request.PickupEndTime,
				request.Tips,
				request.Status,
			).Scan(&existedRequests[i].Id)
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
				Expect(request.Id).To(Equal(existedRequests[0].Id))
				Expect(request.RiderId).To(Equal(existedRequests[0].RiderId))
				Expect(request.RouteId).To(Equal(existedRequests[0].RouteId))
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

	//Describe("ListRequestsByRouteId", func() {
	//	var (
	//		resp    []*ListRequestsByRouteIdResp
	//		routeId int32
	//		err     error
	//	)
	//
	//	JustBeforeEach(func() {
	//		resp, err = ListRequestsByRouteId(routeId)
	//	})
	//
	//	When("requests exist in database", func() {
	//		BeforeEach(func() {
	//			routeId = existedRequests[0].RouteId
	//		})
	//
	//		It("succeeds", func() {
	//			Expect(err).NotTo(HaveOccurred())
	//			Expect(len(resp)).To(Equal(2))
	//		})
	//	})
	//
	//	When("requests do not exist in database", func() {
	//		BeforeEach(func() {
	//			routeId = 0
	//		})
	//
	//		It("succeeds", func() {
	//			Expect(err).NotTo(HaveOccurred())
	//			Expect(len(requests)).To(Equal(0))
	//		})
	//	})
	//})

	Describe("CreateRequest", func() {
		var (
			request *model.Request
			err     error
		)

		newRequest := model.Request{
			RiderId:         -3,
			RouteId:         -3,
			PickupLong:      121.0134308229882,
			PickupLat:       24.79100321524295,
			DropoffLong:     121.01444872393937,
			DropoffLat:      24.79071289283521,
			PickupStartTime: time.Now(),
			PickupEndTime:   time.Now(),
			Tips:            100,
			Status:          constants.RequestStatusPending,
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
				Expect(request.RiderId).To(Equal(newRequest.RiderId))
				Expect(request.RouteId).To(Equal(newRequest.RouteId))
			})
		})
	})

	Describe("UpdateRequestStatus", func() {
		var (
			id  int32
			err error
		)

		JustBeforeEach(func() {
			err = UpdateRequestStatus(id, constants.RequestStatusCompleted)
		})

		When("request exists in database", func() {
			BeforeEach(func() {
				id = existedRequests[0].Id
			})

			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("DeleteRequest", func() {
		var (
			request *model.Request
			err     error
			getErr  error
			id      int32
		)

		JustBeforeEach(func() {
			err = DeleteRequest(id)
			request, getErr = GetRequest(id)
		})

		When("request exists in database", func() {
			BeforeEach(func() {
				id = existedRequests[0].Id
			})

			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(getErr).To(MatchError(ErrRequestNotFound))
				Expect(request).To(BeNil())
			})
		})

		When("request does not exist in database", func() {
			BeforeEach(func() {
				id = 0
			})

			It("fails", func() {
				Expect(getErr).To(MatchError(ErrRequestNotFound))
				Expect(request).To(BeNil())
			})
		})
	})
})

package db

import (
	"context"
	"github.com/CoRide-tw/backend/internal/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testCreateTripSQL = `
	INSERT INTO trips (rider_id, driver_id, request_id, route_id)
	VALUES (
		$1,
		$2, 
		$3, 
		$4
	)
	RETURNING id;
`

const testDeleteTripSQL = `
	DELETE FROM trips WHERE id = $1;
`

var _ = Describe("DBTrip", func() {
	var existedTrips []model.Trip

	BeforeEach(func() {
		existedTrips = []model.Trip{
			{RiderId: -1, DriverId: -1, RequestId: -1, RouteId: -1},
			{RiderId: -1, DriverId: -1, RequestId: -2, RouteId: -2},
			{RiderId: -2, DriverId: -2, RequestId: -3, RouteId: -3},
		}

		for i, existedTrip := range existedTrips {
			err := pgPool.QueryRow(context.Background(), testCreateTripSQL,
				existedTrip.RiderId,
				existedTrip.DriverId,
				existedTrip.RequestId,
				existedTrip.RouteId,
			).Scan(&existedTrips[i].Id)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	AfterEach(func() {
		for _, trip := range existedTrips {
			_, err := pgPool.Exec(context.Background(), testDeleteTripSQL, trip.Id)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	Describe("ListTripByRiderId", func() {
		var (
			trips   []*model.Trip
			riderId int32
			err     error
		)

		JustBeforeEach(func() {
			trips, err = ListTripByRiderId(riderId)
		})

		When("trips exist", func() {
			BeforeEach(func() {
				riderId = -1
			})

			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(len(trips)).To(Equal(2))
			})
		})
	})

	Describe("ListTripByDriverId", func() {
		var (
			trips    []*model.Trip
			driverId int32
			err      error
		)

		JustBeforeEach(func() {
			trips, err = ListTripByDriverId(driverId)
		})

		When("trips exist", func() {
			BeforeEach(func() {
				driverId = -1
			})

			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(len(trips)).To(Equal(2))
			})
		})
	})

	Describe("GetTrip", func() {
		var (
			trip *model.Trip
			id   int32
			err  error
		)

		JustBeforeEach(func() {
			trip, err = GetTrip(id)
		})

		When("trip exists", func() {
			BeforeEach(func() {
				id = existedTrips[0].Id
			})

			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(trip.Id).To(Equal(id))
			})
		})

		When("trip does not exist", func() {
			BeforeEach(func() {
				id = -1
			})

			It("fails", func() {
				Expect(err).NotTo(BeNil())
				Expect(trip).To(BeNil())
			})
		})
	})

	Describe("CreateTrip", func() {
		var (
			trip *model.Trip
			err  error
		)

		newTrip := model.Trip{
			RiderId:   -5,
			DriverId:  -5,
			RequestId: -5,
			RouteId:   -5,
		}

		JustBeforeEach(func() {
			trip, err = CreateTrip(&newTrip)
		})

		When("trip created", func() {
			It("succeed", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(trip.Id).NotTo(Equal(0))
				Expect(trip.RiderId).To(Equal(newTrip.RiderId))
				Expect(trip.DriverId).To(Equal(newTrip.DriverId))
				Expect(trip.RequestId).To(Equal(newTrip.RequestId))
				Expect(trip.RouteId).To(Equal(newTrip.RouteId))
			})
		})
	})
})

package db

import (
	"context"
	. "github.com/CoRide-tw/backend/internal/errors/generated/dberr"
	"github.com/CoRide-tw/backend/internal/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

const testCreateRouteSQL = `
	INSERT INTO routes (driver_id, start_location, end_location, start_time, end_time, capacity)
	VALUES (
		$1, 
		ST_SetSRID(ST_MakePoint($2, $3), 4326), 
		ST_SetSRID(ST_MakePoint($4, $5), 4326), 
		$6, 
		$7, 
		$8
	)

	RETURNING id;
`

const testDeleteRouteSQL = `
	DELETE FROM routes WHERE id = $1;
`

var _ = Describe("DBRequest", func() {
	var (
		existedRoutes                                                  []model.Route
		validStartTime, validEndTime, invalidStartTime, invalidEndTime time.Time
		err                                                            error
	)

	BeforeEach(func() {
		// test data
		validStartTime, err = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
		Expect(err).NotTo(HaveOccurred())
		validEndTime, err = time.Parse(time.RFC3339, "2006-01-02T18:04:05Z")
		Expect(err).NotTo(HaveOccurred())
		invalidStartTime, err = time.Parse(time.RFC3339, "2006-01-03T15:04:05Z")
		Expect(err).NotTo(HaveOccurred())
		invalidEndTime, err = time.Parse(time.RFC3339, "2006-01-03T15:04:05Z")
		Expect(err).NotTo(HaveOccurred())

		existedRoutes = []model.Route{
			// 星巴克關埔店 到 松江屋
			{
				DriverId:  -1,
				StartLong: 121.0134308229882,
				StartLat:  24.79100321524295,
				EndLong:   121.01444872393937,
				EndLat:    24.79071289283521,
				StartTime: validStartTime,
				EndTime:   validEndTime,
				Capacity:  1000,
			},
			// 豐邑商辦大樓 到 路易莎關埔店 (invalid time)
			{
				DriverId:  -1,
				StartLong: 121.01272590458588,
				StartLat:  24.79130028800565,
				EndLong:   121.01537274231576,
				EndLat:    24.790525949323005,
				StartTime: invalidStartTime,
				EndTime:   invalidEndTime,
				Capacity:  1000,
			},
			// 壽司郎 到 契茶小野田
			{
				DriverId:  -1,
				StartLong: 121.01192442739676,
				StartLat:  24.791619557659118,
				EndLong:   121.01631060640568,
				EndLat:    24.78999202940558,
				StartTime: validStartTime,
				EndTime:   validEndTime,
				Capacity:  1000,
			},
		}

		for i, route := range existedRoutes {
			err := pgPool.QueryRow(context.Background(), testCreateRouteSQL,
				route.DriverId,
				route.StartLong,
				route.StartLat,
				route.EndLong,
				route.EndLat,
				route.StartTime,
				route.EndTime,
				route.Capacity,
			).Scan(
				&existedRoutes[i].Id,
			)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	AfterEach(func() {
		for _, route := range existedRoutes {
			_, err := pgPool.Exec(context.Background(), testDeleteRouteSQL, route.Id)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	Describe("GetRoute", func() {
		var (
			route *model.Route
			id    int32
			err   error
		)

		JustBeforeEach(func() {
			route, err = GetRoute(id)
		})

		When("route exists in database", func() {
			BeforeEach(func() {
				id = existedRoutes[0].Id
			})

			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(route.Id).To(Equal(id))
				Expect(route.StartLong).To(Equal(existedRoutes[0].StartLong))
				Expect(route.StartLat).To(Equal(existedRoutes[0].StartLat))
				Expect(route.EndLong).To(Equal(existedRoutes[0].EndLong))
				Expect(route.EndLat).To(Equal(existedRoutes[0].EndLat))
				Expect(route.StartTime.UTC()).To(Equal(existedRoutes[0].StartTime))
				Expect(route.EndTime.UTC()).To(Equal(existedRoutes[0].EndTime))
				Expect(route.CreatedAt).NotTo(BeNil())
				Expect(route.UpdatedAt).NotTo(BeNil())
			})
		})

		When("route does not exist in database", func() {
			BeforeEach(func() {
				id = 0
			})

			It("fails", func() {
				Expect(err).To(MatchError(ErrRouteNotFound))
				Expect(route).To(BeNil())
			})
		})
	})

	Describe("ListNearestRoutes", func() {
		var (
			routes []*model.Route
			err    error
		)

		JustBeforeEach(func() {
			pickupTime, err := time.Parse(time.RFC3339, "2006-01-02T16:04:05Z")
			Expect(err).NotTo(HaveOccurred())
			dropOffTime, err := time.Parse(time.RFC3339, "2006-01-02T17:04:05Z")
			Expect(err).NotTo(HaveOccurred())
			// 在星巴克關埔店跟松江烏之間的兩個點
			pickupLong := 121.01373815586145
			pickupLat := 24.790756765799653
			dropOffLong := 121.01408790650603
			dropOffLat := 24.790713673871583

			routes, err = ListNearestRoutes(pickupLong, pickupLat, dropOffLong, dropOffLat, pickupTime, dropOffTime)
		})

		When("there are routes in database", func() {
			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(routes).NotTo(BeNil())
				Expect(len(routes)).To(BeNumerically(">=", 2))

				Expect(routes[0].Id).To(Equal(existedRoutes[0].Id))
				Expect(routes[1].Id).To(Equal(existedRoutes[2].Id))
			})
		})
	})

	Describe("CreateRoute", func() {
		var (
			route *model.Route
			err   error
		)

		newRoute := model.Route{
			DriverId:  -2,
			StartLong: 121.01272590458588,
			StartLat:  24.79130028800565,
			EndLong:   121.01537274231576,
			EndLat:    24.790525949323005,
			StartTime: validStartTime,
			EndTime:   validEndTime,
			Capacity:  1000,
		}

		JustBeforeEach(func() {
			route, err = CreateRoute(&newRoute)
		})

		When("route is valid", func() {
			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(route.Id).NotTo(BeZero())
				Expect(route.DriverId).To(Equal(newRoute.DriverId))
				Expect(route.StartLong).To(Equal(newRoute.StartLong))
				Expect(route.StartLat).To(Equal(newRoute.StartLat))
				Expect(route.EndLong).To(Equal(newRoute.EndLong))
				Expect(route.EndLat).To(Equal(newRoute.EndLat))
				Expect(route.StartTime).To(Equal(newRoute.StartTime))
				Expect(route.EndTime).To(Equal(newRoute.EndTime))
				Expect(route.Capacity).To(Equal(newRoute.Capacity))
				Expect(route.CreatedAt).NotTo(BeNil())
				Expect(route.UpdatedAt).NotTo(BeNil())
			})
		})
	})

	Describe("DeleteRoute", func() {
		var (
			route  *model.Route
			err    error
			getErr error
			id     int32
		)

		JustBeforeEach(func() {
			err = DeleteRoute(id)
			route, getErr = GetRoute(id)
		})

		When("route exists in database", func() {
			BeforeEach(func() {
				id = existedRoutes[0].Id
			})

			It("succeeds", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(getErr).To(MatchError(ErrRouteNotFound))
				Expect(route).To(BeNil())
			})
		})

		When("route does not exist in database", func() {
			BeforeEach(func() {
				id = 0
			})

			It("fails", func() {
				Expect(getErr).To(MatchError(ErrRouteNotFound))
				Expect(route).To(BeNil())
			})
		})
	})
})

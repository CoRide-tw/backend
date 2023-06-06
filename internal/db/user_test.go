package db

import (
	"context"
	. "github.com/CoRide-tw/backend/internal/errors/generated/dberr"
	"github.com/CoRide-tw/backend/internal/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testCreateUserSQL = `
	INSERT INTO users (name, email, google_id, picture_url)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
`

const testDeleteUserSQL = `
	DELETE FROM users WHERE id = $1;
`

var _ = Describe("DBUser", func() {
	existedUser := model.User{
		Name:       "test",
		Email:      "test",
		GoogleId:   "test",
		PictureUrl: "test",
	}

	BeforeEach(func() {
		err := pgPool.QueryRow(context.Background(), testCreateUserSQL,
			existedUser.Name,
			existedUser.Email,
			existedUser.GoogleId,
			existedUser.PictureUrl,
		).Scan(&existedUser.Id)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		_, err := pgPool.Exec(context.Background(), testDeleteUserSQL, existedUser.Id)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("GetUser", func() {
		var (
			user *model.User
			id   int32
			err  error
		)

		JustBeforeEach(func() {
			user, err = GetUser(id)
		})

		When("user does not exist", func() {
			BeforeEach(func() {
				id = 0
			})

			It("should return error", func() {
				Expect(err).To(MatchError(ErrUserNotFound))
				Expect(user).To(BeNil())
			})
		})

		When("user exists", func() {
			BeforeEach(func() {
				id = existedUser.Id
			})

			It("should return user", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(user.Id).To(Equal(existedUser.Id))
				Expect(user.Name).To(Equal(existedUser.Name))
				Expect(user.Email).To(Equal(existedUser.Email))
				Expect(user.GoogleId).To(Equal(existedUser.GoogleId))
				Expect(user.PictureUrl).To(Equal(existedUser.PictureUrl))
				Expect(user.UpdatedAt).NotTo(BeNil())
				Expect(user.CreatedAt).NotTo(BeNil())
			})
		})
	})

	Describe("CreateUser", func() {
		var (
			user *model.User
			err  error
		)

		newUser := model.User{
			Name:       "test1",
			Email:      "test1",
			GoogleId:   "test1",
			PictureUrl: "test1",
		}

		JustBeforeEach(func() {
			user, err = UpsertUser(&newUser)
		})

		When("user created", func() {
			It("should return user", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(user.Id).NotTo(BeZero())
				Expect(user.Name).To(Equal(newUser.Name))
				Expect(user.Email).To(Equal(newUser.Email))
				Expect(user.GoogleId).To(Equal(newUser.GoogleId))
				Expect(user.PictureUrl).To(Equal(newUser.PictureUrl))
				Expect(user.UpdatedAt).NotTo(BeNil())
				Expect(user.CreatedAt).NotTo(BeNil())
			})
		})
	})

	Describe("UpdateUser", func() {
		var (
			id   int32
			err  error
			user *model.User
		)

		userWithNewEmail := model.User{
			Email: "newEmail",
		}

		JustBeforeEach(func() {
			user, err = UpdateUser(id, &userWithNewEmail)
		})

		When("user does not exist", func() {
			BeforeEach(func() {
				id = 0
			})

			It("should return error", func() {
				Expect(err).To(MatchError(ErrUserNotFound))
				Expect(user).To(BeNil())
			})
		})

		When("user exists", func() {
			BeforeEach(func() {
				id = existedUser.Id
			})

			It("should return user", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(user.Id).To(Equal(existedUser.Id))
				Expect(user.Name).To(Equal(existedUser.Name))
				Expect(user.Email).To(Equal(userWithNewEmail.Email))
				Expect(user.GoogleId).To(Equal(existedUser.GoogleId))
				Expect(user.PictureUrl).To(Equal(existedUser.PictureUrl))
				Expect(user.UpdatedAt).NotTo(BeNil())
				Expect(user.CreatedAt).NotTo(BeNil())
			})
		})
	})

	Describe("DeleteUser", func() {
		var (
			id     int32
			err    error
			getErr error
			user   *model.User
		)

		JustBeforeEach(func() {
			err = DeleteUser(id)
			user, getErr = GetUser(id)
		})

		When("user exists", func() {
			BeforeEach(func() {
				id = existedUser.Id
			})

			It("should delete user", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(getErr).To(MatchError(ErrUserNotFound))
				Expect(user).To(BeNil())
			})
		})
	})
})

package test

import (
	"context"
	"testing"

	"github.com/LaQuannT/astronaut-api/internal/model"
	"github.com/LaQuannT/astronaut-api/internal/service"
	"github.com/stretchr/testify/assert"
)

const plainPwd = "Qwerty_123"

func TestRegisterUser(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Fatalf("Error clearing tables: %v", err)
	}

	ctx := context.TODO()

	t.Run("returns an error for invalid user data", func(t *testing.T) {
		usr := new(model.User)

		usr, err := service.RegisterUser(ctx, userRepo, usr)
		if err == nil {
			t.Error("expected an error registering user")
		}
		assert.Nil(t, usr)
	})

	t.Run("registers a user", func(t *testing.T) {
		usr := &model.User{
			FirstName: "john",
			LastName:  "boy",
			Email:     "john@test.com",
			Password:  plainPwd,
		}

		u, err := service.RegisterUser(ctx, userRepo, usr)
		if err != nil {
			t.Errorf("unexpected error registering user: %v", err)
		}

		assert.Equal(t, usr.Email, u.Email)
		assert.NotEqual(t, plainPwd, u.Password)
		assert.NotEmpty(t, u.APIKey)
	})

	t.Run("return an error if email is already in use", func(t *testing.T) {
		usr := &model.User{
			FirstName: "John",
			LastName:  "doe",
			Email:     "john@test.com",
			Password:  plainPwd,
		}

		usr, err := service.RegisterUser(ctx, userRepo, usr)
		if err == nil {
			t.Error("expected an error registering user")
		}
		assert.Nil(t, usr)
	})
}

func TestSearchUserID(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Fatalf("error clearing tables: %v", err)
	}

	ctx := context.TODO()

	u := &model.User{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane@test.com",
		Password:  plainPwd,
	}

	u, err := service.RegisterUser(ctx, userRepo, u)
	if err != nil {
		t.Fatalf("unexpected error registering a user: %v", err)
	}

	t.Run("returns nil and error for unknown user ID", func(t *testing.T) {
		uid := 87

		usr, err := service.SearchUserID(ctx, userRepo, uid)
		if err == nil {
			t.Error("expected an error searching user ID")
		}

		assert.Nil(t, usr)
	})

	t.Run("returns a user", func(t *testing.T) {
		uid := u.ID

		usr, err := service.SearchUserID(ctx, userRepo, uid)
		if err != nil {
			t.Errorf("unexpected error searching user ID: %v", err)
		}

		assert.Equal(t, u.FirstName, usr.FirstName)
		assert.Equal(t, u.LastName, usr.LastName)
		assert.Equal(t, u.Email, usr.Email)
	})
}

func TestSearchUserEmail(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Fatalf("error clearing tables: %v", err)
	}

	ctx := context.TODO()

	u := &model.User{
		FirstName: "bob",
		LastName:  "Doe",
		Email:     "bob@test.com",
		Password:  plainPwd,
	}

	u, err := service.RegisterUser(ctx, userRepo, u)
	if err != nil {
		t.Fatalf("unexpected error registering a user: %v", err)
	}

	t.Run("returns nil and a error for unknown user email", func(t *testing.T) {
		uEmail := "test@test.com"

		usr, err := service.SearchUserEmail(ctx, userRepo, uEmail)
		if err == nil {
			t.Error("expected an error searching unknown user email")
		}
		assert.Nil(t, usr)
	})

	t.Run("returns a user", func(t *testing.T) {
		uEmail := u.Email

		usr, err := service.SearchUserEmail(ctx, userRepo, uEmail)
		if err != nil {
			t.Errorf("unexpected error searching user email: %v", err)
		}

		assert.Equal(t, u.Email, usr.Email)
	})
}

func TestGetUsers(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Fatalf("error clearing tables: %v", err)
	}

	ctx := context.TODO()

	t.Run("return nil if no users are found", func(t *testing.T) {
		usr, err := service.GetUsers(ctx, userRepo)
		if err != nil {
			t.Errorf("unexpected error getting users: %v", err)
		}
		assert.Nil(t, usr)
	})

	t.Run("returns a list of users", func(t *testing.T) {
		jane := &model.User{
			FirstName: "jane",
			LastName:  "doe",
			Email:     "jane@test.com",
			Password:  plainPwd,
		}

		john := &model.User{
			FirstName: "john",
			LastName:  "doe",
			Email:     "john@test.com",
			Password:  plainPwd,
		}

		for _, u := range []*model.User{jane, john} {
			_, err := service.RegisterUser(ctx, userRepo, u)
			if err != nil {
				t.Fatalf("unexpected error registering users: %v", err)
			}
		}

		usrs, err := service.GetUsers(ctx, userRepo)
		if err != nil {
			t.Errorf("unexpected error getting users: %v", err)
		}
		assert.Len(t, usrs, 2)
	})
}

func TestUpdateUser(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Fatalf("error clearing tables: %v", err)
	}

	ctx := context.TODO()

	u := &model.User{
		FirstName: "test",
		LastName:  "test",
		Email:     "test@test.com",
		Password:  plainPwd,
	}

	u, err := service.RegisterUser(ctx, userRepo, u)
	if err != nil {
		t.Fatalf("unexpected error registering user: %v", err)
	}

	t.Run("returns a error for invalid new user data", func(t *testing.T) {
		usr := &model.User{
			ID:        u.ID,
			FirstName: "",
			LastName:  "",
			Email:     "",
			Password:  "",
		}

		err := service.UpdateUser(ctx, userRepo, usr)
		if err == nil {
			t.Error("expected an error updating user with invalid input")
		}
	})

	t.Run("returns an updated astronaut", func(t *testing.T) {
		usr := &model.User{
			ID:        u.ID,
			FirstName: "james",
			LastName:  "tucker",
			Email:     u.Email,
			Password:  u.Password,
		}

		err := service.UpdateUser(ctx, userRepo, usr)
		if err != nil {
			t.Errorf("unexpected error updating user: %v", err)
		}

		usr, err = service.SearchUserID(ctx, userRepo, u.ID)
		if err != nil {
			t.Errorf("unexpected error geting updated user: %v", err)
		}

		assert.NotEqual(t, u.FirstName, usr.FirstName)
		assert.NotEqual(t, u.LastName, usr.LastName)
	})
}

func TestDeleteUser(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Fatalf("error clearing tables: %v", err)
	}

	ctx := context.TODO()

	t.Run("returns an error when trying to delete unknown user", func(t *testing.T) {
		uid := 998

		if err := service.DeleteUser(ctx, userRepo, uid); err == nil {
			t.Error("expected an error deleting unknown user")
		}
	})

	t.Run("deletes a user", func(t *testing.T) {
		u := &model.User{
			FirstName: "test",
			LastName:  "test",
			Email:     "test@test.com",
			Password:  plainPwd,
		}

		u, err := service.RegisterUser(ctx, userRepo, u)
		if err != nil {
			t.Fatalf("unexpected error registering a user: %v", err)
		}

		err = service.DeleteUser(ctx, userRepo, u.ID)
		if err != nil {
			t.Errorf("unexpected error deleting a user: %v", err)
		}
	})
}

func TestResetPassword(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Fatalf("error clearing tables: %v", err)
	}

	ctx := context.TODO()

	u := &model.User{
		FirstName: "test",
		LastName:  "test",
		Email:     "test@test.com",
		Password:  plainPwd,
	}

	u, err := service.RegisterUser(ctx, userRepo, u)
	if err != nil {
		t.Fatalf("unexpected error registering user: %v", err)
	}

	t.Run("returns error when reseting unknown user password", func(t *testing.T) {
		uid := 89
		pwd := "X3dfdgvdfvgf_"

		err := service.ResetPassword(ctx, userRepo, pwd, uid)
		if err == nil {
			t.Error("expected error reseting unknown user password")
		}
	})

	t.Run("return error for invalid new password", func(t *testing.T) {
		uid := u.ID

		err := service.ResetPassword(ctx, userRepo, "", uid)
		if err == nil {
			t.Error("expected error for invalid new password")
		}
	})

	t.Run("resets a user password", func(t *testing.T) {
		uid := u.ID
		pwd := "X3dfdgvdfvgf_"

		err := service.ResetPassword(ctx, userRepo, pwd, uid)
		if err != nil {
			t.Errorf("unexpected error reseting user password: %v", err)
		}

		usr, err := service.SearchUserID(ctx, userRepo, uid)
		if err != nil {
			t.Errorf("unexpected error fetching user: %v", err)
		}

		assert.NotEqual(t, u.Password, usr.Password)
	})
}

func TestGenerateNewUserAPIKey(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Fatalf("error clearing tables: %v", err)
	}

	ctx := context.TODO()

	u := &model.User{
		FirstName: "john",
		LastName:  "tucker",
		Email:     "jtucker@email.com",
		Password:  plainPwd,
	}

	u, err := service.RegisterUser(ctx, userRepo, u)
	if err != nil {
		t.Errorf("unexpected error registering a user: %v", err)
	}

	t.Run("returns error generating new APIKey for unknown user", func(t *testing.T) {
		uid := 78

		key, err := service.GenerateNewAPIKey(ctx, userRepo, uid)
		if err == nil {
			t.Error("expected error generating new APIKey for unknown user")
		}
		assert.Empty(t, key)
	})

	t.Run("generates new APIKey", func(t *testing.T) {
		uid := u.ID
		u, err := service.SearchUserID(ctx, userRepo, uid)
		if err != nil {
			t.Fatalf("unexpected error geting user: %v", err)
		}

		key, err := service.GenerateNewAPIKey(ctx, userRepo, uid)
		if err != nil {
			t.Errorf("unexpected error generating new APIKey: %v", err)
		}

		assert.NotEqual(t, u.APIKey, key)
	})
}

func TestCreateAdmin(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Fatalf("error clearing tables: %v", err)
	}

	ctx := context.TODO()

	u := &model.User{
		FirstName: "super",
		LastName:  "admin",
		Email:     "admin@email.com",
		Password:  plainPwd,
	}

	u, err := service.RegisterUser(ctx, userRepo, u)
	if err != nil {
		t.Errorf("unexpected error registering a user: %v", err)
	}

	t.Run("returns error when creating an admin with invalid user ID", func(t *testing.T) {
		uid := 1000

		err := service.CreateAdmin(ctx, userRepo, uid)
		if err == nil {
			t.Error("expected error creating admin for unknown user")
		}
	})

	t.Run("creates a admin", func(t *testing.T) {
		uid := u.ID

		err := service.CreateAdmin(ctx, userRepo, uid)
		if err != nil {
			t.Errorf("unexpected error creating a admin: %v", err)
		}
	})
}

func TestRemoveAdmin(t *testing.T) {
	ctx := context.TODO()

	t.Run("returns error trying to remove admin for unknown user ID", func(t *testing.T) {
		uid := 1234

		err := service.RemoveAdmin(ctx, userRepo, uid)
		if err == nil {
			t.Error("expected error removing unknown admin user")
		}
	})

	t.Run("removes a admin user", func(t *testing.T) {
		uid := 1

		err := service.RemoveAdmin(ctx, userRepo, uid)
		if err != nil {
			t.Errorf("expected error removing admin user: %v", err)
		}
	})
}

func TestCheckAdminPermission(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Fatalf("error clearing tables: %v", err)
	}

	ctx := context.TODO()

	jane := &model.User{
		FirstName: "jane",
		LastName:  "doe",
		Email:     "jane@email.com",
		Password:  plainPwd,
	}

	john := &model.User{
		FirstName: "john",
		LastName:  "doe",
		Email:     "john@email.com",
		Password:  plainPwd,
	}

	usrs := []*model.User{jane, john}

	for i, u := range usrs {
		usrs[i] = u
		u, err := service.RegisterUser(ctx, userRepo, usrs[i])
		if err != nil {
			t.Fatalf("unexpected error registering users: %v", err)
		}

		if u.FirstName == "jane" {
			err = service.CreateAdmin(ctx, userRepo, u.ID)
			if err != nil {
				t.Fatalf("error adding admin permission: %v", err)
			}
		}
	}

	t.Run("returns false if user is not an admin", func(t *testing.T) {
		isAdmin, err := service.CheckAdminPermission(ctx, userRepo, john.ID)
		if err != nil {
			t.Errorf("unexpected error checking admin permission: %v", err)
		}

		assert.False(t, isAdmin)
	})

	t.Run("returns true if user is an admin", func(t *testing.T) {
		isAdmin, err := service.CheckAdminPermission(ctx, userRepo, jane.ID)
		if err != nil {
			t.Errorf("unexpected error checking admin permission: %v", err)
		}
		assert.True(t, isAdmin)
	})
}

func TestSearchUserAPIKey(t *testing.T) {
	if err := clearTables(dbConn); err != nil {
		t.Fatalf("error clearing tables: %v", err)
	}
	ctx := context.TODO()

	u := &model.User{
		FirstName: "john",
		LastName:  "doe",
		Email:     "john@email.com",
		Password:  plainPwd,
	}
	u, err := service.RegisterUser(ctx, userRepo, u)
	if err != nil {
		t.Errorf("unexpected error registering a user: %v", err)
	}

	t.Run("returns error trying to find a user with invalid API key", func(t *testing.T) {
		key := "67e505ea-51ee-402a-9a79-cd0c7afede4b"

		usr, err := service.SearchAPIKey(ctx, userRepo, key)
		if err == nil {
			t.Error("expected error trying to find a user with invalid API key")
		}
		assert.Nil(t, usr)
	})

	t.Run("returns a user", func(t *testing.T) {
		usr, err := service.SearchAPIKey(ctx, userRepo, u.APIKey)
		if err != nil {
			t.Errorf("unexpected error searching user API key: %v", err)
		}
		assert.Equal(t, u.Email, usr.Email)
	})
}

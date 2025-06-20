package entity

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

type userEntityTestSuite struct {
	testName string
	id       uuid.UUID
	name     string
	email    string
}

func TestUserEntity(t *testing.T) {
	var tests = []userEntityTestSuite{
		{
			testName: "Valid User",
			id:       uuid.New(),
			name:     "John Doe",
			email:    "john@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			user := User{
				ID:    tt.id,
				Name:  tt.name,
				Email: tt.email,
			}

			if user.Name != tt.name {
				t.Errorf("expected Name %q, got %q", tt.name, user.Name)
			}
			if user.Email != tt.email {
				t.Errorf("expected Email %q, got %q", tt.email, user.Email)
			}
			if user.ID != tt.id {
				t.Errorf("expected ID %v, got %v", tt.id, user.ID)
			}
			if reflect.TypeOf(user.ID) != reflect.TypeOf(uuid.UUID{}) {
				t.Errorf(
					"expected ID type %v, got %v",
					reflect.TypeOf(uuid.UUID{}),
					reflect.TypeOf(user.ID),
				)
			}
		})
	}
}

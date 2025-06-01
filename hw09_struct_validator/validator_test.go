package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid user",
			in: User{
				ID:     "12345678-1234-5678-1234-567812345678",
				Age:    25,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "09876543210"},
			},
			expectedErr: nil,
		},
		{
			name: "invalid age and email",
			in: User{
				ID:     "123",
				Age:    17,
				Email:  "not-an-email",
				Role:   "guest",
				Phones: []string{"12345", "0987654321"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: errValidationLen},
				{Field: "Age", Err: errValidationMin},
				{Field: "Email", Err: errValidationRegexp},
				{Field: "Role", Err: errValidationIn},
			},
		},
		{
			name: "valid app version",
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			name:        "non-struct input",
			in:          "string",
			expectedErr: errNotStruct,
		},
		{
			name: "valid response",
			in: Response{
				Code: 200,
			},
			expectedErr: nil,
		},
		{
			name: "invalid response code",
			in: Response{
				Code: 403,
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: errValidationIn},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tt.name), func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in)

			if tt.expectedErr == nil {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}

			switch expected := tt.expectedErr.(type) {
			case ValidationErrors:
				var actual ValidationErrors
				if !errors.As(err, &actual) {
					t.Errorf("expected ValidationErrors, got %T", err)
					return
				}

				for i, e := range expected {
					if actual[i].Field != e.Field || !errors.Is(actual[i].Err, e.Err) {
						t.Errorf("error mismatch at index %d: expected {Field: %q, Err: %v}, got {Field: %q, Err: %v}",
							i, e.Field, e.Err, actual[i].Field, actual[i].Err)
					}
				}
			case error:
				if !errors.Is(err, expected) {
					t.Errorf("expected error %v, got %v", expected, err)
				}
			default:
				t.Errorf("unknown expected error type: %T", expected)
			}
		})
	}
}

package key

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

// Test_Key_CAPath tests the CAPath function.
func Test_Key_CAPath(t *testing.T) {
	tests := []struct {
		groupName    string
		expectedPath string
	}{
		{
			groupName:    "xa5ly",
			expectedPath: "/certs/xa5ly/ca.pem",
		},
		{
			groupName:    "fah0a",
			expectedPath: "/certs/fah0a/ca.pem",
		},
	}

	for index, test := range tests {
		path := CAPath(test.groupName)

		if !reflect.DeepEqual(test.expectedPath, path) {
			t.Fatalf(
				"%d: expected path does not match returned path\nexpected: %s\nreturned: %s\n",
				index,
				spew.Sdump(test.expectedPath),
				spew.Sdump(path),
			)
		}
	}
}

// Test_Key_CrtPath tests the CrtPath function.
func Test_Key_CrtPath(t *testing.T) {
	tests := []struct {
		groupName    string
		expectedPath string
	}{
		{
			groupName:    "xa5ly",
			expectedPath: "/certs/xa5ly/crt.pem",
		},
		{
			groupName:    "fah0a",
			expectedPath: "/certs/fah0a/crt.pem",
		},
	}

	for index, test := range tests {
		path := CrtPath(test.groupName)

		if !reflect.DeepEqual(test.expectedPath, path) {
			t.Fatalf(
				"%d: expected path does not match returned path\nexpected: %s\nreturned: %s\n",
				index,
				spew.Sdump(test.expectedPath),
				spew.Sdump(path),
			)
		}
	}
}

// Test_Key_KeyPath tests the KeyPath function.
func Test_Key_KeyPath(t *testing.T) {
	tests := []struct {
		groupName    string
		expectedPath string
	}{
		{
			groupName:    "xa5ly",
			expectedPath: "/certs/xa5ly/key.pem",
		},
		{
			groupName:    "fah0a",
			expectedPath: "/certs/fah0a/key.pem",
		},
	}

	for index, test := range tests {
		path := KeyPath(test.groupName)

		if !reflect.DeepEqual(test.expectedPath, path) {
			t.Fatalf(
				"%d: expected path does not match returned path\nexpected: %s\nreturned: %s\n",
				index,
				spew.Sdump(test.expectedPath),
				spew.Sdump(path),
			)
		}
	}
}

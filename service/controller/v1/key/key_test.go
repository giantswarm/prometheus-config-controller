package key

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

// Test_Key_CAPath tests the CAPath function.
func Test_Key_CAPath(t *testing.T) {
	tests := []struct {
		certificateDirectory string
		clusterID            string

		expectedPath string
	}{
		{
			certificateDirectory: "/certs",
			clusterID:            "xa5ly",

			expectedPath: "/certs/xa5ly-ca.pem",
		},
		{
			certificateDirectory: "/certs",
			clusterID:            "fah0a",

			expectedPath: "/certs/fah0a-ca.pem",
		},
		{
			certificateDirectory: "/certificates",
			clusterID:            "fah0a",

			expectedPath: "/certificates/fah0a-ca.pem",
		},
	}

	for index, test := range tests {
		path := CAPath(test.certificateDirectory, test.clusterID)

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
		certificateDirectory string
		clusterID            string

		expectedPath string
	}{
		{
			certificateDirectory: "/certs",
			clusterID:            "xa5ly",

			expectedPath: "/certs/xa5ly-crt.pem",
		},
		{
			certificateDirectory: "/certs",
			clusterID:            "fah0a",

			expectedPath: "/certs/fah0a-crt.pem",
		},
		{
			certificateDirectory: "/certificates",
			clusterID:            "fah0a",

			expectedPath: "/certificates/fah0a-crt.pem",
		},
	}

	for index, test := range tests {
		path := CrtPath(test.certificateDirectory, test.clusterID)

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
		certificateDirectory string
		clusterID            string

		expectedPath string
	}{
		{
			certificateDirectory: "/certs",
			clusterID:            "xa5ly",

			expectedPath: "/certs/xa5ly-key.pem",
		},
		{
			certificateDirectory: "/certs",
			clusterID:            "fah0a",

			expectedPath: "/certs/fah0a-key.pem",
		},
		{
			certificateDirectory: "/certificates",
			clusterID:            "fah0a",

			expectedPath: "/certificates/fah0a-key.pem",
		},
	}

	for index, test := range tests {
		path := KeyPath(test.certificateDirectory, test.clusterID)

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

func Test_LabelSelectorConfigMap(t *testing.T) {
	testCases := []struct {
		name                   string
		expectedSelectorString string
	}{
		{
			name:                   "case 0",
			expectedSelectorString: "app=prometheus",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			selectorString := LabelSelectorConfigMap().String()
			if !reflect.DeepEqual(tc.expectedSelectorString, selectorString) {
				t.Fatalf("tc.expectedSelectorString = %v, want %v", tc.expectedSelectorString, selectorString)
			}

		})
	}
}

func Test_LabelSelectorService(t *testing.T) {
	testCases := []struct {
		name                   string
		expectedSelectorString string
	}{
		{
			name:                   "case 0",
			expectedSelectorString: "app=master,giantswarm.io/cluster",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			selectorString := LabelSelectorService().String()
			if !reflect.DeepEqual(tc.expectedSelectorString, selectorString) {
				t.Fatalf("tc.expectedSelectorString = %v, want %v", tc.expectedSelectorString, selectorString)
			}

		})
	}
}

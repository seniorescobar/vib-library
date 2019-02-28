package library

import (
	"net/url"
	"testing"
)

func TestCheckRequiredParams(t *testing.T) {
	type testCase struct {
		requestParams  url.Values
		requiredParams []string
		expectedErr    error
	}

	tcs := []testCase{
		// no missing params
		testCase{
			url.Values{
				"a": []string{"a1"},
				"b": []string{"b1"},
				"c": []string{"c1"},
			},
			[]string{
				"a",
				"b",
				"c",
			},
			nil,
		},
		// missing param "a"
		testCase{
			url.Values{},
			[]string{
				"a",
			},
			ErrMissingParams,
		},
	}

	for _, tc := range tcs {
		actualErr := checkRequiredParams(tc.requestParams, tc.requiredParams)
		if actualErr != tc.expectedErr {
			t.Errorf("Test failed, expected: '%s', got:  '%s'", tc.expectedErr, actualErr)
		}
	}
}

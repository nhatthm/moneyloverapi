package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnixString_Time(t *testing.T) {
	t.Parallel()

	ts := time.Unix(1619783139, 0)
	ux := UnixString(ts)

	assert.Equal(t, ts, ux.Time())
}

func TestUnixString_Sec(t *testing.T) {
	t.Parallel()

	ts := time.Unix(1619783139, 0)
	ux := UnixString(ts)

	assert.Equal(t, int64(1619783139), ux.Sec())
}

func TestUnixString_String(t *testing.T) {
	t.Parallel()

	ts := time.Unix(1619783139, 0)
	ux := UnixString(ts)

	assert.Equal(t, "1619783139", ux.String())
}

func TestUnixString_MarshalJSON(t *testing.T) {
	t.Parallel()

	ts := time.Unix(1619783139, 0)
	ux := UnixString(ts)

	res, err := json.Marshal(ux)
	require.NoError(t, err)

	assert.Equal(t, []byte(`"1619783139"`), res)
}

func TestUnixString_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		data           []byte
		expectedResult UnixString
		expectedError  string
	}{
		{
			scenario:      "data is nil",
			expectedError: `unexpected end of JSON input`,
		},
		{
			scenario:      "data is not a string",
			data:          []byte(`42`),
			expectedError: `json: cannot unmarshal number into Go value of type string`,
		},
		{
			scenario:      "data is not a unix timestamp string",
			data:          []byte(`"foobar"`),
			expectedError: `strconv.ParseInt: parsing "foobar": invalid syntax`,
		},
		{
			scenario:       "data is a unix timestamp string",
			data:           []byte(`"1612325106"`),
			expectedResult: UnixStringDate(2021, 2, 3, 4, 5, 6, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			var result UnixString

			err := json.Unmarshal(tc.data, &result)

			assert.Equal(t, tc.expectedResult, result)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

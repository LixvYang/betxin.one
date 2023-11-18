package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrToInt64(t *testing.T) {
	tests := []struct {
		input    string
		expected struct {
			result int64
			err    error
		}
	}{
		{"123", struct {
			result int64
			err    error
		}{123, nil}},
		{"-456", struct {
			result int64
			err    error
		}{-456, nil}},
		{"0", struct {
			result int64
			err    error
		}{0, nil}},
	}

	for _, test := range tests {
		result, err := StrToInt64(test.input)
		assert.Equal(t, test.expected.result, result)
		assert.Equal(t, test.expected.err, err)
	}
}

func TestIntToStr(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{123, "123"},
		{-456, "-456"},
		{3.14, "3.14"},
		{"hello", "hello"},
	}

	for _, test := range tests {
		result := IntToStr(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{struct{ Name string }{"John"}, `{"Name":"John"}`},
		{struct{ Age int }{25}, `{"Age":25}`},
		{struct{ Score float64 }{3.14}, `{"Score":3.14}`},
	}

	for _, test := range tests {
		result, err := Marshal(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, result)
	}
}

func TestUnmarshal(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int64  `json:"age"`
	}
	tests := []struct {
		input    string
		expected TestStruct
	}{
		{`{"Name":"John", "age": 12}`, TestStruct{
			Name: "John",
			Age:  12,
		}},
	}

	for _, test := range tests {
		var result TestStruct
		err := Unmarshal(test.input, &result)
		assert.NoError(t, err)
		assert.Equal(t, test.expected.Age, result.Age)
		assert.Equal(t, test.expected.Name, result.Name)
	}
}

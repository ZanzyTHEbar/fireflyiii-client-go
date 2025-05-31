package firefly

import (
	"time"

	"github.com/oapi-codegen/runtime/types"
)

// Helper function to convert string to pointer
func stringPtr(s string) *string {
	return &s
}

// boolPtr returns a pointer to a bool
func boolPtr(b bool) *bool {
	return &b
}

// float64Ptr returns a pointer to a float64
func float64Ptr(f float64) *float64 {
	return &f
}

// int32Ptr returns a pointer to an int32
func int32Ptr(i int) *int32 {
	val := int32(i)
	return &val
}

// timePtr returns a pointer to a time.Time
func timePtr(t time.Time) *time.Time {
	return &t
}

// int32Value returns 0 if the pointer is nil, otherwise returns the value
func int32Value(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

// timeValue returns the zero time if the pointer is nil, otherwise returns the value
func timeValue(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// boolValue returns false if the pointer is nil, otherwise returns the value
func boolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// stringValue returns an empty string if the pointer is nil, otherwise returns the value
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper functions for type conversions
func dateToAPIDate(t *time.Time) *types.Date {
	if t == nil {
		return nil
	}
	date := types.Date{Time: *t}
	return &date
}

func apiDateToTime(d *types.Date) *time.Time {
	if d == nil {
		return nil
	}
	return &d.Time
}

func float32Value(f *float32) float32 {
	if f == nil {
		return 0
	}
	return *f
}
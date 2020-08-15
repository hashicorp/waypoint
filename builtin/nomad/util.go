package nomad

import "time"

// boolToPtr returns the pointer to a boolean
func boolToPtr(b bool) *bool {
	return &b
}

// int8ToPtr returns the pointer to an int8
func int8ToPtr(i int8) *int8 {
	return &i
}

// intToPtr returns the pointer to an int
func intToPtr(i int) *int {
	return &i
}

// uint64ToPtr returns the pointer to an uint64
func uint64ToPtr(u uint64) *uint64 {
	return &u
}

// int64ToPtr returns the pointer to a int64
func int64ToPtr(i int64) *int64 {
	return &i
}

// stringToPtr returns the pointer to a string
func stringToPtr(str string) *string {
	return &str
}

// timeToPtr returns the pointer to a time stamp
func timeToPtr(t time.Duration) *time.Duration {
	return &t
}

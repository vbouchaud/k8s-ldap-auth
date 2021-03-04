package version

import (
	// "reflect"
	"testing"
	// "time"
)

func None() func() {
	return func() {}
}

func Mock() func() {
	return func() {
		APPNAME = "pipeline"
		VERSION = "0.1.0"
		GOVERSION = "1.13"
		BUILDTIME = "2019-10-11T07:34:49Z+00:00"
		COMMITHASH = "13754adcbf"
	}
}

func MockInvalidBuildTime() func() {
	return func() {
		APPNAME = "pipeline"
		VERSION = "0.1.0"
		GOVERSION = "1.13"
		BUILDTIME = "not-10-a 09:da:te+02:00"
		COMMITHASH = "13754adcbf"
	}
}

func Reset() func() {
	return func() {
		APPNAME = "UNKNOWN"
		VERSION = "UNKNOWN"
		GOVERSION = "UNKNOWN"
		BUILDTIME = "UNKNOWN"
		COMMITHASH = "UNKNOWN"
	}
}

func TestVersion(t *testing.T) {
	tests := []struct {
		name   string
		want   string
		before func()
		after  func()
	}{
		{
			name:   "Unset variables",
			want:   "UNKNOWN UNKNOWN (commit UNKNOWN built with go UNKNOWN the UNKNOWN)",
			before: None(),
			after:  None(),
		},
		{
			name:   "With variables",
			want:   "pipeline 0.1.0 (commit 13754adcbf built with go 1.13 the 2019-10-11T07:34:49Z+00:00)",
			before: Mock(),
			after:  Reset(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			if got := Version(); got != tt.want {
				t.Errorf("Version() = %v, want %v", got, tt.want)
			}
			tt.after()
		})
	}
}

// func TestCompiled(t *testing.T) {
//		now := time.Now()
//		expected, err := time.Parse("2006-02-15T15:04:05Z+00:00", "2019-10-11T07:34:49Z+00:00")
//		if err != nil {
//			t.Errorf("Failed to parse expected compiled date, %s", err)
//		}

//		tests := []struct {
//			name   string
//			want   time.Time
//			before func()
//			after  func()
//		}{
//			{
//				name:   "Unset variables",
//				want:   now,
//				before: None(),
//				after:  None(),
//			},
//			{
//				name:   "With invalid buildtime",
//				want:   now,
//				before: MockInvalidBuildTime(),
//				after:  Reset(),
//			},
//			{
//				name:   "With valid buildtime",
//				want:   expected,
//				before: Mock(),
//				after:  Reset(),
//			},
//		}

//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				tt.before()
//				if got := Compiled(); !reflect.DeepEqual(got, tt.want) {
//					t.Errorf("Compiled() = %v, want %v", got, tt.want)
//				}
//				tt.after()
//			})
//		}
// }

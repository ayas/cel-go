package types

import (
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types/ref"
	"testing"
)

func TestTimestamp_Add(t *testing.T) {
	ts := Timestamp{&timestamp.Timestamp{Seconds: 7506}}
	val := ts.Add(Duration{&duration.Duration{Seconds: 3600, Nanos: 1000}})
	if val.ConvertToType(TypeType) != TimestampType {
		t.Error("Could not add duration and timestamp")
	}
	expected := Timestamp{&timestamp.Timestamp{Seconds: 11106, Nanos: 1000}}
	if !expected.Compare(val).Equal(IntZero).(Bool) {
		t.Errorf("Got '%v', expected '%v'", val, expected)
	}
	if !IsError(ts.Add(expected)) {
		t.Error("Cannot add two timestamps together")
	}
}

func TestTimestamp_Subtract(t *testing.T) {
	ts := Timestamp{&timestamp.Timestamp{Seconds: 7506}}
	val := ts.Subtract(Duration{&duration.Duration{Seconds: 3600, Nanos: 1000}})
	if val.ConvertToType(TypeType) != TimestampType {
		t.Error("Could not add duration and timestamp")
	}
	expected := Timestamp{&timestamp.Timestamp{Seconds: 3905, Nanos: 999999000}}
	if !expected.Compare(val).Equal(IntZero).(Bool) {
		t.Errorf("Got '%v', expected '%v'", val, expected)
	}
}

func TestTimestamp_ReceiveGetHours(t *testing.T) {
	// 1970-01-01T02:05:05Z
	ts := Timestamp{&timestamp.Timestamp{Seconds: 7506}}
	hr := ts.Receive(overloads.TimeGetHours, overloads.TimestampToHours, []ref.Value{})
	if !hr.Equal(Int(2)).(Bool) {
		t.Error("Expected 2 hours, got", hr)
	}
	// 1969-12-31T19:05:05Z
	hrTz := ts.Receive(overloads.TimeGetHours, overloads.TimestampToHoursWithTz,
		[]ref.Value{String("America/Phoenix")})
	if !hrTz.Equal(Int(19)).(Bool) {
		t.Error("Expected 19 hours, got", hrTz)
	}
}

func TestTimestamp_ReceiveGetMinutes(t *testing.T) {
	// 1970-01-01T02:05:05Z
	ts := Timestamp{&timestamp.Timestamp{Seconds: 7506}}
	min := ts.Receive(overloads.TimeGetMinutes, overloads.TimestampToMinutes, []ref.Value{})
	if !min.Equal(Int(5)).(Bool) {
		t.Error("Expected 5 minutes, got", min)
	}
	// 1969-12-31T19:05:05Z
	minTz := ts.Receive(overloads.TimeGetMinutes, overloads.TimestampToMinutesWithTz,
		[]ref.Value{String("America/Phoenix")})
	if !minTz.Equal(Int(5)).(Bool) {
		t.Error("Expected 5 minutes, got", minTz)
	}
}

func TestTimestamp_ReceiveGetSeconds(t *testing.T) {
	// 1970-01-01T02:05:05Z
	ts := Timestamp{&timestamp.Timestamp{Seconds: 7506}}
	sec := ts.Receive(overloads.TimeGetSeconds, overloads.TimestampToSeconds, []ref.Value{})
	if !sec.Equal(Int(6)).(Bool) {
		t.Error("Expected 6 seconds, got", sec)
	}
	// 1969-12-31T19:05:05Z
	secTz := ts.Receive(overloads.TimeGetSeconds, overloads.TimestampToSecondsWithTz,
		[]ref.Value{String("America/Phoenix")})
	if !secTz.Equal(Int(6)).(Bool) {
		t.Error("Expected 6 seconds, got", secTz)
	}
}

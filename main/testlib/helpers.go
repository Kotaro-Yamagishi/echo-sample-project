package testlib

import (
	"errors"
	"testing"
	"time"
)

// AssertError エラーの存在をアサート
func AssertError(t *testing.T, err error, expected bool) {
	t.Helper()
	if expected && err == nil {
		t.Error("expected error, got nil")
	}
	if !expected && err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// AssertErrorMessage エラーメッセージをアサート
func AssertErrorMessage(t *testing.T, err error, expectedMsg string) {
	t.Helper()
	if err == nil {
		t.Error("expected error, got nil")
		return
	}
	if err.Error() != expectedMsg {
		t.Errorf("expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// CreateTestTime テスト用の固定時間を作成
func CreateTestTime() time.Time {
	return time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
}

// CreateTestError テスト用のエラーを作成
func CreateTestError(message string) error {
	return errors.New(message)
}

// AssertSliceLength スライスの長さをアサート
func AssertSliceLength(t *testing.T, slice interface{}, expected int) {
	t.Helper()
	switch s := slice.(type) {
	case []interface{}:
		if len(s) != expected {
			t.Errorf("expected slice length %d, got %d", expected, len(s))
		}
	default:
		t.Errorf("unsupported slice type: %T", slice)
	}
}

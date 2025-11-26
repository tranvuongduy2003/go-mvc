package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func StringToInt(s string, defaultValue int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return defaultValue
}

func StringToInt64(s string, defaultValue int64) int64 {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	return defaultValue
}

func StringToFloat64(s string, defaultValue float64) float64 {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return defaultValue
}

func StringToBool(s string, defaultValue bool) bool {
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	return defaultValue
}

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func Remove[T comparable](slice []T, item T) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if v != item {
			result = append(result, v)
		}
	}
	return result
}

func Unique[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	result := make([]T, 0)
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	return result
}

func Map[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = fn(item)
	}
	return result
}

func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func Reduce[T, U any](slice []T, initialValue U, reducer func(U, T) U) U {
	result := initialValue
	for _, item := range slice {
		result = reducer(result, item)
	}
	return result
}

func ToCamelCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	if len(words) == 0 {
		return ""
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		word := words[i]
		if len(word) > 0 {
			result += strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return result
}

func ToSnakeCase(s string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToKebabCase(s string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	kebab := re.ReplaceAllString(s, "${1}-${2}")
	return strings.ToLower(kebab)
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func Truncate(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

func RandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

func RandomInt(min, max int) (int, error) {
	if min > max {
		return 0, fmt.Errorf("min cannot be greater than max")
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()) + min, nil
}

func DeepCopy(src, dst interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("failed to marshal source: %w", err)
	}

	if err := json.Unmarshal(bytes, dst); err != nil {
		return fmt.Errorf("failed to unmarshal to destination: %w", err)
	}

	return nil
}

func IsZero(v interface{}) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

func Coalesce[T comparable](values ...T) T {
	var zero T
	for _, v := range values {
		if v != zero {
			return v
		}
	}
	return zero
}

func Ternary[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

func Ptr[T any](v T) *T {
	return &v
}

func SafeDeref[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

type TimeHelper struct{}

func NewTimeHelper() *TimeHelper {
	return &TimeHelper{}
}

func (h *TimeHelper) Now() time.Time {
	return time.Now().UTC()
}

func (h *TimeHelper) StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func (h *TimeHelper) EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

func (h *TimeHelper) IsToday(t time.Time) bool {
	now := h.Now()
	return h.StartOfDay(t).Equal(h.StartOfDay(now))
}

func (h *TimeHelper) DaysBetween(start, end time.Time) int {
	start = h.StartOfDay(start)
	end = h.StartOfDay(end)
	return int(end.Sub(start).Hours() / 24)
}

func (h *TimeHelper) FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.1fd", d.Hours()/24)
}

func (h *TimeHelper) ParseDuration(s string) (time.Duration, error) {
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid duration format")
	}

	unit := s[len(s)-1:]
	valueStr := s[:len(s)-1]
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid duration value: %w", err)
	}

	switch unit {
	case "s":
		return time.Duration(value * float64(time.Second)), nil
	case "m":
		return time.Duration(value * float64(time.Minute)), nil
	case "h":
		return time.Duration(value * float64(time.Hour)), nil
	case "d":
		return time.Duration(value * float64(24*time.Hour)), nil
	default:
		return time.ParseDuration(s)
	}
}

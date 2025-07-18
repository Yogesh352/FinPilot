package util

func StrPtr(s string) *string {
	return &s
}

func FloatPtr(f float64) *float64 {
	return &f
}

func DerefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
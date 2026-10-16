package tasks

import "strings"

func TaskStatusToUpper(status *string) *string {
	if status == nil {
		return nil
	}
	upper := strings.ToUpper(*status)
	return &upper
}

func TaskStatusToLower(status *string) *string {
	if status == nil {
		return nil
	}
	lower := strings.ToLower(*status)
	return &lower
}

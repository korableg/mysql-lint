package analyzer

import (
	"context"
	"testing"
)

func TestName(t *testing.T) {
	anal := NewAnalyzer()
	err := anal.Analyze(context.Background(), "/home/korableg/git/space307/olymptrade.com/payment/")
	err = anal.Analyze(context.Background(), "/home/korableg/git/space307/olymptrade.com/userspace/")
	if err == nil {
		return
	}
}

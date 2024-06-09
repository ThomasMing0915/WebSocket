package main

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type OTP struct {
	Key    string
	Create time.Time
}

type RetentionMap map[string]OTP

func (m RetentionMap) NewOTP() OTP {
	o := OTP{
		Key:    uuid.NewString(),
		Create: time.Now(),
	}
	m[o.Key] = o
	return o
}

func (rm RetentionMap) VerifyOTP(otp string) bool {
	if _, ok := rm[otp]; !ok {
		return false
	}
	delete(rm, otp)
	return true
}

func (m RetentionMap) Retention(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(400 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			for _, otp := range m {
				if otp.Create.Add(period).Before(time.Now()) {
					delete(m, otp.Key)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func NewRetentionMap(ctx context.Context, retentionPeriod time.Duration) RetentionMap {
	rm := make(RetentionMap)
	go rm.Retention(ctx, retentionPeriod)
	return rm
}

package client

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/seaweedfs/seaweedfs/weed/glog"
	"github.com/seaweedfs/seaweedfs/weed/pb/master_pb"
)

const (
	RenewInterval     = 4 * time.Second
	SafeRenewInterval = 3 * time.Second
	InitLockInterval  = 1 * time.Second
)

type ExclusiveLocker struct {
	token    int64
	lockTsNs int64
	isLocked bool
	cm       *commandEnv
	lockName string
	message  string
}

func NewExclusiveLocker(cm *commandEnv, lockName string) *ExclusiveLocker {
	return &ExclusiveLocker{
		cm:       cm,
		lockName: lockName,
	}
}

func (l *ExclusiveLocker) IsLocked() bool {
	return l.isLocked
}

func (l *ExclusiveLocker) GetToken() (token int64, lockTsNs int64) {
	for time.Unix(0, atomic.LoadInt64(&l.lockTsNs)).Add(SafeRenewInterval).Before(time.Now()) {
		// wait until now is within the safe lock period, no immediate renewal to change the token
		time.Sleep(100 * time.Millisecond)
	}
	return atomic.LoadInt64(&l.token), atomic.LoadInt64(&l.lockTsNs)
}

func (l *ExclusiveLocker) RequestLock(clientName string) {
	if l.isLocked {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// retry to get the lease
	for {
		if err := l.cm.MasterClient.WithClientCustomGetMaster(l.cm.getMasterAddress, false, func(client master_pb.SeaweedClient) error {
			resp, err := client.LeaseAdminToken(ctx, &master_pb.LeaseAdminTokenRequest{
				PreviousToken:    atomic.LoadInt64(&l.token),
				PreviousLockTime: atomic.LoadInt64(&l.lockTsNs),
				LockName:         l.lockName,
				ClientName:       clientName,
			})
			if err == nil {
				atomic.StoreInt64(&l.token, resp.Token)
				atomic.StoreInt64(&l.lockTsNs, resp.LockTsNs)
			}
			return err
		}); err != nil {
			println("lock:", err.Error())
			time.Sleep(InitLockInterval)
		} else {
			break
		}
	}

	l.isLocked = true

	// start a goroutine to renew the lease
	go func() {
		ctx2, cancel2 := context.WithCancel(context.Background())
		defer cancel2()

		for l.isLocked {
			if err := l.cm.MasterClient.WithClientCustomGetMaster(l.cm.getMasterAddress, false, func(client master_pb.SeaweedClient) error {
				resp, err := client.LeaseAdminToken(ctx2, &master_pb.LeaseAdminTokenRequest{
					PreviousToken:    atomic.LoadInt64(&l.token),
					PreviousLockTime: atomic.LoadInt64(&l.lockTsNs),
					LockName:         l.lockName,
					ClientName:       clientName,
					Message:          l.message,
				})
				if err == nil {
					atomic.StoreInt64(&l.token, resp.Token)
					atomic.StoreInt64(&l.lockTsNs, resp.LockTsNs)
					// println("ts", l.lockTsNs, "token", l.token)
				}
				return err
			}); err != nil {
				glog.Errorf("failed to renew lock: %v", err)
				l.isLocked = false
				return
			} else {
				time.Sleep(RenewInterval)
			}

		}
	}()

}

func (l *ExclusiveLocker) ReleaseLock() {
	l.isLocked = false

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l.cm.MasterClient.WithClientCustomGetMaster(l.cm.getMasterAddress, false, func(client master_pb.SeaweedClient) error {
		client.ReleaseAdminToken(ctx, &master_pb.ReleaseAdminTokenRequest{
			PreviousToken:    atomic.LoadInt64(&l.token),
			PreviousLockTime: atomic.LoadInt64(&l.lockTsNs),
			LockName:         l.lockName,
		})
		return nil
	})
	atomic.StoreInt64(&l.token, 0)
	atomic.StoreInt64(&l.lockTsNs, 0)
}

func (l *ExclusiveLocker) SetMessage(message string) {
	l.message = message
}

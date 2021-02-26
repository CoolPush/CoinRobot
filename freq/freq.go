package freq

import (
	"fmt"
	"time"
)

type Freq struct {
}

func NewFreq() *Freq {
	return &Freq{}
}

// Check 检查 每次check 都会对 key ++ true说明检测通过了不需要限频 false需要限频
func (f *Freq) Check(key string, intervalSeconds, limit uint32) bool {
	redisKey := fmt.Sprintf("%s_%d", key, intervalSeconds)
	val, err := freqer.IncrBy(reidsCtx, redisKey, 1).Result()
	if err != nil {
		log.Errorf("err: %v", err)
		return true
	}
	if val < 0 {
		return true
	}
	realValue := uint32(val)
	if realValue == 1 {
		err = freqer.Expire(reidsCtx, redisKey, time.Duration(intervalSeconds)*time.Second).Err()
		if err != nil {
			log.Errorf("err:%v", err)
		}
	}
	if realValue > limit {
		// 收集一下 请求量大的
		if realValue >= 100 {
			log.Warnf("freq key: %s, current cnt %d, interval %d, limit %d",
				key, realValue, intervalSeconds, limit)
		}

		return false
	}

	return true
}

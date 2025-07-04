package exercise

import (
	"sync"
	"testing"
)

// 測試 Lottery 函數
func TestLotteryFunction(t *testing.T) {
	// 重置 currentRewards 到初始狀態
	resetCurrentRewards()

	// 測試抽獎多次
	times := 10
	gifts := Lottery(times)

	// 驗證返回的禮物數量
	if len(gifts) > times {
		t.Errorf("Expected at most %d gifts, got %d", times, len(gifts))
	}

	// 驗證所有禮物都是有效的
	for _, gift := range gifts {
		if gift == nil {
			t.Error("Expected non-nil gift, got nil")
		}
		if gift.ID < 1 || gift.ID > 3 {
			t.Errorf("Invalid gift ID: %d", gift.ID)
		}
	}
}

// 測試 EnterLottery 函數
func TestEnterLottery(t *testing.T) {
	// 重置 currentRewards 到初始狀態
	resetCurrentRewards()

	// 測試並發抽獎
	persons := 3
	times := 5

	// 記錄初始狀態
	initialRewards := make(map[int]int)
	for id, count := range currentRewards {
		initialRewards[id] = count
	}

	// 執行抽獎
	EnterLottery(persons, times)

	// 驗證禮物被消耗
	for id, initialCount := range initialRewards {
		if currentRewards[id] > initialCount {
			t.Errorf("Expected rewards for ID %d to decrease or stay same, got %d (was %d)",
				id, currentRewards[id], initialCount)
		}
	}
}

// 測試邊界情況
func TestEdgeCases(t *testing.T) {
	// 重置 currentRewards 到初始狀態
	resetCurrentRewards()

	// 測試 0 次抽獎
	gifts := Lottery(0)
	if len(gifts) != 0 {
		t.Errorf("Expected 0 gifts for 0 times, got %d", len(gifts))
	}

	// 測試負數次抽獎
	gifts = Lottery(-1)
	if len(gifts) != 0 {
		t.Errorf("Expected 0 gifts for negative times, got %d", len(gifts))
	}

	// 測試 0 人抽獎 - 應該安全處理
	EnterLottery(0, 5)
	// 應該不會出錯

	// 測試負數人抽獎 - 應該安全處理
	EnterLottery(-1, 5)
	// 應該不會出錯
}

// 測試 WaitGroup 行為
func TestWaitGroupBehavior(t *testing.T) {
	// 重置 currentRewards 到初始狀態
	resetCurrentRewards()

	// 測試正常情況
	persons := 3
	times := 2

	// 記錄初始狀態
	initialRewards := make(map[int]int)
	for id, count := range currentRewards {
		initialRewards[id] = count
	}

	// 執行抽獎 - 應該正常完成
	EnterLottery(persons, times)

	// 驗證禮物被消耗
	totalConsumed := 0
	for id, initialCount := range initialRewards {
		consumed := initialCount - currentRewards[id]
		totalConsumed += consumed
	}

	// 應該有禮物被消耗
	if totalConsumed == 0 {
		t.Error("Expected some gifts to be consumed, but none were")
	}
}

// 測試並發安全性（改進版）
func TestConcurrencySafetyImproved(t *testing.T) {
	// 重置 currentRewards 到初始狀態
	resetCurrentRewards()

	// 測試多個 goroutine 同時抽獎
	const numGoroutines = 10
	const timesPerGoroutine = 5

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// 記錄初始狀態
	initialRewards := make(map[int]int)
	for id, count := range currentRewards {
		initialRewards[id] = count
	}

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			Lottery(timesPerGoroutine)
		}()
	}

	wg.Wait()

	// 驗證沒有負數的獎勵數量
	for id, count := range currentRewards {
		if count < 0 {
			t.Errorf("Reward count for ID %d should not be negative, got %d", id, count)
		}
	}

	// 驗證禮物被正確消耗
	totalConsumed := 0
	for id, initialCount := range initialRewards {
		consumed := initialCount - currentRewards[id]
		totalConsumed += consumed
	}

	// 應該有禮物被消耗
	if totalConsumed == 0 {
		t.Error("Expected some gifts to be consumed in concurrent test, but none were")
	}
}

// 測試並發安全性
func TestConcurrencySafety(t *testing.T) {
	// 重置 currentRewards 到初始狀態
	resetCurrentRewards()

	// 測試多個 goroutine 同時抽獎
	const numGoroutines = 20
	const timesPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			Lottery(timesPerGoroutine)
		}()
	}

	wg.Wait()

	// 驗證沒有負數的獎勵數量
	for id, count := range currentRewards {
		if count < 0 {
			t.Errorf("Reward count for ID %d should not be negative, got %d", id, count)
		}
	}
}

// 測試機率分佈（粗略測試）
func TestProbabilityDistribution(t *testing.T) {
	// 重置 currentRewards 到初始狀態
	resetCurrentRewards()

	const numTests = 10000
	giftCounts := make(map[int]int)

	for i := 0; i < numTests; i++ {
		gift := lottery()
		if gift != nil {
			giftCounts[gift.ID]++
		}
	}

	// 驗證每種禮物都有被抽中（粗略測試）
	for _, gift := range gifts {
		if giftCounts[gift.ID] == 0 {
			t.Errorf("Expected gift ID %d to be drawn at least once", gift.ID)
		}
	}
}

// 輔助函數：重置 currentRewards 到初始狀態
func resetCurrentRewards() {
	currentRewardsLock.Lock()
	defer currentRewardsLock.Unlock()

	currentRewards = make(map[int]int)
	for _, gift := range gifts {
		currentRewards[gift.ID] = int(gift.Probability * base)
	}
}

// 基準測試
func BenchmarkLottery(b *testing.B) {
	resetCurrentRewards()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lottery()
	}
}

func BenchmarkLotteryFunction(b *testing.B) {
	resetCurrentRewards()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Lottery(10)
	}
}

func BenchmarkGetGiftWithLock(b *testing.B) {
	resetCurrentRewards()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getGiftWithLock(1)
	}
}

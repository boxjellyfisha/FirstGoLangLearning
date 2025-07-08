package exercise

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

var gifts []Gift
var currentRewards map[int]int // key: Gift.ID, value: current count
var currentRewardsLock sync.Mutex

const base = 100

func init() {
	gifts = []Gift{
		{ID: 1, Name: "SSR: happy corgi", Probability: 1.0 / base},
		{ID: 2, Name: "R: doll", Probability: 2.0 / base},
		{ID: 3, Name: "N: sticker", Probability: 4.0 / base},
	}

	currentRewards = make(map[int]int)
	for _, gift := range gifts {
		currentRewards[gift.ID] = int(gift.Probability * base)
	}
}

type Gift struct {
	ID          int
	Name        string
	Probability float64
}

func getGift(id int) *Gift {
	for _, gift := range gifts {
		if gift.ID == id && currentRewards[id] > 0 {
			currentRewards[id]--
			return &gift
		}
	}
	return nil
}

func getGiftWithLock(id int) *Gift {
	currentRewardsLock.Lock()
	defer currentRewardsLock.Unlock()
	return getGift(id)
}

func lottery() *Gift {
	result := float64(rand.Intn(base)) / base
	for _, gift := range gifts {
		if result < gift.Probability {
			return getGiftWithLock(gift.ID)
		}
	}
	return nil
}

func Lottery(times int) []*Gift {
	var result []*Gift
	for range times {
		gift := lottery()
		if gift != nil {
			result = append(result, gift)
		}
		time.Sleep(10 * time.Millisecond)
	}
	return result
}

func EnterLottery(persion int, times int) {
	// 檢查參數有效性
	if persion <= 0 {
		fmt.Println("Invalid number of persons:", persion)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(persion)

	for i := range persion {
		go func(personID int) {
			defer wg.Done()

			fmt.Println("person " + strconv.Itoa(personID) + " start lottery")
			gifts := Lottery(times)
			if len(gifts) > 0 {
				giftNames := make([]string, len(gifts))
				for j, gift := range gifts {
					giftNames[j] = gift.Name
				}
				fmt.Println("person " + strconv.Itoa(personID) + " got " + strings.Join(giftNames, ", "))
			}
		}(i)
	}
	wg.Wait()

	for id, remaining := range currentRewards {
		fmt.Println(gifts[id-1].Name + " remaining: " + strconv.Itoa(remaining))
	}
}

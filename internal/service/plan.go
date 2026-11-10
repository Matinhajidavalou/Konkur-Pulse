package service

import (
	"errors"
	"math"
)

// AllocateStudyTime زمان مطالعه را بر اساس تراز دروس تقسیم می‌کند
func AllocateStudyTime(subjects map[string]int, totalTime, targetTraz, minTime int) (map[string]int, error) {
	n := len(subjects)
	if n == 0 {
		return nil, errors.New("لیست دروس خالی است")
	}

	// ۱. بررسی شرط اولیه
	if totalTime < n*minTime {
		return nil, errors.New("زمان کل برای اختصاص حداقل زمان به همه دروس کافی نیست")
	}

	remainingTime := float64(totalTime - (n * minTime))

	// ۲. محاسبه وزن‌ها (اختلاف تا هدف)
	weights := make(map[string]float64)
	totalWeight := 0.0

	for name, traz := range subjects {
		weight := float64(targetTraz - traz)
		if weight < 0 {
			weight = 0
		}
		weights[name] = weight
		totalWeight += weight
	}

	// ۳. تخصیص زمان پایه + زمان متغیر
	allocatedTimes := make(map[string]int)
	allocatedSum := 0

	for name := range subjects {
		var share float64
		if totalWeight > 0 {
			share = (weights[name] / totalWeight) * remainingTime
		} else {
			// اگر همه به تراز هدف رسیده‌اند، زمان مساوی تقسیم شود
			share = remainingTime / float64(n)
		}

		// رند کردن به نزدیکترین مضرب ۵
		finalTime := float64(minTime) + share
		roundedTime := int(math.Round(finalTime/5.0) * 5.0)

		allocatedTimes[name] = roundedTime
		allocatedSum += roundedTime
	}

	// ۴. اصلاح خطای گرد کردن (اعمال روی ضعیف‌ترین درس)
	diff := totalTime - allocatedSum
	if diff != 0 {
		var weakestSubject string
		maxWeight := -1.0

		// پیدا کردن درسی که بیشترین اختلاف تراز را دارد
		for name, w := range weights {
			if w > maxWeight {
				maxWeight = w
				weakestSubject = name
			}
		}

		allocatedTimes[weakestSubject] += diff
	}

	return allocatedTimes, nil
}

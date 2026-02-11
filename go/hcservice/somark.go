package hcservice

import (
	"fmt"
	"sync"
)

// SomarkAllocator manages unique somark values per-real with reference counting.
// The same real address across multiple VIPs shares a single somark value.
type SomarkAllocator struct {
	mu           sync.Mutex
	baseSomark   uint32
	maxReals     uint32
	realToSomark map[string]uint32 // realAddr -> somark value
	refCount     map[string]int    // realAddr -> VIP reference count
	nextOffset   uint32            // next allocation offset
	freeList     []uint32          // recycled offsets
}

// NewSomarkAllocator creates a new SomarkAllocator.
//
// Parameters:
//   - baseSomark: The starting somark value. All allocated somarks will be >= this value.
//   - maxReals: The maximum number of unique reals that can be tracked simultaneously.
//
// Returns a new SomarkAllocator instance.
func NewSomarkAllocator(baseSomark, maxReals uint32) *SomarkAllocator {
	return &SomarkAllocator{
		baseSomark:   baseSomark,
		maxReals:     maxReals,
		realToSomark: make(map[string]uint32),
		refCount:     make(map[string]int),
	}
}

// Acquire increments the reference count for a real address and allocates a somark if needed.
//
// Parameters:
//   - realAddr: The real server address (e.g., "192.168.1.1").
//
// Returns:
//   - somark: The allocated somark value.
//   - isNew: True if this is a new allocation (caller should register with katran).
//   - err: Non-nil if capacity is exhausted.
func (a *SomarkAllocator) Acquire(realAddr string) (somark uint32, isNew bool, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if sm, exists := a.realToSomark[realAddr]; exists {
		a.refCount[realAddr]++
		return sm, false, nil
	}

	// Allocate new somark
	var offset uint32
	if len(a.freeList) > 0 {
		offset = a.freeList[len(a.freeList)-1]
		a.freeList = a.freeList[:len(a.freeList)-1]
	} else {
		if a.nextOffset >= a.maxReals {
			return 0, false, fmt.Errorf("somark capacity exhausted (max %d reals)", a.maxReals)
		}
		offset = a.nextOffset
		a.nextOffset++
	}

	sm := a.baseSomark + offset
	a.realToSomark[realAddr] = sm
	a.refCount[realAddr] = 1
	return sm, true, nil
}

// Release decrements the reference count for a real address and recycles the somark if no longer used.
//
// Parameters:
//   - realAddr: The real server address to release.
//
// Returns:
//   - somark: The somark value that was associated with this real.
//   - isLast: True if the reference count reached zero (caller should deregister from katran).
//   - err: Non-nil if the real address is not tracked.
func (a *SomarkAllocator) Release(realAddr string) (somark uint32, isLast bool, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	sm, exists := a.realToSomark[realAddr]
	if !exists {
		return 0, false, fmt.Errorf("real %q not tracked in somark allocator", realAddr)
	}

	a.refCount[realAddr]--
	if a.refCount[realAddr] <= 0 {
		offset := sm - a.baseSomark
		a.freeList = append(a.freeList, offset)
		delete(a.realToSomark, realAddr)
		delete(a.refCount, realAddr)
		return sm, true, nil
	}

	return sm, false, nil
}

// GetSomark returns the somark value for a real address without modifying reference counts.
//
// Parameters:
//   - realAddr: The real server address to look up.
//
// Returns the somark value and true if found, or 0 and false if not tracked.
func (a *SomarkAllocator) GetSomark(realAddr string) (uint32, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	sm, ok := a.realToSomark[realAddr]
	return sm, ok
}

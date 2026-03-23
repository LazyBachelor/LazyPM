package modal

// FocusArea represents a distinct focusable area in the UI
type FocusArea int

const (
	FocusNone FocusArea = iota
	FocusList
	FocusDetail
	FocusColumn0 // Backlog
	FocusColumn1
	FocusColumn2
	FocusColumn3
	FocusColumn4
)

// FocusManager handles focus state across different UI areas.
// It provides a clean separation of focus concerns from modal state.
type FocusManager struct {
	currentArea FocusArea
	areas       map[FocusArea]bool
}

// NewFocusManager creates a new focus manager with all areas disabled by default
func NewFocusManager() *FocusManager {
	return &FocusManager{
		currentArea: FocusNone,
		areas:       make(map[FocusArea]bool),
	}
}

// SetCurrent sets the currently focused area
func (f *FocusManager) SetCurrent(area FocusArea) {
	f.currentArea = area
}

// Current returns the currently focused area
func (f *FocusManager) Current() FocusArea {
	return f.currentArea
}

// IsFocused returns true if the given area is currently focused
func (f *FocusManager) IsFocused(area FocusArea) bool {
	return f.currentArea == area
}

// IsListFocused returns true if any list area is focused
func (f *FocusManager) IsListFocused() bool {
	return f.currentArea == FocusList ||
		f.currentArea == FocusColumn0 ||
		f.currentArea == FocusColumn1 ||
		f.currentArea == FocusColumn2 ||
		f.currentArea == FocusColumn3 ||
		f.currentArea == FocusColumn4
}

// NextColumn moves focus to the next column (for kanban)
func (f *FocusManager) NextColumn() {
	areas := []FocusArea{FocusColumn0, FocusColumn1, FocusColumn2, FocusColumn3, FocusColumn4}
	f.cycleFocus(areas)
}

// PreviousColumn moves focus to the previous column (for kanban)
func (f *FocusManager) PreviousColumn() {
	areas := []FocusArea{FocusColumn4, FocusColumn3, FocusColumn2, FocusColumn1, FocusColumn0}
	f.cycleFocus(areas)
}

// IsDetailFocused returns true if the detail area is focused
func (f *FocusManager) IsDetailFocused() bool {
	return f.currentArea == FocusDetail
}

// EnableArea marks an area as available for focus
func (f *FocusManager) EnableArea(area FocusArea) {
	f.areas[area] = true
}

// DisableArea marks an area as unavailable for focus
func (f *FocusManager) DisableArea(area FocusArea) {
	f.areas[area] = false
	if f.currentArea == area {
		f.currentArea = FocusNone
	}
}

// IsAreaEnabled returns true if the area is enabled
func (f *FocusManager) IsAreaEnabled(area FocusArea) bool {
	return f.areas[area]
}

// Next moves focus to the next enabled area
func (f *FocusManager) Next() {
	areas := []FocusArea{FocusList, FocusDetail}
	f.cycleFocus(areas)
}

// Previous moves focus to the previous enabled area
func (f *FocusManager) Previous() {
	areas := []FocusArea{FocusDetail, FocusList}
	f.cycleFocus(areas)
}

// cycleFocus finds the next enabled area in the given order
func (f *FocusManager) cycleFocus(areas []FocusArea) {
	// Find current position
	startIdx := -1
	for i, area := range areas {
		if area == f.currentArea {
			startIdx = i
			break
		}
	}

	// Search for next enabled area
	for i := 1; i <= len(areas); i++ {
		idx := (startIdx + i) % len(areas)
		if idx < 0 {
			idx += len(areas)
		}
		if f.areas[areas[idx]] {
			f.currentArea = areas[idx]
			return
		}
	}
}

// ToggleDetail toggles between list and detail focus
func (f *FocusManager) ToggleDetail() {
	if f.currentArea == FocusDetail {
		f.currentArea = FocusList
	} else {
		f.currentArea = FocusDetail
	}
}

// Reset clears the current focus
func (f *FocusManager) Reset() {
	f.currentArea = FocusNone
}

// CanHandleKey returns true if the current focus area can handle keyboard input
func (f *FocusManager) CanHandleKey() bool {
	return f.currentArea != FocusNone
}

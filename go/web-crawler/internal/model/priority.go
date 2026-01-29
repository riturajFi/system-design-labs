package model

// Priority indicates crawl importance for scheduling.
type Priority int

const (
	PriorityHigh Priority = iota
	PriorityLow
)

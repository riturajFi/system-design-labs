package content

// Deduper decides whether content has been seen before.
type Deduper interface {
	Seen(hash string) bool
	Mark(hash string)
}

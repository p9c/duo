package status

// Status keeps track of errors on an ongoing basis and hooks into the logger which fills with snapshots of data state for debugging
type Status interface {
	SetStatus(string) Status
	UnsetStatus() Status
}

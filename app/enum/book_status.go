package enum

type BookingStatus string

const (
	BookingPending  BookingStatus = "Pending"
	BookingApproved BookingStatus = "Approved"
	BookingCanceled BookingStatus = "Canceled"
)

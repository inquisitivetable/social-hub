package enums

type PrivacyType int64

const (
	None PrivacyType = iota
	Public
	Private
	SubPrivate
)

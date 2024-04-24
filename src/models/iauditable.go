package models

type IAuditable interface {
	GetID() string
	GetObjectStr() string
	GetObjectName() string
}

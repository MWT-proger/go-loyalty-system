package models

type StatusOrder string

type BaseModeler interface {
	GetType() string
	GetArgsInsert() []any
}

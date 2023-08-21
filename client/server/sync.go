// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

// DB defines an embedded key/value store database interface.
type DB interface {
	Get(namespace, key []byte) (value []byte, err error)
	Set(namespace, key, value []byte) error
	Has(namespace, key []byte) (bool, error)
	Close() error
}

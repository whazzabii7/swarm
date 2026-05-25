package rpr

import "sync"

var requestPool = sync.Pool{
	New: func() any { return new(rawRequest) },
}

var responsePool = sync.Pool{
	New: func() any { return new(Response) },
}

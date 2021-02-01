// +build race

package json

import (
	"sync"
	"unsafe"
)

var setsMu sync.RWMutex

func encodeCompileToGetCodeSet(typeptr uintptr) (*opcodeSet, error) {
	if !existsCachedOpcodeSets {
		return encodeCompileToGetCodeSetSlowPath(typeptr)
	}
	setsMu.RLock()
	if codeSet := cachedOpcodeSets[typeptr-baseTypeAddr]; codeSet != nil {
		setsMu.RUnlock()
		return codeSet, nil
	}
	setsMu.RUnlock()

	// noescape trick for header.typ ( reflect.*rtype )
	copiedType := *(**rtype)(unsafe.Pointer(&typeptr))

	code, err := encodeCompileHead(&encodeCompileContext{
		typ:                      copiedType,
		root:                     true,
		structTypeToCompiledCode: map[uintptr]*compiledCode{},
	})
	if err != nil {
		return nil, err
	}
	code = copyOpcode(code)
	codeLength := code.totalLength()
	codeSet := &opcodeSet{
		code:       code,
		codeLength: codeLength,
	}
	setsMu.Lock()
	cachedOpcodeSets[int(typeptr-baseTypeAddr)] = codeSet
	setsMu.Unlock()
	return codeSet, nil
}

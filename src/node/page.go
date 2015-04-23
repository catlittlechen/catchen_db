package node

import (
	"unsafe"
)

type page struct {
	id    pid
	flag  uint8
	count uint16
	ptr   uintptr
}

const (
	pageHeaderSize = int(unsafe.Offsetof(((*page)(nil)).ptr))
)

func (p *page) nodeData() (nd *nodeData) {
	nd = (*nodeData)(unsafe.Pointer(&p.ptr))
	nd.ptr = p
	return
}

func (p *page) acNodePageElem() (ac *acNodePageElem) {
	ac = (*acNodePageElem)(unsafe.Pointer(&p.ptr))
	ac.ptr = p
	return
}

type pid uint64

type pids []pid

func (p pids) Len() int {
	return len(p)
}

func (p pids) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p pids) Less(i, j int) bool {
	return p[i] < p[j]
}

package node

import (
	"unsafe"
)

const (
	nodePageSize = uint32(unsafe.Sizeof(nodePageElem{}))
	NODE_RED     = 0
	NODE_BLACK   = 1

	NODE_KEY_SMALL = -1
	NODE_KEY_EQUAL = 0
	NODE_KEY_LARGE = 1
)

var (
	treeRoot *nodeRoot
)

type nodeRoot struct {
	node *nodePageElem
}

func (nr *nodeRoot) first(nodeIndex *nodePageElem) (node *nodePageElem) {
	node = nodeIndex

	if node != nil {
		for node.lChild != nil {
			node = node.lChild
		}
	}
	return
}

func (nr *nodeRoot) last(nodeIndex *nodePageElem) (node *nodePageElem) {
	node = nodeIndex

	if node != nil {
		for node.rChild != nil {
			node = node.rChild
		}
	}
	return
}

func (nr *nodeRoot) next(nodeIndex *nodePageElem) (node *nodePageElem) {

	node = nodeIndex
	if node.rChild != nil {
		node = nr.first(node.rChild)
		return
	}

	nodeTmp := node.parent
	for nodeTmp != nil && node == nodeTmp.rChild {
		node = nodeTmp
		nodeTmp = nodeTmp.parent
	}
	node = nodeTmp
	return
}

func (nr *nodeRoot) pre(nodeIndex *nodePageElem) (node *nodePageElem) {

	node = nodeIndex
	if node.lChild != nil {
		node = (nr).last(node.lChild)
		return
	}

	nodeTmp := node.parent
	for nodeTmp != nil && node == nodeTmp.lChild {
		node = nodeTmp
		nodeTmp = nodeTmp.parent
	}
	node = nodeTmp
	return
}

func (nr *nodeRoot) search(key string) (node *nodePageElem) {
	node = nr.node

	for node != nil {
		switch node.compareKey(key) {
		case NODE_KEY_EQUAL:
			return node
		case NODE_KEY_SMALL:
			node = node.rChild
		case NODE_KEY_LARGE:
			node = node.lChild
		}
	}
	return nil
}

func (nr *nodeRoot) leftRotate(node *nodePageElem) {

	nodeTmp := node.rChild
	node.rChild = nodeTmp.lChild
	if nodeTmp.lChild != nil {
		nodeTmp.lChild.parent = node
	}
	nodeTmp.parent = node.parent
	if node.parent == nil {
		nr.node = nodeTmp
	} else {
		if node.parent.lChild == node {
			node.parent.lChild = nodeTmp
		} else {
			node.parent.rChild = nodeTmp
		}
	}

	nodeTmp.lChild = node
	node.parent = nodeTmp
}

func (nr *nodeRoot) rightRotate(node *nodePageElem) {

	nodeTmp := node.lChild
	node.lChild = nodeTmp.rChild
	if nodeTmp.rChild != nil {
		nodeTmp.rChild.parent = node
	}
	nodeTmp.parent = node.parent
	if node.parent == nil {
		nr.node = nodeTmp
	} else {
		if node == node.parent.rChild {
			node.parent.rChild = nodeTmp
		} else {
			node.parent.lChild = nodeTmp
		}
	}

	nodeTmp.rChild = node
	node.parent = nodeTmp
}

func (nr *nodeRoot) insertFixTree(node *nodePageElem) {
	var parent, grandparent, uncle *nodePageElem

	for node.parent.isRed() {
		parent = node.parent
		grandparent = parent.parent

		if parent == grandparent.lChild {
			uncle = grandparent.rChild
			if uncle != nil && uncle.isRed() {
				uncle.setBlack()
				parent.setBlack()
				grandparent.setRed()
				node = grandparent
				continue
			}
			if parent.rChild == node {
				nr.leftRotate(parent)
				parent, node = node, parent
			}
			parent.setBlack()
			grandparent.setRed()
			nr.rightRotate(grandparent)
		} else {
			uncle = grandparent.lChild
			if uncle != nil && uncle.isRed() {
				uncle.setBlack()
				parent.setBlack()
				grandparent.setRed()
				node = grandparent
				continue
			}
			if parent.lChild == node {
				nr.rightRotate(parent)
				parent, node = node, parent
			}
			parent.setBlack()
			grandparent.setRed()
			nr.leftRotate(grandparent)
		}
	}

	nr.node.setBlack()
}

func (nr *nodeRoot) insert(node *nodePageElem) {
	var nodeX, nodeY *nodePageElem

	for nodeX != nil {
		nodeY = nodeX
		if node.compare(nodeX) {
			nodeX = nodeX.lChild
		} else {
			nodeX = nodeX.rChild
		}
	}
	node.parent = nodeY

	if nodeY != nil {
		if node.compare(nodeY) {
			nodeY.lChild = node
		} else {
			nodeY.rChild = node
		}
	} else {
		nr.node = node
	}

	node.setRed()
	nr.insertFixTree(node)
}

func (nr *nodeRoot) createNode(key, value string, parent, lChild, rChild *nodePageElem) (node *nodePageElem) {
	return
}

func (nr *nodeRoot) insertNode(key, value string) bool {

	node := nr.search(key)
	if node != nil {
		return node.setValue(value)
	}
	if node = nr.createNode(key, value, nil, nil, nil); node != nil {
		nr.insert(node)
		return true
	}
	return false
}

func (nr *nodeRoot) deleteFixTree(node, parent *nodePageElem) {
	for (node != nil || node.isBlack()) && node != nr.node {
		if parent.lChild == node {
			other := parent.rChild
			if other.isRed() {
				other.setBlack()
				parent.setRed()
				nr.leftRotate(parent)
				other = parent.rChild
			}
			if (other.lChild == nil || other.lChild.isBlack()) && (other.rChild == nil || other.rChild.isBlack()) {
				other.setRed()
				node = parent
				parent = node.parent
			} else {
				if other.rChild == nil || other.rChild.isBlack() {
					other.lChild.setBlack()
					other.setRed()
					nr.rightRotate(other)
					other = parent.rChild
				}
				if parent.isBlack() {
					other.setBlack()
				} else {
					other.setRed()
				}
				parent.setBlack()
				other.rChild.setBlack()
				nr.leftRotate(parent)
				node = nr.node
				break
			}
		} else {
			other := parent.lChild
			if other.isRed() {
				other.setBlack()
				parent.setRed()
				nr.rightRotate(parent)
				other = parent.lChild
			}
			if (other.lChild == nil || other.lChild.isBlack()) && (other.rChild == nil || other.rChild.isBlack()) {
				other.setRed()
				node = parent
				parent = node.parent
			} else {
				if other.lChild == nil || other.lChild.isBlack() {
					other.rChild.setBlack()
					other.setRed()
					nr.leftRotate(other)
					other = parent.lChild
				}
				if parent.isRed() {
					other.setRed()
				} else {
					other.setBlack()
				}
				parent.setBlack()
				other.lChild.setBlack()
				nr.rightRotate(parent)
				node = nr.node
				break
			}
		}
	}
	if node != nil {
		node.setBlack()
	}
}

func (nr *nodeRoot) delet(node *nodePageElem) {
	var child, parent, replace *nodePageElem
	var color bool

	if node.lChild != nil && node.rChild != nil {
		replace = node.rChild
		for replace.lChild != nil {
			replace = replace.lChild
		}
		if node.parent != nil {
			if node.parent.lChild == node {
				node.parent.lChild = replace
			} else {
				node.parent.rChild = replace
			}
		} else {
			nr.node = replace
		}

		child = replace.rChild
		parent = replace.parent
		color = replace.isRed()

		if parent == node {
			parent = replace
		} else {
			if child != nil {
				child.parent = parent
			}
			parent.lChild = child
			replace.rChild = node.rChild
			node.rChild.parent = replace
		}
		replace.parent = node.parent
		replace.colorType = node.colorType
		replace.lChild = node.lChild
		node.lChild.parent = replace

		if color == false {
			nr.deleteFixTree(child, parent)
		}
		//TODO:
		//释放空间
		return
	}
	if node.lChild != nil {
		child = node.lChild
	} else {
		child = node.rChild
	}

	parent = node.parent
	color = node.isRed()

	if child != nil {
		child.parent = parent
	}

	if parent != nil {
		if parent.lChild == node {
			parent.lChild = child
		} else {
			parent.rChild = child
		}
	} else {
		nr.node = child
	}
	if color == false {
		nr.deleteFixTree(child, parent)
	}
	//TODO:
	//释放空间
}

func (nr *nodeRoot) deleteNode(key string) {
	if node := nr.search(key); node != nil {
		nr.delet(node)
	}
}

type nodePageElem struct {
	colorType bool
	lChild    *nodePageElem
	rChild    *nodePageElem
	parent    *nodePageElem
	keySize   uint32
	valueSize uint32
}

func (n *nodePageElem) isRed() bool {
	return n.colorType
}

func (n *nodePageElem) isBlack() bool {
	return !n.colorType
}

func (n *nodePageElem) setRed() {
	n.colorType = true
}

func (n *nodePageElem) setBlack() {
	n.colorType = false
}

func (n *nodePageElem) compare(node *nodePageElem) bool {
	return true
}

func (n *nodePageElem) compareKey(key string) int {
	return NODE_KEY_EQUAL
}

func (n *nodePageElem) key() []byte {
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(n))
	return buf[nodePageSize : nodePageSize+n.keySize]
}

func (n *nodePageElem) value() []byte {
	buf := (*[maxAlloacSize]byte)(unsafe.Pointer(n))
	return buf[nodePageSize+n.keySize : nodePageSize+n.keySize+n.valueSize]
}

func (n *nodePageElem) setKey(key string) bool {
	return true
}

func (n *nodePageElem) setValue(value string) bool {
	return true
}

func init() {
	treeRoot = new(nodeRoot)
}

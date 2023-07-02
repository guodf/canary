package tree

type color bool

const (
	red   color = false
	black color = true
)

type ICompare interface {
	Compare(value interface{}) int
}

type rbTree struct {
	root  *rbTreeNode
	count int
}

type rbTreeNode struct {
	value      interface{}
	color      color
	leftNode   *rbTreeNode
	rightNode  *rbTreeNode
	parentNode *rbTreeNode
}

func NewRBTree() *rbTree {
	return &rbTree{}
}

func (tree *rbTree) Add(value ICompare) {
	newNode := &rbTreeNode{
		value:      value,
		color:      red,
		leftNode:   nil,
		rightNode:  nil,
		parentNode: nil,
	}
	if tree.root == nil {
		tree.root = newNode
	} else {
		tree.root.add(newNode)
	}
	newNode.fixColor()
	tree.count += 1
	tree.resetRoot()
}

func (tree *rbTree) resetRoot() {
	tree.root = tree.root.resetRoot()
}

func (tree *rbTree) Del(value interface{}) {
	if tree.root == nil {
		return
	}
	pNode := tree.root.parentNode
	// 如果无父，说明是根节点且为黑色
	if pNode == nil || pNode.value.(ICompare).Compare(value) == 0 {
		tree.root = nil
		return
	}
	nodes := tree.root.find(value.(ICompare), false)
	for _, node := range nodes {
		node.del()
	}
}

func (cNode *rbTreeNode) add(node *rbTreeNode) {
	r := compare(cNode, node)
	if r > 0 {
		if cNode.leftNode == nil {
			cNode.leftNode = node
			node.parentNode = cNode
		} else {
			cNode.leftNode.add(node)
		}
	} else {
		if cNode.rightNode == nil {

			cNode.rightNode = node
			node.parentNode = cNode
		} else {
			cNode.rightNode.add(node)
		}
	}
}

func (cNode *rbTreeNode) fixColor() {
	pNode := cNode.parentNode
	if pNode == nil {
		cNode.color = black
		return
	}
	gNode := pNode.parentNode
	if gNode == nil {
		return
	}
	// 1. 若当前节点为红色，且叔父节点亦为红色时需要修复
	if gNode.leftNode != nil && gNode.rightNode != nil && gNode.leftNode.color == gNode.rightNode.color && gNode.leftNode.color == red {
		gNode.color = red
		gNode.leftNode.color = black
		gNode.rightNode.color = black
		gNode.fixColor()
		return
	}
	// 父子节点同为红色，叔叔为黑色
	if pNode.color == cNode.color && cNode.color == red {
		// 需要爷爷节点参与旋转的情况：祖孙三代在一条直线上，先将父染为黑色，再将爷爷染为红色 ，该种情况旋转后无需修复颜色
		//   2. 同为右子树，左旋
		//   3. 同为左子树，右旋旋
		// 不需要爷爷几点参与旋转的情况：祖孙三代不在一条直线上,直接旋转父子即可
		//   4. 子为右子树，左旋
		//   5. 子为左子树，右旋
		pIsRight := isRightNode(pNode)
		cIsRight := isRightNode(cNode)
		if pIsRight && cIsRight { // 2. 同为右子树，左旋
			pNode.color = black
			gNode.color = red
			pNode.leftRotate()
		} else if !pIsRight && !cIsRight { //3. 同为左子树，右旋旋
			pNode.color = black
			gNode.color = red
			pNode.rightRotate()
		} else if cIsRight { //4. 子为右子树，左旋
			cNode.leftRotate()
			pNode.fixColor()
		} else { // 5. 子为左子树，右旋
			cNode.rightRotate()
			pNode.fixColor()
		}
	}
}

// 左旋
func (cNode *rbTreeNode) leftRotate() {
	pNode := cNode.parentNode
	gNode := pNode.parentNode
	if gNode != nil {
		if isRightNode(pNode) {
			gNode.rightNode = cNode
		} else {
			gNode.leftNode = cNode
		}
	}
	cNode.parentNode = pNode.parentNode
	pNode.rightNode = cNode.leftNode
	cNode.leftNode = pNode
	pNode.parentNode = cNode
}

// 右旋
func (cNode *rbTreeNode) rightRotate() {
	pNode := cNode.parentNode
	gNode := pNode.parentNode
	if gNode != nil {
		if isRightNode(pNode) {
			gNode.rightNode = cNode
		} else {
			gNode.leftNode = cNode
		}
	}
	cNode.parentNode = pNode.parentNode
	pNode.leftNode = cNode.rightNode
	cNode.rightNode = pNode
	pNode.parentNode = cNode
}

// 重置root
func (cNode *rbTreeNode) resetRoot() *rbTreeNode {
	if cNode.parentNode != nil {
		return cNode.parentNode.resetRoot()
	}
	return cNode
}

func (cNode *rbTreeNode) find(value ICompare, all bool) []*rbTreeNode {
	var nodes []*rbTreeNode
	tmpNode := cNode
	// 非递归查找
	for {
		r := tmpNode.value.(ICompare).Compare(value)
		if r == 0 {
			nodes = append(nodes, tmpNode)
			// 是否查找所有的相同的值
			if all {
				tmpNode = tmpNode.rightNode
			} else {
				tmpNode = nil
			}
		} else if r > 0 {
			tmpNode = cNode.leftNode
		} else {
			tmpNode = cNode.rightNode
		}
		if tmpNode == nil {
			break
		}
	}
	return nodes
}

func (cNode *rbTreeNode) del() {
	// 1. 如果是红节点，且无子
	if isRedLonely(cNode) {
		cNode.cutOff()
		return
	}
	// 2. 如果有子，则转换成无子的情况处理,从它的前驱或者后继中招一个红色的，如果没有红色的，随便使用一个
	// 只所以要找红色的是因为红色的无子节点删除不会破坏平衡，就可以避免平衡修复，但是为了寻找红色会多出一次查找
	predNode := cNode.pred()
	if predNode != nil {
		// 仅交换值，保留颜色
		cNode.value, predNode.value = predNode.value, cNode.value
		// 转化为删除这个孤节点
		predNode.del()
		return
	}
	succNode := cNode.succ()
	if succNode != nil {
		// 仅交换值，保留颜色
		cNode.value, succNode.value = succNode.value, cNode.value
		// 转化为删除这个孤节点
		succNode.del()
		return
	}
	// 如果以上条件都不满足，则说明当前节点一定时黑色的

	pNode := cNode.parentNode
	isRight := isRightNode(cNode)
	var bNode *rbTreeNode
	//此时需考虑以下几种情况：
	if isRight {
		bNode = pNode.leftNode
	} else {
		bNode = pNode.rightNode
	}
	if bNode.color == red {
		// 1. 兄为红：
		// 		特点：这种情况下，父一定为黑，兄一定有两个子且为黑色
		// 		修正：父兄交换颜色，若兄在右则父兄左旋，若兄在左父兄右旋
		pNode.color, bNode.color = bNode.color, pNode.color
		if isRight {
			bNode.leftRotate()
		} else {
			bNode.rightRotate()
		}
	} else if pNode.color == red {
		// 2. 父为红:
		//		特点：这种情况下兄弟一定为黑色，侄子一定为红色
		// 2.1 兄无子：
		// 修正：交换父子颜色或者旋转父子，若兄在右则左旋父子，若兄在左则右旋父子
		if !bNode.hasChildren() {
			pNode.color, bNode.color = bNode.color, pNode.color
		} else {
			//2.2.1 兄在右，兄有右子
			//		修正：父和右侄子染黑，兄染红，左旋父兄
			//2.2.2 兄在右，兄仅有左子
			//		修正：左侄子染黑，兄染红，先右旋兄与左侄子，再左旋父与新右子
			//2.3.1 兄在左，兄有左子
			//		修正：父和左侄子染黑，兄染红，右旋父兄
			//2.3.2 兄在左，兄仅有右子
			//		修正：右侄子染黑，兄染红，先左旋兄与右侄子，再右旋父与新左子
			if !isRight {
				if bNode.rightNode != nil {
					pNode.color = black
					bNode.rightNode.color = black
					bNode.color = red
					bNode.leftRotate()
				} else {
					bNode.leftNode.color = black
					bNode.color = red
					bNode.leftNode.rightRotate()
					pNode.rightNode.leftRotate()
				}
			} else {
				if bNode.leftNode != nil {
					pNode.color = black
					bNode.leftNode.color = black
					bNode.color = red
					bNode.rightRotate()
				} else {
					bNode.rightNode.color = black
					bNode.color = red
					bNode.rightNode.leftRotate()
					bNode.leftNode.rightRotate()
				}
			}

		}
	} else {
		// 父子皆为黑
		// 3.1  兄无子
		// 		修正：兄染红，以父为当前节点染红兄弟（递归），直到根节点或者兄本身就是红时终止
		if !bNode.hasChildren() {
			bNode.color = red
			pNodeTmp := pNode
			gNodeTmp := pNode.parentNode
			for {
				if gNodeTmp.parentNode != nil {
					var bNodeTmp *rbTreeNode
					if isRightNode(pNodeTmp) {
						bNodeTmp = gNodeTmp.leftNode
					} else {
						bNodeTmp = gNodeTmp.rightNode
					}
					if bNodeTmp.color == red {
						break
					}
					bNodeTmp.color = red
					if gNodeTmp.color == red {
						gNodeTmp.color = black
						break
					}
					pNodeTmp = gNodeTmp
					gNodeTmp = gNodeTmp.parentNode
				}
			}
		} else {
			// 3.2.1 兄在右，兄有右子
			//		修正：右侄子染黑，左旋父兄
			// 3.2.2 兄在右，兄仅有左子
			//		修正：左侄子染黑，先左旋兄与其左子，再右旋父与右子
			// 3.3.1 兄在左，兄有左子
			//		修正：左侄子染黑，右旋父兄
			// 3.3.2 兄在左，兄仅有右子
			//		修正：右侄子染黑，先左旋兄与其右子，再右旋父与新左子
			if !isRight {
				if bNode.rightNode != nil {
					bNode.rightNode.color = black
					bNode.leftRotate()
				} else {
					bNode.leftNode.color = black
					bNode.leftNode.leftRotate()
					pNode.rightNode.rightRotate()
				}
			} else {
				if bNode.leftNode != nil {
					bNode.leftNode.color = black
					bNode.rightRotate()
				} else {
					bNode.rightNode.color = black
					bNode.rightNode.leftRotate()
					pNode.leftNode.rightRotate()
				}
			}
		}
	}

	cNode.cutOff()
}

// 查找前驱
func (cNode *rbTreeNode) pred() *rbTreeNode {
	node := cNode.leftNode
	if node == nil {
		return nil
	}
	rNode := node.rightNode
	for {
		if rNode == nil {
			return node
		}
		rNode = rNode.rightNode
		node = rNode
	}
}

// 查找后继
func (cNode *rbTreeNode) succ() *rbTreeNode {
	node := cNode.rightNode
	if node == nil {
		return nil
	}
	lNode := node.leftNode
	for {
		if lNode == nil {
			return node
		}
		lNode = lNode.leftNode
		node = lNode
	}
}

// 判断有无后代
func (cNode *rbTreeNode) hasChildren() bool {
	return cNode.leftNode != nil || cNode.rightNode != nil
}

// 隔断父子关系
func (cNode *rbTreeNode) cutOff() {
	if cNode.parentNode != nil {
		if isRightNode(cNode) {
			cNode.parentNode.rightNode = nil
		} else {
			cNode.parentNode.leftNode = nil
		}
		cNode.parentNode = nil
	}
}

// 判断两个节点的值
func compare(node1 *rbTreeNode, node2 *rbTreeNode) int {
	return node1.value.(ICompare).Compare(node2.value)
}

// 判断是否为右子树，不为右子树则必然为左子树，不考虑节点为根的情况
func isRightNode(node *rbTreeNode) bool {
	// 根节点既可以是左子树也可以是右子树
	if node.parentNode == nil {
		return true
	}
	// 父的右子树为nil，则该子一定不是右子树
	if node.parentNode.rightNode == nil {
		return false
	}
	return compare(node.parentNode.rightNode, node) == 0
}

// 红色孤节点
func isRedLonely(node *rbTreeNode) bool {
	return node.color == red && node.leftNode == nil && node.rightNode == nil
}

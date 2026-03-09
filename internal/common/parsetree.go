package common

import "fmt"

type ParseTreeNode struct {
	InnerToken Token
	ChildNodes []ParseTreeNode
}

func (n ParseTreeNode) Display(start, increase string) {
	fmt.Println(
		start,
		NameMapWithTokenKind[n.InnerToken.TokenKind],
		n.InnerToken.Token,
	)
	for _, t := range n.ChildNodes {
		t.Display(start+increase, increase)
	}
}

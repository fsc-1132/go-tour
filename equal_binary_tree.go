package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk ���� tree t �����е�ֵ�� tree ���͵� channel ch��
func Walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	Walk(t.Left, ch)
	ch <- t.Value
	Walk(t.Right, ch)
}

// Same ����� t1 �� t2 �Ƿ�����ͬ��ֵ��
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for i := 0; i < 10; i++ {
		x, y := <-ch1, <-ch2
		fmt.Println(x, y)
		if x != y {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(Same(tree.New(1), tree.New(2)))
	fmt.Println(Same(tree.New(1), tree.New(1)))
}

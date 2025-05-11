package customList

// Структура Linked List - это ее корень
type LinkedList struct {
	Root *Node
}

// Передаем по указателям, т.к. нужно будет менять значения
type Node struct {
	Next     *Node
	Data     *NodeData
	InAction bool
}

// Данные конкретной Node
type NodeData struct {
	Value       float64
	Operation   rune
	IsOperation bool
}

// Создание нового LinkedList с nil root
func NewLinkedList() *LinkedList {
	return &LinkedList{
		nil,
	}
}

// Добавление нового элемента в Linked List
func (list *LinkedList) Add(data *NodeData) error {
	list.Root = &Node{
		Next: list.Root,
		Data: data,
	}
	return nil
}

package tree
/* 
#include <libxml/tree.h>
*/
import "C"
import "unsafe"

type Element struct {
	*XmlNode
}

func (node *Element) ElementType() int {
	elem := (*C.xmlElement)(unsafe.Pointer(node.ptr()))
	return int(elem.etype)
}

func (node *Element) new(ptr *C.xmlNode) *Element {
	if ptr == nil {
		return nil
	}
	return NewNode(unsafe.Pointer(ptr), node.Doc()).(*Element)
}

func (node *Element) NextElement() *Element {
	return node.new(C.xmlNextElementSibling(node.NodePtr))
}

func (node *Element) PrevElement() *Element {
	return node.new(C.xmlPreviousElementSibling(node.NodePtr))
}

func (node *Element) FirstElement() *Element {
	return node.new(C.xmlFirstElementChild(node.NodePtr))
}

func (node *Element) LastElement() *Element {
	return node.new(C.xmlLastElementChild(node.NodePtr))
}

func (node *Element) Clear() {
	// Remember, as we delete them, the last one moves to the front
	child := node.First()
	for child != nil {
		child.Remove()
		child.Free()
		child = node.First()
	}
}

func (node *Element) Content() string {
	child := node.First()
	output := ""
	for child != nil {
		output = output + child.DumpHTML()
		child = child.Next()
	}
	return output
}

func (node *Element) SetContent(content string) {
	node.Clear()
	node.AppendContent(content)
}

func (node *Element) AppendContent(content string) {
	newDoc := XmlParseFragmentWithOptions(content, "", "", 
        XML_PARSE_RECOVER | 
        XML_PARSE_NONET|
        XML_PARSE_NOERROR|
        XML_PARSE_NOWARNING)

	defer newDoc.Free()
	child := newDoc.RootElement().First()
	for child != nil {
		//need to save the next sibling before appending it,
		//because once it loses its link to the next sibling in its original tree once appended to the new doc
		nextChild := child.Next()
		node.AppendChildNode(child)
		child = nextChild
	}
}

func (node *Element) PrependContent(content string) {
	newDoc := XmlParseFragmentWithOptions(content, "", "", 
        XML_PARSE_RECOVER | 
        XML_PARSE_NONET|
        XML_PARSE_NOERROR|
        XML_PARSE_NOWARNING)

	defer newDoc.Free()
	child := newDoc.RootElement().Last()
	for child != nil {
		prevChild := child.Prev()
		node.PrependChildNode(child)
		child = prevChild
	}
}

func (node *Element) AddContentAfter(content string) {
    newDoc := XmlParseFragmentWithOptions(content, "", "", 
        XML_PARSE_RECOVER | 
        XML_PARSE_NONET|
        XML_PARSE_NOERROR|
        XML_PARSE_NOWARNING)
    defer newDoc.Free()
	child := newDoc.Parent().Last()
	for child != nil {
		node.AddNodeAfter(child)
		child = child.Prev()
	}
}
func (node *Element) AddContentBefore(content string) {
    newDoc := XmlParseFragmentWithOptions(content, "", "", 
        XML_PARSE_RECOVER | 
        XML_PARSE_NONET|
        XML_PARSE_NOERROR|
        XML_PARSE_NOWARNING)
    defer newDoc.Free()

	child := newDoc.Parent().First()
	for child != nil {
		node.AddNodeBefore(child)
		child = child.Next()
	}
}

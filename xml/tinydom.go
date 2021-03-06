/*
Package tinydom	实现了一个简单的XML的DOM树构造工具.
*/
package tinydom

import (
    "bytes"
    "encoding/xml"
    "errors"
    "io"
)

//  XMLAttribute    是一个元素的属性的接口
type XMLAttribute interface {
    Name() string
    Value() string
    SetValue(string)
}

//  XMLNode 定义了XML所有节点的基础设施，提供了基本的元素遍历、增删等操作,也提供了逆向转换能力.
type XMLNode interface {
    ToElement() XMLElement
    ToText() XMLText
    ToComment() XMLComment
    ToDocument() XMLDocument
    ToProcInst() XMLProcInst
    ToDirective() XMLDirective

    Value() string
    SetValue(newValue string)

    GetDocument() XMLDocument

    NoChildren() bool
    Parent() XMLNode
    FirstChild() XMLNode
    LastChild() XMLNode
    PreviousSibling() XMLNode
    NextSibling() XMLNode
    FirstChildElement(name string) XMLElement
    LastChildElement(name string) XMLElement
    PreviousSiblingElement(name string) XMLElement
    NextSiblingElement(name string) XMLElement

    InsertEndChild(node XMLNode) XMLNode
    InsertFirstChild(node XMLNode) XMLNode
    InsertAfterChild(afterThis XMLNode, addThis XMLNode) XMLNode
    DeleteChildren()
    DeleteChild(node XMLNode)
    Accept(visitor XMLVisitor) bool

    //  被迫入侵的接口
    setParent(node XMLNode)
    setPrev(node XMLNode)
    setNext(node XMLNode)

    unlink(child XMLNode)
}

//  XMLElement  提供了访问XML基本节点元素的能力
//
//  Name、SetName其实是Value和SetValue的别名，目的是为了使得接口更加符合直观理解。
//
//  Text、SetText的作用是设置<node>与</node>之间的文字，虽然文字都是有XMLText对象来承载的，但是通常来说直接在XMLElement中访问会更加方便。
//
//  FindAttribute和ForeachAttribute分别用于查找特定的XML节点的属性和遍历XML属性列表。
//
//  Attribute、SetAttribute、DeleteAttribute用于读取和删除属性。
type XMLElement interface {
    XMLNode

    Name() string
    SetName(name string)

    FindAttribute(name string) XMLAttribute
    ForeachAttribute(callback func(attribute XMLAttribute) int) int

    AttributeCount() int
    Attribute(name string, def string) string
    SetAttribute(name string, value string) XMLAttribute
    DeleteAttribute(name string) XMLAttribute
    ClearAttributes()

    Text() string
    SetText(text string)
}

//  XMLText 提供了对XML元素间文本的封装
type XMLText interface {
    XMLNode
    SetCDATA(isCData bool)
    CDATA() bool
}

type XMLComment interface {
    XMLNode
    Comment() string
    SetComment(string)
}

type XMLProcInst interface {
    XMLNode
    Target() string
    Instruction() string
}

type XMLDirective interface {
    XMLNode
}

type XMLDocument interface {
    XMLNode
}

type XMLVisitor interface {
    VisitEnterDocument(XMLDocument) bool
    VisitExitDocument(XMLDocument) bool

    VisitEnterElement(XMLElement) bool
    VisitExitElement(XMLElement) bool

    VisitProcInst(XMLProcInst) bool
    VisitText(XMLText) bool
    VisitComment(XMLComment) bool
    VisitDirective(XMLDirective) bool
}

type XMLHandle interface {
    //Root() XMLCursor
    //
    Parent() XMLHandle
    FirstChild() XMLHandle
    LastChild() XMLHandle
    PreviousSibling() XMLHandle
    NextSibling() XMLHandle
    FirstChildElement(name string) XMLHandle
    LastChildElement(name string) XMLHandle
    PreviousSiblingElement(name string) XMLHandle
    NextSiblingElement(name string) XMLHandle

    ToNode() XMLNode
    ToElement() XMLElement
    ToText() XMLText
    ToComment() XMLComment
    ToDocument() XMLDocument
    ToProcInst() XMLProcInst
    ToDirective() XMLDirective
}

//=========================================================

type xmlAttributeImpl struct {
    name  string
    value string
}

func (this *xmlAttributeImpl) Name() string {
    return this.name
}

func (this *xmlAttributeImpl) Value() string {
    return this.value
}

func (this *xmlAttributeImpl) SetValue(newValue string) {
    this.value = newValue
}

//==================================================================

type xmlNodeImpl struct {
    impl XMLNode

    document XMLDocument
    parent   XMLNode
    value    string

    firstChild XMLNode
    lastChild  XMLNode

    prev XMLNode
    next XMLNode
}

func (this *xmlNodeImpl) getDocument() XMLDocument {
    return this.document
}

func (this *xmlNodeImpl) setParent(node XMLNode) {
    this.parent = node
}

func (this *xmlNodeImpl) setPrev(node XMLNode) {
    this.prev = node
}

func (this *xmlNodeImpl) setNext(node XMLNode) {
    this.next = node
}

func (this *xmlNodeImpl) ToElement() XMLElement {
    return nil
}

func (this *xmlNodeImpl) ToText() XMLText {
    return nil
}

func (this *xmlNodeImpl) ToComment() XMLComment {
    return nil
}

func (this *xmlNodeImpl) ToDocument() XMLDocument {
    return nil
}

func (this *xmlNodeImpl) ToProcInst() XMLProcInst {
    return nil
}

func (this *xmlNodeImpl) ToDirective() XMLDirective {
    return nil
}

func (this *xmlNodeImpl) Value() string {
    return this.value
}

func (this *xmlNodeImpl) SetValue(newValue string) {
    this.value = newValue
}

func (this *xmlNodeImpl) GetDocument() XMLDocument {
    return this.document
}

func (this *xmlNodeImpl) Parent() XMLNode {
    return this.parent
}

func (this *xmlNodeImpl) NoChildren() bool {
    return nil == this.firstChild
}

func (this *xmlNodeImpl) FirstChild() XMLNode {
    return this.firstChild
}

func (this *xmlNodeImpl) LastChild() XMLNode {
    return this.lastChild
}

func (this *xmlNodeImpl) PreviousSibling() XMLNode {
    return this.prev
}

func (this *xmlNodeImpl) NextSibling() XMLNode {
    return this.next
}

func (this *xmlNodeImpl) FirstChildElement(name string) XMLElement {
    for item := this.firstChild; nil != item; item = item.NextSibling() {
        elem := item.ToElement()
        if nil == elem {
            continue
        }

        if ("" == name) || (elem.Name() == name) {
            return elem
        }
    }

    return nil
}

func (this *xmlNodeImpl) LastChildElement(name string) XMLElement {

    for item := this.lastChild; nil != item; item = item.PreviousSibling() {
        elem := item.ToElement()
        if nil == elem {
            continue
        }

        if ("" == name) || (elem.Name() == name) {
            return elem
        }
    }

    return nil
}

func (this *xmlNodeImpl) PreviousSiblingElement(name string) XMLElement {
    for item := this.prev; nil != item; item = item.PreviousSibling() {
        elem := item.ToElement()
        if nil == elem {
            continue
        }

        if ("" == name) || (elem.Name() == name) {
            return elem
        }
    }

    return nil
}

func (this *xmlNodeImpl) NextSiblingElement(name string) XMLElement {

    for item := this.next; nil != item; item = item.NextSibling() {
        elem := item.ToElement()
        if nil == elem {
            continue
        }

        if ("" == name) || (elem.Name() == name) {
            return elem
        }
    }

    return nil
}

func (this *xmlNodeImpl) unlink(child XMLNode) {
    if child == this.firstChild {
        this.firstChild = this.firstChild.NextSibling()
    }

    if child == this.lastChild {
        this.lastChild = this.lastChild.PreviousSibling()
    }

    if nil != child.PreviousSibling() {
        child.PreviousSibling().setNext(child.NextSibling())
    }

    if nil != child.NextSibling() {
        child.NextSibling().setPrev(child.PreviousSibling())
    }

    child.setParent(nil)
}

func (this *xmlNodeImpl) InsertEndChild(addThis XMLNode) XMLNode {
    if addThis.GetDocument() != this.document {
        return nil
    }

    if nil != addThis.Parent() {
        addThis.Parent().unlink(addThis)
    }

    if nil != this.lastChild {
        this.lastChild.setNext(addThis)
        addThis.setPrev(this.lastChild)
        this.lastChild = addThis
        addThis.setNext(nil)
    } else {
        this.firstChild = addThis
        this.lastChild = addThis

        addThis.setPrev(nil)
        addThis.setNext(nil)
    }

    addThis.setParent(this.impl)
    return addThis
}

func (this *xmlNodeImpl) InsertFirstChild(addThis XMLNode) XMLNode {
    if addThis.GetDocument() != this.document {
        return nil
    }

    if nil != addThis.Parent() {
        addThis.Parent().unlink(addThis)
    }

    if nil != this.firstChild {
        this.firstChild.setPrev(addThis)
        addThis.setNext(this.firstChild)
        this.firstChild = addThis
        addThis.setPrev(nil)
    } else {
        this.firstChild = addThis
        this.lastChild = addThis

        addThis.setPrev(nil)
        addThis.setNext(nil)
    }

    addThis.setParent(this.impl)
    return addThis
}

func (this *xmlNodeImpl) InsertAfterChild(afterThis XMLNode, addThis XMLNode) XMLNode {
    if addThis.GetDocument() != this.document {
        return nil
    }

    if afterThis.Parent() != this.impl {
        return nil
    }

    if afterThis.NextSibling() == nil {
        return this.InsertEndChild(addThis)
    }

    if nil != addThis.Parent() {
        addThis.Parent().unlink(addThis)
    }

    addThis.setPrev(afterThis)
    addThis.setNext(afterThis.NextSibling())
    afterThis.NextSibling().setPrev(addThis)
    afterThis.setNext(addThis)
    addThis.setParent(this.impl)

    return addThis
}

func (this *xmlNodeImpl) DeleteChildren() {
    for nil != this.firstChild {
        this.DeleteChild(this.firstChild)
    }

    this.firstChild = nil
    this.lastChild = nil
}

func (this *xmlNodeImpl) DeleteChild(node XMLNode) {
    this.unlink(node)
}

func (this *xmlNodeImpl) Accept(visitor XMLVisitor) bool {
    return false
}

//------------------------------------------------------------------

type xmlElementImpl struct {
    xmlNodeImpl

    //rootAttribute XMLAttribute
    attributes map[string]XMLAttribute
}

func (this *xmlElementImpl) ToElement() XMLElement {
    return this
}

func (this *xmlElementImpl) Accept(visitor XMLVisitor) bool {

    if visitor.VisitEnterElement(this) {
        for node := this.FirstChild(); nil != node; node = node.NextSibling() {
            if !node.Accept(visitor) {
                break
            }
        }
    }

    return visitor.VisitExitElement(this)
}

func (this *xmlElementImpl) Name() string {
    return this.Value()
}

func (this *xmlElementImpl) SetName(name string) {
    this.SetValue(name)
}

func (this *xmlElementImpl) FindAttribute(name string) XMLAttribute {
    if nil == this.attributes {
        return nil
    }

    attr, ok := this.attributes[name]
    if !ok {
        return nil
    }

    return attr
}

func (this *xmlElementImpl) AttributeCount() int {
    if nil == this.attributes {
        return 0
    }
    return len(this.attributes)
}

func (this *xmlElementImpl) Attribute(name string, def string) string {
    if nil == this.attributes {
        return def
    }

    attr, ok := this.attributes[name]
    if !ok {
        return def
    }

    return attr.Value()
}

func (this *xmlElementImpl) SetAttribute(name string, value string) XMLAttribute {
    if nil == this.attributes {
        this.attributes = make(map[string]XMLAttribute)
        attr := newAttribute(name, value)
        this.attributes[name] = attr
        return attr
    }

    attr, ok := this.attributes[name]
    if ok {
        attr.SetValue(value)
        return attr
    }

    attr = newAttribute(name, value)
    this.attributes[name] = attr
    return attr
}

func (this *xmlElementImpl) DeleteAttribute(name string) XMLAttribute {
    attr := this.FindAttribute(name)
    if nil == attr {
        return nil
    }
    delete(this.attributes, name)
    return attr
}

func (this *xmlElementImpl) Text() string {
    if text := this.FirstChild(); (nil != text) && (nil != text.ToText()) {
        return text.Value()
    }

    return ""
}

func (this *xmlElementImpl) SetText(inText string) {
    if node := this.FirstChild(); (nil != node) && (nil != node.ToText()) {
        node.SetValue(inText)
    } else {
        theText := NewText(this.getDocument(), inText)
        this.InsertFirstChild(theText)
    }
}

func (this *xmlElementImpl) ForeachAttribute(callback func(attribute XMLAttribute) int) int {
    if nil == this.attributes {
        return 0
    }

    for _, value := range this.attributes {
        if ret := callback(value); 0 != ret {
            return ret
        }
    }

    return 0
}

func (this *xmlElementImpl) ClearAttributes() {
    this.attributes = nil
}

//------------------------------------------------------------------

type xmlCommentImpl struct {
    xmlNodeImpl
}

func (this *xmlCommentImpl) ToComment() XMLComment {
    return this
}

func (this *xmlCommentImpl) Comment() string {
    return this.value
}

func (this *xmlCommentImpl) SetComment(newComment string) {
    this.value = newComment
}

func (this *xmlCommentImpl) Accept(visitor XMLVisitor) bool {
    return visitor.VisitComment(this)
}

//------------------------------------------------------------------

type xmlProcInstImpl struct {
    xmlNodeImpl
    instruction string
}

func (this *xmlProcInstImpl) ToProcInst() XMLProcInst {
    return this
}

func (this *xmlProcInstImpl) Accept(visitor XMLVisitor) bool {
    return visitor.VisitProcInst(this)
}

func (this *xmlProcInstImpl) Target() string {
    return this.value
}

func (this *xmlProcInstImpl) Instruction() string {
    return this.instruction
}

//------------------------------------------------------------------

type xmlDocumentImpl struct {
    xmlNodeImpl
}

func (this *xmlDocumentImpl) ToDocument() XMLDocument {
    return this
}

func (this *xmlDocumentImpl) Accept(visitor XMLVisitor) bool {

    if visitor.VisitEnterDocument(this) {
        for node := this.FirstChild(); nil != node; node = node.NextSibling() {
            if !node.Accept(visitor) {
                break
            }
        }
    }

    return visitor.VisitExitDocument(this)
}

//------------------------------------------------------------------

type xmlTextImpl struct {
    xmlNodeImpl
    cdata bool
}

func (this *xmlTextImpl) ToText() XMLText {
    return this
}
func (this *xmlTextImpl) Accept(visitor XMLVisitor) bool {
    return visitor.VisitText(this)
}
func (this *xmlTextImpl) SetCDATA(isCData bool) {
    this.cdata = isCData
}
func (this *xmlTextImpl) CDATA() bool {
    return this.cdata
}

//------------------------------------------------------------------

type xmlDirectiveImpl struct {
    xmlNodeImpl
}

func (this *xmlDirectiveImpl) ToDirective() XMLDirective {
    return this
}

func (this *xmlDirectiveImpl) Accept(visitor XMLVisitor) bool {
    return visitor.VisitDirective(this)
}

//------------------------------------------------------------------

//	NewText	创建一个新的XMLText对象
func NewText(document XMLDocument, text string) XMLText {
    node := new(xmlTextImpl)
    node.impl = node
    node.document = document
    node.value = text
    return node
}

//	XMLComment	创建一个新的XMLComment对象
func NewComment(document XMLDocument, comment string) XMLComment {
    node := new(xmlCommentImpl)
    node.impl = node
    node.document = document
    node.value = comment
    return node
}

//	NewElement	创建一个新的XMLElement对象
func NewElement(document XMLDocument, name string) XMLElement {
    node := new(xmlElementImpl)
    node.impl = node
    node.document = document
    node.value = name
    node.attributes = make(map[string]XMLAttribute)
    return node
}

//	NewProcInst	创建一个新的XMLProcInst对象
func NewProcInst(document XMLDocument, target string, inst string) XMLProcInst {
    node := new(xmlProcInstImpl)
    node.impl = node
    node.document = document
    node.value = target
    node.instruction = inst
    return node
}

//	NewDirective	创建一个新的XMLDirective对象
func NewDirective(document XMLDocument, directive string) XMLDirective {
    node := new(xmlDirectiveImpl)
    node.impl = node
    node.document = document
    node.value = directive
    return node
}

//	newAttribute	创建一个新的XMLAttribute对象.
//	name和value分别用于指定属性的名称和值
func newAttribute(name string, value string) XMLAttribute {
    attr := new(xmlAttributeImpl)
    attr.name = name
    attr.value = value
    return attr
}

//	NewDocument	创建一个全新的XMLDocument对象
func NewDocument() XMLDocument {
    node := new(xmlDocumentImpl)
    node.impl = node
    node.document = node
    return node
}

//	LoadDocument	从rd流中读取XML码流并构建成XMLDocument对象
func LoadDocument(rd io.Reader) (XMLDocument, error) {
    doc := NewDocument()
    var parent XMLNode = doc
    decoder := xml.NewDecoder(rd)
    var token xml.Token
    var err error
    rootElemExist := false
    for token, err = decoder.Token(); nil == err; token, err = decoder.Token() {
        switch token.(type) {
        case xml.StartElement:
            startElement := token.(xml.StartElement)

            //  一个XML文档只允许有唯一一个根节点
            if doc == parent {
                if rootElemExist {
                    return nil, errors.New("Root element has been exist:" + startElement.Name.Local)
                }

                //  标记一下根节点已经存在了
                rootElemExist = true
            }

            node := NewElement(doc, startElement.Name.Local)
            for _, item := range startElement.Attr {
                if nil != node.FindAttribute(item.Name.Local) {
                    return nil, errors.New("Attributes have the same name:" + item.Name.Local)
                }
                node.SetAttribute(item.Name.Local, item.Value)
            }
            parent.InsertEndChild(node)
            parent = node

        case xml.EndElement:
            //endElement := token.(xml.EndElement)
            parent = parent.Parent()
        case xml.Comment:
            comment := token.(xml.Comment)
            node := NewComment(doc, string(comment))
            parent.InsertEndChild(node)
        case xml.Directive:
            directive := token.(xml.Directive)
            node := NewDirective(doc, string(directive))
            parent.InsertEndChild(node)
        case xml.ProcInst:
            procInst := token.(xml.ProcInst)
            node := NewProcInst(doc, procInst.Target, string(procInst.Inst))
            parent.InsertEndChild(node)
        case xml.CharData:
            charData := token.(xml.CharData)
            shortCharData := bytes.TrimSpace(charData)
            if (nil != shortCharData) && (len(shortCharData) > 0) {
                if doc == parent {
                    return nil, errors.New("Text should be in the element")
                }

                node := NewText(doc, string(charData))
                parent.InsertEndChild(node)
            }
        default:
            return nil, errors.New("Unsupported token type")
        }
    }

    if (nil == err) || (io.EOF == err) {
        //  不能是空文档
        if nil == doc.FirstChildElement("") {
            return nil, errors.New("XML document missing the root element")
        }

        return doc, nil
    }

    return nil, err
}

//------------------------------------------------------------------
type xmlSimplePrinter struct {
    writer io.Writer
}

func NewSimplePrinter(writer io.Writer) XMLVisitor {
    visitor := new(xmlSimplePrinter)
    visitor.writer = writer
    return visitor
}

func (this *xmlSimplePrinter) VisitEnterDocument(node XMLDocument) bool {
    return true
}

func (this *xmlSimplePrinter) VisitExitDocument(node XMLDocument) bool {
    return true
}

func (this *xmlSimplePrinter) VisitEnterElement(node XMLElement) bool {
    io.WriteString(this.writer, "<")
    io.WriteString(this.writer, node.Name())

    node.ForeachAttribute(func(attribute XMLAttribute) int {
        io.WriteString(this.writer, ` `)
        io.WriteString(this.writer, attribute.Name())
        io.WriteString(this.writer, `="`)
        xml.EscapeText(this.writer, []byte(attribute.Value()))
        io.WriteString(this.writer, `"`)
        return 0
    })

    if node.NoChildren() {
        io.WriteString(this.writer, "/>")

        return true
    }

    io.WriteString(this.writer, ">")
    return true
}

func (this *xmlSimplePrinter) VisitExitElement(node XMLElement) bool {
    if node.NoChildren() {
        return true
    }

    io.WriteString(this.writer, "</")
    io.WriteString(this.writer, node.Name())
    io.WriteString(this.writer, ">")
    return true
}

func (this *xmlSimplePrinter) VisitProcInst(node XMLProcInst) bool {
    io.WriteString(this.writer, "<?")
    io.WriteString(this.writer, node.Target())
    io.WriteString(this.writer, " ")
    io.WriteString(this.writer, node.Instruction())
    io.WriteString(this.writer, "?>")
    io.WriteString(this.writer, "\n")
    return true
}

func (this *xmlSimplePrinter) VisitText(node XMLText) bool {
    if node.CDATA() {
        io.WriteString(this.writer, "<![CDATA[")
        io.WriteString(this.writer, node.Value())
        io.WriteString(this.writer, "]]")
        return true
    }

    xml.EscapeText(this.writer, []byte(node.Value()))
    return true
}

func (this *xmlSimplePrinter) VisitComment(node XMLComment) bool {
    io.WriteString(this.writer, "<!--")
    xml.EscapeText(this.writer, []byte(node.Value()))
    io.WriteString(this.writer, "-->")
    return true
}

func (this *xmlSimplePrinter) VisitDirective(node XMLDirective) bool {
    io.WriteString(this.writer, "<!")
    xml.EscapeText(this.writer, []byte(node.Value()))
    io.WriteString(this.writer, ">")
    return true
}

//------------------------------------------------------------------

type xmlHandleImpl struct {
    node XMLNode
}

func NewHandle(node XMLNode) XMLHandle {
    handle := new(xmlHandleImpl)
    handle.node = node
    return handle
}

func (this *xmlHandleImpl) Parent() XMLHandle {
    if nil == this.node {
        return this
    }

    return NewHandle(this.node.Parent())
}

func (this *xmlHandleImpl) FirstChild() XMLHandle {
    if nil == this.node {
        return this
    }

    return NewHandle(this.node.FirstChild())
}

func (this *xmlHandleImpl) LastChild() XMLHandle {
    if nil == this.node {
        return this
    }

    return NewHandle(this.node.LastChild())
}

func (this *xmlHandleImpl) PreviousSibling() XMLHandle {
    if nil == this.node {
        return this
    }

    return NewHandle(this.node.PreviousSibling())
}

func (this *xmlHandleImpl) NextSibling() XMLHandle {
    if nil == this.node {
        return this
    }

    return NewHandle(this.node.NextSibling())
}

func (this *xmlHandleImpl) FirstChildElement(name string) XMLHandle {
    if nil == this.node {
        return this
    }

    return NewHandle(this.node.FirstChildElement(name))
}

func (this *xmlHandleImpl) LastChildElement(name string) XMLHandle {
    if nil == this.node {
        return this
    }

    return NewHandle(this.node.LastChildElement(name))
}

func (this *xmlHandleImpl) PreviousSiblingElement(name string) XMLHandle {
    if nil == this.node {
        return this
    }

    return NewHandle(this.node.PreviousSiblingElement(name))
}

func (this *xmlHandleImpl) NextSiblingElement(name string) XMLHandle {
    if nil == this.node {
        return this
    }

    return NewHandle(this.node.NextSiblingElement(name))
}

func (this *xmlHandleImpl) ToNode() XMLNode {
    return this.node
}

func (this *xmlHandleImpl) ToElement() XMLElement {
    if nil == this.node {
        return nil
    }

    return this.node.ToElement()
}

func (this *xmlHandleImpl) ToText() XMLText {
    if nil == this.node {
        return nil
    }

    return this.node.ToText()
}

func (this *xmlHandleImpl) ToComment() XMLComment {
    if nil == this.node {
        return nil
    }

    return this.node.ToComment()
}

func (this *xmlHandleImpl) ToDocument() XMLDocument {
    if nil == this.node {
        return nil
    }

    return this.node.ToDocument()
}

func (this *xmlHandleImpl) ToProcInst() XMLProcInst {
    if nil == this.node {
        return nil
    }

    return this.node.ToProcInst()
}

func (this *xmlHandleImpl) ToDirective() XMLDirective {
    if nil == this.node {
        return nil
    }

    return this.node.ToDirective()
}

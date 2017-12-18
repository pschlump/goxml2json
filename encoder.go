package xml2json

import (
	"bytes"
	"io"
	"sort"
	"unicode/utf8"
)

// An Encoder writes JSON objects to an output stream.
type Encoder struct {
	w               io.Writer
	err             error
	contentPrefix   string
	attributePrefix string
	indent          bool
	indentText      string
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:               w,
		contentPrefix:   contentPrefix,
		attributePrefix: attrPrefix,
		indent:          false,
		indentText:      "",
	}
}

func (enc *Encoder) SetAttributePrefix(prefix string) *Encoder {
	enc.attributePrefix = prefix
	return enc
}

func (enc *Encoder) SetContentPrefix(prefix string) *Encoder {
	enc.contentPrefix = prefix
	return enc
}

func (enc *Encoder) SetIndent(s string) *Encoder {
	enc.indent = true
	enc.indentText = s
	return enc
}

func (enc *Encoder) EncodeWithCustomPrefixes(root *Node, contentPrefix string, attributePrefix string) error {
	enc.contentPrefix = contentPrefix
	enc.attributePrefix = attributePrefix
	return enc.Encode(root)
}

// Encode writes the JSON encoding of v to the stream
func (enc *Encoder) Encode(root *Node) error {
	if enc.err != nil {
		return enc.err
	}
	if root == nil {
		return nil
	}

	enc.err = enc.format(root, 0)

	// Terminate each value with a newline.  This makes the output look a little nicer
	// when debugging, and some kind of space is required if the encoded value was a number,
	// so that the reader knows there aren't more digits coming.
	enc.write("\n")

	return enc.err
}

// xyzzy004 - comment
func (enc *Encoder) format(curNode *Node, lvl int) error {
	var indentN = func(n int) {
		if enc.indent {
			for ii := 0; ii < n; ii++ {
				enc.write(enc.indentText)
			}
		}
	}
	if curNode.HasChildren() {
		enc.write("{")
		if enc.indent {
			enc.write("\n")
		}

		// xyzzy005 - must sort names before print?  Attributes must be in order for compare.

		// Add data as an additional attibute (if any)
		if len(curNode.Data) > 0 {
			indentN(lvl + 1)
			enc.write(`"`, enc.contentPrefix, "content", `": `, sanitiseString(curNode.Data), ", ")
			if enc.indent {
				enc.write("\n")
			}
		}

		sl := make([]string, 0, len(curNode.Children))
		for label := range curNode.Children {
			sl = append(sl, label)
		}
		// fmt.Printf("sl->%s<-\n", sl)
		if len(sl) > 1 {
			// fmt.Printf("Must sort")
			sort.Strings(sl)
		}
		// fmt.Printf("sorted: sl->%s<-\n", sl)

		com := ""
		// for label, children := range curNode.Children {
		for ii := range sl {
			label, children := sl[ii], curNode.Children[sl[ii]]
			enc.write(com)
			indentN(lvl + 1)
			enc.write(`"`, label, `": `)

			if len(children) > 1 {
				// Array
				// xyzzy005 - may need to sort?
				enc.write("[") // xyzzy006 - need to estimate if length is less than X- then one line - else - multi-line
				com1 := ""
				for _, ch := range children {
					enc.write(com1)
					enc.format(ch, lvl+2)
					com1 = ", "
				}
				enc.write("]")
			} else {
				// Map
				enc.format(children[0], lvl+1)
			}

			if enc.indent {
				com = ",\n"
			} else {
				com = ", "
			}
		}

		enc.write("\n")
		indentN(lvl)
		enc.write("}")
	} else {
		// TODO : Extract data type
		enc.write(sanitiseString(curNode.Data))
	}

	return nil
}

// xyzzy004 - comment
func (enc *Encoder) write(s ...string) {
	for _, ss := range s {
		enc.w.Write([]byte(ss))
	}
}

// https://golang.org/src/encoding/json/encode.go?s=5584:5627#L788
var hex = "0123456789abcdef"

// xyzzy008 - test
// xyzzy004 - comment
// see also: https://golang.org/src/html/escape.go
func sanitiseString(s string) string {
	var buf bytes.Buffer

	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' && b != '&' { // xyzzy009 - test for Unicode - test
				i++
				continue
			}
			if start < i {
				buf.WriteString(s[start:i])
			}
			switch b {
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			case '\t':
				buf.WriteByte('\\')
				buf.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \n and \r,
				// as well as <, > and &. The latter are escaped because they
				// can lead to security holes when user-controlled strings
				// are rendered into JSON and served to some browsers.
				buf.WriteString(`\u00`)
				buf.WriteByte(hex[b>>4])
				buf.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\u202`)
			buf.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		buf.WriteString(s[start:])
	}
	buf.WriteByte('"')
	return buf.String()
}

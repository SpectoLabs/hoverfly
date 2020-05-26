// パッケージ内の関数をテストするために、XmlPartialMatchと同じなパッケージを使う。
package matchers

import (
	"testing"

	. "github.com/beevik/etree"
	. "github.com/onsi/gomega"
)

/////////////////////////////////////////
// isLeaf

func Test_isLeaf_1(t *testing.T) {
	RegisterTestingT(t)
	doc := NewDocument()
	doc.ReadFromString("<xml>aaa</xml>")
	Expect(isLeaf(doc.Root())).To(BeTrue())
}

func Test_isLeaf_2(t *testing.T) {
	RegisterTestingT(t)
	doc := NewDocument()
	doc.ReadFromString("<xml></xml>")
	Expect(isLeaf(doc.Root())).To(BeTrue())
}

func Test_isLeaf_3(t *testing.T) {
	RegisterTestingT(t)
	doc := NewDocument()
	doc.ReadFromString("<xml />")
	Expect(isLeaf(doc.Root())).To(BeTrue())
}
/*
func Test_isLeaf_4(t *testing.T) {
	RegisterTestingT(t)
	doc := NewDocument()
	doc.ReadFromString("<xml x-match=\"allow-unknown-children\"></xml>")
	Expect(isLeaf(doc.Root())).To(BeFalse())
}

func Test_isLeaf_5(t *testing.T) {
	RegisterTestingT(t)
	doc := NewDocument()
	doc.ReadFromString("<xml x-match=\"allow-unknown-children\">aaa</xml>")
	Expect(isLeaf(doc.Root())).To(BeFalse())
}

func Test_isLeaf_6(t *testing.T) {
	RegisterTestingT(t)
	doc := NewDocument()
	doc.ReadFromString("<xml x-match-times=\"2\"></xml>")
	Expect(isLeaf(doc.Root())).To(BeTrue())
}

func Test_isLeaf_7(t *testing.T) {
	RegisterTestingT(t)
	doc := NewDocument()
	doc.ReadFromString("<xml x-match-times=\"2\">aaa</xml>")
	Expect(isLeaf(doc.Root())).To(BeTrue())
}
*/

func Test_isLeaf_8(t *testing.T) {
	RegisterTestingT(t)
	doc := NewDocument()
	doc.ReadFromString("<xml attr=\"xxx\"></xml>")
	Expect(isLeaf(doc.Root())).To(BeTrue())
}

func Test_isLeaf_9(t *testing.T) {
	RegisterTestingT(t)
	doc := NewDocument()
	doc.ReadFromString("<xml attr=\"xxx\">aaa</xml>")
	Expect(isLeaf(doc.Root())).To(BeTrue())
}

func Test_isLeaf_10(t *testing.T) {
	RegisterTestingT(t)
	doc := NewDocument()
	doc.ReadFromString("<xml><a></a></xml>")
	Expect(isLeaf(doc.Root())).To(BeFalse())
}

func Test_isLeaf_11(t *testing.T) {
	RegisterTestingT(t)
	doc := NewDocument()
	doc.ReadFromString("<xml attr=\"xxx\"><a>xxx</a></xml>")
	Expect(isLeaf(doc.Root())).To(BeFalse())
}
/*
func Test_isLeaf_12(t *testing.T) {
	RegisterTestingT(t)
	doc := NewDocument()
	doc.ReadFromString("<xml x-match-times=\"2\"><a>xxx</a></xml>")
	Expect(isLeaf(doc.Root())).To(BeFalse())
}

func Test_isLeaf_13(t *testing.T) {
	RegisterTestingT(t)
	doc := NewDocument()
	doc.ReadFromString("<xml x-match=\"allow-unknown-children\"><a>xxx</a></xml>")
	Expect(isLeaf(doc.Root())).To(BeFalse())
}
*/
/////////////////////////////////////////////

func Test_compareValue_1(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("exact", "exact")).To(BeTrue())
}

func Test_compareValue_2(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("exact", "tcaxe")).To(BeFalse())
}

func Test_compareValue_3(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("", "")).To(BeTrue())
}

func Test_compareValue_4(t *testing.T) {
	RegisterTestingT(t)
	// 全角スペース	vs 半角スペース	
	Expect(compareValue("　", " ")).To(BeFalse())
}

func Test_compareValue_5(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("   ", "  ")).To(BeFalse())
}

func Test_compareValue_6(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{ignore}}", "aaa")).To(BeTrue())
}

func Test_compareValue_7(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("   {{ignore}}", "aaa")).To(BeTrue())
}

func Test_compareValue_8(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{ignore}}     ", "aaa")).To(BeTrue())
}

func Test_compareValue_9(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{   ignore}}", "aaa")).To(BeTrue())
}

func Test_compareValue_10(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{ignore    }}", "aaa")).To(BeTrue())
}

func Test_compareValue_11(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{ {ignore}}", "aaa")).To(BeFalse())
}

func Test_compareValue_12(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{ {ignore}}", "{ {ignore}}")).To(BeTrue())
}

func Test_compareValue_13(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{ignore} }", "aaa")).To(BeFalse())
}

func Test_compareValue_14(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{ignore} }", "{{ignore} }")).To(BeTrue())
}

func Test_compareValue_15(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{ig  nore}}", "aaa")).To(BeFalse())
}

func Test_compareValue_16(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{ig  nore}}", "{{ig  nore}}")).To(BeTrue())
}

func Test_compareValue_17(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{Ignore}}", "aaa")).To(BeFalse())
}

func Test_compareValue_18(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{Ignore}}", "{{Ignore}}")).To(BeTrue())
}

func Test_compareValue_19(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{IGNORE}}", "aaa")).To(BeFalse())
}

func Test_compareValue_20(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{IGNORE}}", "{{IGNORE}}")).To(BeTrue())
}

func Test_compareValue_21(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("aaa{{ignore}}", "aaa")).To(BeFalse())
}

func Test_compareValue_22(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("aaa{{ignore}}", "aaaa")).To(BeFalse())
}

func Test_compareValue_23(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("aaa{{ignore}}", "aaa{{ignore}}")).To(BeTrue())
}

func Test_compareValue_24(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{ignore}}aaa", "aaa")).To(BeFalse())
}

func Test_compareValue_25(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{ignore}}aaa", "aaaa")).To(BeFalse())
}

func Test_compareValue_26(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{ignore}}aaa", "{{ignore}}aaa")).To(BeTrue())
}

func Test_compareValue_27(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{regex:\\d{5}}}", "12345")).To(BeTrue())
}

func Test_compareValue_28(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("   {{regex:\\d{5}}}", "12345")).To(BeTrue())
}

func Test_compareValue_29(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{regex:\\d{5}}}   ", "12345")).To(BeTrue())
}

func Test_compareValue_30(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{  regex:\\d{5}}}", "12345")).To(BeTrue())
}

func Test_compareValue_31(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{regex:  \\d{5}}}", "12346")).To(BeFalse())
}

func Test_compareValue_32(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{regex:  \\d{5}}}", "  12345")).To(BeTrue())
}

func Test_compareValue_33(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{regex:\\d{5}  }}", "12345")).To(BeFalse())
}

func Test_compareValue_34(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{regex:\\d{5}  }}", "12345  ")).To(BeTrue())
}

func Test_compareValue_35(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{regex:\\d{5}  \\d{2}}}", "12345  67")).To(BeTrue())
}

func Test_compareValue_36(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{regex:\\d{5}  \\d{2}}}", "1234567")).To(BeFalse())
}

func Test_compareValue_37(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{ {regex:\\d{5}}}", "12345")).To(BeFalse())
}

func Test_compareValue_38(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{ {regex:\\d{5}}}", "{ {regex:\\d{5}}}")).To(BeTrue())
}

func Test_compareValue_39(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{regex:\\d{5}} }", "12345")).To(BeFalse())
}

func Test_compareValue_40(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{regex:\\d{5}} }", "{{regex:\\d{5}} }")).To(BeTrue())
}

func Test_compareValue_41(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{re　gex:\\d{5}}}", "12345")).To(BeFalse())
}

func Test_compareValue_42(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{re　gex:\\d{5}}}", "{{re　gex:\\d{5}}}")).To(BeTrue())
}

func Test_compareValue_43(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{Regex:\\d{5}}}", "12345")).To(BeFalse())
}

func Test_compareValue_44(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{Regex:\\d{5}}}", "{{Regex:\\d{5}}}")).To(BeTrue())
}

func Test_compareValue_45(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{REGEX:\\d{5}}}", "12345")).To(BeFalse())
}

func Test_compareValue_46(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{REGEX:\\d{5}}}", "{{REGEX:\\d{5}}}")).To(BeTrue())
}

func Test_compareValue_47(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("aaa{{regex:\\d{5}}}", "12345")).To(BeFalse())
}

func Test_compareValue_48(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("aaa{{regex:\\d{5}}}", "aaa12345")).To(BeFalse())
}

func Test_compareValue_49(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("aaa{{regex:\\d{5}}}", "aaa{{regex:\\d{5}}}")).To(BeTrue())
}

func Test_compareValue_50(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{regex:\\d{5}}}aaa", "12345")).To(BeFalse())
}

func Test_compareValue_51(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{regex:\\d{5}}}aaa", "12345aaa")).To(BeFalse())
}

func Test_compareValue_52(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{regex:\\d{5}}}aaa", "{{regex:\\d{5}}}aaa")).To(BeTrue())
}

func Test_compareValue_53(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{{regex:\\d{5}}}}", "12345")).To(BeFalse())
}

func Test_compareValue_54(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{{regex:\\d{5}}}}", "{12345}")).To(BeFalse())
}

func Test_compareValue_55(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{{regex:\\d{5}}}}", "{{{regex:\\d{5}}}}")).To(BeTrue())
}

func Test_compareValue_56(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{regex:\\d{5}++}}", "12345")).To(BeFalse())
}

func Test_compareValue_57(t *testing.T) {
	RegisterTestingT(t)
	Expect(compareValue("{{regex:\\d{5}++}}", "{{regex:\\d{5}++}}")).To(BeFalse())
}

////////////////////////////////////////////////

// compareTree

func Test_compareTree_1(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_2(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>test</aaa>							
			<mismatch>check</mismatch>							
		</xml>								
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_3(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>
			<bbb>check</bbb>
			<aaa>test</aaa>
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_4(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<bbb>check</bbb>

		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_5(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>
			<ccc>test</ccc>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

/*

func Test_compareTree_6(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml x-match="allow-unknown-children">					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>	
			<bbb>check</bbb>							
			<aaa>test</aaa>							
			<ccc>test</ccc>							
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_7(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml x-match="allow-unknown-children">					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>	
			 							
			<aaa>test</aaa>							
			<ccc>test</ccc>							
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_8(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa x-match-times="2">test</aaa>								
			<bbb>check</bbb>
		</xml>											
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>test</aaa>							
			<aaa>test</aaa>							
			<bbb>check</bbb>							
		</xml>								
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_9(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa x-match-times="2">test</aaa>								
			<bbb>check</bbb>								
								
		</xml>													
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>test</aaa>							
			<aaa>test</aaa>							
			<aaa>test</aaa>							
			<bbb>check</bbb>							
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_10(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa x-match-times="0">test</aaa>								
			<bbb>check</bbb>								
									
		</xml>													
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
								
								
			<bbb>check</bbb>							
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_11(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa x-match-times="-1">test</aaa>								
			<bbb>check</bbb>								

		</xml>				
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
						
			<bbb>check</bbb>							
		</xml>								
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_12(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa x-match-times="*">test</aaa>								
			<bbb>check</bbb>								
									
		</xml>
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>test</aaa>							
			<aaa>test</aaa>							
			<bbb>check</bbb>							
		</xml>								
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_13(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa x-match-times="*">test</aaa>								
			<bbb>check</bbb>								
									
									
		</xml>													
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>test</aaa>							
			<aaa>test</aaa>							
			<aaa>test</aaa>							
			<bbb>check</bbb>							
		</xml>								
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_14(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa x-match-times="*">test</aaa>								
			<bbb>check</bbb>								
									
		</xml>
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
								
								
			<bbb>check</bbb>							
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_15(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa x-match-times="2">								
				<bbb x-match-times="3">x</bbb>							
			</aaa>								
		</xml>									
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>							
				<bbb>x</bbb>						
				<bbb>x</bbb>						
				<bbb>x</bbb>						
			</aaa>							
			<aaa>							
				<bbb>x</bbb>						
				<bbb>x</bbb>						
				<bbb>x</bbb>						
			</aaa>							
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_16(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa x-match-times="2">								
				<bbb x-match-times="3">x</bbb>							
			</aaa>								
		</xml>		
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>							
				<bbb>x</bbb>						
				<bbb>x</bbb>						
				<bbb>x</bbb>						
			</aaa>							
			<aaa>							
										
				<bbb>x</bbb>						
				<bbb>x</bbb>						
			</aaa>							
		</xml>								
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_17(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa x-match-times="2">								
				<bbb x-match-times="3">							
					{{regex:[A-Z]\d{3}}}						
				</bbb>							
			</aaa>								
		</xml>		
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>							
				<bbb>A123</bbb>						
				<bbb>B234</bbb>						
				<bbb>C345</bbb>						
			</aaa>							
			<aaa>							
				<bbb>D456</bbb>						
				<bbb>E567</bbb>						
				<bbb>F678</bbb>						
			</aaa>							
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_18(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa x-match-times="2">								
				<bbb x-match-times="3">							
					{{regex:[A-Z]\d{3}}}						
				</bbb>							
			</aaa>								
		</xml>														
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>							
				<bbb>A123</bbb>						
				<bbb>B234</bbb>						
				<bbb>C345</bbb>						
			</aaa>							
			<aaa>							
				<bbb>D456</bbb>						
				<bbb>E567</bbb>						
				<bbb>F6789</bbb>						
			</aaa>							
		</xml>								
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_19(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa x-match-times="2">								
				<bbb x-match-times="3">							
					{{ignore}}						
				</bbb>							
			</aaa>								
		</xml>
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>							
				<bbb>A123</bbb>						
				<bbb>B234</bbb>						
				<bbb>C345</bbb>						
			</aaa>							
			<aaa>							
				<bbb>D456</bbb>						
				<bbb>E567</bbb>						
				<bbb>F678</bbb>						
			</aaa>							
		</xml>								
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_20(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa x-match-times="2">								
				<bbb x-match-times="3">							
					{{ignore}}						
				</bbb>							
			</aaa>								
		</xml>		
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>							
				<bbb>A123</bbb>						
				<bbb>B234</bbb>						
				<bbb>C345</bbb>						
			</aaa>							
			<aaa>							
				<bbb>D456</bbb>						
				<bbb>E567</bbb>						
				<bbb>F6789</bbb>						
			</aaa>							
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_21(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa>exact</aaa>								
			<aaa>{{regex:[a-z]{5}}}</aaa>								
			<aaa>{{ignore}}</aaa>								
		</xml>														
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>abcde</aaa>							
			<aaa>exact</aaa>							
			<aaa>test</aaa>							
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

*/


func Test_compareTree_22(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa>{{regex:[a-z]{5}}}</aaa>								
			<aaa>exact</aaa>								
			<aaa>{{ignore}}</aaa>								
		</xml>	
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>abcde</aaa>							
			<aaa>exact</aaa>							
			<aaa>test</aaa>							
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_23(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa>{{ignore}}</aaa>								
			<aaa>{{regex:[a-z]{5}}}</aaa>								
			<aaa>exact</aaa>								
		</xml>	
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>abcde</aaa>							
			<aaa>exact</aaa>							
			<aaa>test</aaa>							
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_24(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>test</aaa>				
			<bbb>mismatch</bbb>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_25(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>　</aaa>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa> </aaa>								
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_26(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>   </aaa>								
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>  </aaa>		
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_27(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{ignore}}</aaa>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>xxx</aaa>					
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_28(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>  {{ignore}}</aaa>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>xxx</aaa>					
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_29(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{ignore}}  </aaa>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>xxx</aaa>					
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_30(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{ {ignore}}</aaa>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>xxx</aaa>					
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_31(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{ignore} }</aaa>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>xxx</aaa>					
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}


func Test_compareTree_32(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{ig nore}}</aaa>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>xxx</aaa>					
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}


func Test_compareTree_33(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>prefix{{ignore}}</aaa>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>prefix</aaa>					
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_34(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{ignore}}suffix</aaa>				
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>suffix</aaa>					
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_35(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{regex:\d{5}}}</aaa>	
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_36(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>   {{regex:\d{5}}}</aaa>	
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_37(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{regex:\d{5}}}   </aaa>	
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_38(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{   regex:\d{5}}}</aaa>	
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_39(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{regex:\d{5}   }}</aaa>	
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_40(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{regex:\d{5}   }}</aaa>	
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>12345   </aaa>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_41(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{regex:   \d{5}}}</aaa>	
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}


func Test_compareTree_42(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{regex:   \d{5}}}</aaa>	
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>   12345</aaa>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_43(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{regex   :\d{5}}}</aaa>	
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_44(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{regex:\d{5} \d{2}}}</aaa>	
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>1234567</aaa>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_45(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{regex:\d{5} \d{2}}}</aaa>	
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>12345 67</aaa>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeTrue())
}

func Test_compareTree_46(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>prefix{{regex:\d{5}}}</aaa>	
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>prefix12345</aaa>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_47(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{regex:\d{5}}}suffix</aaa>	
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>12345suffix</aaa>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_48(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{regex:\d{5}++}}</aaa>	
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_49(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>					
			<aaa>{{regex:\d{5}++}}</aaa>	
		</xml>					
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>					
			<aaa>{{regex:\d{5}++}}</aaa>				
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}

func Test_compareTree_50(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa>exact</aaa>								
			<aaa>{{regex:[a-z]{5}}}</aaa>								
			<aaa>{{ignore}}</aaa>								
		</xml>		
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>abcde</aaa>							
			<aaa>exact</aaa>							
										
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}


func Test_compareTree_51(t *testing.T) {
	RegisterTestingT(t)
	
	expect := NewDocument()
	expect.ReadFromString(`
		<xml>									
			<aaa>exact</aaa>								
			<aaa>{{regex:[a-z]{5}}}</aaa>								
			<aaa>ignore</aaa>								
											
		</xml>	
	`)
	actual := NewDocument()
	actual.ReadFromString(`
		<xml>								
			<aaa>abcde</aaa>							
			<aaa>exact</aaa>							
			<aaa>any</aaa>							
			<aaa>exact</aaa>							
		</xml>
	`)

	Expect(compareTree(expect.Root(), actual.Root())).To(BeFalse())
}




//////////////////////////////////////////////////////

// XMLPartialMatch

func Test_XmlPartialMatch_1(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
	xml							
		<aaa>test</aaa>						
		bbb check</bbb>			
	`
	
	actual :=`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_2(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`
	
	actual :=`
		xml								
			<aaa>test</aaa>							
			bbb check</bbb>	
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}





func Test_XmlPartialMatch_3(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_4(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`
	
	actual :=`
		<xml>								
			<aaa>test</aaa>							
			<mismatch>check</mismatch>							
		</xml>								
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_5(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`
	
	actual :=`
		<xml>
			<bbb>check</bbb>
			<aaa>test</aaa>
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_6(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`
	
	actual :=`
		<xml>								
			<bbb>check</bbb>

		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_7(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>
			<ccc>test</ccc>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

/*

func Test_XmlPartialMatch_8(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml x-match="allow-unknown-children">					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`
	
	actual :=`
		<xml>	
			<bbb>check</bbb>							
			<aaa>test</aaa>							
			<ccc>test</ccc>							
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_9(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml x-match="allow-unknown-children">					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`
	
	actual :=`
		<xml>	
			 							
			<aaa>test</aaa>							
			<ccc>test</ccc>							
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_10(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa x-match-times="2">test</aaa>								
			<bbb>check</bbb>
		</xml>											
	`
	
	actual :=`
		<xml>								
			<aaa>test</aaa>							
			<aaa>test</aaa>							
			<bbb>check</bbb>							
		</xml>								
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_11(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa x-match-times="2">test</aaa>								
			<bbb>check</bbb>								
								
		</xml>													
	`
	
	actual :=`
		<xml>								
			<aaa>test</aaa>							
			<aaa>test</aaa>							
			<aaa>test</aaa>							
			<bbb>check</bbb>							
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_12(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa x-match-times="0">test</aaa>								
			<bbb>check</bbb>								
									
		</xml>													
	`
	
	actual :=`
		<xml>								
								
								
			<bbb>check</bbb>							
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_13(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa x-match-times="-1">test</aaa>								
			<bbb>check</bbb>								

		</xml>				
	`
	
	actual :=`
		<xml>								
						
			<bbb>check</bbb>							
		</xml>								
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_14(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa x-match-times="*">test</aaa>								
			<bbb>check</bbb>								
									
		</xml>
	`
	
	actual :=`
		<xml>								
			<aaa>test</aaa>							
			<aaa>test</aaa>							
			<bbb>check</bbb>							
		</xml>								
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_15(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa x-match-times="*">test</aaa>								
			<bbb>check</bbb>								
									
									
		</xml>													
	`
	
	actual :=`
		<xml>								
			<aaa>test</aaa>							
			<aaa>test</aaa>							
			<aaa>test</aaa>							
			<bbb>check</bbb>							
		</xml>								
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_16(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa x-match-times="*">test</aaa>								
			<bbb>check</bbb>								
									
		</xml>
	`
	
	actual :=`
		<xml>								
								
								
			<bbb>check</bbb>							
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_17(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa x-match-times="2">								
				<bbb x-match-times="3">x</bbb>							
			</aaa>								
		</xml>									
	`
	
	actual :=`
		<xml>								
			<aaa>							
				<bbb>x</bbb>						
				<bbb>x</bbb>						
				<bbb>x</bbb>						
			</aaa>							
			<aaa>							
				<bbb>x</bbb>						
				<bbb>x</bbb>						
				<bbb>x</bbb>						
			</aaa>							
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_18(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa x-match-times="2">								
				<bbb x-match-times="3">x</bbb>							
			</aaa>								
		</xml>		
	`
	
	actual :=`
		<xml>								
			<aaa>							
				<bbb>x</bbb>						
				<bbb>x</bbb>						
				<bbb>x</bbb>						
			</aaa>							
			<aaa>							
										
				<bbb>x</bbb>						
				<bbb>x</bbb>						
			</aaa>							
		</xml>								
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_19(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa x-match-times="2">								
				<bbb x-match-times="3">							
					{{regex:[A-Z]\d{3}}}						
				</bbb>							
			</aaa>								
		</xml>		
	`
	
	actual :=`
		<xml>								
			<aaa>							
				<bbb>A123</bbb>						
				<bbb>B234</bbb>						
				<bbb>C345</bbb>						
			</aaa>							
			<aaa>							
				<bbb>D456</bbb>						
				<bbb>E567</bbb>						
				<bbb>F678</bbb>						
			</aaa>							
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_20(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa x-match-times="2">								
				<bbb x-match-times="3">							
					{{regex:[A-Z]\d{3}}}						
				</bbb>							
			</aaa>								
		</xml>														
	`
	
	actual :=`
		<xml>								
			<aaa>							
				<bbb>A123</bbb>						
				<bbb>B234</bbb>						
				<bbb>C345</bbb>						
			</aaa>							
			<aaa>							
				<bbb>D456</bbb>						
				<bbb>E567</bbb>						
				<bbb>F6789</bbb>						
			</aaa>							
		</xml>								
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_21(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa x-match-times="2">								
				<bbb x-match-times="3">							
					{{ignore}}						
				</bbb>							
			</aaa>								
		</xml>
	`
	
	actual :=`
		<xml>								
			<aaa>							
				<bbb>A123</bbb>						
				<bbb>B234</bbb>						
				<bbb>C345</bbb>						
			</aaa>							
			<aaa>							
				<bbb>D456</bbb>						
				<bbb>E567</bbb>						
				<bbb>F678</bbb>						
			</aaa>							
		</xml>								
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_22(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa x-match-times="2">								
				<bbb x-match-times="3">							
					{{ignore}}						
				</bbb>							
			</aaa>								
		</xml>		
	`
	
	actual :=`
		<xml>								
			<aaa>							
				<bbb>A123</bbb>						
				<bbb>B234</bbb>						
				<bbb>C345</bbb>						
			</aaa>							
			<aaa>							
				<bbb>D456</bbb>						
				<bbb>E567</bbb>						
				<bbb>F6789</bbb>						
			</aaa>							
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

*/


func Test_XmlPartialMatch_23(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa>exact</aaa>								
			<aaa>{{regex:[a-z]{5}}}</aaa>								
			<aaa>{{ignore}}</aaa>								
		</xml>																				
	`
	
	actual :=`
		<xml>								
			<aaa>abcde</aaa>							
			<aaa>exact</aaa>							
			<aaa>test</aaa>							
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_24(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa>{{regex:[a-z]{5}}}</aaa>								
			<aaa>exact</aaa>								
			<aaa>{{ignore}}</aaa>								
		</xml>
	`
	
	actual :=`
		<xml>								
			<aaa>abcde</aaa>							
			<aaa>exact</aaa>							
			<aaa>test</aaa>							
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_25(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>									
			<aaa>{{ignore}}</aaa>								
			<aaa>{{regex:[a-z]{5}}}</aaa>								
			<aaa>exact</aaa>								
		</xml>
	`
	
	actual :=`
		<xml>								
			<aaa>abcde</aaa>							
			<aaa>exact</aaa>							
			<aaa>test</aaa>							
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_26(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>test</aaa>				
			<bbb>check</bbb>				
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>test</aaa>				
			<bbb>mismatch</bbb>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_27(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>　</aaa>				
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa> </aaa>								
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_28(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>   </aaa>								
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>  </aaa>		
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_29(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{ignore}}</aaa>				
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>xxx</aaa>					
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_30(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>  {{ignore}}</aaa>				
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>xxx</aaa>					
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_31(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{ignore}}  </aaa>				
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>xxx</aaa>					
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_32(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{ {ignore}}</aaa>				
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>xxx</aaa>					
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_33(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{ignore} }</aaa>				
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>xxx</aaa>					
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}


func Test_XmlPartialMatch_34(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{ig nore}}</aaa>				
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>xxx</aaa>					
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}


func Test_XmlPartialMatch_35(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>prefix{{ignore}}</aaa>				
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>prefix</aaa>					
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_36(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{ignore}}suffix</aaa>				
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>suffix</aaa>					
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_37(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{regex:\d{5}}}</aaa>	
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_38(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>   {{regex:\d{5}}}</aaa>	
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_39(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{regex:\d{5}}}   </aaa>	
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_40(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{   regex:\d{5}}}</aaa>	
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_41(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{regex:\d{5}   }}</aaa>	
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_42(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{regex:\d{5}   }}</aaa>	
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>12345   </aaa>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_43(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{regex:   \d{5}}}</aaa>	
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}


func Test_XmlPartialMatch_44(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{regex:   \d{5}}}</aaa>	
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>   12345</aaa>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_45(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{regex   :\d{5}}}</aaa>	
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_46(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{regex:\d{5} \d{2}}}</aaa>	
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>1234567</aaa>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_47(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{regex:\d{5} \d{2}}}</aaa>	
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>12345 67</aaa>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeTrue())
}

func Test_XmlPartialMatch_48(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>prefix{{regex:\d{5}}}</aaa>	
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>prefix12345</aaa>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_49(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{regex:\d{5}}}suffix</aaa>	
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>12345suffix</aaa>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_50(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{regex:\d{5}++}}</aaa>	
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>12345</aaa>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_51(t *testing.T) {
	RegisterTestingT(t)
	
	
	expect :=`
		<xml>					
			<aaa>{{regex:\d{5}++}}</aaa>	
		</xml>					
	`
	
	actual :=`
		<xml>					
			<aaa>{{regex:\d{5}++}}</aaa>				
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}


func Test_XmlPartialMatch_52(t *testing.T) {
	RegisterTestingT(t)
	
	expect := `
		<xml>									
			<aaa>exact</aaa>								
			<aaa>{{regex:[a-z]{5}}}</aaa>								
			<aaa>{{ignore}}</aaa>								
		</xml>		
	`
	actual :=  `
		<xml>								
			<aaa>abcde</aaa>							
			<aaa>exact</aaa>							
										
		</xml>
	`

	Expect(XmlPartialMatch(expect , actual)).To(BeFalse())
}


func Test_XmlPartialMatch_53(t *testing.T) {
	RegisterTestingT(t)
	
	expect := `
		<xml>									
			<aaa>exact</aaa>								
			<aaa>{{regex:[a-z]{5}}}</aaa>								
			<aaa>ignore</aaa>								
											
		</xml>	
	`
	actual :=  `
		<xml>								
			<aaa>abcde</aaa>							
			<aaa>exact</aaa>							
			<aaa>any</aaa>							
			<aaa>exact</aaa>							
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}

func Test_XmlPartialMatch_54(t *testing.T) {
	RegisterTestingT(t)
	
	expect := 53
	actual :=  `
		<xml>								
			<aaa>abcde</aaa>							
			<aaa>exact</aaa>							
			<aaa>any</aaa>							
			<aaa>exact</aaa>							
		</xml>
	`

	Expect(XmlPartialMatch(expect, actual)).To(BeFalse())
}
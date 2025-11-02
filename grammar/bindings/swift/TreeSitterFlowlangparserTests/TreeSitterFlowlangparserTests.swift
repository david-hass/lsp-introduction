import XCTest
import SwiftTreeSitter
import TreeSitterFlowlangparser

final class TreeSitterFlowlangparserTests: XCTestCase {
    func testCanLoadGrammar() throws {
        let parser = Parser()
        let language = Language(language: tree_sitter_flowlangparser())
        XCTAssertNoThrow(try parser.setLanguage(language),
                         "Error loading flowlang parser grammar")
    }
}

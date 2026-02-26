//go:build !arm64 || purego

package operations

var mulConstImpl = mulConstGeneric
var mulAddConstImpl = mulAddConstGeneric
var mulConstNibbleImpl = mulConstNibbleGeneric
var mulAddConstNibbleImpl = mulAddConstNibbleGeneric
var mulConstTableImpl = mulConstTableGeneric
var mulAddConstTableImpl = mulAddConstTableGeneric
var xorAssignSliceImpl = xorAssignSliceGeneric

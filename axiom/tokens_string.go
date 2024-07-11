// Code generated by "stringer -type=Action -linecomment -output=tokens_string.go"; DO NOT EDIT.

package axiom

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[emptyAction-0]
	_ = x[ActionCreate-1]
	_ = x[ActionRead-2]
	_ = x[ActionUpdate-3]
	_ = x[ActionDelete-4]
}

const _Action_name = "createreadupdatedelete"

var _Action_index = [...]uint8{0, 0, 6, 10, 16, 22}

func (i Action) String() string {
	if i >= Action(len(_Action_index)-1) {
		return "Action(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Action_name[_Action_index[i]:_Action_index[i+1]]
}
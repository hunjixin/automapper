package automapper

import "errors"

var (
	ErrNotStruct       = errors.New("type must be struct")
	ErrMapperNotDefine = errors.New("mapping not define")
	ErrLengthNotMatch  = errors.New("length of slice or array not match")
)

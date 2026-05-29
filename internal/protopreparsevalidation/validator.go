package protopreparsevalidation

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type RepeatedFieldLimitError struct {
	Path  string
	Count int
	Limit int
}

func (e *RepeatedFieldLimitError) Error() string {
	return fmt.Sprintf("%s repeated field count %d exceeds limit %d", e.Path, e.Count, e.Limit)
}

// Validators use slices indexed by protobuf field number for fast dispatch. Cap
// the field number so an unexpectedly sparse schema cannot allocate a huge slice
// during Compile.
const maxAllowedFieldNumber protowire.Number = 128

// Validator validates raw protobuf bytes against a compiled set of repeated
// field limits and nested message traversals.
//
// Validator is stateful and reuses internal count buffers across Validate calls.
// It is not safe for concurrent use and should be owned by a single caller.
//
// Assumptions and limitations:
//   - Proto3 only: group wire types (StartGroup/EndGroup) are rejected outright.
//   - Field number dispatch is O(1) via sparse array indexed by field number; this
//     trades memory for speed and assumes field numbers are reasonably dense.
//

type Validator struct {
	fieldsByNumber      []*fieldValidator
	fieldCountsByNumber []int
}

type fieldValidator struct {
	// path is the protobuf field FullName used in validation errors.
	path string

	// hasLimit is true iff fieldLimits contains an explicit entry for this field.
	// Invariants: hasLimit implies the field is a non-packed repeated field, and
	// limit is non-negative. A zero limit is valid and rejects the first occurrence.
	hasLimit bool
	limit    int

	// rejectDuplicates is true for singular message fields that we traverse.
	// Canonical protobuf serialization will not emit duplicate singular fields, and
	// rejecting them avoids modeling protobuf's singular-message merge semantics.
	// Invariant: rejectDuplicates implies child != nil.
	rejectDuplicates bool

	// child validates nested message bytes. It is non-nil only for message fields
	// with count limits somewhere below them.
	child *Validator
}

// The compiler keeps explicit DFS state so it can reject recursive descriptor
// graphs instead of constructing cyclic validators.
type compileState uint8

const (
	compileStateUnvisited compileState = iota
	compileStateVisiting
	compileStateDone
)

type validatorCompiler struct {
	fieldLimits map[protoreflect.FullName]int
	reachable   map[protoreflect.FullName]protoreflect.FieldDescriptor
	states      map[protoreflect.FullName]compileState
	compiled    map[protoreflect.FullName]*Validator
}

type numberedFieldValidator struct {
	number    protowire.Number
	validator *fieldValidator
}

// Compile builds a reusable raw-message Validator from a root descriptor and a
// table of bounded repeated fields keyed by protobuf FullName. Each limit must
// be a non-negative integer.
//

// The returned Validator is stateful and must not be used concurrently.
//
// Assumptions and Limitations:
//   - No recursive schemas: the compiler rejects cyclic message descriptor graphs.
//   - No map fields: map fields are rejected at compile time because their wire
//     representation (repeated synthetic entry messages) would require special handling,
//     and we don't use them in our schemas.
//   - No repeated packable scalars (except bool): packed encoding merges multiple
//     values into a single length-delimited record, making wire-level element counting
//     inaccurate. Rather than adding complexity, we forbid these field types.
//   - No field numbers above maxAllowedFieldNumber to keep memory footprint small.
//   - See also the comment on [Validator].
func Compile(root protoreflect.MessageDescriptor, fieldLimits map[protoreflect.FullName]int) (*Validator, error) {
	compiler := newValidatorCompiler(fieldLimits)
	validator, err := compiler.compileMessage(root)
	if err != nil {
		return nil, err
	}
	if err := compiler.validateFieldLimits(root); err != nil {
		return nil, err
	}
	if validator == nil {
		return nil, fmt.Errorf("protopreparsevalidation: no validation rules reachable from root %s", root.FullName())
	}
	return validator, nil
}

// Validate scans each message slice once, enforcing repeated-field limits and
// recursively validating bounded nested message subtrees before full unmarshal.
func (validator *Validator) Validate(message []byte) error {
	if len(validator.fieldsByNumber) == 0 {
		return nil
	}

	clear(validator.fieldCountsByNumber)

	return scan(message, func(number protowire.Number, typ protowire.Type, value []byte) error {
		if int(number) >= len(validator.fieldsByNumber) {
			return nil
		}

		fieldValidator := validator.fieldsByNumber[number]
		if fieldValidator == nil {
			return nil
		}

		if fieldValidator.hasLimit || fieldValidator.rejectDuplicates {
			validator.fieldCountsByNumber[number]++
			fieldCount := validator.fieldCountsByNumber[number]
			if fieldValidator.rejectDuplicates && fieldCount > 1 {
				// This is technically valid in protobuf, but the canonical protobuf serialization that we use
				// will never emit duplicate singular fields.
				// This makes our validator easier to implement, because we do not need to worry about merging
				// multiple occurrences of a singular message field into one.
				return fmt.Errorf("%s singular message field appears multiple times", fieldValidator.path)
			}
			if fieldValidator.hasLimit && fieldCount > fieldValidator.limit {
				return &RepeatedFieldLimitError{
					fieldValidator.path,
					fieldCount,
					fieldValidator.limit,
				}
			}
		}

		if fieldValidator.child == nil {
			return nil
		}

		// fieldValidator.child != nil implies field.Kind() == MessageKind (the compiler
		// only recurses into message-typed fields). In proto3, message fields, whether
		// singular or repeated, always use BytesType on the wire.
		if typ != protowire.BytesType {
			return fmt.Errorf("%s has wire type %d, want %d", fieldValidator.path, typ, protowire.BytesType)
		}

		return fieldValidator.child.Validate(value)
	})
}

func newValidatorCompiler(fieldLimits map[protoreflect.FullName]int) *validatorCompiler {
	return &validatorCompiler{
		fieldLimits,
		map[protoreflect.FullName]protoreflect.FieldDescriptor{},
		map[protoreflect.FullName]compileState{},
		map[protoreflect.FullName]*Validator{},
	}
}

// After the DFS compile finishes, reachable contains exactly the fields that can
// appear under the chosen root descriptor. Any policy entry not seen during the
// walk is therefore invalid for this validator.
func (compiler *validatorCompiler) validateFieldLimits(root protoreflect.MessageDescriptor) error {
	for fieldName, limit := range compiler.fieldLimits {
		if limit < 0 {
			return fmt.Errorf("protopreparsevalidation: field %s has a negative limit %d", fieldName, limit)
		}

		field, reachable := compiler.reachable[fieldName]
		if !reachable {
			return fmt.Errorf("protopreparsevalidation: field %s is not reachable from root %s or does not exist", fieldName, root.FullName())
		}

		if field.IsMap() {
			return fmt.Errorf("protopreparsevalidation: field %s is a map field; bounded map fields are not currently supported", fieldName)
		}
		if !field.IsList() {
			return fmt.Errorf("protopreparsevalidation: field %s is not repeated and cannot be bounded", fieldName)
		}
		if field.IsPacked() {
			return fmt.Errorf("protopreparsevalidation: field %s uses packed encoding; wire-level counting would be inaccurate", fieldName)
		}
	}

	for fieldName, field := range compiler.reachable {
		if field.IsMap() {
			return fmt.Errorf("protopreparsevalidation: map field %s is reachable from root %s; map fields are not supported", fieldName, root.FullName())
		}
		if !field.IsList() {
			continue
		}
		if field.IsPacked() {
			if field.Kind() == protoreflect.BoolKind {

				continue
			}
			return fmt.Errorf("protopreparsevalidation: repeated packed field %s (kind %s) is reachable from root %s; only repeated bool is allowed among packable types", fieldName, field.Kind(), root.FullName())
		}
		if _, covered := compiler.fieldLimits[fieldName]; !covered {
			return fmt.Errorf("protopreparsevalidation: repeated field %s is reachable from root %s but has no limit", fieldName, root.FullName())
		}
	}

	return nil
}

// compileMessage recursively performs a DFS from the root descriptor and returns nil
// for message subtrees with no validation work so they do not need to be traversed
// at validation time.
func (compiler *validatorCompiler) compileMessage(message protoreflect.MessageDescriptor) (*Validator, error) {
	messageName := message.FullName()
	switch compiler.states[messageName] {
	case compileStateUnvisited:
		// continue with compilation, it's the first time we encounter this message
		break
	case compileStateVisiting:
		return nil, fmt.Errorf("protopreparsevalidation: recursive message graph not supported: %s", messageName)
	case compileStateDone:
		return compiler.compiled[messageName], nil
	}
	compiler.states[messageName] = compileStateVisiting

	validator, err := compiler.compileFields(message.Fields())
	if err != nil {
		return nil, err
	}

	compiler.states[messageName] = compileStateDone
	if validator != nil {
		compiler.compiled[messageName] = validator
	}
	return validator, nil
}

// compileFields folds field descriptors into a validator. Fields that are
// reachable but have no validation work below them are skipped.
func (compiler *validatorCompiler) compileFields(fields protoreflect.FieldDescriptors) (*Validator, error) {
	validator := &Validator{}

	// We only need to allocate the state for counting fields if at least
	// one of the fields' validators requires it.
	hasFieldCounts := false

	for index := 0; index < fields.Len(); index++ {
		compiledField, err := compiler.compileField(fields.Get(index))
		if err != nil {
			return nil, err
		}
		if compiledField == nil {
			continue
		}

		if compiledField.validator.hasLimit || compiledField.validator.rejectDuplicates {
			hasFieldCounts = true
		}

		if err := validator.setFieldValidator(compiledField.number, compiledField.validator); err != nil {
			return nil, err
		}
	}

	if len(validator.fieldsByNumber) == 0 {
		return nil, nil
	}

	if hasFieldCounts {
		validator.fieldCountsByNumber = make([]int, len(validator.fieldsByNumber))
	}
	return validator, nil
}

// compileField returns nil for reachable fields that do not need a runtime rule.
func (compiler *validatorCompiler) compileField(field protoreflect.FieldDescriptor) (*numberedFieldValidator, error) {
	compiler.reachable[field.FullName()] = field

	// Map fields are not supported.
	if field.IsMap() {
		return nil, fmt.Errorf("protopreparsevalidation: map field %s is not supported", field.FullName())
	}

	limit, hasLimit := compiler.fieldLimits[field.FullName()]
	var child *Validator

	// Recurse (only through non-map) message fields. Other field kinds can still be
	// bounded, but they never contribute additional nested validation work.
	if field.Kind() == protoreflect.MessageKind {
		var err error
		child, err = compiler.compileMessage(field.Message())
		if err != nil {
			return nil, err
		}
	}
	if !hasLimit && child == nil {
		return nil, nil
	}

	return &numberedFieldValidator{
		protowire.Number(field.Number()),
		&fieldValidator{
			string(field.FullName()),
			hasLimit,
			limit,
			child != nil && !field.IsList(),
			child,
		},
	}, nil
}

// Field numbers are used as sparse array indices so validation dispatch stays
// O(1) during the raw scan without a per-message map lookup.
func (validator *Validator) setFieldValidator(number protowire.Number, field *fieldValidator) error {
	if number > maxAllowedFieldNumber {
		return fmt.Errorf("protopreparsevalidation: field %s has number %d, maximum allowed field number is %d", field.path, number, maxAllowedFieldNumber)
	}
	if int(number) >= len(validator.fieldsByNumber) {
		expanded := make([]*fieldValidator, number+1)
		copy(expanded, validator.fieldsByNumber)
		validator.fieldsByNumber = expanded
	}
	validator.fieldsByNumber[number] = field
	return nil
}

// scan iterates over the raw wire fields of a protobuf message.
func scan(message []byte, visit func(protowire.Number, protowire.Type, []byte) error) error {
	for len(message) > 0 {
		number, typ, n := protowire.ConsumeTag(message)
		if n < 0 {
			return protowire.ParseError(n)
		}

		fieldValue := message[n:]
		switch typ {
		case protowire.BytesType:
			value, valueLen := protowire.ConsumeBytes(fieldValue)
			if valueLen < 0 {
				return protowire.ParseError(valueLen)
			}
			if err := visit(number, typ, value); err != nil {
				return err
			}
			message = fieldValue[valueLen:]
		case protowire.VarintType, protowire.Fixed32Type, protowire.Fixed64Type:
			valueLen := protowire.ConsumeFieldValue(number, typ, fieldValue)
			if valueLen < 0 {
				return protowire.ParseError(valueLen)
			}
			if err := visit(number, typ, nil); err != nil {
				return err
			}
			message = fieldValue[valueLen:]
		case protowire.StartGroupType, protowire.EndGroupType:
			// Groups are a deprecated proto2 feature not used by our proto3 schemas.
			// We reject them outright to be conservative.
			return fmt.Errorf("unsupported group wire type %d at field %d", typ, number)
		default:
			// Fail closed
			return fmt.Errorf("unknown wire type %d at field %d", typ, number)
		}
	}
	return nil
}

package conversion

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

// Field and Value are the common source representation used by reflection and
// generated field accessors. They describe source data rather than native blocks.
type Field struct {
	Name  string
	Tag   Tag
	Value Value
}

// Value is a source value read by reflection or generated accessors.
// Kind is "scalar", "object", "list", or "dynamic".
type Value struct {
	Kind   string
	Text   string
	Empty  bool
	Nil    bool
	Fields []Field
	Items  []Value
	Err    error
	shape  *Shape
}

// Generated is implemented by generated accessors in the caller's package.
type Generated interface {
	ConvertagoFields(platform string) ([]Field, error)
}

// Shape describes source types without depending on reflection or go/types.
type Shape struct {
	Kind   string
	Fields []ShapeField
	Elem   *Shape
	Err    error
}

// ShapeField describes one tagged field of a source shape.
type ShapeField struct {
	Name  string
	Tag   Tag
	Value *Shape
	Err   error
}

type fieldPlan struct {
	name  string
	index int
	tag   Tag
	value *valuePlan
	err   error
}

type valuePlan struct {
	typeOf    reflect.Type
	shape     *Shape
	elem      *valuePlan
	fields    []fieldPlan
	marshaler bool
	stringer  bool
}

type planKey struct {
	typeOf   reflect.Type
	platform string
}

type cachedPlan struct {
	plan *valuePlan
	err  error
}

var sourcePlans sync.Map
var rawPlans sync.Map
var textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
var stringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

// Fields reads the tagged fields of a source struct, preferring generated accessors.
func Fields(input any, platform string) ([]Field, error) {
	profile, err := Lookup(platform)
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(input)
	base := v
	rootPointers := make(map[any]bool)
	for base.IsValid() && base.Kind() == reflect.Pointer {
		if base.IsNil() {
			return nil, Error(platform, "$", "invalid_source", "source must be a non-nil struct")
		}
		pointer := base.Interface()
		if rootPointers[pointer] {
			return nil, Error(platform, "$", "invalid_source", "cyclic source value")
		}
		rootPointers[pointer] = true
		base = base.Elem()
	}
	if !base.IsValid() || base.Kind() != reflect.Struct {
		return nil, Error(platform, "$", "invalid_source", "source must be a struct or a pointer to a struct")
	}
	if generated, ok := input.(Generated); ok {
		return generated.ConvertagoFields(platform)
	}
	if generated, ok := base.Interface().(Generated); ok {
		return generated.ConvertagoFields(platform)
	}
	entry := loadPlan(base.Type(), profile)
	if entry.err != nil {
		return nil, entry.err
	}
	value := readValue(entry.plan, base, profile, &State{})
	if value.Err != nil {
		return nil, Error(platform, "$", "invalid_source", value.Err.Error())
	}
	return value.Fields, nil
}

func loadPlan(typeOf reflect.Type, profile Profile) cachedPlan {
	key := planKey{typeOf, profile.Platform}
	if found, ok := sourcePlans.Load(key); ok {
		return found.(cachedPlan)
	}
	plan := compileObjectPlan(typeOf, profile)
	entry := cachedPlan{plan: plan, err: ValidateShape(profile, plan.shape)}
	actual, _ := sourcePlans.LoadOrStore(key, entry)
	return actual.(cachedPlan)
}

// A root struct is always a source object, even if it has a text method.
// Nested fields may still explicitly use that method as their representation.
func compileObjectPlan(typeOf reflect.Type, profile Profile) *valuePlan {
	plan := &valuePlan{typeOf: typeOf, shape: &Shape{Kind: "object"}}
	compileFields(plan, profile, make(map[reflect.Type]*valuePlan))
	return plan
}

func compileFields(plan *valuePlan, profile Profile, memo map[reflect.Type]*valuePlan) {
	typeOf := plan.typeOf
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		raw := field.Tag.Get(profile.Platform)
		if raw == "" || raw == "-" {
			continue
		}
		tag, err := Parse(profile, raw)
		if field.PkgPath != "" {
			err = fmt.Errorf("tagged field must be exported")
		}
		child := compilePlan(field.Type, profile, memo)
		plan.fields = append(plan.fields, fieldPlan{field.Name, i, tag, child, err})
		plan.shape.Fields = append(plan.shape.Fields, ShapeField{field.Name, tag, child.shape, err})
	}
}

func rawPlan(typeOf reflect.Type, profile Profile) *valuePlan {
	key := planKey{typeOf, profile.Platform}
	if found, ok := rawPlans.Load(key); ok {
		return found.(*valuePlan)
	}
	plan := compilePlan(typeOf, profile, make(map[reflect.Type]*valuePlan))
	actual, _ := rawPlans.LoadOrStore(key, plan)
	return actual.(*valuePlan)
}

func compilePlan(typeOf reflect.Type, profile Profile, memo map[reflect.Type]*valuePlan) *valuePlan {
	if found := memo[typeOf]; found != nil {
		return found
	}
	plan := &valuePlan{typeOf: typeOf, shape: &Shape{}}
	memo[typeOf] = plan
	if typeOf.Kind() == reflect.Interface {
		plan.shape.Kind = "dynamic"
		return plan
	}
	if typeOf.Implements(textMarshalerType) || typeOf.Implements(stringerType) {
		plan.shape.Kind = "scalar"
		plan.marshaler = typeOf.Implements(textMarshalerType)
		plan.stringer = !plan.marshaler
		return plan
	}
	switch typeOf.Kind() {
	case reflect.Pointer:
		plan.elem = compilePlan(typeOf.Elem(), profile, memo)
		plan.shape.Kind = "pointer"
		plan.shape.Elem = plan.elem.shape
	case reflect.Struct:
		plan.shape.Kind = "object"
		compileFields(plan, profile, memo)
	case reflect.Slice, reflect.Array:
		plan.shape.Kind = "list"
		plan.elem = compilePlan(typeOf.Elem(), profile, memo)
		plan.shape.Elem = plan.elem.shape
	case reflect.Interface:
		plan.shape.Kind = "dynamic"
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		plan.shape.Kind = "scalar"
	default:
		plan.shape.Err = fmt.Errorf("unsupported source type %s", typeOf)
	}
	return plan
}

func readValue(plan *valuePlan, v reflect.Value, profile Profile, state *State) Value {
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		if v.IsNil() {
			return Value{Kind: baseKind(plan.shape), Nil: true, Empty: true}
		}
	}
	if v.Kind() == reflect.Pointer || (v.Kind() == reflect.Slice && v.Len() != 0) {
		pointer, kind := v.Interface(), "pointer"
		if v.Kind() == reflect.Slice {
			pointer, kind = v.Index(0).Addr().Interface(), "slice"
		}
		if !state.Enter(pointer, kind) {
			return Value{Kind: baseKind(plan.shape), Err: fmt.Errorf("cyclic source value")}
		}
		defer state.Leave(pointer, kind)
	}
	if plan.marshaler {
		value := Marshaled(v.Interface().(encoding.TextMarshaler))
		value.Empty = value.Text == "" && v.Kind() != reflect.Pointer
		return value
	}
	if plan.stringer {
		text := v.Interface().(fmt.Stringer).String()
		return Value{Kind: "scalar", Text: text, Empty: text == "" && v.Kind() != reflect.Pointer}
	}
	switch v.Kind() {
	case reflect.Pointer:
		value := readValue(plan.elem, v.Elem(), profile, state)
		value.Empty = false
		return value
	case reflect.Interface:
		child := rawPlan(v.Elem().Type(), profile)
		value := readValue(child, v.Elem(), profile, state)
		value.shape = child.shape
		value.Empty = false
		return value
	case reflect.Struct:
		fields := make([]Field, 0, len(plan.fields))
		for _, field := range plan.fields {
			value := Value{Err: field.err}
			if field.err == nil {
				value = readValue(field.value, v.Field(field.index), profile, state)
			}
			fields = append(fields, Field{field.name, field.tag, value})
		}
		return Object(fields)
	case reflect.Slice, reflect.Array:
		items := make([]Value, v.Len())
		for i := range items {
			items[i] = readValue(plan.elem, v.Index(i), profile, state)
		}
		return Value{Kind: "list", Empty: v.Len() == 0, Items: items}
	case reflect.String:
		return Value{Kind: "scalar", Text: v.String(), Empty: v.Len() == 0}
	case reflect.Bool:
		return Value{Kind: "scalar", Text: strconv.FormatBool(v.Bool()), Empty: !v.Bool()}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Value{Kind: "scalar", Text: strconv.FormatInt(v.Int(), 10), Empty: v.Int() == 0}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return Value{Kind: "scalar", Text: strconv.FormatUint(v.Uint(), 10), Empty: v.Uint() == 0}
	case reflect.Float32, reflect.Float64:
		return Value{Kind: "scalar", Text: strconv.FormatFloat(v.Float(), 'g', -1, v.Type().Bits()), Empty: v.Float() == 0}
	default:
		return Value{Err: plan.shape.Err}
	}
}

func baseKind(shape *Shape) string {
	seen := make(map[*Shape]bool)
	for shape.Kind == "pointer" && !seen[shape] {
		seen[shape] = true
		shape = shape.Elem
	}
	return shape.Kind
}

// Object builds a tagged object value, which is empty when all its fields are empty.
func Object(fields []Field) Value {
	value := Value{Kind: "object", Empty: true, Fields: fields}
	for _, field := range fields {
		if !field.Value.Empty {
			value.Empty = false
		}
	}
	return value
}

// Marshaled reads a scalar value from MarshalText, retaining its error.
func Marshaled(input encoding.TextMarshaler) Value {
	text, err := input.MarshalText()
	return Value{Kind: "scalar", Text: string(text), Empty: len(text) == 0, Err: err}
}

// Dynamic is the reflection fallback for interface-valued generated fields.
func Dynamic(input any, platform string) Value {
	return (&State{}).Dynamic(input, platform)
}

// Dynamic reads an interface-valued field by reflection within this traversal.
func (state *State) Dynamic(input any, platform string) Value {
	profile, err := Lookup(platform)
	if err != nil {
		return Value{Err: err}
	}
	v := reflect.ValueOf(input)
	if !v.IsValid() {
		return Value{Kind: "dynamic", Nil: true, Empty: true}
	}
	plan := rawPlan(v.Type(), profile)
	value := readValue(plan, v, profile, state)
	value.shape = plan.shape
	if !value.Nil {
		value.Empty = false
	}
	return value
}

package generate

import (
	"fmt"
	"go/types"
	"reflect"

	"github.com/JSYoo5B/convertago/internal/conversion"
)

func structTag(raw, platform string) string { return reflect.StructTag(raw).Get(platform) }

func objectShape(t types.Type, profile conversion.Profile) *conversion.Shape {
	shape := &conversion.Shape{Kind: "object"}
	shapeFields(shape, t.Underlying().(*types.Struct), profile, make(map[types.Type]*conversion.Shape))
	return shape
}

func shapeFields(shape *conversion.Shape, underlying *types.Struct, profile conversion.Profile, memo map[types.Type]*conversion.Shape) {
	for i := 0; i < underlying.NumFields(); i++ {
		raw := structTag(underlying.Tag(i), profile.Platform)
		if raw == "" || raw == "-" {
			continue
		}
		field := underlying.Field(i)
		tag, err := conversion.Parse(profile, raw)
		if !field.Exported() {
			err = fmt.Errorf("tagged field must be exported")
		}
		shape.Fields = append(shape.Fields, conversion.ShapeField{
			Name: field.Name(), Tag: tag, Value: shapeOf(field.Type(), profile, memo), Err: err,
		})
	}
}

func shapeOf(t types.Type, profile conversion.Profile, memo map[types.Type]*conversion.Shape) *conversion.Shape {
	if shape := memo[t]; shape != nil {
		return shape
	}
	shape := &conversion.Shape{}
	memo[t] = shape
	if _, ok := t.Underlying().(*types.Interface); ok {
		shape.Kind = "dynamic"
		return shape
	}
	if types.Implements(t, marshalerInterface) || types.Implements(t, stringerInterface) {
		shape.Kind = "scalar"
		return shape
	}
	switch underlying := t.Underlying().(type) {
	case *types.Pointer:
		shape.Kind, shape.Elem = "pointer", shapeOf(underlying.Elem(), profile, memo)
	case *types.Struct:
		shape.Kind = "object"
		shapeFields(shape, underlying, profile, memo)
	case *types.Slice:
		shape.Kind, shape.Elem = "list", shapeOf(underlying.Elem(), profile, memo)
	case *types.Array:
		shape.Kind, shape.Elem = "list", shapeOf(underlying.Elem(), profile, memo)
	case *types.Basic:
		if underlying.Info()&(types.IsString|types.IsBoolean|types.IsInteger|types.IsFloat) != 0 {
			shape.Kind = "scalar"
		} else {
			shape.Err = fmt.Errorf("unsupported source type %s", t)
		}
	default:
		shape.Err = fmt.Errorf("unsupported source type %s", t)
	}
	return shape
}

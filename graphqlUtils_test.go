package resource

import (
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertEqualNoWhitespace(t *testing.T, expected, actual string) {
	whitespace := regexp.MustCompile(`\s+`)
	assert.Equal(t, strings.Trim(whitespace.ReplaceAllString(expected, " "), " "), strings.Trim(whitespace.ReplaceAllString(actual, " "), " "))
}

func TestSingleFieldInsertIntoQueryType(t *testing.T) {
	combined, err := combineSchemas(`
	type Query {
	}`, `
	type Query {
		hello: String
	}`, []Resolver{
		{
			FieldName: "hello",
			TypeName:  "Query",
		},
	}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `type Query { hello: String }`, combined)
}

func TestSingleFieldNotInResolverSet(t *testing.T) {
	combined, err := combineSchemas(`
	type Query {
	}`, `
	type Query {
		hello: String
	}`, []Resolver{}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `type Query {}`, combined)
}

func TestSingleFieldInUnknownType(t *testing.T) {
	combined, err := combineSchemas(`
	type Query {
	}`, `
	type Foo {
		hello: String
	}
	type Query {
		foo: Foo
	}`, []Resolver{
		{
			FieldName: "hello",
			TypeName:  "Foo",
		},
	}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `type Query {}`, combined)
}

func TestSingleFieldWithType(t *testing.T) {
	combined, err := combineSchemas(`
	type Query {
	}`, `
	type Foo {
		hello: String
	}
	type Query {
		foo: Foo
	}`, []Resolver{
		{
			FieldName: "foo",
			TypeName:  "Query",
		},
	}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `type Query { foo: Foo } type Foo { hello: String }`, combined)
}

func TestSingleFieldWithTypeAndNestedType(t *testing.T) {
	combined, err := combineSchemas(`
	type Query {
	}`, `
	type Bar {
		hello: String
	}
	type Foo {
		hello: Bar
	}
	type Query {
		foo: Foo
	}`, []Resolver{
		{
			FieldName: "foo",
			TypeName:  "Query",
		},
	}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `type Query { foo: Foo } type Foo { hello: Bar } type Bar { hello: String }`, combined)
}

func TestFieldWithInterfaceReturn(t *testing.T) {
	combined, err := combineSchemas(`
	type Query {
	}`, `
	type Bar {
		hello: String
	}
	type Baz {
		hello: String
	}
	interface IFoo {
		hello: String
		baz: Baz
	}
	type Foo implements IFoo {
		hello: String
		baz: Baz
		bar: Bar
	}
	type Query {
		foo: IFoo
	}`, []Resolver{
		{
			FieldName: "foo",
			TypeName:  "Query",
		},
	}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `
	type Query {
		foo: IFoo
	}
    interface IFoo {
		hello: String
		baz: Baz
	}
	type Baz {
		hello: String
	}
	type Foo implements IFoo {
		hello: String
		baz: Baz
		bar: Bar
	}
	type Bar {
		hello: String
	}
	`, combined)
}

func TestFieldWithArguments(t *testing.T) {
	combined, err := combineSchemas(`
	type Query {
	}`, `
	type Query {
		foo(bar: String): String
	}`, []Resolver{
		{
			FieldName: "foo",
			TypeName:  "Query",
		},
	}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `type Query { foo(bar: String): String }`, combined)
}

func TestFieldWithInputType(t *testing.T) {
	combined, err := combineSchemas(`
	type Query {
	}`, `
	input BarInput {
		bar: String
	}
	input FooInput {
		bar: BarInput
	}
	type Query {
		foo(inp: FooInput): String
	}`, []Resolver{
		{
			FieldName: "foo",
			TypeName:  "Query",
		},
	}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `type Query { foo(inp: FooInput): String } input FooInput { bar: BarInput } input BarInput { bar: String }`, combined)
}

func TestFieldWithEnum(t *testing.T) {
	combined, err := combineSchemas(`
	type Query {
	}`, `
	enum BarEnum {
		FOO
		BAR
	}
	type Foo {
		bar: BarEnum
	}
	type Query {
		foo: Foo
	}`, []Resolver{
		{
			FieldName: "foo",
			TypeName:  "Query",
		},
	}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `type Query { foo: Foo } type Foo { bar: BarEnum } enum BarEnum { FOO BAR }`, combined)
}

func TestFieldWithEnumInInput(t *testing.T) {
	combined, err := combineSchemas(`
	type Query {
	}`, `
	enum BarEnum {
		FOO
		BAR
	}
	input FooInput {
		bar: BarEnum
	}
	type Query {
		foo(inp: FooInput): String
	}`, []Resolver{
		{
			FieldName: "foo",
			TypeName:  "Query",
		},
	}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `type Query { foo(inp: FooInput): String } input FooInput { bar: BarEnum } enum BarEnum { FOO BAR }`, combined)
}

func TestFieldWithUnion(t *testing.T) {
	combined, err := combineSchemas(`
	type Query {
	}`, `
	type Foo {
		bar: String
	}
	type Baz {
		baz: String
	}
	union BarUnion = Foo | Baz
	type Query {
		foo: BarUnion
	}`, []Resolver{
		{
			FieldName: "foo",
			TypeName:  "Query",
		},
	}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `type Query { foo: BarUnion } union BarUnion = Foo | Baz type Foo { bar: String } type Baz { baz: String }`, combined)
}

func TestExistingType(t *testing.T) {
	combined, err := combineSchemas(`
	type Foo {
		previous: String
	}
	type Query {
		foo: Foo
	}`, `
	type Foo {
		hello: String
	}
	type Query {
		foo: Foo
	}`, []Resolver{
		{
			FieldName: "foo",
			TypeName:  "Query",
		},
	}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `type Foo { hello: String } type Query { foo: Foo }`, combined)
}

func TestMultipleModificationsOnExistingType(t *testing.T) {
	combined, err := combineSchemas(`
	type Foo {
		previous: String
	}
	type Query {
		foo: Foo
		bar: String
	}`, `
	type Foo {
		hello: String
	}
	type Query {
		foo: Foo
		baz: String
	}`, []Resolver{
		{
			FieldName: "foo",
			TypeName:  "Query",
		},
	}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `type Foo { hello: String } type Query { foo: Foo bar: String }`, combined)
}

func TestMultipleUpdatedResolvers(t *testing.T) {
	combined, err := combineSchemas(`
	type Foo {
		previous: String
	}
	type Query {
		foo: Foo
		bar: String
	}`, `
	type Foo {
		hello: String
	}
	type Bar {
		hello: String
	}
	type Query {
		foo: Foo
		bar: Bar
		baz: String
	}`, []Resolver{
		{
			FieldName: "foo",
			TypeName:  "Query",
		},
		{
			FieldName: "bar",
			TypeName:  "Query",
		},
	}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `type Foo { hello: String } type Query { foo: Foo bar: Bar } type Bar { hello: String }`, combined)
}

func TestChangeFromScalarToDifferentScalar(t *testing.T) {
	combined, err := combineSchemas(`
	scalar Foo
	type Query {
		foo: Foo
	}`, `
	scalar Bar
	type Query {
		foo: Bar
	}`, []Resolver{
		{
			FieldName: "foo",
			TypeName:  "Query",
		},
	}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `scalar Foo type Query { foo: Bar } scalar Bar`, combined)
}

func TestChangeFromScalarToObject(t *testing.T) {
	combined, err := combineSchemas(`
	scalar Foo
	type Query {
		foo: Foo
	}`, `
	type Foo {
		hello: String
	}
	type Query {
		foo: Foo
	}`, []Resolver{
		{
			FieldName: "foo",
			TypeName:  "Query",
		},
	}, log.Default())

	assert.Nil(t, err)
	assertEqualNoWhitespace(t, `type Foo { hello: String } type Query { foo: Foo }`, combined)
}

// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/google/cel-go/common/operators"
	expr "github.com/google/cel-spec/proto/v1/syntax"
)

type TestExpr struct {
	Expr       *expr.Expr
	SourceInfo *expr.SourceInfo
}

func (t *TestExpr) Info(location string) *expr.SourceInfo {
	info := proto.Clone(t.SourceInfo).(*expr.SourceInfo)
	info.Location = location
	return info
}

var (
	// program with no instructions.
	Empty = &TestExpr{
		&expr.Expr{},

		&expr.SourceInfo{
			LineOffsets: []int32{},
			Positions:   map[int64]int32{}}}

	// [1, 1u, 1.0].exists(x, type(x) == uint)
	Exists = &TestExpr{
		ExprComprehension(1,
			"x",
			ExprList(8,
				ExprLiteral(2, int64(0)),
				ExprLiteral(3, int64(1)),
				ExprLiteral(4, int64(2)),
				ExprLiteral(5, int64(3)),
				ExprLiteral(6, int64(4)),
				ExprLiteral(7, uint64(5))),
			"_accu_",
			ExprLiteral(9, false),
			ExprCall(10,
				operators.LogicalNot,
				ExprIdent(11, "_accu_")),
			ExprCall(12,
				operators.Equals,
				ExprCall(13,
					"type",
					ExprIdent(14, "x")),
				ExprIdent(15, "uint")),
			ExprIdent(16, "_accu_")),

		&expr.SourceInfo{
			LineOffsets: []int32{0},
			Positions: map[int64]int32{
				0:  12,
				1:  0,
				2:  1,
				3:  4,
				4:  8,
				5:  0,
				6:  18,
				7:  18,
				8:  18,
				9:  18,
				10: 18,
				11: 20,
				12: 20,
				13: 28,
				14: 28,
				15: 28,
				16: 28}}}

	// elems.exists(x, type(x) == uint)
	ExistsWithInput = &TestExpr{
		ExprComprehension(1,
			"x",
			ExprIdent(2, "elems"),
			"_accu_",
			ExprLiteral(3, false),
			ExprCall(4,
				operators.LogicalNot,
				ExprIdent(5, "_accu_")),
			ExprCall(6,
				operators.Equals,
				ExprCall(7,
					"type",
					ExprIdent(8, "x")),
				ExprIdent(9, "uint")),
			ExprIdent(10, "_accu_")),

		&expr.SourceInfo{
			LineOffsets: []int32{0},
			Positions: map[int64]int32{
				0:  12,
				1:  0,
				2:  1,
				3:  4,
				4:  8,
				5:  0,
				6:  18,
				7:  18,
				8:  18,
				9:  18,
				10: 18}}}

	// {"hello": "world".size(),
	//  "dur": duration.Duration{10},
	//  "ts": timestamp.Timestamp{1000},
	//  "null": null,
	//  "bytes": b"bytes-string"}
	DynMap = &TestExpr{
		ExprMap(17,
			ExprEntry(2,
				ExprLiteral(1, "hello"),
				ExprMemberCall(3,
					"size",
					ExprLiteral(4, "world"))),
			ExprEntry(12,
				ExprLiteral(11, "null"),
				ExprLiteral(13, structpb.NullValue_NULL_VALUE)),
			ExprEntry(15,
				ExprLiteral(14, "bytes"),
				ExprLiteral(16, []byte("bytes-string")))),

		&expr.SourceInfo{
			LineOffsets: []int32{},
			Positions:   map[int64]int32{}}}

	// a && TestProto{c: true}.c
	LogicalAnd = &TestExpr{
		ExprCall(2, operators.LogicalAnd,
			ExprIdent(1, "a"),
			ExprSelect(6,
				ExprType(5, "TestProto",
					ExprField(4, "c", ExprLiteral(5, true))),
				"c")),
		&expr.SourceInfo{
			LineOffsets: []int32{},
			Positions:   map[int64]int32{}}}

	// a ? b < 1.0 : c == ["hello"]
	Conditional = &TestExpr{
		ExprCall(9, operators.Conditional,
			ExprIdent(1, "a"),
			ExprCall(3,
				operators.Less,
				ExprIdent(2, "b"),
				ExprLiteral(4, 1.0)),
			ExprCall(6,
				operators.Equals,
				ExprIdent(5, "c"),
				ExprList(8, ExprLiteral(7, "hello")))),
		&expr.SourceInfo{
			LineOffsets: []int32{},
			Positions:   map[int64]int32{}}}

	// a.b.c
	Select = &TestExpr{
		ExprSelect(3,
			ExprSelect(2,
				ExprIdent(1, "a"),
				"b"),
			"c"),
		&expr.SourceInfo{
			LineOffsets: []int32{},
			Positions:   map[int64]int32{}}}

	// a == 42
	Equality = &TestExpr{
		ExprCall(2,
			operators.Equals,
			ExprIdent(1, "a"),
			ExprLiteral(3, int64(42))),
		&expr.SourceInfo{
			LineOffsets: []int32{},
			Positions:   map[int64]int32{}}}

	// a == 42
	TypeEquality = &TestExpr{
		ExprCall(4,
			operators.Equals,
			ExprCall(1, "type",
				ExprIdent(2, "a")),
			ExprIdent(3, "uint")),
		&expr.SourceInfo{
			LineOffsets: []int32{},
			Positions:   map[int64]int32{}}}
)

func ExprIdent(id int64, name string) *expr.Expr {
	return &expr.Expr{id, &expr.Expr_IdentExpr{
		&expr.Expr_Ident{name}}}
}

func ExprSelect(id int64, operand *expr.Expr, field string) *expr.Expr {
	return &expr.Expr{id,
		&expr.Expr_SelectExpr{
			&expr.Expr_Select{operand, field, false}}}
}

func ExprLiteral(id int64, value interface{}) *expr.Expr {
	var literal *expr.Literal
	switch value.(type) {
	case bool:
		literal = &expr.Literal{&expr.Literal_BoolValue{value.(bool)}}
	case int64:
		literal = &expr.Literal{&expr.Literal_Int64Value{
			value.(int64)}}
	case uint64:
		literal = &expr.Literal{&expr.Literal_Uint64Value{
			value.(uint64)}}
	case float64:
		literal = &expr.Literal{&expr.Literal_DoubleValue{
			value.(float64)}}
	case string:
		literal = &expr.Literal{&expr.Literal_StringValue{
			value.(string)}}
	case structpb.NullValue:
		literal = &expr.Literal{&expr.Literal_NullValue{
			value.(structpb.NullValue)}}
	case []byte:
		literal = &expr.Literal{&expr.Literal_BytesValue{
			value.([]byte)}}
	default:
		panic("literal type not implemented")
	}
	return &expr.Expr{id, &expr.Expr_LiteralExpr{LiteralExpr: literal}}
}

func ExprCall(id int64, function string, args ...*expr.Expr) *expr.Expr {
	return &expr.Expr{id,
		&expr.Expr_CallExpr{
			&expr.Expr_Call{nil, function, args}}}
}

func ExprMemberCall(id int64, function string, target *expr.Expr, args ...*expr.Expr) *expr.Expr {
	return &expr.Expr{id,
		&expr.Expr_CallExpr{
			&expr.Expr_Call{target, function, args}}}
}

func ExprList(id int64, elements ...*expr.Expr) *expr.Expr {
	return &expr.Expr{id,
		&expr.Expr_ListExpr{
			&expr.Expr_CreateList{elements}}}
}

func ExprMap(id int64, entries ...*expr.Expr_CreateStruct_Entry) *expr.Expr {
	return &expr.Expr{id, &expr.Expr_StructExpr{
		&expr.Expr_CreateStruct{Entries: entries}}}
}

func ExprType(id int64, messageName string,
	entries ...*expr.Expr_CreateStruct_Entry) *expr.Expr {
	return &expr.Expr{id, &expr.Expr_StructExpr{
		&expr.Expr_CreateStruct{
			messageName, entries}}}
}

func ExprEntry(id int64, key *expr.Expr,
	value *expr.Expr) *expr.Expr_CreateStruct_Entry {
	return &expr.Expr_CreateStruct_Entry{id,
		&expr.Expr_CreateStruct_Entry_MapKey{key},
		value}
}

func ExprField(id int64, field string,
	value *expr.Expr) *expr.Expr_CreateStruct_Entry {
	return &expr.Expr_CreateStruct_Entry{id,
		&expr.Expr_CreateStruct_Entry_FieldKey{field},
		value}
}

func ExprComprehension(id int64,
	iterVar string, iterRange *expr.Expr,
	accuVar string, accuInit *expr.Expr,
	loopCondition *expr.Expr, loopStep *expr.Expr,
	resultExpr *expr.Expr) *expr.Expr {
	return &expr.Expr{id, &expr.Expr_ComprehensionExpr{
		&expr.Expr_Comprehension{
			iterVar,
			iterRange,
			accuVar,
			accuInit,
			loopCondition,
			loopStep,
			resultExpr}}}
}

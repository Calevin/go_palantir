// Code generated by ent, DO NOT EDIT.

package token

import (
	"entgo.io/ent/dialect/sql"
	"github.com/Calevin/go_palantir/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldID, id))
}

// File applies equality check predicate on the "file" field. It's identical to FileEQ.
func File(v string) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldFile, v))
}

// Line applies equality check predicate on the "line" field. It's identical to LineEQ.
func Line(v int) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldLine, v))
}

// Order applies equality check predicate on the "order" field. It's identical to OrderEQ.
func Order(v int) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldOrder, v))
}

// Token applies equality check predicate on the "token" field. It's identical to TokenEQ.
func Token(v string) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldToken, v))
}

// FileEQ applies the EQ predicate on the "file" field.
func FileEQ(v string) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldFile, v))
}

// FileNEQ applies the NEQ predicate on the "file" field.
func FileNEQ(v string) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldFile, v))
}

// FileIn applies the In predicate on the "file" field.
func FileIn(vs ...string) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldFile, vs...))
}

// FileNotIn applies the NotIn predicate on the "file" field.
func FileNotIn(vs ...string) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldFile, vs...))
}

// FileGT applies the GT predicate on the "file" field.
func FileGT(v string) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldFile, v))
}

// FileGTE applies the GTE predicate on the "file" field.
func FileGTE(v string) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldFile, v))
}

// FileLT applies the LT predicate on the "file" field.
func FileLT(v string) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldFile, v))
}

// FileLTE applies the LTE predicate on the "file" field.
func FileLTE(v string) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldFile, v))
}

// FileContains applies the Contains predicate on the "file" field.
func FileContains(v string) predicate.Token {
	return predicate.Token(sql.FieldContains(FieldFile, v))
}

// FileHasPrefix applies the HasPrefix predicate on the "file" field.
func FileHasPrefix(v string) predicate.Token {
	return predicate.Token(sql.FieldHasPrefix(FieldFile, v))
}

// FileHasSuffix applies the HasSuffix predicate on the "file" field.
func FileHasSuffix(v string) predicate.Token {
	return predicate.Token(sql.FieldHasSuffix(FieldFile, v))
}

// FileEqualFold applies the EqualFold predicate on the "file" field.
func FileEqualFold(v string) predicate.Token {
	return predicate.Token(sql.FieldEqualFold(FieldFile, v))
}

// FileContainsFold applies the ContainsFold predicate on the "file" field.
func FileContainsFold(v string) predicate.Token {
	return predicate.Token(sql.FieldContainsFold(FieldFile, v))
}

// LineEQ applies the EQ predicate on the "line" field.
func LineEQ(v int) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldLine, v))
}

// LineNEQ applies the NEQ predicate on the "line" field.
func LineNEQ(v int) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldLine, v))
}

// LineIn applies the In predicate on the "line" field.
func LineIn(vs ...int) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldLine, vs...))
}

// LineNotIn applies the NotIn predicate on the "line" field.
func LineNotIn(vs ...int) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldLine, vs...))
}

// LineGT applies the GT predicate on the "line" field.
func LineGT(v int) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldLine, v))
}

// LineGTE applies the GTE predicate on the "line" field.
func LineGTE(v int) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldLine, v))
}

// LineLT applies the LT predicate on the "line" field.
func LineLT(v int) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldLine, v))
}

// LineLTE applies the LTE predicate on the "line" field.
func LineLTE(v int) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldLine, v))
}

// OrderEQ applies the EQ predicate on the "order" field.
func OrderEQ(v int) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldOrder, v))
}

// OrderNEQ applies the NEQ predicate on the "order" field.
func OrderNEQ(v int) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldOrder, v))
}

// OrderIn applies the In predicate on the "order" field.
func OrderIn(vs ...int) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldOrder, vs...))
}

// OrderNotIn applies the NotIn predicate on the "order" field.
func OrderNotIn(vs ...int) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldOrder, vs...))
}

// OrderGT applies the GT predicate on the "order" field.
func OrderGT(v int) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldOrder, v))
}

// OrderGTE applies the GTE predicate on the "order" field.
func OrderGTE(v int) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldOrder, v))
}

// OrderLT applies the LT predicate on the "order" field.
func OrderLT(v int) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldOrder, v))
}

// OrderLTE applies the LTE predicate on the "order" field.
func OrderLTE(v int) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldOrder, v))
}

// TokenEQ applies the EQ predicate on the "token" field.
func TokenEQ(v string) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldToken, v))
}

// TokenNEQ applies the NEQ predicate on the "token" field.
func TokenNEQ(v string) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldToken, v))
}

// TokenIn applies the In predicate on the "token" field.
func TokenIn(vs ...string) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldToken, vs...))
}

// TokenNotIn applies the NotIn predicate on the "token" field.
func TokenNotIn(vs ...string) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldToken, vs...))
}

// TokenGT applies the GT predicate on the "token" field.
func TokenGT(v string) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldToken, v))
}

// TokenGTE applies the GTE predicate on the "token" field.
func TokenGTE(v string) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldToken, v))
}

// TokenLT applies the LT predicate on the "token" field.
func TokenLT(v string) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldToken, v))
}

// TokenLTE applies the LTE predicate on the "token" field.
func TokenLTE(v string) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldToken, v))
}

// TokenContains applies the Contains predicate on the "token" field.
func TokenContains(v string) predicate.Token {
	return predicate.Token(sql.FieldContains(FieldToken, v))
}

// TokenHasPrefix applies the HasPrefix predicate on the "token" field.
func TokenHasPrefix(v string) predicate.Token {
	return predicate.Token(sql.FieldHasPrefix(FieldToken, v))
}

// TokenHasSuffix applies the HasSuffix predicate on the "token" field.
func TokenHasSuffix(v string) predicate.Token {
	return predicate.Token(sql.FieldHasSuffix(FieldToken, v))
}

// TokenEqualFold applies the EqualFold predicate on the "token" field.
func TokenEqualFold(v string) predicate.Token {
	return predicate.Token(sql.FieldEqualFold(FieldToken, v))
}

// TokenContainsFold applies the ContainsFold predicate on the "token" field.
func TokenContainsFold(v string) predicate.Token {
	return predicate.Token(sql.FieldContainsFold(FieldToken, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Token) predicate.Token {
	return predicate.Token(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Token) predicate.Token {
	return predicate.Token(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Token) predicate.Token {
	return predicate.Token(sql.NotPredicates(p))
}

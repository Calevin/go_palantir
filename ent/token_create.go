// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Calevin/go_palantir/ent/file"
	"github.com/Calevin/go_palantir/ent/token"
)

// TokenCreate is the builder for creating a Token entity.
type TokenCreate struct {
	config
	mutation *TokenMutation
	hooks    []Hook
}

// SetLine sets the "line" field.
func (tc *TokenCreate) SetLine(i int) *TokenCreate {
	tc.mutation.SetLine(i)
	return tc
}

// SetOrder sets the "order" field.
func (tc *TokenCreate) SetOrder(i int) *TokenCreate {
	tc.mutation.SetOrder(i)
	return tc
}

// SetToken sets the "token" field.
func (tc *TokenCreate) SetToken(s string) *TokenCreate {
	tc.mutation.SetToken(s)
	return tc
}

// SetFileID sets the "file" edge to the File entity by ID.
func (tc *TokenCreate) SetFileID(id int) *TokenCreate {
	tc.mutation.SetFileID(id)
	return tc
}

// SetNillableFileID sets the "file" edge to the File entity by ID if the given value is not nil.
func (tc *TokenCreate) SetNillableFileID(id *int) *TokenCreate {
	if id != nil {
		tc = tc.SetFileID(*id)
	}
	return tc
}

// SetFile sets the "file" edge to the File entity.
func (tc *TokenCreate) SetFile(f *File) *TokenCreate {
	return tc.SetFileID(f.ID)
}

// Mutation returns the TokenMutation object of the builder.
func (tc *TokenCreate) Mutation() *TokenMutation {
	return tc.mutation
}

// Save creates the Token in the database.
func (tc *TokenCreate) Save(ctx context.Context) (*Token, error) {
	return withHooks(ctx, tc.sqlSave, tc.mutation, tc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TokenCreate) SaveX(ctx context.Context) *Token {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tc *TokenCreate) Exec(ctx context.Context) error {
	_, err := tc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tc *TokenCreate) ExecX(ctx context.Context) {
	if err := tc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tc *TokenCreate) check() error {
	if _, ok := tc.mutation.Line(); !ok {
		return &ValidationError{Name: "line", err: errors.New(`ent: missing required field "Token.line"`)}
	}
	if _, ok := tc.mutation.Order(); !ok {
		return &ValidationError{Name: "order", err: errors.New(`ent: missing required field "Token.order"`)}
	}
	if _, ok := tc.mutation.Token(); !ok {
		return &ValidationError{Name: "token", err: errors.New(`ent: missing required field "Token.token"`)}
	}
	return nil
}

func (tc *TokenCreate) sqlSave(ctx context.Context) (*Token, error) {
	if err := tc.check(); err != nil {
		return nil, err
	}
	_node, _spec := tc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	tc.mutation.id = &_node.ID
	tc.mutation.done = true
	return _node, nil
}

func (tc *TokenCreate) createSpec() (*Token, *sqlgraph.CreateSpec) {
	var (
		_node = &Token{config: tc.config}
		_spec = sqlgraph.NewCreateSpec(token.Table, sqlgraph.NewFieldSpec(token.FieldID, field.TypeInt))
	)
	if value, ok := tc.mutation.Line(); ok {
		_spec.SetField(token.FieldLine, field.TypeInt, value)
		_node.Line = value
	}
	if value, ok := tc.mutation.Order(); ok {
		_spec.SetField(token.FieldOrder, field.TypeInt, value)
		_node.Order = value
	}
	if value, ok := tc.mutation.Token(); ok {
		_spec.SetField(token.FieldToken, field.TypeString, value)
		_node.Token = value
	}
	if nodes := tc.mutation.FileIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   token.FileTable,
			Columns: []string{token.FileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(file.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.file_tokens = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// TokenCreateBulk is the builder for creating many Token entities in bulk.
type TokenCreateBulk struct {
	config
	err      error
	builders []*TokenCreate
}

// Save creates the Token entities in the database.
func (tcb *TokenCreateBulk) Save(ctx context.Context) ([]*Token, error) {
	if tcb.err != nil {
		return nil, tcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(tcb.builders))
	nodes := make([]*Token, len(tcb.builders))
	mutators := make([]Mutator, len(tcb.builders))
	for i := range tcb.builders {
		func(i int, root context.Context) {
			builder := tcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TokenMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, tcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, tcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tcb *TokenCreateBulk) SaveX(ctx context.Context) []*Token {
	v, err := tcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcb *TokenCreateBulk) Exec(ctx context.Context) error {
	_, err := tcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcb *TokenCreateBulk) ExecX(ctx context.Context) {
	if err := tcb.Exec(ctx); err != nil {
		panic(err)
	}
}

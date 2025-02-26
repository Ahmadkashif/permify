package commands

import (
	"context"
	"errors"
	`fmt`

	`github.com/Permify/permify/internal/repositories`
	"github.com/Permify/permify/pkg/helper"
	"github.com/Permify/permify/pkg/logger"
	base "github.com/Permify/permify/pkg/pb/base/v1"
)

// LookupSchemaCommand -
type LookupSchemaCommand struct {
	// repositories
	schemaReader repositories.SchemaReader
	// logger
	logger logger.Interface
}

// NewLookupSchemaCommand -
func NewLookupSchemaCommand(schemaReader repositories.SchemaReader, l logger.Interface) *LookupSchemaCommand {
	return &LookupSchemaCommand{
		schemaReader: schemaReader,
		logger:       l,
	}
}

// Execute -
func (command *LookupSchemaCommand) Execute(ctx context.Context, request *base.PermissionLookupSchemaRequest) (*base.PermissionLookupSchemaResponse, error) {

	var err error

	if request.GetSchemaVersion() == "" {
		var ver string
		ver, err = command.schemaReader.HeadVersion(ctx)
		if err != nil {
			return nil, err
		}
		request.SchemaVersion = ver
	}

	var en *base.EntityDefinition
	en, _, err = command.schemaReader.ReadSchemaDefinition(ctx, request.GetEntityType(), request.GetSchemaVersion())
	if err != nil {
		return nil, err
	}

	response := &base.PermissionLookupSchemaResponse{
		ActionNames: []string{},
	}
	for _, action := range en.GetActions() {
		var can bool
		can, err = command.l(ctx, request, action.Child)
		if err != nil {
			return nil, err
		}
		if can {
			response.ActionNames = append(response.ActionNames, action.Name)
		}
	}
	return response, nil
}

// SchemaLookupDecision -
type SchemaLookupDecision struct {
	Can bool
	Err error
}

// sendSchemaLookupDecision -
func sendSchemaLookupDecision(can bool, err error) SchemaLookupDecision {
	return SchemaLookupDecision{
		Can: can,
		Err: err,
	}
}

// SchemaLookupFunction -
type SchemaLookupFunction func(ctx context.Context, lookupChan chan<- SchemaLookupDecision)

// SchemaLookupCombiner .
type SchemaLookupCombiner func(ctx context.Context, functions []SchemaLookupFunction) SchemaLookupDecision

// c -
func (command *LookupSchemaCommand) l(ctx context.Context, request *base.PermissionLookupSchemaRequest, child *base.Child) (bool, error) {
	var fn SchemaLookupFunction
	switch child.Type.(type) {
	case *base.Child_Rewrite:
		fn = command.lookupRewrite(ctx, request, child.GetRewrite())
	case *base.Child_Leaf:
		fn = command.lookupLeaf(ctx, request, child.GetLeaf())
	}

	if fn == nil {
		return false, errors.New(base.ErrorCode_ERROR_CODE_UNDEFINED_CHILD_KIND.String())
	}

	result := schemaLookupUnion(ctx, []SchemaLookupFunction{fn})
	return result.Can, result.Err
}

// lookupRewrite -
func (command *LookupSchemaCommand) lookupRewrite(ctx context.Context, request *base.PermissionLookupSchemaRequest, rewrite *base.Rewrite) SchemaLookupFunction {
	switch rewrite.GetRewriteOperation() {
	case *base.Rewrite_OPERATION_UNION.Enum():
		return command.setChild(ctx, request, rewrite.GetChildren(), schemaLookupUnion)
	case *base.Rewrite_OPERATION_INTERSECTION.Enum():
		return command.setChild(ctx, request, rewrite.GetChildren(), schemaLookupIntersection)
	default:
		return schemaLookupFail(errors.New(base.ErrorCode_ERROR_CODE_UNDEFINED_CHILD_TYPE.String()))
	}
}

// checkLeaf -
func (command *LookupSchemaCommand) lookupLeaf(ctx context.Context, request *base.PermissionLookupSchemaRequest, leaf *base.Leaf) SchemaLookupFunction {
	switch leaf.GetType().(type) {
	case *base.Leaf_TupleToUserSet:
		return command.lookup(ctx, fmt.Sprintf("%s.%s", leaf.GetTupleToUserSet().GetTupleSet().GetRelation(), leaf.GetTupleToUserSet().GetComputed().GetRelation()), request, leaf.GetExclusion())
	case *base.Leaf_ComputedUserSet:
		return command.lookup(ctx, leaf.GetComputedUserSet().GetRelation(), request, leaf.GetExclusion())
	default:
		return schemaLookupFail(errors.New(base.ErrorCode_ERROR_CODE_UNDEFINED_CHILD_TYPE.String()))
	}
}

// set -
func (command *LookupSchemaCommand) setChild(ctx context.Context, request *base.PermissionLookupSchemaRequest, children []*base.Child, combiner SchemaLookupCombiner) SchemaLookupFunction {
	var functions []SchemaLookupFunction
	for _, child := range children {
		switch child.GetType().(type) {
		case *base.Child_Rewrite:
			functions = append(functions, command.lookupRewrite(ctx, request, child.GetRewrite()))
		case *base.Child_Leaf:
			functions = append(functions, command.lookupLeaf(ctx, request, child.GetLeaf()))
		default:
			return schemaLookupFail(errors.New(base.ErrorCode_ERROR_CODE_UNDEFINED_CHILD_KIND.String()))
		}
	}

	return func(ctx context.Context, resultChan chan<- SchemaLookupDecision) {
		resultChan <- combiner(ctx, functions)
	}
}

// check -
func (command *LookupSchemaCommand) lookup(ctx context.Context, relation string, request *base.PermissionLookupSchemaRequest, exclusion bool) SchemaLookupFunction {
	return func(ctx context.Context, lookupChan chan<- SchemaLookupDecision) {
		var err error
		if exclusion {
			lookupChan <- sendSchemaLookupDecision(!helper.InArray(relation, request.GetRelationNames()), err)
			return
		}
		lookupChan <- sendSchemaLookupDecision(helper.InArray(relation, request.GetRelationNames()), err)
		return
	}
}

// union -
func schemaLookupUnion(ctx context.Context, functions []SchemaLookupFunction) SchemaLookupDecision {
	if len(functions) == 0 {
		return sendSchemaLookupDecision(true, nil)
	}

	lookupChan := make(chan SchemaLookupDecision, len(functions))
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, fn := range functions {
		go fn(childCtx, lookupChan)
	}

	for i := 0; i < len(functions); i++ {
		select {
		case result := <-lookupChan:
			if result.Err == nil && result.Can {
				return sendSchemaLookupDecision(true, nil)
			}
			if result.Err != nil {
				return sendSchemaLookupDecision(false, result.Err)
			}
		case <-ctx.Done():
			return sendSchemaLookupDecision(false, errors.New(base.ErrorCode_ERROR_CODE_CANCELLED.String()))
		}
	}

	return sendSchemaLookupDecision(false, nil)
}

// intersection -
func schemaLookupIntersection(ctx context.Context, functions []SchemaLookupFunction) SchemaLookupDecision {
	if len(functions) == 0 {
		return sendSchemaLookupDecision(true, nil)
	}

	lookupChan := make(chan SchemaLookupDecision, len(functions))
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, fn := range functions {
		go fn(childCtx, lookupChan)
	}

	for i := 0; i < len(functions); i++ {
		select {
		case result := <-lookupChan:
			if result.Err == nil && !result.Can {
				return sendSchemaLookupDecision(false, nil)
			}
			if result.Err != nil {
				return sendSchemaLookupDecision(false, result.Err)
			}
		case <-ctx.Done():
			return sendSchemaLookupDecision(false, errors.New(base.ErrorCode_ERROR_CODE_CANCELLED.String()))
		}
	}

	return sendSchemaLookupDecision(true, nil)
}

// schemaLookupFail -
func schemaLookupFail(err error) SchemaLookupFunction {
	return func(ctx context.Context, decisionChan chan<- SchemaLookupDecision) {
		decisionChan <- sendSchemaLookupDecision(false, err)
	}
}

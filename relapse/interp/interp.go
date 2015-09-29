//  Copyright 2015 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package interp

import (
	"fmt"
	"github.com/katydid/katydid/expr/compose"
	"github.com/katydid/katydid/relapse/ast"
	"github.com/katydid/katydid/serialize"
	"io"
	"log"
)

//This is a naive implementation and it does not handle left recursion
func Interpret(g *relapse.Grammar, tree serialize.Parser) bool {
	refs := relapse.NewRefsLookup(g)
	res := refs["main"]
	res = refs["main"]
	err := tree.Next()
	for err == nil {
		name, nameErr := tree.String()
		if nameErr != nil {
			panic(nameErr)
		}
		log.Printf("Interpret = %s given input %s", res, name)
		res = sderiv(refs, res, tree)
		err = tree.Next()
	}
	if err != io.EOF {
		panic(err)
	}
	log.Printf("Interpret Final = %s", res)
	return Nullable(refs, res)
}

//TODO improve nullable for left recursion using fix points
// https://github.com/kennknowles/go-yid/blob/master/src/yid/nullable.go
//This is a naive implementation and it does not handle left recursion
func Nullable(refs relapse.RefLookup, p *relapse.Pattern) bool {
	typ := p.GetValue()
	switch v := typ.(type) {
	case *relapse.Empty:
		return true
	case *relapse.EmptySet:
		return false
	case *relapse.TreeNode:
		return false
	case *relapse.LeafNode:
		return false
	case *relapse.Concat:
		return Nullable(refs, v.GetLeftPattern()) && Nullable(refs, v.GetRightPattern())
	case *relapse.Or:
		return Nullable(refs, v.GetLeftPattern()) || Nullable(refs, v.GetRightPattern())
	case *relapse.And:
		return Nullable(refs, v.GetLeftPattern()) && Nullable(refs, v.GetRightPattern())
	case *relapse.ZeroOrMore:
		return true
	case *relapse.Reference:
		return Nullable(refs, refs[v.GetName()])
	case *relapse.Not:
		return !(Nullable(refs, v.GetPattern()))
	case *relapse.ZAny:
		return true
	case *relapse.WithSomeOr:
		return Nullable(refs, v.GetLeftPattern()) || Nullable(refs, v.GetRightPattern())
	case *relapse.WithSomeAnd:
		return Nullable(refs, v.GetLeftPattern()) && Nullable(refs, v.GetRightPattern())
	case *relapse.WithSomeTreeNode:
		return Nullable(refs, v.GetPattern())
	}
	panic(fmt.Sprintf("unknown pattern typ %T", typ))
}

func evalName(nameExpr *relapse.NameExpr, name string) bool {
	typ := nameExpr.GetValue()
	switch v := typ.(type) {
	case *relapse.Name:
		return name == v.GetName()
	case *relapse.AnyName:
		return true
	case *relapse.AnyNameExcept:
		return !evalName(v.GetExcept(), name)
	case *relapse.NameChoice:
		return evalName(v.GetLeft(), name) || evalName(v.GetRight(), name)
	}
	panic(fmt.Sprintf("unknown nameExpr typ %T", typ))
}

func derivTreeNode(refs relapse.RefLookup, p *relapse.TreeNode, tree serialize.Parser) *relapse.Pattern {
	name, nameErr := tree.String()
	if nameErr != nil {
		panic(nameErr)
	}
	matched := evalName(p.GetName(), name)
	if !matched {
		return relapse.NewEmptySet()
	}
	tree.Down()
	res := p.GetPattern()
	err := tree.Next()
	for err == nil {
		if !tree.IsLeaf() {
			name1, nameErr1 := tree.String()
			if nameErr1 != nil {
				panic(nameErr1)
			}
			log.Printf("derivTreeNode = %s given input %s", res, name1)
		}
		res = sderiv(refs, res, tree)
		err = tree.Next()
	}
	if err != io.EOF {
		panic(err)
	}
	log.Printf("derivTreeNode Final = %s", res)
	tree.Up()
	if !Nullable(refs, res) {
		return relapse.NewEmptySet()
	}
	return relapse.NewEmpty()
}

func sderiv(refs relapse.RefLookup, p *relapse.Pattern, tree serialize.Parser) *relapse.Pattern {
	d := deriv(refs, p, tree)
	log.Printf("sderiv %s -> %s", p, d)
	return Simplify(refs, d)
}

func deriv(refs relapse.RefLookup, p *relapse.Pattern, tree serialize.Parser) *relapse.Pattern {
	typ := p.GetValue()
	switch v := typ.(type) {
	case *relapse.Empty:
		return relapse.NewEmptySet()
	case *relapse.EmptySet:
		return relapse.NewEmptySet()
	case *relapse.TreeNode:
		return derivTreeNode(refs, v, tree)
	case *relapse.LeafNode:
		f, err := compose.NewBool(v.GetExpr())
		if err != nil {
			panic(err)
		}
		if !tree.IsLeaf() {
			return relapse.NewEmptySet()
		}
		res, err := f.Eval(tree)
		if err != nil {
			return relapse.NewEmptySet()
		}
		if res {
			return relapse.NewEmpty()
		}
		return relapse.NewEmptySet()
	case *relapse.Concat:
		leftDeriv := relapse.NewConcat(sderiv(refs, v.GetLeftPattern(), tree), v.GetRightPattern())
		if Nullable(refs, v.GetLeftPattern()) {
			return relapse.NewOr(
				leftDeriv,
				sderiv(refs, v.GetRightPattern(), tree.Copy()),
			)
		} else {
			return leftDeriv
		}
	case *relapse.Or:
		return relapse.NewOr(
			sderiv(refs, v.GetLeftPattern(), tree),
			sderiv(refs, v.GetRightPattern(), tree.Copy()),
		)
	case *relapse.And:
		return relapse.NewAnd(
			sderiv(refs, v.GetLeftPattern(), tree),
			sderiv(refs, v.GetRightPattern(), tree.Copy()),
		)
	case *relapse.ZeroOrMore:
		return relapse.NewConcat(sderiv(refs, v.Pattern, tree), relapse.NewZeroOrMore(v.Pattern))
	case *relapse.Reference:
		return sderiv(refs, refs[v.GetName()], tree)
	case *relapse.Not:
		return relapse.NewNot(sderiv(refs, v.GetPattern(), tree))
	case *relapse.ZAny:
		return deriv(refs, relapse.NewNot(relapse.NewEmptySet()), tree)
	case *relapse.WithSomeOr:
		left := relapse.NewConcat(relapse.NewZAny(), v.GetLeftPattern(), relapse.NewZAny())
		right := relapse.NewConcat(relapse.NewZAny(), v.GetRightPattern(), relapse.NewZAny())
		newor := relapse.NewOr(left, right)
		return deriv(refs, newor, tree)
	case *relapse.WithSomeAnd:
		left := relapse.NewConcat(relapse.NewZAny(), v.GetLeftPattern(), relapse.NewZAny())
		right := relapse.NewConcat(relapse.NewZAny(), v.GetRightPattern(), relapse.NewZAny())
		newand := relapse.NewAnd(left, right)
		return deriv(refs, newand, tree)
	case *relapse.WithSomeTreeNode:
		newp := relapse.NewConcat(relapse.NewZAny(), v.GetPattern(), relapse.NewZAny())
		return deriv(refs, newp, tree)
	}
	panic(fmt.Sprintf("unknown typ %T", typ))
}

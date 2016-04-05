//  Copyright 2016 Walter Schulze
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
	"github.com/katydid/katydid/parser"
	"github.com/katydid/katydid/relapse/ast"
	"github.com/katydid/katydid/relapse/compose"
	nameexpr "github.com/katydid/katydid/relapse/name"
	"io"
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

//This is a naive implementation and it does not handle left recursion
func Interpret(g *ast.Grammar, parser parser.Interface) bool {
	refs := ast.NewRefLookup(g)
	finals := deriv(refs, []*ast.Pattern{refs["main"]}, parser)
	return Nullable(refs, finals[0])
}

func escapable(patterns []*ast.Pattern) bool {
	for _, pattern := range patterns {
		if pattern.ZAny != nil {
			continue
		}
		if pattern.Not != nil {
			if pattern.GetNot().GetPattern().ZAny != nil {
				continue
			}
		}
		return true
	}
	return false
}

func deriv(refs map[string]*ast.Pattern, patterns []*ast.Pattern, tree parser.Interface) []*ast.Pattern {
	var resPatterns []*ast.Pattern = patterns
	for {
		if !escapable(resPatterns) {
			return resPatterns
		}
		if err := tree.Next(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		childPatterns := derivCalls(refs, resPatterns, tree)
		if tree.IsLeaf() {
			//do nothing
		} else {
			tree.Down()
			zchild, zi := zip(childPatterns)
			zchild = deriv(refs, zchild, tree)
			childPatterns = unzip(zchild, zi)
			tree.Up()
		}
		resPatterns = derivReturns(refs, resPatterns, childPatterns)
		resPatterns = simps(refs, resPatterns)
	}
	return resPatterns
}

func simps(refs map[string]*ast.Pattern, patterns []*ast.Pattern) []*ast.Pattern {
	for i := range patterns {
		patterns[i] = Simplify(refs, patterns[i])
	}
	return patterns
}

func zip(patterns []*ast.Pattern) ([]*ast.Pattern, []int) {
	zipped := ast.Set(patterns)
	ast.Sort(zipped)
	indexes := make([]int, len(patterns))
	for i, pattern := range patterns {
		indexes[i] = ast.Index(zipped, pattern)
	}
	return zipped, indexes
}

func unzip(patterns []*ast.Pattern, indexes []int) []*ast.Pattern {
	res := make([]*ast.Pattern, len(indexes))
	for i, index := range indexes {
		res[i] = patterns[index]
	}
	return res
}

func derivCalls(refs map[string]*ast.Pattern, patterns []*ast.Pattern, label parser.Value) []*ast.Pattern {
	res := []*ast.Pattern{}
	for _, pattern := range patterns {
		ps := derivCall(refs, pattern, label)
		ps = simps(refs, ps)
		res = append(res, ps...)
	}
	return res
}

func derivCall(refs map[string]*ast.Pattern, p *ast.Pattern, label parser.Value) []*ast.Pattern {
	typ := p.GetValue()
	switch v := typ.(type) {
	case *ast.Empty:
		return []*ast.Pattern{}
	case *ast.ZAny:
		return []*ast.Pattern{}
	case *ast.TreeNode:
		b := nameexpr.NameToFunc(v.GetName())
		f, err := compose.NewBoolFunc(b)
		if err != nil {
			panic(err)
		}
		eval, err := f.Eval(label)
		if err != nil {
			panic(err)
		}
		if eval {
			return []*ast.Pattern{v.GetPattern()}
		}
		return []*ast.Pattern{ast.NewNot(ast.NewZAny())}
	case *ast.LeafNode:
		b, err := compose.NewBool(v.GetExpr())
		if err != nil {
			panic(err)
		}
		f, err := compose.NewBoolFunc(b)
		if err != nil {
			panic(err)
		}
		eval, err := f.Eval(label)
		if err != nil {
			panic(err)
		}
		if eval {
			return []*ast.Pattern{ast.NewEmpty()}
		}
		return []*ast.Pattern{ast.NewNot(ast.NewZAny())}
	case *ast.Concat:
		l := derivCall(refs, v.GetLeftPattern(), label)
		if !Nullable(refs, v.GetLeftPattern()) {
			return l
		}
		r := derivCall(refs, v.GetRightPattern(), label)
		return append(l, r...)
	case *ast.Or:
		return derivCall2(refs, v.GetLeftPattern(), v.GetRightPattern(), label)
	case *ast.And:
		return derivCall2(refs, v.GetLeftPattern(), v.GetRightPattern(), label)
	case *ast.Interleave:
		return derivCall2(refs, v.GetLeftPattern(), v.GetRightPattern(), label)
	case *ast.ZeroOrMore:
		return derivCall(refs, v.GetPattern(), label)
	case *ast.Reference:
		return derivCall(refs, refs[v.GetName()], label)
	case *ast.Not:
		return derivCall(refs, v.GetPattern(), label)
	case *ast.Contains:
		return derivCall(refs, ast.NewConcat(ast.NewZAny(), ast.NewConcat(v.GetPattern(), ast.NewZAny())), label)
	case *ast.Optional:
		return derivCall(refs, ast.NewOr(v.GetPattern(), ast.NewEmpty()), label)
	}
	panic(fmt.Sprintf("unknown pattern typ %T", typ))
}

func derivCall2(refs map[string]*ast.Pattern, left, right *ast.Pattern, label parser.Value) []*ast.Pattern {
	l := derivCall(refs, left, label)
	r := derivCall(refs, right, label)
	return append(l, r...)
}

func derivReturns(refs map[string]*ast.Pattern, originals []*ast.Pattern, evaluated []*ast.Pattern) []*ast.Pattern {
	res := make([]*ast.Pattern, len(originals))
	rest := evaluated
	for i, original := range originals {
		res[i], rest = derivReturn(refs, original, rest)
	}
	return res
}

func derivReturn(refs map[string]*ast.Pattern, p *ast.Pattern, patterns []*ast.Pattern) (*ast.Pattern, []*ast.Pattern) {
	typ := p.GetValue()
	switch v := typ.(type) {
	case *ast.Empty:
		return ast.NewNot(ast.NewZAny()), patterns
	case *ast.ZAny:
		return ast.NewZAny(), patterns
	case *ast.TreeNode:
		if Nullable(refs, patterns[0]) {
			return ast.NewEmpty(), patterns[1:]
		}
		return ast.NewNot(ast.NewZAny()), patterns[1:]
	case *ast.LeafNode:
		if Nullable(refs, patterns[0]) {
			return ast.NewEmpty(), patterns[1:]
		}
		return ast.NewNot(ast.NewZAny()), patterns[1:]
	case *ast.Concat:
		l, leftRest := derivReturn(refs, v.GetLeftPattern(), patterns)
		leftConcat := ast.NewConcat(l, v.GetRightPattern())
		if !Nullable(refs, v.GetLeftPattern()) {
			return leftConcat, leftRest
		}
		r, rightRest := derivReturn(refs, v.GetRightPattern(), leftRest)
		return ast.NewOr(leftConcat, r), rightRest
	case *ast.Or:
		l, leftRest := derivReturn(refs, v.GetLeftPattern(), patterns)
		r, rightRest := derivReturn(refs, v.GetRightPattern(), leftRest)
		return ast.NewOr(l, r), rightRest
	case *ast.And:
		l, leftRest := derivReturn(refs, v.GetLeftPattern(), patterns)
		r, rightRest := derivReturn(refs, v.GetRightPattern(), leftRest)
		return ast.NewAnd(l, r), rightRest
	case *ast.Interleave:
		l, leftRest := derivReturn(refs, v.GetLeftPattern(), patterns)
		r, rightRest := derivReturn(refs, v.GetRightPattern(), leftRest)
		return ast.NewOr(ast.NewInterleave(l, v.GetRightPattern()), ast.NewInterleave(r, v.GetLeftPattern())), rightRest
	case *ast.ZeroOrMore:
		c, rest := derivReturn(refs, v.GetPattern(), patterns)
		return ast.NewConcat(c, p), rest
	case *ast.Reference:
		return derivReturn(refs, refs[v.GetName()], patterns)
	case *ast.Not:
		c, rest := derivReturn(refs, v.GetPattern(), patterns)
		return ast.NewNot(c), rest
	case *ast.Contains:
		return derivReturn(refs, ast.NewConcat(ast.NewZAny(), ast.NewConcat(v.GetPattern(), ast.NewZAny())), patterns)
	case *ast.Optional:
		return derivReturn(refs, ast.NewOr(v.GetPattern(), ast.NewEmpty()), patterns)
	}
	panic(fmt.Sprintf("unknown pattern typ %T", typ))
}

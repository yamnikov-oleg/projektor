package main

import (
	"go/ast"
	"go/parser"
	"math"
	"strconv"
  "strings"
)

func PerformCalc(query string) LaunchEntriesList {
	expr, err := parser.ParseExpr(query)
	if err != nil {
		return nil
	}

	value, ok := EvaluateGoExpr(expr)
	if !ok {
		return nil
	}

	return LaunchEntriesList{NewCalcLaunchEntry(value)}
}

func EvaluateGoExpr(expr ast.Expr) (val float64, ok bool) {
	switch expr := expr.(type) {
	case *ast.ParenExpr:
		return EvaluateGoExpr(expr.X)

	case *ast.BinaryExpr:
		return EvaluateGoBinaryExpr(expr.X, expr.Op.String(), expr.Y)

	case *ast.BasicLit:
		val, err := strconv.ParseFloat(expr.Value, 64)
		if err != nil {
			return 0, false
		}
		return val, true

	case *ast.CallExpr:
		return EvaluateGoCallExpr(expr.Fun, expr.Args)

	case *ast.Ident:
    val, ok := SupportedConstants[strings.ToLower(expr.String())]
		return val, ok

  case *ast.UnaryExpr:
		// Only negation operator i supported
    if expr.Op.String() != "-" {
			return 0, false
		}
		val, ok := EvaluateGoExpr(expr.X)
		if !ok {
			return 0, false
		}
		return -val, true

	default:
		return 0, false
	}
}

func EvaluateGoBinaryExpr(x ast.Expr, op string, y ast.Expr) (float64, bool) {
	valX, ok := EvaluateGoExpr(x)
	if !ok {
		return 0, false
	}

	valY, ok := EvaluateGoExpr(y)
	if !ok {
		return 0, false
	}

	switch op {
	case "+":
		return valX + valY, true
	case "-":
		return valX - valY, true
	case "*":
		return valX * valY, true
	case "/":
		return valX / valY, true
	case "^":
		return math.Pow(valX, valY), true
	default:
		return 0, false
	}
}

type CalcFunc struct {
	Names    []string
	ArgCount int
	Eval     func(...float64) (float64, bool)
}

func (c *CalcFunc) Matches(name string, argCount int) bool {
	if c.ArgCount != argCount {
		return false
	}

	for _, n := range c.Names {
		if n == name {
			return true
		}
	}
	return false
}

func Factorial(x int64) int64 {
	if x <= 1 {
		return 1
	}
	return x * Factorial(x-1)
}

func EvaluateGoCallExpr(fun ast.Expr, args []ast.Expr) (float64, bool) {
	funIndent, ok := fun.(*ast.Ident)
	if !ok {
		return 0, false
	}

	funName := strings.ToLower(funIndent.String())

	var calcFunc *CalcFunc
	for _, cf := range SupportedCalcFuncs {
		if cf.Matches(funName, len(args)) {
			calcFunc = &cf
			break
		}
	}

	if calcFunc == nil {
		return 0, false
	}

	argVals := []float64{}
	for _, arg := range args {
		val, ok := EvaluateGoExpr(arg)
		if !ok {
			return 0, false
		}
		argVals = append(argVals, val)
	}

	return calcFunc.Eval(argVals...)
}

var SupportedCalcFuncs []CalcFunc = []CalcFunc{
	{
		[]string{"sqrt"}, 1,
		func(args ...float64) (float64, bool) {
			return math.Sqrt(args[0]), true
		},
	},
	{
		[]string{"power", "pow"}, 2,
		func(args ...float64) (float64, bool) {
			return math.Pow(args[0], args[1]), true
		},
	},
	{
		[]string{"root", "rt"}, 2,
		func(args ...float64) (float64, bool) {
			return math.Pow(args[0], 1/args[1]), true
		},
	},
	{
		[]string{"sin"}, 1,
		func(args ...float64) (float64, bool) {
			return math.Sin(args[0]), true
		},
	},
	{
		[]string{"cos"}, 1,
		func(args ...float64) (float64, bool) {
			return math.Cos(args[0]), true
		},
	},
	{
		[]string{"tan", "tg"}, 1,
		func(args ...float64) (float64, bool) {
			return math.Tan(args[0]), true
		},
	},
	{
		[]string{"cot", "ctg"}, 1,
		func(args ...float64) (float64, bool) {
			return 1 / math.Tan(args[0]), true
		},
	},
	{
		[]string{"fact"}, 1,
		func(args ...float64) (float64, bool) {
			return float64(Factorial(int64(args[0]))), true
		},
	},
	{
		[]string{"log", "ln"}, 1,
		func(args ...float64) (float64, bool) {
			return math.Log(args[0]), true
		},
	},
	{
		[]string{"log2"}, 1,
		func(args ...float64) (float64, bool) {
			return math.Log2(args[0]), true
		},
	},
	{
		[]string{"log10"}, 1,
		func(args ...float64) (float64, bool) {
			return math.Log10(args[0]), true
		},
	},
	{
		[]string{"log"}, 2,
		func(args ...float64) (float64, bool) {
			return math.Log(args[1])/math.Log(args[0]), true
		},
	},
}

var SupportedConstants = map[string]float64{
	"pi": math.Pi,
	"e":  math.E,
}

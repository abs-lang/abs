package parser

import (
	"fmt"
	"testing"

	"github.com/abs-lang/abs/ast"
	"github.com/abs-lang/abs/lexer"
	"github.com/abs-lang/abs/token"
)

func TestAssignStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"x = 5;", "x", 5},
		{"y = true;", "y", true},
		{"foobar = y;", "foobar", "y"},
		{"x, y = [1, 2];", "", nil},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]

		switch tt.expectedValue.(type) {
		case string:
			if !testAssignStatement(t, stmt, tt.expectedIdentifier) {
				continue
			}

			val := stmt.(*ast.AssignStatement).Value
			if !testLiteralExpression(t, val, tt.expectedValue) {
				continue
			}
		case nil:
			if len(stmt.(*ast.AssignStatement).Names) == 0 {
				t.Fatalf("stmt.Names does not have any value")
			}
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return", nil},
		{"return;", nil},
		{"return 5", 5},
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.returnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
		if testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func TestNumberLiteralExpression(t *testing.T) {
	prefixTests := []struct {
		input string
		value float64
	}{
		{"5;", 5},
		{"5.5", 5.5},
		{"5.5555555", 5.5555555},
		{"5_000", 5000},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		literal, ok := stmt.Expression.(*ast.NumberLiteral)
		if !ok {
			t.Fatalf("exp not *ast.NumberLiteral. got=%T", stmt.Expression)
		}
		if literal.Value != tt.value {
			t.Errorf("literal.Value not %v. got=%v", tt.value, literal.Value)
		}

		if literal.TokenLiteral() != fmt.Sprintf("%v", tt.value) {
			t.Errorf("number.TokenLiteral not %v. got=%s", tt.value, literal.TokenLiteral())
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 ~ 5;", 5, "~", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 ** 5;", 5, "**", 5},
		{"5 <=> 5;", 5, "<=>", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
		{"1 && 2", 1, "&&", 2},
		{"2 && 1", 2, "&&", 1},
		{"1 || 2", 1, "||", 2},
		{"2 || 1", 2, "||", 1},
		{"1 .. 10", 1, "..", 10},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue,
			tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - ff",
			"(((a + (b * c)) + (d / e)) - ff)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / ff + g)",
			"add((((a + b) + ((c * d) / ff)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean,
				boolean.Value)
		}
	}
}

func TestIfIfElseExpression(t *testing.T) {
	input := `if x < y { x } else if x > y { y } else if x == y { z }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if len(exp.Scenarios) != 3 {
		t.Errorf("expression does not have 4 scenarios. got=%d\n",
			len(exp.Scenarios))
	}

	// First scenario
	scenario := exp.Scenarios[0]

	if !testInfixExpression(t, scenario.Condition, "x", "<", "y") {
		return
	}

	if len(scenario.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(scenario.Consequence.Statements))
	}

	consequence, ok := scenario.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			scenario.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	// Second scenario
	scenario = exp.Scenarios[1]

	if !testInfixExpression(t, scenario.Condition, "x", ">", "y") {
		return
	}

	if len(scenario.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(scenario.Consequence.Statements))
	}

	consequence, ok = scenario.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			scenario.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "y") {
		return
	}

	// Third scenario
	scenario = exp.Scenarios[2]

	if !testInfixExpression(t, scenario.Condition, "x", "==", "y") {
		return
	}

	if len(scenario.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(scenario.Consequence.Statements))
	}

	consequence, ok = scenario.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			scenario.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "z") {
		return
	}
}

func TestIfIfElseElseExpression(t *testing.T) {
	input := `if x < y { x } else if x > y { y } else if x == y { z } else { a }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if len(exp.Scenarios) != 4 {
		t.Errorf("expression does not have 4 scenarios. got=%d\n",
			len(exp.Scenarios))
	}

	// First scenario
	scenario := exp.Scenarios[0]

	if !testInfixExpression(t, scenario.Condition, "x", "<", "y") {
		return
	}

	if len(scenario.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(scenario.Consequence.Statements))
	}

	consequence, ok := scenario.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			scenario.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	// Second scenario
	scenario = exp.Scenarios[1]

	if !testInfixExpression(t, scenario.Condition, "x", ">", "y") {
		return
	}

	if len(scenario.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(scenario.Consequence.Statements))
	}

	consequence, ok = scenario.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			scenario.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "y") {
		return
	}

	// Third scenario
	scenario = exp.Scenarios[2]

	if !testInfixExpression(t, scenario.Condition, "x", "==", "y") {
		return
	}

	if len(scenario.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(scenario.Consequence.Statements))
	}

	consequence, ok = scenario.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			scenario.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "z") {
		return
	}

	// Fourth scenario, the else
	scenario = exp.Scenarios[3]

	if !testBooleanLiteral(t, scenario.Condition, true) {
		return
	}

	if len(scenario.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(scenario.Consequence.Statements))
	}

	consequence, ok = scenario.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			scenario.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "a") {
		return
	}
}

func TestIfExpression(t *testing.T) {
	input := `if x < y { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if len(exp.Scenarios) != 1 {
		t.Errorf("expression does not have 4 scenarios. got=%d\n",
			len(exp.Scenarios))
	}

	// First scenario
	scenario := exp.Scenarios[0]

	if !testInfixExpression(t, scenario.Condition, "x", "<", "y") {
		return
	}

	if len(scenario.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(scenario.Consequence.Statements))
	}

	consequence, ok := scenario.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			scenario.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
}
func TestMoreComplexIfExpression(t *testing.T) {
	input := `x = 1; y = 1; if x > y {
	x
} else {
	y
}

echo(1)
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 4, len(program.Statements))
	}

	stmt, ok := program.Statements[2].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if len(exp.Scenarios) != 2 {
		t.Errorf("expression does not have 2 scenarios. got=%d\n", len(exp.Scenarios))
	}

	// First scenario
	scenario := exp.Scenarios[0]

	if !testInfixExpression(t, scenario.Condition, "x", ">", "y") {
		return
	}

	if len(scenario.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(scenario.Consequence.Statements))
	}

	consequence, ok := scenario.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			scenario.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	return
}

func TestWhileExpression(t *testing.T) {
	input := `
while (x > y) {
	x
}	
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.WhileExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.WhileExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", ">", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
}

func TestForExpression(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`for x = 0; x < y; x = x + 1 {x}`},
		{`for x = 0; x < y; k = increment(k) { x};`},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ForExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ForExpression. got=%T", stmt.Expression)
		}

		if exp.Identifier != "x" {
			t.Errorf("wrong identifier in for loop. got=%s\n", exp.Identifier)
		}

		_, ok = exp.Starter.(*ast.AssignStatement)
		if !ok {
			t.Fatalf("Starter is not ast.AssignStatement. got=%T", exp.Starter)
		}

		if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
			continue
		}

		_, ok = exp.Closer.(*ast.AssignStatement)
		if !ok {
			t.Fatalf("Closer is not ast.AssignStatement. got=%T", exp.Closer)
		}

		block, ok := exp.Block.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Block.Statements[0])
		}

		if !testIdentifier(t, block.Expression, "x") {
			continue
		}
	}
}

func TestForInExpression(t *testing.T) {
	tests := []struct {
		input string
		key   string
		val   string
	}{
		{`for k, v in y {
	x
}`, "k", "v"},
		{`for v in y {
	x
}`, "", "v"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ForInExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ForExpression. got=%T", stmt.Expression)
		}

		if exp.Key != tt.key {
			t.Errorf("wrong key in for loop. got=%s\n", exp.Key)
		}

		if exp.Value != tt.val {
			t.Errorf("wrong val in for loop. got=%s\n", exp.Value)
		}

		if len(exp.Block.Statements) != 1 {
			t.Errorf("block is not 1 statements. got=%d\n", len(exp.Block.Statements))
		}

		block, ok := exp.Block.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Block.Statements[0])
		}

		if !testIdentifier(t, block.Expression, "x") {
			return
		}
	}
}

func TestForElseExpression(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`for x in [] { x } else { y }`},
		{`for x in {} { x } else { y }`},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ForInExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ForExpression. got=%T", stmt.Expression)
		}

		if len(exp.Block.Statements) != 1 {
			t.Errorf("block is not 1 statements. got=%d\n", len(exp.Block.Statements))
		}

		block, ok := exp.Block.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Block.Statements[0])
		}

		if !testIdentifier(t, block.Expression, "x") {
			return
		}

		if len(exp.Alternative.Statements) != 1 {
			t.Errorf("Alternative is not 1 statements. got=%d\n", len(exp.Block.Statements))
		}

		alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Alternative statements[0] is not ast.ExpressionStatement. got=%T", exp.Block.Statements[0])
		}

		if !testIdentifier(t, alternative.Expression, "y") {
			return
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `f(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestCommandParsing(t *testing.T) {
	input := `$(curl icanhazip.com -X POST)`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	command, ok := stmt.Expression.(*ast.CommandExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CommandExpression. got=%T",
			stmt.Expression)
	}

	testCommand(t, command, "curl icanhazip.com -X POST")
}

func TestBacktickCommandParsing(t *testing.T) {
	input := "`curl icanhazip.com -X POST`"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	command, ok := stmt.Expression.(*ast.CommandExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CommandExpression. got=%T",
			stmt.Expression)
	}

	testCommand(t, command, "curl icanhazip.com -X POST")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "f() {};", expectedParams: []string{}},
		{input: "f(x) {};", expectedParams: []string{"x"}},
		{input: "f(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Arguments[i].String())
			}
		}
	}
}

func TestMethodExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedMethod string
		expectedObj    string
		expectedArgs   []string
	}{
		{
			input:          "test.method(1, 2 * 3, 4 + 5);",
			expectedObj:    "test",
			expectedMethod: "method",
			expectedArgs:   []string{"1", "(2 * 3)", "(4 + 5)"},
		},
		{
			input:          "a.method_name(1);",
			expectedObj:    "a",
			expectedMethod: "method_name",
			expectedArgs:   []string{"1"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expression.(*ast.MethodExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.MethodExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Object, tt.expectedObj) {
			continue
		}

		if !testIdentifier(t, exp.Method, tt.expectedMethod) {
			continue
		}

		if len(exp.Arguments) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Arguments[i].String())
			}
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestNullLiteralExpression(t *testing.T) {
	input := `null;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.NullLiteral)
	if !ok {
		t.Fatalf("exp not *ast.NullLiteral. got=%T", stmt.Expression)
	}

	if literal.TokenLiteral() != "null" {
		t.Errorf("literal.Value not %q. got=%s", "null", literal.TokenLiteral())
	}
}

func TestParsingEmptyArrayLiterals(t *testing.T) {
	input := "[]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 0 {
		t.Errorf("len(array.Elements) not 0. got=%d", len(array.Elements))
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testNumberLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingIndexRangeExpressions(t *testing.T) {
	input := "myArray[99 : 101]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !indexExp.IsRange {
		t.Fatalf("exp is not range")
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	testNumberLiteral(t, indexExp.Index, 99)
	testNumberLiteral(t, indexExp.End, 101)
}

func TestParsingIndexRangeWithoutStartExpressions(t *testing.T) {
	input := "myArray[: 101]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !indexExp.IsRange {
		t.Fatalf("exp is not range")
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	testNumberLiteral(t, indexExp.Index, 0)
	testNumberLiteral(t, indexExp.End, 101)
}

func TestParsingIndexRangeWithoutEndExpressions(t *testing.T) {
	input := "myArray[99 : ]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !indexExp.IsRange {
		t.Fatalf("exp is not range")
	}

	testNumberLiteral(t, indexExp.Index, 99)

	if indexExp.End != nil {
		t.Fatalf("range end is not nil. got=%T", indexExp.End)
	}
}

func TestParsingProperty(t *testing.T) {
	input := "var.prop"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	propExpr, ok := stmt.Expression.(*ast.PropertyExpression)

	if !ok {
		t.Fatalf("exp not *ast.PropertyExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, propExpr.Object, "var") {
		return
	}

	if !testIdentifier(t, propExpr.Property, "prop") {
		return
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	expected := map[string]float64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[literal.String()]
		testNumberLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsBooleanKeys(t *testing.T) {
	input := `{true: 1, false: 2}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	expected := map[string]float64{
		"true":  1,
		"false": 2,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		boolean, ok := key.(*ast.Boolean)
		if !ok {
			t.Errorf("key is not ast.BooleanLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[boolean.String()]
		testNumberLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsNumberKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	expected := map[string]float64{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		Number, ok := key.(*ast.NumberLiteral)
		if !ok {
			t.Errorf("key is not ast.NumberLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[Number.String()]

		testNumberLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}

		testFunc(value)
	}
}

func testAssignStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != token.ASSIGN {
		t.Errorf("s.TokenLiteral not '%s'. got=%q", token.ASSIGN, s.TokenLiteral())
		return false
	}

	stmt, ok := s.(*ast.AssignStatement)
	if !ok {
		t.Errorf("s not *ast.AssignStatement. got=%T", s)
		return false
	}

	if stmt.Name.Value != name {
		t.Errorf("stmt.Name.Value not '%s'. got=%s", name, stmt.Name.Value)
		return false
	}

	if stmt.Name.TokenLiteral() != name {
		t.Errorf("stmt.Name.TokenLiteral() not '%s'. got=%s",
			name, stmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testNumberLiteral(t, exp, float64(v))
	case float64:
		return testNumberLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	case nil:
		return testNullLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testNumberLiteral(t *testing.T, il ast.Expression, value float64) bool {
	number, ok := il.(*ast.NumberLiteral)
	if !ok {
		t.Errorf("il not *ast.NumberLiteral. got=%T", il)
		return false
	}

	if number.Value != value {
		t.Errorf("number.Value not %f. got=%f", value, number.Value)
		return false
	}

	if number.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("number.TokenLiteral not %v. got=%s", value, number.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testCommand(t *testing.T, exp ast.Expression, value string) bool {
	command, ok := exp.(*ast.CommandExpression)
	if !ok {
		t.Errorf("exp not *ast.CommandExpression. got=%T", exp)
		return false
	}

	if command.Value != value {
		t.Errorf("command.Value not %s. got=%s", value, command.Value)
		return false
	}

	if command.TokenLiteral() != value {
		t.Errorf("command.TokenLiteral not %s. got=%s", value,
			command.TokenLiteral())
		return false
	}

	return true
}

func testMethod(t *testing.T, exp ast.Expression, value string) bool {
	command, ok := exp.(*ast.CommandExpression)
	if !ok {
		t.Errorf("exp not *ast.CommandExpression. got=%T", exp)
		return false
	}

	if command.Value != value {
		t.Errorf("command.Value not %s. got=%s", value, command.Value)
		return false
	}

	if command.TokenLiteral() != value {
		t.Errorf("command.TokenLiteral not %s. got=%s", value,
			command.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}

func testNullLiteral(t *testing.T, exp ast.Expression, value interface{}) bool {
	nl, ok := exp.(*ast.NullLiteral)
	if !ok {
		t.Errorf("exp not *ast.NullLiteral. got=%T", exp)
		return false
	}

	if nl.TokenLiteral() != "null" {
		t.Errorf("nl.TokenLiteral not %t. got=%s", value, nl.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

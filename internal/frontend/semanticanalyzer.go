package frontend

import (
	"fmt"

	"github.com/SamJohn04/simple-lang-compiler/internal/common"
)

// converts the Parse Tree to AST
func SemanticAnalyzer(
	input common.ParseTreeNode,
) (common.ProgramAST, []common.IdentifierInformation, error) {
	identifiers := []common.IdentifierInformation{}
	program, identifiers, err := lowerProgram(input, identifiers)
	return program, identifiers, err
}

func lowerProgram(
	input common.ParseTreeNode,
	identifiers []common.IdentifierInformation,
) (common.ProgramAST, []common.IdentifierInformation, error) {
	programAST := common.ProgramAST{
		Instructions: []common.InstructionAST{},
	}

	current := input
	// it will exit when current has no children
	for len(current.ChildNodes) > 0 {
		if current.InnerToken.TokenKind != common.TokenBlock {
			return common.ProgramAST{}, identifiers, semanticInternalError(
				"instruction is expected to be a token block",
			)
		}
		if len(current.ChildNodes) != 2 {
			return common.ProgramAST{}, identifiers, semanticInternalError(
				"instruction has neither 0 nor 2 blocks",
			)
		}

		var instructionAST common.InstructionAST
		var err error

		// lower the instruction
		instructionAST, identifiers, err = lowerInstruction(
			current.ChildNodes[0],
			identifiers,
		)
		if err != nil {
			return programAST, identifiers, err
		}
		if instructionAST != nil {
			programAST.Instructions = append(programAST.Instructions, instructionAST)
		}

		// update current
		current = current.ChildNodes[1]
	}
	return programAST, identifiers, nil
}

func lowerInstruction(
	instruction common.ParseTreeNode,
	identifiers []common.IdentifierInformation,
) (common.InstructionAST, []common.IdentifierInformation, error) {
	if len(instruction.ChildNodes) == 0 {
		return nil, identifiers, semanticInternalError(
			"instruction structure is not expected to be empty",
		)
	}

	switch instruction.ChildNodes[0].InnerToken.TokenKind {
	case common.TokenIdent:
		// reassignment
		return lowerReassignment(instruction, identifiers)

	case common.TokenLet:
		// assignment
		return lowerAssignment(instruction, identifiers)

	case common.TokenIf:
		// if
		return lowerIfStatement(instruction, identifiers)

	case common.TokenWhile:
		// while
		return lowerWhileStatement(instruction, identifiers)

	case common.TokenOutput:
		// output
		output, err := lowerOutputStatement(instruction, identifiers)
		return output, identifiers, err

	default:
		return nil, identifiers, semanticInternalError(
			fmt.Sprintf(
				"instruction first child is of a different type (%v)",
				common.NameMapWithTokenKind[instruction.ChildNodes[0].InnerToken.TokenKind],
			),
		)
	}
}

func lowerReassignment(
	instruction common.ParseTreeNode,
	identifiers []common.IdentifierInformation,
) (common.AssignmentAST, []common.IdentifierInformation, error) {
	if len(instruction.ChildNodes) != 4 {
		return common.AssignmentAST{}, identifiers, semanticInternalError(
			"identifier instruction not having expected length",
		)
	}
	childIdentifier := instruction.ChildNodes[0]
	index := find(identifiers, childIdentifier.InnerToken.Token)
	if index < 0 {
		return common.AssignmentAST{}, identifiers, semanticError(
			"identifier used before being declared",
		)
	}

	information := identifiers[index]
	if !information.Mutable {
		return common.AssignmentAST{}, identifiers, semanticError(
			"identifier that was not declared as mutable being mutated",
		)
	}
	childArrayUsage, err := lowerArrayUsage(instruction.ChildNodes[1], identifiers)
	if err != nil {
		return common.AssignmentAST{}, identifiers, err
	}
	childEquals := instruction.ChildNodes[2]
	if childEquals.InnerToken.TokenKind != common.TokenAssignment {
		return common.AssignmentAST{}, identifiers, semanticInternalError("'=' expected")
	}
	childR, err := lowerRelation(instruction.ChildNodes[3], identifiers)
	if err != nil {
		return common.AssignmentAST{}, identifiers, err
	}

	return common.AssignmentAST{
		AssignToIdentifier: index,
		ArrayValues:        childArrayUsage,
		AssignValue:        childR,
	}, identifiers, nil
}

func lowerArrayUsage(
	arrayInstruction common.ParseTreeNode, identifiers []common.IdentifierInformation,
) ([]common.ExpressionAST, error) {
	arrays := []common.ExpressionAST{}
	for len(arrayInstruction.ChildNodes) > 0 {
		if len(arrayInstruction.ChildNodes) != 4 {
			return []common.ExpressionAST{}, semanticInternalError(
				"array length not 0 or 4",
			)
		}
		if arrayInstruction.ChildNodes[0].InnerToken.TokenKind != common.TokenOpenSquareBraces ||
			arrayInstruction.ChildNodes[2].InnerToken.TokenKind != common.TokenCloseSquareBraces {
			return []common.ExpressionAST{}, semanticInternalError(
				"mismatching open and close square braces",
			)
		}
		childE, err := lowerE(arrayInstruction.ChildNodes[1], identifiers)
		if err != nil {
			return []common.ExpressionAST{}, err
		}
		arrays = append(arrays, childE)
		arrayInstruction = arrayInstruction.ChildNodes[3]
	}
	return arrays, nil
}

func lowerAssignment(
	instruction common.ParseTreeNode, identifiers []common.IdentifierInformation,
) (common.InstructionAST, []common.IdentifierInformation, error) {
	if len(instruction.ChildNodes) != 2 {
		return common.AssignmentAST{}, identifiers, semanticInternalError(
			"let instruction not having 2 length",
		)
	}
	childInstruction, identifiers, err := lowerAssignmentAfterLet(instruction.ChildNodes[1], identifiers)
	if err != nil {
		return common.AssignmentAST{}, identifiers, err
	}
	return childInstruction, identifiers, nil
}

func lowerAssignmentAfterLet(
	instruction common.ParseTreeNode, identifiers []common.IdentifierInformation,
) (common.InstructionAST, []common.IdentifierInformation, error) {
	if len(instruction.ChildNodes) != 3 {
		return common.AssignmentAST{}, identifiers, semanticInternalError(
			"instruction after let should have length 3",
		)
	}
	switch instruction.ChildNodes[0].InnerToken.TokenKind {
	case common.TokenIdent:
		// v = R
		childIdentifier := instruction.ChildNodes[0]
		index := find(identifiers, childIdentifier.InnerToken.Token)
		if index >= 0 {
			return nil, identifiers, semanticError("identifier already declared")
		}
		if instruction.ChildNodes[1].InnerToken.TokenKind != common.TokenAssignment {
			return common.AssignmentAST{}, identifiers, semanticInternalError(
				"instruction after let does not have assignment",
			)
		}
		childR, err := lowerRelation(instruction.ChildNodes[2], identifiers)
		if err != nil {
			return common.AssignmentAST{}, identifiers, err
		}
		identifiers = append(identifiers, common.IdentifierInformation{
			IdentifierName: childIdentifier.InnerToken.Token,
			Mutable:        false,
		})
		return common.AssignmentAST{
			AssignToIdentifier: len(identifiers) - 1,
			ArrayValues:        []common.ExpressionAST{},
			AssignValue:        childR,
		}, identifiers, nil

	case common.TokenMutable:
		// mut v ...
		childIdentifier := instruction.ChildNodes[1]
		index := find(identifiers, childIdentifier.InnerToken.Token)
		if index >= 0 {
			return nil, identifiers, semanticError("identifier already declared")
		}
		childExpression, identifiers, err := lowerMutableAssignment(
			instruction.ChildNodes[2], childIdentifier, identifiers,
		)
		if err != nil {
			return nil, identifiers, err
		}
		return childExpression, identifiers, nil

	default:
		return nil, identifiers, semanticInternalError("unexpected token after let")
	}
}

func lowerMutableAssignment(
	instruction, identifier common.ParseTreeNode, identifiers []common.IdentifierInformation,
) (common.InstructionAST, []common.IdentifierInformation, error) {
	identifiers = append(identifiers, common.IdentifierInformation{
		IdentifierName: identifier.InnerToken.Token,
		Mutable:        true,
	})
	if len(instruction.ChildNodes) == 0 {
		// simply declaring the variable
		return nil, identifiers, nil
	}
	if len(instruction.ChildNodes) != 2 {
		return nil, identifiers, semanticInternalError("mut should have 0 or 2 children")
	}
	childEquals := instruction.ChildNodes[0]
	if childEquals.InnerToken.TokenKind != common.TokenAssignment {
		return nil, identifiers, semanticInternalError("'=' expected")
	}
	childExpression, err := lowerRelation(instruction.ChildNodes[1], identifiers)
	if err != nil {
		return nil, identifiers, err
	}
	return common.AssignmentAST{
		AssignToIdentifier: len(identifiers) - 1,
		ArrayValues:        []common.ExpressionAST{},
		AssignValue:        childExpression,
	}, identifiers, nil
}

func lowerIfStatement(
	instruction common.ParseTreeNode, identifiers []common.IdentifierInformation,
) (common.IfStatementAST, []common.IdentifierInformation, error) {
	ifStatement := common.IfStatementAST{
		IfExpressions: []common.IfExpression{},
	}
	if len(instruction.ChildNodes) == 0 {
		return ifStatement, identifiers, semanticInternalError("if does not have children")
	}
	for len(instruction.ChildNodes) > 0 {
		var childProgram common.ProgramAST
		var err error

		switch instruction.ChildNodes[0].InnerToken.TokenKind {
		case common.TokenIf:
			if len(instruction.ChildNodes) != 6 {
				return ifStatement, identifiers, semanticInternalError(
					"if has unexpected number of children",
				)
			}
			childR, err := lowerRelation(instruction.ChildNodes[1], identifiers)
			if err != nil {
				return ifStatement, identifiers, err
			}
			if instruction.ChildNodes[2].InnerToken.TokenKind != common.TokenOpenCurly ||
				instruction.ChildNodes[4].InnerToken.TokenKind != common.TokenCloseCurly {
				return ifStatement, identifiers, semanticInternalError(
					"if does not have opening or closing braces",
				)
			}

			childProgram, identifiers, err = lowerProgram(instruction.ChildNodes[3], identifiers)
			if err != nil {
				return ifStatement, identifiers, err
			}
			ifStatement.IfExpressions = append(ifStatement.IfExpressions, common.IfExpression{
				Condition: childR,
				Program:   childProgram,
			})

		case common.TokenOpenCurly:
			if len(instruction.ChildNodes) != 3 {
				return ifStatement, identifiers, semanticInternalError(
					"if has unexpected number of children",
				)
			}
			if instruction.ChildNodes[2].InnerToken.TokenKind != common.TokenCloseCurly {
				return ifStatement, identifiers, semanticInternalError(
					"if does not have closing braces",
				)
			}
			childProgram, identifiers, err = lowerProgram(instruction.ChildNodes[1], identifiers)
			if err != nil {
				return ifStatement, identifiers, err
			}
			// else { ... } is turned into an else if (true) { ... }
			ifStatement.IfExpressions = append(ifStatement.IfExpressions, common.IfExpression{
				Condition: common.Literal{
					Value:    "true",
					Datatype: common.TypedBool,
				},
				Program: childProgram,
			})
			return ifStatement, identifiers, nil

		default:
			return ifStatement, identifiers, semanticInternalError("unexpected token at if")
		}

		// only TokenIf reaches here
		instruction = instruction.ChildNodes[5]
		if len(instruction.ChildNodes) == 0 {
			return ifStatement, identifiers, nil
		}
		if len(instruction.ChildNodes) != 2 ||
			instruction.ChildNodes[0].InnerToken.TokenKind != common.TokenElse {
			return ifStatement, identifiers, semanticInternalError("unexpected token length at else")
		}
		instruction = instruction.ChildNodes[1]
	}
	return ifStatement, identifiers, nil
}

func lowerWhileStatement(
	instruction common.ParseTreeNode, identifiers []common.IdentifierInformation,
) (common.WhileStatementAST, []common.IdentifierInformation, error) {
	if len(instruction.ChildNodes) != 5 {
		return common.WhileStatementAST{}, identifiers, semanticInternalError(
			"unexpected length of while",
		)
	}
	childR, err := lowerRelation(instruction.ChildNodes[1], identifiers)
	if err != nil {
		return common.WhileStatementAST{}, identifiers, err
	}
	if instruction.ChildNodes[2].InnerToken.TokenKind != common.TokenOpenCurly ||
		instruction.ChildNodes[4].InnerToken.TokenKind != common.TokenCloseCurly {
		return common.WhileStatementAST{}, identifiers, semanticInternalError(
			"open or closing brace missing",
		)
	}
	childProgram, identifiers, err := lowerProgram(instruction.ChildNodes[3], identifiers)
	if err != nil {
		return common.WhileStatementAST{}, identifiers, err
	}
	return common.WhileStatementAST{
		Condition: childR,
		Program:   childProgram,
	}, identifiers, nil
}

func lowerOutputStatement(
	instruction common.ParseTreeNode, identifiers []common.IdentifierInformation,
) (common.OutputStatementAST, error) {
	if len(instruction.ChildNodes) != 3 {
		return common.OutputStatementAST{}, semanticInternalError("output not having 3 children")
	}
	outputStatement := common.OutputStatementAST{
		Arguments: []common.ExpressionAST{
			common.Literal{
				Value: instruction.ChildNodes[1].InnerToken.Token,
				Datatype: common.StringDatatype{
					HasKnownLength: true,
					CharacterCount: len(instruction.ChildNodes[1].InnerToken.Token) - 2,
				},
			},
		},
	}
	childC := instruction.ChildNodes[2]
	for len(childC.ChildNodes) > 0 {
		if len(childC.ChildNodes) != 2 {
			return outputStatement, semanticInternalError("output continuation not having 0 or 2 children")
		}
		childR, err := lowerRelation(childC.ChildNodes[0], identifiers)
		if err != nil {
			return outputStatement, err
		}
		outputStatement.Arguments = append(outputStatement.Arguments, childR)
		childC = childC.ChildNodes[1]
	}
	return outputStatement, nil
}

func lowerRelation(
	relationInstruction common.ParseTreeNode, identifiers []common.IdentifierInformation,
) (common.ExpressionAST, error) {
	if len(relationInstruction.ChildNodes) != 2 {
		return nil, semanticInternalError("Relation expected to have two children")
	}
	expression, err := lowerRa(relationInstruction.ChildNodes[0], identifiers)
	if err != nil {
		return nil, err
	}
	expression, err = lowerRz(relationInstruction.ChildNodes[1], expression, identifiers)
	if err != nil {
		return nil, err
	}
	return expression, nil
}

func lowerRa(
	relationInstruction common.ParseTreeNode, identifiers []common.IdentifierInformation,
) (common.ExpressionAST, error) {
	if len(relationInstruction.ChildNodes) != 2 {
		return nil, semanticInternalError("Relation (Ra) expected to have two children")
	}
	expression, err := lowerRb(relationInstruction.ChildNodes[0], identifiers)
	if err != nil {
		return nil, err
	}
	expression, err = lowerRy(relationInstruction.ChildNodes[1], expression, identifiers)
	if err != nil {
		return nil, err
	}
	return expression, nil
}

func lowerRz(
	relationInstruction common.ParseTreeNode, calculationsUntilNow common.ExpressionAST,
	identifiers []common.IdentifierInformation,
) (common.ExpressionAST, error) {
	if len(relationInstruction.ChildNodes) == 0 {
		return calculationsUntilNow, nil
	}
	if len(relationInstruction.ChildNodes) != 3 {
		return nil, semanticInternalError("Rz expected to have 0 or 3 children")
	}
	if relationInstruction.ChildNodes[0].InnerToken.TokenKind != common.TokenOr {
		return nil, semanticInternalError("|| expected in Rz")
	}
	secondOperand, err := lowerRa(relationInstruction.ChildNodes[1], identifiers)
	if err != nil {
		return nil, err
	}
	binaryRelation := common.BinaryExpression{
		Operator:      common.BinaryOr,
		FirstOperand:  calculationsUntilNow,
		SecondOperand: secondOperand,
	}
	expression, err := lowerRz(relationInstruction.ChildNodes[2], binaryRelation, identifiers)
	if err != nil {
		return nil, err
	}
	return expression, nil
}

func lowerRb(
	relationInstruction common.ParseTreeNode, identifiers []common.IdentifierInformation,
) (common.ExpressionAST, error) {
	if len(relationInstruction.ChildNodes) != 2 {
		return nil, semanticInternalError("unexpected number of children in Rb")
	}
	if relationInstruction.ChildNodes[0].InnerToken.TokenKind == common.TokenNot {
		childCalculations, err := lowerRelation(
			relationInstruction.ChildNodes[1], identifiers,
		)
		if err != nil {
			return nil, err
		}
		return common.UnaryExpression{
			Operator: common.UnaryNot,
			Operand:  childCalculations,
		}, nil
	}
	firstExpression, err := lowerE(relationInstruction.ChildNodes[0], identifiers)
	if err != nil {
		return nil, err
	}
	if len(relationInstruction.ChildNodes[1].ChildNodes) == 0 {
		return firstExpression, nil
	}
	if len(relationInstruction.ChildNodes[1].ChildNodes) != 2 {
		return nil, semanticInternalError("unexpected number of children in child of Rb")
	}

	secondExpression, err := lowerE(relationInstruction.ChildNodes[1].ChildNodes[1], identifiers)
	if err != nil {
		return nil, err
	}

	expression := common.BinaryExpression{
		FirstOperand:  firstExpression,
		SecondOperand: secondExpression,
	}

	switch relationInstruction.ChildNodes[1].ChildNodes[0].InnerToken.TokenKind {
	case common.TokenRelationalEquals:
		expression.Operator = common.BinaryRelationalEquals

	case common.TokenRelationalNotEquals:
		expression.Operator = common.BinaryRelationalNotEquals

	case common.TokenRelationalGreaterThan:
		expression.Operator = common.BinaryRelationalGreaterThan

	case common.TokenRelationalGreaterThanOrEquals:
		expression.Operator = common.BinaryRelationalGreaterThanOrEquals

	case common.TokenRelationalLesserThan:
		expression.Operator = common.BinaryRelationalLesserThan

	case common.TokenRelationalLesserThanOrEquals:
		expression.Operator = common.BinaryRelationalLesserThanOrEquals

	default:
		return nil, semanticInternalError("unknown operand in relation")
	}
	return expression, nil
}

func lowerRy(
	relationInstruction common.ParseTreeNode, calculationsUntilNow common.ExpressionAST,
	identifiers []common.IdentifierInformation,
) (common.ExpressionAST, error) {
	if len(relationInstruction.ChildNodes) == 0 {
		return calculationsUntilNow, nil
	}
	if len(relationInstruction.ChildNodes) != 3 {
		return nil, semanticInternalError("Ry expected to have 0 or 3 children")
	}
	if relationInstruction.ChildNodes[0].InnerToken.TokenKind != common.TokenAnd {
		return nil, semanticInternalError("&& expected in Ry")
	}
	secondOperand, err := lowerRb(relationInstruction.ChildNodes[1], identifiers)
	if err != nil {
		return nil, err
	}
	binaryRelation := common.BinaryExpression{
		Operator:      common.BinaryAnd,
		FirstOperand:  calculationsUntilNow,
		SecondOperand: secondOperand,
	}
	expression, err := lowerRy(relationInstruction.ChildNodes[2], binaryRelation, identifiers)
	if err != nil {
		return nil, err
	}
	return expression, nil
}

func lowerE(
	expression common.ParseTreeNode, identifiers []common.IdentifierInformation,
) (common.ExpressionAST, error) {
	if len(expression.ChildNodes) != 2 {
		expression.Display("", ">")
		return nil, semanticInternalError("expression does not have two children")
	}
	childT, err := lowerT(expression.ChildNodes[0], identifiers)
	if err != nil {
		return nil, err
	}
	childE1, err := lowerE1(expression.ChildNodes[1], childT, identifiers)
	if err != nil {
		return nil, err
	}
	return childE1, nil
}

func lowerT(
	expression common.ParseTreeNode, identifiers []common.IdentifierInformation,
) (common.ExpressionAST, error) {
	if len(expression.ChildNodes) != 2 {
		return nil, semanticInternalError("expression T does not have two children")
	}
	childF, err := lowerF(expression.ChildNodes[0], identifiers)
	if err != nil {
		return nil, err
	}
	childT1, err := lowerT1(expression.ChildNodes[1], childF, identifiers)
	if err != nil {
		return nil, err
	}
	return childT1, nil
}

func lowerE1(
	expression common.ParseTreeNode, calculationsUntilNow common.ExpressionAST,
	identifiers []common.IdentifierInformation,
) (common.ExpressionAST, error) {
	if len(expression.ChildNodes) == 0 {
		return calculationsUntilNow, nil
	}
	if len(expression.ChildNodes) != 3 {
		return nil, semanticInternalError("E1 has unexpected number of array elements")
	}

	secondExpression, err := lowerT(expression.ChildNodes[1], identifiers)
	if err != nil {
		return nil, err
	}

	binaryExpression := common.BinaryExpression{
		FirstOperand:  calculationsUntilNow,
		SecondOperand: secondExpression,
	}

	switch expression.ChildNodes[0].InnerToken.TokenKind {
	case common.TokenExpressionAdd:
		binaryExpression.Operator = common.BinaryPlus

	case common.TokenExpressionSub:
		binaryExpression.Operator = common.BinaryMinus

	default:
		return nil, semanticInternalError("unexpected operator in E1")
	}

	return binaryExpression, nil
}

func lowerF(
	expression common.ParseTreeNode, identifiers []common.IdentifierInformation,
) (common.ExpressionAST, error) {
	if len(expression.ChildNodes) == 0 {
		return nil, semanticInternalError("expression F needs at least 1 element")
	}
	switch expression.ChildNodes[0].InnerToken.TokenKind {
	case common.TokenOpenSquareBraces:
		if len(expression.ChildNodes) != 3 {
			return nil, semanticInternalError("expression for parsing arrays need 3 elements")
		}
		if expression.ChildNodes[2].InnerToken.TokenKind != common.TokenCloseSquareBraces {
			return nil, semanticInternalError("expected ]")
		}
		array, err := lowerArrayExpression(expression.ChildNodes[1], identifiers)
		if err != nil {
			return nil, err
		}
		return array, nil

	case common.TokenExpressionSub:
		if len(expression.ChildNodes) != 2 {
			return nil, semanticInternalError("expression for unary minus needs 2 elements")
		}
		childExpression, err := lowerF(expression.ChildNodes[1], identifiers)
		if err != nil {
			return nil, err
		}
		return common.UnaryExpression{
			Operator: common.UnaryMinus,
			Operand:  childExpression,
		}, nil

	case common.TokenBlock:
		if len(expression.ChildNodes) != 1 {
			return nil, semanticInternalError("expression block should have no siblings")
		}
		childExpression, err := lowerRelation(expression.ChildNodes[0], identifiers)
		if err != nil {
			return nil, err
		}
		return childExpression, nil

	case common.TokenIdent:
		childIdentifier, err := lowerIdentifierAfterDeclaration(expression, identifiers)
		if err != nil {
			return nil, err
		}
		return childIdentifier, nil

	case common.TokenLiteralInt:
		if len(expression.ChildNodes) != 1 {
			return nil, semanticInternalError("F should have no siblings")
		}
		return common.Literal{
			Value:    expression.ChildNodes[0].InnerToken.Token,
			Datatype: common.TypedInt,
		}, nil

	case common.TokenLiteralChar:
		if len(expression.ChildNodes) != 1 {
			return nil, semanticInternalError("F should have no siblings")
		}
		return common.Literal{
			Value:    expression.ChildNodes[0].InnerToken.Token,
			Datatype: common.TypedChar,
		}, nil

	case common.TokenLiteralBool:
		if len(expression.ChildNodes) != 1 {
			return nil, semanticInternalError("F should have no siblings")
		}
		return common.Literal{
			Value:    expression.ChildNodes[0].InnerToken.Token,
			Datatype: common.TypedBool,
		}, nil

	case common.TokenLiteralFloat:
		if len(expression.ChildNodes) != 1 {
			return nil, semanticInternalError("F should have no siblings")
		}
		return common.Literal{
			Value:    expression.ChildNodes[0].InnerToken.Token,
			Datatype: common.TypedFloat,
		}, nil

	case common.TokenLiteralString:
		if len(expression.ChildNodes) != 1 {
			return nil, semanticInternalError("F should have no siblings")
		}
		return common.Literal{
			Value: expression.ChildNodes[0].InnerToken.Token,
			Datatype: common.StringDatatype{
				HasKnownLength: true,
				CharacterCount: len(expression.InnerToken.Token) - 2,
			},
		}, nil

	case common.TokenInput:
		if len(expression.ChildNodes) != 1 {
			return nil, semanticInternalError("input statement should have no siblings")
		}
		return common.InputExpression{}, nil

	default:
		fmt.Println(common.NameMapWithTokenKind[expression.ChildNodes[0].InnerToken.TokenKind])
		return nil, semanticInternalError("unexpected token type in F")
	}
}

func lowerT1(
	expression common.ParseTreeNode, calculationsUntilNow common.ExpressionAST,
	identifiers []common.IdentifierInformation,
) (common.ExpressionAST, error) {
	if len(expression.ChildNodes) == 0 {
		return calculationsUntilNow, nil
	}
	if len(expression.ChildNodes) != 3 {
		return nil, semanticInternalError("T1 has unexpected number of elements")
	}

	secondExpression, err := lowerF(expression.ChildNodes[1], identifiers)
	if err != nil {
		return nil, err
	}

	binaryOperation := common.BinaryExpression{
		FirstOperand:  calculationsUntilNow,
		SecondOperand: secondExpression,
	}

	switch expression.ChildNodes[0].InnerToken.TokenKind {
	case common.TokenExpressionMul:
		binaryOperation.Operator = common.BinaryMul

	case common.TokenExpressionDiv:
		binaryOperation.Operator = common.BinaryDiv

	case common.TokenExpressionModulo:
		binaryOperation.Operator = common.BinaryModulo

	default:
		return nil, semanticInternalError("unexpected operand in T1")
	}

	return binaryOperation, nil
}

func lowerArrayExpression(
	arrayExpression common.ParseTreeNode, identifiers []common.IdentifierInformation,
) (common.ArrayExpression, error) {
	array := common.ArrayExpression{
		Elements: []common.ExpressionAST{},
	}
	for len(arrayExpression.ChildNodes) > 0 {
		if len(arrayExpression.ChildNodes) != 2 {
			return array, semanticInternalError("expected 2 children while parsing array")
		}
		element, err := lowerRelation(arrayExpression.ChildNodes[0], identifiers)
		if err != nil {
			return array, err
		}
		array.Elements = append(array.Elements, element)
		arrayExpression = arrayExpression.ChildNodes[1]
	}
	return array, nil
}

func lowerIdentifierAfterDeclaration(
	input common.ParseTreeNode, identifiers []common.IdentifierInformation,
) (common.Identifier, error) {
	if len(input.ChildNodes) != 2 {
		return common.Identifier{}, semanticInternalError("identifier expected to be two elements")
	}
	if input.ChildNodes[0].InnerToken.TokenKind != common.TokenIdent ||
		input.ChildNodes[1].InnerToken.TokenKind != common.TokenBlock {
		return common.Identifier{}, semanticInternalError("identifier and block expected")
	}
	arrayUsage, err := lowerArrayUsage(input.ChildNodes[1], identifiers)
	if err != nil {
		return common.Identifier{}, err
	}
	index := find(identifiers, input.ChildNodes[0].InnerToken.Token)
	if index < 0 {
		return common.Identifier{}, semanticError("identifier used before being declared")
	}
	return common.Identifier{
		Id:          index,
		ArrayValues: arrayUsage,
	}, nil
}

func semanticError(message string) *common.CompilationError {
	return &common.CompilationError{
		PointOfFailure: "Semantic Analyzer",
		Message:        message,
	}
}

func semanticInternalError(message string) *common.InternalError {
	return &common.InternalError{
		PointOfFailure: "Semantic Analyzer",
		Message:        message,
	}
}

func find(
	identifiers []common.IdentifierInformation, lookingFor string,
) int {
	for index, identifier := range identifiers {
		if identifier.IdentifierName == lookingFor {
			return index
		}
	}
	return -1
}

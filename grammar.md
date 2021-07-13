ExpressionList = Expression* .
Expression = Assignment .
Assignment = LogicAnd ("=" LogicAnd)? .
LogicAnd = LogicOr ("and" LogicAnd)? .
LogicOr = Equality ("or" LogicOr)? .
Equality = Comparison ((("!" "=") | ("=" "=")) Equality)? .
Comparison = Addition (((">" "=") | ">" | ("<" "=") | "<") Comparison)? .
Addition = Multiplication (("-" | "+") Addition)? .
Multiplication = Unary (("/" | "*" | "%") Multiplication)? .
Unary = (("!" | "-") Unary) | Primary .
Primary = If | DataLiteral | ("(" Expression ")") | Call | ForKeyValue | ForValue | For | ForWhile | Return | Break | Continue | <float> | <int> | <string> | "true" | "false" | "nil" | <ident> .
If = "if" Expression "{" Expression* "}" ElseIf* ("else" "{" Expression* "}")? .
ElseIf = "else" "if" Expression "{" Expression* "}" ElseIf* .
DataLiteral = FunctionLiteral | ListLiteral | DictLiteral .
FunctionLiteral = (("(" (<ident> ("," <ident>)*)? ")") | <ident>) "=" ">" (("{" Expression* "}") | Expression) .
ListLiteral = "[" (Expression ("," Expression)*)? "]" .
DictLiteral = "{" (DictEntry ("," DictEntry)* ","?)? "}" .
DictEntry = (<ident> | Expression) ":" Expression .
Call = (<ident> | ("(" Expression ")")) CallChain .
CallChain = (("(" (Expression ("," Expression)*)? ")") | ("." <ident>) | ("[" Expression "]")) CallChain? .
ForKeyValue = ("for" <ident> "," <ident> "in" (<ident> | Expression) "{" Expression* "}") .
ForValue = ("for" <ident> "in" (<ident> | Expression) "{" Expression* "}") .
For = "for" (Assignment ";" Expression ";" Expression "{" Expression* "}") .
ForWhile = "for" (Expression? "{" Expression* "}") .
Return = ("return" Expression) .
Break = "break" .
Continue = "continue" .

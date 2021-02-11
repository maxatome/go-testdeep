%{
// Copyright (c) 2020, 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package json

type Placeholder struct {
	Num  int
	Name string
}

type Operator struct {
	Name string
	Params   []interface{}
}

type member struct {
	key   string
	value interface{}
}

func finalize(l yyLexer, value interface{}) {
	l.(*json).value = value
}
%}

%union {
  object map[string]interface{}
  member member
  array  []interface{}
  string string
  value  interface{}
}

%start json

%token <value>   TRUE FALSE NULL NUMBER PLACEHOLDER OPERATOR_SHORTCUT
%token <string>  STRING OPERATOR

%type <object>   object members
%type <member>   member
%type <array>    array elements op_params
%type <value>    json value operator

%%

json: value
                {
                  $$ = $1
                  finalize(yylex, $$)
                }

value:
    object      { $$ = $1 }
  | array       { $$ = $1 }
  | operator    { $$ = $1 }
  | STRING      { $$ = $1 }
  | NUMBER
  | TRUE
  | FALSE
  | NULL
  | PLACEHOLDER
  ;

object: '{' members '}'
                {
                  $$ = $2
                }

members:
                {
                  $$ = map[string]interface{}{}
                }
  | member
                {
                  $$ = map[string]interface{}{
                    $1.key: $1.value,
                  }
                }
  | member ',' members
                {
                  $3[$1.key] = $1.value
                  $$ = $3
                }

member: STRING ':' value
                {
                  $$ = member{
                    key:   $1,
                    value: $3,
                  }
                }

array: '[' elements ']'
                {
                  $$ = $2
                }

elements:
                {
                  $$ = []interface{}{}
                }
  | value
                {
                  $$ = []interface{}{$1}
                }
  | elements ',' value
                {
                  $$ = append($1, $3)
                }

op_params: elements
                {
                  $$ = $1
                }
  | elements ','
                {
                  $$ = $1
                }

operator:
    OPERATOR_SHORTCUT
  | OPERATOR  '(' op_params ')'
                {
                  j := yylex.(*json)
                  opPos := j.popPos()
                  op, err := j.getOperator(Operator{Name: $1, Params: $3}, opPos)
                  if err != nil {
                    j.fatal(err.Error(), opPos)
                    return 1
                  }
                  $$ = op
                }
%%

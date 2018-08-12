#!/usr/bin/env perl

# Copyright (c) 2018, Maxime Soulé
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

use strict;
use warnings;
use autodie;
use 5.010;

die "usage $0 [-h] DIR" if @ARGV == 0 or $ARGV[0] =~ /^--?h/;

my $HEADER = <<'EOH';
// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// DO NOT EDIT!!! AUTOMATICALLY GENERATED!!!
EOH

chop(my $ARGS_COMMENT = <<'EOC');
//
// "args..." are optional and allow to name the test. This name is
// logged as is in case of failure. If len(args) > 1 and the first
// item of args is a string and contains a '%' rune then fmt.Fprintf
// is used to compose the name, else args are passed to fmt.Fprint.
EOC

# These functions are variadics, but only with one possible param. In
# this case, discard the variadic property and use a default value for
# this optional parameter.
my %IGNORE_VARIADIC = (Between   => 'BoundsInIn',
		       N         => 0,
		       Re        => 'nil',
		       TruncTime => 0);

my $dir = shift;

opendir(my $dh, $dir);

my %funcs;

while (readdir $dh)
{
    if (/^td_.*\.go\z/ and not /_test.go\z/)
    {
        open(my $fh, '<', "$dir/$_");
        while (defined(my $line = <$fh>))
        {
            if ($line =~ /^func ([A-Z]\w*)\((.*?)\) TestDeep \{$/)
            {
		my $func = $1;
		if ($func ne 'Ignore')
		{
		    my @args;
		    foreach my $arg (split(/, /, $2))
		    {
			my %arg;
			@arg{qw(name type)} = split(/ /, $arg, 2);
			if ($arg{variadic} = $arg{type} =~ s/^\.{3}//)
			{
			    if (exists $IGNORE_VARIADIC{$func})
			    {
				$arg{default} = $IGNORE_VARIADIC{$func};
				delete $arg{variadic};
			    }
			}

			push(@args, \%arg);
		    }
		    $funcs{$func}{args} = \@args;
		}
            }
        }
        close $fh;
    }
}

closedir($dh);

my $funcs_contents = my $t_contents = <<EOH;
$HEADER
package testdeep

import (
\t"time"
)
EOH

foreach my $func (sort keys %funcs)
{
    my $func_name = "Cmp$func";

    my $cmp_args = 'got interface{}';
    my $call_args = '';

    foreach my $arg (@{$funcs{$func}{args}})
    {
	$call_args .= ', ' unless $call_args eq '';
	$call_args .= $arg->{name};

	$cmp_args .= ", $arg->{name} ";

	if ($arg->{variadic})
	{
	    $call_args .= '...';
	    $cmp_args .= '[]';
	}

	$cmp_args .= $arg->{type};
    }

    $funcs_contents .= <<EOF;

// Cmp$func is a shortcut for:
//
//   CmpDeeply(t, got, $func($call_args), args...)
EOF

    $t_contents .= <<EOF;

// $func is a shortcut for:
//
//   t.CmpDeeply(got, $func($call_args), args...)
EOF


    my $func_comment;
    my $last_arg = $funcs{$func}{args}[-1];
    if (exists $last_arg->{default})
    {
	$func_comment .= <<EOF;
//
// $func() optional parameter "$last_arg->{name}" is here mandatory.
// $last_arg->{default} value should be passed to mimic its absence in
// original $func() call.
EOF
    }

    $func_comment .= <<EOF;
//
// Returns true if the test is OK, false if it fails.
EOF

    $funcs_contents .= $func_comment . <<EOF;
$ARGS_COMMENT
func Cmp$func(t TestingT, $cmp_args, args ...interface{}) bool {
\tt.Helper()
\treturn CmpDeeply(t, got, $func($call_args), args...)
}
EOF

    $t_contents .= $func_comment . <<EOF;
$ARGS_COMMENT
func (t *T)$func($cmp_args, args ...interface{}) bool {
\tt.Helper()
\treturn t.CmpDeeply(got, $func($call_args), args...)
}
EOF
}

my $examples = do { open(my $efh, '<', 'example_test.go'); local $/; <$efh> };
my $funcs_reg = join('|', keys %funcs);

my($imports) = ($examples =~ /^(import \(.+?^\))$/ms);

while ($examples =~ /^func Example($funcs_reg)(_\w+)?\(\) \{\n(.*?)^\}/gms)
{
    push(@{$funcs{$1}{examples}}, { name => $2 // '', code => $3 });
}

{
    open(my $fh, "| gofmt -s > '$dir/cmp_funcs.go'");
    print $fh $funcs_contents;
    close $fh;
    say "$dir/cmp_funcs.go generated";
}

{
    open(my $fh, "| gofmt -s > '$dir/t.go'");
    print $fh $t_contents;
    close $fh;
    say "$dir/t.go generated";
}


my $funcs_test_contents = <<EOH;
$HEADER
package testdeep

$imports
EOH

my $t_test_contents = $funcs_test_contents;

my($rep, $reb, $rec);
$rep = qr/\( [^()]* (?:(??{ $rep }) [^()]* )* \)/x; # recursively matches (...)
$reb = qr/\[ [^][]* (?:(??{ $reb }) [^][]* )* \]/x; # recursively matches [...]
$rec = qr/\{ [^{}]* (?:(??{ $rec }) [^{}]* )* \}/x; # recursively matches {...}

my $rparam =qr/"(?:\\.|[^"]+)*"            # "string"
              |'(?:\\.|[^']+)*'            # 'char'
	      |`[^`]*`                     # `string`
              |&[a-zA-Z_]\w*(?:\.\w+)?(?:$rec)? # &Struct{...}, &variable
              |\[[^][]*\]\w+$rec           # []Array{...}
	      |map${reb}\w+$rec            # map[...]Type{...}
              |func\([^)]*\)[^{]+$rec      # func fn (...) ... { ... }
              |[a-zA-Z_]\w*(?:\.\w+)?(?:$rec|$rep)? # Str{...}, Fn(...), pkg.var
	      |[\w.*+-\/]+                 # 123*pkg.var...
	      |$rep$rep                    # (type)(value)
              /x;

sub extract_params
{
    my($func, $params_str) = @_;
    my $str = substr($params_str, 1, -1);

    $str ne '' or return;

    my @params;
    for (;;)
    {
	if ($str =~ /\G\s*($rparam)\s*(,|\z)/omsgx)
	{
	    push(@params, $1);
	    $2 or return @params;
	}
	else
	{
	    die "Cannot extract params from $func: $params_str\n"
	}
    }
}

foreach my $func (sort keys %funcs)
{
    my $args = $funcs{$func}{args};

    foreach my $example (@{$funcs{$func}{examples}})
    {
	my $name = $example->{name};

	foreach my $info ([ "Cmp$func(t, ", "Cmp$func", \$funcs_test_contents ],
			  [ "t.$func(",     "T_$func",  \$t_test_contents ])
	{
	    (my $code = $example->{code}) =~
		s%CmpDeeply\(t,\s+($rparam),\s+$func($rep)%
                  my @params = extract_params("$func$name", $2);
                  my $repl = $info->[0] . $1;
                  for (my $i = 0; $i < @$args; $i++)
                  {
                      $repl .= ', ';
                      if ($args->[$i]{variadic})
                      {
                          if (defined $params[$i])
                          {
                              $repl .= '[]' . $args->[$i]{type} . '{'
                                     . join(', ', @params[$i .. $#params])
                                     . '}';
                          }
                          else
                          {
                              $repl .= 'nil';
                          }
                          last
                      }
                      $repl .= $params[$i]
                          // $args->[$i]{default}
		  	   // die("Internal error, no param: "
		  		  . "$func$name -> #$i/$args->[$i]{name}!\n");
                  }
                  $repl
                  %egs;

	    ${$info->[2]} .= <<EOF;

func Example$info->[1]$name() {
$code}
EOF
	}
    }
}

{
    open(my $fh, "| gofmt -s > '$dir/cmp_funcs_test.go'");
    print $fh $funcs_test_contents;
    close $fh;
    say "$dir/cmp_funcs_test.go generated";
}

{
    $t_test_contents =~ s/t := &testing\.T\{\}/t := NewT(&testing\.T\{\})/g;
    $t_test_contents =~ s/CmpDeeply\(t,/t.CmpDeeply(/g;

    open(my $fh, "| gofmt -s > '$dir/t_test.go'");
    print $fh $t_test_contents;
    close $fh;
    say "$dir/t_test.go generated";
}

#
# Check "args..." comment is the same everywhere it needs to be
my @args_errors;
#foreach my $go_file (qw(cmp_funcs.go cmp_funcs_misc.go equal.go t.go))
foreach my $go_file (do { opendir(my $dh, $dir);
			  grep /(?<!_test)\.go\z/, readdir $dh })
{
    my $contents = do { local $/; open(my $fh, '<', "$dir/$go_file"); <$fh> };

    while ($contents =~ m,\n((?://[^\n]*\n)*)
	                    func\ ([A-Z]\w+|\(t\ \*T\)\ [A-Z]\w+)($rep),xg)
    {
	my($comment, $func, $params) = ($1, $2, $3);

	if ($params =~ /\Qargs ...interface{})\E\z/)
	{
	    chomp $comment;
	    if (substr($comment, - length($ARGS_COMMENT)) ne $ARGS_COMMENT)
	    {
		push(@args_errors, "$go_file: $func");
	    }
	}
    }
}
if (@args_errors)
{
    die "*** At least one args comment is missing or not conform:\n- "
	. join("\n- ", @args_errors)
	. "\n";
}

#$funcs_test_contents !~ /CmpDeeply/
#    or die "At least one CmpDeeply() occurrence has not been replaced!\n";

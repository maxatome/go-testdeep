#!/usr/bin/env perl

# Copyright (c) 2018-2022, Maxime Soulé
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

use strict;
use warnings;
use autodie;
use 5.010;

use IPC::Open2;

die "usage $0 [-h]\n" if @ARGV != 0;

(my $REPO_DIR = $0) =~ s,/[^/]+\z,/..,;
-d $REPO_DIR or die "Cannot find repository directory ($REPO_DIR)\n";

my $SITE_REPO_DIR = "$REPO_DIR/../go-testdeep-site";
unless (-d $SITE_REPO_DIR)
{
    if ($ENV{PROD_SITE})
    {
        die "*** Cannot PROD_SITE as $SITE_REPO_DIR not found!\n";
    }

    warn "*** WARNING: cannot find $SITE_REPO_DIR. Disabling site upgrade.\n";
    undef $SITE_REPO_DIR;
}

my $DIR = "$REPO_DIR/td";
-d $DIR or die "Cannot find td/ directory ($DIR)\n";

my $HEADER = <<'EOH';
// Copyright (c) 2018-2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// DO NOT EDIT!!! AUTOMATICALLY GENERATED!!!
EOH

my $args_comment_src = <<'EOC';

%arg{args...} are optional and allow to name the test. This name is
used in case of failure to qualify the test. If %code{len(args) > 1} and
the first item of %arg{args} is a string and contains a '%' rune then
%godoc{fmt.Fprintf} is used to compose the name, else %arg{args} are passed to
%godoc{fmt.Fprint}. Do not forget it is the name of the test, not the
reason of a potential failure.
EOC

my $ARGS_COMMENT_GD = doc2godoc($args_comment_src);
my $ARGS_COMMENT_MD = doc2md($args_comment_src);


# These functions are variadics, but only with one possible param. In
# this case, discard the variadic property and use a default value for
# this optional parameter.
my %IGNORE_VARIADIC = (Between   => 'td.BoundsInIn',
                       N         => 0,
                       Re        => 'nil',
                       TruncTime => 0,
                       # These operators accept several StructFields,
                       # but we want only one here
                       Struct    => 'nil',
                       SStruct   => 'nil');

# Smuggler operators (automatically filled)
my %SMUGGLER_OPERATORS;

# These operators should be renamed when used as *T method
my %RENAME_METHOD = (Lax => 'CmpLax');

# These operators do not have *T method nor Cmp shortcut
my %ONLY_OPERATORS = map { $_ => 1 } qw(Catch Delay Ignore Tag);

my @INPUT_LABELS = qw(nil bool str int float cplx
                      array slice map struct ptr if chan func);
my %INPUTS;
@INPUTS{@INPUT_LABELS} = ();

opendir(my $dh, $DIR);

my(%funcs, %operators, %consts, %forbiddenOpsInJSON);

while (readdir $dh)
{
    if (/^td_.*\.go\z/ and not /_test.go\z/)
    {
        my $contents = do { local $/; open(my $fh, '<', "$DIR/$_"); <$fh> };

        # Load the operators forbidden inside JSON()
        if ($_ eq 'td_json.go')
        {
            $contents =~ /^var forbiddenOpsInJSON = map\[string\]string\{(.*?)^\}/ms
                or die "$_: forbiddenOpsInJSON map not found\n";
            @forbiddenOpsInJSON{$1 =~ /"([^"]+)":/g} = ();
        }

        while ($contents =~ /^const \(\n(.+)^\)\n/gms)
        {
            @consts{$1 =~ /^\t([A-Z]\w+)/mg} = ();
        }

        my %ops;
        while ($contents =~ m,^// summary\((\w+)\): (.*\n(?://.*\n)*),gm)
        {
            my($op, $summary) = ($1, $2);
            $summary =~ s,^// input\(.*,,sm;
            $ops{$op} = process_summary($summary =~ s,\n(?://|\z),,gr);
        }

        my %inputs;
        while ($contents =~ m,^// input\((\w+)\): (.*\n(?://.*\n)*),gm)
        {
            my $op = $1;
            foreach my $in (split(/\s*,\s*/, $2 =~ s,\n(?://|\z),,gr))
            {
                if ($in eq 'all')
                {
                    @{$inputs{$op}}{keys %INPUTS} = ('✓') x keys %INPUTS;
                    next;
                }
                if ($in =~ /^(\w+)\((.*)\)\z/)
                {
                    $inputs{$op}{$1} = process_summary($2);
                    $in = $1;
                }
                else
                {
                    $inputs{$op}{$in} = '✓';
                }
                exists $INPUTS{$in} or die "$_: input($op) unknow input '$in'\n";
                $inputs{$op}{if} //= '✓'; # interface
            }
        }

        my $num_smugglers = keys %SMUGGLER_OPERATORS;

        while ($contents =~ m,^(// ([A-Z]\w*) .*\n(?://.*\n)*)func \2\((.*?)\) TestDeep \{\n,gm)
        {
            exists $ops{$2} or die "$_: no summary($2) found\n";
            exists $inputs{$2} or die "$_: no input($2) found\n";

            my($doc, $func, $params) = ($1, $2, $3);

            if ($doc =~ /is a smuggler operator/)
            {
                $SMUGGLER_OPERATORS{$func} = 1;
            }

            my @args;
            foreach my $arg (split(/, /, $params))
            {
                my %arg;
                @arg{qw(name type)} = split(/ /, $arg, 2);
                if (defined $arg{type}
                    and $arg{variadic} = $arg{type} =~ s/^\.{3}//)
                {
                    if (exists $IGNORE_VARIADIC{$func})
                    {
                        $arg{default} = $IGNORE_VARIADIC{$func};
                        delete $arg{variadic};
                    }
                }

                push(@args, \%arg);
            }
            my $last_type;
            foreach my $arg (reverse @args)
            {
                if (defined(my $arg_type = $arg->{type}) and not $arg->{variadic})
                {
                    if (defined $last_type and $arg_type eq $last_type)
                    {
                        delete $arg->{type};
                    }
                    $last_type = $arg_type;
                }
            }

	    $funcs{$func}{args} = \@args unless $ONLY_OPERATORS{$func};

            die "TAB detected in $func operator documentation\n" if $doc =~ m,(?<!^//)\t,m;

            $operators{$func} = {
                name      => $func,
                summary   => delete $ops{$func},
                input     => delete $inputs{$func},
                doc       => $doc,
                signature => "func $func($params) TestDeep",
                args      => \@args,
            };
        }

        if (%ops)
        {
            die "$_: summary found without operator definition: "
                . join(', ', keys %ops) . "\n";
        }

        if (%inputs)
        {
            die "$_: input found without operator definition: "
                . join(', ', keys %inputs) . "\n";
        }

        if ($contents =~ m,^\ttdSmugglerBase(?! // ignored),m
            and $num_smugglers == keys %SMUGGLER_OPERATORS)
        {
            die "$_: this file should contain at least one smuggler operator\n";
        }
    }
}

closedir($dh);

%funcs or die "No TestDeep golang source file found!\n";

my $funcs_contents = my $t_contents = <<EOH;
$HEADER
package td

import (
\t"time"
)

EOH

$funcs_contents .= <<EOV;
// allOperators lists the ${\scalar(keys(%funcs) + keys %ONLY_OPERATORS)} operators.
// nil means not usable in JSON().
var allOperators = map[string]any{
  ${\join('', map
              { qq("$_": ) . (exists $forbiddenOpsInJSON{$_} ? 'nil' : $_) . ",\n" }
              sort keys(%funcs), keys %ONLY_OPERATORS)}
}

EOV

my @sorted_funcs = sort { lc($a) cmp lc($b) } keys %funcs;

foreach my $func (@sorted_funcs)
{
    my $func_name = "Cmp$func";
    my $method_name = $RENAME_METHOD{$func} // $func;

    my $cmp_args = 'got';
    if (@{$funcs{$func}{args}})
    {
        foreach my $arg (@{$funcs{$func}{args}})
        {
            if (defined $arg->{type})
            {
                if ($arg->{type} ne 'any' or $arg->{variadic})
                {
                    $cmp_args .= ' any';
                }
                last
            }
        }
    }
    else
    {
        $cmp_args .= ' any';
    }

    my $call_args = '';
    my @cmpt_args;

    foreach my $arg (@{$funcs{$func}{args}})
    {
        $call_args .= ', ' unless $call_args eq '';
        $call_args .= $arg->{name};
        push(@cmpt_args, { name => $arg->{name} });

        $cmp_args .= ", $arg->{name} ";

        if ($arg->{variadic})
        {
            $call_args .= '...';
            $cmp_args .= '[]';
        }

        $cmp_args .= $arg->{type} if defined $arg->{type};
    }

    my $cmp_doc = <<EOF;
Cmp$func is a shortcut for:

\ttd.Cmp(t, got, td.$func($call_args), args...)

EOF

    $funcs_contents .= "\n" . go_comment($cmp_doc) . <<EOF;
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#$func for details.
EOF
    $cmp_doc .= <<EOF; # operator doc
See above for details.
EOF

    my $t_doc = <<EOF;
$method_name is a shortcut for:

\tt.Cmp(got, td.$func($call_args), args...)

EOF
    $t_contents .= "\n" . go_comment($t_doc) . <<EOF;
// See https://pkg.go.dev/github.com/maxatome/go-testdeep/td#$func for details.
EOF
    $t_doc .= <<EOF; # operator doc
See above for details.
EOF

    my $func_comment;
    my $last_arg = $funcs{$func}{args}[-1];
    if (exists $last_arg->{default})
    {
        $func_comment .= <<EOF;

$func() optional parameter "$last_arg->{name}" is here mandatory.
$last_arg->{default} value should be passed to mimic its absence in
original $func() call.
EOF
    }

    $func_comment .= <<EOF;

Returns true if the test is OK, false if it fails.
EOF
    my $method_comment = $func_comment;
    $func_comment .= <<EOF;

If "t" is a *T then its Config is inherited.
EOF
    $operators{$func}{cmp}{name} = "Cmp$func";
    $operators{$func}{cmp}{doc} = $cmp_doc . $func_comment . $ARGS_COMMENT_MD;
    $operators{$func}{cmp}{signature} = my $cmp_sig =
        "func Cmp$func(t TestingT, $cmp_args, args ...any) bool";
    $operators{$func}{cmp}{args} = \@cmpt_args;
    $funcs_contents .= go_comment($func_comment . $ARGS_COMMENT_GD) . <<EOF;
$cmp_sig {
\tt.Helper()
\treturn Cmp(t, got, $func($call_args), args...)
}
EOF

    $operators{$func}{t}{name} = $method_name;
    $operators{$func}{t}{doc} = $t_doc . $method_comment . $ARGS_COMMENT_MD;
    $operators{$func}{t}{signature} = my $t_sig =
        "func (t *T) $method_name($cmp_args, args ...any) bool";
    $operators{$func}{t}{args} = \@cmpt_args;
    $t_contents .= go_comment($method_comment . $ARGS_COMMENT_GD) . <<EOF;
$t_sig {
\tt.Helper()
\treturn t.Cmp(got, $func($call_args), args...)
}
EOF
}

my $examples = do { open(my $efh, '<', "$DIR/example_test.go"); local $/; <$efh> };
my $funcs_reg = join('|', @sorted_funcs);

my($imports) = ($examples =~ /^(import \(.+?^\))$/ms);

while ($examples =~ /^func Example($funcs_reg)(_\w+)?\(\) \{\n(.*?)^\}$/gms)
{
    push(@{$funcs{$1}{examples}}, { name => $2 // '', code => $3 });
}

{
    open(my $fh, "| gofmt -s > '$DIR/cmp_funcs.go'");
    print $fh $funcs_contents;
    close $fh;
    say "$DIR/cmp_funcs.go generated";
}

{
    open(my $fh, "| gofmt -s > '$DIR/t.go'");
    print $fh $t_contents;
    close $fh;
    say "$DIR/t.go generated";
}


my $funcs_test_contents = <<EOH;
$HEADER
package td_test

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
              |&?\[[^][]*\](?:any|\w+)$rec # []Array{...}
              |\[\]byte\("[^"]+"\)         # []byte("...")
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

foreach my $func (@sorted_funcs)
{
    my $args = $funcs{$func}{args};
    my $method = $RENAME_METHOD{$func} // $func;

    foreach my $example (@{$funcs{$func}{examples}})
    {
        my $name = $example->{name};

        foreach my $info ([ "td.Cmp$func(t, ", "Cmp$func", \$funcs_test_contents ],
                          [ "t.$method(",   "T_$method",\$t_test_contents ])
        {
            (my $code = $example->{code}) =~
                s%td\.Cmp\(t,\s+($rparam),\s+td\.$func($rep)%
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
    # Cmp* examples
    open(my $fh, "| gofmt -s > '$DIR/example_cmp_test.go'");
    print $fh $funcs_test_contents;
    close $fh;
    say "$DIR/example_cmp_test.go generated";
}

{
    # T.* examples
    $t_test_contents =~ s/t := &testing\.T\{\}/t := td.NewT(&testing\.T\{\})/g;
    $t_test_contents =~ s/td\.Cmp\(t,/t.Cmp(/g;

    open(my $fh, "| gofmt -s > '$DIR/example_t_test.go'");
    print $fh $t_test_contents;
    close $fh;
    say "$DIR/example_t_test.go generated";
}

#
# Check "args..." comment is the same everywhere it needs to be
my @args_errors;
$ARGS_COMMENT_GD = go_comment($ARGS_COMMENT_GD);
foreach my $go_file (do { opendir(my $dh, $DIR);
                          grep /(?<!_test)\.go\z/, readdir $dh })
{
    my $contents = do { local $/; open(my $fh, '<', "$DIR/$go_file"); <$fh> };

    while ($contents =~ m,\n((?://[^\n]*\n)*)
                            func\ ([A-Z]\w+|\(t\ \*T\)\ [A-Z]\w+)($rep),xg)
    {
        my($comment, $func, $params) = ($1, $2, $3);

        next if ($func eq '(t *T) CmpDeeply'
                 or $func eq 'CmpDeeply'
                 or $func =~ /^\(t \*T\) (?:Log|Error|Fatal)Trace\z/);

        if ($params =~ /\Qargs ...any)\E\z/)
        {
            #chomp $comment;
            if (substr($comment, - length($ARGS_COMMENT_GD)) ne $ARGS_COMMENT_GD)
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

my $common_links = do
{
    my $base_url = 'https://pkg.go.dev/github.com/maxatome/go-testdeep';
    my $td_url = "$base_url/td";

    # Specific types and functions
    join("\n", map "[`$_`]: $td_url#$_", qw(T TestDeep Cmp))
        . "\n\n"
        # Helpers
        . join("\n", map "[`$_`]: $base_url/helpers/$_",
               qw(tdhttp tdsuite tdutil))
        . "\n\n"
        # Specific links
        . "[`BeLax` config flag]: $td_url#ContextConfig\n"
        . "[`error`]: https://pkg.go.dev/builtin/#error\n"
        . "\n\n"
        # Foreign types
        . join("\n", map "[`$_->[0]`]: https://pkg.go.dev/$_->[1]",
               [ 'fmt.Stringer' => 'fmt/#Stringer' ],
               [ 'time.Time'    => 'time/#Time' ],
               [ 'math.NaN'     => 'math/#NaN' ])
        . "\n";
};

my @sorted_operators = sort { lc($a) cmp lc($b) } keys %operators;

my $md_links = do
{
    $common_links
        . join("\n", map qq([`$_`]: {{< ref "$_" >}}), @sorted_operators)
        . "\n\n"
        # Cmp* functions
        . join("\n", map qq([`Cmp$_`]: {{< ref "$_#cmp\L$_\E-shortcut" >}}),
                     @sorted_funcs)
        . "\n\n"
        # T.Cmp* methods
        . join("\n", map
               {
                   my $m = $RENAME_METHOD{$_} // $_;
                   qq([`T.$m`]: {{< ref "$_#t\L$m\E-shortcut" >}})
               }
               @sorted_funcs);
};

my $gh_links = do
{
    my $td_url = 'https://go-testdeep.zetta.rocks/operators/';
    $common_links
        . join("\n", map qq([`$_`]: $td_url\L$_/), @sorted_operators)
        . "\n\n"
        # Cmp* functions
        . join("\n", map qq([`Cmp$_`]: $td_url\L$_/#cmp$_-shortcut), @sorted_funcs)
        . "\n\n"
        # T.Cmp* methods
        . join("\n", map
               {
                   my $m = $RENAME_METHOD{$_} // $_;
                   qq([`T.$m`]: $td_url\L$_/#t$m-shortcut)
               }
               @sorted_funcs);
};

# README.md
{
    my $readme = do { local $/; open(my $fh, '<', "$REPO_DIR/README.md"); <$fh> };

    # Links
    $readme =~ s{(<!-- links:begin -->).*(<!-- links:end -->)}
                {$1\n$gh_links\n$2}s;

    open(my $fh, '>', "$REPO_DIR/README.md.new");
    print $fh $readme;
    close $fh;
    rename "$REPO_DIR/README.md.new", "$REPO_DIR/README.md";
    say "$REPO_DIR/README.md modified";
}

# Hugo
if (defined $SITE_REPO_DIR)
{
    my $op_examples = do { local $/;
                           open(my $fh, '<', "$DIR/example_test.go");
                           <$fh> };

    # Reload generated examples so they are properly gofmt'ed
    my $cmp_examples = do { local $/;
                            open(my $fh, '<', "$DIR/example_cmp_test.go");
                            <$fh> };
    my $t_examples = do { local $/;
                          open(my $fh, '<', "$DIR/example_t_test.go");
                          <$fh> };

    foreach my $operator (@sorted_operators)
    {
        # Rework each operator doc
        my $doc = process_doc($operators{$operator});

        open(my $fh, '>', "$SITE_REPO_DIR/docs_src/content/operators/$operator.md");
        print $fh <<EOM;
---
title: "$operator"
weight: 10
---

```go
$operators{$operator}{signature}
```

$doc

> See also [<i class='fas fa-book'></i> $operator godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#$operator).

EOM

        my @examples;
        my $re = qr/^func Example${operator}(?:_(\w+))?\(\) \{\n(.+?)^\}$/ms;
        while ($op_examples =~ /$re/g)
        {
            my $name = ucfirst($1 // 'Base');

            push(@examples, <<EOE);
{{%expand "$name example" %}}```go
${2}
```{{% /expand%}}
EOE
        }
        if (@examples)
        {
            printf $fh "### Example%s\n\n", @examples > 1 ? 's' : '';
            print $fh @examples;
        }

        if (my $cmp = $operators{$operator}{cmp})
        {
            my $doc = process_doc($cmp);
            print $fh <<EOM;
## $cmp->{name} shortcut

```go
$cmp->{signature}
```

$doc

> See also [<i class='fas fa-book'></i> $cmp->{name} godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#$cmp->{name}).

EOM

            @examples = ();
            my $re = qr/func ExampleCmp${operator}(?:_(\w+))?\(\) \{\n(.+?)^\}$/ms;
            while ($cmp_examples =~ /$re/g)
            {
                my $name = ucfirst($1 // 'Base');

                push(@examples, <<EOE);
{{%expand "$name example" %}}```go
${2}
```{{% /expand%}}
EOE
            }
            if (@examples)
            {
                printf $fh "### Example%s\n\n", @examples > 1 ? 's' : '';
                print $fh @examples;
            }
        }

        if (my $t = $operators{$operator}{t})
        {
            my $doc = process_doc($t);
            print $fh <<EOM;
## T.$t->{name} shortcut

```go
$t->{signature}
```

$doc

> See also [<i class='fas fa-book'></i> T.$t->{name} godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.$t->{name}).

EOM

            @examples = ();
            my $re = qr/func ExampleT_$t->{name}(?:_(\w+))?\(\) \{\n(.+?)^\}$/ms;
            while ($t_examples =~ /$re/g)
            {
                my $name = ucfirst($1 // 'Base');

                push(@examples, <<EOE);
{{%expand "$name example" %}}```go
${2}
```{{% /expand%}}
EOE
            }
            if (@examples)
            {
                printf $fh "### Example%s\n\n", @examples > 1 ? 's' : '';
                print $fh @examples;
            }
        }
        close $fh;
    }

    # Dump operators
    {
        my $op_list_file = "$SITE_REPO_DIR/docs_src/content/operators/_index.md";
        my $op_list = do { local $/; open(my $fh, '<', $op_list_file); <$fh> };

        $op_list =~ s{(<!-- operators:begin -->).*(<!-- operators:end -->)}
                     {
                         "$1\n"
                             . join('',
                                    map qq![`$_`]({{< ref "$_" >}})\n: $operators{$_}{summary}\n\n!,
                                    @sorted_operators)
                             . $2
                     }se or die "operators tags not found in $op_list_file\n";

        $op_list =~ s{(<!-- smugglers:begin -->).*(<!-- smugglers:end -->)}
                     {
                         "$1\n"
                             . join('',
                                    map qq![`$_`]({{< ref "$_" >}})\n: $operators{$_}{summary}\n\n!,
                                    sort { lc($a) cmp lc($b) }
                                    keys %SMUGGLER_OPERATORS)
                             . "$md_links\n$2"
                     }se or die "smugglers tags not found in $op_list_file\n";

        open(my $fh, '>', "$op_list_file.new");
        print $fh $op_list;
        close $fh;
        rename "$op_list_file.new", $op_list_file;
    }

    # Dump matrices
    {
        my $matrix_file = "$SITE_REPO_DIR/docs_src/content/operators/matrix.md";
        my $matrix = do { local $/; open(my $fh, '<', $matrix_file); <$fh> };

        my $header = <<'EOH';

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
EOH

        $matrix =~ s{(<!-- op-go-matrix:begin -->).*(<!-- op-go-matrix:end -->)}
                    {
                        my $repl = "$1\n";
                        my $num = 0;
                        foreach my $op (@sorted_operators)
                        {
                            $repl .= $header if $num++ % 10 == 0;
                            $repl .= "| [`$op`]";
                            for my $label (@INPUT_LABELS)
                            {
                                $repl .= " | "
                                    . ($operators{$op}{input}{$label} // '✗');
                            }
                            $repl .= " | [`$op`] |\n";
                        }
                        "$repl\n$md_links\n$2"
                    }se or die "op-go-matrix tags not found in $matrix_file\n";

        my %by_input;
        while (my($op, $info) = each %operators)
        {
            while (my($label, $comment) = each %{$operators{$op}{input}})
            {
                $by_input{$label}{$op} = $comment;
            }
        }
        $matrix =~ s{(<!-- go-(\w+)-matrix:begin -->).*(<!-- go-\2-matrix:end -->)}
                    {
                        my $repl = "$1\n";
                        foreach my $op (sort keys %{$by_input{$2}})
                        {
                            my $comment = $by_input{$2}{$op};
                            next if $comment eq 'todo';
                            if ($comment eq '✓')
                            {
                                next if $2 eq 'if';
                                $comment = '';
                            }
                            elsif ($2 eq 'if')
                            {
                                $comment =~ s/^✓ \+/ →/;
                            }
                            else
                            {
                                substr($comment, 0, 0, ' only ');
                            }
                            $repl .= "- [`$op`]$comment\n";
                        }
                        $repl . $3
                    }gse or die "go-op-matrix tags not found in $matrix_file\n";

        open(my $fh, '>', "$matrix_file.new");
        print $fh $matrix;
        close $fh;
        rename "$matrix_file.new", $matrix_file;
    }

    # tdhttp example
    {
        my $example = do
        {
            open(my $fh, '<', "$REPO_DIR/helpers/tdhttp/example_test.go");
            local $/;
            <$fh>
        };

        my($import) = $example =~ /^(import \(.*?^\))$/ms;
        $import or die "tdhttp example, import not found!\n";
        $example =~ s/.*^func Example\(\) \{\n\tt := &testing.T\{\}\n\n//ms
                                                                                            or die "tdhttp example, func Example() not found!\n";

        $example =~ s/fmt\.Printf/t.Logf/g
            or die "tdhttp example, fmt.Printf not found\n";
        $example =~ s/fmt\.Println/t.Log/g
            or die "tdhttp example, fmt.Println not found\n";
        $example =~ s,\n\t// Output:\n.*,},s
            or die "tdhttp example, Output: not found\n";

        my $md_file = "$SITE_REPO_DIR/docs_src/content/helpers/_index.md";

        my $final = do { open(my $fh, '<', $md_file); local $/; <$fh> } =~
            s{(<!-- tdhttp:begin -->).*(<!-- tdhttp:end -->)}
             <$1
{{%expand "Main example" %}}```go
package myapi

$import

func TestMyAPI(t *testing.T) {
$example
```{{% /expand%}}
$2>rs or die "tdhttp example not found in $md_file!";

        open(my $out, '>', "$md_file.new");
        print $out $final;
        close $out;
        rename "$md_file.new", $md_file;
    }

    # Final publish
    if ($ENV{PROD_SITE})
    {
        # Delegate to go-testdeep-site repository
        chdir $SITE_REPO_DIR;
        exec './publish.sh';
    }
}


# ""      → "//"
# "\txxx" → "//\txxx"
# "xxx"   → "// xxx"
sub go_comment
{
    shift =~ s{^(.?)}
              {
                  $1 eq ''
                  ? '//'
                  : substr($1, 0, 1) eq "\t" ? "//$1" : "// $1"
              }egmr
}

sub doc2godoc
{
    my $doc = shift;

    state $repl = { arg   => sub { qq("$_[0]") },
                    code  => sub { $_[0] },
                    godoc => sub { $_[0] } };

    $doc =~ s/%([a-z]+)\{([^}]+)\}/($repl->{$1} or die $1)->($2)/egr;
}

sub doc2md
{
    my $doc = shift;

    state $repl = { arg   => sub { "*$_[0]*" },
                    code  => sub { "`$_[0]`" },
                    godoc => sub
                    {
                        my($pkg, $fn) = split('\.', $_[0], 2);
                        "[`$_[0]`](https://pkg.go.dev/$pkg/#$fn)"
                    } };

    $doc =~ s/%([a-z]+)\{([^}]+)\}/($repl->{$1} or die $1)->($2)/egr;
}

sub process_summary
{
    my $text = shift;

    $text =~ s/(time\.Time|fmt.Stringer|error)/[`$1`]/g;
    $text =~ s/BeLax config flag/[`BeLax` config flag]/;
    $text =~ s/(\[\]byte|\bnil\b)/`$1`/g;

    return $text;
}

sub process_doc
{
    my $op = shift;

    my $doc = $op->{doc};

    $doc =~ s,^// ?,,mg if $doc =~ m,^//,;

    $doc =~ s/\n{3,}/\n\n/g;

    my($inEx, $inBul);
    $doc =~ s{^(?:(\n?\S)
                 |(\n?)(\s+)(\S+))}
             <
                if (defined $1)
                {
                    if ($inEx)     { $inEx = ''; "```\n$1" }
                    elsif ($inBul) { $inBul = ''; "\n$1" }
                    else { $1 }
                }
                else
                {
                    my($nl, $indent, $beg) = ($2, $3, $4);
                    if ($inEx) { $nl . substr($indent, length($inEx)) . $beg }
                    elsif ($inBul) { $nl . substr($indent, length($inBul)) . $beg }
                    elsif ($beg =~ /^---/) { $inEx = $indent; "$nl```\n$beg" }
                    elsif ($beg =~ /^-/)   { $inBul = $indent; "\n-" }
                    else                   { $inEx = $indent; "$nl```go\n$beg" }
                }
             >gemx;
    $doc .= "```\n" if $inEx;

    my @codes;
    $doc =~ s/^(```go\n.*?^```\n)/push(@codes, $1); "CODE<$#codes>"/gems;

    $doc =~ s<
        (\$\^[A-Za-z]+)                    # $1
      | (\b(${\join('|', grep !/^JSON/, keys %operators)}
           |JSON(?!\s+(?:value|data|filename|object|representation|specification)))
        (?:\([^)]*\)|\b))                  # $2 $3
      | ((?:(?:\[\])+|\*+|\b)(?:bool\b
                               |u?int(?:\*|(?:8|16|32|64)?\b)
                               |float(?:\*|(?:32|64)\b)
                               |complex(?:\*|(?:64|128)\b)
                               |string\b
                               |rune\b
                               |byte\b
                               |any)
        |\(\*byte\)\(nil\)
        |\bmap\[string\]any
        |\b(?:len|cap)\(\)
        |\bnil\b
        |\$(?:\d+|[a-zA-Z_]\w*))           # $4
      | ((?:\b|\*)fmt\.Stringer
        |\breflect\.Type
        |\bregexp\.MustCompile
        |\*regexp\.Regexp
        |\btime\.[A-Z][a-zA-Z]+
        |\bjson\.(?:Unm|M)arshal
        |\bio\.Reader
        |\bioutil\.Read(?:All|File))\b     # $5
      | (\berror\b)                        # $6
      | (\bTypeBehind(?:\(\)|\b))          # $7
      | \b(${\join('|', keys %consts)})\b  # $8
      | \b(smuggler\s+operator)\b          # $9
      | \b(TestDeep\s+operators?)\b        # $10
      | (?|\*?(SmuggledGot)|(Flatten))\b   # $11
      >{
           if ($1)
           {
               "`$1`"
           }
           elsif ($2)
           {
               qq![`$2`]({{< ref "$3" >}})!
           }
           elsif ($4)
           {
               "`$4`"
           }
           elsif ($5)
           {
               my $all = $5;
               my($pkg, $fn) = split('\.', $all, 2);
               $pkg =~ s/^\*//;
               "[`$all`](https://pkg.go.dev/$pkg/#$fn)"
           }
           elsif ($6)
           {
               "[`$6`](https://pkg.go.dev/builtin/#error)"
           }
           elsif ($7)
           {
               qq![`$7`]({{< ref "operators#typebehind-method" >}})!
           }
           elsif ($8)
           {
               "[`$8`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsKind)"
           }
           elsif ($9)
           {
               qq![$9]({{< ref "operators#smuggler-operators" >}})!
           }
           elsif ($10)
           {
               qq![$10]({{< ref "operators" >}})!
           }
           elsif ($11)
           {
               qq![`$11`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#$11)!
           }
       }geox;

    if ($op->{args} and @{$op->{args}})
    {
        $doc =~ s/"(${\join('|', map quotemeta($_->{name}),
                            @{$op->{args}})})"/*$1*/g;
    }

    return $doc =~ s/^CODE<(\d+)>/go_format($op, $codes[$1])/egmr;
}

sub go_format
{
    my($operator, $code) = @_;

    $code =~ s/^```go\n// or return $code;
    $code =~ s/\n```\n\z//;

    my $pid = open2(my $fmt_out, my $fmt_in, 'gofmt', '-s');

    my $root;
    if ($code =~ /^func/)
    {
        $root = 1;
        print $fmt_in <<EOM;
package x

//line $operator->{name}.go:1
$code
EOM
    }
    else
    {
        print $fmt_in <<EOM;
package x

//line $operator->{name}.go:1
func x() {
$code
}
EOM
    }
    close $fmt_in;

    my $new_code = do { <$fmt_out>; <$fmt_out>; <$fmt_out>; # skip 1st 3 lines
                        local $/; <$fmt_out> };
    chomp($new_code);
    unless ($root)
    {
        $new_code =~ s/[^\t]+//;
        $new_code =~ s/\n\}\z//;
        $new_code =~ s/^\t//gm;
    }

    waitpid $pid, 0;
    if ($? != 0)
    {
        die <<EOD
gofmt of following example for function $operator->{name} failed:
$code
EOD
    }

    $new_code =~ s/^(\t+)/"  " x length $1/gme;

    if ($new_code ne $code)
    {
        die <<EOD;
Code example function $operator->{name} is not correctly indented:
$code
------------------ should be ------------------
$new_code
EOD
    }

    return "```go\n$new_code\n```\n";
}

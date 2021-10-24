#!/usr/bin/env perl

# Copyright (c) 2018, Maxime Soulé
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

use strict;
use warnings;
use autodie;
use 5.014;

use HTTP::Tiny;

my $SPEW_BASE_URL =
    'https://raw.githubusercontent.com/davecgh/go-spew/master/spew/';

foreach my $file (qw(bypass.go bypasssafe.go))
{
    my $resp = HTTP::Tiny::->new->get("$SPEW_BASE_URL$file");
    $resp->{success} or die "Failed to retrieve $file!\n";

    unless ($resp->{content} =~ s/^package \Kspew$/dark/m)
    {
        die "'package spew' line not found in $file!\n";
    }

    open(my $fh, '>', "internal/dark/$file");

    say $fh <<EOH;
// DO NOT EDIT!!! AUTOMATICALLY COPIED FROM
// https://github.com/davecgh/go-spew/blob/master/spew/$file
EOH

    my %ops = (',' => ' && ', ' ' => ' || ');
    $resp->{content} =~
        s{^(?=// \+build (.*))}
         {"//go:build " . ($1 =~ s!([ ,])!$ops{$1}!gr) . "\n"}em;
    print $fh $resp->{content};

    close $fh;
}

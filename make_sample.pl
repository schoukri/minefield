#!/usr/bin/perl

use strict;
use warnings;

use Getopt::Long;

GetOptions(
    'num=i' => \my $num,
    'quad=i' => \my $quad,
    ) || die "can't get options: $!";

$num //= 1000;
$quad //= 10;

sub rand_coord {
    my $n = rand($quad);
    if (int(rand(2)) % 2 == 0) {
	$n *= -1;
    }
    return sprintf("%.6f", $n);
}

sub rand_power {
    return sprintf("%.6f", rand(2) * rand(2))
}

for (1..$num) {
    print join("\t", rand_coord(), rand_coord(), rand_power()), "\n";
}


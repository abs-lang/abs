#!/usr/bin/env bats

## These tests are run by Bash Automated Testing System (BATS)
## ref. https://github.com/bats-core/bats-core
## ref. https://testanything.org/
## ref. https://opensource.com/article/19/2/testing-bash-bats
##
## Notes:
## 1) Existing ABS scripts in tests/*.abs do not need to be altered for BATS
## 2) You can see error case printf output by changing [status -eq ??] tests
## 3) See Test Anything Protocol (TAP) output using $ bats --tap <testFile>
## 4) Tests are written using bash [ test ] and $(expr) syntax.
## 5) $var strings must be quoted. E.g. "$output" or "${lines[0]}"
##
## Install BATS 
## cd ~/git
## git clone https://github.com/bats-core/bats-core.git
## cd ~/git/bats-core
## ./install.sh ~/bats
## sudo ln -sf ~/bats/bin/bats /usr/local/bin/bats
##
## Test ABS with BATS
## cd ~/git/abs
## bats tests/test-abs.bats

@test "abs --version prints version number" {
    run abs --version
    printf '# VERSION: %s\n' "$output" >&2
    [ $status -eq 0 ]
    [ $(expr "$output" : "[0-9][0-9.][0-9.]*") -ne 0 ]
}

@test "Load ABS Init File for tests/test-builtins.abs" {
    # setup
    mkdir -p $BATS_TMPDIR/abs
    ABS_INIT_FILE=$BATS_TMPDIR/abs/.absrc
    cp ~/git/abs/tests/test-absrc.abs $ABS_INIT_FILE
    export ABS_INIT_FILE
    # run test
    run abs tests/test-builtins.abs
    # teardown
    rm -rf $BATS_TMPDIR/abs
    # test results
    printf 'ABS_INIT_FILE: %s\n' $ABS_INIT_FILE
    printf '# %s\n' "$output" >&2
    [ $status -eq 0 ]
    [ "${lines[0]}" = "ABS_INTERACTIVE: false" ]
    [ "${lines[2]}" = "chdir path/to/nowhere: no such file or directory" ]
}

@test "tests/test-parser.abs" {
    run abs tests/test-parser.abs
    printf '# %s\n' "$output" >&2
    # parser errors are expected
    [ $status -eq 99 ]
    # match one or more parser error lines
    [ "${lines[5]}" = "	no prefix parse function for '%' found" ]
    [ "${lines[6]}" = "	[13:4]	b %% c" ]
}

@test "tests/test-eval.abs" {
    for i in 1 2 3; do
        run abs tests/test-eval.abs $i
        printf '# %s\n' "${lines[0]}" >&2
        printf '# %s\n' "${lines[1]}" >&2
        # eval errors are expected
        [ $status -eq 99 ]
        case $i in 
        1)  [ "${lines[0]}" = "ERROR: type mismatch: STRING + NUMBER" ]
            [ "${lines[1]}" = "	[8:11]	    s = s + 1   # this is a comment" ]
            ;;
        2)  [ "${lines[0]}" = "ERROR: invalid property 'junk' on type ARRAY" ]
            [ "${lines[1]}" = "	[14:6]	    a.junk" ]
            ;;
        3)  [ "${lines[0]}" = "ERROR: index operator not supported: f(x) {x} on HASH" ]
            [ "${lines[1]}" = '	[19:20]	    {"name": "Abs"}[f(x) {x}];  ' ]
            ;;
        esac
    done
}

@test "tests/test-hash-funcs.abs" {
    run abs tests/test-hash-funcs.abs
    printf '# %s\n' "$output" >&2
    # no errors are expected
    [ $status -eq 0 ]
    # match abs expression results only, ignore other lines
    [ "${lines[3]}" = "{a: 1, b: 2, c: {x: 10, y: 20}}" ]
    [ "${lines[5]}" = "[a, b, c]" ]
    [ "${lines[7]}" = "[a, b, c]" ]
    [ "${lines[17]}" = "{a: 1}" ]
    [ "${lines[19]}" = "{b: 2, c: {x: 10, y: 20}}" ]
    [ "${lines[21]}" = "{c: {x: 10, y: 20}}" ]
    [ "${lines[23]}" = "{b: 2}" ]
    [ "${lines[25]}" = "null" ]
    [ "${lines[27]}" = "{b: 2}" ]
    # need to use a regex with hash.values() or hash.items() because they may present in any order
    [ $(expr "${lines[9]}" : ".*{x: 10, y: 20}.*") -ne 0 ]
    [ $(expr "${lines[11]}" : ".*{x: 10, y: 20}.*") -ne 0 ]
    [ $(expr "${lines[13]}" : ".*{x: 10, y: 20}.*") -ne 0 ]
    [ $(expr "${lines[15]}" : ".*{x: 10, y: 20}.*") -ne 0 ]
 }



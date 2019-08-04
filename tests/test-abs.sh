#! /bin/bash
# test the abs parser and eval error location

DIR=~/go/src/github.com/abs-lang/abs
ABS=`command -v abs`

DEBUG=$1

if [[ "$DEBUG" == "-d" ]]; then
    ABS=$DIR/builds/abs
fi
if [ -z $ABS ]; then
    echo "Cannot locate abs binary; exiting"
    exit 1
fi

cd $DIR
LINE="======================================="
echo $LINE
echo "Test Parser"
FILE=tests/test-parser.abs
echo $FILE
$ABS $FILE
echo "Exit code: $?"
echo

echo $LINE
echo "Test Eval()"
FILE=tests/test-eval.abs
echo $FILE

for i in 1 2 3; do
    $ABS $FILE $i
    echo "Exit code: $?"
    echo
done

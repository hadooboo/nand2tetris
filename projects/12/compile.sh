#!/bin/bash

class=("Array" "Keyboard" "Math" "Memory" "Output" "Screen" "String" "Sys")

if [[ $# -ne 1 || ! " ${class[*]} " =~ " $1 " ]]; then
    echo "Usage: $0 [${class[*]}]"
    exit
fi

testDir="$1Test"
testFile="$1.jack"
tempDir="tmp_$testDir"
cp -r "$testDir" "$tempDir"
for item in "${class[@]}"; do
    cp "../../tools/OS/$item.vm" "$tempDir"
done
cp "$testFile" "$tempDir"
../../tools/JackCompiler.sh "$tempDir"

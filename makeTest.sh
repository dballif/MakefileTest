#!/bin/bash

#Args: Makefile to be tested

#Variables that can be changed
#outDir - the directory where outputs will be stored
outDir=makeTestOut
#threads - the number of threads Make will use. Not all systems can handle 8,
#   but if possible, it is useful to run it with a high number to flush out any
#   dependency issues. This will sometimes break other programs so beware.
threads=8

#Goal: To test the makefile to make sure it is working
#Level 1: All Targets in the Makefile are runnable and output no errors
#Level 2: Standard targets do correct things (Correct files are created or deleted)
#Level 3: Non standard, but full targets perform there correct functions
#Level 4: Intermediate targets all do the correct things

#Boilerplate
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

#Level 1:
# Get list of dependencies using Grep (Will probably require some cool regex)
# Run them all
# On error, exit and explode

#cat $1 | awk -F':' '/^[a-zA-Z0-9][^$#\\\t=]*:([^=]|$)/ {split($1,A,/ /);for(i in A)print A[i]}'

LIST=`cat $1 | awk -F':' '/^[a-zA-Z0-9][^$#\/\t=]*:([^=]|$)/ {split($1,A,/ /);for(i in A)print A[i]}' | sort -u | grep -v Makefile | grep -v local.mk`

#Print list string
#echo ($LIST)

#Convert list string to array
array=($LIST)

#Print all Targets
# for element in "${array[@]}"
# do
#     echo "$element"
# done

#Create directory for output files
mkdir -p $outDir

#Run each Target to make sure it works
for element in "${array[@]}"
do
	#Save any error output into a variable
	echo "[RUNNING TARGET: ${element}]"
    make $element -j$threads > $outDir/$element.output 2>&1
	#Check for any "***" which make uses to display errors
	testResult=$(cat $outDir/$element.output | grep "Error" | grep "Stop." | grep -v "ignore")
	#Check if the target had any errors
	if [ -z "$testResult" ]
	then
		echo -e "$element: ${GREEN}PASS${NC}"
	else
		echo -e "$element: ${RED}FAIL${NC}"
		echo "See $outDir/$element.output for log"
		overallResult=fail
	fi
done

echo

if [ -z "$overallResult" ]
then
	echo -e "Overall: ${GREEN}Pass${NC}"
	echo "All targets ran without generating an error"
else
	echo -e "Overall: ${RED}Fail${NC}"
fi
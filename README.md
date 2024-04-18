# MakefileTest
Originally a Bash Script that would just run all the targets in a Makefile to see if any failed, this project has evolved into a Go program that will allow for a more modular approach to testing a Makefile.

## Features
 - A program that can take a Makefile, and runs all the targets.
 - Takes a json file containing the targets to be run and the desired output/file for each one
 - JSON Params: 
      - Name (string)
      - targetToRun (string)
      - filesCreated (string)
      - filesDeleted (string)
      - searchForFailureInOutput (bool)
 
 Note: More features will probably be added later, this is going to be the basic functionality though.

 ## JSON Format
 
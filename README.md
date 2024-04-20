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

## Future Features
A -g flag that parses the Makefile (given by -m) and creates a JSON accordingly. It will gather a list of all the targets and each target will have a test created. Name will just be "<Target> Test". If the test has a phony, the fileCreated test will be left empty. If "rm" or some version of that is found, it will try to add a deleted file (should I look for clean targets too?). searchForFailureInOutput will be set to true. This will be a method that is run up front and will cause the application to return so no testing is actually run.

 ## JSON Format
 
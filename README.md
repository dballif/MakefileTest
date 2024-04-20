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

## Flags
| Flag  | Description                                             |
| ----- | ------------------------------------------------------- |
|  -h   | Print the help message                                  |
|  -m   | The path to the Makefile to be tested                   |
|  -f   | The JSON File containing information regarding the test |
|  -g   | Generate a JSON from the Makefile                       |

## Automatic JSON Generation (-g)
This feature is still  undergoing testing.  

The program will parse the Makefile specified by "-m" and attempt to create a JSON file as a starting basis for a test file. It will create a test for each target within the Makefile with the following features:
- Name will just be "<Target> Test"
- If the test has a phony, the fileCreated test will be left empty
- If "rm" or some version of that is found, it will try to add a deleted file
- searchForFailureInOutput will be set to true by default
     
**Note:** This will be a method that is run up front and will cause the application to return so no testing is actually run

## Future Features
- Change the output failure parsing to just look for any stderr instead of just parsing for "fail" keyword

 ## JSON Format - Example
 ```
 {
    "testTargets": [
        {
            "name": "Test Target 1",
            "target": "target1",
            "filesCreated": "target1test,target2test",
            "filesDeleted": "",
            "ignoreFailure": true
        },
        {
            "name": "Test Target 2",
            "target": "target1Clean",
            "filesCreated": "",
            "filesDeleted": "target1test,target2test",
            "ignoreFailure": false
        },
        {
            "name": "Test Fail",
            "target": "failTarget",
            "filesCreated": "",
            "filesDeleted": "",
            "ignoreFailure": false
        }
    ]
}
```
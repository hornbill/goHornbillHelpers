# Hornbill Helpers

A Go module containing a number of helper functions that are used by Hornbill Go open source tools.

## CalculateTimeDuration
Calculates a new date/time when giving a starting point and a duration.

### Input Parameters
1. **time.Time**: A starting date & time
2. **string**: a duration string in the following format:
- For adding a period of time from the start time: **P1D2H3M4S** This will add 1 day (1D), 2 hours (2H), 3 minutes (3H) and 4 seconds (4S) to the provided time
- For subtracting a period of time from the start time: **-P1D2H3M4S** - This will subtract 1 day (1D), 2 hours (2H), 3 minutes (3H) and 4 seconds (4S) from the provided time

### Output Parameters
1. **time.Time**: The calculated date & time
2. **int**: The number of seconds between the starting and calculated date/times.

## ConvFloatToStorage
Takes a float64 value, returns a human readable storage capacity (kB, MB, GB, TB, PB) string.

### Input Parameters
1. **float64**: Float to convert in to a capacity string

### Output Parameters
1. **string**: Formatted storage capacity string

## Logger
Create or append to a log file in to the **/log** folder of the folder where the Go code is executed.

### Input Parameters
1. **int**: Log entry type:
- **1**: [DEBUG] CLI output GREEN
- **2**: [MESSAGE] CLI output GREEN
- **3**: No type prefix, CLI output GREEN
- **4**: [ERROR] CLI output RED
- **5**: [WARNING] CLI output YELLOW
- **6**: No type prefix, CLI output YELLOW
2. **string**: Log entry string
3. **bool**: Should the log entry also be output to the CLI
4. **string**: Log file name

##ConfirmResponse
CLI prompt to user, expects a user response of:
- Fuzzy yes (y, Y, yes, Yes, YES) or input param provided string
- Fuzzy no (n, N, no, No, NO)
Function does not return until an expected response is given.

### Input Parameters
1. **string**: An overriding replacement for the fuzzy yes expected response
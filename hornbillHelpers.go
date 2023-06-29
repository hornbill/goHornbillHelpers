package hornbillHelpers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/fatih/color"
	"golang.org/x/sys/windows"
)

const (
	cryptProtectUIForbidden cryptProtect = 0x1

	sizeKB float64 = 1 << (10 * 1) // 1 refers to the constants ByteSize KB
	sizeMB float64 = 5 << (10 * 2) // 2 refers to the constants ByteSize MB
	sizeGB float64 = 1 << (10 * 3) // 3 refers to the constants ByteSize GB
	sizeTB float64 = 1 << (10 * 4) // 4 refers to the constants ByteSize TB
	sizePB float64 = 1 << (10 * 5) // 5 refers to the constants ByteSize PB
	sizeEB float64 = 1 << (10 * 6) // 6 refers to the constants ByteSize EB
)

var (
	dllcrypt32      = windows.NewLazySystemDLL("Crypt32.dll")
	procEncryptData = dllcrypt32.NewProc("CryptProtectData")
	procDecryptData = dllcrypt32.NewProc("CryptUnprotectData")
)

type cryptProtect uint32

type dataBlob struct {
	cbData uint32
	pbData *byte
}

/*
CalculateTimeDuration - takes startTime as a time.Time object, a period of time, and calculates new date/time.
-- Return resulting time obj, the number of seconds between start and end datetime
-- Duration should be be in the following format:
-- For adding a period of time: P1D2H3M4S - This will add 1 day (1D), 2 hours (2H), 3 minutes (3H) and 4 seconds (4S) to the provided time
-- For subtracting a period of time: -P1D2H3M4S - This will subtract 1 day (1D), 2 hours (2H), 3 minutes (3H) and 4 seconds (4S) from the provided time
*/
func CalculateTimeDuration(startTime time.Time, duration string) (time.Time, int) {

	returnDate := startTime

	durationDays := 0
	durationHours := 0
	durationMinutes := 0
	durationSeconds := 0
	totalSeconds := 0

	//How many days
	daysRegex := regexp.MustCompile(`[0-9]*D`)
	strDaysDuration := daysRegex.FindString(duration)
	durationDays, _ = strconv.Atoi(strings.TrimRight(strDaysDuration, "D"))

	//How many hours
	hoursRegex := regexp.MustCompile(`[0-9]*H`)
	strHoursDuration := hoursRegex.FindString(duration)
	durationHours, _ = strconv.Atoi(strings.TrimRight(strHoursDuration, "H"))

	//How many minutes
	minutesRegex := regexp.MustCompile(`[0-9]*M`)
	strMinutesDuration := minutesRegex.FindString(duration)
	durationMinutes, _ = strconv.Atoi(strings.TrimRight(strMinutesDuration, "M"))

	//How many seconds
	secondsRegex := regexp.MustCompile(`[0-9]*S`)
	strSecondsDuration := secondsRegex.FindString(duration)
	durationSeconds, _ = strconv.Atoi(strings.TrimRight(strSecondsDuration, "S"))

	//Add time
	if duration[0:1] == "P" {
		totalSeconds = durationSeconds + (durationMinutes * 60) + (durationHours * 60 * 60) + (durationDays * 60 * 60 * 24)
		timeSeconds := time.Second * time.Duration(totalSeconds)
		returnDate = startTime.Add(timeSeconds)
	}

	//Subtract time
	if duration[0:2] == "-P" {
		totalSeconds = durationSeconds + (durationMinutes * 60) + (durationHours * 60 * 60) + (durationDays * 60 * 60 * 24)
		timeSeconds := time.Second * time.Duration(totalSeconds)
		returnDate = startTime.Add(-timeSeconds)
	}

	return returnDate, totalSeconds
}

// ConvFloatToStorage - takes given float64 value, returns a human readable storage capacity string
func ConvFloatToStorage(floatNum float64) (strReturn string) {
	if floatNum >= sizePB {
		strReturn = fmt.Sprintf("%.2fPB", floatNum/sizePB)
	} else if floatNum >= sizeTB {
		strReturn = fmt.Sprintf("%.2fTB", floatNum/sizeTB)
	} else if floatNum >= sizeGB {
		strReturn = fmt.Sprintf("%.2fGB", floatNum/sizeGB)
	} else if floatNum >= sizeMB {
		strReturn = fmt.Sprintf("%.2fMB", floatNum/sizeMB)
	} else if floatNum >= sizeKB {
		strReturn = fmt.Sprintf("%.2fKB", floatNum/sizeKB)
	} else {
		strReturn = fmt.Sprintf("%vB", int(floatNum))
	}
	return
}

// Logger -- function to create or append to a log file
func Logger(t int, s string, outputToCLI bool, fileName string) {
	cwd, _ := os.Getwd()
	logPath := cwd + "/log"
	logFileName := logPath + "/" + fileName

	//-- If Folder Does Not Exist then create it
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err := os.Mkdir(logPath, 0777)
		if err != nil {
			color.Red("Error Creating Log Folder %q: %s \r", logPath, err)
			os.Exit(101)
		}
	}

	//-- Open Log File
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	// don't forget to close it
	if err != nil {
		//We didnt manage to open the log file so exit the function
		color.Red("Error Opening Log File %q: %s \r", logFileName, err)
		return
	}
	defer f.Close()
	if err != nil {
		color.Red("Error Creating Log File %q: %s \n", logFileName, err)
		os.Exit(100)
	}
	// assign it to the standard logger
	log.SetOutput(f)
	var errorLogPrefix string
	//-- Create Log Entry
	switch t {
	case 1:
		errorLogPrefix = "[DEBUG] "
		if outputToCLI {
			color.Set(color.FgGreen)
			defer color.Unset()
		}
	case 2:
		errorLogPrefix = "[MESSAGE] "
		if outputToCLI {
			color.Set(color.FgGreen)
			defer color.Unset()
		}
	case 3:
		if outputToCLI {
			color.Set(color.FgGreen)
			defer color.Unset()
		}
	case 4:
		errorLogPrefix = "[ERROR] "
		if outputToCLI {
			color.Set(color.FgRed)
			defer color.Unset()
		}
	case 5:
		errorLogPrefix = "[WARNING] "
		if outputToCLI {
			color.Set(color.FgYellow)
			defer color.Unset()
		}
	case 6:
		if outputToCLI {
			color.Set(color.FgYellow)
			defer color.Unset()
		}
	case 7:
		if outputToCLI {
			color.Set(color.FgBlue)
			defer color.Unset()
		}
	}
	if outputToCLI {
		fmt.Printf("%v \n", errorLogPrefix+s)
	}
	log.Println(errorLogPrefix + s)
}

// ConfirmResponse - prompts user, expects a fuzzy yes (or a provided string) or no response, does not continue until this is given
func ConfirmResponse(confirm string) bool {
	var cmdResponse string
	_, errResponse := fmt.Scanln(&cmdResponse)
	if errResponse != nil {
		log.Fatal(errResponse)
	}
	if confirm != "" {
		if cmdResponse == confirm {
			return true
		} else if cmdResponse == "n" || cmdResponse == "no" || cmdResponse == "N" || cmdResponse == "No" || cmdResponse == "NO" {
			return false
		} else {
			color.Red("Please enter " + confirm + " or no to continue:")
			return ConfirmResponse(confirm)
		}
	}

	if cmdResponse == "y" || cmdResponse == "yes" || cmdResponse == "Y" || cmdResponse == "Yes" || cmdResponse == "YES" {
		return true
	} else if cmdResponse == "n" || cmdResponse == "no" || cmdResponse == "N" || cmdResponse == "No" || cmdResponse == "NO" {
		return false
	} else {
		color.Red("Please enter yes or no to continue:")
		return ConfirmResponse(confirm)
	}
}

// Encrypt - takes mandatory input string and entropy, returns b64 encoded string
func Encrypt(inputString, entropy string) (string, error) {
	if inputString == "" {
		return "", errors.New("inputString input is mandatory")
	}
	if entropy == "" {
		return "", errors.New("entropy input is mandatory")
	}
	encrypted, err := encryptBytes([]byte(inputString), cryptProtectUIForbidden, entropy)
	b64Enc := base64.StdEncoding.EncodeToString(encrypted)
	return b64Enc, err
}

// Decrypt - takes mandatory b64 encoded input string and entropy, returns  encoded string
func Decrypt(inputString, entropy string) (string, error) {
	if inputString == "" {
		return "", errors.New("inputString input is mandatory")
	}
	if entropy == "" {
		return "", errors.New("entropy input is mandatory")
	}
	raw, err := base64.StdEncoding.DecodeString(inputString)
	if err != nil {
		return "", err
	}
	b, err := decryptBytes(raw, cryptProtectUIForbidden, entropy)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func encryptBytes(data []byte, cf cryptProtect, entropy string) ([]byte, error) {
	var (
		outblob dataBlob
		r       uintptr
		err     error
	)

	r, _, err = procEncryptData.Call(uintptr(unsafe.Pointer(newBlob(data))), 0, uintptr(unsafe.Pointer(newBlob([]byte(entropy)))), 0, 0, uintptr(cf), uintptr(unsafe.Pointer(&outblob)))

	if r == 0 {
		return nil, err
	}
	enc := outblob.toByteArray()
	return enc, outblob.free()
}

// DecryptBytes decrypts a byte array returning a byte array
func decryptBytes(data []byte, cf cryptProtect, entropy string) ([]byte, error) {
	var (
		outblob dataBlob
		r       uintptr
		err     error
	)

	r, _, err = procDecryptData.Call(uintptr(unsafe.Pointer(newBlob(data))), 0, uintptr(unsafe.Pointer(newBlob([]byte(entropy)))), 0, 0, uintptr(cf), uintptr(unsafe.Pointer(&outblob)))
	if r == 0 {
		return nil, err
	}

	dec := outblob.toByteArray()
	outblob.zeroMemory()
	return dec, outblob.free()
}

func (b *dataBlob) toByteArray() []byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
	return d
}

func newBlob(d []byte) *dataBlob {
	if len(d) == 0 {
		return &dataBlob{}
	}
	return &dataBlob{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

func (b *dataBlob) free() error {
	_, err := windows.LocalFree(windows.Handle(unsafe.Pointer(b.pbData)))
	if err != nil {
		return err
	}

	return nil
}

func (b *dataBlob) zeroMemory() {
	zeros := make([]byte, b.cbData)
	copy((*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:], zeros)
}

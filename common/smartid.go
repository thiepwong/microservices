package common

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// LuhnValidate function to check luhn string
func LuhnValidate(luhnString string) bool {
	checksumMod := calculateChecksum(luhnString, false) % 10

	return checksumMod == 0
}

func calculateChecksum(luhnString string, double bool) int {
	source := strings.Split(luhnString, "")
	checksum := 0

	for i := len(source) - 1; i > -1; i-- {
		t, _ := strconv.ParseInt(source[i], 10, 8)
		n := int(t)

		if double {
			n = n * 2
		}
		double = !double

		if n >= 10 {
			n = n - 9
		}

		checksum += n
	}

	return checksum
}

// GenerateSmartID function Generate an account id base in Luhn algorithm
func GenerateSmartID(systemCode int, nodeCode int, size int) (id uint64, err error) {
	if (systemCode > 9 || systemCode < 1) || (nodeCode < 0 || nodeCode > 255) {

		err = errors.New("System Code is invalid")
		return 0, err
	}

	var nodeCodeStr string
	if nodeCode >= 0 && nodeCode <= 9 {
		nodeCodeStr = "00" + strconv.Itoa(nodeCode)
	}
	if nodeCode >= 10 && nodeCode <= 99 {
		nodeCodeStr = "0" + strconv.Itoa(nodeCode)
	}

	if nodeCode > 99 {
		nodeCodeStr = strconv.Itoa(nodeCode)
	}
	randomCode := randomString(size - 6)
	accountSeed := strconv.Itoa(systemCode) + nodeCodeStr + strconv.Itoa(generateControlDigit(randomCode)) + randomCode
	accountSeed += strconv.Itoa(generateControlDigit(accountSeed))
	fmt.Println(accountSeed)
	accountID, err := strconv.ParseUint(accountSeed, 10, 64)
	if err != nil {
		return 0, err
	}
	return accountID, err
}

func randomString(size int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	source := make([]int, size)

	for i := 0; i < size; i++ {
		source[i] = rand.Intn(9)
	}

	return integersToString(source)
}

func generateControlDigit(luhnString string) int {
	controlDigit := calculateChecksum(luhnString, true) % 10

	if controlDigit != 0 {
		controlDigit = 10 - controlDigit
	}

	return controlDigit
}

func integersToString(integers []int) string {
	result := make([]string, len(integers))

	for i, number := range integers {
		result[i] = strconv.Itoa(number)
	}

	return strings.Join(result, "")
}

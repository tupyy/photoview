package scanner_utils

import (
	"log"
	"os"

	"github.com/photoview/photoview/api/repositories"
)

func FileExists(testPath string) bool {
	_, err := repositories.GetDataRepository().Stat(testPath)

	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		// unexpected error logging
		log.Printf("Error: checking for file existence (%s): %s", testPath, err)
		return false
	}
	return true
}

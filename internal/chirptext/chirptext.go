package chirptext

import "strings"

func replaceChirpInput(input string) string {
	// replace all instances of words kerfuffle, sharbert and fornax with  ****
	splitText := strings.Split(input, " ")
	finalText := ""

	for _, word := range splitText {
		//need to validate if the word is  kerfuffle, sharbert or fornax and replace it with ****
		lowerCaseWord := strings.ToLower(word)
		if lowerCaseWord == "kerfuffle" || lowerCaseWord == "sharbert" || lowerCaseWord == "fornax" {
			finalText += "**** "
		} else {
			finalText += word + " "
		}
	}
	return strings.Trim(finalText, " ")
}

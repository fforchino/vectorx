package intents

import (
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"math/rand"
	"time"
)

/**********************************************************************************************************************/
/*                                            PILLS OF WISDOM                                                         */
/**********************************************************************************************************************/

func PillsOfWisdom_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"tell me something"}
	utterances[LOCALE_ITALIAN] = []string{"dimmi qualcosa"}
	utterances[LOCALE_SPANISH] = []string{"dime algo"}
	utterances[LOCALE_FRENCH] = []string{"dis moi quelque chose"}
	utterances[LOCALE_GERMAN] = []string{"erzähle mir etwas"}

	var intent = IntentDef{
		IntentName: "extended_intent_pills_of_wisdom",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    pillsOfWisdom,
	}
	*intentList = append(*intentList, intent)

	return nil
}

func pillsOfWisdom(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_GREETING_HELLO
	sdk_wrapper.SayText(getRandomSentence())
	return returnIntent
}

func getRandomSentence() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	sentences := [5][10]string{
		// English
		[10]string{
			"Fortune favors the bold",
			"I think, therefore I am",
			"Time is money",
			"I came, I saw, I conquered",
			"When life gives you lemons, make lemonade",
			"Practice makes perfect",
			"Knowledge is power",
			"Have no fear of perfection, you'll never reach it",
			"No pain no gain",
			"That which does not kill us makes us stronger",
		},
		// Italian
		[10]string{
			"La fortuna aiuta gli audaci",
			"Penso quindi sono",
			"Il tempo è denaro",
			"Sono venuto, ho visto, ho conquistato",
			"Quando la vita ti dà i limoni, fai limonata",
			"La pratica rende perfetti",
			"Sapere è potere",
			"Non temere la perfezione, perché non la raggiungerai mai",
			"Senza dolore non c'è vincitore",
			"Ciò che non ci uccide ci rende più forti",
		},
		// Spanish
		[10]string{
			"La fortuna favorece a los atrevidos",
			"Pienso, luego existo",
			"El tiempo es dinero",
			"Vine, mire, conquiste",
			"Cuando la vida te da limones, haz limonada",
			"La práctica hace la perfección",
			"El conocimiento es poder",
			"No temas a la perfección, tanon nunca la alcanzarás",
			"Sin dolor no hay ganancia",
			"Lo que no nos mata nos hace más fuertes",
		},
		// French
		[10]string{
			"La fortune sourit aux audacieux",
			"Je pense donc je suis",
			"Le temps, c'est de l'argent",
			"Je suis venu, j'ai vu, j'ai vaincu",
			"Lorsque la vie vous donne des citrons, faites de la limonade",
			"C'est en forgeant qu'on devient forgeron",
			"Savoir est le pouvoir",
			"Ne crains pas la perfection, vous ne l'atteindrez jamais autant",
			"On a rien sans rien",
			"Ce qui ne nous tue pas nous rend plus fort",
		},
		// German
		[10]string{
			"Dem Mutigen gehört die Welt",
			"Ich denke, also bin ich",
			"Zeit ist Geld",
			"Ich kam, ich sah, ich eroberte",
			"Wenn das Leben dir Zitronen gibt, mach Limonade daraus",
			"Übung macht den Meister",
			"Wissen ist Macht",
			"Fürchte keine Perfektion, du wirst es nie so sehr erreichen",
			"Kein Schmerz kein Gewinn",
			"Was uns nicht tötet, macht uns stärker",
		},
	}

	num := 0
	currentLanguage := sdk_wrapper.GetLanguage()
	if currentLanguage == sdk_wrapper.LANGUAGE_ITALIAN {
		num = 1
	} else if currentLanguage == sdk_wrapper.LANGUAGE_SPANISH {
		num = 2
	} else if currentLanguage == sdk_wrapper.LANGUAGE_FRENCH {
		num = 3
	} else if currentLanguage == sdk_wrapper.LANGUAGE_GERMAN {
		num = 4
	}

	phrase := sentences[num][r1.Intn(10)]
	return phrase
}

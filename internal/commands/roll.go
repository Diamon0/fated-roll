package commands

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"slices"
	"strconv"

	"github.com/expr-lang/expr"
)

type rollParams struct {
	Dice               int
	Faces              int
	KeepHighest        int
	KeepLowest         int
	ModifierExpression string
}

// oooo scary (no, seriously, scary, so many capture groups).
//
// Indexes:
//
// 0 is the number of dice.
//
// 1 is the number of faces of the dice.
//
// 2 is the keep-highest number.
//
// 3 is the keep-lowest number.
//
// 4 is the expression.
//
// If you use FindStringSubmatch, add 1 to all of them.
//
// If I stick to regex, maybe make a custom type?
var regex *regexp.Regexp

func init() {
	// If you change this, remember to also change the description of the variable
	regex = regexp.MustCompile(
		`^(\d+)` +
			`(?:d(\d+))` +
			`?(?:kh(\d+))` +
			`?(?:kl(\d+))` +
			`?(.*)$`,
	)
}

// For now, I am just gonna slap this together with flex tape.
// Yes, it does several passes through the string, I know.
//
// TODO: For when we have telemetry, log the errors and stuff, rn we only use the top-level error we create.
func RollParse(str string) (rollParams, error) {
	params := rollParams{}

	// Quick and dirty
	matches := regex.FindStringSubmatch(str)

	if matches[1] == "" {
		params.Dice = 1
	} else {
		num, err := strconv.Atoi(matches[1])
		if err != nil {
			return params, fmt.Errorf("Failed to parse number of dice")
		}
		if num < 1 {
			return params, fmt.Errorf("Number of dice must be higher than 0")
		}
		params.Dice = num
	}

	if matches[2] == "" {
		params.Faces = 20
	} else {
		num, err := strconv.Atoi(matches[2])
		if err != nil {
			return params, fmt.Errorf("Failed to parse number of faces")
		}
		if num <= 1 {
			return params, fmt.Errorf("Number of faces must be higher than 1 (Just use a calculator at this point lol)")
		}
		params.Faces = num
	}

	if matches[3] != "" {
		num, err := strconv.Atoi(matches[3])
		if err != nil {
			return params, fmt.Errorf("Failed to parse keep highest number")
		}
		if num < 1 {
			return params, fmt.Errorf("Amount of highest-rolling dice to keep must be more than 0")
		}
		params.KeepHighest = num
	}

	if matches[4] != "" {
		num, err := strconv.Atoi(matches[4])
		if err != nil {
			return params, fmt.Errorf("Failed to parse keep lowest number")
		}
		if num < 1 {
			return params, fmt.Errorf("Amount of lowest-rolling dice to keep must be more than 0")
		}
		params.KeepLowest = num
	}

	if matches[5] != "" {
		params.ModifierExpression = matches[5]
	}

	return params, nil
}

func Roll(args []string) string {
	params, err := RollParse(args[0])
	if err != nil {
		return "Sorry, we got the following error: " + err.Error()
	}

	rolls := []int64{}
	rollsText := ""

	for dice := range params.Dice {
		roll, err := rand.Int(rand.Reader, big.NewInt(int64(params.Faces)+2))
		if err != nil {
			return "Sorry, our roll engine seems to be experiencing issues..."
		}
		roll.Add(roll, big.NewInt(1))

		rolls = append(rolls, roll.Int64())
		if dice == 0 {
			rollsText = rollsText + roll.String()
		} else {
			rollsText = rollsText + ", " + roll.String()
		}
	}

	slices.Sort(rolls)

	if params.KeepHighest != 0 && params.KeepLowest != 0 {
		if len(rolls) >= (params.KeepHighest + params.KeepLowest) {
			// WARN: Yes, we are dropping the unused rolls into the ether
			rolls = append(rolls[0:params.KeepLowest], rolls[len(rolls)-params.KeepHighest:]...)
		}
	} else if params.KeepHighest != 0 {
		rolls = rolls[len(rolls)-params.KeepHighest:]
	} else if params.KeepLowest != 0 {
		rolls = rolls[0:params.KeepLowest]
	}

	var sum int64 = 0
	for num := range rolls {
		sum += rolls[num]
	}

	sumString := strconv.FormatInt(sum, 10)

	message := ""

	if len(args) > 1 {
		message += args[1]
	} else {
		message += "**Result:** "
	}

	message += strconv.Itoa(params.Dice) + "d" + strconv.Itoa(params.Faces) + " "

	message += fmt.Sprintf("(%v) ", rollsText)

	if params.ModifierExpression != "" {
		message += params.ModifierExpression + "\n"
		evaluator, err := expr.Compile(sumString + "" + params.ModifierExpression)
		if err != nil {
			return "Sorry, our roll engine seems to be experiencing issues..."
		}

		total, err := expr.Run(evaluator, nil)
		if err != nil {
			return "Sorry, our roll engine seems to be experiencing issues..."
		}

		message += "**Total:** " + fmt.Sprint(total)

	} else {
		message += "\n"
		message += "**Total:**" + sumString
	}

	return message
}

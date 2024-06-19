package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"time"
)

type card struct {
	cardIndex int
	suitValue int
}

func (c *card) suit() string {
	switch c.suitValue {
	case 0:
		return "hearts"
	case 1:
		return "clubs"
	case 2:
		return "diamonds"
	case 3:
		return "spades"
	default:
		panic("Random number out of range")
	}
}

func (c *card) name() string {
	switch c.cardIndex {
	case 1:
		return "ace"
	case 11:
		return "jack"
	case 12:
		return "queen"
	case 13:
		return "king"
	default:
		return strconv.Itoa(c.cardIndex)
	}
}

func (c *card) value() int {
	if c.cardIndex > 10 {
		return 10
	} else {
		return c.cardIndex
	}
}

func contains(s []card, e card) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func drawCard(drawnCards *[]card, player string, visible bool) card {

	newCard := card{rand.IntN(13) + 1, rand.IntN(4)}
	time.Sleep(500 * time.Millisecond)

	for {
		newCard := card{rand.IntN(13) + 1, rand.IntN(4)}
		if !contains(*drawnCards, newCard) {
			*drawnCards = append(*drawnCards, newCard)
			break
		}
	}
	if player == "dealer" {
		if visible {
			fmt.Println("The dealer has drawn a", newCard.name(), "of", newCard.suit())
		} else {
			fmt.Println("The dealer has drawn a card")
		}
	} else {
		fmt.Println("You have drawn a", newCard.name(), "of", newCard.suit())

	}

	return newCard
}

func playDealersHand(dealerCards *[]card, drawnCards *[]card) int {

	stick := false
	cards := *dealerCards
	fmt.Println("Dealer revels their second card as", cards[1].name(), "of", cards[1].suit())

	dealerTotal := calculateTotal(*dealerCards)
	for !stick {
		if dealerTotal > 17 {
			stick = true
		} else {
			newCard := drawCard(drawnCards, "dealer", true)
			dealerTotal += newCard.value()
		}
	}

	return dealerTotal

}

func determineWinner(dealerTotal int, playerTotal int) bool {
	return playerTotal <= 21 && playerTotal > dealerTotal || dealerTotal > 21
}

func calculateTotal(cards []card) int {
	total := 0
	for _, card := range cards {
		total += card.value()
	}
	return total
}

func getUserBet(remainingChips int) int {
	validBet := false
	bet := 0
	for !validBet {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("How much would you like to bet? Remaining chips: %d\n", remainingChips)
		scanner.Scan()
		testBet, err := strconv.Atoi(scanner.Text())
		if err != nil || testBet > remainingChips {
			fmt.Printf("Invalid input! Please input a whole number below %d\n", remainingChips)
		} else {
			validBet = true
			bet = testBet
		}

	}
	return bet
}

func main() {

	var drawnCards []card
	var playerCards []card
	var dealerCards []card
	playing := true

	chipsTotal := 1000

	for playing {
		playerTotal := 0
		playerCards = nil
		dealerCards = nil

		bet := getUserBet(chipsTotal)

		playerCards = append(playerCards, drawCard(&drawnCards, "player", true))
		dealerCards = append(dealerCards, drawCard(&drawnCards, "dealer", true))
		playerCards = append(playerCards, drawCard(&drawnCards, "player", true))
		dealerCards = append(dealerCards, drawCard(&drawnCards, "dealer", false))

		playerTotal = calculateTotal(playerCards)

		for playerTotal < 22 {
			fmt.Printf("Current total: %d \n", playerTotal)

			options := "hit, stand"
			scanner := bufio.NewScanner(os.Stdin)
			if playerCards[0].cardIndex == playerCards[1].cardIndex {
				options += ", split"
			}
			fmt.Printf("Would you like to %s? ", options)
			scanner.Scan()
			text := scanner.Text()
			if text == "hit" {
				newCard := drawCard(&drawnCards, "player", true)
				playerTotal += newCard.value()
			} else if text == "stand" {
				fmt.Printf("You stand on %d\n", playerTotal)
				time.Sleep(500 * time.Millisecond)
				dealerTotal := playDealersHand(&dealerCards, &drawnCards)
				time.Sleep(1 * time.Second)

				if determineWinner(dealerTotal, playerTotal) {
					fmt.Printf("The player has won with %d, over dealers %d\n", playerTotal, dealerTotal)
					chipsTotal += bet
				} else {
					fmt.Printf("The dealer has won with %d, over players %d\n", dealerTotal, playerTotal)
					chipsTotal -= bet
				}
				break
			}
		}
		if playerTotal > 21 {
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("You went bust on %d. Better luck next time\n", playerTotal)
			chipsTotal -= bet
		}
		fmt.Printf("You now have %d chips left\n", chipsTotal)
		time.Sleep(1000 * time.Millisecond)

		if chipsTotal > 0 {
			fmt.Printf("Would you like to play again?\n")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			ans := scanner.Text()
			if ans != "yes" {
				playing = false
			}

		} else {
			fmt.Println("Sorry, you've run out of money :(")
			os.Exit(0)
		}

	}
}

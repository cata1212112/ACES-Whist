package ACES

type Round struct {
	sumOfBids     int
	trump         Card
	first         Card
	winningCard   Card
	winningPlayer *Player
}

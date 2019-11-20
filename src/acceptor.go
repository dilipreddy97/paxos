package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Acceptor struct {
	pid              string
	leaderFacingPort string
	currentBallot    int
	accepted         map[int]string
}

//Acceptor passive responds to incoming messages from scouts/commanders
//On p1a, acceptor adopts received ballot if it exceeds current ballot; responds back with current ballot number
//On p2a, acceptor accepts b,s,p if b >= current ballot; replies with current ballot
func (self *Acceptor) Run() {
	fmt.Println("LEADER FACING PORT: " + self.leaderFacingPort)
	lLeader, _ := net.Listen(CONNECT_TYPE, CONNECT_HOST+":"+self.leaderFacingPort)

	defer lLeader.Close()

	for {
		connLeader, err := lLeader.Accept()
		if err != nil {
			fmt.Println("error from acceptor: ")
			fmt.Println(err)
			continue
		}
		reader := bufio.NewReader(connLeader)
		message, _ := reader.ReadString('\n')
	
		message = strings.TrimSuffix(message, "\n")
		messageSlice := strings.Split(message, ",")
		keyWord := messageSlice[0]

		retMessage := ""
		switch keyWord {
		case "p1a":

			receivedBallot := messageSlice[2] // b
			receivedBallotInt, _ := strconv.Atoi(receivedBallot)
			if receivedBallotInt > self.currentBallot {
				fmt.Println("ACCEPTED NEW BALLOT: " + receivedBallot)
				self.currentBallot = receivedBallotInt
			}
			retMessage += "p1b," + self.pid + "," + strconv.Itoa(self.currentBallot)
			acceptedStr := ""
			for _, accepted := range self.accepted {
				acceptedStr += "," + accepted
			}
			retMessage += acceptedStr
			connLeader.Write([]byte(retMessage + "\n"))
			if crashStage == "p1b" {
				os.Exit(1)
			}

		case "p2a":
		
			pval := messageSlice[2]
			fmt.Println(pval)
			pvalSlice := strings.Split(pval, " ")
			receivedBallotInt, _ := strconv.Atoi(pvalSlice[0])
			receivedSlotInt, _ := strconv.Atoi(pvalSlice[1])
			if receivedBallotInt >= self.currentBallot {
				self.currentBallot = receivedBallotInt
				self.accepted[receivedSlotInt] = pval
				//self.accepted = append(self.accepted, pval)
			}
			fmt.Println(self.pid + " ACCEPTOR ACCEPTED")
			fmt.Println(self.accepted)
			retMessage += "p2b," + self.pid + "," + strconv.Itoa(self.currentBallot) + "\n"
			fmt.Println(retMessage)
			connLeader.Write([]byte(retMessage))

			if crashStage == "p2b" {
				os.Exit(1)
			}

		default:
			retMessage += "Invalid keyword, must be p1a or p2a or ping"
			connLeader.Write([]byte(retMessage + "\n"))

		}
		connLeader.Close()
	}

}


package main

import (
	"bufio"
	"net"
	"strings"
	"strconv"
)

type Acceptor struct {
	pid              string
	leaderFacingPort string
	currentBallot	 int
	accepted         []string
}


func (self *Acceptor) Run() {
    lLeader, error := net.Listen(CONNECT_TYPE, CONNECT_HOST+":"+self.leaderFacingPort)

    defer lLeader.Close()
	msg := ""
	connLeader, error := lLeader.Accept()
	reader := bufio.NewReader(connLeader)


	for {

		message, _ := reader.ReadString('\n')

		message = strings.TrimSuffix(message, "\n")
		messageSlice := strings.Split(message, ",")
		keyWord := messageSlice[0]

		retMessage := ""
		switch keyWord {
			case "p1a":
				leaderId := messageSlice[1]  // lambda
				receivedBallot:= messageSlice[2]  // b
				receivedBallotInt, _ := strconv.Atoi(receivedBallot)
				if receivedBallotInt > self.currentBallot {
					self.currentBallot = receivedBallotInt
				}
				retMessage += "p1b," + self.pid + "," + strconv.Itoa(self.currentBallot)  
				acceptedStr := ""
				for _, accepted := range self.accepted {
					acceptedStr += "," + accepted 
				}  
				retMessage += acceptedStr
				connLeader.Write([]byte(retMessage))
			
				
			case "p2a":
				leaderId := messageSlice[1]  // lambda
				pval := messageSlice[2]
				pvalSlice := strings.Split(pval, " ")
				receivedBallotInt, _ := strconv.Atoi(pvalSlice[0])
				if receivedBallotInt >= self.currentBallot {
					self.currentBallot = receivedBallotInt
					self.accepted = append(self.accepted, pval)
				}
				retMessage += "p2b," + self.pid + "," + strconv.Itoa(self.currentBallot)
				connLeader.Write([]byte(retMessage))
			default:	
				retMessage += "Invalid keyword, must be p1a or p2a"
				connLeader.Write([]byte(retMessage))
			
		}


	}


}

//Heartbeat should hhappen here


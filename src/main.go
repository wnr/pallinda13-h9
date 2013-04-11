/**
 * Created with IntelliJ IDEA.
 * User: lucas
 * Date: 2013-04-07
 * Time: 9:38 PM
 */
package main

import (
	"matching"
	"server"
	"client"
	"julia"
)

func main() {
	go matching.Main()
	go julia.Main()
	go server.Main()
	go client.Main()

	select {}
}

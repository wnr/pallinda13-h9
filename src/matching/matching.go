// http://www.nada.kth.se/~snilsson/concurrency/

//+Vad händer om man tar bort go-kommandot från Seek-anropet i main-funktionen?
//-Meddelanden skulle skickas på följande sätt: Anna -> Bob. Cody -> Dave. Eva ->
//+Vad händer om man byter deklarationen wg := new(sync.WaitGroup) mot var wg sync.WaitGroup och parametern wg *sync.WaitGroup mot wg sync.WaitGroup?
//-det blir deadlock. Eftersom wg skickas som värde (istället för pekare) till Seek så kommer wg.Done() endast att ändra wg lokalt i varje Seek.
//-det innebär att Main's wg aldrig kommer ändras, och wg.Wait ser till att Main-tråden ständigt väntar.
//+Vad händer om man tar bort bufferten på kanalen match?
//-eftersom det är ett udda antal personer som skickar, kommer ett meddelande skickas som ej kommer att mottas. Det innebär
//-att en seek-rutin kommer att blocka tills någon tar emot meddelandet, och Main kommer ej förbi wg.Wait().
//+Vad händer om man tar bort default-fallet från case-satsen i main-funktionen?
//-ingenting eftersom det alltid kommer finnas ett meddelande att läsa, då antalet personer är udda. Om antalet personer
//-istället hade varit jämnt, hade det blivit en deadlock. default-fallet är till för att programmet ska avslutas
//-om det inte finns något meddelande över.

package matching

import (
	"fmt"
	"sync"
)

// This programs demonstrates how a channel can be used for sending and
// receiving by any number of goroutines. It also shows how  the select
// statement can be used to choose one out of several communications.
func Main() {
	people := []string{"Anna", "Bob", "Cody", "Dave", "Eva", "lucas"}
	match := make(chan string, 1) // Make room for one unmatched send.
	wg := new(sync.WaitGroup)
	wg.Add(len(people))
	for _, name := range people {
		go Seek(name, match, wg)
	}
	wg.Wait()
	select {
	case name := <-match:
		fmt.Printf("No one received %s's message.\n", name)
	default:
		// There was no pending send operation.
	}
}

// Seek either sends or receives, whichever possible, a name on the match
// channel and notifies the wait group when done.
func Seek(name string, match chan string, wg *sync.WaitGroup) {
	select {
	case peer := <-match:
		fmt.Printf("%s sent a message to %s.\n", peer, name)
	case match <- name:
		// Wait for someone to receive my message.
	}
	wg.Done()
}

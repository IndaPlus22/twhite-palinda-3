// http://www.nada.kth.se/~snilsson/concurrency/

/* 
1. What happens if you remove the go-command from the Seek call in
the main function?

No Goroutines are created and the program will not work.


2. What happens if you switch the declaration wg := new(sync.WaitGroup)
to var wg sync.WaitGroup and the parameter wg *sync.WaitGroup to 
wg sync.WaitGroup?

Changing the parameter to wg sync.WaitGroup will not work because the
WaitGroup needs to be passed by reference. Thus, the Seek function will not be able
to update the WaitGroup.

3. What happens if you remove the buffer on the channel match?

The program will crash at the main routine. Since 
the channel hasn't got a buffer the message it expects the
message to be read instantaneously, but it's not. This will cause
a deadlock.


4. What happens if you remove the default-case from the case-statement 
in the main function?

The program should run without any problems. The default
case is only executed if there is no pending send operation. Since
there is always a pending send operation the default case will never
be executed.

*/
package main

import (
	"fmt"
	"sync"
)

// This programs demonstrates how a channel can be used for sending and
// receiving by any number of goroutines. It also shows how the select
// statement can be used to choose one out of several communications.
func main() {
	people := []string{"Anna", "Bob", "Cody", "Dave", "Eva"}
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

# Thread Replication Proposal
A proposal attempting to solve the problem of being able to replicate a single thread across multiple boards.
 
### Changes to Thread
Here is the proposed new Thread structure:
```go
package main

type Thread struct {
        Name        string `json:"name"`
        Desc        string `json:"description"`
        MasterBoard string `json:"master_board"`
        Hash        string `json:"hash" enc:"-"`
}
```
**Notice the extra field; `MasterBoard`.** This will store the hex representation of the master-board's public key.

Note that "master-board" in this context represents the thread's master-board (owner and creator of the thread).

**Any changes made to a thread, needs to go through the thread's master-board's master-node** (via RPC).

### Thread Duplication

As the actual posts of the 
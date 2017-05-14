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

Note that "master-board", in this context, represents the thread's master-board (owner and creator of the thread).

**Any changes made to a thread, needs to go through the thread's master-board's master-node** (via RPC).

### Thread Replication

As the actual posts of the thread is "referenced" from another type (`ThreadPage`), when we duplicate a `Thread`, the associated `ThreadPage` will need to be synced.
 
Hence, for each board, the master-node of that board needs to check for changes of replicated threads.

Furthermore, when a client wants to make changes to a thread, it does it via RPC to the thread's master-node.
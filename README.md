# Word of wisdom.


Good day!

It was an interesting task.  I would like to give her more time.
The essence of the task is to make a tcp server with protection against DDoS attacks based on the pow algorithm.
I have not solved this problem before, so I decided to do a little research.
I found out that such solutions exist, for example, the Tor project has a proposal [torproject](https://github.com/torproject/torspec/blob/main/proposals/327-pow-over-intro.txt)
in which they suggest using EquiX to protect against DDoS attacks. [equix](https://github.com/tevador/equix)
Well, I thought it was not a bad option, but I did not find the implementation of this algorithm on golang, and decided to study it a little more. I found out that the algorithm is based on the problem of birthdays. I searched and found, as it seemed to me, an algorithm with the same advantages and disadvantages called Equihash and found several examples of its implementation on golang, from which I put together my own version.
Equihash is an asymmetric proof-of-work algorithm based on a computationally hard generalized birthday problem, which requires a lot of memory to generate a proof, but is instant to verify. Equihash is adapted as the PoW in [Zcash](https://z.cash/) a public implementation of the cryptocurrency protocol [Zerocash](http://zerocash-project.org/paper). It is possible to use Equihash in TLS as a [client puzzle](https://tools.ietf.org/html/draft-nygren-tls-client-puzzles-00).


Equihash has two parameters: **N** (width in bits) and **K** (length), which determine the complexity of the underlying problem and thus the memory and time complexity of the Equihash PoW. The underlying hash function is Blake2b, but any collision-resistant hash function would work too.

The time complexity is proportional to K2^{N/(K+1)}, and memory complexity to 2^{K+N/(K+1)}. The proof size is 2^{K}(1+N/(K+1))+192 bits. Verification requires 2^K hashes and XORs.

K and N are pre-installed on the client and on the server. This is a simplification adopted by me in order not to inflate the implementation code. For a real task, I would transmit these values over the network so that it would be possible to adjust the complexity.

The algorithm of action is simple, when connected, the client receives a challenge string, for which he must calculate the solution and send it back to the server. If the decision was correct, the server will issue a quote. Otherwise, the client will be disconnected.

### Docker

```bash
# Run client and server just for demo
$ docker-compose up

# Run server listening on 8080 port
$ docker-compose up -d server
```

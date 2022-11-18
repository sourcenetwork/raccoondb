
The racoondb abstraction ideally manages all the indexes.

Users could manually define / declare indexes and they would be responsible for using the correct index in a lookup?
maybe so

memdb uses a schema and i guess the internals decides which index to use?

overall though, the racoondb project isr eally interesting and its tapping into all the databases and distributed system stuff i've read about recently

one thing is that it would require concurrency control


Ok, so the biggest picture:

We have an abstraction of a KV store which allows for get, set, filter, iter, delete.

Internaly it would be great if this abstraction was able to choose which index to use

The design of hashicorps memdb is interesting

it works with structured data

the racoondb abstraction difference is that it would work with graphs
the minimal functionaly should at the very least contain a bfs / dfs algs?

nevertheless the point is that nodes will be stored somehow

the node storage should be indexable
the way in which data is stored in zanzibar is kinda wild
it stores the source and target node in the same entry, as opposed to the id
i guess that saves a lookup in the end?
ah, okay, so it's clever because the zanzibar graph uses virtual nodes,
no that is irrelevant, or should be.
i mean, no it isn't
filtervering by a virtual node would return nothing 

so really, why is the zanzibar graph stored as is? seems like it just adds some overhead?
actually it saves from overheads.
it saves doing a lookup each time at the cost of more storage.
it's a clever scheme, and if the nodes are relatively small in size it is potentially good

so ok,

graph abstraction on top of db store
edges store node data directly - just a zanzibar optimization
digraph stores directed edges
graph stores undirected edges, which means a backward index

nodes should be indexable
storing nodes separately from edges would be simpler
storing complete edges allows for more sophisticated indexes, indexes that goes accross tuples

it's all a bunch of tradeoffs, decisions should be made according what would best fit our use case

golden rule of engineering:
trade off between abstraction / generalization and performance

an interesting projects, lots of possible paths to choose from.

we need the notion of transaction to turn this into a lib, probably
also, we would require concurrency control for each row - a la two pahse commit

okay, lots of thoughts on this front, nothing is well defined. i'll just work on the specific stuff for our use case and write whatever i have in mind here

also consider the hashing method for node, there are tradeoffs there

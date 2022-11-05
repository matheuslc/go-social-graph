# Context
This project is a twitter-like app. So, we will have maaany relations between data, as we want to people talk to each other etc.

# Decision

The first option that came in my mind was PostgreSQL, which is a powerful database that I think is easy and cheap to run for so much time. When PostgeSQL start to get slow, you can also optimize your queries, indexes, check if some read-replicas meet your app access patterns and if so you can scale it vertically or even horizontally (harder).

But the thing is that is hard to create these queries that get data together, mainly when you will have too much relations. So, to make it easier to do the queries and create deeper connections, I've chosen a graph database. Using a Graph database also make way easier to do recommendations features, what I plan to do.

Also, we don't need consistency requiriments, and the schema can be pretty open.

# Consequences
I'll need to learn a new technologie, but that's the plan! 

# Status
Approved
# portal

Portal is supply benchmark status view and enqueue benchmark for the players.

## Features

* Sign in / sign out for the Team
* Assign webapp server configuration
* Enqueue benchmark
* Observe current benchmark status

#### Queue message format

Send to benchmark:

| Message                     |                    |
| ---                         | ---                |
| Data []byte                 | empty              |
| Attribute map[string]string | "team_id": TEAM_ID |

Receive from benchmark:

| Message                     |                                                         |
| ---                         | ---                                                     |
| Data []byte                 | directly json from the benchmark                        |
| Attribute map[string]string | "team_id": TEAM_ID,<br>"cretaed_at": CREATE_TIME        |

#### Queue flow

1. Send queue
2. Receive queue
3. Saves the receive queue message to Database
4. Show result message in portal client via the Database

#### Memo

* `/queues` algorithm
  * call `/enqueue` and save `msg_id`
  * call `{pubsub}/stat/{subscription}` for the msg_id set
    * need implements msg_id set
  * highlighting my team_id that match the msg_id set and the queues table

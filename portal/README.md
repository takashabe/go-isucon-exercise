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
| Attribute map[string]string | "team_id": TEAM_ID,<br>"cretaed_at": UNIX_TIME_SECONDS} |

#### Queue flow

1. Send queue
2. Receive queue
3. Saves the receive queue message to Database
4. Show result message in portal client via the Database

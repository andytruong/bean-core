Session life cycle
====

## Login

```
            login                         Create session
anonymous ---------------------> access ---------------------> DB
                                   |
            session.{jwt, token}   |
     user <------------------------+
```

## Active session

```
            payload with JWT              Valid JWT
    user ---------------------> access ---------------------> endpoint
                                   |                            |
            expired JWT            |                            |
         <-------------------------+                            |
                                                                |
            refresh JWT                                         |
    user ---------------------> access                          |
                                   |                            |
            JWT                    |                            |
         <-------------------------+                            |
                                                                |
            payload with JWT              Valid JWT             |
    user ---------------------> access ---------------------> endpoint
                                                                |
            response                                            |
         <------------------------------------------------------+

```

## Logout

```
            logout                      DELETE session
    user ---------------------> access ---------------------> DB
                                   |
                                   |
anonymous <------------------------+
```

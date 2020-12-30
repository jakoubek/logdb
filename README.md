# logdb

**logdb** is a minimal HTTP webserver app in Go with a request logger structure. The server counts requests to different routes in an in-memory instance and periodically syncs the log structure to a JSON file.

Because logging requests in a webserver is write-heavy with only few reads using a conventional database seemed not suitable to me.

The request logging logic logs all requests to an in-memory struct that gets written to a JSON file every 5 minutes.

```json
{
  "server_started": "2020-12-30T10:10:28.2853956Z",
  "requests": 9,
  "lastrequest": "2020-12-30T10:11:13.6521293Z",
  "routes": [
    {
      "route": "/",
      "requests": 7,
      "lastrequest": "2020-12-30T10:11:13.6521293Z"
    },
    {
      "route": "/notfound",
      "requests": 2,
      "lastrequest": "2020-12-30T10:11:12.4754029Z"
    }
  ]
}
```

On server startup the programs reads the JSON file and populates the in-memory struct.


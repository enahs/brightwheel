# Brightwheel Takehome

## Requirements

- requires go 1.22+

## run

```shell
$ go run main.go
```

## API

Get Cumulative

---

returns cumulative count for device ID
returns `404 Not Found` if no data exists for `id` .

request

```
GET /v1/devices/{id}/cumulative

curl localhost:8888/v1/devices/xyz12/cumulative
```

response

```JSON
{"count":6}
```

Latest Timestamp

---

returns latest timestamp for device ID
returns `404 Not Found` if no data exists for `id` .

request

```
GET /v1/devices/{id}/latest

curl localhost:8888/v1/devices/xyz123/latest
```

response

```JSON
{"latest_timestamp":"2021-09-29T16:09:15+01:00"}
```

Store Readings

---

stores readings for device ID specified in body
requires device ID to be specified.
Ignores subsequent readings for the same time.

returns `400 Bad request` if no `id` is specified in the request body.
returns `409 Unprocessible entity` if datum exists for `timestamp`.

request

```
curl  -v -X POST -d '{"id":"xyz123", "readings":[{"timestamp": "2021-09-29T16:08:15+01:00","count": 2}]}'  localhost:8888/v1/devices
```

## successful response

201 Created

```json
{ "success": true }
```

## Discussion

#### data structure

We essentially use a hashmap.
We store a pointer to device data.
As we inter data, we check if an entry exists for the given timestamp.
If not, we inter it, incrementing the sum on the global data reference and the latest timestamp.

This gives us O(1) lookups for any given data point, the sum and the timestamp.

#### Potential improvements

- middleware
  - auth
    some sort of bearer auth scheme seems appropriate here.
  - logging
    log requests and errors
  - recover
    catch panics, print stack, return error.
- environment
  move port to env variables, anything else.
- versioning
  URLs are versioned so client devices can be supported by different versions of the API
- structure
  all one file since this is a simple enough project.
- models + controllers  
  break out separate models/controllers into their own folder
  fully typed output for controller input/output
  validation on input object

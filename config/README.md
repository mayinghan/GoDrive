# GoDrive
A netdisk developed by Go

## Dev Manual
-  [Error handling](#error-handler)
	+ [Internal server error](#internal-server-error)
	+ [Client error](#client-error)

## Error handling
Guide of handling internal server or client error.
### Internal server error
Internal server error happens when the user's request is legal and the error happens due to operations in the server (i.e. errors when query databases/redis, errors of some Go build-in functions).

To handle these internal server errors, a [panic](https://blog.golang.org/defer-panic-and-recover) of an error string should be generated. 
```go
func handler(c *gin.Context) {
	_, err := QueryDB("abc") // err should implement the error interface
	if err != nil {
		panic(err.Error())
    // or panic("short description of the error"), using your own description is not recommanded
	}
}
```
Make sure the ```err``` implemented the golang "error" interface which has an ```Error()``` function which returns the error description string.

The error then will be handled by the error handling middleware.
### Client error
Client error happens when user's request is illegal (i.e. no user authenticated, user requesting unavailable resources (a.k.a 404 error), user request's request method is not provided, user request body's fields is wrong). 

To handle client error, return a JSON object directly without throw a panic. The response should have a status code of 200 to make the front end dev easier. 
In the response body, set a ```code``` field to 1 and a ```msg``` field to a short description of the error.
```go
func handler(c *gin.Context) {
	input := c.Query("number") // assuming user can only request number > 10
	if input <= 10 {
		// invalid input
		c.JSON(200, gin.H{
			"code": 1,
			"msg": "input should be larger than 10",
		}
	}
}
```


# MD5 calculator, GOlang

## Task:
The task is  to write a web service in the Go language, which allows to calculate the MD5 hash from a file located on the Internet. Downloading and calculation should occur in the background. 

## API:


* POST request to /submit with url parameters. To this request, the service must create a task with an identifier, according to which the user can learn about the state of its execution. As a response to this request, the service issues a query ID to the user.
* GET request to /check with id parameter. To this request, the service should return the task status to the user according to the id specified by the user. States - "the problem does not exist", "task in work", "task is completed", "task failed". If the task is completed, then in addition it is necessary to specify the url of the document and its calculated MD5 hash in the response. The status of the response code must be consistent with the response itself (404 if the task does not exist and so on)

## Examples of usage:
The following commands should be written in another terminal when the server is on:
```sh
>>> curl -X POST -d "url=http://site.com/file.txt" http://localhost:8000/submit
	{"id":"0e4fac17-f367-4807-8c28-8a059a2f82ac"}
>>> curl -X GET http://localhost:8000/check?id=0e4fac17-f367-4807-8c28-8a059a2f82ac
	{"status":"running"}
>>> curl -X GET http://localhost:8000/check?id=0e4fac17-f367-4807-8c28-8a059a2f82ac
	{"md5":"f4afe93ad799484b1d512cc20e93efd1","status":"done","url":"http://site.com/file.txt"}
```
## Run instructions:
* To run the program you need the GOlang compiler.
* Clone this repository and go to its folder.
* Run `go run serv1.go`
* Open new terminal and POST or GET command with port number, given after server run (default - 8080).
* You can use one of folo

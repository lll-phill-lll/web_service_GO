# MD5 calculator, GOlang

## Task:
The task is  to write a web service in the Go language, which allows to calculate the MD5 hash from a file located on the Internet. Downloading and calculation should occur in the background. 

## API:


* POST request to /submit with url parameters. To this request, the service must create a task with an identifier, according to which the user can learn about the state of its execution. As a response to this request, the service issues a query ID to the user.
* GET request to /check with id parameter. To this request, the service should return the task status to the user according to the id specified by the user. States - "the problem does not exist", "task in work", "task is completed", "task failed". If the task is completed, then in addition it is necessary to specify the url of the document and its calculated MD5 hash in the response. The status of the response code must be consistent with the response itself (404 if the task does not exist and so on)

## Examples of usage:
The following commands should be written in another terminal when the server is on:
```sh
>>> curl -X POST -d "url=https://www.google.com/robots.txt" http://localhost:8080/submit
	{"id":"0e4fac17-f367-4807-8c28-8a059a2f82ac"}
>>> curl -X GET http://localhost:8080/check?id=0e4fac17-f367-4807-8c28-8a059a2f82ac
	{"status":"running"}
>>> curl -X GET http://localhost:8080/check?id=0e4fac17-f367-4807-8c28-8a059a2f82ac
	{"md5":"09b67eacc3e50b9e34dcffb3771ef11e","status":"done","url":"https://www.google.com/robots.txt"}
```
## Run instructions:
* To run the program you need the GOlang compiler.
* Clone this repository and go to its folder.
* Run `go run serv1.go`
* Open new terminal and use POST or GET command with port number, given after server run (default - 8080). After entering command POST, you get the key that you can enter with GET command, to see the result of POST function.
* You can use one of following links to check the correctness of service. All the files in given links are also avaliable in the folder "test_files", you can download them and check for example [here](http://onlinemd5.com/).

## Tests:
While testing change the IDs, they are unique.
* 

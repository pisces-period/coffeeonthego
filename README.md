# COFFEE-ON-THE-GO: GOLANG HTCPCP IMPLEMENTATION
------------------------------
This is a simple and straightforward adaptation of the [HTCPCP] protocol (RFC 2324) implementation using Golang, Docker containers and MVC design pattern. 

WEBDAV PROPFIND and XML structures are not used - instead, your coffee is encapsulated in a delicious JSON format.

Speaking of capsules, this app includes a MongoDB container to store the configuration of all your coffee capsules at no cost.

Coffee pots have been replaced by a coffee-machine.

### STARTING THE COFFEE-MACHINE

Thanks to the magic of Docker containers, deploying a coffee-machine is very simple.

A _Dockerfile_ is provided, along with a _docker-compose.yml_ file.

The _Dockerfile_ uses a multi-stage approach to building all the necessary golang binaries into a super lightweight docker image (cca 22 MB).

The _docker-compose.yml_ file declares 2 services, one for the coffee-on-the-go app and one for the mongo-db.

To deploy and start your very own coffee-machine, run the following command at the root folder of this repository:
```
docker-compose up -d
```
To verify the appropriate installation, run the following command:
```
docker container ls
```

You should see an output similar to the following:
```
CONTAINER ID        IMAGE                 COMMAND                  CREATED             STATUS              PORTS                  NAMES
f760adbaf0e5        coffeeonthego:v1.0        "./..."                  About an hour ago   Up About an hour    0.0.0.0:8080->80/tcp   coffee-on-the-go
6084b7c26084        mongo:3.4.23-xenial   "docker-entrypoint.s…"   About an hour ago   Up About an hour    27017/tcp              mongo-on-the-go
```
As you can see, the Mongo-DB container is named _mongo-on-the-go_ and runs on port 27017. The coffee-machine container is named _coffee-on-the-go_ and is mapped to the host's port 8080.

If you wish to customize configurations such as the container names and assigned ports, you can change these settings in the _docker-compose.yml_ file.

### BREWING COFFEE

As per the HTCPCP protocol implementation, both BREW and POST reques methods are supported. For the purpose of this documentation, they are to be taken as synonyms.

Coffee capsules undergo a slightly different process of preparation, but to maintain compatibility with the HTCPCP protocol we are keeping the BREW/POST nomenclature.

The coffee-machine has 2 possible states when it comes to processing capsules:

* start - add the capsule to the coffee-machine
* stop - stops the coffee-machine

Web forms are too formal - brewing coffee should be as easy as C U R L.

The BREW request body must contain a _"coffee-message"_ field with either _"start"_ or _"stop"_ value, which determines what action the coffee-machine should take. 

You can send as many BREW requests as you want (as long as you don't exceed your daily caffein limit).

Below you can find an example of a valid BREW request (the _"flavor"_ field is optional and defaults to _"traditional"_, if none is provided):

```curl
curl -X BREW http://localhost:8080/coffee/brew \
-H 'Content-Type: application/json' \
-d '{"flavor":"vanilla-sky", "coffee-message":"start"}'
```

An associated ID is generated for you whenever you send a request to start brewing coffee.

When you feel like your coffee is ready, to stop the coffee-machine, you must first provide the ID of the coffee capsule, along with a _"stop"_ coffee-message.

Failing to do so yields no result, and can potentially damage the coffee-machine and cause leakings.

Below you can find an example of a valid BREW request to stop the coffee-machine (both the _"id"_ and _"coffee-message"_ fields are required):

```curl
curl -X BREW http://localhost:8080/coffee/brew \
-H 'Content-Type: application/json' \
-d '{"id": "${bson_id}", "coffee-message":"stop"}'
```

### GETTING INFORMATION ABOUT COFFEE

As previously mentioned, all the relevant coffee meta-data is stored in a Mongo-DB database.

The db name is _coffeeshop_ and the collection name is _coffee_.

Coffee capsule data is stored in the following way:

| field              | description                  | optional  |
|:------------------:|:----------------------------:|:---------:|
| ID                 | capsule id                   | no        |
| Flavor             | capsule flavor               | yes       |
| PreparationState   | preparation state            | yes       |
| CoffeeMessage      | sets the preparation state   | no        |

The _PreparationState_ field is controlled by the _CoffeeMessage_ field and any attempt by you to set the _PreparationState_ via request body is overriden by the _CoffeeMessage_.

To get meta-data for a particular coffee capsule, all you need to do is send a GET request, setting the _id_ as a parameter.

Below you can find an example of a valid GET request (needless to say, the _"id"_ parameter is required):

```
curl http://localhost:8080/coffee/get?id={$bson_id}
```
Conversely, if you do not provide any id, the coffee-machine returns a list of all coffee capsules which have been requested.

Below you can find an example of such transaction:
```
curl http://localhost:8080/coffee/get
```
For a better user experience, use your web browser instead of CURL. Fresh go templates are served with your coffee.

### TEAPOTS ###

As you know, teapots differ from coffee pots and coffee machines.

The _/coffee/_ route is the only safe way to brew and get coffee.

Any attempt to use a _/teapot/_ route results in http status code 418.

Below you can find an example of an unadverted attempt to use teapots to brew coffee:
```
curl http://localhost:8080/teapot/
```
Again, for a better experience, use your web browser instead of CURL.

### INTERACTING WITH THE MONGO-DATABASE

If you want to check in real time how the data is populated into the database, you can connect to the database container through the following command:
```
docker exec -it mongo-on-the-go /bin/bash
```
Then, once inside the container, run the following command:
```
mongo -u ${user} -p --authenticationDatabase admin
```
Provide the password (as specified by the _docker-compose.yml_ file) and TA DA! You are connected to the database. 

Up next, let's do a simple query.
```mongo
use coffeeshop
db.coffee.find().pretty()
```
The command above should output all the coffee capsules recorded into the database and confirm that it is, indeed, working.

Read on for more insightful Mongo commands.

### FURTHER READING

Check out this article [HERE] for further information on how to package multiple libraries in Golang.

Check out [THIS] article for more interesting mongo commands and database queries.


[HTCPCP]: https://tools.ietf.org/html/rfc2324 
[HERE]: https://ieftimov.com/post/golang-package-multiple-binaries/
[THIS]: https://docs.mongodb.com/manual/reference/command/

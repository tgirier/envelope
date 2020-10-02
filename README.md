![envelope](https://i.imgur.com/ZzxTAKL.gif)
Envelope is a basic chat server using TCP websockets.

Envelope broadcasts messages to all connected users.

Built with :heart: in Paris

A huge thanks to @bitfield for his mentoring on this project.

# Installation

Envelope can be run in two different ways :

* Using the provided Dockerfile to build a **Docker image**
* Importing the package and use the **API** within another **Go Program**

## Docker

*Requires Docker to be installed on your machine.*

To build a Docker image and run it locally:

1. Build the image:
```bash
docker build https://github.com/tgirier/envelope.git
```
2. Run the container:
```bash
docker run -d -p 8080:8080 [YOUR-IMAGE-ID]
```

## Go API

To import the API within your Go program, simply add the following statement to your go package:
```Go
import "github.com/tgirier/chat"
```

# Usage

## Connect with nc

If to connect to your envelope server using nc, simply run:
```bash
nc <HOST> <PORT>
```
*Example - if your envelope server is running locally on port 8080:*
```bash
nc localhost 8080
```
If the connection is successful, the following welcome message will be displayed:
```bash
Welcome to envelope! Please enter your username:
```
Enter your username, press **Enter**. The notification will be boradcasted to all connected users:
```bash
<USERNAME> joined envelope
```
Then, simply enter your message and press **Enter** to send it to all connected users.

## Connect with the provided Go client

*Requires the envelope package to be imported within your Go package*

To create a new client connected to a given envelope server use ConnectClient. It takes the address of the target envelope server as a string:
```Go
client, err := ConnectClient(<SERVER_ADDRESS>)
```
*Example - For a local server running on 8080:*
```Go
client, err := ConnectClient("localhost:8080")
```

Next, the server will send the welcome message and ask for the username. To handle it, simply use the Read method on the newly created client:
```Go
client.Read()
```
Next, the server will expect a username. Send the username using the Send method:
```Go
client.Send("<USERNAME>"+"\n")
```
:warning: **All strings sent through the Send method must be ended by a newline character. Otherwise the server will not consider that the message is over.** :warning:

Then the client can send messages or read them using the provided methods.
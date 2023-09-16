# Ping-Pong Project

#### Setup

To install & setup this project, you need to have golang (version 1.21) installed & initiate the module through golang. The environment that this program runs on is through linux (specifically Ubuntu 20.04LTS). This is vital since the environment used to run this program is only compatable through linux (in terms of the install script). If you want to run this program through a different environment, you have to download golang by other means and initiate the module manually.

Dependencies:

- Go (Version `1.21`)
- Snap

Go Libraries:

- net
- yaml (`V3`)
- fmt
- os
- time
- math
- sort

Installation Steps:

1. Intall through the install script (`install.sh`).
2. Check to make sure that the golang version you have is `1.21` or higher.
3. Run `pingpong` and check out the features that are available.

Quick Fixes:

1. If you don't have a working module (meaning it wasn't initiated properly when you cloned the repo) try running `go get main`, it should install all of the dependencies associated with the `go.mod` file.
   - To grab associated dependencies: `go get main`
   - To list dependencies: `less go.mod` or `cat go.mod | less`

2. If you did not install golang version `1.21`, you need to goto the website and install it by whatever means is necessary from the source.
3. If the `install.sh` script isn't working because the permissions did not carry through, run `chmod +x install.sh`
   - Make install script executable: `chmod +x install.sh`
   - To run the install script: `./install.sh` or `sh install.sh`

4. If you do not have snap try installing golang through the `apt` install environment or building from the source.
5. If you cannot bind to a port (as pong role) or have the connection established (as a ping role) try running the program as `sudo` or as `root`.
   - Run program as sudo: `sudo go run main -r <ping or pong>`

#### Config File

The config file is a `config.yaml` file in the CWD (current working directory). The config file will have a structure formatted as such:

```
---
Roles:
  # Ping role message
  Ping: "ping_IBNHRCEG7UOQR9JAVN2JS0UQ1GM45AEKBCK3HN2PXMO7F9UHKKG26UBLRJVAGJQLCZQPZ7HDI6X74F74OK9489YW4NRSJ5VTHDF3S9ENT1305M3J3A78GXO8NSG"
  # Pong role message
  Pong: "pong_IMJ9S1KJD8FZYCKF8ILPXES1RTR685OTYOJYKK5EDMU1AO2XVX8Z2X8VKXKJ1RS2EC06N8A01IM5AKB5PB135Q81NXXVN6QNGLRTCX6I5P566ZS9I7BGB0L3BMG"
Host: "<Server-IP-Or-DNS-Name-Here>"
Ports:
  # Socket port
  Socket: "9988"
```

#### Pong Role (Server)

Pong role is just a server role that waits until a ping host pings the pong host and they continue to communicate until the run executions is done (10000 trials). This experiment is used to determine the network capabilities of the VM's in the cloud environment and determine their interconnected delays that may be necessary to understand when designing infrastructure in the cloud or building apps using the cloud.

```sh
# If you want to run the server through the go environment
sudo go run main -r pong
```

```sh
# If you want to run the server through a built environment
go build main
mv main pingpong
sudo ./pingpong -r pong
```

#### Ping Role (Client)

The ping role is a client role that sends a ping message to a pong host and has a timer for the round trip time (RTT) between the ping packet transmission and the pong packet transmission. This is vital for the use of understanding stats for a infrastructure that is going to operate with end users. When the program is run as a ping role, the results of the RTT data is stored in two files `rtt_socket_output_stats.txt` & `rtt_socket_output_file.txt`.

Before running the command, you need to edit the `config.yaml` file to properly associate the data to the pong server's configuration.

```sh
# If you want to run the client through the go environment
go run main -r ping
```

```sh
# If you want to run the client through a built environment
go build main
mv main pingpong
./pingpong -r ping
```
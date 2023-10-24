# KVDB Benchmark Project

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
- flag
- strconv
- strings

Installation Steps:

1. Intall through the install script (`install.sh`).
2. Check to make sure that the golang version you have is `1.21` or higher.
3. Run `kvdb` and check out the features that are available.

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

   - Run program as sudo: `sudo go run main -r <s or c>`
      - sudo go run main -r c -t ll
      - `sudo go run main -r s -id <int (1 - max number of nodes)>`

#### Config File

The config file is a `config.yaml` file in the CWD (current working directory). The config file will have a structure formatted as such:

```yaml
---
Hosts:
  - "host1-ip"
  - "host2-ip"
Ports:
  # KV Storage Service port
  Socket: "9988"
```

#### Server Role

The server role is essentially hosting the KVDB (Key-Value DB) service. It processes the commands sent from the client and executes them in an orderly fashion by storing the command in a DB log. After all clients send a command, the server executes the commands sequentially based on Client connection priority. 

```sh
# If you want to run the server through the go environment
sudo go run main -r s -id $VALUE
```

```sh
# If you want to run the server through a built environment
go build main
mv main kvdb
sudo ./kvdb -r s -id $VALUE
```

#### Client Role

The ping role is a client role that sends a ping message to a pong host and has a timer for the round trip time (RTT) between the ping packet transmission and the pong packet transmission. This is vital for the use of understanding stats for a infrastructure that is going to operate with end users. When the program is run as a ping role, the results of the RTT data is stored in two files `rtt_socket_output_stats.txt` & `rtt_socket_output_file.txt`.

Before running the command, you need to edit the `config.yaml` file to properly associate the data to the pong server's configuration.

```sh
# If you want to run the client through the go environment
go run main -r c -t ll
```

```sh
# If you want to run the server through a built environment
go build main
mv main kvdb
sudo ./kvdb -r c -t ll
```
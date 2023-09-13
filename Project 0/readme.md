# Ping-Pong Project

#### Setup

To install & setup this project, you need to have golang (version 1.21) installed & initiate the module through golang. The environment that this program runs on is through linux (specifically Ubuntu 20.04LTS). This is vital since the environment used to run this program is only compatable through linux (in terms of the install script). If you want to run this program through a different environment, you have to download golang by other means and initiate the module manually.

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

#### Pong Role (Server)

Pong role is just a server role that waits until a ping host pings the pong host and they continue to communicate until the run executions is done (10000 trials). This experiment is used to determine the network capabilities of the VM's in the cloud environment and determine their interconnected delays that may be necessary to understand when designing infrastructure in the cloud or building apps using the cloud.

```sh
# If you want to run the server through the go environemtn
go run main -r pong
```

```sh
# If you want to run the server through a built environment
go build main
mv main pingpong
./pingpong -r pong
```

#### Ping Role (Client)

When that happends you can read the output results from the `rtt_socket_output_stats.txt` & `rtt_socket_output_file.txt`

Pong role is just a server role that waits until a ping host pings the pong host and they continue to communicate until the run executions is done (10000 trials). This experiment is used to determine the network capabilities of the VM's in the cloud environment and determine their interconnected delays that may be necessary to understand when designing infrastructure in the cloud or building apps using the cloud.
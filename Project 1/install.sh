#!/bin/sh

# install golang
echo "Installing golang"
sudo snap install go --classic # running the install through snap
# sudo apt install golang-1.20

# # Check go version
clear
echo "Checking installed go version"
go version
echo "\n\n"

# Getting associated dependencies and modules for the program
go get main
go mod init main
go mod tidy

# Run the program
go run main -r s -id 1

# Let the user know they have to edit the config.yaml file
echo "\n\n\n\tDon't forget to edit the \"config.yaml\" file so that way you can associate the ping (client) to the pong (server) host."
echo "\tThe data associated from this experiment are listed as \"rtt_socket_output_stats.txt\" (stats) \n\talong the list of objects that were used to calculate the stats \"rtt_socket_output_file.txt\"."
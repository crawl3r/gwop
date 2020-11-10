# A script designed for setting up the environment to be able to use gwop.
# As gwop uses Go to build new implants, this script will check whether or not Go is installed and configured on the target box ($PATH specific checks)
# If not, it will download Go, install it and set up the ~/go environment and alter the users $PATH

# TODO: Check to see if Go is currently installed on the box
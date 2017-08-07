# spyder
API Testing and Request Framework


## API Testing and Requests

Spyder gives an easy interface to make and test API endpoints from within the terminal.
It uses simple JSON config files which are ment to be versioned. It also has built in
scripting support for dynamically modifying endpoint configuration and performing certain
types of tasks, such as appending authentication, persisting tokens, etc.


## Start a project

Spyder expects the configuration for a set of endpoints to live in its own folder structure.
This folder can then be versioned to allow for easy maintaince and team synchronization. 
The `init` command will generate the folder structure for you. Just run it in an empty folder.

```
spyder init
```

## Configuring an endpoint

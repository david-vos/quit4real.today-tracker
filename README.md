# Quit4Real.today
This project is the backend for https://Quit4Real.today
The goal is simple. Track a user to see *If* they play a game they said they would not.

# Suggestions
If you have any suggestions dont be scared to create an issue/PR. Ill be looking around at other systems to support in a bit.
Right now steam will be the main priority but I would like riot ( but there api keys suck ) and Blizzard. I have yet to look into Ubisoft and Xbox :eyes:

# Setup
```bash
cp config/.example.config.go config/config.go
go build -v
go run quit4real.today
```

# Docker compose
just do `docker compose up --build -d` 

# Desgin
![image](https://github.com/user-attachments/assets/306bab47-b1e1-4190-8939-6f9bc8f32799)
This is the idea behind it.

Endpoints can only talk to Services.
Services can talk to one another
Services can talk to an any Command or Query Handler
Commands and Queries cannot talk to each other. They do belong to eah other. Meaning that a FailRepository can only be called by an FailCommandHandler or FailQueryHandler
The Repository will persist it in SQLite or any other future projection


# VERY BETA
This project is currently stil a WIP.
Everything is currently subject to change. That includes the data design and reponse types
Right now I am bussy with the BE and will soon start working on the FE afther at least steam is added.

# Help me
If you want to help me just create a issue or PR. Dont be shy I dont bite

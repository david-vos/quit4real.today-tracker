# Quit4Real.today
This project is the backend for https://Quit4Real.today
The goal is simple. Track a user to see *If* they play a game they said they would not.

# Setup
```bash
cp config/.example.config.go config/config.go
go build
nohup go run project > output.log 2>&1 &
```
**Close the BE**
```bash
ps aux | grep 'go run'
```
```bash
kill <PID>
```

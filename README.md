# Project Management API

## Prerequisites
- Installed and **turned ON (daemon)** Docker


## Setup Instructions

```
git clone https://github.com/canyouhearthemusic/project-management.git

cd project-management
```
___
macOS/Linux:
```
cp .env.example .env
```
Windows:
```
copy .env.example .env
```

Then configure `.env` file by yourself.
___
And finally run the app.
```
make build && make up
```

## Endpoints (`/api/v1`)
https://project-management-82r5.onrender.com/swagger/index.html 


## Health status check
https://project-management-82r5.onrender.com/api/v1/heartbeat

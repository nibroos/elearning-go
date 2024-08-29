<h1 align="center">
    <br>
  Backend E-Learning
  <br>
</h1>

## 🚀 Quick Start
### Developement Environment
On `/` dir, Run `make copy-env`, Modify to suit your environment, focus on these key, you can leave others as it is. The key name is explanatory itself.
```bash
# MAIN APP PORT
GATEWAY_PORT=4050

# DATABASE
POSTGRES_PORT=5432
```

> Without docker, you need to install [air-verse](https://github.com/air-verse/air) to activate the hot reloading.

### 🐳 Docker :: Container Platform

[Docker](https://docs.docker.com/get-docker/) Install.

- On the root folder, Starts the containers in the background and leaves them running : `docker-compose -f docker/docker-compose-dev.yml up --build -d`
- Stops containers and removes containers, networks, volumes, and images : `docker-compose down`

## 🛎 Available Commands each Service

Change bash directory to each server.
> ${arg} means replace all of it match your args without space
- Create mirgration : `make migrate create name=${your_migration_name}`
- Run migration : `make migrate-up`
- Stepback migraiton: `make migrate-down`
- Generate proto file, leave the proto args blank if you want to generate all proto file: `make proto ${your-proto.proto}`. If its fail, run this command on specific service. for example, in /users-service/ run bash `export PATH="$PATH:$(go env GOPATH)/bin`

## 💎 The Package Features

<p>
  <img src="https://img.shields.io/badge/-Docker-2496ED?style=for-the-badge&logo=Docker&logoColor=fff" />&nbsp;&nbsp;
  <img src="https://img.shields.io/badge/-NGINX-269539?style=for-the-badge&logo=NGINX&logoColor=fff" />
  <img src="https://img.shields.io/badge/-Go-1185F4?style=for-the-badge&logo=Go&logoColor=fff" />
</p>
<p>
<img src="https://img.shields.io/badge/-PostgreSQL-336791?style=for-the-badge&logo=PostgreSQL&logoColor=fff" />&nbsp;&nbsp;
</p>

### 📗 API Document
All endpoints stored in  `-.json`

<h1 align="center">
    <br>
  Features
  <br>
</h1>

Feel free to ask if you have any questions or need more details!

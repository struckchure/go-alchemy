# Alchemy (Go)

**Alchemy** provides ready-to-use components for backend developers, empowering them to build applications quickly and efficiently.

---

## ✨ Philosophy

If frontend developers can do something as simple as:

```sh
npx shadcn add button
```

Then backend developers should be able to do:

```sh
alchemy add authentication
```

Alchemy bridges this gap, providing a modular and intuitive toolset for backend development.

---

## 🛠️ Installation

For Linux and MacOS

```sh
curl -fsSL https://raw.githubusercontent.com/struckchure/go-alchemy/main/scripts/install.sh | bash
```

For Windows

```sh
irm https://raw.githubusercontent.com/struckchure/go-alchemy/main/scripts/install.ps1 | iex
```

Or download binaries from [release page](https://github.com/struckchure/go-alchemy/releases)

---

## 🚀 Usage

### Initialize a New Alchemy Project

> You current directory must be a Golang project with `go.mod` for things to work properly.

Run the following command to initialize a new project:

```sh
alchemy init
```

You’ll be prompted to provide details interactively:

```plaintext
? Provide alchemy component root:  .
? Choose ORM:  Prisma
? Choose Database Provider:  PostgreSQL
? Provision Database with Docker Compose:  Yes
```

Alchemy will:

- Generate a Docker Compose file for the database.
- Configure your ORM (e.g., Prisma).
- Set up a new Prisma project.
- Update your Go dependencies.

**Example Output:**

```plaintext
✨ Alchemy config has been generated!
🛠️ Updating Go dependencies ...
🥂 You're all set!
```

To start the database service, run:

```sh
$ docker compose up -d
```

---

### Add Components to Your Project

Alchemy allows you to add modular backend components.

#### **Interactively Add a Module**

```sh
$ alchemy add
```

**Example Interactive Flow:**

```plaintext
? Component Category Authentication
? Select Components Login
Creating Authentication.Login component
  + prisma/schema.prisma
  + dao/user.go
  + services/authentication.go
+ Authentication.Login
```

#### **Add All Components from a Module**

```sh
$ alchemy add Authentication.All
```

This will add all components related to the `Authentication` module.

#### **Add a Specific Component**

```sh
$ alchemy add Authentication.Login
```

This will add only the `Login` component from the `Authentication` module.

> Component name is not case sensitive `Authentication.Login` is the same as `authentication.login`

---

## 📂 Example Project Structure

After running the commands, your project might look like this:

```
.
├── prisma
│   ├── schema.prisma
├── dao
│   ├── user.go
├── services
│   ├── authentication.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
```

---

## 🧑‍💻 Contributions

Contributions are welcome! Feel free to open issues or submit pull requests to enhance Alchemy.

---

## 📖 License

This project is licensed under the [MIT License](LICENSE).

---

## 🛟 Support

If you encounter any issues, please open an issue in the [GitHub repository](https://github.com/struckchure/go-alchemy).

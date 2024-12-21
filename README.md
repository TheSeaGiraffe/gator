## `gator`

This is an RSS feed aggregator created as part of the "Build a Blog Aggregator" course on
[Boot.dev.](https://www.boot.dev/) It's a CLI application that allows users to manage
their own personal collection of RSS feeds. Users are able to add any number of feeds,
store them in a local database, and browse through them at their leisure.

## Installation

`gator` requires Postgres and Go. Brief installation instructions for both will be
provided in the following subsections.

### Postgres

#### Docker

By far the easiest way to install Postgres is to simply pull the latest image from docker
hub:

```bash
docker pull postgres
```

Of course, this assumes that you already have Docker installed on your system.

#### Local installation

Most Linux distributions include Postgres in their repositories. Here we assume that
you're using Ubuntu where Postgres can be installed with the following command:

```bash
sudo apt install postgres
```

If you use another distribution, simply use the built-in package manager to install
Postgres.

If you're using a Mac, you can install Postgres with Homebrew:

```bash
brew install postgresql
```

After the installation is finished make sure to start the Postgres service:

```bash
brew services start postgresql
```

If you're using Windows, we recommend setting up a dev environment in WSL2 and installing
Postgres there.

### Go

As with Postres, Go is available in the repositories for most Linux distributions.
Assuming that you're using Ubuntu you can install Go with the following command:

```bash
sudo apt install go
```

For Mac, you can use Homebrew:

```bash
brew install go
```

For Windows, we again recommend setting up a development environment in WSL2 and
installing Go there.

### `gator`

Once you've installed all of the prerequisites you can install `gator` with `go install`:

```bash
go install github.com/TheSeaGiraffe/gator
```

`gator` should automatically be added to the system path.

## Usage

Before you can use `gator` you must first create a `.gatorconfig.json` file in your home
directory with the following contents:

```json
{
  "db_url": "db_connection_string"
}
```

Where `db_connection_string` is the connection string that will be used to connect to the
local Postgres DB. It should look something like this:

```
postgres://my_user:user_pass@db_hostname:db_port/db_name?sslmode=disable
```

where:

- `my_user`: The user that you registered when creating the database
- `user_pass`: The password for the user
- `db_hostname`: The name of the host that you're using for your DB. Assuming that you're
  using a local DB it should be `localhost`.
- `db_port`: The port used to access your DB. This will most likely be the default port of
  `5432`
- `db_name`: The name of the database that you'll be using for `gator`

You must then register a username. This can be done with the `gator register` command:

```bash
gator register john
```

Afterwards, you can begin adding feeds with the `addfeed` command:

```bash
gator addfeed "World News - The Guardian" "https://www.theguardian.com/world/rss"
```

Once you've started following a few feeds you can pull posts from the feed with the `agg`
command:

```bash
gator agg
```

Note that this command is meant to run continuously in the background. Additionally, it
takes an optional "duration" parameter of the form "`duration_amount duration_units`"
where `duration_amount` is an integer value and `duration_units` is a time unit such as
`h`, `m`, or `s`. If no duration is specified then `agg` defaults to 5 minutes. In order
to keep from DDoS-ing a feed it is recommended to use a duration that is not too short.
Once `agg` is running, users can then browse the saved posts with

```bash
gator browse
```

The `browse` has an optional "limit" parameter that specifies the maximum number of posts
to display. Once you are finished, you can stop the running `agg` process with `Ctrl+C`.

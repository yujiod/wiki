# A smart wiki software by Go language.

* [Go](http://golang.org) Programming Language
* [Revel](http://revel.github.io/) A high-productivity web framework for the Go language.
* [GORM](https://github.com/jinzhu/gorm/) The fantastic ORM library for Golang, aims to be developer friendly.

[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

## Features

* Show and edit pages by everybody
* Markdown editor
* Bracket link
    * ex) Make link of "Some Page": [[Some Page]]
* Revision and Diff

## Planned features

* Attachment file on local disk, Google Cloud Storae and Amazon S3.
* User authentication and permissions
* Detect editing page collision
* and more...

# Using on Docker

When you start container exposed port 9000, after open url http://localhost:9000/ from your browser.

## Default, SQLite3 in container.

    docker run -d -p 9000:9000 yujiod/wiki
    # same above
    docker run -d -p 9000:9000 -e DB_DRIVER=sqlite3 -e DB_SOURCE="./wiki.db" yujiod/wiki

## Using MySQL

    docker run -d -p 9000:9000 -e DB_DRIVER=mysql -e DB_SOURCE="dbuser:dbpass@tcp(hostname:3306)/dbname?charset=utf8" yujiod/wiki

The Data Source Name, see [Go-MySQL-Driver](https://github.com/go-sql-driver/mysql).

## Using PostgreSQL

    docker run -d -p 9000:9000 -e DB_DRIVER=postgres -e DB_SOURCE="host=hostname user=dbuser dbname=dbpass sslmode=disable" yujiod/wiki

The Connection String Parameters, see [pq](https://github.com/lib/pq).

# Using from source

You need [Go](http://golang.org), install it before.

    go get github.com/yujiod/wiki/app
    go get github.com/revel/cmd/revel
    revel run github.com/yujiod/wiki

## Build, install and run.

    # Installing to /usr/lcoal/wiki.
    revel build github.com/yujiod/wiki /usr/local/wiki
    /usr/local/wiki/run.sh

## Packaging for deploy

    # Creating wiki.tar.gz
    revel package github.com/yujiod/wiki

# License

Released under the [MIT License](http://www.opensource.org/licenses/MIT).

# A smart wiki software by Go language.

* [Go](http://golang.org) Programming Language
* [Revel](http://revel.github.io/) A high-productivity web framework for the Go language.
* [GORM](https://github.com/jinzhu/gorm/) The fantastic ORM library for Golang, aims to be developer friendly.

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

# Using on Dcoker

Default, sqlite3 in container disk.

    docker run -d -p 9000:9000 yujiod/wiki
    # same above
    docker run -d -p 9000:9000 -e DB_DRIVER=sqlite3 -e DB_SOURCE="./wiki.db" yujiod/wiki

Using MySQL

    docker run -d -p 9000:9000 -e DB_DRIVER=mysql -e DB_SOURCE="dbuser:dbpass@/dbname?charset=utf8" yujiod/wiki

Using PostgreSQL

    docker run -d -p 9000:9000 -e DB_DRIVER=postgres -e DB_SOURCE="user=dbuser dbname=dbpass sslmode=disable" yujiod/wiki

And open url http://localhost:9000/ from your browser.

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

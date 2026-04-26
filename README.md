# DEV.ENV

A simple management tool for application developoment services like databases, search engines etc.

## Content

## Introduction

When I am working on apps there are allways these same tools like MySQL, Redis, ElasticSearch and so on.
Since it feels useless to give every project its own docker compose file which spins up its private containers I looked for a better way to manage this.

There are some tools that come very close, but I do like a tool that does it in a way I would like to see it.

So that is where the root of this project started.

## Features

- Spawn a container once and use it in multiple projects through the the feature that the services provides to seperate it from other projects.
- Manage which services you use in your project by defining them in one configuration file, the `dev.env.yaml` file.
- Takes care of preparing the services to be used in your project.
- Prepares a environment file in your projects with a template based on the type of project you are working on.
- Powered by docker.

## Setup

### Prerequisits

- docker

### Instalation

TODO!

#### Linux

TODO!

#### MacOS

TODO!

#### Windows

TODO!

## CLI Reference

### init

Initialize the configuration of the project.

```
dev init
```

Note: With the autodetection of your project type it will sugest certain configurations.

### start

Starts up the services that are required for this project if they are not yet started.
After that it will make sure all the preperations are done and if not yet it will prompt you to generate the environment file adjecent to you project type.

```
dev start
```

Note: If you already have been working on a other project it will not spawn any new containers.

### stop

Stops the containers required for this project.

```
dev stop
```

Note: It will not stop containers from running if another project is still using them.
It will only stop the container when all projects using projects are stopped.

### self-update

Updates DEV.ENV to the most recent update.

```
dev self-update
```

## Project Configuration

This is what a `dev.env.yaml` could look like:

```
project: some-project-name
type: laravel
services:
  - image: mysql # will resolve to the mysql:8.0 docker image
    tag: 8.0
  - redis # will resolve to the redis:latest docker image
```

### Reference

TODO!

## DEV.ENV Configuration

Aside from all the projects their personal configuration there is also a location where all the project information is stored.
This is done in the `.dev.env/` dir in your home directory.

This directory contains the folowing:

- The docker-compose file where we manage all of our services.
- The directory of all of our projects and the status (running or not) of them.
- User specific preferences.
- tbd.

### Reference

TODO!

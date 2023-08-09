# Golang fixture data loader

## Goal

Create a database with the type structs defined for your entities in your project.
Verify that all the fields for your yaml entities match with local structs defined in your project.
Insert all the entities you described in yaml files

## Description

This package will load yaml file, verify if each entity match with local structs of your project and insert rows in your database.

The entities in yaml files are named fixture.

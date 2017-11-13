Drupal
======

The following directory is an example of bootstrapping a Drupal project with the following files:

* docker-compose.yml - Basic configuration with php, mysql, solr and mailhog
* m8s.yml - For running build steps
* Makefile - Workflow declarations

Requirements:

* Drupal is checked out into the "app" directory
* All mysql, solr and mailhog connection strings should use 127.0.0.1
* You use in conjunction with one of our Pipeline docs
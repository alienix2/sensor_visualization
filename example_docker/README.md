# Using the example's docker

This example is based on the image available in the repository of **mosquitto go auth** (<https://github.com/iegomez/mosquitto-go-auth>) as well as the official go docker image, mariadb image and my own repository to obtain sensor_info (<https://github.com/alienix2/sensor_info>)

## How to run

To run the example you must have docker installed on your machine as well as the current user configured to use docker.

There is an automated script *full_autostart.sh* that will automatically generated self signed certificates as well as clone the sensor_info repo and start the dockers for the web_app as well as the dockers for the omnisubscriber, publisher and subscriber.

Otherwise you can manually start the dockers with the following options:

- docker compose up (to start the web_app, database and the mosquitto broker)
- docker compose up omnisubscriber (to start the omnisubscriber)
- docker compose up publisher (to start the publisher)
- docker compose up subscriber (to start the subscriber)

*Note:* When doing things manually you should provide coherent mosquitto config as well as signed certificates in the /certificates folder. The script
*certificate_creation.sh* can be run by hand to generate the certificates.

*Note:* If you use the full_autostart.sh script, pls note that the download of packages used to build the go projects might take a while so you can check the *docker logs* to see if the build process finished before trying to access the web_app.

*Note:* The web_app will be available at the address: <https://localhost:45678>, the browser will likely complain about the self-signed certificate

## Default ports

- **MariaDB:** 12345 (external) 3306 (internal)
- **Mosquitto:** (over tls): 8883 (external) 8883 (internal)
- **Web_app:** 45678 (external) 4000 (internal)

## Users and pass

The default users are:

- <alice@example.com>
- <omnisub@example.com>
- <bob@example.com>
- <charlie@example.com>

All have the password: password

## Command messages

The messages to which the publisher will react are:

`{"command": "turn_on"}` and `{"command": "turn_off"}`

Additionally notes can be added in the form `{"notes": "some note"}`

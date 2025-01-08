# sensor_visualization

This project aims to provide a simple way to visualize sensor data produced using the data structure defined inside my other other project: <https://github.com/alienix2/sensor_info> and used in the **omnisubscriber** example.

The project can be run using the command:

`go run ./cmd/web/` in the main folder of the project. The default port is 4000, to access the web app locally go to the address: <http://localhost:4000>

There are some flags that can be used to configure the web app, to see their description run the project using the command: `go run ./cmd/web/ -h`

In the current state, the web_app won't start if the dsn isn't reachable as it is thought mainly as a way to visualize the sensor's data and therefore if the data isn't accessible, it's usefulness is very limited.

The tests can be run using the command:

`go test ./...` in the main folder of the project.

The web_app provides the possibility to authenticate the user following the rules defined in the go-auth-plugin: <https://github.com/iegomez/mosquitto-go-auth>

The web_app is intended to be used for **visualization and message sending only** for the moment, there is no possibility to add new users, topics and ACL using the web_app.

For the management of the database structure refer to the go-auth-plugin and the sensor_info project. The files inside the *example_docker/* folder can give you insight on how a minimum database configuration can be done to make the web app compliant with the go-auth-plugin.

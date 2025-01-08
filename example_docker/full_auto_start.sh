set -e

echo "Starting cration of the certifications"
./certificate_creation.sh

echo "Starting the docker containers for the web app, mqtt broker and mariadb"
docker-compose up -d

echo "Cloning the sensor_info repository"
if [ -d "sensor_info" ]; then
  echo "sensor_info already exists, deleting and cloning again"
  rm -rf sensor_info
fi
git clone https://github.com/alienix2/sensor_info.git

echo "Starting the docker containers for the subscriber, publisher and omnisubscriber"
docker-compose up -d subscriber publisher omnisubscriber

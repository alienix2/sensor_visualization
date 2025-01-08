#!/bin/bash

mkdir -p certifications

openssl genpkey -algorithm RSA -out certifications/ca.key -pkeyopt rsa_keygen_bits:2048
openssl req -x509 -new -nodes -key certifications/ca.key -sha256 -days 3650 -out certifications/ca.crt -subj "/C=US/ST=State/L=City/O=Organization/OU=OrgUnit/CN=CA"

create_cert() {
  local name=$1
  local subj=$2

  openssl genpkey -algorithm RSA -out certifications/$name.key -pkeyopt rsa_keygen_bits:2048

  openssl req -new -key certifications/$name.key -out certifications/$name.csr -subj "$subj"

  # Create config file for SAN
  cat >certifications/$name.cnf <<EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = req_ext
[req_distinguished_name]
[req_ext]
subjectAltName = @alt_names
[alt_names]
DNS.1 = mqtt_broker
EOF

  openssl x509 -req -in certifications/$name.csr -CA certifications/ca.crt -CAkey certifications/ca.key -CAcreateserial -out certifications/$name.crt -days 365 -sha256 -extfile certifications/$name.cnf -extensions req_ext
}

create_cert "publisher" "/C=US/ST=State/L=City/O=Organization/OU=OrgUnit/CN=publisher"
create_cert "subscriber" "/C=US/ST=State/L=City/O=Organization/OU=OrgUnit/CN=subscriber"
create_cert "mosquitto" "/C=US/ST=State/L=City/O=Organization/OU=OrgUnit/CN=mosquitto"

echo "Certificates and keys have been created in the 'certifications' directory."

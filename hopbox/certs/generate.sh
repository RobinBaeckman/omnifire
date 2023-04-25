# generate private CA key and certificate
openssl genrsa -out ca-key.pem 4096
openssl req -x509 -new -key ca-key.pem -out ca-cert.pem -days 3650 -config ca-cert.cnf

# generate traefik proxy cert
openssl genrsa -out proxy-key.pem 2048
openssl req -new -key proxy-key.pem -out proxy.csr -config proxy-cert.cnf
openssl x509 -req -in proxy.csr -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out proxy-cert.pem -days 365 -extensions v3_req -extfile proxy-cert.cnf

# generate client cert
openssl genrsa -out client-key.pem 2048
openssl req -new -key client-key.pem -out client.csr -config client-cert.cnf
openssl x509 -req -in client.csr -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -days 365 -extensions v3_req -extfile client-cert.cnf


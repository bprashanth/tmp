#!/bin/bash
mkdir ~/SSLCA/root/
cd ~/SSLCA/root/
openssl genrsa -aes256 -out rootca.key 2048
openssl req -sha256 -new -x509 -days 1826 -key rootca.key -out rootca.crt
touch certindex
echo 1000 > certserial
echo 1000 > crlnumber
echo '
[ ca ]
default_ca = myca

[ crl_ext ]
issuerAltName=issuer:copy
authorityKeyIdentifier=keyid:always

 [ myca ]
 dir = ./
 new_certs_dir = $dir
 unique_subject = no
 certificate = $dir/rootca.crt
 database = $dir/certindex
 private_key = $dir/rootca.key
 serial = $dir/certserial
 default_days = 730
 default_md = sha1
 policy = myca_policy
 x509_extensions = myca_extensions
 crlnumber = $dir/crlnumber
 default_crl_days = 730

 [ myca_policy ]
 commonName = supplied
 stateOrProvinceName = supplied
 countryName = optional
 emailAddress = optional
 organizationName = supplied
 organizationalUnitName = optional

 [ myca_extensions ]
 basicConstraints = critical,CA:TRUE
 keyUsage = critical,any
 subjectKeyIdentifier = hash
 authorityKeyIdentifier = keyid:always,issuer
 keyUsage = digitalSignature,keyEncipherment,cRLSign,keyCertSign
 extendedKeyUsage = serverAuth
 crlDistributionPoints = @crl_section
 subjectAltName  = @alt_names
 authorityInfoAccess = @ocsp_section

 [ v3_ca ]
 basicConstraints = critical,CA:TRUE,pathlen:0
 keyUsage = critical,any
 subjectKeyIdentifier = hash
 authorityKeyIdentifier = keyid:always,issuer
 keyUsage = digitalSignature,keyEncipherment,cRLSign,keyCertSign
 extendedKeyUsage = serverAuth
 crlDistributionPoints = @crl_section
 subjectAltName  = @alt_names
 authorityInfoAccess = @ocsp_section

 [alt_names]
 DNS.0 = Sparkling Intermidiate CA 1
 DNS.1 = Sparkling CA Intermidiate 1

 [crl_section]
 URI.0 = http://pki.sparklingca.com/SparklingRoot.crl
 URI.1 = http://pki.backup.com/SparklingRoot.crl

 [ocsp_section]
 caIssuers;URI.0 = http://pki.sparklingca.com/SparklingRoot.crt
 caIssuers;URI.1 = http://pki.backup.com/SparklingRoot.crt
 OCSP;URI.0 = http://pki.sparklingca.com/ocsp/
 OCSP;URI.1 = http://pki.backup.com/ocsp/' > ca.conf

openssl genrsa -out intermediate1.key 2048
openssl req -new -sha256 -key intermediate1.key -out intermediate1.csr
openssl ca -batch -config ca.conf -notext -in intermediate1.csr -out intermediate1.crt
mkdir ~/SSLCA/intermediate1/
cd ~/SSLCA/intermediate1/
cp ~/SSLCA/root/intermediate1.key ./
cp ~/SSLCA/root/intermediate1.crt ./
touch certindex
echo 1000 > certserial
echo 1000 > crlnumber
echo '
[ ca ]
default_ca = myca

[ crl_ext ]
issuerAltName=issuer:copy
authorityKeyIdentifier=keyid:always

 [ myca ]
 dir = ./
 new_certs_dir = $dir
 unique_subject = no
 certificate = $dir/intermediate1.crt
 database = $dir/certindex
 private_key = $dir/intermediate1.key
 serial = $dir/certserial
 default_days = 365
 default_md = sha1
 policy = myca_policy
 x509_extensions = myca_extensions
 crlnumber = $dir/crlnumber
 default_crl_days = 365

 [ myca_policy ]
 commonName = supplied
 stateOrProvinceName = supplied
 countryName = optional
 emailAddress = optional
 organizationName = supplied
 organizationalUnitName = optional

 [ myca_extensions ]
 basicConstraints = critical,CA:FALSE
 keyUsage = critical,any
 subjectKeyIdentifier = hash
 authorityKeyIdentifier = keyid:always,issuer
 keyUsage = digitalSignature,keyEncipherment
 extendedKeyUsage = serverAuth
 crlDistributionPoints = @crl_section
 subjectAltName  = @alt_names
 authorityInfoAccess = @ocsp_section

 [alt_names]
 DNS.0 = bitrot.in
 DNS.1 = bitrot.com

 [crl_section]
 URI.0 = http://pki.sparklingca.com/SparklingIntermidiate1.crl
 URI.1 = http://pki.backup.com/SparklingIntermidiate1.crl

 [ocsp_section]
 caIssuers;URI.0 = http://pki.sparklingca.com/SparklingIntermediate1.crt
 caIssuers;URI.1 = http://pki.backup.com/SparklingIntermediate1.crt
 OCSP;URI.0 = http://pki.sparklingca.com/ocsp/
 OCSP;URI.1 = http://pki.backup.com/ocsp/' > ca.conf

mkdir enduser-certs
openssl genrsa -out enduser-certs/enduser-example.com.key 2048
openssl req -new -sha256 -key enduser-certs/enduser-example.com.key -out enduser-certs/enduser-example.com.csr
openssl ca -batch -config ca.conf -notext -in enduser-certs/enduser-example.com.csr -out enduser-certs/enduser-example.com.crt

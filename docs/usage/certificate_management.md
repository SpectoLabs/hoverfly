# HTTPS support & certificate management
Hoverfly ships with a default certificate (`cert.pem` in the repository root directory). To use Hoverfly with HTTPS traffic, you will need to add this default certificate to your trust store. 

Using Hoverfly, you may also generate new certificates.  

## Generating and using certificates

To enable support for HTTPS services, Hoverfly can generate public and private keys. To generate a key pair, use the `-generate-ca-cert` flag:

    ./hoverfly -generate-ca-cert
    
This will create `cert.pem` and a `key.pem` files in your current directory. Next time you run Hoverfly, you can tell it to use these certificate and key files using the `-cert` and `-key` flags:

    ./hoverfly -cert cert.pem -key key.pem
    
You will then need to add the `cert.pem` file to your trusted certificates. Alternatively, you can turn off certificate verification. For example, to make insecure requests with cURL, you can use the `-k` flag:

    curl https://www.bbc.co.uk --proxy http://${HOVERFLY_HOST}:8500 -k

## Turn off verification when capturing or modifying traffic

You can tell Hoverfly to ignore untrusted certificates when capturing or modifying traffic in two ways.

1. Use the `-tls-verification=false` flag on startup:

       ./hoverfly -tls-verification=false
       
2. Set the `HoverflyTlsVerification` environment variable:

       export HoverflyTlsVerification=false      
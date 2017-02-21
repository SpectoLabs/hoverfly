.. _configuressl:

Configuring SSL in Hoverfly
===========================

In some cases, you may not wish to use Hoverfly's default SSL certificate. Hoverfly allows 
you to generate a new certificate and key.

The following command will start a Hoverfly process and create new ``cert.pem`` and ``key.pem``
files in the current working directory. These newly-created files will be loaded into the
running Hoverfly instance.

.. literalinclude:: hoverfly-generate-ca.sh
   :language: sh 

Optionally, you can provide a custom certificate name and authority:

.. literalinclude:: hoverfly-generate-ca-custom-name.sh
   :language: sh 

Once you have generated ``cert.pem`` and ``key.pem`` files with Hoverfly, you can use hoverctl 
to start an instance of Hoverfly using these files.

.. literalinclude:: hoverctl-start-with-custom-cert.sh
   :language: sh

.. note::
   Both a certificate and a key file must be supplied. The files must be in unencrypted PEM format.

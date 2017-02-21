.. _configuressl:

Configuring SSL in Hoverfly
===========================

When running Hoverfly, you may not want to use the default SSL certificate being used by Hoverfly.
You may want to geenrate a new certificate and key for your testing;

Generating a new certificate
----------------------------

.. literalinclude:: hoverfly-generate-ca.sh
   :language: sh 

**Cert name and cert org are optional**

.. literalinclude:: hoverctl-start-with-custom-cert.sh
   :language: sh

**Certificiate and key need to be provided together, cannot be supplied individually**

**Should be unecrypted PEM files**

# -*- mode: ruby -*-
# vi: set ft=ruby :

# This script downloads a pre-build hoverfly binary because building hoverfly
# in a vagrant guest means dealing with Go dependencies. Life is too short.

$bootstrapScript = <<SCRIPT
sudo apt-get update
sudo apt-get install -y redis-server
sudo mv /etc/redis/redis.conf /etc/redis/redis.conf.old
echo "bind 0.0.0.0" | sudo tee /etc/redis/redis.conf
cat /etc/redis/redis.conf.old | grep -v bind | sudo tee -a /etc/redis/redis.conf
sudo service redis-server restart
wget –quiet https://storage.googleapis.com/hoverfly-binaries/hoverfly_v0.3_linux_amd64
chmod +x hoverfly_v0.3_linux_amd64
ln -s hoverfly_v0.3_linux_amd64 hoverfly
SCRIPT

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "ubuntu/trusty64"
  config.vm.network "forwarded_port", guest: 8888, host: 8888
  config.vm.network "forwarded_port", guest: 8500, host: 8500
  config.vm.provision "shell", privileged: false, inline: $bootstrapScript
end

# -*- mode: ruby -*-
# vi: set ft=ruby :

# This script downloads a pre-built hoverfly binary because building hoverfly
# in a vagrant guest means dealing with Go dependencies. Life is too short.

$bootstrapScript = <<SCRIPT
sudo apt-get update
sudo apt-get install -y redis-server
sudo mv /etc/redis/redis.conf /etc/redis/redis.conf.old
echo "bind 0.0.0.0" | sudo tee /etc/redis/redis.conf
cat /etc/redis/redis.conf.old | grep -v bind | sudo tee -a /etc/redis/redis.conf
sudo service redis-server restart
wget https://storage.googleapis.com/hoverfly-binaries/hoverfly_v0.4_linux_amd64
mv hoverfly_v0.4_linux_amd64 hoverfly && chmod +x hoverfly
ln -s /vagrant/static /home/vagrant/static
nohup ./hoverfly > hoverfly.log 2>&1 & echo $! > hoverfly_pid &
SCRIPT

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "ubuntu/trusty64"
  config.vm.network "forwarded_port", guest: 8888, host: 8888
  config.vm.network "forwarded_port", guest: 8500, host: 8500
  config.vm.provision "shell", privileged: false, inline: $bootstrapScript
end

# -*- mode: ruby -*-
# vi: set ft=ruby :

$bootstrapScript = <<SCRIPT

# Install and configure Redis

sudo apt-get update
sudo apt-get install -y redis-server
sudo mv /etc/redis/redis.conf /etc/redis/redis.conf.old
echo "bind 0.0.0.0" | sudo tee /etc/redis/redis.conf
cat /etc/redis/redis.conf.old | grep -v bind | sudo tee -a /etc/redis/redis.conf
sudo service redis-server restart

# Download a pre-built hoverfly binary to avoid dealing with Go dependencies

wget https://storage.googleapis.com/hoverfly-binaries/latest/hoverfly_v0.4.2_linux_amd64
mv hoverfly_v0.4.2_linux_amd64 hoverfly && chmod +x hoverfly

# Symlink the webUI static files

ln -s /vagrant/static /home/vagrant/static

# Start hoverfly in the background in "virtualize" mode, pass logging output
# to hoverfly.log and save the hoverfly PID into a file so you can kill the process
# with 'kill -SIGTERM $(cat hoverfly_pid)'

nohup ./hoverfly > hoverfly.log 2>&1 & echo $! > hoverfly_pid &

SCRIPT

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "ubuntu/trusty64"
  config.vm.network "forwarded_port", guest: 8888, host: 8888
  config.vm.network "forwarded_port", guest: 8500, host: 8500
  config.vm.provision "shell", privileged: false, inline: $bootstrapScript
end

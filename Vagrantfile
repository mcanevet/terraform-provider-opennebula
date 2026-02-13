# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  # Ubuntu 22.04 LTS - use generic/ubuntu2204 for libvirt compatibility
  config.vm.box = "generic/ubuntu2204"

  config.vm.hostname = "opennebula-test"

  # Forward OpenNebula ports to host
  config.vm.network "forwarded_port", guest: 2633, host: 2633  # XML-RPC API
  config.vm.network "forwarded_port", guest: 2474, host: 2474  # Flow API
  config.vm.network "forwarded_port", guest: 9869, host: 9869  # Sunstone (optional)

  # Allocate resources for OpenNebula + KVM
  config.vm.provider "libvirt" do |libvirt|
    libvirt.memory = 4096
    libvirt.cpus = 2
    libvirt.nested = true  # Enable nested virtualization
  end

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "4096"
    vb.cpus = 2
    vb.name = "opennebula-test"
    vb.customize ["modifyvm", :id, "--nested-hw-virt", "on"]
  end

  # Provision OpenNebula using miniONE
  config.vm.provision "shell", inline: <<-SHELL
    set -e

    echo "📦 Installing prerequisites..."
    apt-get update
    apt-get install -y wget curl

    echo "📥 Downloading miniONE..."
    wget -q https://github.com/OpenNebula/minione/releases/latest/download/minione -O /tmp/minione
    chmod +x /tmp/minione

    echo "🚀 Installing OpenNebula 6.10 (this takes ~5 minutes)..."
    bash /tmp/minione --version 6.10 --password opennebula --yes

    echo "✅ OpenNebula installed successfully!"
    echo ""
    echo "Credentials:"
    echo "  Username: oneadmin"
    echo "  Password: opennebula"
    echo ""
    echo "Endpoints from host machine:"
    echo "  API:   http://localhost:2633/RPC2"
    echo "  Flow:  http://localhost:2474"
    echo "  UI:    http://localhost:9869"
  SHELL

  # Display connection info on vagrant up
  config.vm.post_up_message = <<-MSG
    ✅ OpenNebula VM is ready!

    Test connection from host:
      export OPENNEBULA_ENDPOINT="http://localhost:2633/RPC2"
      export OPENNEBULA_USERNAME="oneadmin"
      export OPENNEBULA_PASSWORD="opennebula"
      export OPENNEBULA_FLOW_ENDPOINT="http://localhost:2474"
      export TF_ACC=1

      # Run migration tests
      go test ./opennebula -v -tags=acceptance -run TestMigrate -timeout 30m

    Manage VM:
      vagrant ssh      # Access VM
      vagrant halt     # Stop VM
      vagrant destroy  # Delete VM
  MSG
end

Vagrant.configure("2") do |config|
  config.vm.box = "generic/alpine38"
  config.vm.provision "file", source: "snowcapper", destination: "snowcapper"
  config.vm.provision "shell", inline: "
  apk update &&
  apk add ca-certificates &&
  sudo ./snowcapper
  "
  config.vbguest.auto_update = false
end

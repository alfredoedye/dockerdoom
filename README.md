# dockerdoom
Tired of having to manually stop docker containers? 

Your suffering has ended. You can now have fun while killing all the containers.




# How to use

docker run --rm=true -p 5900:5900 -v /var/run/docker.sock:/var/run/docker.sock --name=dockerdoom doomc




#build the go app
. export GOOS=linux
. export GOARCH=386
. go build -name dockerdoom .


#build the docker image
docker build -t dockerdoom .




# Doom Cheat Codes. 

http://doom.wikia.com/wiki/Doom_cheat_codes


idspispopd

idkfa




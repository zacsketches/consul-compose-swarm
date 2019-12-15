# This dockerfile expects to find the built binary in a local /bin 
# directory accessible relative to thd Dockerfile.  Since this project
# does not include CI then this bin folder should be created locally
#      $ mkdir bin
#
# Then the golang executable should be built and installed into the bin
# directory.  Note that we build for linux regardless of the host operating
# system because the container is linux 
#      $ env GOOS=linux GOARCH=amd64 go build -o bin/webserver
#
# To build the container use this command line foo
#      $ docker build -t consul-example-webserver .
#
# To run the container use this foo which will keep the container from 
# persisting on the local machine
#      $ docker run -it --rm -p 8080:8080 zacsketches/consul-example-webserver
#
# And finally, to push the container login and tag then push
#      $ docker tag consul-example-webserver <username>/consul-example-webserver
#      $ docker push <username>/consul-example-webserver

FROM golang
COPY bin/webserver $GOPATH/bin/webserver

EXPOSE 8080

CMD [ "webserver" ]
# CMD [ "/bin/sh" ]


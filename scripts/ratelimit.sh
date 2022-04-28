hey -n 10000 http://whoami.docker.localhost/

docker run --rm  williamyeh/hey  \
      -n 1000 -c 100 http://whoami.docker.localhost/
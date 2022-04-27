curl --request POST \
  --url http://next.docker.localhost/api/create \
  --header 'Content-Type: application/json' \
  --data '{
	"title":"Naan"
}'


hey -n 100 -m POST -H "Content-Type: application/json" -d '{"title":"Pommes"}' http://next.docker.localhost/api/create
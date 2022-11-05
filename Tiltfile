# For more on Extensions, see: https://docs.tilt.dev/extensions.html
load('ext://restart_process', 'docker_build_with_restart')
load('ext://helm_remote', 'helm_remote')

compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api ./cmd/main.go'

# Compile locally and push to docker
local_resource(
  'go-compile',
  compile_cmd,
  deps=['./cmd/main.go', './'],
  ignore=['./build'],
  auto_init=False,
  labels=['api']
)

docker_build_with_restart(
  'social-graph-api-image',
  '.',
  entrypoint=['/app/build/api'],
  dockerfile='./Dockerfile.tilt',
  only=[
    './build',
    './',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./', '/app/'),
  ],
)

k8s_yaml('./infra/service/k8s.yml')
k8s_resource('social-graph-api', 
  auto_init=False,
  port_forwards=['3010:3010', '3011:3011'],
  resource_deps=['go-compile'],
  links=['neo4j-standalone'],
  labels=['api']
)

# Neo4j using Helm
helm_remote(
  'neo4j-standalone',
  repo_name='neo4j',
  repo_url='https://helm.neo4j.com/neo4j/',
  values=['./infra/database/values.yaml']
)

k8s_resource(
    'neo4j-standalone',
    auto_init=False,
    port_forwards=['7474:7474'],
    labels=['neo4j']
)
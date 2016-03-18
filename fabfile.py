from fabric.api import run, local, env, put
from fabric.context_managers import cd

def prod():
	env.hosts = ['soccerlc.com']
	env.user = 'root'

def build_server():
	local('docker build -t server .')
	local('docker run -it -v "$PWD":/usr/src/myapp -w /usr/src/myapp server go build')
	local('tar -cvf server.tar -C . .')

def deploy():
	build_server()
	run("service soccerlc stop")
	run("rm -r /srv/*")
	put("server.tar","./work/src/soccerlcprod/server.tar")
	run('tar -xf /root/work/src/soccerlcprod/server.tar -C /srv')
	run("service soccerlc restart")

from fabric.api import run, local, env, put
from fabric.context_managers import cd

def prod():
	env.hosts = ['soccerlc.com']
	env.user = 'root'

def build_server():
	#local('docker build -t server .')
	#local('docker run -it -v "$PWD":/usr/src/myapp -w /usr/src/myapp server go build')
	local('tar -cvf server.tar ../soccerlc')

def deploy():
	build_server()
	put("server.tar","./work/src/jess/joel.tar")
	run('tar -xf /root/work/src/jess/joel.tar -C /root/work/src/jess')
	#run('killall /root/work/bin/soccerlc')
	with cd('/root/work/src/jess/soccerlc'):
		run('go install .; nohup $GOPATH/bin/soccerlc')

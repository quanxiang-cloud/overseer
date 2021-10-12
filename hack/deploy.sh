#! /bin/bash
IMG=controller:latest
myReposity="false"
cd ../
ABS_DIR="$PWD"
cd hack/

while true
do
  case $1 in
  --my-reposity)
    myReposity=$2
    shift 2
    ;;
    --)
      shift 1
      break
      ;;
  esac
done

#build

docker build -t myReposity/${IMG} -f ../ .
docker push myReposity${IMG}

##install
cd $ABS_DIR/bin/
curl -s "https://raw.githubusercontent.com/\
kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"  | bash
cd ../hack
../bin/kustomize build ../config/crd | kubectl apply -f -
 #deploy
cd ../config/manager && ../../bin/kustomize edit set image controller=myReposity/${IMG}
../../bin/kustomize build ../../config/default | kubectl apply -f -

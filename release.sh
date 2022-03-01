#! /usr/bin/env bash
set -e

if [ -z "$1" ]; then
  echo "No version supplied"
  exit 1
fi

version=$1

if ! [[ $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  >&2 echo "Invalid Version: $version"
  exit 1
fi


for action in `find . -type f -name action.yaml`; do
  action_name=$(basename $(dirname $action))
  sed -i -E -e "s/${action_name}:[0-9]+\.[0-9]+\.[0-9]+/${action_name}:${version}/g" $action
  git add $action
done

git commit -m "build(release): v${version}"
git tag -m "v${version}" v${version}

major_version=`echo $version | cut -d . -f 1`
git tag -f -m "v${major_version}" v${major_version}

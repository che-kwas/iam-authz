#!/usr/bin/bash

# iam-auth smoke testing

APISERVER_ADDR=localhost:8000
APISERVER_ADDR=localhost:8010

Header="-HContent-Type: application/json"
CCURL="curl -s -XPOST" # Create
UCURL="curl -s -XPUT" # Update
RCURL="curl -s -XGET" # Retrieve
DCURL="curl -s -XDELETE" # Delete

test::admin_login()
{
  basicToken="-HAuthorization: Basic Y2hlOmNoZS1rd2FzLmdpdGVlLmlv"
  ${CCURL} "${basicToken}" http://${APISERVER_ADDR}/login | grep -Po '(?<=token":")(.+)(?=")'
}

test::tom_login()
{
  basicToken="-HAuthorization: Basic dG9tOnRvbXRvbQ=="
  ${CCURL} "${basicToken}" http://${APISERVER_ADDR}/login | grep -Po '(?<=token":")(.+)(?=")'
}

test::user()
{
  echo -e '\033[32m/v1/user test begin========\033[0m'

  token="-HAuthorization: Bearer $(test::admin_login)"

  # 1. 如果有tom、jerry、john用户先清空
  echo -e '\033[32m1. delete users\033[0m'
  ${DCURL} "${token}" http://${APISERVER_ADDR}/v1/users/tom; echo
  ${DCURL} "${token}" http://${APISERVER_ADDR}/v1/users/jerry; echo
  ${DCURL} "${token}" http://${APISERVER_ADDR}/v1/users/john; echo

  # 2. 创建tom、jerry、john用户
  echo -e '\033[32m2. create users\033[0m'
  ${CCURL} "${Header}" "${token}" http://${APISERVER_ADDR}/v1/users \
    -d'{"password":"tomtom","username":"tom","email":"tom@gmail.com","phone":"1812884xxxx"}'; echo
  ${CCURL} "${Header}" "${token}" http://${APISERVER_ADDR}/v1/users \
    -d'{"password":"jerryjerry","username":"jerry","email":"jerry@gmail.com","phone":"1812884xxxx"}'; echo
  ${CCURL} "${Header}" "${token}" http://${APISERVER_ADDR}/v1/users \
    -d'{"password":"johnjohn","username":"john","email":"john@gmail.com","phone":"1812884xxxx"}'; echo

  # tomToken="-HAuthorization: Bearer $(test::tom_login)"

  # # 3. 列出所有用户
  # echo -e '\033[32m3.1 tom cannot list users\033[0m'
  # ${RCURL} "${tomToken}" "http://${APISERVER_ADDR}/v1/users?offset=0&limit=10"; echo
  # echo -e '\033[32m3.2 admin list users\033[0m'
  # ${RCURL} "${token}" "http://${APISERVER_ADDR}/v1/users?offset=0&limit=10"; echo

  # # 4. 获取tom用户的详细信息
  # echo -e '\033[32m4.1 get user without login\033[0m'
  # ${RCURL} http://${APISERVER_ADDR}/v1/users/tom; echo
  # echo -e '\033[32m4.2 tom get user\033[0m'
  # ${RCURL} "${tomToken}" http://${APISERVER_ADDR}/v1/users/tom; echo

  # # 5. 修改tom用户
  # echo -e '\033[32m5. update tom\033[0m'
  # ${UCURL} "${Header}" "${token}" http://${APISERVER_ADDR}/v1/users/tom \
  #   -d'{"username":"tom","email":"tom_modified@gmail.com","phone":"1812884xxxx"}'; echo

  # # 6. 删除tom用户
  # echo -e '\033[32m6. delete tom\033[0m'
  # ${DCURL} "${token}" http://${APISERVER_ADDR}/v1/users/tom; echo

  # # 7. 批量删除用户
  # echo -e '\033[32m7. delete users\033[0m'
  # ${DCURL} "${token}" "http://${APISERVER_ADDR}/v1/users?name=jerry&name=john"; echo

  echo -e '\033[32m/v1/user test end==========\033[0m'
}

test::secret()
{
  echo -e '\033[32m/v1/secret test begin========\033[0m'

  token="-HAuthorization: Bearer $(test::admin_login)"

  # 1. 如果有secret0密钥先清空
  echo -e '\033[32m1. delete secret\033[0m'
  ${DCURL} "${token}" http://${APISERVER_ADDR}/v1/secrets/secret0; echo

  # 2. 创建secret0密钥
  echo -e '\033[32m2. create secret\033[0m'
  ${CCURL} "${Header}" "${token}" http://${APISERVER_ADDR}/v1/secrets \
    -d'{"metadata":{"name":"secret0"},"expires":0,"description":"admin secret"}'; echo

  # # 3. 列出所有密钥
  # echo -e '\033[32m3. list secrets\033[0m'
  # ${RCURL} "${token}" http://${APISERVER_ADDR}/v1/secrets; echo

  # # 4. 获取secret0密钥的详细信息
  # echo -e '\033[32m4. get secret\033[0m'
  # ${RCURL} "${token}" http://${APISERVER_ADDR}/v1/secrets/secret0; echo

  # # 5. 修改secret0密钥
  # echo -e '\033[32m5. update secret\033[0m'
  # ${UCURL} "${Header}" "${token}" http://${APISERVER_ADDR}/v1/secrets/secret0 \
  #   -d'{"expires":0,"description":"admin secret(modified)"}'; echo

  # # 6. 删除secret0密钥
  # echo -e '\033[32m6. delete secret\033[0m'
  # ${DCURL} "${token}" http://${APISERVER_ADDR}/v1/secrets/secret0; echo

  echo -e '\033[32m/v1/secret test end==========\033[0m'
}

test::policy()
{
  echo -e '\033[32m/v1/policy test begin========\033[0m'

  token="-HAuthorization: Bearer $(test::admin_login)"

  # 1. 如果有policy0策略先清空
  echo -e '\033[32m1. delete policy\033[0m'
  ${DCURL} "${token}" http://${APISERVER_ADDR}/v1/policies/policy0; echo

  # 2. 创建policy0策略
  echo -e '\033[32m2. create policy\033[0m'
  ${CCURL} "${Header}" "${token}" http://${APISERVER_ADDR}/v1/policies \
    -d'{"metadata":{"name":"policy0"},"policy":{"description":"One policy to rule them all.","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}'; echo

  # # 3. 列出所有策略
  # echo -e '\033[32m3. list policies\033[0m'
  # ${RCURL} "${token}" http://${APISERVER_ADDR}/v1/policies; echo

  # # 4. 获取policy0策略的详细信息
  # echo -e '\033[32m4. get policy\033[0m'
  # ${RCURL} "${token}" http://${APISERVER_ADDR}/v1/policies/policy0; echo

  # # 5. 修改policy0策略
  # echo -e '\033[32m5. update policy\033[0m'
  # ${UCURL} "${Header}" "${token}" http://${APISERVER_ADDR}/v1/policies/policy0 \
  #   -d'{"policy":{"description":"One policy to rule them all(modified).","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}'; echo

  # # 6. 删除policy0策略
  # echo -e '\033[32m6. delete policy\033[0m'
  # ${DCURL} "${token}" http://${APISERVER_ADDR}/v1/policies/policy0; echo

  echo -e '\033[32m/v1/policy test end==========\033[0m'
}

test::refresh_logout()
{
  echo -e '\033[32mrefresh_logout begin========\033[0m'

  token="-HAuthorization: Bearer $(test::admin_login)"

  echo -e '\033[32m1. refresh token\033[0m'
  ${CCURL} "${token}" http://${APISERVER_ADDR}/refresh; echo
  echo -e '\033[32m2. logout \033[0m'
  ${CCURL} "${token}" http://${APISERVER_ADDR}/logout; echo

  echo -e '\033[32mrefresh_logout end==========\033[0m'
}

test::create_admin()
{
  ${CCURL} "${Header}" http://${APISERVER_ADDR}/v1/users \
    -d'{"password":"che-kwas.gitee.io","username":"che","email":"che@kwas.com","phone":"17700001111","isAdmin":true}'; echo
}

test::delete_admin()
{
  token="-HAuthorization: Bearer $(test::admin_login)"
  ${DCURL} "${token}" http://${APISERVER_ADDR}/v1/users/che; echo
}

# test::delete_admin
# test::create_admin
test::user
test::secret
test::policy
# test::refresh_logout

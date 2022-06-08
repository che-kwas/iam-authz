#!/usr/bin/bash

# iam-authz smoke testing

APISERVER_ADDR=localhost:8000
AUTHZSERVER_ADDR=localhost:8010

Header="-HContent-Type: application/json"
CCURL="curl -s -XPOST" # Create
UCURL="curl -s -XPUT" # Update
RCURL="curl -s -XGET" # Retrieve
DCURL="curl -s -XDELETE" # Delete

test::admin_login()
{
  basicToken="-HAuthorization: Basic Y2hlOmNoZS1rd2FzLmdpdGVlLmlv"
  ${CCURL} "${basicToken}" http://${APISERVER_ADDR}/login | grep -Po 'token[" :]+\K[^"]+'
}

test::authz()
{
  echo -e '\033[32m/v1/auth test begin========\033[0m'

  token="-HAuthorization: Bearer $(test::admin_login)"

  # 1. 如果有 policies 先清空
  echo -e '\033[32m1. delete policies\033[0m'
  ${DCURL} "${token}" http://${APISERVER_ADDR}/v1/policies/name=authzpolicy1&name=authzpolicy2; echo

  # 2. 创建 policies
  echo -e '\033[32m2. create policies\033[0m'
  ${CCURL} "${Header}" "${token}" http://${APISERVER_ADDR}/v1/policies \
    -d'{"metadata":{"name":"authzpolicy1"},"policy":{"description":"One policy to rule them all.","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}'; echo
  ${CCURL} "${Header}" "${token}" http://${APISERVER_ADDR}/v1/policies \
    -d'{"metadata":{"name":"authzpolicy2"},"policy":{"description":"Deny peter delete","subjects":["users:peter"],"actions":["delete"],"effect":"deny","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}'; echo

  # 3. 如果有 authzsecret 先清空
  echo -e '\033[32m3. delete authzsecret\033[0m'
  ${DCURL} "${token}" http://${APISERVER_ADDR}/v1/secrets/authzsecret; echo

  # 4. 创建 authzsecret 策略
  echo -e '\033[32m4. create authzsecret\033[0m'
  ${CCURL} "${Header}" "${token}" http://${APISERVER_ADDR}/v1/secrets -d'{"metadata":{"name":"authzsecret"},"expires":0,"description":"admin secret"}'; echo

  # 5. 生成由 authzsecret 签名的token
  echo -e '\033[32m5. get authzsecret token\033[0m'
  authzToken=$(${RCURL} "${Header}" "${token}" http://${APISERVER_ADDR}/v1/secrets/authzsecret/token | grep -Po 'token[" :]+\K[^"]+')
  echo "authzToken=${authzToken}"

  # 注意这里要sleep 2s 等待 iam-authz 将新建的密钥同步到其内存中
  sleep 2

  # 6. 测试授权通过
  echo -e '\033[32m6. authz granted\033[0m'
  $CCURL "${Header}" -H"Authorization: Bearer ${authzToken}" http://${AUTHZSERVER_ADDR}/v1/authz \
    -d'{"subject":"users:maria","action":"delete","resource":"resources:articles:ladon-introduction","context":{"remoteIPAddress":"192.168.0.5"}}'; echo

  # 7. 测试授权未通过 - 禁止peter delete
  echo -e '\033[32m7. authz rejected by authzpolicy2\033[0m'
  $CCURL "${Header}" -H"Authorization: Bearer ${authzToken}" http://${AUTHZSERVER_ADDR}/v1/authz \
    -d'{"subject":"users:peter","action":"delete","resource":"resources:articles:ladon-introduction","context":{"remoteIPAddress":"192.168.0.5"}}'; echo

  # 8. 测试授权未通过 - 没有对应的策略（CIDR不符合condition）
  echo -e '\033[32m8. authz rejected by no policy\033[0m'
  $CCURL "${Header}" -H"Authorization: Bearer ${authzToken}" http://${AUTHZSERVER_ADDR}/v1/authz \
    -d'{"subject":"users:maria","action":"delete","resource":"resources:articles:ladon-introduction","context":{"remoteIPAddress":"10.17.0.5"}}'; echo

  echo -e '\033[32m/v1/auth test end==========\033[0m'
}

test::authz

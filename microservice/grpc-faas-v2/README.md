features:

1. file transfer
2. build and deploy container through faas-cli with transfered data
3. remove functions 
4. test functions 

service:

1. upload
2. deploy (upload -> build -> push -> deploy)
3. remove 
4. trigger (invoke)
5. schedule 

concerns:

1. user seperation
  {
    I. users have each unique file path
    II. /tmp/faasCode/{USER_NAMESPACE}/{FUNCTION_NAME.tar.gz}
  } 